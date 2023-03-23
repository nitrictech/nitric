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

package gateway_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/cloud/gcp/mocks/provider"
	cloudrun_plugin "github.com/nitrictech/nitric/cloud/gcp/runtime/gateway"
	mock_worker "github.com/nitrictech/nitric/core/mocks/worker"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/events"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/worker/pool"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Http", func() {
	defer GinkgoRecover()

	ctrl := gomock.NewController(GinkgoT())
	pool := pool.NewProcessPool(&pool.ProcessPoolOptions{})
	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	// Set this to loopback to ensure its not public in our CI/Testing environments
	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	mockHandler := mock_worker.NewMockWorker(ctrl)
	mockHandler.EXPECT().HandlesTrigger(gomock.Any()).AnyTimes().Return(true)

	err := pool.AddWorker(mockHandler)
	Expect(err).To(BeNil())

	provider := mock_provider.NewMockGcpProvider(ctrl)

	httpPlugin, err := cloudrun_plugin.New(provider)
	Expect(err).To(BeNil())

	// Run on a non-blocking thread
	go func(gw gateway.GatewayService) {
		defer GinkgoRecover()
		_ = gw.Start(pool)
	}(httpPlugin)

	// Delay to allow the HTTP server to correctly start
	// FIXME: Should block on channels...
	time.Sleep(500 * time.Millisecond)

	When("Invoking the GCP HTTP Gateway", func() {
		When("with a HTTP request", func() {
			It("Should be handled successfully", func() {
				payload := []byte("Test")

				var capturedRequest *v1.TriggerRequest

				By("Handling exactly 1 request")
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(arg0 interface{}, arg1 interface{}) (*v1.TriggerResponse, error) {
					capturedRequest = arg1.(*v1.TriggerRequest)

					return &v1.TriggerResponse{
						Data: []byte("success"),
						Context: &v1.TriggerResponse_Http{
							Http: &v1.HttpResponseContext{
								Status: 200,
							},
						},
					}, nil
				})

				request, err := http.NewRequest("POST", fmt.Sprintf("%s/test", gatewayUrl), bytes.NewReader(payload))
				Expect(err).To(BeNil())
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				request.Header.Add("Cookie", "test1=testcookie1")
				request.Header.Add("Cookie", "test2=testcookie2")
				resp, err := http.DefaultClient.Do(request)

				responseBody := make([]byte, 0)

				if err == nil {
					if bytes, err := io.ReadAll(resp.Body); err == nil {
						responseBody = bytes
					}
				}

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Preserving the original requests method")
				Expect(capturedRequest.GetHttp().Method).To(Equal("POST"))

				By("Preserving the original requests path")
				Expect(capturedRequest.GetHttp().Path).To(Equal("/test"))

				By("Preserving the original requests body")
				Expect(capturedRequest.Data).To(BeEquivalentTo([]byte("Test")))

				By("Preserving the original requests headers")
				Expect(capturedRequest.GetHttp().Headers["User-Agent"].Value[0]).To(Equal("Test"))
				Expect(capturedRequest.GetHttp().Headers["X-Nitric-Request-Id"].Value[0]).To(Equal("1234"))
				Expect(capturedRequest.GetHttp().Headers["X-Nitric-Payload-Type"].Value[0]).To(Equal("Test Payload"))
				Expect(capturedRequest.GetHttp().Headers["Cookie"].Value[0]).To(Equal("test1=testcookie1; test2=testcookie2"))

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
				var capturedRequest *v1.TriggerRequest

				By("Handling exactly 1 request")
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(arg0 interface{}, arg1 interface{}) (*v1.TriggerResponse, error) {
					capturedRequest = arg1.(*v1.TriggerRequest)

					return &v1.TriggerResponse{
						Data: []byte("success"),
						Context: &v1.TriggerResponse_Topic{
							Topic: &v1.TopicResponseContext{
								Success: true,
							},
						},
					}, nil
				})

				request, err := http.NewRequest("POST", fmt.Sprintf("%s/x-nitric-topic/test", gatewayUrl), bytes.NewReader(payloadBytes))
				Expect(err).To(BeNil())
				request.Header.Add("Content-Type", "application/json")
				resp, err := http.DefaultClient.Do(request)
				Expect(err).To(BeNil())
				responseBody, _ := io.ReadAll(resp.Body)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Passing through the published message data")
				Expect(capturedRequest.Data).To(BeEquivalentTo(eventBytes))

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("success"))
			})
		})
	})
})
