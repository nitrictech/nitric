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

package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	"github.com/nitrictech/nitric/cloud/gcp/runtime/core"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/core/pkg/plugins/events"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

type PubsubEventService struct {
	events.UnimplementedeventsPlugin
	core.GcpProvider
	client      ifaces_pubsub.PubsubClient
	tasksClient ifaces_cloudtasks.CloudtasksClient
}

func (s *PubsubEventService) ListTopics(ctx context.Context) ([]string, error) {
	newErr := errors.ErrorsWithScope("PubsubEventService.ListTopics", nil)
	iter := s.client.Topics(ctx)

	var topics []string
	for topic, err := iter.Next(); !errors.Is(err, iterator.Done); topic, err = iter.Next() {
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error retrieving topics",
				err,
			)
		}

		topics = append(topics, topic.ID())
	}

	return topics, nil
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
	pubsubTopic := s.client.Topic(topic)

	_, err := pubsubTopic.Publish(ctx, msg).Get(ctx)
	return err
}

func (s *PubsubEventService) publishDelayed(ctx context.Context, topic string, delay int, pubsubMsg *pubsub.Message) error {
	saEmail, err := s.GetServiceAccountEmail()
	if err != nil {
		return err
	}

	projectId, err := s.GetProjectID()
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

	// Delay the message publishing
	_, err = s.tasksClient.CreateTask(ctx, &tasks.CreateTaskRequest{
		Parent: utils.GetEnv("DELAY_QUEUE_NAME", ""),
		Task: &tasks.Task{
			MessageType: &tasks.Task_HttpRequest{
				HttpRequest: &tasks.HttpRequest{
					AuthorizationHeader: &tasks.HttpRequest_OauthToken{
						OauthToken: &tasks.OAuthToken{
							ServiceAccountEmail: saEmail,
						},
					},
					HttpMethod: tasks.HttpMethod_POST,
					Url:        fmt.Sprintf("https://pubsub.googleapis.com/v1/projects/%s/topics/%s:publish", projectId, topic),
					// TODO: Add message body with attributes
					Body: jsonBody,
				},
			},
			// schedule for the future
			ScheduleTime: timestamppb.New(timestamppb.Now().AsTime().Add(time.Duration(delay) * time.Second)),
		},
	})

	return err
}

func (s *PubsubEventService) Publish(ctx context.Context, topic string, delay int, event *events.NitricEvent) error {
	newErr := errors.ErrorsWithScope(
		"PubsubEventService.Publish",
		map[string]interface{}{
			"topic": topic,
			"event": event,
		},
	)

	eventBytes, err := json.Marshal(event)
	if err != nil {
		return newErr(
			codes.Internal,
			"error marshalling event payload",
			err,
		)
	}

	attributes := propagation.MapCarrier{
		"x-nitric-topic": topic,
	}

	propagator.CloudTraceFormatPropagator{}.Inject(ctx, attributes)

	pubsubMsg := &pubsub.Message{
		Attributes: attributes,
		Data:       eventBytes,
	}

	if delay > 0 {
		err = s.publishDelayed(ctx, topic, delay, pubsubMsg)
	} else {
		err = s.publish(ctx, topic, pubsubMsg)
	}

	if err != nil {
		return newErr(
			codes.Internal,
			fmt.Sprintf("error publishing message: %s", err.Error()),
			err,
		)
	}

	return nil
}

func New(provider core.GcpProvider) (events.EventService, error) {
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
		GcpProvider: provider,
		client:      ifaces_pubsub.AdaptPubsubClient(client),
		tasksClient: tasksClient,
	}, nil
}

func NewWithClient(provider core.GcpProvider, client ifaces_pubsub.PubsubClient, tasksClient ifaces_cloudtasks.CloudtasksClient) (events.EventService, error) {
	return &PubsubEventService{
		GcpProvider: provider,
		client:      client,
		tasksClient: tasksClient,
	}, nil
}
