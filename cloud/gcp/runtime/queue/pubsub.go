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

package queue

import (
	"context"
	"errors"
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"
	"google.golang.org/grpc/codes"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"cloud.google.com/go/pubsub"
	pubsubbase "cloud.google.com/go/pubsub/apiv1"
	pubsubpb "cloud.google.com/go/pubsub/apiv1/pubsubpb"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	ifaces_pubsub "github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	"github.com/nitrictech/nitric/core/pkg/logger"

	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
)

type PubsubQueueService struct {
	queuespb.UnimplementedQueuesServer
	// queue.UnimplementedQueuePlugin
	client              ifaces_pubsub.PubsubClient
	newSubscriberClient func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error)
	projectId           string
	cache               map[string]ifaces_pubsub.Topic
}

// Retrieves the Nitric "Queue Topic" for the specified queue (PubSub Topic).
//
// This retrieves the default Nitric Queue for the Topic based on tagging conventions.
func (s *PubsubQueueService) getPubsubTopicFromName(queue string) (ifaces_pubsub.Topic, error) {
	if s.cache == nil {
		topics := s.client.Topics(context.Background())
		s.cache = make(map[string]ifaces_pubsub.Topic)
		stackID := env.GetNitricStackID()
		for {
			t, err := topics.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding queue: %s; %w", queue, err)
			}

			labels, err := t.Labels(context.TODO())
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding queue labels: %s; %w", queue, err)
			}

			resType, hasType := labels[tags.GetResourceTypeKey(stackID)]

			if name, ok := labels[tags.GetResourceNameKey(stackID)]; ok && name == queue && hasType && resType == "queue" {
				s.cache[name] = t
			}
		}
	}

	if t, ok := s.cache[queue]; ok {
		return t, nil
	}

	return nil, fmt.Errorf("queue not found")
}

// Retrieves the Nitric "Queue Subscription" for the specified queue (PubSub Topic).
//
// GCP PubSub requires a Subscription in order to Pull messages from a Topic.
// we use this behavior to implement queues.
//
// This retrieves the default Nitric Pull subscription for the Topic base on convention.
func (s *PubsubQueueService) getQueueSubscription(ctx context.Context, queueName string) (ifaces_pubsub.Subscription, error) {
	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	topic, err := s.getPubsubTopicFromName(queueName)
	if err != nil {
		return nil, err
	}

	subsIt := topic.Subscriptions(ctx)

	for {
		sub, err := subsIt.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve pull subscription for topic: %s\n%w", topic.ID(), err)
		}

		labels, err := sub.Labels(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve pull subscription labels for topic: %s\n%w", topic.ID(), err)
		}

		resourceType, hasType := labels[tags.GetResourceTypeKey(env.GetNitricStackID())]
		if name, ok := labels[tags.GetResourceNameKey(env.GetNitricStackID())]; hasType && ok && resourceType == string(resources.Queue) {
			if name == queueName {
				return sub, nil
			}
		}
	}

	return nil, fmt.Errorf("pull subscription not found, pull subscribers may not be configured for this topic")
}

func (s *PubsubQueueService) Enqueue(ctx context.Context, req *queuespb.QueueEnqueueRequest) (*queuespb.QueueEnqueueResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("PubsubQueueService.Enqueue")

	// We'll be using pubsub with pull subscribers to facilitate queue functionality
	topic, err := s.getPubsubTopicFromName(req.QueueName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"queue not found",
			err,
		)
	}

	// SendBatch once we've published all tasks to the client
	results := make([]ifaces_pubsub.PublishResult, 0)
	failedTasks := make([]*queuespb.FailedEnqueueMessage, 0)
	publishedTasks := make([]*queuespb.QueueMessage, 0)

	attributes := propagation.MapCarrier{}

	propagator.CloudTraceFormatPropagator{}.Inject(ctx, attributes)

	for _, task := range req.Messages {
		t := task
		if taskBytes, err := proto.Marshal(t); err == nil {
			msg := ifaces_pubsub.AdaptPubsubMessage(&pubsub.Message{
				Data:       taskBytes,
				Attributes: attributes,
			})

			results = append(results, topic.Publish(ctx, msg))
			publishedTasks = append(publishedTasks, t)
		} else {
			failedTasks = append(failedTasks, &queuespb.FailedEnqueueMessage{
				Message: t,
				Details: "Error unmarshalling message for queue",
			})
		}
	}

	for idx, result := range results {
		// Iterate over the results to check for successful publishing...
		if _, err := result.Get(ctx); err != nil {
			// Add this to our failures list in our results...
			failedTasks = append(failedTasks, &queuespb.FailedEnqueueMessage{
				Message: publishedTasks[idx],
				Details: err.Error(),
			})
		}
	}

	return &queuespb.QueueEnqueueResponse{
		FailedMessages: failedTasks,
	}, nil
}

