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

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// RouteWorker - Worker representation for an http api route handler
type SubscriptionWorker struct {
	topic string
	Delegate
	adapter.Adapter
}

var _ Worker = &SubscriptionWorker{}

func (s *SubscriptionWorker) Topic() string {
	return s.topic
}

func (s *SubscriptionWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	if topic := trigger.GetTopic(); topic != nil {
		return topic.Topic == s.topic
	}

	return false
}

func (s *SubscriptionWorker) HandleTrigger(ctx context.Context, trigger *v1.TriggerRequest) (*v1.TriggerResponse, error) {
	if trigger.GetTopic() == nil {
		return nil, fmt.Errorf("subscription workers cannot handle HTTP requests")
	}

	return s.Adapter.HandleTrigger(ctx, trigger)
}

type SubscriptionWorkerOptions struct {
	Topic string
}

// Package private method
// Only a pool may create a new faas worker
func NewSubscriptionWorker(adapter adapter.Adapter, opts *SubscriptionWorkerOptions) *SubscriptionWorker {
	return &SubscriptionWorker{
		topic:   opts.Topic,
		Adapter: adapter,
	}
}
