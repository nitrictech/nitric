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

package events_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/utils"
)

type LocalEventService struct {
	events.UnimplementedeventsPlugin
	subscriptions map[string][]string
	client        LocalHttpeventsClient
}

// Interface for methods utilised by
// The local pubsub plugin for http events
type LocalHttpeventsClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Publish a message to a given topic
func (s *LocalEventService) Publish(topic string, event *events.NitricEvent) error {
	newErr := errors.ErrorsWithScope(
		"LocalEventService.Publish",
		map[string]interface{}{
			"topic": topic,
			"event": event,
		},
	)

	requestId := event.ID
	payloadType := event.PayloadType
	payload := event.Payload

	marshaledPayload, err := json.Marshal(payload)
	contentType := http.DetectContentType(marshaledPayload)

	if err != nil {
		return newErr(
			codes.Internal,
			"error marshalling event payload",
			err,
		)
	}

	if targets, ok := s.subscriptions[topic]; ok {
		fmt.Println(fmt.Sprintf("Publishing event to: %s", targets))
		for _, target := range targets {
			httpRequest, _ := http.NewRequest("POST", target, bytes.NewReader(marshaledPayload))

			httpRequest.Header.Add("Content-Type", contentType)
			httpRequest.Header.Add("x-nitric-request-id", requestId)
			httpRequest.Header.Add("x-nitric-source", topic)
			httpRequest.Header.Add("x-nitric-source-type", triggers.TriggerType_Subscription.String())
			httpRequest.Header.Add("x-nitric-payload-type", payloadType)

			// Call the target
			res, err := s.client.Do(httpRequest)
			if err != nil {
				fmt.Println(err)
				return newErr(
					codes.Internal,
					"unable to send message",
					err,
				)
			}
			if res.StatusCode < 200 || res.StatusCode >= 300 {
				buf := new(bytes.Buffer)
				_, _ = buf.ReadFrom(res.Body)
				body := buf.String()
				// TODO: Think about dead-letter functionality for these failed subscribers.
				// Just log failed delivery of events, since a single receiver failing to process an event wouldn't be an error in a cloud service.
				fmt.Println(fmt.Sprintf("Failed to publish event to %s\nStatus Code: %v\n%s", target, res.StatusCode, body))
			}
		}
	} else {
		return newErr(
			codes.NotFound,
			"unable to find subscriber for topic",
			nil,
		)
	}

	return nil
}

// Get a list of available topics
func (s *LocalEventService) ListTopics() ([]string, error) {
	keys := []string{}

	for key := range s.subscriptions {
		keys = append(keys, key)
	}

	return keys, nil
}

// Create new Dev EventService
func New() (events.EventService, error) {
	localSubscriptions := utils.GetEnv("LOCAL_SUBSCRIPTIONS", "{}")

	tmpSubs := make(map[string][]string)
	subs := make(map[string][]string)

	json.Unmarshal([]byte(localSubscriptions), &tmpSubs)

	for key, val := range tmpSubs {
		subs[strings.ToLower(key)] = val
	}

	return &LocalEventService{
		subscriptions: subs,
		client:        http.DefaultClient,
	}, nil
}

func NewWithClientAndSubs(client LocalHttpeventsClient, subs map[string][]string) (events.EventService, error) {
	return &LocalEventService{
		subscriptions: subs,
		client:        client,
	}, nil
}
