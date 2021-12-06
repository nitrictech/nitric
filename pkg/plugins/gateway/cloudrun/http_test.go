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

package cloudrun_plugin_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/pkg/plugins/events"
	cloudrun_plugin "github.com/nitrictech/nitric/pkg/plugins/gateway/cloudrun"
	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/worker"
	mock_worker "github.com/nitrictech/nitric/tests/mocks/worker"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Http", func() {
	pool := worker.NewProcessPool(&worker.ProcessPoolOptions{})
	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	// Set this to loopback to ensure its not public in our CI/Testing environments
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
	httpPlugin, _ := cloudrun_plugin.New()
	// Run on a non-blocking thread
	go (httpPlugin.Start)(pool)

	// Delay to allow the HTTP server to correctly start
	// FIXME: Should block on channels...
	time.Sleep(500 * time.Millisecond)

	AfterEach(func() {
		mockHandler.Reset()
	})

	When("Invoking the GCP HTTP Gateway", func() {
		When("with a HTTP request", func() {

			It("Should be handled successfully", func() {
				request, err := http.NewRequest("POST", fmt.Sprintf("%s/test", gatewayUrl), bytes.NewReader([]byte("Test")))
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				request.Header.Add("Cookie", "test1=testcookie1")
				request.Header.Add("Cookie", "test2=testcookie2")
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
				Expect(mockHandler.ReceivedRequests).To(HaveLen(1))

				handledRequest := mockHandler.ReceivedRequests[0]
				By("Preserving the original requests method")
				Expect(handledRequest.Method).To(Equal("POST"))

				By("Preserving the original requests path")
				Expect(handledRequest.Path).To(Equal("/test"))

				streamRead := handledRequest.Body
				By("Preserving the original requests body")
				Expect(streamRead).To(BeEquivalentTo([]byte("Test")))

				By("Preserving the original requests headers")
				Expect(string(handledRequest.Header["User-Agent"][0])).To(Equal("Test"))
				Expect(string(handledRequest.Header["X-Nitric-Request-Id"][0])).To(Equal("1234"))
				Expect(string(handledRequest.Header["X-Nitric-Payload-Type"][0])).To(Equal("Test Payload"))
				Expect(handledRequest.Header["Cookie"]).To(Equal([]string{"test1=testcookie1; test2=testcookie2"}))

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("success"))
			})
		})

		When("From a subcription with a NitricEvent", func() {
			eventPayload := map[string]interface{}{
				"Test": "Test",
			}
			eventBytes, _ := json.Marshal(&events.NitricEvent{
				ID:          "1234",
				PayloadType: "Test Payload",
				Payload:     eventPayload,
			})

			b64Event := base64.StdEncoding.EncodeToString(eventBytes)

			payloadBytes, _ := json.Marshal(&map[string]interface{}{
				"subscription": "test",
				"message": map[string]interface{}{
					"attributes": map[string]string{
						"x-nitric-topic": "test",
					},
					"id":   "test",
					"data": b64Event,
				},
			})

			It("Should handle the event successfully", func() {
				request, err := http.NewRequest("POST", gatewayUrl, bytes.NewReader(payloadBytes))
				request.Header.Add("Content-Type", "application/json")
				resp, err := http.DefaultClient.Do(request)
				responseBody, _ := ioutil.ReadAll(resp.Body)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Handling exactly 1 request")
				Expect(mockHandler.ReceivedEvents).To(HaveLen(1))

				handledEvent := mockHandler.ReceivedEvents[0]

				By("Passing through the pubsub message ID")
				Expect(handledEvent.ID).To(Equal("test"))

				By("Extracting the topic name from the subscription")
				Expect(handledEvent.Topic).To(Equal("test"))

				By("Passing through the published message data")
				Expect(handledEvent.Payload).To(BeEquivalentTo(eventBytes))

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("success"))
			})
		})
	})
})
