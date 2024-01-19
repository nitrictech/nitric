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

// Bucket - Implements deployments of Nitric Buckets using AWS S3
func (a *NitricAwsPulumiProvider) Bucket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Bucket) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	bucket, err := s3.NewBucket(ctx, name, &s3.BucketArgs{
		Tags: pulumi.ToStringMap(common.Tags(a.stackId, name, resources.Bucket)),
	}, opts...)
	if err != nil {
		return err
	}

	a.buckets[name] = bucket

	if len(config.Listeners) > 0 {
		notificationName := fmt.Sprintf("notification-%s", name)
		notification, err := createNotification(ctx, notificationName, &S3NotificationArgs{
			StackID:   a.stackId,
			Location:  a.region,
			Bucket:    bucket,
			Lambdas:   a.lambdas,
			Listeners: config.Listeners,
		}, opts...)
		if err != nil {
			return err
		}

		a.bucketNotifications[name] = notification
	}

	return nil
}
