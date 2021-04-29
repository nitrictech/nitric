package membrane_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/membrane"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/triggers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockEventingServer struct {
	sdk.UnimplementedEventingPlugin
}

type MockStorageServer struct {
	sdk.UnimplementedStoragePlugin
}

type MockKeyValueServer struct {
	sdk.UnimplementedKeyValuePlugin
}

type MockQueueServer struct {
	sdk.UnimplementedQueuePlugin
}

type MockAuthServer struct {
	sdk.UnimplementedAuthPlugin
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
	triggers []triggers.Trigger
	// store responses for inspection
	responses []*triggers.HttpResponse
	started   bool
}

func (gw *MockGateway) Start(handler handler.TriggerHandler) error {
	// Spy on the mock gateway
	gw.responses = make([]*triggers.HttpResponse, 0)

	gw.started = true
	if gw.triggers != nil {
		for _, trigger := range gw.triggers {
			if s, ok := trigger.(*triggers.HttpRequest); ok {
				resp, err := handler.HandleHttpRequest(s)

				if err != nil {
					gw.responses = append(gw.responses, &triggers.HttpResponse{
						StatusCode: 500,
						Body:       []byte(err.Error()),
					})
				} else {
					gw.responses = append(gw.responses, resp)
				}
			} else if s, ok := trigger.(*triggers.Event); ok {
				handler.HandleEvent(s)
			}
		}
	}

	// Successfully end
	return nil
}

