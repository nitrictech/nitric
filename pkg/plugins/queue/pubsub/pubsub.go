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

package pubsub_queue_service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nitric-dev/membrane/pkg/ifaces/pubsub"

	"cloud.google.com/go/pubsub"
	pubsubbase "cloud.google.com/go/pubsub/apiv1"
	"github.com/nitric-dev/membrane/pkg/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"
)

type PubsubQueueService struct {
	sdk.UnimplementedQueuePlugin
	client              ifaces_pubsub.PubsubClient
	newSubscriberClient func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error)
	projectId           string
	messages            []*pubsub.Message
}

// TODO: clearly document the reason for this subscription.
// Get the default Nitric Queue Subscription name for a given queue name.
func generateQueueSubscription(queue string) string {
	return fmt.Sprintf("%s-nitricqueue", queue)
}

func (s *PubsubQueueService) Send(queue string, task sdk.NitricTask) error {
	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	ctx := context.TODO()
	topic := s.client.Topic(queue)

	if exists, err := topic.Exists(ctx); !exists || err != nil {
		return fmt.Errorf("Queue: %s does not exist", queue)
	}

	if taskBytes, err := json.Marshal(task); err == nil {
		msg := ifaces_pubsub.AdaptPubsubMessage(&pubsub.Message{
			Data: taskBytes,
		})

		result := topic.Publish(ctx, msg)

		if _, err := result.Get(ctx); err != nil {
			return fmt.Errorf("Error getting result: %v", err)
		}
	} else {
		return fmt.Errorf("Error marshalling task: %v", err)
	}

	return nil
}

func (s *PubsubQueueService) SendBatch(queue string, tasks []sdk.NitricTask) (*sdk.SendBatchResponse, error) {
	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	ctx := context.TODO()
	topic := s.client.Topic(queue)

	if exists, err := topic.Exists(ctx); !exists || err != nil {
		return nil, fmt.Errorf("Queue: %s does not exist", queue)
	}

	// SendBatch once we've published all tasks to the client
	// TODO: We may want to revisit this, and chunk up our publishing in a way that makes more sense...
	results := make([]ifaces_pubsub.PublishResult, 0)
	failedTasks := make([]*sdk.FailedTask, 0)
	publishedTasks := make([]sdk.NitricTask, 0)

	for _, task := range tasks {
		if taskBytes, err := json.Marshal(task); err == nil {
			msg := ifaces_pubsub.AdaptPubsubMessage(&pubsub.Message{
				Data: taskBytes,
			})

			results = append(results, topic.Publish(ctx, msg))
			publishedTasks = append(publishedTasks, task)
		} else {
			failedTasks = append(failedTasks, &sdk.FailedTask{
				Task:    &task,
				Message: "Error unmarshalling message for queue",
			})
		}
	}

	for idx, result := range results {
		// Iterate over the results to check for successful publishing...
		if _, err := result.Get(ctx); err != nil {
			// Add this to our failures list in our results...
			failedTasks = append(failedTasks, &sdk.FailedTask{
				Task:    &publishedTasks[idx],
				Message: err.Error(),
			})
		}
	}

	return &sdk.SendBatchResponse{
		FailedTasks: failedTasks,
	}, nil
}

// Retrieves the Nitric "Queue Subscription" for the specified queue (PubSub Topic).
//
// GCP PubSub requires a Subscription in order to Pull messages from a Topic.
// we use this behavior to emulate a queue.
//
// This retrieves the default Nitric Pull subscription for the Topic base on convention.
func (s *PubsubQueueService) getQueueSubscription(queue string) (ifaces_pubsub.Subscription, error) {
	ctx := context.Background()

	topic := s.client.Topic(queue)
	subsIt := topic.Subscriptions(ctx)

	for {
		sub, err := subsIt.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve pull subscription for topic: %s\n%s", topic.ID(), err)
		}
		queueSubName := generateQueueSubscription(queue)
		if sub.ID() == queueSubName {
			return sub, nil
		}
	}

	return nil, fmt.Errorf("pull subscription not found, pull subscribers may not be configured for this topic")
}

