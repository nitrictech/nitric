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
	"github.com/nitrictech/nitric/core/pkg/plugins/gateway"
	"github.com/nitrictech/nitric/core/pkg/triggers"
	"github.com/nitrictech/nitric/core/pkg/worker"
	mock_worker "github.com/nitrictech/nitric/core/tests/mocks/worker"
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
			Body:       []byte("Testing Response"),
			StatusCode: 200,
		},
	})
	err := pool.AddWorker(mockHandler)
	Expect(err).To(BeNil())

	ctrl := gomock.NewController(GinkgoT())
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
		_ = gw.Start(pool)
	}(httpPlugin)

	// Delay to allow the HTTP server to correctly start
	// FIXME: Should block on channels...
	time.Sleep(1000 * time.Millisecond)

	AfterEach(func() {
		mockHandler.Reset()
	})

	When("Invoking the Azure AppService HTTP Gateway", func() {
		When("with a standard Nitric Request", func() {
			It("Should be handled successfully", func() {
				request, _ := http.NewRequest("POST", fmt.Sprintf("%s/test/", gatewayUrl), bytes.NewReader([]byte("Test")))
				request.Header.Add("x-nitric-request-id", "1234")
				request.Header.Add("x-nitric-payload-type", "Test Payload")
				request.Header.Add("User-Agent", "Test")
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				By("Handling exactly 1 request")
				Expect(mockHandler.ReceivedRequests).To(HaveLen(1))

				handledRequest := mockHandler.ReceivedRequests[0]

				By("Having the provided path")
				Expect(handledRequest.Path).To((Equal("/test/")))
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

				By("Not invoking the nitric application")
				Expect(mockHandler.ReceivedRequests).To(BeEmpty())

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

				requestBody, err := json.Marshal(evt)
				Expect(err).To(BeNil())
				request, err := http.NewRequest("POST", gatewayUrl, bytes.NewReader(requestBody))
				Expect(err).To(BeNil())
				request.Header.Add("aeg-event-type", "Notification")
				_, _ = http.DefaultClient.Do(request)

				By("Passing the event to the Nitric Application")
				Expect(mockHandler.ReceivedEvents).To(HaveLen(1))

				event := mockHandler.ReceivedEvents[0]
				By("Having the provided requestId")
				Expect(event.ID).To(Equal("1234"))

				By("Having the provided topic")
				Expect(event.Topic).To(Equal("test"))

				By("Having the provided payload")
				Expect(event.Payload).To(BeEquivalentTo(payloadBytes))
			})
		})
	})
})
