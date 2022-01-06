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

package worker

import (
	"fmt"
	"net/http"

	"github.com/nitrictech/nitric/pkg/triggers"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

// RouteWorker - Worker representation for an http api route handler
type SubscriptionWorker struct {
	topic string
	GrpcWorker
}

func (s *SubscriptionWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return false
}

func (s *SubscriptionWorker) HandlesEvent(trigger *triggers.Event) bool {
	return trigger.Topic == s.topic
}

func (s *SubscriptionWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	// Generate an ID here
	return nil, fmt.Errorf("subscription workers cannot handle HTTP requests")
}

func (s *SubscriptionWorker) HandleEvent(trigger *triggers.Event) error {
	// Generate an ID here
	ID, returnChan := s.newTicket()
	triggerRequest := &pb.TriggerRequest{
		Data:     trigger.Payload,
		MimeType: http.DetectContentType(trigger.Payload),
		Context: &pb.TriggerRequest_Topic{
			Topic: &pb.TopicTriggerContext{
				Topic: trigger.Topic,
				// FIXME: Add missing fields here...
			},
		},
	}

	// construct the message
	message := &pb.ServerMessage{
		Id: ID,
		Content: &pb.ServerMessage_TriggerRequest{
			TriggerRequest: triggerRequest,
		},
	}

	// send the message
	err := s.send(message)

	if err != nil {
		// There was an error enqueuing the message
		return err
	}

	// wait for the response
	// FIXME: Need to handle timeouts here...
	response := <-returnChan

	topic := response.GetTopic()

	if topic == nil {
		// Fatal error in this case
		// We don't have the correct response type for this handler
		return fmt.Errorf("fatal: Error handling event, incorrect response received from function")
	}

	if topic.GetSuccess() {
		return nil
	}

	return fmt.Errorf("error ocurred handling the event")
}

type SubscriptionWorkerOptions struct {
	Topic string
}

// Package private method
// Only a pool may create a new faas worker
func NewSubscriptionWorker(stream pb.FaasService_TriggerStreamServer, opts *SubscriptionWorkerOptions) *SubscriptionWorker {
	return &SubscriptionWorker{
		topic:      opts.Topic,
		GrpcWorker: NewGrpcListener(stream),
	}
}
