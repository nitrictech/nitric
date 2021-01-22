package ifaces

import (
	"context"
	"time"
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
	ID() string
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
