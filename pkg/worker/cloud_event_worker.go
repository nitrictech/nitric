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

	"github.com/nitrictech/nitric/pkg/triggers"

	pb "github.com/nitrictech/nitric/pkg/api/nitric/v1"
)

// RouteWorker - Worker representation for an http api route handler
type CloudEventWorker struct {
	sources    []string
	eventTypes []string
	GrpcWorker
}

func (s *CloudEventWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return false
}

func (s *CloudEventWorker) HandlesEvent(trigger *triggers.Event) bool {
	return false
}

func (s *CloudEventWorker) HandlesCloudEvent(trigger *triggers.CloudEvent) bool {
	sourceMatch := false
	eventTypeMatch := false

	// filter on defined sources (will eval true if at least one source matches or no sources were specified)
	if sourceMatch = len(s.sources) == 0; sourceMatch == false {
		for _, source := range s.sources {
			if sourceMatch = source == trigger.Event.Source; sourceMatch {
				break
			}
		}
	}

	// filter on defined event types (will eval true if at least one event type matches or no event types were specified)
	if eventTypeMatch = len(s.eventTypes) == 0; eventTypeMatch == false {
		for _, eventType := range s.eventTypes {
			if eventTypeMatch = eventType == trigger.Event.Type; eventTypeMatch {
				break
			}
		}
	}

	return sourceMatch && eventTypeMatch
}

func (s *CloudEventWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	// Generate an ID here
	return nil, fmt.Errorf("cloud event workers cannot handle HTTP requests")
}

func (s *CloudEventWorker) HandleEvent(trigger *triggers.CloudEvent) bool {
	return false
}

type CloudEventWorkerOptions struct {
	Sources    []string
	EventTypes []string
}

// Package private method
// Only a pool may create a new faas worker
func NewCloudEventWorker(stream pb.FaasService_TriggerStreamServer, opts *CloudEventWorkerOptions) *CloudEventWorker {
	return &CloudEventWorker{
		sources:    opts.Sources,
		eventTypes: opts.EventTypes,
		GrpcWorker: NewGrpcListener(stream),
	}
}
