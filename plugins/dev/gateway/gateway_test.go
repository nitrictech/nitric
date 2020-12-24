package gateway_plugin_test

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"time"

	gateway_plugin "github.com/nitric-dev/membrane/plugins/dev/gateway"
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

const GATEWAY_ADDRESS = "127.0.0.1:9002"

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
	go (gateway.Start)(handler.handle)
	// FIXME: Update gateway to block on channel...
	time.Sleep(500 * time.Millisecond)

	When("Recieving standard HTTP requests", func() {
		When("The request contains standard nitric headers", func() {
			payload := []byte("Test")
			request, _ := http.NewRequest("POST", gatewayUrl, bytes.NewReader(payload))

			request.Header.Add("x-nitric-request-id", "1234")
			request.Header.Add("x-nitric-payload-type", "test-payload")
			request.Header.Add("User-Agent", "Test")

			It("should succesfully pass on the request", func() {
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Passing through exactly 1 request")
				Expect(handler.handledRequests).To(HaveLen(1))

				handledRequest := handler.handledRequests[0]

				By("Passing setting the nitric source type to REQUEST")
				Expect(handledRequest.Context.SourceType).To(Equal(sdk.Request))

				By("Passing through the User-Agent as the nitric source")
				Expect(handledRequest.Context.Source).To(Equal("Test"))

				By("Passing through the sent nitric request id")
				Expect(handledRequest.Context.RequestId).To(Equal("1234"))

				By("Passing through the sent nitric payload type")
				Expect(handledRequest.Context.PayloadType).To(Equal("test-payload"))

				By("Passing through the provided payload")
				Expect(handledRequest.Payload).To(BeEquivalentTo(payload))
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

				By("Passing through exactly 1 request")
				Expect(handler.handledRequests).To(HaveLen(1))

				handledRequest := handler.handledRequests[0]

				By("Passing setting the nitric source type to SUBSCRIPTION")
				Expect(handledRequest.Context.SourceType).To(Equal(sdk.Subscription))

				By("Passing through the source header as the nitric source")
				Expect(handledRequest.Context.Source).To(Equal("test-topic"))

				By("Passing through the sent nitric request id")
				Expect(handledRequest.Context.RequestId).To(Equal("1234"))

				By("Passing through the sent nitric payload type")
				Expect(handledRequest.Context.PayloadType).To(Equal("test-payload"))

				By("Passing through the provided payload")
				Expect(handledRequest.Payload).To(BeEquivalentTo(payload))
			})
		})
	})
})
