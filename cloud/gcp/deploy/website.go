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
	"path"
	"path/filepath"
	"sort"
	"strings"

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/certificatemanager"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/dns"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ensureValidSubdomain(domain string, subdomain string) error {
	domain = strings.ToLower(strings.TrimSuffix(domain, "."))
	subdomain = strings.ToLower(strings.TrimSuffix(subdomain, "."))

	if subdomain == domain || strings.HasSuffix(subdomain, "."+domain) {
		return nil
	}

	return fmt.Errorf("%s is not a valid subdomain of %s", subdomain, domain)
}

// Deploy a cloud CDN entrypoint
func (a *NitricGcpPulumiProvider) deployEntrypoint(ctx *pulumi.Context) error {
	if a.GcpConfig.CdnDomain.ZoneName == "" {
		return fmt.Errorf("a valid DNS zone is required to deploy websites to GCP")
	}

	if a.GcpConfig.CdnDomain.DomainName == "" {
		return fmt.Errorf("a valid domain name is required to deploy websites to GCP")
	}

	pathRules := compute.URLMapPathMatcherPathRuleArray{}

	// Add deployed API gateways to the URLMap
	for apiName, api := range a.ApiGateways {
		neg, err := compute.NewRegionNetworkEndpointGroup(ctx, fmt.Sprintf("%s-apigw-neg", apiName), &compute.RegionNetworkEndpointGroupArgs{
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

		bs, err := compute.NewBackendService(ctx, fmt.Sprintf("%s-apigw-bs", apiName), &compute.BackendServiceArgs{
			Backends: compute.BackendServiceBackendArray{
				compute.BackendServiceBackendArgs{
					Group: neg.SelfLink,
				},
			},
			Protocol: pulumi.String("HTTPS"),
		}, nil)
		if err != nil {
			return err
		}

		pr := compute.URLMapPathMatcherPathRuleArgs{
			Service: bs.ID(),
			Paths:   pulumi.StringArray{pulumi.Sprintf("/api/%s/*", apiName)},
			RouteAction: compute.URLMapPathMatcherPathRuleRouteActionArgs{
				UrlRewrite: compute.URLMapPathMatcherPathRuleRouteActionUrlRewriteArgs{
					PathPrefixRewrite: pulumi.String("/"),
					HostRewrite:       api.DefaultHostname,
				},
			},
		}

		pathRules = append(pathRules, pr)
	}

	cdnPolicy := compute.BackendBucketCdnPolicyArgs{}

	if a.GcpConfig.CdnDomain.ClientTtl != nil {
		cdnPolicy.ClientTtl = pulumi.Int(*a.GcpConfig.CdnDomain.ClientTtl)
	}

	if a.GcpConfig.CdnDomain.DefaultTtl != nil {
		cdnPolicy.DefaultTtl = pulumi.Int(*a.GcpConfig.CdnDomain.DefaultTtl)
	}

	var defaultService pulumi.StringOutput
	for sitePath, siteBucket := range a.WebsiteBuckets {
		normalizedName := strings.Replace(sitePath, "/", "", -1)
		if normalizedName == "" {
			normalizedName = "default"
		}

		backend, err := compute.NewBackendBucket(ctx, fmt.Sprintf("%s-site-bucket", normalizedName), &compute.BackendBucketArgs{
			BucketName:      siteBucket.Name,
			EnableCdn:       pulumi.Bool(true),
			CompressionMode: pulumi.String("AUTOMATIC"),
			CdnPolicy:       cdnPolicy,
		})
		if err != nil {
			return err
		}

		if sitePath == "/" {
			defaultService = backend.SelfLink
		} else {
			pr := compute.URLMapPathMatcherPathRuleArgs{
				Service: backend.ID(),
				Paths: pulumi.StringArray{
					pulumi.String(path.Join("/", sitePath)),
					pulumi.String(path.Join("/", sitePath, "./*")),
				},
				RouteAction: compute.URLMapPathMatcherPathRuleRouteActionArgs{
					UrlRewrite: compute.URLMapPathMatcherPathRuleRouteActionUrlRewriteArgs{
						PathPrefixRewrite: pulumi.String("/"),
					},
				},
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
	httpsUrlMap, err := compute.NewURLMap(ctx, "https-site-url-map", &compute.URLMapArgs{
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

	// If a domain is specified in the config, then lookup to see if there is a GCP managed zone for it
	managedZone, err := dns.LookupManagedZone(ctx, &dns.LookupManagedZoneArgs{
		Name: a.GcpConfig.CdnDomain.ZoneName,
	})
	if err != nil {
		return err
	}

	// Add root zone, to ensure reliable comparisons (i.e. trailing dot, e.g. example.com.)
	subDomain := strings.ToLower(a.GcpConfig.CdnDomain.DomainName)
	if !strings.HasSuffix(subDomain, ".") {
		subDomain = subDomain + "."
	}

	err = ensureValidSubdomain(managedZone.DnsName, subDomain)
	if err != nil {
		return fmt.Errorf("CDN domain '%s' is not valid for zone '%s': %w", strings.TrimSuffix(subDomain, "."), managedZone.Name, err)
	}

	// Create root DNS record for the IP address
	_, err = dns.NewRecordSet(ctx, "cdn-dns-record", &dns.RecordSetArgs{
		Name:        pulumi.String(subDomain),
		ManagedZone: pulumi.String(managedZone.Name),
		Type:        pulumi.String("A"),
		Rrdatas:     pulumi.StringArray{ip.Address},
		Ttl:         pulumi.IntPtr(300),
	})
	if err != nil {
		return err
	}

	_, err = dns.NewRecordSet(ctx, "www-cdn-dns-record", &dns.RecordSetArgs{
		Name:        pulumi.String(fmt.Sprintf("www.%s", subDomain)),
		ManagedZone: pulumi.String(managedZone.Name),
		Type:        pulumi.String("A"),
		Rrdatas:     pulumi.StringArray{ip.Address},
		Ttl:         pulumi.IntPtr(300),
	})
	if err != nil {
		return err
	}

	// The certificate will use Load Balancer authorization (as opposed to DNS auth).
	sslCert, err := certificatemanager.NewCertificate(ctx, "cdn-cert", &certificatemanager.CertificateArgs{
		Scope: pulumi.String("DEFAULT"),
		Managed: certificatemanager.CertificateManagedArgs{
			Domains: pulumi.StringArray{
				// Removing trailing dot (root zone), it's unsupported by certificate manager
				pulumi.String(strings.TrimSuffix(subDomain, ".")),
			},
		},
	})
	if err != nil {
		return err
	}

	certMap, err := certificatemanager.NewCertificateMapResource(ctx, "cert-map", &certificatemanager.CertificateMapResourceArgs{})
	if err != nil {
		return err
	}

	_, err = certificatemanager.NewCertificateMapEntry(ctx, "default", &certificatemanager.CertificateMapEntryArgs{
		Name:        pulumi.String("cdn-cert-map-entry"),
		Description: pulumi.String("CDN Certificate Map Entry"),
		Map:         certMap.Name,
		Certificates: pulumi.StringArray{
			sslCert.ID(),
		},
		Matcher: pulumi.String("PRIMARY"),
	})
	if err != nil {
		return err
	}

	// Create an HTTP proxy to route requests to the URLMap.
	// https://www.pulumi.com/registry/packages/gcp/api-docs/compute/targethttpsproxy/#target-https-proxy-certificate-manager-certificate
	httpsProxy, err := compute.NewTargetHttpsProxy(ctx, "http-proxy", &compute.TargetHttpsProxyArgs{
		// CertificateManagerCertificates: pulumi.StringArray{pulumi.Sprintf("//certificatemanager.googleapis.com/%v", sslCert.ID())},
		CertificateMap: pulumi.Sprintf("//certificatemanager.googleapis.com/%v", certMap.ID()),
		UrlMap:         httpsUrlMap.SelfLink,
	})
	if err != nil {
		return err
	}

	// Create a GlobalForwardingRule rule to route requests to the HTTP proxy.
	_, err = compute.NewGlobalForwardingRule(ctx, "http-forwarding-rule", &compute.GlobalForwardingRuleArgs{
		IpAddress:  ip.Address,
		IpProtocol: pulumi.String("TCP"),
		PortRange:  pulumi.String("443"),
		Target:     httpsProxy.SelfLink,
	})
	if err != nil {
		return err
	}

	if a.GcpConfig.CdnDomain.SkipCacheInvalidation {
		return nil
	}

	sortedMd5Result := a.websiteFileMd5Outputs.ToArrayOutput().ApplyT(func(arr []interface{}) string {
		md5Strings := []string{}
		for _, md5 := range arr {
			if md5Str, ok := md5.(string); ok && md5Str != "" {
				md5Strings = append(md5Strings, md5Str)
			}
		}

		sort.Strings(md5Strings)

		return strings.Join(md5Strings, "")
	}).(pulumi.StringOutput)

	// Invalidate the CDN Cache
	_, err = local.NewCommand(ctx, "invalidate-cache", &local.CommandArgs{
		Create:  pulumi.Sprintf("gcloud compute url-maps invalidate-cdn-cache %s --path '/*' --async", httpsUrlMap.Name),
		Logging: local.LoggingStdoutAndStderr,
		Triggers: pulumi.Array{
			sortedMd5Result,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

// Website - Implements the Website deployment method for the GCP provider
func (a *NitricGcpPulumiProvider) Website(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Website) error {
	if a.GcpConfig.CdnDomain.DomainName == "" {
		return fmt.Errorf("website deployments to GCP require a domain name to be configured in the stack file.")
	}

	var err error

	indexDoc := config.GetIndexDocument()
	if indexDoc == "" {
		indexDoc = "index.html"
	}

	errorDoc := config.GetErrorDocument()
	if errorDoc == "" {
		errorDoc = "404.html"
	}

	a.WebsiteBuckets[config.BasePath], err = storage.NewBucket(ctx, fmt.Sprintf("%s-site", name), &storage.BucketArgs{
		Location: pulumi.String(a.Region),
		Website: &storage.BucketWebsiteArgs{
			MainPageSuffix: pulumi.String(indexDoc),
			NotFoundPage:   pulumi.String(errorDoc),
		},
	})
	if err != nil {
		return err
	}

	_, err = storage.NewBucketIAMBinding(ctx, fmt.Sprintf("%s-site-bucket-iam", name), &storage.BucketIAMBindingArgs{
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

		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		relativePath, err := filepath.Rel(config.GetLocalDirectory(), path)
		if err != nil {
			return err
		}

		// Clean the relative path to ensure it is URL-safe and cross platform
		// This is important so files from a Windows host don't use backslashes for bucket keys
		cleanedRelativePath := filepath.ToSlash(relativePath)

		siteFile, err := storage.NewBucketObject(ctx, fmt.Sprintf("%s-%s", name, cleanedRelativePath), &storage.BucketObjectArgs{
			Bucket:      a.WebsiteBuckets[config.BasePath].Name,
			Name:        pulumi.String(cleanedRelativePath),
			Source:      pulumi.NewFileAsset(path),
			ContentType: pulumi.String(contentType),
		})
		if err != nil {
			return err
		}

		a.websiteFileMd5Outputs = append(a.websiteFileMd5Outputs, siteFile.Md5hash)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
