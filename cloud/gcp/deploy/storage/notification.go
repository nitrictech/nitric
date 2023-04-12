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

package storage

import (
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudStorageNotification struct {
	pulumi.ResourceState

	Name         string
	Notification *storage.Notification
}

type CloudStorageNotificationArgs struct {
	Location  string
	StackID   pulumi.StringInput

	Bucket *CloudStorageBucket
	Config *v1.BucketNotificationConfig
	Function *exec.CloudRunner
}

func EventTypeToStorageEventType(eventType *v1.EventType) []string {
	switch *eventType {
	case v1.EventType_All:
		return []string{"OBJECT_FINALIZE", "OBJECT_DELETE"}
	case v1.EventType_Created:
		return []string{"OBJECT_FINALIZE"}
	case v1.EventType_Deleted:
		return []string{"OBJECT_DELETE"}
	default:
		return []string{}
	}
}

func NewCloudStorageNotification(ctx *pulumi.Context, name string, args *CloudStorageNotificationArgs, opts ...pulumi.ResourceOption) (*CloudStorageNotification, error) {
	res := &CloudStorageNotification{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:bucket:GCPCloudStorageNotification", name, res, opts...)
	if err != nil {
		return nil, err
	}

	topic, err := pubsub.NewTopic(ctx, name+"-notificationtopic", &pubsub.TopicArgs{
		Labels: common.Tags(ctx, args.StackID, name),
	}, opts...)
	if err != nil {
		return nil, err
	}

	invokerAccount, err := serviceaccount.NewAccount(ctx, name+"subacct", &serviceaccount.AccountArgs{
		// accountId accepts a max of 30 chars, limit our generated name to this length
		AccountId: pulumi.String(utils.StringTrunc(name, 30-8) + "subacct"),
	}, append(opts, pulumi.Parent(res))...)
	if err != nil {
		return nil, errors.WithMessage(err, "invokerAccount "+name)
	}

	// Apply permissions for the above account to the newly deployed cloud run service
	_, err = cloudrun.NewIamMember(ctx, name+"-subrole", &cloudrun.IamMemberArgs{
		Member:   pulumi.Sprintf("serviceAccount:%s", invokerAccount.Email),
		Role:     pulumi.String("roles/run.invoker"),
		Service:  args.Function.Service.Name,
		Location: args.Function.Service.Location,
	}, append(opts, pulumi.Parent(res))...)
	if err != nil {
		return nil, errors.WithMessage(err, "iam member "+name)
	}

	_, err = pubsub.NewSubscription(ctx, name, &pubsub.SubscriptionArgs{
		Topic:              topic.Name,
		AckDeadlineSeconds: pulumi.Int(300),
		RetryPolicy: pubsub.SubscriptionRetryPolicyArgs{
			MinimumBackoff: pulumi.String("15s"),
			MaximumBackoff: pulumi.String("600s"),
		},
		PushConfig: pubsub.SubscriptionPushConfigArgs{
			OidcToken: pubsub.SubscriptionPushConfigOidcTokenArgs{
				ServiceAccountEmail: invokerAccount.Email,
			},
			PushEndpoint: pulumi.Sprintf("%s/x-nitric-notification/bucket/%s", args.Function.Url, args.Bucket.Name),
		},
		ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
			Ttl: pulumi.String(""),
		},
	}, append(opts, pulumi.Parent(topic))...)
	if err != nil {
		return nil, errors.WithMessage(err, "subscription "+name+"-sub")
	}

	res.Notification, err = storage.NewNotification(ctx, name, &storage.NotificationArgs{
		Bucket:           args.Bucket.CloudStorage.Name,
		PayloadFormat:    pulumi.String("JSON_API_V1"),
		Topic:            topic.ID(),
		EventTypes:       pulumi.ToStringArray(EventTypeToStorageEventType(&args.Config.EventType)),
		ObjectNamePrefix: pulumi.String(args.Config.EventFilter),
	}, append(opts, pulumi.DependsOn([]pulumi.Resource{topic}))...)
	if err != nil {
		return nil, err
	}

	return res, nil
}