package gateway_plugin_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	gateway_plugin "github.com/nitric-dev/membrane/plugins/dev/gateway"
	"github.com/nitric-dev/membrane/sources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockHandler struct {
	// store the recieved requests for testing
	requests []*sources.HttpRequest
	events   []*sources.Event
	// provide fixed mock response for testing
	// respondsWith *sdk.NitricResponse
}

func (m *MockHandler) HandleEvent(evt *sources.Event) error {
	if m.events == nil {
		m.events = make([]*sources.Event, 0)
	}

	m.events = append(m.events, evt)

	return nil
}

func (m *MockHandler) HandleHttpRequest(r *sources.HttpRequest) *http.Response {
	if m.requests == nil {
		m.requests = make([]*sources.HttpRequest, 0)
	}

	// Read and re-created a new read stream here...
	body, _ := ioutil.ReadAll(r.Body)
	r.Body = ioutil.NopCloser(bytes.NewReader(body))

	m.requests = append(m.requests, r)

	return &http.Response{
		StatusCode: 200,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte("success"))),
	}
}

func (m *MockHandler) resetRequests() {
	m.requests = make([]*sources.HttpRequest, 0)
	m.events = make([]*sources.Event, 0)
}

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Gateway", func() {

	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	handler := &MockHandler{}
	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	gateway, _ := gateway_plugin.New()

	AfterEach(func() {
		handler.resetRequests()
	})

	// Start the gatewat on a seperate thread so it doesn't block the tests...
	go (gateway.Start)(handler)
	// FIXME: Update gateway to block on channel...
	time.Sleep(500 * time.Millisecond)

	When("Recieving standard HTTP requests", func() {
		When("The request contains standard nitric headers", func() {
			payload := []byte("Test")
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test", gatewayUrl), bytes.NewReader(payload))

			request.Header.Add("x-nitric-request-id", "1234")
			request.Header.Add("x-nitric-payload-type", "test-payload")
			request.Header.Add("User-Agent", "Test")

			It("should succesfully pass on the request", func() {
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Passing through exactly 1 request")
				Expect(handler.requests).To(HaveLen(1))

				handledRequest := handler.requests[0]

				By("Preserving the original request method")
				Expect(handledRequest.Method).To(Equal("POST"))

				By("Preserving the original request URL")
				Expect(handledRequest.Path).To(Equal("/test"))

				// FIXME: Wierd bug occuring in tests,
				// need to validate genuine runtime behaviour here...
				// Seems like the original request stream is closed
				// before we can actually properly assess it
				body, _ := ioutil.ReadAll(handledRequest.Body)
				By("Preserving the original request Body")
				Expect(string(body)).To(Equal("Test"))
			})
		})
		// TODO: Handle cases of missing nitric headers
		// TODO: Handle cases of other non POST methods
	})

	When("Recieving requests from a topic subscription", func() {
		When("The request contains standard nitric headers", func() {
			payload := []byte("Test")
			request, _ := http.NewRequest("POST", gatewayUrl, bytes.NewReader(payload))

			request.Header.Add("x-nitric-request-id", "1234")
			request.Header.Add("x-nitric-payload-type", "test-payload")
			request.Header.Add("x-nitric-source-type", "SUBSCRIPTION")
			request.Header.Add("x-nitric-source", "test-topic")

			It("should succesfully pass on the event", func() {
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Passing through exactly 1 event")
				Expect(handler.events).To(HaveLen(1))

				evt := handler.events[0]

				By("Preserving the provided payload")
				Expect(evt.Payload).To(BeEquivalentTo(payload))

				By("Extracting the provided topic")
				Expect(evt.Topic).To(Equal("test-topic"))

				By("Extracting the provided ID")
				Expect(evt.ID).To(Equal("1234"))
			})
		})
	})
})
