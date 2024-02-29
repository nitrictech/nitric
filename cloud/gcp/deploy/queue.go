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
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// type PubSubTopic struct {
// 	pulumi.ResourceState

// 	Name         string
// 	PubSub       *pubsub.Topic
// 	Subscription *pubsub.Subscription
// }

// type PubSubTopicArgs struct {
// 	Location string
// 	StackID  string
// 	Queue    *v1.Queue
// }

func (p *NitricGcpPulumiProvider) Queue(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Queue) error {
	var err error
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	resourceLabels := common.Tags(p.StackId, name, resources.Queue)

	p.Queues[name], err = pubsub.NewTopic(ctx, fmt.Sprintf("%s-queue", name), &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(resourceLabels),
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return err
	}

	p.QueueSubscriptions[name], err = pubsub.NewSubscription(ctx, fmt.Sprintf("%s-nitricqueue", name), &pubsub.SubscriptionArgs{
		Topic:  p.Queues[name].Name,
		Labels: pulumi.ToStringMap(resourceLabels),
		ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
			Ttl: pulumi.String(""),
		},
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return err
	}

	return nil
}
