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
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// RouteWorker - Worker representation for an http api route handler
type ScheduleWorker struct {
	key string
	adapter.Adapter
}

var _ Worker = &ScheduleWorker{}

func (s *ScheduleWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	if topic := trigger.GetTopic(); topic != nil {
		return ScheduleKeyToTopicName(s.key) == topic.Topic
	}

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

func (s *ScheduleWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	if trigger.GetTopic() == nil {
		return nil, fmt.Errorf("cannot handle given trigger")
	}

	return s.Adapter.HandleTrigger(ctx, trigger)
}

type ScheduleWorkerOptions struct {
	Key string
}

// Package private method
// Only a pool may create a new faas worker
func NewScheduleWorker(adapter adapter.Adapter, opts *ScheduleWorkerOptions) *ScheduleWorker {
	return &ScheduleWorker{
		key:     opts.Key,
		Adapter: adapter,
	}
}
