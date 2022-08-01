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

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_worker "github.com/nitrictech/nitric/mocks/worker"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	"github.com/nitrictech/nitric/pkg/plugins/gateway"
	cloudrun_plugin "github.com/nitrictech/nitric/pkg/plugins/gateway/cloudrun"
	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/worker"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Http", func() {
	pool := worker.NewProcessPool(&worker.ProcessPoolOptions{
		MinWorkers: 0,
	})
	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	// Set this to loopback to ensure its not public in our CI/Testing environments
	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	httpPlugin, err := cloudrun_plugin.New()
	Expect(err).To(BeNil())

	// Run on a non-blocking thread
	go func(gw gateway.GatewayService) {
		_ = gw.Start(pool)
	}(httpPlugin)

	// Delay to allow the HTTP server to correctly start
	// FIXME: Should block on channels...
	time.Sleep(500 * time.Millisecond)

	When("Invoking the GCP HTTP Gateway", func() {
		Context("with a HTTP request", func() {
			var wrkr *worker.RouteWorker
			var hndlr *mock_worker.MockAdapter
			var ctrl *gomock.Controller

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				hndlr = mock_worker.NewMockAdapter(ctrl)
				wrkr = worker.NewRouteWorker(hndlr, &worker.RouteWorkerOptions{
					Api:     "test",
					Path:    "/test",
					Methods: []string{"POST"},
				})
				_ = pool.AddWorker(wrkr)
			})

			AfterEach(func() {
				_ = pool.RemoveWorker(wrkr)
				ctrl.Finish()
			})

			It("Should be handled successfully", func() {
				By("Calling the worker with expected request")
				hndlr.EXPECT().HandleHttpRequest(gomock.AssignableToTypeOf(&triggers.HttpRequest{})).Return(&triggers.HttpResponse{
					Body:       []byte("success"),
					StatusCode: 200,
				}, nil).Times(1)

				request, err := http.NewRequest("POST", fmt.Sprintf("%s/test", gatewayUrl), bytes.NewReader([]byte("Test")))
				Expect(err).To(BeNil())
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				request.Header.Add("Cookie", "test1=testcookie1")
				request.Header.Add("Cookie", "test2=testcookie2")
				resp, err := http.DefaultClient.Do(request)

				responseBody := make([]byte, 0)

				if err == nil {
					if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
						responseBody = bytes
					}
				}

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("success"))
			})
		})

		Context("From a subcription with a NitricEvent", func() {
			eventPayload := map[string]interface{}{
				"Test": "Test",
			}
			eventBytes, _ := json.Marshal(&events.NitricEvent{
				ID:          "1234",
				PayloadType: "Test Payload",
				Payload:     eventPayload,
			})

			pBytes, _ := json.Marshal(eventPayload)

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

			var wrkr *worker.SubscriptionWorker
			var hndlr *mock_worker.MockAdapter
			var ctrl *gomock.Controller

			BeforeEach(func() {
				ctrl = gomock.NewController(GinkgoT())
				hndlr = mock_worker.NewMockAdapter(ctrl)
				wrkr = worker.NewSubscriptionWorker(hndlr, &worker.SubscriptionWorkerOptions{
					Topic: "test",
				})
				_ = pool.AddWorker(wrkr)
			})

			AfterEach(func() {
				_ = pool.RemoveWorker(wrkr)
				ctrl.Finish()
			})

			It("Should handle the event successfully", func() {
				By("Calling the handler with the expected request")
				hndlr.EXPECT().HandleEvent(&triggers.Event{
					ID:      "1234",
					Topic:   "test",
					Payload: pBytes,
				}).Times(1).Return(nil)

				request, err := http.NewRequest("POST", fmt.Sprintf("%s/x-nitric-subscription/test", gatewayUrl), bytes.NewReader(payloadBytes))
				Expect(err).To(BeNil())
				request.Header.Add("Content-Type", "application/json")
				resp, err := http.DefaultClient.Do(request)
				Expect(err).To(BeNil())
				responseBody, _ := ioutil.ReadAll(resp.Body)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("success"))
			})
		})
	})
})
