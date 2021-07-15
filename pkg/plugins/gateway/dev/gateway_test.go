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

package gateway_plugin_test

import (
	"bytes"
	"fmt"
	gateway_plugin "github.com/nitric-dev/membrane/pkg/plugins/gateway/dev"
	"github.com/nitric-dev/membrane/pkg/triggers"
	"github.com/nitric-dev/membrane/pkg/worker"
	mock_worker "github.com/nitric-dev/membrane/tests/mocks/worker"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Gateway", func() {
	pool := worker.NewProcessPool(&worker.ProcessPoolOptions{})

	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	mockHandler := mock_worker.NewMockWorker(&mock_worker.MockWorkerOptions{
		ReturnHttp: &triggers.HttpResponse{
			Body:       []byte("success"),
			StatusCode: 200,
		},
	})
	pool.AddWorker(mockHandler)

	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	gateway, _ := gateway_plugin.New()

	AfterEach(func() {
		mockHandler.Reset()
	})

	// Start the gatewat on a seperate thread so it doesn't block the tests...
	go (gateway.Start)(pool)
	// FIXME: Update gateway to block on channel...
	time.Sleep(500 * time.Millisecond)

	When("Receiving standard HTTP requests", func() {
		When("The request contains standard nitric headers", func() {
			payload := []byte("Test")
			request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test", gatewayUrl), bytes.NewReader(payload))

			request.Header.Add("x-nitric-request-id", "1234")
			request.Header.Add("x-nitric-payload-type", "test-payload")
			request.Header.Add("User-Agent", "Test")

			It("should successfully pass on the request", func() {
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Passing through exactly 1 request")
				Expect(mockHandler.ReceivedRequests).To(HaveLen(1))

				handledRequest := mockHandler.ReceivedRequests[0]

				By("Preserving the original request method")
				Expect(handledRequest.Method).To(Equal("POST"))

				By("Preserving the original request URL")
				Expect(handledRequest.Path).To(Equal("/test"))

				// FIXME: Weird bug occurring in tests,
				// need to validate genuine runtime behaviour here...
				// Seems like the original request stream is closed
				// before we can actually properly assess it
				By("Preserving the original request Body")
				Expect(string(handledRequest.Body)).To(Equal("Test"))
			})
		})
		// TODO: Handle cases of missing nitric headers
		// TODO: Handle cases of other non POST methods
	})

	When("Receiving requests from a topic subscription", func() {
		When("The request contains standard nitric headers", func() {
			payload := []byte("Test")
			request, _ := http.NewRequest("POST", gatewayUrl, bytes.NewReader(payload))

			request.Header.Add("x-nitric-request-id", "1234")
			request.Header.Add("x-nitric-payload-type", "test-payload")
			request.Header.Add("x-nitric-source-type", "SUBSCRIPTION")
			request.Header.Add("x-nitric-source", "test-topic")

			It("should successfully pass on the event", func() {
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Passing through exactly 1 event")
				Expect(mockHandler.ReceivedEvents).To(HaveLen(1))

				evt := mockHandler.ReceivedEvents[0]

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
