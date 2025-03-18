// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"fmt"
	"io/fs"
	"mime"
	"path/filepath"
	"strings"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"

	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Deploy a cloud cdn entrypoint
func (a *NitricGcpPulumiProvider) deployEntrypoint(ctx *pulumi.Context) error {
	pathRules := compute.URLMapPathMatcherPathRuleArray{}

	// Add deployed API gatewayss to the URLMap
	for name, api := range a.ApiGateways {
		neg, err := compute.NewRegionNetworkEndpointGroup(ctx, fmt.Sprintf("%s-apigw-neg", name), &compute.RegionNetworkEndpointGroupArgs{
			NetworkEndpointType: pulumi.String("SERVERLESS"),
			Region:              api.Region,
			ServerlessDeployment: compute.RegionNetworkEndpointGroupServerlessDeploymentArgs{
				Platform: pulumi.String("apigateway.googleapis.com"),
				Resource: api.GatewayId,
			},
		})
		if err != nil {
			return err
		}

		bs, err := compute.NewBackendService(ctx, fmt.Sprintf("%s-apigw-bs", name), &compute.BackendServiceArgs{
			Backends: compute.BackendServiceBackendArray{
				compute.BackendServiceBackendArgs{
					Group: neg.SelfLink,
				},
			},
			// EnableCdn: pulumi.Bool(true),
			Protocol: pulumi.String("HTTPS"),
		}, nil)
		if err != nil {
			return err
		}

		pr := compute.URLMapPathMatcherPathRuleArgs{
			Service:     bs.ID(),
			Paths:       pulumi.StringArray{pulumi.Sprintf("/apis/%s/*", name)},
			RouteAction: compute.URLMapPathMatcherPathRuleRouteActionArgs{},
		}

		pathRules = append(pathRules, pr)
	}
	var defaultService pulumi.StringOutput
	for name, website := range a.WebsiteBuckets {
		normalizedName := strings.Replace(name, "/", "", -1)
		if normalizedName == "" {
			normalizedName = "default"
		}

		backend, err := compute.NewBackendBucket(ctx, fmt.Sprintf("%s-site-bucket", normalizedName), &compute.BackendBucketArgs{
			BucketName: website.Name,
			EnableCdn:  pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		if name == "/" {
			defaultService = backend.SelfLink

			// pr := compute.URLMapPathMatcherPathRuleArgs{
			// 	Service: backend.ID(),
			// 	Paths:   pulumi.StringArray{pulumi.String("./*")},
			// }

			// pathRules = append(pathRules, pr)
		} else {
			pr := compute.URLMapPathMatcherPathRuleArgs{
				Service: backend.ID(),
				Paths:   pulumi.StringArray{pulumi.String(filepath.Join("/", name, "./*"))},
			}

			pathRules = append(pathRules, pr)
		}
	}

	// Provision a global IP address for the CDN.
	ip, err := compute.NewGlobalAddress(ctx, "ip", nil)
	if err != nil {
		return err
	}

	// Create a URLMap to route requests to the storage bucket.
	urlMap, err := compute.NewURLMap(ctx, "url-map", &compute.URLMapArgs{
		DefaultService: defaultService,
		HostRules: compute.URLMapHostRuleArray{
			compute.URLMapHostRuleArgs{
				Hosts:       pulumi.StringArray{pulumi.String("*")},
				PathMatcher: pulumi.String("all-paths"),
			},
		},
		PathMatchers: compute.URLMapPathMatcherArray{
			compute.URLMapPathMatcherArgs{
				Name:           pulumi.String("all-paths"),
				DefaultService: defaultService,
				PathRules:      pathRules,
			},
		},
	})
	if err != nil {
		return err
	}

	// Create an HTTP proxy to route requests to the URLMap.
	httpProxy, err := compute.NewTargetHttpProxy(ctx, "http-proxy", &compute.TargetHttpProxyArgs{
		UrlMap: urlMap.SelfLink,
	})
	if err != nil {
		return err
	}

	// Create a GlobalForwardingRule rule to route requests to the HTTP proxy.
	_, err = compute.NewGlobalForwardingRule(ctx, "http-forwarding-rule", &compute.GlobalForwardingRuleArgs{
		IpAddress:  ip.Address,
		IpProtocol: pulumi.String("TCP"),
		PortRange:  pulumi.String("80"),
		Target:     httpProxy.SelfLink,
	})
	if err != nil {
		return err
	}

	return nil
}

// Website - Implements the Website deployment method for the GCP provider
func (a *NitricGcpPulumiProvider) Website(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Website) error {
	var err error

	a.WebsiteBuckets[config.BasePath], err = storage.NewBucket(ctx, "websites", &storage.BucketArgs{
		Location: pulumi.String(a.Region),
		Website: &storage.BucketWebsiteArgs{
			MainPageSuffix: pulumi.String(config.IndexDocument),
			NotFoundPage:   pulumi.String(config.ErrorDocument),
		},
	})
	if err != nil {
		return err
	}

	_, err = storage.NewBucketIAMBinding(ctx, "bucket-iam-binding", &storage.BucketIAMBindingArgs{
		Bucket: a.WebsiteBuckets[config.BasePath].Name,
		Role:   pulumi.String("roles/storage.objectViewer"),
		Members: pulumi.StringArray{
			pulumi.String("allUsers"),
		},
	})
	if err != nil {
		return err
	}

	err = filepath.WalkDir(config.GetLocalDirectory(), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Get file info to check for special types
		info, err := d.Info()
		if err != nil {
			return err
		}

		// Skip non-regular files (e.g., symlinks, sockets, devices)
		if info.Mode()&fs.ModeType != 0 {
			return nil
		}

		// Determine the content type based on the file extension
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		relativePath, err := filepath.Rel(config.GetLocalDirectory(), path)
		if err != nil {
			return err
		}

		storage.NewBucketObject(ctx, path, &storage.BucketObjectArgs{
			Bucket:      a.WebsiteBuckets[config.BasePath].Name,
			Name:        pulumi.String(relativePath),
			Source:      pulumi.NewFileAsset(path),
			ContentType: pulumi.String(contentType),
		})
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
