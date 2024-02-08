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
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func GetSubName(serviceName string, topicName string) string {
	return fmt.Sprintf("%s-%s-sub", serviceName, topicName)
}

func (p *NitricGcpPulumiProvider) Topic(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Topic) error {
	var err error
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	p.topics[name], err = pubsub.NewTopic(ctx, name, &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(common.Tags(p.stackId, name, resources.Topic)),
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return err
	}

	for _, sub := range config.Subscriptions {
		targetService := p.cloudRunServices[sub.GetService()]

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
