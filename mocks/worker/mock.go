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
