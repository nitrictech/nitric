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

package eventgrid_service

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid"
	"github.com/Azure/azure-sdk-for-go/services/eventgrid/2018-01-01/eventgrid/eventgridapi"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/date"

	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	"github.com/nitrictech/nitric/pkg/providers/azure/core"
)

type EventGridEventService struct {
	events.UnimplementedeventsPlugin
	client   eventgridapi.BaseClientAPI
	provider core.AzProvider
}

func (s *EventGridEventService) ListTopics() ([]string, error) {
	newErr := errors.ErrorsWithScope(
		"EventGrid.ListTopics",
		map[string]interface{}{
			"list": "topics",
		},
	)

	topics, err := s.provider.GetResources(core.AzResource_Topic)

	if err != nil {
		return nil, newErr(codes.Internal, "unable to retrieve topic list", err)
	}

	topicsList := make([]string, 0, len(topics))
	for t := range topics {
		topicsList = append(topicsList, t)
	}

	return topicsList, nil
}

func (s *EventGridEventService) nitricEventsToAzureEvents(topic string, events []*events.NitricEvent) ([]eventgrid.Event, error) {
	var azureEvents []eventgrid.Event
	for _, event := range events {
		dataVersion := "1.0"
		azureEvents = append(azureEvents, eventgrid.Event{
			ID:          &event.ID,
			Data:        event.Payload,
			EventType:   &event.PayloadType,
			Subject:     &topic,
			EventTime:   &date.Time{time.Now()},
			DataVersion: &dataVersion,
		})
	}

	return azureEvents, nil
}

func (s *EventGridEventService) Publish(topic string, event *events.NitricEvent) error {
	newErr := errors.ErrorsWithScope(
		"EventGrid.Publish",
		map[string]interface{}{
			"topic": topic,
		},
	)

	topics, err := s.provider.GetResources(core.AzResource_Topic)
	if err != nil {
		return newErr(
			codes.NotFound,
			fmt.Sprintf("unable to find topic %s: %v", topic, err),
			err,
		)
	}

	t, ok := topics[topic]
	if !ok {
		return newErr(
			codes.NotFound,
			fmt.Sprintf("topic %s does not exist", topic),
			err,
		)
	}

	// TODO: Determine correctness of availability zone in endpoint hostname
	topicHostName := fmt.Sprintf("%s.%s-1.eventgrid.azure.net", t.Name, t.Location)
	fmt.Println("topic host name is", topicHostName)

	eventToPublish, err := s.nitricEventsToAzureEvents(topicHostName, []*events.NitricEvent{event})
	if err != nil {
		return newErr(
			codes.Internal,
			"error marshalling event",
			err,
		)
	}

	result, err := s.client.PublishEvents(context.TODO(), topicHostName, eventToPublish)
	if err != nil {
		return newErr(
			codes.Internal,
			"error publishing event",
			err,
		)
	}

	if result.StatusCode < 200 || result.StatusCode >= 300 {
		return newErr(
			codes.Internal,
			"returned non 200 status code",
			fmt.Errorf(result.Status),
		)
	}

	return nil
}

func New(provider core.AzProvider) (events.EventService, error) {
	//Get the event grid token, using the event grid resource endpoint
	spt, err := provider.ServicePrincipalToken("https://eventgrid.azure.net")
	if err != nil {
		return nil, fmt.Errorf("error authenticating event grid client: %v", err.Error())
	}

	client := eventgrid.New()
	client.Authorizer = autorest.NewBearerAuthorizer(spt)

	return &EventGridEventService{
		provider: provider,
		client:   client,
	}, nil
}

func NewWithClient(provider core.AzProvider, client eventgridapi.BaseClientAPI) (events.EventService, error) {
	return &EventGridEventService{
		client:   client,
		provider: provider,
	}, nil
}
