package http_service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	http_plugin "github.com/nitric-dev/membrane/plugins/azure/gateway/http"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/sources"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

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

var _ = Describe("Http", func() {
	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	// Set this to loopback to ensure its not public in our CI/Testing environments
	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	mockHandler := &MockHandler{}
	httpPlugin, _ := http_plugin.New()
	// Run on a non-blocking thread
	go (httpPlugin.Start)(mockHandler)

	// Delay to allow the HTTP server to correctly start
	// FIXME: Should block on channels...
	time.Sleep(500 * time.Millisecond)

	AfterEach(func() {
		mockHandler.resetRequests()
	})

	When("Invoking the GCP HTTP Gateway", func() {
		When("with a standard Nitric Request", func() {

			It("Should be handled successfully", func() {
				request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test/", gatewayUrl), bytes.NewReader([]byte("Test")))
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				_, err := http.DefaultClient.Do(request)

				//var responseBody = make([]byte, 0)

				//if err == nil {
				//	if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				//		responseBody = bytes
				//	}
				//}

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Handling exactly 1 request")
				Expect(mockHandler.requests).To(HaveLen(1))

				handledRequest := mockHandler.requests[0]

				By("Having the provided path")
				Expect(handledRequest.Path).To((Equal("/test/")))
			})
		})

		When("With a SubscriptionValidation event", func() {
			It("Should return the provided validation code", func() {
				validationCode := "test"
				evt := []eventgrid.Event{
					eventgrid.Event{
						Data: eventgrid.SubscriptionValidationEventData{
							ValidationCode: &validationCode,
						},
					},
				}

				requestBody, _ := json.Marshal(evt)
				request, _ := http.NewRequest("POST", gatewayUrl, bytes.NewReader([]byte(requestBody)))
				request.Header.Add("aeg-event-type", "SubscriptionValidation")
				resp, _ := http.DefaultClient.Do(request)

				By("Not invoking the nitric application")
				Expect(mockHandler.requests).To(BeEmpty())

				By("Returning a 200 response")
				Expect(resp.StatusCode).To(Equal(200))

				By("Containing the provided validation code")
				var respEvt eventgrid.SubscriptionValidationResponse
				bytes, _ := ioutil.ReadAll(resp.Body)
				json.Unmarshal(bytes, &respEvt)
				Expect(*respEvt.ValidationResponse).To(BeEquivalentTo(validationCode))
			})
		})

		When("With a Notification event", func() {
			It("Should successfully handle the notification", func() {
				testPayload := map[string]interface{}{
					"Test": "Test",
				}

				// testPayloadBytes, _ := json.Marshal(testPayload)
				testTopic := "test"
				testID := "1234"
				evt := []eventgrid.Event{
					eventgrid.Event{
						ID:    &testID,
						Topic: &testTopic,
						Data: sdk.NitricEvent{
							RequestId:   "1234",
							PayloadType: "test-payload",
							Payload:     testPayload,
						},
					},
				}

				requestBody, _ := json.Marshal(evt)
				request, _ := http.NewRequest("POST", gatewayUrl, bytes.NewReader([]byte(requestBody)))
				request.Header.Add("aeg-event-type", "Notification")
				_, _ = http.DefaultClient.Do(request)

				By("Passing the event to the Nitric Application")
				Expect(mockHandler.events).To(HaveLen(1))

				event := mockHandler.events[0]
				By("Having the provided requestId")
				Expect(event.ID).To(Equal("1234"))
			})
		})
	})
})
