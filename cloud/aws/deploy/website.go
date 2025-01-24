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

	"github.com/nitrictech/nitric/cloud/aws/deploy/embeds"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createBucket - creates a new S3 bucket in AWS and tags it.
func (a *NitricAwsPulumiProvider) createWebsiteBucket(ctx *pulumi.Context) error {
	var err error

	name := "public-website-bucket"

	a.publicWebsiteBucket, err = s3.NewBucket(ctx, name, nil)
	if err != nil {
		return err
	}

	a.Buckets[name] = a.publicWebsiteBucket

	return nil
}

func (a *NitricAwsPulumiProvider) deployCloudfrontDistribution(ctx *pulumi.Context) error {
	origins := cloudfront.DistributionOriginArray{}
	var defaultCacheBehaviour *cloudfront.DistributionDefaultCacheBehaviorArgs = nil
	orderedCacheBeviours := cloudfront.DistributionOrderedCacheBehaviorArray{}

	oai, err := cloudfront.NewOriginAccessIdentity(ctx, "oai", &cloudfront.OriginAccessIdentityArgs{
		Comment: pulumi.String("OAI for accessing S3 bucket"),
	})
	if err != nil {
		return err
	}

	policy := pulumi.All(a.publicWebsiteBucket.Arn, oai.IamArn).ApplyT(func(args []interface{}) (string, error) {
		bucketID, bucketIdOk := args[0].(string)
		oaiPath, oaiPathOk := args[1].(string)
		if !bucketIdOk || !oaiPathOk {
			return "", fmt.Errorf("failed to get bucket ID or OAI path")
		}
		return fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"AWS": "%s"
					},
					"Action": "s3:GetObject",
					"Resource": "%s/*"
				}
			]
		}`, oaiPath, bucketID), nil
	}).(pulumi.StringOutput)
	_, err = s3.NewBucketPolicy(ctx, "publicBucketPolicy", &s3.BucketPolicyArgs{
		Bucket: a.publicWebsiteBucket.Bucket,
		Policy: policy,
	})
	if err != nil {
		return err
	}
	// We conventionally route to nitric resources from this distribution to create a single entry point
	// for the entire stack. e.g. /api/main/* will route to a nitric api named "main"
	apiRewriteFun, err := cloudfront.NewFunction(ctx, "api-url-rewrite-function", &cloudfront.FunctionArgs{
		Comment: pulumi.String("Rewrite API URLs routed to nitric services"),
		Code:    embeds.GetApiUrlRewriteFunction(),
		Runtime: pulumi.String("cloudfront-js-1.0"),
	})
	if err != nil {
		return err
	}

	rewriteFun, err := cloudfront.NewFunction(ctx, "url-rewrite-function", &cloudfront.FunctionArgs{
		Comment: pulumi.String("Rewrite URLs to default index document"),
		Code:    embeds.GetUrlRewriteFunction(),
		Runtime: pulumi.String("cloudfront-js-1.0"),
	})
	if err != nil {
		return err
	}

	// TODO handle multiple websites
	if a.publicWebsiteBucket != nil {
		origins = append(origins, &cloudfront.DistributionOriginArgs{
			DomainName: a.publicWebsiteBucket.BucketRegionalDomainName,
			OriginId:   pulumi.String("publicOrigin"),
			S3OriginConfig: &cloudfront.DistributionOriginS3OriginConfigArgs{
				OriginAccessIdentity: oai.CloudfrontAccessIdentityPath,
			},
		})
		defaultCacheBehaviour = &cloudfront.DistributionDefaultCacheBehaviorArgs{
			TargetOriginId:       pulumi.String("publicOrigin"),
			ViewerProtocolPolicy: pulumi.String("redirect-to-https"),
			AllowedMethods: pulumi.StringArray{
				pulumi.String("GET"),
				pulumi.String("HEAD"),
			},
			CachedMethods: pulumi.StringArray{
				pulumi.String("GET"),
				pulumi.String("HEAD"),
			},
			ForwardedValues: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{
				QueryString: pulumi.Bool(false),
				Cookies: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
					Forward: pulumi.String("none"),
				},
			},
			MinTtl:     pulumi.Int(0),
			DefaultTtl: pulumi.Int(3600),
			MaxTtl:     pulumi.Int(86400),
			FunctionAssociations: cloudfront.DistributionDefaultCacheBehaviorFunctionAssociationArray{
				&cloudfront.DistributionDefaultCacheBehaviorFunctionAssociationArgs{
					EventType:   pulumi.String("viewer-request"),
					FunctionArn: rewriteFun.Arn,
				},
			},
		}
	}
	// For each API forward to the appropriate API gateway
	for name, api := range a.Apis {
		apiDomainName := api.ApiEndpoint.ApplyT(func(endpoint string) string {
			return strings.Replace(endpoint, "https://", "", 1)
		}).(pulumi.StringOutput)
		origins = append(origins, &cloudfront.DistributionOriginArgs{
			DomainName: apiDomainName,
			OriginId:   pulumi.String(name),
			CustomOriginConfig: &cloudfront.DistributionOriginCustomOriginConfigArgs{
				OriginReadTimeout:    pulumi.Int(30),
				OriginProtocolPolicy: pulumi.String("https-only"),
				OriginSslProtocols: pulumi.StringArray{
					pulumi.String("TLSv1.2"),
					pulumi.String("SSLv3"),
				},
				HttpPort:  pulumi.Int(80),
				HttpsPort: pulumi.Int(443),
			},
		})
		orderedCacheBeviours = append(orderedCacheBeviours,
			&cloudfront.DistributionOrderedCacheBehaviorArgs{
				PathPattern: pulumi.Sprintf("api/%s/*", name),
				// rewrite the URL to the nitric service
				FunctionAssociations: cloudfront.DistributionOrderedCacheBehaviorFunctionAssociationArray{
					&cloudfront.DistributionOrderedCacheBehaviorFunctionAssociationArgs{
						EventType:   pulumi.String("viewer-request"),
						FunctionArn: apiRewriteFun.Arn,
					},
				},
				AllowedMethods: pulumi.ToStringArray([]string{"GET", "HEAD", "OPTIONS", "PUT", "POST", "PATCH", "DELETE"}),
				CachedMethods:  pulumi.ToStringArray([]string{"GET", "HEAD", "OPTIONS"}),
				TargetOriginId: pulumi.String(name),
				ForwardedValues: &cloudfront.DistributionOrderedCacheBehaviorForwardedValuesArgs{
					QueryString: pulumi.Bool(true),
					Cookies: &cloudfront.DistributionOrderedCacheBehaviorForwardedValuesCookiesArgs{
						Forward: pulumi.String("all"),
					},
					// Headers: pulumi.ToStringArray([]string{"*"}),
				},
				ViewerProtocolPolicy: pulumi.String("https-only"),
			},
		)
	}

	// Deploy a CloudFront distribution for the S3 bucket
	cdn, err := cloudfront.NewDistribution(ctx, "distribution", &cloudfront.DistributionArgs{
		Origins:               origins,
		Enabled:               pulumi.Bool(true),
		DefaultCacheBehavior:  defaultCacheBehaviour,
		DefaultRootObject:     pulumi.String("index.html"), // TODO use root website index document
		OrderedCacheBehaviors: orderedCacheBeviours,
		Restrictions: &cloudfront.DistributionRestrictionsArgs{
			GeoRestriction: &cloudfront.DistributionRestrictionsGeoRestrictionArgs{
				RestrictionType: pulumi.String("none"),
			},
		},
		ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
			CloudfrontDefaultCertificate: pulumi.Bool(true),
		},
	})
	if err != nil {
		return err
	}

	ctx.Export("cdn", pulumi.Sprintf("https://%s", cdn.DomainName))

	return nil
}

// Bucket - Implements deployments of Nitric Buckets using AWS S3
func (a *NitricAwsPulumiProvider) Website(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Website) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{a.publicWebsiteBucket})}
	//tags := common.Tags(a.StackId, name, resources.Website)

	// indexPath := filepath.ToSlash(filepath.Join(name, config.IndexDocument))
	// errorPrefixedPath := filepath.ToSlash(filepath.Join(name, config.ErrorDocument))

	// TODO probably not needed
	test, err := s3.NewBucketWebsiteConfigurationV2(ctx, "website-config-"+name, &s3.BucketWebsiteConfigurationV2Args{
		Bucket: a.publicWebsiteBucket.Bucket,
		IndexDocument: s3.BucketWebsiteConfigurationV2IndexDocumentArgs{
			Suffix: pulumi.String(config.IndexDocument),
		},
		ErrorDocument: s3.BucketWebsiteConfigurationV2ErrorDocumentArgs{
			Key: pulumi.String(config.ErrorDocument),
		},
	}, opts...)
	if err != nil {
		return err
	}

	// print test.WebsiteEndpoint with pulumi run
	ctx.Export(fmt.Sprintf("website-%s-endpoint", name), test.WebsiteEndpoint)

	cleanedPath := filepath.ToSlash(filepath.Clean(config.OutputDirectory))

	// _, err = synced.NewS3BucketFolder(ctx, "bucket-folder-"+name, &synced.S3BucketFolderArgs{
	// 	Path:       pulumi.String(cleanedPath),
	// 	BucketName: a.publicWebsiteBucket.Bucket,
	// }, opts...)
	// if err != nil {
	// 	return err
	// }
	// Enumerate the public directory in pwd and upload all files to the public bucket
	// This will be the source for our cloudfront distribution
	err = filepath.WalkDir(cleanedPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip directories
		if d.IsDir() {
			return nil
		}
		// Determine the content type based on the file extension
		contentType := mime.TypeByExtension(filepath.Ext(path))
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		// Generate the object key to include the folder structure
		var objectKey string

		filePath := path[len(cleanedPath):]

		arn := filepath.ToSlash(filepath.Join(name, filePath))

		// If the base path is not the root, include it in the object key
		if config.BasePath == "/" {
			objectKey = filepath.ToSlash(filePath)
		} else {
			objectKey = filepath.ToSlash(filepath.Join(config.BasePath, filePath))
		}

		_, err = s3.NewBucketObject(ctx, arn, &s3.BucketObjectArgs{
			Bucket:      a.publicWebsiteBucket.Bucket,
			Source:      pulumi.NewFileAsset(path),
			ContentType: pulumi.String(contentType),
			Key:         pulumi.String(objectKey),
		}, opts...)
		return err
	})
	if err != nil {
		return err
	}

	return nil
}
