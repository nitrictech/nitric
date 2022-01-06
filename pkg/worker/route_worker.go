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
	"strings"

	"github.com/nitrictech/nitric/pkg/utils"

	"github.com/nitrictech/nitric/pkg/triggers"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
	"github.com/valyala/fasthttp"
)

// RouteWorker - Worker representation for an http api route handler
type RouteWorker struct {
	methods []string
	path    string

	GrpcWorker
}

func (s *RouteWorker) extractPathParams(trigger *triggers.HttpRequest) (map[string]string, error) {
	requestPathSegments := utils.SplitPath(trigger.Path)
	pathSegments := utils.SplitPath(s.path)
	params := make(map[string]string)

	// TODO: Filter for trailing/leading slashes
	if len(requestPathSegments) != len(pathSegments) {
		return nil, fmt.Errorf("path template mismatch")
	}

	for i, p := range pathSegments {
		if !strings.HasPrefix(p, ":") && p != requestPathSegments[i] {
			return nil, fmt.Errorf("path template mismatch")
		} else if strings.HasPrefix(p, ":") {
			params[strings.Replace(p, ":", "", 1)] = requestPathSegments[i]
		}
	}

	return params, nil
}

func (s *RouteWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	_, err := s.extractPathParams(trigger)

	return err == nil
}

func (s *RouteWorker) HandlesEvent(trigger *triggers.Event) bool {
	return false
}

func (s *RouteWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	// Generate an ID here
	ID, returnChan := s.newTicket()

	var mimeType string = ""
	if trigger.Header != nil && len(trigger.Header["Content-Type"]) > 0 {
		mimeType = trigger.Header["Content-Type"][0]
	}

	if mimeType == "" {
		mimeType = http.DetectContentType(trigger.Body)
	}

	headers := make(map[string]*pb.HeaderValue)
	headersOld := make(map[string]string)
	for k, v := range trigger.Header {
		if v != nil {
			headers[k] = &pb.HeaderValue{
				Value: v,
			}
			if len(v) > 0 {
				headersOld[k] = v[0]
			}
		}
	}

	query := make(map[string]*pb.QueryValue)
	queryOld := make(map[string]string)
	for k, v := range trigger.Query {
		if v != nil {
			query[k] = &pb.QueryValue{
				Value: v,
			}
			if len(v) > 0 {
				queryOld[k] = v[0]
			}
		}
	}

	params, err := s.extractPathParams(trigger)

	if err != nil {
		return nil, err
	}

	triggerRequest := &pb.TriggerRequest{
		Data:     trigger.Body,
		MimeType: mimeType,
		Context: &pb.TriggerRequest_Http{
			Http: &pb.HttpTriggerContext{
				Path:           trigger.Path,
				Method:         trigger.Method,
				QueryParams:    query,
				QueryParamsOld: queryOld,
				Headers:        headers,
				HeadersOld:     headersOld,
				PathParams:     params,
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
	err = s.send(message)

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

func (s *RouteWorker) HandleEvent(trigger *triggers.Event) error {
	return fmt.Errorf("route workers cannot handle events")
}

type RouteWorkerOptions struct {
	Path    string
	Methods []string
}

// Package private method
// Only a pool may create a new faas worker
func NewRouteWorker(stream pb.FaasService_TriggerStreamServer, opts *RouteWorkerOptions) *RouteWorker {
	return &RouteWorker{
		path:       opts.Path,
		methods:    opts.Methods,
		GrpcWorker: NewGrpcListener(stream),
	}
}
