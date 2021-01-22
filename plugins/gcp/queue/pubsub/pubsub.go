package pubsub_queue_plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/gcp/adapters"
	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"golang.org/x/oauth2/google"
)

type PubsubPlugin struct {
	sdk.UnimplementedQueuePlugin
	client ifaces.PubsubClient
}

func (s *PubsubPlugin) Push(queue string, events []*sdk.NitricEvent) (*sdk.PushResponse, error) {
	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	ctx := context.TODO()
	topic := s.client.Topic(queue)

	if exists, err := topic.Exists(ctx); !exists || err != nil {
		return nil, fmt.Errorf("Queue: %s does not exist", queue)
	}

	// Push once we've published all messages to the client
	// TODO: We may want to revisit this, and chunk up our publishing in a way that makes more
	// sense...
	results := make([]ifaces.PublishResult, 0)
	failedMessages := make([]*sdk.NitricEvent, 0)
	publishedMessages := make([]*sdk.NitricEvent, 0)

	for _, evt := range events {
		if eventBytes, err := json.Marshal(evt); err == nil {
			msg := adapters.AdaptPubsubMessage(&pubsub.Message{
				Data: eventBytes,
			})

			results = append(results, topic.Publish(ctx, msg))
			publishedMessages = append(publishedMessages, evt)
		} else {
			failedMessages = append(failedMessages, evt)
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

// New - Constructs a new GCP pubsub client with defaults
func New() (sdk.QueuePlugin, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}

	client, clientError := pubsub.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("pubsub client error: %v", clientError)
	}

	return &PubsubPlugin{
		client: adapters.AdaptPubsubClient(client),
	}, nil
}

func NewWithClient(client ifaces.PubsubClient) sdk.QueuePlugin {
	return &PubsubPlugin{
		client: client,
	}
}
