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
	"context"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"mime"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/nitrictech/nitric/cloud/aws/deploy/embeds"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/cloudfront"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"

	"github.com/aws/aws-sdk-go-v2/config"
	awscloudfront "github.com/aws/aws-sdk-go-v2/service/cloudfront"
	awscloudfronttypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
)

// createBucket - creates a new S3 bucket in AWS and tags it.
func (a *NitricAwsPulumiProvider) createWebsiteBucket(ctx *pulumi.Context) error {
	var err error

	tags := common.Tags(a.StackId, a.websiteBucketName, resources.Website)

	a.publicWebsiteBucket, err = s3.NewBucket(ctx, a.websiteBucketName, &s3.BucketArgs{
		Tags: pulumi.ToStringMap(tags),
	})
	if err != nil {
		return err
	}

	return nil
}

func fileETag(ctx context.Context, client *awss3.Client, bucketName, key string) (string, error) {
	headObjectOutput, err := client.HeadObject(ctx, &awss3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		var notFound *s3types.NotFound

		// If the file does not exist, return an empty string so we can ignore checking it
		if errors.As(err, &notFound) {
			return "", nil
		}

		// Otherwise, return the error
		return "", err
	}

	// Trim the ETag to remove the quotes
	etag := strings.Trim(*headObjectOutput.ETag, "\"")

	return etag, nil
}

func (a *NitricAwsPulumiProvider) getS3Client() (*awss3.Client, error) {
	// Load the AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(a.Region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create an S3 client
	return awss3.NewFromConfig(cfg), nil
}

