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
type ScheduleWorker struct {
	key string
	GrpcWorker
}

func (s *ScheduleWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return false
}

func (s *ScheduleWorker) HandlesEvent(trigger *triggers.Event) bool {
	// TODO: Determine if event should be handled by using convention for
	// this schedule
	return false
}

func (s *ScheduleWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	// Generate an ID here
	return nil, fmt.Errorf("schedule workers cannot handle HTTP requests")
}

func (s *ScheduleWorker) HandleEvent(trigger *triggers.Event) error {
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

type ScheduleWorkerOptions struct {
	Key string
}

// Package private method
// Only a pool may create a new faas worker
func NewScheduleWorker(stream pb.FaasService_TriggerStreamServer, opts *ScheduleWorkerOptions) *ScheduleWorker {
	return &ScheduleWorker{
		key:        opts.Key,
		GrpcWorker: NewGrpcListener(stream),
	}
}
