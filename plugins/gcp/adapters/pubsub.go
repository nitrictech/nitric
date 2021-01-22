package adapters

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
)

// AdaptClient adapts a pubsub.Client so that it satisfies the Client
// interface.
func AdaptPubsubClient(c *pubsub.Client) ifaces.PubsubClient {
	return pubsubClient{c}
}

// AdaptMessage adapts a pubsub.Message so that it satisfies the Message
// interface.
func AdaptPubsubMessage(msg *pubsub.Message) ifaces.Message {
	return message{msg}
}

type (
	pubsubClient  struct{ *pubsub.Client }
	topic         struct{ *pubsub.Topic }
	subscription  struct{ *pubsub.Subscription }
	message       struct{ *pubsub.Message }
	publishResult struct{ *pubsub.PublishResult }
	topicIterator struct{ *pubsub.TopicIterator }
)

func (c pubsubClient) Topic(id string) ifaces.Topic {
	return topic{c.Client.Topic(id)}
}

func (c pubsubClient) Topics(ctx context.Context) ifaces.TopicIterator {
	return topicIterator{c.Client.Topics(ctx)}
}

func (c topicIterator) Next() (ifaces.Topic, error) {
	t, err := c.TopicIterator.Next()
	return topic{t}, err
}

// func (c client) CreateSubscription(ctx context.Context, id string, cfg SubscriptionConfig) (Subscription, error) {
// 	s, err := c.Client.CreateSubscription(ctx, id, cfg.toPS())
// 	if err != nil {
// 		return nil, err
// 	}
// 	return subscription{s}, nil
// }

func (t topic) String() string {
	return t.Topic.String()
}

func (t topic) Publish(ctx context.Context, msg ifaces.Message) ifaces.PublishResult {
	return publishResult{t.Topic.Publish(ctx, msg.(message).Message)}
}

func (t topic) Exists(ctx context.Context) (bool, error) {
	return t.Topic.Exists(ctx)
}

func (t topic) ID() string {
	return t.Topic.ID()
}

// func (s subscription) Exists(ctx context.Context) (bool, error) {
// 	return s.Subscription.Exists(ctx)
// }

// func (s subscription) Receive(ctx context.Context, f func(ctx context.Context, msg Message)) error {
// 	return s.Subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
// 		f(ctx, AdaptMessage(msg))
// 	})
// }

// func (s subscription) Delete(ctx context.Context) error {
// 	return s.Subscription.Delete(ctx)
// }

func (m message) ID() string {
	return m.Message.ID
}

func (m message) Data() []byte {
	return m.Message.Data
}

func (m message) Attributes() map[string]string {
	return m.Message.Attributes
}

func (m message) PublishTime() time.Time {
	return m.Message.PublishTime
}

func (r publishResult) Get(ctx context.Context) (serverID string, err error) {
	return r.PublishResult.Get(ctx)
}

// func (cfg SubscriptionConfig) toPS() pubsub.SubscriptionConfig {
// 	return pubsub.SubscriptionConfig{
// 		Topic:               cfg.Topic.(topic).Topic,
// 		PushConfig:          cfg.PushConfig,
// 		AckDeadline:         cfg.AckDeadline,
// 		RetainAckedMessages: cfg.RetainAckedMessages,
// 		RetentionDuration:   cfg.RetentionDuration,
// 		Labels:              cfg.Labels,
// 	}
// }
