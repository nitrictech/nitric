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
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/eventgrid/eventgrid"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/mocks/provider"
	mock_worker "github.com/nitrictech/nitric/mocks/worker"
	"github.com/nitrictech/nitric/pkg/plugins/gateway"
	http_service "github.com/nitrictech/nitric/pkg/plugins/gateway/appservice"
	"github.com/nitrictech/nitric/pkg/providers/azure/core"
	"github.com/nitrictech/nitric/pkg/triggers"
	"github.com/nitrictech/nitric/pkg/worker"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Http", func() {
	pool := worker.NewProcessPool(&worker.ProcessPoolOptions{})

	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	// Set this to loopback to ensure its not public in our CI/Testing environments
	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	ctrl := gomock.NewController(GinkgoT())
	provider := mock_provider.NewMockAzProvider(ctrl)

	provider.EXPECT().GetResources(core.AzResource_Topic).AnyTimes().Return(map[string]core.AzGenericResource{
		"test": {
			Name:     "test",
			Type:     "topic",
			Location: "eastus2",
		},
	}, nil)

	httpPlugin, err := http_service.New()
	Expect(err).To(BeNil())
	// Run on a non-blocking thread
	go func(gw gateway.GatewayService) {
		_ = gw.Start(pool)
	}(httpPlugin)

	time.Sleep(500 * time.Millisecond)

	When("Invoking the Azure AppService HTTP Gateway", func() {
		Context("with a standard Nitric Request", func() {
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
				By("Receiving the expected request")
				hndlr.EXPECT().HandleHttpRequest(gomock.Any()).Return(&triggers.HttpResponse{
					StatusCode: 200,
				}, nil).Times(1)

				request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test/", gatewayUrl), bytes.NewReader([]byte("Test")))
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())
			})
		})

		Context("With a SubscriptionValidation event", func() {
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
				request, err := http.NewRequest("POST", fmt.Sprintf("%s/x-nitric-subscription/test", gatewayUrl), bytes.NewReader(requestBody))
				Expect(err).To(BeNil())
				request.Header.Add("aeg-event-type", "SubscriptionValidation")
				resp, err := http.DefaultClient.Do(request)
				Expect(err).To(BeNil())

				By("Returning a 200 response")
				Expect(resp.StatusCode).To(Equal(200))

				By("Containing the provided validation code")
				var respEvt eventgrid.SubscriptionValidationResponse
				bytes, err := ioutil.ReadAll(resp.Body)
				Expect(err).To(BeNil())

				err = json.Unmarshal(bytes, &respEvt)
				Expect(err).To(BeNil())

				Expect(*respEvt.ValidationResponse).To(BeEquivalentTo(validationCode))
			})
		})

		Context("With a Notification event", func() {
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

			It("Should successfully handle the notification", func() {
				payload := map[string]string{
					"testing": "test",
				}
				payloadBytes, _ := json.Marshal(payload)
				testTopic := "test"
				testID := "1234"
				evt := []eventgrid.Event{
					{
						ID:    &testID,
						Topic: &testTopic,
						Data:  payload,
					},
				}

				By("The handler receiving the request")
				hndlr.EXPECT().HandleEvent(&triggers.Event{
					ID:      "1234",
					Topic:   "test",
					Payload: payloadBytes,
				}).Times(1).Return(nil)

				requestBody, err := json.Marshal(evt)
				Expect(err).To(BeNil())
				request, err := http.NewRequest("POST", fmt.Sprintf("%s/x-nitric-subscription/test", gatewayUrl), bytes.NewReader(requestBody))
				Expect(err).To(BeNil())
				request.Header.Add("aeg-event-type", "Notification")
				_, _ = http.DefaultClient.Do(request)
			})
		})
	})
})