func (s *PubsubQueueService) Dequeue(ctx context.Context, req *queuespb.QueueDequeueRequest) (*queuespb.QueueDequeueResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("PubsubQueueService.Dequeue")

	// Find the generic pull subscription for the provided topic (queue)
	queueSubscription, err := s.getQueueSubscription(ctx, req.QueueName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"could not find queue subscription",
			err,
		)
	}

	// Using base client, so that asynchronous message acknowledgement can take place without needing to keep messages
	// in a stateful service. Standard PubSub go library doesn't provide access to the 'acknowledge' ID of the messages
	// or an independent acknowledge function. It's only provided as a method on message objects.
	client, err := s.newSubscriberClient(ctx)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to create subscriber client",
			err,
		)
	}
	defer client.Close()

	// Retrieve the requested number of messages from the subscription (queue)
	pubsubRequest := pubsubpb.PullRequest{
		Subscription: queueSubscription.String(),
		MaxMessages:  req.GetDepth(),
	}
	res, err := client.Pull(ctx, &pubsubRequest)
	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == grpccodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this queue?", err)
		}

		return nil, newErr(
			codes.Internal,
			"failed to pull messages",
			err,
		)
	}

	// An empty list is returned from PubSub if no messages are available
	// we return our own empty list in turn.
	if len(res.ReceivedMessages) == 0 {
		return &queuespb.QueueDequeueResponse{
			Messages: []*queuespb.ReceivedMessage{},
		}, nil
	}

	// Convert the PubSub messages into Nitric tasks
	var tasks []*queuespb.ReceivedMessage
	for _, m := range res.ReceivedMessages {
		// var nitricTask queuespb.ReceivedTask
		var queueMessage queuespb.QueueMessage
		err := proto.Unmarshal(m.Message.Data, &queueMessage)
		if err != nil {
			// This item could be immediately requeued.
			// However, that risks the unprocessable items being reprocessed immediately,
			// causing a loop where the receiver frequently attempts to receive the same item.
			logger.Errorf("failed to deserialize queue item payload: %s", err.Error())
			continue
		}

		tasks = append(tasks, &queuespb.ReceivedMessage{
			Message: &queueMessage,
			LeaseId: m.AckId,
		})
	}

	return &queuespb.QueueDequeueResponse{
		Messages: tasks,
	}, nil
}

func (s *PubsubQueueService) Complete(ctx context.Context, req *queuespb.QueueCompleteRequest) (*queuespb.QueueCompleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("PubsubQueueService.Complete")

	// Find the generic pull subscription for the provided topic (queue)
	queueSubscription, err := s.getQueueSubscription(ctx, req.QueueName)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			"could not find queue subscription",
			err,
		)
	}

	// Using base client, so that asynchronous message acknowledgement can take place without needing to keep messages
	// in a stateful service. Standard PubSub go library is stateful and don't provide access to the acknowledge ID of
	// the messages or an independent acknowledge function. It's only provided as a method on message objects.
	client, err := s.newSubscriberClient(ctx)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"failed to create subscriberclient",
			err,
		)
	}
	defer client.Close()

	// Acknowledge the queue item, so it's removed from the queue
	pubsubRequest := pubsubpb.AcknowledgeRequest{
		Subscription: queueSubscription.String(),
		AckIds:       []string{req.LeaseId},
	}
	err = client.Acknowledge(ctx, &pubsubRequest)
	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == grpccodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to the queue?",
				err)
		}

		return nil, newErr(
			codes.Internal,
			"failed to de-queue task",
			err,
		)
	}

	return &queuespb.QueueCompleteResponse{}, nil
}

// adaptNewClient - Adapts the pubsubbase.NewSubscriberClient func to one that implements the SubscriberClient
// interface. This is used to enable substitution of the base pubsub client, primarily for mocking support.
func adaptNewClient(f func(context.Context, ...option.ClientOption) (*pubsubbase.SubscriberClient, error)) func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
	return func(c context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error) {
		return f(c, opts...)
	}
}

// New - Constructs a new GCP pubsub client with defaults
func New() (*PubsubQueueService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %w", credentialsError)
	}
	client, clientError := pubsub.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("pubsub client error: %w", clientError)
	}

	return &PubsubQueueService{
		client:              ifaces_pubsub.AdaptPubsubClient(client),
		newSubscriberClient: adaptNewClient(pubsubbase.NewSubscriberClient),
		projectId:           credentials.ProjectID,
	}, nil
}

func NewWithClient(client ifaces_pubsub.PubsubClient) *PubsubQueueService {
	return &PubsubQueueService{
		client:              client,
		newSubscriberClient: nil,
	}
}

// *pubsubbase.SubscriberClient
func NewWithClients(client ifaces_pubsub.PubsubClient, subscriberClientGenerator func(ctx context.Context, opts ...option.ClientOption) (ifaces_pubsub.SubscriberClient, error)) *PubsubQueueService {
	return &PubsubQueueService{
		client:              client,
		newSubscriberClient: subscriberClientGenerator,
	}
}
