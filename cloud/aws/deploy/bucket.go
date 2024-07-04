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
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
)

type S3NotificationArgs struct {
	Location string
	StackID  string

	Bucket    *s3.Bucket
	Listeners []*deploymentspb.BucketListener
	Lambdas   map[string]*lambda.Function
}

func eventTypeToStorageEventType(eventType *storagepb.BlobEventType) []string {
	switch *eventType {
	case storagepb.BlobEventType_Created:
		return []string{"s3:ObjectCreated:*"}
	case storagepb.BlobEventType_Deleted:
		return []string{"s3:ObjectRemoved:*"}
	default:
		return []string{}
	}
}

// createNotification creates an AWS S3 bucket notification, containing all target lambda functions and their filters
func createNotification(ctx *pulumi.Context, name string, args *S3NotificationArgs, opts ...pulumi.ResourceOption) (*s3.BucketNotification, error) {
	invokePerms := map[string]pulumi.Resource{}
	notificationTargetLambdas := s3.BucketNotificationLambdaFunctionArray{}

	for _, listener := range args.Listeners {
		// Get the deployed service
		funcName := listener.GetService()
		lambdaFunc, ok := args.Lambdas[funcName]
		if !ok {
			return nil, fmt.Errorf("invalid service %s given for bucket subscription", funcName)
		}

		// Don't create duplicate permissions
		if invokePerms[funcName] == nil {
			perm, err := lambda.NewPermission(ctx, name+"-"+funcName, &lambda.PermissionArgs{
				Action:    pulumi.String("lambda:InvokeFunction"),
				Function:  lambdaFunc.Arn,
				Principal: pulumi.String("s3.amazonaws.com"),
				SourceArn: args.Bucket.Arn,
			}, opts...)
			if err != nil {
				return nil, fmt.Errorf("unable to create lambda invoke permission: %w", err)
			}

			invokePerms[funcName] = perm
		}

		if listener.Config.KeyPrefixFilter == "*" {
			listener.Config.KeyPrefixFilter = ""
		}

		notificationTargetLambdas = append(notificationTargetLambdas, s3.BucketNotificationLambdaFunctionArgs{
			LambdaFunctionArn: lambdaFunc.Arn,
			Events: pulumi.ToStringArray(
				eventTypeToStorageEventType(&listener.Config.BlobEventType),
			),
			FilterPrefix: pulumi.String(listener.Config.KeyPrefixFilter),
		}.ToBucketNotificationLambdaFunctionOutput())
	}

	notificationOptions := append([]pulumi.ResourceOption{pulumi.DependsOn(lo.Values(invokePerms))}, opts...)

	notification, err := s3.NewBucketNotification(ctx, name, &s3.BucketNotificationArgs{
		Bucket:          args.Bucket.ID(),
		LambdaFunctions: notificationTargetLambdas,
	}, notificationOptions...)
	if err != nil {
		return nil, fmt.Errorf("unable to create bucket notification: %w", err)
	}

	return notification, nil
}

// extractBucketName - extracts the bucket name from an S3 ARN.
func extractBucketName(arn string) (string, error) {
	s3ArnRegex := regexp.MustCompile(`(?i)^arn:aws:s3:::([^/]+)`)

	matches := s3ArnRegex.FindStringSubmatch(arn)
	if len(matches) < 2 {
		return "", fmt.Errorf("invalid S3 bucket ARN: %s", arn)
	}

	bucketName := matches[1]

	if bucketName == "" {
		return "", fmt.Errorf("invalid S3 bucket ARN: bucket name could not be extracted from %s", arn)
	}

	return bucketName, nil
}

// importBucket - tags an existing bucket in AWS and adds it to the stack.
func importBucket(ctx *pulumi.Context, name string, importIdentifier string, opts []pulumi.ResourceOption, tags map[string]string, tagClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI) (*s3.Bucket, error) {
	// Allow bucket names or ARNs as import identifiers
	bucketName, err := extractBucketName(importIdentifier)
	if err != nil {
		bucketName = importIdentifier
	}

	bucketLookup, err := s3.LookupBucket(ctx, &s3.LookupBucketArgs{
		Bucket: bucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("unable to lookup imported S3 bucket %s: %w", bucketName, err)
	}

	_, err = tagClient.TagResources(&resourcegroupstaggingapi.TagResourcesInput{
		ResourceARNList: aws.StringSlice([]string{bucketLookup.Arn}),
		Tags:            aws.StringMap(tags),
	})
	if err != nil {
		return nil, fmt.Errorf("unable to tag imported S3 bucket %s: %w", bucketName, err)
	}

	// nitric didn't create this resource, so it shouldn't delete it either.
	allOpts := append(opts, pulumi.RetainOnDelete(true))

	bucket, err := s3.GetBucket(
		ctx,
		name,
		pulumi.ID(bucketLookup.Id),
		nil,
		allOpts...,
	)
	if err != nil {
		return nil, fmt.Errorf("unable to import S3 bucket %s: %w", bucketName, err)
	}

	return bucket, nil
}

// createBucket - creates a new S3 bucket in AWS and tags it.
func createBucket(ctx *pulumi.Context, name string, opts []pulumi.ResourceOption, tags map[string]string) (*s3.Bucket, error) {
	bucket, err := s3.NewBucket(ctx, name, &s3.BucketArgs{
		Tags: pulumi.ToStringMap(tags),
	}, opts...)
	if err != nil {
		return nil, err
	}

	return bucket, nil
}

// Bucket - Implements deployments of Nitric Buckets using AWS S3
func (a *NitricAwsPulumiProvider) Bucket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Bucket) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}
	tags := common.Tags(a.StackId, name, resources.Bucket)

	var err error
	var bucket *s3.Bucket

	importArn := ""
	if a.AwsConfig.Import.Buckets != nil {
		importArn = a.AwsConfig.Import.Buckets[name]
	}

	if importArn != "" {
		bucket, err = importBucket(ctx, name, importArn, opts, tags, a.ResourceTaggingClient)
	} else {
		bucket, err = createBucket(ctx, name, opts, tags)
	}

	if err != nil {
		return err
	}

	a.Buckets[name] = bucket

	if len(config.Listeners) > 0 {
		notificationName := fmt.Sprintf("notification-%s", name)
		notification, err := createNotification(ctx, notificationName, &S3NotificationArgs{
			StackID:   a.StackId,
			Location:  a.Region,
			Bucket:    bucket,
			Lambdas:   a.Lambdas,
			Listeners: config.Listeners,
		}, opts...)
		if err != nil {
			return err
		}

		a.BucketNotifications[name] = notification
	}

	return nil
}
