package pubsub_plugin

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

type PubsubPlugin struct {
	sdk.UnimplementedEventingPlugin
	client *pubsub.Client
}

func (s *PubsubPlugin) GetTopics() ([]string, error) {
	iter := s.client.Topics(context.TODO())

	var topics []string
	for {
		topic, err := iter.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("Error retrieving topics %v", err)
		}

		topics = append(topics, topic.ID())
	}

	return topics, nil
}

func (s *PubsubPlugin) Publish(topic string, event *sdk.NitricEvent) error {
	ctx := context.TODO()

	eventBytes, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("Payload marshalling error: %v", err)
	}

	pubsubTopic := s.client.Topic(topic)

	msg := &pubsub.Message{
		Data: eventBytes,
	}

	if _, err := pubsubTopic.Publish(ctx, msg).Get(ctx); err != nil {
		return fmt.Errorf("Payload marshalling error: %v", err)
	}

	return nil
}

func New() (sdk.EventingPlugin, error) {
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
		client: client,
	}, nil
}

func NewWithClient(client *pubsub.Client) (sdk.EventingPlugin, error) {
	return &PubsubPlugin{
		client: client,
	}, nil
}