var _ = Describe("Membrane", func() {
	Context("New", func() {
		Context("Tolerate Missing Services is enabled", func() {
			When("The gateway plugin is missing", func() {
				It("Should still fail to create", func() {
					m, err := membrane.New(&membrane.MembraneOptions{
						SuppressLogs:            true,
						TolerateMissingServices: true,
					})
					Expect(err).Should(HaveOccurred())
					Expect(m).To(BeNil())
				})
			})

			When("The gateway plugin is present", func() {
				mockGateway := &MockGateway{}
				mbraneOpts := membrane.MembraneOptions{
					SuppressLogs:            true,
					GatewayPlugin:           mockGateway,
					TolerateMissingServices: true,
				}
				It("Should successfully create the membrane server", func() {
					m, err := membrane.New(&mbraneOpts)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(m).ToNot(BeNil())
				})
			})
		})

		Context("Tolerate Missing Services is disabled", func() {
			When("Only the gateway plugin is present", func() {
				mockGateway := &MockGateway{}
				mbraneOpts := membrane.MembraneOptions{
					TolerateMissingServices: false,
					SuppressLogs:            true,
					GatewayPlugin:           mockGateway,
				}
				It("Should fail to create", func() {
					m, err := membrane.New(&mbraneOpts)
					Expect(err).Should(HaveOccurred())
					Expect(m).To(BeNil())
				})
			})

			When("All plugins are present", func() {
				mockEventingServer := &MockEventingServer{}
				mockKeyValueServer := &MockKeyValueServer{}
				mockStorageServer := &MockStorageServer{}
				mockQueueServer := &MockQueueServer{}
				mockAuthServer := &MockAuthServer{}

				mockGateway := &MockGateway{}
				mbraneOpts := membrane.MembraneOptions{
					TolerateMissingServices: false,
					SuppressLogs:            true,
					GatewayPlugin:           mockGateway,
					EventingPlugin:          mockEventingServer,
					KvPlugin:                mockKeyValueServer,
					StoragePlugin:           mockStorageServer,
					QueuePlugin:             mockQueueServer,
					AuthPlugin:              mockAuthServer,
				}

				It("Should successfully create the membrane server", func() {
					m, err := membrane.New(&mbraneOpts)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(m).ToNot(BeNil())
				})
			})
		})
	})

	Context("Starting the server", func() {
		Context("That tolerates missing adapters", func() {
			When("The Gateway plugin is available and working", func() {
				mockGateway := &MockGateway{}

				membrane, _ := membrane.New(&membrane.MembraneOptions{
					GatewayPlugin:           mockGateway,
					SuppressLogs:            true,
					TolerateMissingServices: true,
				})

				It("Start should not error", func() {
					err := membrane.Start()
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("Mock Gateways start method should have been called", func() {
					Expect(mockGateway.started).To(BeTrue())
				})
			})
		})

		When("The configured service port is already consumed", func() {
			mockGateway := &MockGateway{}
			var lis net.Listener

			membrane, _ := membrane.New(&membrane.MembraneOptions{
				GatewayPlugin:           mockGateway,
				SuppressLogs:            true,
				TolerateMissingServices: true,
				ServiceAddress:          "localhost:9005",
			})

			BeforeEach(func() {
				lis, _ = net.Listen("tcp", "localhost:9005")
			})

			AfterEach(func() {
				lis.Close()
			})

			It("Should return an error", func() {
				err := membrane.Start()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Could not listen"))
			})
		})
	})

	Context("Starting the child process", func() {
		var mockGateway *MockGateway
		var mb *membrane.Membrane
		When("The configured command exists", func() {
			BeforeEach(func() {
				mockGateway = &MockGateway{}

				mb, _ = membrane.New(&membrane.MembraneOptions{
					ChildAddress:            "localhost:8081",
					ChildCommand:            "echo",
					GatewayPlugin:           mockGateway,
					ChildTimeoutSeconds:     1,
					TolerateMissingServices: true,
					SuppressLogs:            true,
				})
			})

			When("There is nothing listening on ChildAddress", func() {
				It("Should return an error", func() {
					err := mb.Start()
					Expect(err).Should(HaveOccurred())
				})
			})

			When("There is something listening on childAddress", func() {
				BeforeEach(func() {
					go (func() {
						http.ListenAndServe(fmt.Sprintf("localhost:8081"), nil)
					})()
				})

				AfterEach(func() {

				})

				It("Should wait for the service to start", func() {
					err := mb.Start()
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		When("The configured command does not exist", func() {
			BeforeEach(func() {
				mockGateway = &MockGateway{}

				mb, _ = membrane.New(&membrane.MembraneOptions{
					ChildAddress:            "localhost:808",
					ChildCommand:            "fakecommand",
					GatewayPlugin:           mockGateway,
					TolerateMissingServices: true,
					SuppressLogs:            true,
				})
			})

			It("Should return an error", func() {
				err := mb.Start()
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Operating in FaaS Mode", func() {
		Context("Handling A Single HttpRequest", func() {
			var mockGateway *MockGateway
			var mb *membrane.Membrane
			BeforeEach(func() {
				mockGateway = &MockGateway{
					triggers: []triggers.Trigger{
						&triggers.HttpRequest{
							Body:   []byte("Test Payload"),
							Path:   "/test/",
							Header: make(map[string]string),
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
					err := mb.Start()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(mockGateway.responses).To(HaveLen(1))

					response := mockGateway.responses[0]

					By("Having the 500 HTTP error code")
					Expect(response.StatusCode).To(Equal(500))

					By("Containing a Body with the error message from the gateway")
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
					err := mb.Start()
					Expect(err).ShouldNot(HaveOccurred())
					Expect(mockGateway.responses).To(HaveLen(1))

					response := mockGateway.responses[0]

					By("The handler recieving exactly one request")
					Expect(handlerFunction.requests).To(HaveLen(1))

					request := handlerFunction.requests[0]

					By("By consuming the path of the request")
					Expect(request.URL.String()).To(ContainSubstring("/"))

					By("Having the 200 HTTP status code")
					Expect(response.StatusCode).To(Equal(200))

					By("Having a Content-Type returned by the handler")
					Expect(response.Header.Peek("Content-Type")).To(ContainSubstring("text/plain"))

					By("Containing a Body with handler response")
					Expect(string(response.Body)).To(ContainSubstring("Hello World!"))
				})
			})
		})
	})
})
