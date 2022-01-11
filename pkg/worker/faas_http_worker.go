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
	"net"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
	"google.golang.org/protobuf/encoding/protojson"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
	"github.com/nitrictech/nitric/pkg/triggers"
)

// A Nitric HTTP worker
type FaasHttpWorker struct {
	address string
}

var METHOD_TYPE = []byte("POST")

func (s *FaasHttpWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return true
}

func (s *FaasHttpWorker) HandlesEvent(trigger *triggers.Event) bool {
	return true
}

// HandleEvent - Handles an event from a subscription by converting it to an HTTP request.
func (h *FaasHttpWorker) HandleEvent(trigger *triggers.Event) error {
	address := fmt.Sprintf("http://%s", h.address)
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	// Release resources after finishing
	defer func() {
		request.Reset()
		response.Reset()
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	triggerRequest := &pb.TriggerRequest{
		Data:     trigger.Payload,
		MimeType: http.DetectContentType(trigger.Payload),
		Context: &pb.TriggerRequest_Topic{
			Topic: &pb.TopicTriggerContext{
				Topic: trigger.Topic,
			},
		},
	}

	if jsonData, err := protojson.Marshal(triggerRequest); err == nil {
		fmt.Println(fmt.Sprintf("Membrane receieved event:\n%s", string(jsonData)))
		request.Header.SetContentType("application/json")
		request.SetBody(jsonData)
		request.SetRequestURI(address)

		err := fasthttp.Do(request, response)

		if err != nil {
			return fmt.Errorf("Function request failed")
		}

		// Response body should contain an instance of triggerResponse
		var triggerResponse pb.TriggerResponse
		err = protojson.Unmarshal(response.Body(), &triggerResponse)

		if err != nil {
			return err
		}

		topic := triggerResponse.GetTopic()

		if topic != nil {
			if topic.Success {
				return nil
			}

			return fmt.Errorf("topic context indicated processing was unsuccessful")
		}

		return fmt.Errorf("response from function did not contain topic context")
	} else {
		return fmt.Errorf("error marshalling request. Details: %v", err)
	}
}

// HandleHttpRequest - Handles an HTTP request by forwarding it as an HTTP request.
func (h *FaasHttpWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	address := fmt.Sprintf("http://%s", h.address)
	request := fasthttp.AcquireRequest()
	response := fasthttp.AcquireResponse()

	// Release resources after finishing
	defer func() {
		request.Reset()
		response.Reset()
		fasthttp.ReleaseRequest(request)
		fasthttp.ReleaseResponse(response)
	}()

	var mimeType string = ""
	if trigger.Header != nil && len(trigger.Header["Content-Type"]) > 0 {
		mimeType = trigger.Header["Content-Type"][0]
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(trigger.Body)
	}

	headers := make(map[string]*pb.HeaderValue)
	for k, v := range trigger.Header {
		headers[k] = &pb.HeaderValue{
			Value: v,
		}
	}

	query := make(map[string]*pb.QueryValue)
	for k, v := range trigger.Query {
		query[k] = &pb.QueryValue{
			Value: v,
		}
	}

	triggerRequest := &pb.TriggerRequest{
		Data:     trigger.Body,
		MimeType: mimeType,
		Context: &pb.TriggerRequest_Http{
			Http: &pb.HttpTriggerContext{
				Path:        trigger.Path,
				Headers:     headers,
				Method:      trigger.Method,
				QueryParams: query,
			},
		},
	}

	if jsonData, err := protojson.Marshal(triggerRequest); err == nil {
		request.Header.SetContentType("application/json")
		request.SetBody(jsonData)
		request.SetRequestURI(address)

		err := fasthttp.Do(request, response)

		if err != nil {
			return nil, err
		}

		// Response body should contain an instance of triggerResponse
		var triggerResponse pb.TriggerResponse
		err = protojson.Unmarshal(response.Body(), &triggerResponse)

		if err != nil {
			return nil, err
		}

		return triggers.FromTriggerResponse(&triggerResponse)
	} else {
		return nil, err
	}
}

// Creates a new FaasHttpWorker
// Will wait to ensure that the provided address is dialable
// before proceeding
func NewFaasHttpWorker(address string) (*FaasHttpWorker, error) {
	// Dial the child port to see if it's open and ready...
	maxWaitTime := time.Duration(5) * time.Second
	// Longer poll times, e.g. 200 milliseconds results in slow lambda cold starts (15s+)
	pollInterval := time.Duration(15) * time.Millisecond

	var waitedTime = time.Duration(0)
	for {
		conn, _ := net.Dial("tcp", address)
		if conn != nil {
			conn.Close()
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				return nil, fmt.Errorf("Unable to dial http worker, does it expose a http server at: %s?", address)
			}
		}
	}
	// Dial the provided address to ensure its availability
	return &FaasHttpWorker{
		address: address,
	}, nil
}