// Website - Implements the Website deployment method for the AWS provider
func (a *NitricAwsPulumiProvider) Website(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Website) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{a.publicWebsiteBucket})}

	localDir, ok := config.AssetSource.(*deploymentspb.Website_LocalDirectory)
	if !ok {
		return fmt.Errorf("unsupported asset source type for website: %s", name)
	}

	cleanedPath := filepath.ToSlash(filepath.Clean(localDir.LocalDirectory))

	if config.BasePath == "/" {
		a.websiteIndexDocument = config.IndexDocument
		a.websiteErrorDocument = config.ErrorDocument
	}

	// get the S3 client for reading the ETag of existing files
	client, err := a.getS3Client()
	if err != nil {
		return err
	}

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

		existingTag := a.publicWebsiteBucket.Bucket.ApplyT(func(bucket string) (string, error) {
			// Check if the object already exists
			existingETag, err := fileETag(context.TODO(), client, bucket, strings.TrimPrefix(objectKey, "/"))
			if err != nil {
				return "", err
			}

			return existingETag, nil
		})

		obj, err := s3.NewBucketObject(ctx, arn, &s3.BucketObjectArgs{
			Bucket:      a.publicWebsiteBucket.Bucket,
			Source:      pulumi.NewFileAsset(path),
			ContentType: pulumi.String(contentType),
			Key:         pulumi.String(objectKey),
		}, opts...)
		if err != nil {
			return err
		}

		keyToInvalidate := pulumi.All(obj.Etag, existingTag).ApplyT(func(args []any) (string, error) {
			newEtag, newEtagOk := args[0].(string)
			existingEtag, existingEtagOk := args[1].(string)

			if !newEtagOk || !existingEtagOk {
				return "", fmt.Errorf("failed to assert ETag types")
			}

			// if an existing ETag is present and it is different from the new ETag, return the key to invalidate
			if existingEtag != "" && newEtag != existingEtag {
				return objectKey, nil
			}

			return "", nil
		}).(pulumi.StringOutput)

		a.websiteChangedFileOutputs = append(a.websiteChangedFileOutputs, keyToInvalidate)

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (a *NitricAwsPulumiProvider) deployCloudfrontDistribution(ctx *pulumi.Context) error {
	origins := cloudfront.DistributionOriginArray{}
	orderedCacheBehaviors := cloudfront.DistributionOrderedCacheBehaviorArray{}

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

	// Add the public bucket as an origin
	origins = append(origins, &cloudfront.DistributionOriginArgs{
		DomainName: a.publicWebsiteBucket.BucketRegionalDomainName,
		OriginId:   pulumi.String("publicOrigin"),
		S3OriginConfig: &cloudfront.DistributionOriginS3OriginConfigArgs{
			OriginAccessIdentity: oai.CloudfrontAccessIdentityPath,
		},
	})

	// Default cache behavior for the public bucket
	defaultCacheBehavior := &cloudfront.DistributionDefaultCacheBehaviorArgs{
		TargetOriginId:       pulumi.String("publicOrigin"),
		ViewerProtocolPolicy: pulumi.String("redirect-to-https"),
		AllowedMethods: pulumi.StringArray{
			pulumi.String("GET"),
			pulumi.String("HEAD"),
			pulumi.String("OPTIONS"),
		},
		CachedMethods: pulumi.StringArray{
			pulumi.String("GET"),
			pulumi.String("HEAD"),
			pulumi.String("OPTIONS"),
		},
		ForwardedValues: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesArgs{
			QueryString: pulumi.Bool(false),
			Cookies: &cloudfront.DistributionDefaultCacheBehaviorForwardedValuesCookiesArgs{
				Forward: pulumi.String("none"),
			},
		},
		// could be added to stack config in the future
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

	// Sort the APIs by name
	sortedApiKeys := lo.Keys(a.Apis)
	slices.Sort(sortedApiKeys)

	// For each API forward to the appropriate API gateway
	for _, name := range sortedApiKeys {
		api := a.Apis[name]

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

		orderedCacheBehaviors = append(orderedCacheBehaviors,
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

	name := fmt.Sprintf("%s-cdn", a.StackId)

	tags := common.Tags(a.StackId, name, resources.Website)

	// Deploy a CloudFront distribution for the S3 bucket
	a.Distribution, err = cloudfront.NewDistribution(ctx, name, &cloudfront.DistributionArgs{
		Origins:               origins,
		Enabled:               pulumi.Bool(true),
		DefaultCacheBehavior:  defaultCacheBehavior,
		DefaultRootObject:     pulumi.String(a.websiteIndexDocument),
		OrderedCacheBehaviors: orderedCacheBehaviors,
		Restrictions: &cloudfront.DistributionRestrictionsArgs{
			GeoRestriction: &cloudfront.DistributionRestrictionsGeoRestrictionArgs{
				RestrictionType: pulumi.String("none"),
			},
		},
		Tags: pulumi.ToStringMap(tags),
		ViewerCertificate: &cloudfront.DistributionViewerCertificateArgs{
			CloudfrontDefaultCertificate: pulumi.Bool(true),
		},
		CustomErrorResponses: cloudfront.DistributionCustomErrorResponseArray{
			&cloudfront.DistributionCustomErrorResponseArgs{
				ErrorCode:        pulumi.Int(404),
				ResponseCode:     pulumi.Int(200),
				ResponsePagePath: pulumi.String(fmt.Sprintf("/%v", a.websiteErrorDocument)),
			},
			// Redirect all 403 errors to the error page, s3 by default will return a 403 for missing files
			&cloudfront.DistributionCustomErrorResponseArgs{
				ErrorCode:        pulumi.Int(403),
				ResponseCode:     pulumi.Int(200),
				ResponsePagePath: pulumi.String(fmt.Sprintf("/%v", a.websiteErrorDocument)),
			},
		},
	})
	if err != nil {
		return err
	}

	// apply invalidation on the distribution when files change
	pulumi.All(a.Distribution.ID().ToStringOutput(), a.websiteChangedFileOutputs.ToStringArrayOutput()).ApplyT(func(args []interface{}) error {
		cdnID := args[0].(string)
		websiteChangedFileKeys := []string{}

		// Filter out empty strings from the array
		for _, key := range args[1].([]string) {
			if key != "" {
				websiteChangedFileKeys = append(websiteChangedFileKeys, key)
			}
		}

		if len(websiteChangedFileKeys) > 0 {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(a.Region))
			if err != nil {
				return fmt.Errorf("failed to load AWS config: %w", err)
			}

			// Create CloudFront client
			client := awscloudfront.NewFromConfig(cfg)

			quantity, err := SafeInt32(len(websiteChangedFileKeys))
			if err != nil {
				return err
			}

			input := awscloudfront.CreateInvalidationInput{
				DistributionId: &cdnID,
				InvalidationBatch: &awscloudfronttypes.InvalidationBatch{
					CallerReference: aws.String(time.Now().Format("2006-01-02 15:04:05")),
					Paths: &awscloudfronttypes.Paths{
						Quantity: aws.Int32(quantity),
						Items:    websiteChangedFileKeys,
					},
				},
			}

			_, err = client.CreateInvalidation(context.TODO(), &input)
			if err != nil {
				return fmt.Errorf("failed to create CloudFront invalidation: %w", err)
			}
		}

		return nil
	})

	ctx.Export("cdn", pulumi.Sprintf("https://%s", a.Distribution.DomainName))

	return nil
}

// SafeInt32 - Safely convert an int to an int32
func SafeInt32(n int) (int32, error) {
	if n > math.MaxInt32 {
		return 0, fmt.Errorf("value exceeds int32 limit: %d", n)
	}

	return int32(n), nil //#nosec G115 -- n is checked to be within int32 range
}
