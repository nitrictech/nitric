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
	http_plugin "github.com/nitric-dev/membrane/plugins/gateway/appservice"
	"github.com/nitric-dev/membrane/triggers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

type MockHandler struct {
	// store the recieved requests for testing
	requests []*triggers.HttpRequest
	events   []*triggers.Event
	// provide fixed mock response for testing
	// respondsWith *sdk.NitricResponse
}

func (m *MockHandler) HandleEvent(evt *triggers.Event) error {
	if m.events == nil {
		m.events = make([]*triggers.Event, 0)
	}

	fmt.Println("HandlingEvent")

	m.events = append(m.events, evt)

	return nil
}

func (m *MockHandler) HandleHttpRequest(r *triggers.HttpRequest) (*triggers.HttpResponse, error) {
	if m.requests == nil {
		m.requests = make([]*triggers.HttpRequest, 0)
	}

	fmt.Println("HandlingRequest")

	// Read and re-created a new read stream here...
	// body := r.Body
	// r.Body = ioutil.NopCloser(bytes.NewReader(body))

	m.requests = append(m.requests, r)

	return &triggers.HttpResponse{
		StatusCode: 200,
		Body:       []byte("success"),
	}, nil
}

func (m *MockHandler) resetRequests() {
	m.requests = make([]*triggers.HttpRequest, 0)
	m.events = make([]*triggers.Event, 0)
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
	time.Sleep(1000 * time.Millisecond)

	AfterEach(func() {
		mockHandler.resetRequests()
	})

	When("Invoking the Azure AppService HTTP Gateway", func() {
		When("with a standard Nitric Request", func() {

			It("Should be handled successfully", func() {
				request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test/", gatewayUrl), bytes.NewReader([]byte("Test")))
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				_, err := http.DefaultClient.Do(request)

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
					{
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
				payload := map[string]string{
					"testing": "test",
				}
				payloadBytes, _ := json.Marshal(payload)
				testTopic := "test"
				testID := "1234"
				evt := []eventgrid.Event{
					{
						ID:    &testID,
						Topic: &testTopic,
						Data:  payload,
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

				By("Having the provided topic")
				Expect(event.Topic).To(Equal("test"))

				By("Having the provided payload")
				Expect(event.Payload).To(BeEquivalentTo(payloadBytes))
			})
		})
	})
})
