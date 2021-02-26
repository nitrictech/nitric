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
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockHandler struct {
	// store the recieved requests for testing
	handledRequests []*sdk.NitricRequest
	// provide fixed mock response for testing
	respondsWith *sdk.NitricResponse
}

const GATEWAY_ADDRESS = "127.0.0.1:9001"

func (m *MockHandler) handle(request *sdk.NitricRequest) *sdk.NitricResponse {
	if m.handledRequests == nil {
		// Initialize the handled requests array
		m.handledRequests = make([]*sdk.NitricRequest, 0)
	}

	m.handledRequests = append(m.handledRequests, request)

	if m.respondsWith != nil {
		return m.respondsWith
	}

	// If there is no configured mock response, we'll return a default one
	return &sdk.NitricResponse{
		Headers: map[string]string{
			"Content-Type": "text/plain",
		},
		Status: 200,
		Body:   []byte("Test"),
	}
}

func (m *MockHandler) resetRequests() {
	m.handledRequests = make([]*sdk.NitricRequest, 0)
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
	go (httpPlugin.Start)(mockHandler.handle)

	// Delay to allow the HTTP server to correctly start
	// FIXME: Should block on channels...
	time.Sleep(500 * time.Millisecond)

	AfterEach(func() {
		mockHandler.resetRequests()
	})

	When("Invoking the GCP HTTP Gateway", func() {
		When("with a standard Nitric Request", func() {

			It("Should be handled successfully", func() {
				request, err := http.NewRequest("POST", gatewayUrl, bytes.NewReader([]byte("Test")))
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				resp, err := http.DefaultClient.Do(request)

				var responseBody = make([]byte, 0)

				if err == nil {
					if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
						responseBody = bytes
					}
				}

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Handling exactly 1 request")
				Expect(mockHandler.handledRequests).To(HaveLen(1))

				handledRequest := mockHandler.handledRequests[0]

				By("Having the provided RequestId")
				Expect(handledRequest.Context.RequestId).To(Equal("1234"))

				By("Having the provided payload type")
				Expect(handledRequest.Context.PayloadType).To(Equal("Test Payload"))

				By("Have the correct source type")
				Expect(handledRequest.Context.SourceType).To(Equal(sdk.Request))

				By("Have the correct provided source")
				Expect(handledRequest.Context.Source).To(Equal("Test"))

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("Test"))
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
				Expect(mockHandler.handledRequests).To(BeEmpty())

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

				testPayloadBytes, _ := json.Marshal(testPayload)
				testTopic := "test"
				evt := []eventgrid.Event{
					eventgrid.Event{
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
				resp, _ := http.DefaultClient.Do(request)

				By("Passing the event to the Nitric Application")
				Expect(mockHandler.handledRequests).To(HaveLen(1))

				handledRequest := mockHandler.handledRequests[0]
				By("Having the provided requestId")
				Expect(handledRequest.Context.RequestId).To(Equal("1234"))

				By("Having the provided payload type")
				Expect(handledRequest.Context.PayloadType).To(Equal("test-payload"))

				By("Having the provided payload")
				Expect(handledRequest.Payload).To(BeEquivalentTo(testPayloadBytes))

				By("Sending back a 200 OK status")
				Expect(resp.StatusCode).To(Equal(200))
			})
		})
	})
})
