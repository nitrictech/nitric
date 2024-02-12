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
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid/eventgridapi"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	topicpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

type EventGridEventService struct {
	client   eventgridapi.BaseClientAPI
	provider resource.AzProvider
}

var _ topicpb.TopicsServer = &EventGridEventService{}

func (s *EventGridEventService) nitricEventToAzureEvent(topic string, payload *topicpb.Message) (*eventgrid.Event, error) {
	dataVersion := "1.0"
	eventType := "nitric"

	uid := uuid.New()
	id := uid.String()

	msgBytes, err := proto.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshalling event payload, this is a bug in nitric. %w", err)
	}

	return &eventgrid.Event{
		ID:          &id,
		Data:        msgBytes,
		EventType:   &eventType,
		Subject:     &topic,
		EventTime:   &date.Time{Time: time.Now()},
		DataVersion: &dataVersion,
	}, nil
}

func (s *EventGridEventService) Publish(ctx context.Context, req *topicpb.TopicPublishRequest) (*topicpb.TopicPublishResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("EventGrid.Publish")

	if req.Delay.AsDuration() != time.Duration(0) {
		return nil, newErr(codes.Unimplemented, "delayed messages with eventgrid are unsupported", nil)
	}

	topics, err := s.provider.GetResources(ctx, resource.AzResource_Topic)
	if err != nil {
		return nil, newErr(
			codes.NotFound,
			fmt.Sprintf("unable to find topic %s", req.TopicName),
			err,
		)
	}

	t, ok := topics[req.TopicName]
	if !ok {
		return nil, newErr(
			codes.NotFound,
			fmt.Sprintf("topic %s does not exist", req.TopicName),
			err,
		)
	}

	topicHostName := fmt.Sprintf("%s.%s-1.eventgrid.azure.net", t.Name, t.Location)

	eventToPublish, err := s.nitricEventToAzureEvent(topicHostName, req.Message)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error marshalling event",
			err,
		)
	}

	result, err := s.client.PublishEvents(ctx, topicHostName, []eventgrid.Event{*eventToPublish})
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error publishing event",
			err,
		)
	}

	if result.StatusCode < 200 || result.StatusCode >= 300 {
		return nil, newErr(
			codes.Internal,
			"returned non 200 status code",
			fmt.Errorf(result.Status),
		)
	}

	return &topicpb.TopicPublishResponse{}, nil
}

func New(provider resource.AzProvider) (*EventGridEventService, error) {
	// Get the event grid token, using the event grid resource endpoint
	spt, err := provider.ServicePrincipalToken("https://eventgrid.azure.net")
	if err != nil {
		return nil, fmt.Errorf("error authenticating event grid client: %w", err)
	}

	client := eventgrid.New()
	client.Authorizer = autorest.NewBearerAuthorizer(spt)

	return &EventGridEventService{
		provider: provider,
		client:   client,
	}, nil
}

func NewWithClient(provider resource.AzProvider, client eventgridapi.BaseClientAPI) (*EventGridEventService, error) {
	return &EventGridEventService{
		client:   client,
		provider: provider,
	}, nil
}
