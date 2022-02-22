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

	pb "github.com/nitrictech/nitric/pkg/api/nitric/v1"
)

// RouteWorker - Worker representation for an http api route handler
type ScheduleWorker struct {
	key string
	GrpcWorker
}

func (s *ScheduleWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return false
}

// ScheduleKeyToTopicName - converts a schedule description to a name for a topic
// e.g. "Prune Customer Orders" -> "prune-customer-orders"
func ScheduleKeyToTopicName(key string) string {
	return strings.ToLower(strings.ReplaceAll(key, " ", "-"))
}

func (s *ScheduleWorker) Key() string {
	return s.key
}

func (s *ScheduleWorker) HandlesEvent(trigger *triggers.Event) bool {
	return ScheduleKeyToTopicName(s.key) == trigger.Topic
}

func (s *ScheduleWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	// Generate an ID here
	return nil, fmt.Errorf("schedule workers cannot handle HTTP requests")
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
