// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apis

import (
	"fmt"
	"slices"
	"strings"
	"sync"

	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	"github.com/nitrictech/nitric/core/pkg/workers"
)

type WorkerConnection = workers.WorkerRequestBroker[*apispb.ServerMessage, *apispb.ClientMessage]

type RouteWorker struct {
	routeMatcher string
	methods      []string
	connection   *WorkerConnection
}

// slashSplitter - used to split strings, with the same output regardless of leading or trailing slashes
// e.g - strings.FieldsFunc("/one/two/three/", f) == strings.FieldsFunc("/one/two/three", f) == strings.FieldsFunc("one/two/three", f) == ["one" "two" "three"]
func slashSplitter(c rune) bool {
	return c == '/'
}

type ApiRequestHandler interface {
	apispb.ApiServer
	HandleRequest(apiName string, request *apispb.ServerMessage) (*apispb.ClientMessage, error)
	WorkerCount() int
}

// SplitPath - splits a path into its component parts, ignoring leading or trailing slashes.
// e.g - SplitPath("/one/two/three/") == SplitPath("/one/two/three") == SplitPath("one/two/three") == ["one" "two" "three"]
func splitPath(p string) []string {
	return strings.FieldsFunc(p, slashSplitter)
}

func extractPathParams(route string, requestPath string) (map[string]string, error) {
	requestPathSegments := splitPath(requestPath)
	pathSegments := splitPath(route)
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

func (r *RouteWorker) isSupportedRequest(httpRequest *apispb.HttpRequest) bool {
	if !slices.Contains[[]string](r.methods, httpRequest.GetMethod()) {
		return false
	}

	_, err := extractPathParams(r.routeMatcher, httpRequest.GetPath())

	return err == nil
}

type ApiName = string

type RouteWorkerManager struct {
	routeWorkerMap map[ApiName][]*RouteWorker
	lock           sync.RWMutex
}

var _ apispb.ApiServer = &RouteWorkerManager{}

func (s *RouteWorkerManager) WorkerCount() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	total := 0
	for _, workers := range s.routeWorkerMap {
		total += len(workers)
	}

	return total
}

// registerRouteHandler registers a worker by the routes and methods it handles.
func (s *RouteWorkerManager) registerRouteHandler(apiName string, path string, methods []string, worker *WorkerConnection) (*RouteWorker, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, exists := s.routeWorkerMap[apiName]; !exists {
		s.routeWorkerMap[apiName] = []*RouteWorker{}
	}

	rw := &RouteWorker{
		routeMatcher: path,
		methods:      methods,
		connection:   worker,
	}
	s.routeWorkerMap[apiName] = append(s.routeWorkerMap[apiName], rw)

	return rw, nil
}

func (s *RouteWorkerManager) unregisterRouteHandler(apiName string, worker *RouteWorker) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.routeWorkerMap[apiName] = slices.DeleteFunc[[]*RouteWorker](s.routeWorkerMap[apiName], func(rw *RouteWorker) bool {
		return rw == worker
	})

	if len(s.routeWorkerMap[apiName]) == 0 {
		delete(s.routeWorkerMap, apiName)
	}
}

func (s *RouteWorkerManager) Serve(stream apispb.Api_ServeServer) error {
	initRequest, err := stream.Recv()
	if err != nil {
		return err
	}

	if initRequest.GetRegistrationRequest() == nil {
		return fmt.Errorf("first request must be an init request")
	}

	wrkr := workers.NewWorkerRequestBroker[*apispb.ServerMessage, *apispb.ClientMessage](stream)

	// Get routing info
	apiName := initRequest.GetRegistrationRequest().Api
	path := initRequest.GetRegistrationRequest().Path
	methods := initRequest.GetRegistrationRequest().Methods

	routeWorker, err := s.registerRouteHandler(apiName, path, methods, wrkr)
	if err != nil {
		return err
	}

	defer s.unregisterRouteHandler(apiName, routeWorker)

	// send ack of registration
	err = stream.Send(&apispb.ServerMessage{
		Content: &apispb.ServerMessage_RegistrationResponse{
			RegistrationResponse: &apispb.RegistrationResponse{},
		},
	})
	if err != nil {
		return err
	}

	return wrkr.Run()
}

func (s *RouteWorkerManager) HandleRequest(apiName string, request *apispb.ServerMessage) (*apispb.ClientMessage, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	if request.Id == "" {
		request.Id = workers.GenerateUniqueId()
	}

	possibleHandlers, ok := s.routeWorkerMap[apiName]
	if !ok {
		return nil, fmt.Errorf("no routes registered for api %s", apiName)
	}

	// Handlers are applied using Highlander rules (THERE CAN BE ONLY ONE!!!)
	var theOneTrueHandler *RouteWorker = nil
	for _, handler := range possibleHandlers {
		if handler.isSupportedRequest(request.GetHttpRequest()) {
			// Praise the one true handler 🙌
			theOneTrueHandler = handler
			break
		}
	}

	if theOneTrueHandler == nil {
		return nil, fmt.Errorf("no worker registered for Api %s on route: %s - %s", apiName, request.GetHttpRequest().GetMethod(), request.GetHttpRequest().GetPath())
	}

	if request.GetHttpRequest().GetPathParams() == nil || len(request.GetHttpRequest().GetPathParams()) < 1 {
		pathParams, err := extractPathParams(theOneTrueHandler.routeMatcher, request.GetHttpRequest().GetPath())
		if err != nil {
			return nil, err
		}
		request.GetHttpRequest().PathParams = pathParams
	}

	resp, err := theOneTrueHandler.connection.Send(request)
	if err != nil {
		return nil, err
	}

	return *resp, nil
}

func New() *RouteWorkerManager {
	return &RouteWorkerManager{
		routeWorkerMap: map[string][]*RouteWorker{},
		lock:           sync.RWMutex{},
	}
}
