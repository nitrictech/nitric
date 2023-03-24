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

package adapter

import (
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type GrpcAdapter struct {
	stream v1.FaasService_TriggerStreamServer
	// Response channels for this worker
	responseQueueLock sync.Locker
	responseQueue     map[string]chan *v1.TriggerResponse
}

var _ Adapter = &GrpcAdapter{}

// newTicket - Generates a request/response ID and response channel
// for the requesting thread to wait on
func (s *GrpcAdapter) newTicket() (string, chan *v1.TriggerResponse) {
	s.responseQueueLock.Lock()
	defer s.responseQueueLock.Unlock()

	ID := uuid.New().String()
	responseChan := make(chan *v1.TriggerResponse)

	s.responseQueue[ID] = responseChan

	return ID, responseChan
}

// resolveTicket - Retrieves a response channel from the queue for
// the given ID and removes the entry from the map
func (s *GrpcAdapter) resolveTicket(ID string) (chan *v1.TriggerResponse, error) {
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

func (gwb *GrpcAdapter) send(msg *v1.ServerMessage) error {
	return gwb.stream.Send(msg)
}

func (gwb *GrpcAdapter) Start(errchan chan error) {
	for {
		msg, err := gwb.stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				// return will close stream from server side
				log.Println("exit")
			}
			log.Printf("received error %v", err)

			errchan <- err
			return
		}

		if msg.GetInitRequest() != nil {
			log.Default().Println("Received init request from worker")
			err = gwb.stream.Send(&v1.ServerMessage{
				Content: &v1.ServerMessage_InitResponse{
					InitResponse: &v1.InitResponse{},
				},
			})
			if err != nil {
				log.Default().Printf("send error %v", err)
			}
			continue
		}

		// Load the response channel and delete its map key reference
		val, err := gwb.resolveTicket(msg.GetId())
		if err != nil {
			err = errors.WithMessage(err, "Fatal: FaaS Worker in bad state closing stream: "+msg.GetId())
			log.Default().Println(err.Error())
			errchan <- err
			return
		}
		// For now assume this is a trigger response...
		response := msg.GetTriggerResponse()
		// Write the response the the waiting recipient
		val <- response
	}
}

func (s *GrpcAdapter) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	ID, returnChan := s.newTicket()

	// construct the message
	message := &v1.ServerMessage{
		Id: ID,
		Content: &v1.ServerMessage_TriggerRequest{
			TriggerRequest: trigger,
		},
	}

	// send the message
	err := s.send(message)
	if err != nil {
		// There was an error enqueuing the message
		return nil, err
	}

	triggerResponse := <-returnChan

	return triggerResponse, nil
}

func NewGrpcAdapter(stream v1.FaasService_TriggerStreamServer) *GrpcAdapter {
	return &GrpcAdapter{
		stream:            stream,
		responseQueueLock: &sync.Mutex{},
		responseQueue:     make(map[string]chan *v1.TriggerResponse),
	}
}
