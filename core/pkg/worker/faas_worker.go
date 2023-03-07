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
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker/adapter"
)

// FaasWorker
// Worker representation for a Nitric FaaS function using gRPC
type FaasWorker struct {
	adapter.Adapter
}

var _ Worker = &FaasWorker{}

func (s *FaasWorker) HandlesTrigger(trigger *v1.TriggerRequest) bool {
	return true
}

// NewFaasWorker - Create a new FaaS worker
func NewFaasWorker(adapter adapter.Adapter) *FaasWorker {
	return &FaasWorker{
		Adapter: adapter,
	}
}
