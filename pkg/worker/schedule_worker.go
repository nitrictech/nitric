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
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/nitrictech/nitric/pkg/triggers"

	"github.com/google/uuid"
	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

// RouteWorker - Worker representation for an http api route handler
type ScheduleWorker struct {
	key string
	// gRPC Stream for this worker
	stream pb.FaasService_TriggerStreamServer
	// Response channels for this worker
	responseQueueLock sync.Mutex
	responseQueue     map[string]chan *pb.TriggerResponse
}

// newTicket - Generates a request/response ID and response channel
// for the requesting thread to wait on
func (s *ScheduleWorker) newTicket() (string, chan *pb.TriggerResponse) {
	s.responseQueueLock.Lock()
	defer s.responseQueueLock.Unlock()

	ID := uuid.New().String()
	responseChan := make(chan *pb.TriggerResponse)

	s.responseQueue[ID] = responseChan

	return ID, responseChan
}

// resolveTicket - Retrieves a response channel from the queue for
// the given ID and removes the entry from the map
func (s *ScheduleWorker) resolveTicket(ID string) (chan *pb.TriggerResponse, error) {
	s.responseQueueLock.Lock()
	defer func() {
		delete(s.responseQueue, ID)
		s.responseQueueLock.Unlock()
	}()

	if s.responseQueue[ID] == nil {
		return nil, fmt.Errorf("attempted to resolve ticket that does not exist")
	}

	return s.responseQueue[ID], nil
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
	err := s.stream.Send(message)

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

// listen
func (s *ScheduleWorker) Listen(errchan chan error) {
	// Listen for responses
	for {
		msg, err := s.stream.Recv()

		if err != nil {
			if err == io.EOF {
				// return will close stream from server side
				log.Println("exit")
			}
			if err != nil {
				log.Printf("received error %v", err)
			}

			errchan <- err
			break
		}

		if msg.GetInitRequest() != nil {
			errchan <- fmt.Errorf("init request recieved during runtime, exiting")
			break
		}

		// Load the response channel and delete its map key reference
		if val, err := s.resolveTicket(msg.GetId()); err == nil {
			// For now assume this is a trigger response...
			response := msg.GetTriggerResponse()
			// Write the response the the waiting recipient
			val <- response
		} else {
			fmt.Println("fatal: FaaS Worker in bad state closing stream: ", msg.GetId())
			errchan <- fmt.Errorf("fatal: FaaS Worker in bad state closing stream! %v", msg.GetId())
			break
		}
	}
}

type ScheduleWorkerOptions struct {
	Key string
}

// Package private method
// Only a pool may create a new faas worker
func NewScheduleWorker(stream pb.FaasService_TriggerStreamServer, opts *ScheduleWorkerOptions) *ScheduleWorker {
	return &ScheduleWorker{
		key:               opts.Key,
		stream:            stream,
		responseQueueLock: sync.Mutex{},
		responseQueue:     make(map[string]chan *pb.TriggerResponse),
	}
}