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
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (p *NitricGcpPulumiProvider) Bucket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Bucket) error {
	var err error
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	resourceLabels := common.Tags(p.stackId, name, resources.Bucket)

	p.buckets[name], err = storage.NewBucket(ctx, name, &storage.BucketArgs{
		Location: pulumi.String(p.region),
		Labels:   pulumi.ToStringMap(resourceLabels),
	}, opts...)
	if err != nil {
		return err
	}

	for _, listener := range config.Listeners {
		if err := p.newCloudStorageNotification(ctx, parent, name, listener); err != nil {
			return err
		}
	}

	return nil
}

func (p *NitricGcpPulumiProvider) newCloudStorageNotification(ctx *pulumi.Context, parent pulumi.Resource, bucketName string, listener *deploymentspb.BucketListener) error {
	var err error
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	name := bucketName + "-" + listener.GetService()

	if listener == nil || listener.Config == nil {
		return fmt.Errorf("invalid config provided for bucket notification")
	}

	topic, err := pubsub.NewTopic(ctx, name+"-topic", &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(common.Tags(p.stackId, name, resources.Bucket)),
	}, opts...)
	if err != nil {
		return err
	}

	targetService, ok := p.cloudRunServices[listener.GetService()]
	if !ok {
		return fmt.Errorf("unable to find target service for bucket listener: %s", listener.GetService())
	}

	targetBucket, ok := p.buckets[bucketName]
	if !ok {
		return fmt.Errorf("unable to find target bucket for bucket listener: %s", bucketName)
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
				ServiceAccountEmail: targetService.Invoker.Email,
			},
			PushEndpoint: pulumi.Sprintf("%s/x-nitric-notification/bucket/%s?token=%s", targetService.Url, bucketName, targetService.EventToken),
		},
		ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
			Ttl: pulumi.String(""),
		},
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, "subscription "+name)
	}

	// Give the cloud storage service account publishing permissions
	gcsAccount, err := storage.GetProjectServiceAccount(ctx, nil, nil)
	if err != nil {
		return err
	}

	binding, err := pubsub.NewTopicIAMBinding(ctx, name+"-binding", &pubsub.TopicIAMBindingArgs{
		Topic: topic.ID(),
		Role:  pulumi.String("roles/pubsub.publisher"),
		Members: pulumi.StringArray{
			pulumi.String(fmt.Sprintf("serviceAccount:%v", gcsAccount.EmailAddress)),
		},
	})
	if err != nil {
		return errors.WithMessage(err, "topic binding "+name)
	}

	prefix := listener.Config.KeyPrefixFilter
	if prefix == "*" {
		prefix = ""
	}

	_, err = storage.NewNotification(ctx, name, &storage.NotificationArgs{
		Bucket:           targetBucket.Name,
		PayloadFormat:    pulumi.String("JSON_API_V1"),
		Topic:            topic.ID(),
		EventTypes:       pulumi.ToStringArray(notificationTypeToStorageEventType(listener.Config.BlobEventType)),
		ObjectNamePrefix: pulumi.String(prefix),
	}, append(opts, pulumi.DependsOn([]pulumi.Resource{binding}))...)
	if err != nil {
		return errors.WithMessage(err, "storage notification "+name)
	}

	return nil
}

func notificationTypeToStorageEventType(eventType storagepb.BlobEventType) []string {
	switch eventType {
	case storagepb.BlobEventType_Created:
		return []string{"OBJECT_FINALIZE"}
	case storagepb.BlobEventType_Deleted:
		return []string{"OBJECT_DELETE"}
	default:
		return []string{}
	}
}
