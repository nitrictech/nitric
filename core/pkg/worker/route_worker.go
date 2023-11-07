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
	"context"
	"fmt"
	"strings"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// RouteWorker - Worker representation for an http api route handler
type RouteWorker struct {
	api     string
	methods []string
	path    string

	adapter.Adapter
}

var _ Worker = &RouteWorker{}

// Api - Retrieve the name of the API this
// route worker was registered for
func (s *RouteWorker) Api() string {
	return s.api
}

func (s *RouteWorker) ExtractPathParams(path string) (map[string]string, error) {
	requestPathSegments := utils.SplitPath(path)
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

func (s *RouteWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	if http := trigger.GetHttp(); http != nil {
		if !s.hasMethod(http.Method) {
			return false
		}

		_, err := s.ExtractPathParams(http.Path)

		return err == nil
	}

	return false
}

func (s *RouteWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	if http := trigger.GetHttp(); http != nil {
		params, err := s.ExtractPathParams(http.Path)
		if err != nil {
			return nil, err
		}

		http.PathParams = params

		return s.Adapter.HandleTrigger(ctx, trigger)
	}

	return nil, fmt.Errorf("Router Worker does not handle non-HTTP triggers")
}

type RouteWorkerOptions struct {
	Api     string
	Path    string
	Methods []string
}

// Package private method
// Only a pool may create a new faas worker
func NewRouteWorker(adapter adapter.Adapter, opts *RouteWorkerOptions) *RouteWorker {
	return &RouteWorker{
		api:     opts.Api,
		path:    opts.Path,
		methods: opts.Methods,
		Adapter: adapter,
	}
}
