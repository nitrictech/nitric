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

	pubsubv1 "cloud.google.com/go/pubsub/apiv1"
	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// tagTopic - applies nitric tags to an existing topic in GCP and adds it to the stack.
func tagTopic(ctx *pulumi.Context, name string, projectId string, topicName string, tags map[string]string, client *pubsubv1.PublisherClient, opts []pulumi.ResourceOption) (*pubsub.Topic, error) {
	topicLookup, err := pubsub.LookupTopic(ctx, &pubsub.LookupTopicArgs{
		Project: &projectId,
		Name:    topicName,
	})
	if err != nil {
		return nil, err
	}

	_, err = client.UpdateTopic(ctx.Context(), &pubsubpb.UpdateTopicRequest{
		Topic: &pubsubpb.Topic{
			Name:   topicLookup.Id,
			Labels: tags,
		},
		UpdateMask: &fieldmaskpb.FieldMask{
			Paths: []string{"labels"},
		},
	})
	if err != nil {
		return nil, err
	}

	topic, err := pubsub.GetTopic(
		ctx,
		name,
		pulumi.ID(topicLookup.Id),
		nil,
		// nitric didn't create this resource, so it shouldn't delete it either.
		append(opts, pulumi.RetainOnDelete(true))...,
	)
	if err != nil {
		return nil, err
	}
	return topic, nil
}

func createTopic(ctx *pulumi.Context, name string, stackId string, tags map[string]string, opts []pulumi.ResourceOption) (*pubsub.Topic, error) {
	return pubsub.NewTopic(ctx, name, &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(tags),
	}, opts...)
}

func GetSubName(serviceName string, topicName string) string {
	return fmt.Sprintf("%s-%s-sub", serviceName, topicName)
}

func (p *NitricGcpPulumiProvider) Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error {
	var err error
	var topic *pubsub.Topic
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	if gcpName, ok := p.GcpConfig.Import.Topics[name]; ok {
		topic, err = tagTopic(ctx, name, p.GcpConfig.ProjectId, gcpName, common.Tags(p.StackId, name, resources.Topic), p.PubsubClient, opts)
	} else {
		topic, err = createTopic(ctx, name, p.StackName, common.Tags(p.StackId, name, resources.Topic), opts)
	}

	if err != nil {
		return err
	}

	p.Topics[name] = topic

	for _, sub := range config.Subscriptions {
		targetService := p.CloudRunServices[sub.GetService()]

		_, err := pubsub.NewSubscription(ctx, GetSubName(targetService.Name, name), &pubsub.SubscriptionArgs{
			Topic:              p.Topics[name].Name, // The GCP topic name
			AckDeadlineSeconds: pulumi.Int(300),
			RetryPolicy: pubsub.SubscriptionRetryPolicyArgs{
				MinimumBackoff: pulumi.String("15s"),
				MaximumBackoff: pulumi.String("600s"),
			},
			PushConfig: pubsub.SubscriptionPushConfigArgs{
				OidcToken: pubsub.SubscriptionPushConfigOidcTokenArgs{
					ServiceAccountEmail: targetService.Invoker.Email,
				},
				// https://cloud.google.com/appengine/docs/flexible/writing-and-responding-to-pub-sub-messages?tab=go#top
				PushEndpoint: pulumi.Sprintf("%s/x-nitric-topic/%s?token=%s", targetService.Url, name, targetService.EventToken),
			},
			ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
				Ttl: pulumi.String(""),
			},
		}, p.WithDefaultResourceOptions(opts...)...)
		if err != nil {
			return errors.WithMessage(err, "subscription "+name+"-sub")
		}
	}

	return nil
}
