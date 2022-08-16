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

import "github.com/nitrictech/nitric/pkg/triggers"

// FaasWorker
// Worker representation for a Nitric FaaS function using gRPC
type FaasWorker struct {
	Adapter
}

var _ Worker = &FaasWorker{}

func (s *FaasWorker) HandlesHttpRequest(trigger *triggers.HttpRequest) bool {
	return true
}

func (s *FaasWorker) HandlesEvent(trigger *triggers.Event) bool {
	return true
}

// NewFaasWorker - Create a new FaaS worker
func NewFaasWorker(adapter Adapter) *FaasWorker {
	return &FaasWorker{
		Adapter: adapter,
	}
}
