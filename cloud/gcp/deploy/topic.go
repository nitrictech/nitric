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
	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PubSubTopic struct {
	pulumi.ResourceState

	Name   string
	PubSub *pubsub.Topic
}

type PubSubTopicArgs struct {
	Location  string
	StackID   string
	ProjectId string

	Topic *v1.Topic
}

func NewPubSubTopic(ctx *pulumi.Context, name string, args *PubSubTopicArgs, opts ...pulumi.ResourceOption) (*PubSubTopic, error) {
	res := &PubSubTopic{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:topic:GCPPubSubTopic", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.PubSub, err = pubsub.NewTopic(ctx, name, &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(common.Tags(args.StackID, name, resources.Topic)),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

type PubSubSubscription struct {
	pulumi.ResourceState

	Name string

	Subscription *pubsub.Subscription
}

type PubSubSubscriptionArgs struct {
	Function *exec.CloudRunner
	Topic    *PubSubTopic
}

func GetSubName(executionName string, topicName string) string {
	return fmt.Sprintf("%s-%s-sub", executionName, topicName)
}

func (p *NitricGcpPulumiProvider) Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error {
	var err error
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	p.topics[name], err = pubsub.NewTopic(ctx, name, &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(common.Tags(p.stackId, name, resources.Topic)),
	})
	if err != nil {
		return err
	}

	for _, sub := range config.Subscriptions {
		targetService := p.cloudRunServices[sub.GetExecutionUnit()]

		_, err := pubsub.NewSubscription(ctx, GetSubName(targetService.Name, name), &pubsub.SubscriptionArgs{
			Topic:              p.topics[name].Name, // The GCP topic name
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
				PushEndpoint: pulumi.Sprintf("%s/x-nitric-topic/%s?token=%s", targetService.Url, p.topics[name].Name, targetService.EventToken),
			},
			ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
				Ttl: pulumi.String(""),
			},
		}, opts...)
		if err != nil {
			return errors.WithMessage(err, "subscription "+name+"-sub")
		}
	}

	return nil
}
