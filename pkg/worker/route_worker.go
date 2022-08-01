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
	"strings"

	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/utils"
)

// RouteWorker - Worker representation for an http api route handler
type RouteWorker struct {
	api     string
	methods []string
	path    string

	Adapter
}

var _ Worker = &RouteWorker{}

// Api - Retrieve the name of the API this
// route worker was registered for
func (s *RouteWorker) Api() string {
	return s.api
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

func (s *RouteWorker) hasMethod(method string) bool {
	for _, m := range s.methods {
		if method == m {
			return true
		}
	}

	return false
}

func (s *RouteWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	if !s.hasMethod(trigger.Method) {
		return false
	}

	_, err := s.extractPathParams(trigger)

	return err == nil
}

func (s *RouteWorker) HandlesEvent(trigger *triggers.Event) bool {
	return false
}

func (s *RouteWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	params, err := s.extractPathParams(trigger)

	if err != nil {
		return nil, err
	}

	trigger.Params = params

	return s.Adapter.HandleHttpRequest(trigger)
}

func (s *RouteWorker) HandleEvent(trigger *triggers.Event) error {
	return fmt.Errorf("route workers cannot handle events")
}

type RouteWorkerOptions struct {
	Api     string
	Path    string
	Methods []string
}

// Package private method
// Only a pool may create a new faas worker
func NewRouteWorker(handler Adapter, opts *RouteWorkerOptions) *RouteWorker {
	return &RouteWorker{
		api:     opts.Api,
		path:    opts.Path,
		methods: opts.Methods,
		Adapter: handler,
	}
}
