package http_service_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

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
	})
})
