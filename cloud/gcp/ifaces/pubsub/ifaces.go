// Copyright 2021 Nitric Pty Ltd.
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

package ifaces_pubsub

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/googleapis/gax-go/v2"
)

type PubsubClient interface {
	Topic(id string) Topic
	Topics(ctx context.Context) TopicIterator
}

type TopicIterator interface {
	Next() (Topic, error)
}

type Topic interface {
	String() string
	Publish(ctx context.Context, msg Message) PublishResult
	Exists(ctx context.Context) (bool, error)
	Subscriptions(ctx context.Context) SubscriptionIterator
	ID() string
	Labels(ctx context.Context) (map[string]string, error)
}

type SubscriptionIterator interface {
	Next() (Subscription, error)
}

type Subscription interface {
	ID() string
	String() string
	Labels(ctx context.Context) (map[string]string, error)
}

type Message interface {
	ID() string
	Data() []byte
	Attributes() map[string]string
	PublishTime() time.Time
	Ack()
	Nack()
}

type PublishResult interface {
	Get(ctx context.Context) (serverID string, err error)
}

type SubscriberClient interface {
	Close() error
	Pull(ctx context.Context, req *pubsubpb.PullRequest, opts ...gax.CallOption) (*pubsubpb.PullResponse, error)
	Acknowledge(ctx context.Context, req *pubsubpb.AcknowledgeRequest, opts ...gax.CallOption) error
}
