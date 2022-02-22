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

	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	v1 "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/errors"
	"github.com/nitrictech/nitric/pkg/triggers"
)

type GrpcWorker interface {
	Worker
	Listen(chan error)
}

type grpcWorkerBase struct {
	stream v1.FaasService_TriggerStreamServer
	// Response channels for this worker
	responseQueueLock sync.Locker
	responseQueue     map[string]chan *v1.TriggerResponse
}

// newTicket - Generates a request/response ID and response channel
// for the requesting thread to wait on
func (s *grpcWorkerBase) newTicket() (string, chan *v1.TriggerResponse) {
	s.responseQueueLock.Lock()
	defer s.responseQueueLock.Unlock()

	ID := uuid.New().String()
	responseChan := make(chan *v1.TriggerResponse)

	s.responseQueue[ID] = responseChan

	return ID, responseChan
}

// resolveTicket - Retrieves a response channel from the queue for
// the given ID and removes the entry from the map
func (s *grpcWorkerBase) resolveTicket(ID string) (chan *v1.TriggerResponse, error) {
	s.responseQueueLock.Lock()
	defer func() {
		delete(s.responseQueue, ID)
		s.responseQueueLock.Unlock()
	}()

	if s.responseQueue[ID] == nil {
		return nil, fmt.Errorf("attempted to resolve ticket that does not exist!")
	}

	return s.responseQueue[ID], nil
}

func (gwb *grpcWorkerBase) send(msg *v1.ServerMessage) error {
	return gwb.stream.Send(msg)
}

func (gwb *grpcWorkerBase) Listen(errchan chan error) {
	for {
		msg, err := gwb.stream.Recv()

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
			log.Default().Println("Received init request from worker")
			gwb.stream.Send(&v1.ServerMessage{
				Content: &v1.ServerMessage_InitResponse{
					InitResponse: &v1.InitResponse{},
				},
			})
			continue
		}

		// Load the response channel and delete its map key reference
		if val, err := gwb.resolveTicket(msg.GetId()); err == nil {
			// For now assume this is a trigger response...
			response := msg.GetTriggerResponse()
			// Write the response the the waiting recipient
			val <- response
		} else {
			err := errors.Fatal("FaaS Worker in bad state closing stream: " + msg.GetId())
			log.Default().Println(err.Error())
			errchan <- err
			break
		}
	}
}

func (s *grpcWorkerBase) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return true
}

func (s *grpcWorkerBase) HandlesEvent(trigger *triggers.Event) bool {
	return true
}

func (s *grpcWorkerBase) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	// Generate an ID here
	ID, returnChan := s.newTicket()

	var mimeType string = ""
	if trigger.Header != nil && len(trigger.Header["Content-Type"]) > 0 {
		mimeType = trigger.Header["Content-Type"][0]
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(trigger.Body)
	}

	headers := make(map[string]*v1.HeaderValue)
	headersOld := make(map[string]string)
	for k, v := range trigger.Header {
		if v != nil {
			headers[k] = &v1.HeaderValue{
				Value: v,
			}
			if len(v) > 0 {
				headersOld[k] = v[0]
			}
		}
	}

	query := make(map[string]*v1.QueryValue)
	queryOld := make(map[string]string)
	for k, v := range trigger.Query {
		if v != nil {
			query[k] = &v1.QueryValue{
				Value: v,
			}
			if len(v) > 0 {
				queryOld[k] = v[0]
			}
		}
	}

	triggerRequest := &v1.TriggerRequest{
		Data:     trigger.Body,
		MimeType: mimeType,
		Context: &v1.TriggerRequest_Http{
			Http: &v1.HttpTriggerContext{
				Path:           trigger.Path,
				Method:         trigger.Method,
				QueryParams:    query,
				QueryParamsOld: queryOld,
				Headers:        headers,
				HeadersOld:     headersOld,
				PathParams:     trigger.Params,
			},
		},
	}

	// construct the message
	message := &v1.ServerMessage{
		Id: ID,
		Content: &v1.ServerMessage_TriggerRequest{
			TriggerRequest: triggerRequest,
		},
	}

	// send the message
	err := s.send(message)

	if err != nil {
		// There was an error enqueuing the message
		return nil, err
	}

	// wait for the response
	triggerResponse := <-returnChan

	httpResponse := triggerResponse.GetHttp()

	if httpResponse == nil {
		return nil, fmt.Errorf("fatal: Error handling event, incorrect response received from function")
	}

	fasthttpHeader := &fasthttp.ResponseHeader{}

	for key, val := range httpResponse.GetHeaders() {
		headerList := val.Value
		if key == "Set-Cookie" || key == "Cookie" {
			for _, v := range headerList {
				fasthttpHeader.Add(key, v)
			}
		} else if len(headerList) > 0 {
			fasthttpHeader.Set(key, headerList[0])
		}
	}

	response := &triggers.HttpResponse{
		Body: triggerResponse.Data,
		// No need to worry about integer truncation
		// as this should be a HTTP status code...
		StatusCode: int(httpResponse.Status),
		Header:     fasthttpHeader,
	}

	return response, nil
}

func (s *grpcWorkerBase) HandleEvent(trigger *triggers.Event) error {
	// Generate an ID here
	ID, returnChan := s.newTicket()
	triggerRequest := &v1.TriggerRequest{
		Data:     trigger.Payload,
		MimeType: http.DetectContentType(trigger.Payload),
		Context: &v1.TriggerRequest_Topic{
			Topic: &v1.TopicTriggerContext{
				Topic: trigger.Topic,
				// FIXME: Add missing fields here...
			},
		},
	}

	// construct the message
	message := &v1.ServerMessage{
		Id: ID,
		Content: &v1.ServerMessage_TriggerRequest{
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
		return fmt.Errorf("Fatal: Error handling event, incorrect response received from function")
	}

	if topic.GetSuccess() {
		return nil
	}

	return fmt.Errorf("Error ocurred handling the event")
}

func NewGrpcListener(stream v1.FaasService_TriggerStreamServer) *grpcWorkerBase {
	return &grpcWorkerBase{
		stream:            stream,
		responseQueueLock: &sync.Mutex{},
		responseQueue:     make(map[string]chan *v1.TriggerResponse),
	}
}
