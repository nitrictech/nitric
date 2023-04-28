// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bucket

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
)

func eventTypeToStorageEventType(eventType *v1.BucketNotificationType) []string {
	switch *eventType {
	case v1.BucketNotificationType_All:
		return []string{"s3:ObjectCreated:*", "s3:ObjectRemoved:*"}
	case v1.BucketNotificationType_Created:
		return []string{"s3:ObjectCreated:*"}
	case v1.BucketNotificationType_Deleted:
		return []string{"s3:ObjectRemoved:*"}
	default:
		return []string{}
	}
}

type S3Notification struct {
	pulumi.ResourceState

	Name         string
	Notification *s3.BucketNotification
}

type S3NotificationArgs struct {
	Location string
	StackID  pulumi.StringInput

	Bucket       *S3Bucket
	Notification []*deploy.BucketNotificationTarget
	Functions    map[string]*exec.LambdaExecUnit
}

func NewS3Notification(ctx *pulumi.Context, name string, args *S3NotificationArgs, opts ...pulumi.ResourceOption) (*S3Notification, error) {
	res := &S3Notification{
		Name: name,
	}
	err := ctx.RegisterComponentResource("nitric:bucket:AWSS3Notification", name, res, opts...)
	if err != nil {
		return nil, err
	}

	invokePerms := map[string]pulumi.Resource{}
	bucketNotifications := s3.BucketNotificationLambdaFunctionArray{}

	for _, notification := range args.Notification {
		// Get the deployed execution unit
		funcName := notification.GetExecutionUnit()
		unit, ok := args.Functions[funcName]
		if !ok {
			return nil, fmt.Errorf("invalid execution unit %s given for bucket subscription", funcName)
		}

		// Don't create duplicate permissions
		if invokePerms[funcName] == nil {
			perm, err := lambda.NewPermission(ctx, name+"-"+funcName, &lambda.PermissionArgs{
				Action:    pulumi.String("lambda:InvokeFunction"),
				Function:  unit.Function.Arn,
				Principal: pulumi.String("s3.amazonaws.com"),
				SourceArn: args.Bucket.S3.Arn,
			})
			if err != nil {
				return nil, fmt.Errorf("unable to create lambda invoke permission: %w", err)
			}

			invokePerms[funcName] = perm
		}

		if notification.Config.EventFilter == "*" {
			notification.Config.EventFilter = ""
		}

		// Append notification
		bucketNotifications = append(bucketNotifications, s3.BucketNotificationLambdaFunctionArgs{
			LambdaFunctionArn: unit.Function.Arn,
			Events: pulumi.ToStringArray(
				eventTypeToStorageEventType(&notification.Config.EventType),
			),
			FilterPrefix: pulumi.String(notification.Config.EventFilter),
		}.ToBucketNotificationLambdaFunctionOutput())
	}

	res.Notification, err = s3.NewBucketNotification(ctx, name, &s3.BucketNotificationArgs{
		Bucket:          args.Bucket.S3.ID(),
		LambdaFunctions: bucketNotifications,
	}, pulumi.DependsOn(lo.Values(invokePerms)))
	if err != nil {
		return nil, fmt.Errorf("unable to create bucket notification: %w", err)
	}

	return res, nil
}