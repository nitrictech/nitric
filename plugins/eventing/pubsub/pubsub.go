package pubsub_service

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	ifaces "github.com/nitric-dev/membrane/ifaces/pubsub"
	"github.com/nitric-dev/membrane/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

type PubsubEventService struct {
	sdk.UnimplementedEventingPlugin
	client ifaces.PubsubClient
}

func (s *PubsubEventService) ListTopics() ([]string, error) {
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

func (s *PubsubEventService) Publish(topic string, event *sdk.NitricEvent) error {
	ctx := context.TODO()

	eventBytes, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("Payload marshalling error: %v", err)
	}

	pubsubTopic := s.client.Topic(topic)

	msg := ifaces.AdaptPubsubMessage(&pubsub.Message{
		Attributes: map[string]string{
			"x-nitric-topic": topic,
		},
		Data: eventBytes,
	})

	if _, err := pubsubTopic.Publish(ctx, msg).Get(ctx); err != nil {
		return fmt.Errorf("Payload marshalling error: %v", err)
	}

	return nil
}

func New() (sdk.EventService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}

	client, clientError := pubsub.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("pubsub client error: %v", clientError)
	}

	return &PubsubEventService{
		client: ifaces.AdaptPubsubClient(client),
	}, nil
}

func NewWithClient(client ifaces.PubsubClient) (sdk.EventService, error) {
	return &PubsubEventService{
		client: client,
	}, nil
}
