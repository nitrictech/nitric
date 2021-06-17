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

package worker_mocks

import (
	"github.com/nitric-dev/membrane/triggers"
)

type MockWorkerOptions struct {
	ReturnHttp *triggers.HttpResponse
	HttpError  error
	eventError error
}

// MockWorker - A mock worker interface for testing
type MockWorker struct {
	returnHttp       *triggers.HttpResponse
	httpError        error
	eventError       error
	RecievedEvents   []*triggers.Event
	RecievedRequests []*triggers.HttpRequest
}

func (m *MockWorker) HandleEvent(trigger *triggers.Event) error {
	m.RecievedEvents = append(m.RecievedEvents, trigger)

	return m.eventError
}

func (m *MockWorker) HandleHttpRequest(trigger *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	m.RecievedRequests = append(m.RecievedRequests, trigger)

	return m.returnHttp, m.httpError
}

func (m *MockWorker) Reset() {
	m.RecievedEvents = make([]*triggers.Event, 0)
	m.RecievedRequests = make([]*triggers.HttpRequest, 0)
}

func NewMockWorker(opts *MockWorkerOptions) *MockWorker {
	return &MockWorker{
		httpError:        opts.HttpError,
		returnHttp:       opts.ReturnHttp,
		eventError:       opts.eventError,
		RecievedEvents:   make([]*triggers.Event, 0),
		RecievedRequests: make([]*triggers.HttpRequest, 0),
	}
}
