package pubsub_queue_plugin

import (
	"context"
	"encoding/json"

	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

type TasksPlugin struct {
	sdk.UnimplementedQueuePlugin
	client *pubsub.Client
}

func (s *TasksPlugin) Push(queue string, events []*sdk.NitricEvent) (*sdk.PushResponse, error) {
	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	topic := s.client.Topic(queue)
	ctx := context.TODO()

	// Push once we've published all messages to the client
	// TODO: We may want to revisit this, and chunk up our publishing in a way that makes more
	// sense...
	topic.PublishSettings.CountThreshold = len(events)
	results := make([]*pubsub.PublishResult, 0)
	failedMessages := make([]*sdk.NitricEvent, 0)
	publishedMessages := make([]*sdk.NitricEvent, 0)

	for _, evt := range events {
		if eventBytes, err := json.Marshal(evt); err == nil {
			msg := &pubsub.Message{
				Data: eventBytes,
			}

			results = append(results, topic.Publish(ctx, msg))
			publishedMessages = append(publishedMessages, evt)
		} else {
			// TODO: Append a publishing error to the error results here...
			failedMessages = append(failedMessages, evt)
			// Reduce the publish threshold by one
			// as the client will never see this message
			topic.PublishSettings.CountThreshold--
		}
	}

	for idx, result := range results {
		// Iterate over the results to check for successful publishing...
		if _, err := result.Get(ctx); err != nil {
			// Add this to our failures list in our results...
			failedMessages = append(failedMessages, publishedMessages[idx])
		}
	}

	return &sdk.PushResponse{
		FailedMessages: failedMessages,
	}, nil
}
