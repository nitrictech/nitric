package ifaces_pubsub

import (
	"context"
	"time"

	"github.com/googleapis/gax-go/v2"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

type PubsubClient interface {
	// CreateTopic(ctx context.Context, topicID string) (Topic, error)
	Topic(id string) Topic
	Topics(ctx context.Context) TopicIterator
	// CreateSubscription(ctx context.Context, id string, cfg SubscriptionConfig) (Subscription, error)
	// Subscription(id string) Subscription
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
}

type SubscriptionIterator interface {
	Next() (Subscription, error)
}

type Subscription interface {
	ID() string
	String() string
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
