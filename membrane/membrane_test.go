package membrane_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/nitric-dev/membrane/membrane"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockEventingServer struct {
	sdk.UnimplementedEventingPlugin
}

type MockStorageServer struct {
	sdk.UnimplementedStoragePlugin
}

type MockDocumentsServer struct {
	sdk.UnimplementedDocumentsPlugin
}

type MockFunction struct {
	// Records the requests that its recieved for later inspection
	requests []*http.Request
	// Returns a fixed HTTP response
	response *http.Response
}

func (m *MockFunction) handler(rw http.ResponseWriter, req *http.Request) {

	if m.requests == nil {
		m.requests = make([]*http.Request, 0)
	}

	m.requests = append(m.requests, req)

	for key, value := range m.response.Header {
		rw.Header().Add(key, strings.Join(value, ""))
	}
	rw.WriteHeader(m.response.StatusCode)

	var rBody []byte = nil
	if m.response.Body != nil {
		rBody, _ = ioutil.ReadAll(m.response.Body)
	}

	rw.Write(rBody)
}

type MockGateway struct {
	sdk.UnimplementedGatewayPlugin
	// The nitric requests to process
	requests []*sdk.NitricRequest
	// store responses for inspection
	responses []*sdk.NitricResponse
	started   bool
}

func (gw *MockGateway) Start(handler sdk.GatewayHandler) error {
	// Spy on the mock gateway
	gw.responses = make([]*sdk.NitricResponse, 0)

	gw.started = true
	if gw.requests != nil {
		for _, request := range gw.requests {
			gw.responses = append(gw.responses, handler(request))
		}
	}

	// Successfully end
	return nil
}

var _ = Describe("Membrane", func() {
	Context("Starting the server", func() {
		Context("That tolerates missing services", func() {
			When("It is missing the gateway plugin", func() {
				membrane, _ := membrane.New(&membrane.MembraneOptions{
					TolerateMissingServices: true,
					SuppressLogs:            true,
				})
				It("Start should Panic", func() {
					Expect(membrane.Start).To(Panic())
				})
			})

			When("The Gateway plugin is available and working", func() {
				mockGateway := &MockGateway{}

				membrane, _ := membrane.New(&membrane.MembraneOptions{
					GatewayPlugin:           mockGateway,
					SuppressLogs:            true,
					TolerateMissingServices: true,
				})

				It("Start should not Panic", func() {
					Expect(membrane.Start).ToNot(Panic())
				})

				It("Mock Gateways start method should have been called", func() {
					Expect(mockGateway.started).To(BeTrue())
				})
			})
		})

		Context("That does not tolerate missing services", func() {
			mockGateway := &MockGateway{}
			When("It is missing the eventing plugin", func() {
				membrane, _ := membrane.New(&membrane.MembraneOptions{
					TolerateMissingServices: false,
					SuppressLogs:            true,
					GatewayPlugin:           mockGateway,
				})
				It("Start should Panic", func() {
					Expect(membrane.Start).To(Panic())
				})
			})

			When("It is missing the documents plugin", func() {
				mockEventingServer := &MockEventingServer{}

				membrane, _ := membrane.New(&membrane.MembraneOptions{
					EventingPlugin:          mockEventingServer,
					GatewayPlugin:           mockGateway,
					SuppressLogs:            true,
					TolerateMissingServices: false,
				})

				It("Start should Panic", func() {
					Expect(membrane.Start).To(Panic())
				})
			})

			When("It is missing the storage plugin", func() {
				mockEventingServer := &MockEventingServer{}
				mockDocumentsServer := &MockDocumentsServer{}

				membrane, _ := membrane.New(&membrane.MembraneOptions{
					EventingPlugin:          mockEventingServer,
					DocumentsPlugin:         mockDocumentsServer,
					GatewayPlugin:           mockGateway,
					SuppressLogs:            true,
					TolerateMissingServices: false,
				})

				It("Start should Panic", func() {
					Expect(membrane.Start).To(Panic())
				})
			})
		})
	})

	Context("Handling A Single Gateway Request", func() {
		var mockGateway *MockGateway
		var mb *membrane.Membrane
		BeforeEach(func() {
			mockGateway = &MockGateway{
				requests: []*sdk.NitricRequest{
					&sdk.NitricRequest{
						Context: &sdk.NitricContext{
							RequestId:   "1234",
							PayloadType: "test-payload",
							Source:      "test",
							SourceType:  sdk.Request,
						},
						ContentType: "text/plain",
						Payload:     []byte("Test Payload"),
					},
				},
			}

			mb, _ = membrane.New(&membrane.MembraneOptions{
				ChildAddress:            "localhost:8080",
				GatewayPlugin:           mockGateway,
				TolerateMissingServices: true,
				SuppressLogs:            true,
			})
		})

		When("There is no function available", func() {
			It("Should recieve a single error response", func() {
				Expect(mb.Start).ToNot(Panic())
				Expect(mockGateway.responses).To(HaveLen(1))

				response := mockGateway.responses[0]

				By("Having the 503 HTTP error code")
				Expect(response.Status).To(Equal(503))

				By("Having a Content-Type of text/plain")
				Expect(response.Headers["Content-Type"]).To(Equal("text/plain"))

				By("Containing a Body with the encountered error message")
				Expect(string(response.Body)).To(ContainSubstring("connection refused"))
			})
		})

		When("There is a function available to recieve", func() {
			var handlerFunction *MockFunction
			BeforeEach(func() {
				handlerFunction = &MockFunction{
					response: &http.Response{
						StatusCode: 200,
						Header: http.Header{
							"Content-Type": []string{"text/plain"},
						},
						// Note: This can only be read once!
						Body: ioutil.NopCloser(bytes.NewReader([]byte("Hello World!"))),
					},
				}
				// Setup the function handler here...
				http.HandleFunc("/", handlerFunction.handler)
				go (func() {
					http.ListenAndServe(fmt.Sprintf("localhost:8080"), nil)
				})()

				// FIXME: This is expensive! Need to wait for the server to start...
				time.Sleep(200 * time.Millisecond)
			})

			It("The request should be successfully handled", func() {
				Expect(mb.Start).ToNot(Panic())
				Expect(mockGateway.responses).To(HaveLen(1))

				response := mockGateway.responses[0]

				By("The handler recieving exactly one request")
				Expect(handlerFunction.requests).To(HaveLen(1))

				request := handlerFunction.requests[0]

				By("The NitricRequest being translated to a HTTP request")
				Expect(request.Header.Get("x-nitric-request-id")).To(Equal("1234"))
				Expect(request.Header.Get("x-nitric-payload-type")).To(Equal("test-payload"))
				Expect(request.Header.Get("x-nitric-source")).To(Equal("test"))
				Expect(request.Header.Get("x-nitric-source-type")).To(Equal("REQUEST"))

				// body, _ := ioutil.ReadAll(request.Body)

				// By("Passing through the given body")
				// Expect(string(body)).To(Equal("Test Payload"))

				By("Passing through the computed content-type")
				Expect(request.Header.Get("Content-Type")).To(ContainSubstring("text/plain"))

				By("Having the 200 HTTP status code")
				Expect(response.Status).To(Equal(200))

				By("Having a Content-Type returned by the handler")
				Expect(response.Headers["Content-Type"]).To(ContainSubstring("text/plain"))

				By("Containing a Body with handler response")
				Expect(string(response.Body)).To(ContainSubstring("Hello World!"))
			})
		})
	})
})
