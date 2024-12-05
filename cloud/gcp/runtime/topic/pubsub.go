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

package topic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	grpccodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	tasks "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"cloud.google.com/go/pubsub"
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/propagator"
	"go.opentelemetry.io/otel/propagation"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"

	ifaces_cloudtasks "github.com/nitrictech/nitric/cloud/gcp/ifaces/cloudtasks"
	ifaces_pubsub "github.com/nitrictech/nitric/cloud/gcp/ifaces/pubsub"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/resource"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	"github.com/nitrictech/nitric/core/pkg/help"
	topicpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

type PubsubEventService struct {
	resource.GcpResourceResolver
	client      ifaces_pubsub.PubsubClient
	tasksClient ifaces_cloudtasks.CloudtasksClient
	cacheLock   sync.Mutex
	cache       map[string]ifaces_pubsub.Topic
}

var _ topicpb.TopicsServer = &PubsubEventService{}

func (s *PubsubEventService) getPubsubTopicFromName(topic string) (ifaces_pubsub.Topic, error) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	if s.cache == nil {
		topics := s.client.Topics(context.Background())
		s.cache = make(map[string]ifaces_pubsub.Topic)
		for {
			t, err := topics.Next()
			if errors.Is(err, iterator.Done) {
				break
			}
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding topic: %s; %w", topic, err)
			}

			labels, err := t.Labels(context.TODO())
			if err != nil {
				return nil, fmt.Errorf("an error occurred finding topic labels: %s; %w", topic, err)
			}

			resType, hasType := labels[tags.GetResourceTypeKey(env.GetNitricStackID())]

			if name, ok := labels[tags.GetResourceNameKey(env.GetNitricStackID())]; ok && hasType && name == topic && resType == "topic" {
				s.cache[name] = t
			}
		}
	}

	if topic, ok := s.cache[topic]; ok {
		return topic, nil
	}

	return nil, fmt.Errorf("topic not found")
}

type httpPubsubMessage struct {
	Attributes map[string]string `json:"attributes"`
	Data       []byte            `json:"data"`
}

type httpPubsubMessages struct {
	Messages []httpPubsubMessage `json:"messages"`
}

func (s *PubsubEventService) publish(ctx context.Context, topic string, pubsubMsg *pubsub.Message) error {
	msg := ifaces_pubsub.AdaptPubsubMessage(pubsubMsg)
	pubsubTopic, err := s.getPubsubTopicFromName(topic)
	if err != nil {
		return err
	}

	_, err = pubsubTopic.Publish(ctx, msg).Get(ctx)

	return err
}

func (s *PubsubEventService) publishDelayed(ctx context.Context, topic string, delay time.Duration, pubsubMsg *pubsub.Message) error {
	delayTo := timestamppb.Now().AsTime().Add(delay)

	saEmail, err := s.GetServiceAccountEmail()
	if err != nil {
		return err
	}

	body := httpPubsubMessages{
		Messages: []httpPubsubMessage{{
			Attributes: pubsubMsg.Attributes,
			Data:       pubsubMsg.Data,
		}},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	pubsubTopic, err := s.getPubsubTopicFromName(topic)
	if err != nil {
		return err
	}

	// Delay the message publishing
	_, err = s.tasksClient.CreateTask(ctx, &tasks.CreateTaskRequest{
		Parent: env.DELAY_QUEUE_NAME.String(),
		Task: &tasks.Task{
			MessageType: &tasks.Task_HttpRequest{
				HttpRequest: &tasks.HttpRequest{
					AuthorizationHeader: &tasks.HttpRequest_OauthToken{
						OauthToken: &tasks.OAuthToken{
							ServiceAccountEmail: saEmail,
						},
					},
					HttpMethod: tasks.HttpMethod_POST,
					Url:        fmt.Sprintf("https://pubsub.googleapis.com/v1/%s:publish", pubsubTopic.String()),
					Body:       jsonBody,
				},
			},
			// schedule for the future
			ScheduleTime: timestamppb.New(delayTo),
		},
	})

	return err
}

func (s *PubsubEventService) Publish(ctx context.Context, req *topicpb.TopicPublishRequest) (*topicpb.TopicPublishResponse, error) {
	delay := req.Delay.AsDuration()
	newErr := grpc_errors.ErrorsWithScope("PubsubEventService.Publish")

	messageBytes, err := proto.Marshal(req.Message)
	if err != nil {
		return nil, newErr(
			codes.Unknown,
			fmt.Sprintf("unable to serialize message. %s", help.BugInNitricHelpText()),
			err,
		)
	}

	// allows ctx to include the name of the source topic.
	attributes := propagation.MapCarrier{
		"x-nitric-topic": req.TopicName,
	}

	propagator.CloudTraceFormatPropagator{}.Inject(ctx, attributes)

	pubsubMsg := &pubsub.Message{
		Attributes: attributes,
		Data:       messageBytes,
	}

	if delay > 0 {
		err = s.publishDelayed(ctx, req.TopicName, delay, pubsubMsg)
	} else {
		err = s.publish(ctx, req.TopicName, pubsubMsg)
	}

	if err != nil {
		errStatus, _ := status.FromError(err)
		if errStatus.Code() == grpccodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this topic?", err,
			)
		}

		return nil, newErr(
			codes.Internal,
			fmt.Sprintf("error publishing message: %s", err.Error()),
			err,
		)
	}

	return &topicpb.TopicPublishResponse{}, nil
}

func New(provider resource.GcpResourceResolver) (topicpb.TopicsServer, error) {
	ctx := context.Background()

	credentials, err := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if err != nil {
		return nil, fmt.Errorf("GCP credentials error: %w", err)
	}

	client, err := pubsub.NewClient(ctx, credentials.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("pubsub client error: %w", err)
	}

	tasksClient, err := cloudtasks.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("cloudtasks client error: %w", err)
	}

	return &PubsubEventService{
		GcpResourceResolver: provider,
		client:              ifaces_pubsub.AdaptPubsubClient(client),
		tasksClient:         tasksClient,
	}, nil
}

func NewWithClient(provider resource.GcpResourceResolver, client ifaces_pubsub.PubsubClient, tasksClient ifaces_cloudtasks.CloudtasksClient) (topicpb.TopicsServer, error) {
	return &PubsubEventService{
		GcpResourceResolver: provider,
		client:              client,
		tasksClient:         tasksClient,
	}, nil
}
