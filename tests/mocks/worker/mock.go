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
	triggers2 "github.com/nitrictech/nitric/pkg/triggers"
)

type MockWorkerOptions struct {
	ReturnHttp *triggers2.HttpResponse
	HttpError  error
	eventError error
}

// MockWorker - A mock worker interface for testing
type MockWorker struct {
	returnHttp       *triggers2.HttpResponse
	httpError        error
	eventError       error
	ReceivedEvents   []*triggers2.Event
	ReceivedRequests []*triggers2.HttpRequest
}

func (m *MockWorker) HandleEvent(trigger *triggers2.Event) error {
	m.ReceivedEvents = append(m.ReceivedEvents, trigger)

	return m.eventError
}

func (m *MockWorker) HandleHttpRequest(trigger *triggers2.HttpRequest) (*triggers2.HttpResponse, error) {
	m.ReceivedRequests = append(m.ReceivedRequests, trigger)

	return m.returnHttp, m.httpError
}

func (m *MockWorker) Reset() {
	m.ReceivedEvents = make([]*triggers2.Event, 0)
	m.ReceivedRequests = make([]*triggers2.HttpRequest, 0)
}

func NewMockWorker(opts *MockWorkerOptions) *MockWorker {
	return &MockWorker{
		httpError:        opts.HttpError,
		returnHttp:       opts.ReturnHttp,
		eventError:       opts.eventError,
		ReceivedEvents:   make([]*triggers2.Event, 0),
		ReceivedRequests: make([]*triggers2.HttpRequest, 0),
	}
}