// Receives a collection of tasks off a given queue.
func (s *PubsubQueueService) Receive(options sdk.ReceiveOptions) ([]sdk.NitricTask, error) {
	err := options.Validate()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	// Find the generic pull subscription for the provided topic (queue)
	queueSubscription, err := s.getQueueSubscription(options.QueueName)
	if err != nil {
		return nil, err
	}

	// Using base client, so that asynchronous message acknowledgement can take place without needing to keep messages
	// in a stateful service. Standard PubSub go library do not provide access to the acknowledge ID of the messages or
	// an independent acknowledge function. It's only provided as a method on message objects.
	client, err := s.newSubscriberClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create pubsub client.\n%s", err)
	}
	defer client.Close()

	// Retrieve the requested number of messages from the subscription (queue)
	req := pubsubpb.PullRequest{
		Subscription: queueSubscription.String(),
		MaxMessages:  int32(*options.Depth),
	}
	res, err := client.Pull(ctx, &req)
	if err != nil {
		// TODO: catch standard grpc errors, like NotFound.
		return nil, fmt.Errorf("failed to pull pubsub messages.\n%s", err)
	}

	// An empty list is returned from PubSub if no messages are available
	// we return our own empty list in turn.
	if len(res.ReceivedMessages) == 0 {
		return []sdk.NitricTask{}, nil
	}

	// Convert the PubSub messages into Nitric tasks
	var tasks []sdk.NitricTask
	for _, m := range res.ReceivedMessages {
		var nitricTask sdk.NitricTask
		err := json.Unmarshal(m.Message.Data, &nitricTask)
		if err != nil {
			// TODO: append error to error list and Nack the message.
			continue
		}

		tasks = append(tasks, sdk.NitricTask{
			ID:          nitricTask.ID,
			Payload:     nitricTask.Payload,
			PayloadType: nitricTask.PayloadType,
			LeaseID:     nitricTask.LeaseID,
		})
	}

	return tasks, nil
}

// Completes a previously popped queue item
func (s *PubsubQueueService) Complete(queue string, leaseId string) error {
	ctx := context.Background()

	// Find the generic pull subscription for the provided topic (queue)
	queueSubscription, err := s.getQueueSubscription(queue)
	if err != nil {
		return err
	}

	// Using base client, so that asynchronous message acknowledgement can take place without needing to keep messages
	// in a stateful service. Standard PubSub go library is stateful and don't provide access to the acknowledge ID of
	// the messages or an independent acknowledge function. It's only provided as a method on message objects.
	client, err := s.newSubscriberClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to create pubsub client.\n%s", err)
	}
	defer client.Close()

	// Acknowledge the queue item so it's removed from the queue
	req := pubsubpb.AcknowledgeRequest{
		Subscription: queueSubscription.String(),
		AckIds:       []string{leaseId},
	}
	err = client.Acknowledge(ctx, &req)
	if err != nil {
		// TODO: catch standard grpc errors, like NotFound.
		return fmt.Errorf("failed to complete queue item.\n%s", err)
	}

	return nil
}

// adaptNewClient - Adapts the pubsubbase.NewSubscriberClient func to one that implements the SubscriberClient
// interface. This is used to enable substitution of the base pubsub client, primarily for mocking support.
func adaptNewClient(f func(context.Context, ...option.ClientOption) (*pubsubbase.SubscriberClient, error)) func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
	return func(c context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
		return f(c, opts...)
	}
}

// New - Constructs a new GCP pubsub client with defaults
func New() (sdk.QueueService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}
	client, clientError := pubsub.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("pubsub client error: %v", clientError)
	}

	return &PubsubQueueService{
		client: ifaces_pubsub.AdaptPubsubClient(client),
		// TODO: replace this with a better mechanism for mocking the client.
		newSubscriberClient: adaptNewClient(pubsubbase.NewSubscriberClient),
		projectId:           credentials.ProjectID,
	}, nil
}

func NewWithClient(client ifaces_pubsub.PubsubClient) sdk.QueueService {
	return &PubsubQueueService{
		client:              client,
		newSubscriberClient: nil,
	}
}

//*pubsubbase.SubscriberClient
func NewWithClients(client ifaces_pubsub.PubsubClient, subscriberClientGenerator func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error)) sdk.QueueService {
	return &PubsubQueueService{
		client:              client,
		newSubscriberClient: subscriberClientGenerator,
	}
}
