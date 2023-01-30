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

package http_service_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/cloud/azure/mocks/provider"
	"github.com/nitrictech/nitric/cloud/azure/runtime/core"
	http_service "github.com/nitrictech/nitric/cloud/azure/runtime/gateway"
	mock_pool "github.com/nitrictech/nitric/core/mocks/pool"
	mock_worker "github.com/nitrictech/nitric/core/mocks/worker"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Http", func() {
	ctrl := gomock.NewController(GinkgoT())
	pool := mock_pool.NewMockWorkerPool(ctrl)

	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	// Set this to loopback to ensure its not public in our CI/Testing environments
	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	provider := mock_provider.NewMockAzProvider(ctrl)

	provider.EXPECT().GetResources(gomock.Any(), core.AzResource_Topic).AnyTimes().Return(map[string]core.AzGenericResource{
		"test": {
			Name:     "test",
			Type:     "topic",
			Location: "eastus2",
		},
	}, nil)

	httpPlugin, err := http_service.New(provider)
	Expect(err).To(BeNil())
	// Run on a non-blocking thread
	go func(gw gateway.GatewayService) {
		defer GinkgoRecover()
		_ = gw.Start(pool)
	}(httpPlugin)

	// Delay to allow the HTTP server to correctly start
	// FIXME: Should block on channels...
	time.Sleep(1000 * time.Millisecond)

	When("Invoking the Azure AppService HTTP Gateway", func() {
		When("with a standard Nitric Request", func() {
			payload := []byte("Test")
			ctrl := gomock.NewController(GinkgoT())
			mockHandler := mock_worker.NewMockWorker(ctrl)

			It("Should be handled successfully", func() {
				By("Returning the expected worker")
				pool.EXPECT().GetWorker(gomock.Any()).Return(mockHandler, nil)

				var capturedRequest *v1.TriggerRequest

				By("Handling exactly 1 request")
				// TODO: Capture and validate payload
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), gomock.Any()).DoAndReturn(func(arg0 interface{}, arg1 interface{}) (*v1.TriggerResponse, error) {
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

				By("Generating a HTTP trigger")
				request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test/", gatewayUrl), bytes.NewReader(payload))
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Having the provided path")
				Expect(capturedRequest.GetHttp().Path).To((Equal("/test/")))

				By("Having expected headers")
				Expect(capturedRequest.GetHttp().Headers["X-Nitric-Request-Id"].Value[0]).To((Equal("1234")))
				Expect(capturedRequest.GetHttp().Headers["X-Nitric-Payload-Type"].Value[0]).To((Equal("Test Payload")))
				Expect(capturedRequest.GetHttp().Headers["User-Agent"].Value[0]).To((Equal("Test")))

				ctrl.Finish()
			})
		})

		When("With a SubscriptionValidation event", func() {
			It("Should return the provided validation code", func() {
				validationCode := "test"
				payload := map[string]string{
					"ValidationCode": validationCode,
				}
				evt := []eventgrid.Event{
					{
						Data: payload,
					},
				}

				requestBody, err := json.Marshal(evt)
				Expect(err).To(BeNil())
				request, err := http.NewRequest("POST", gatewayUrl, bytes.NewReader(requestBody))
				Expect(err).To(BeNil())
				request.Header.Add("aeg-event-type", "SubscriptionValidation")
				resp, err := http.DefaultClient.Do(request)
				Expect(err).To(BeNil())

				By("Returning a 200 response")
				Expect(resp.StatusCode).To(Equal(200))

				By("Containing the provided validation code")
				var respEvt eventgrid.SubscriptionValidationResponse
				bytes, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				err = json.Unmarshal(bytes, &respEvt)
				Expect(err).To(BeNil())

				Expect(*respEvt.ValidationResponse).To(BeEquivalentTo(validationCode))
			})
		})

		When("With a Notification event", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockHandler := mock_worker.NewMockWorker(ctrl)

			It("Should successfully handle the notification", func() {
				payload := map[string]string{
					"testing": "test",
				}
				payloadBytes, _ := json.Marshal(payload)
				testTopic := "test"

				By("Returning the expected worker")
				pool.EXPECT().GetWorker(gomock.Any()).AnyTimes().Return(mockHandler, nil)

				By("Handling exactly 1 request")
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), &v1.TriggerRequest{
					Data: payloadBytes,
					Context: &v1.TriggerRequest_Topic{
						Topic: &v1.TopicTriggerContext{
							Topic: testTopic,
						},
					},
				})

				testID := "1234"
				evt := []eventgrid.Event{
					{
						ID:    &testID,
						Topic: &testTopic,
						Data:  payload,
					},
				}

				requestBody, err := json.Marshal(evt)
				Expect(err).To(BeNil())
				request, err := http.NewRequest("POST", gatewayUrl, bytes.NewReader(requestBody))
				Expect(err).To(BeNil())
				request.Header.Add("aeg-event-type", "Notification")
				_, _ = http.DefaultClient.Do(request)
			})
		})
	})
})
