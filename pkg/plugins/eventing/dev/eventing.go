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

package eventing_service

import (
	"bytes"
	"encoding/json"
	"fmt"
	triggers2 "github.com/nitric-dev/membrane/pkg/triggers"
	utils2 "github.com/nitric-dev/membrane/pkg/utils"
	"net/http"
	"strings"

	"github.com/nitric-dev/membrane/pkg/sdk"
)

type LocalEventService struct {
	sdk.UnimplementedEventingPlugin
	subscriptions map[string][]string
	client        LocalHttpEventingClient
}

// Interface for methods utilised by
// The local pubsub plugin for http eventing
type LocalHttpEventingClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Publish a message to a given topic
func (s *LocalEventService) Publish(topic string, event *sdk.NitricEvent) error {
	requestId := event.ID
	payloadType := event.PayloadType
	payload := event.Payload

	marshaledPayload, err := json.Marshal(payload)
	contentType := http.DetectContentType(marshaledPayload)

	if err != nil {
		return err
	}

	if targets, ok := s.subscriptions[topic]; ok {
		fmt.Println(fmt.Sprintf("Publishing event to: %s", targets))
		for _, target := range targets {
			httpRequest, _ := http.NewRequest("POST", target, bytes.NewReader(marshaledPayload))

			httpRequest.Header.Add("Content-Type", contentType)
			httpRequest.Header.Add("x-nitric-request-id", requestId)
			httpRequest.Header.Add("x-nitric-source", topic)
			httpRequest.Header.Add("x-nitric-source-type", triggers2.TriggerType_Subscription.String())
			httpRequest.Header.Add("x-nitric-payload-type", payloadType)

			// Call the target
			res, err := s.client.Do(httpRequest)
			if err != nil {
				fmt.Println(err)
				return err
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
		return fmt.Errorf("No subscription found for %s in %v", topic, s.subscriptions)
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
func New() (sdk.EventService, error) {
	localSubscriptions := utils2.GetEnv("LOCAL_SUBSCRIPTIONS", "{}")

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

func NewWithClientAndSubs(client LocalHttpEventingClient, subs map[string][]string) (sdk.EventService, error) {
	return &LocalEventService{
		subscriptions: subs,
		client:        client,
	}, nil
}
