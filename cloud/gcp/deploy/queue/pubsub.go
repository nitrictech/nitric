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

package queue

import (
	"fmt"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PubSubTopic struct {
	pulumi.ResourceState

	Name         string
	PubSub       *pubsub.Topic
	Subscription *pubsub.Subscription
}

type PubSubTopicArgs struct {
	Location string
	StackID  string
	Queue    *v1.Queue
}

func NewPubSubTopic(ctx *pulumi.Context, name string, args *PubSubTopicArgs, opts ...pulumi.ResourceOption) (*PubSubTopic, error) {
	res := &PubSubTopic{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:queue:GCPPubSubTopic", name, res, opts...)
	if err != nil {
		return nil, err
	}

	topicTags := common.Tags(args.StackID, name)
	topicTags["x-nitric-type"] = "queue"

	res.PubSub, err = pubsub.NewTopic(ctx, name, &pubsub.TopicArgs{
		Labels: pulumi.ToStringMap(topicTags),
	})
	if err != nil {
		return nil, err
	}

	subscriptionTags := common.Tags(args.StackID, name)
	subscriptionTags["x-nitric-type"] = "queue-subscription"

	res.Subscription, err = pubsub.NewSubscription(ctx, fmt.Sprintf("%s-nitricqueue", name), &pubsub.SubscriptionArgs{
		Topic:  res.PubSub.Name,
		Labels: pulumi.ToStringMap(subscriptionTags),
		ExpirationPolicy: &pubsub.SubscriptionExpirationPolicyArgs{
			Ttl: pulumi.String(""),
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
