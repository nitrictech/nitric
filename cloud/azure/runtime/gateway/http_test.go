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

package http_service

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
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/cloud/azure/mocks/provider"
	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	mock_apis "github.com/nitrictech/nitric/core/mocks/workers/apis"
	mock_http "github.com/nitrictech/nitric/core/mocks/workers/http"
	mock_topics "github.com/nitrictech/nitric/core/mocks/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"github.com/nitrictech/nitric/test"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Http", func() {
	ctrl := gomock.NewController(GinkgoT())

	testEvtToken := "test"
	os.Setenv("EVENT_TOKEN", testEvtToken)

	// pool := mock_pool.NewMockWorkerPool(ctrl)

	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)
	// Set this to loopback to ensure its not public in our CI/Testing environments
	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	provider := mock_provider.NewMockAzResourceResolver(ctrl)

	provider.EXPECT().GetResources(gomock.Any(), resource.AzResource_Topic).AnyTimes().Return(map[string]resource.AzGenericResource{
		"test": {
			Name:     "test",
			Type:     "topic",
			Location: "eastus2",
		},
	}, nil)

	gatewayOptions := &gateway.GatewayStartOpts{
		ApiPlugin: nil,
	}

	httpPlugin, err := New(provider)
	Expect(err).To(BeNil())
	// Run on a non-blocking thread
	go func(gw gateway.GatewayService) {
		defer GinkgoRecover()
		_ = gw.Start(gatewayOptions)
	}(httpPlugin)

	// Delay to allow the HTTP server to correctly start, ideally we could block on a channel instead of waiting a fixed time
	time.Sleep(1 * time.Second)

	When("Invoking the Azure AppService HTTP Gateway", func() {
		When("with a standard Nitric Request", func() {
			payload := []byte("Test")
			ctrl := gomock.NewController(GinkgoT())

			mockManager := mock_apis.NewMockApiRequestHandler(ctrl)
			mockHttpManager := mock_http.NewMockHttpRequestHandler(ctrl)
			gatewayOptions.ApiPlugin = mockManager
			gatewayOptions.HttpPlugin = mockHttpManager

			It("Should be handled successfully", func() {
				By("Handling exactly 1 request")

				mockRequest := &apispb.ServerMessage{
					Content: &apispb.ServerMessage_HttpRequest{
						HttpRequest: &apispb.HttpRequest{
							Method: "POST",
							Path:   "test/",
							Headers: map[string]*apispb.HeaderValue{
								"Content-Length":  {Value: []string{"4"}},
								"User-Agent":      {Value: []string{"Go-http-client/1.1"}},
								"X-Forwarded-For": {Value: []string{"127.0.0.1:9001"}},
								"Accept-Encoding": {Value: []string{"gzip"}},
							},
							Body: payload,
						},
					},
				}

				mockManager.EXPECT().HandleRequest("test", test.ProtoEq(mockRequest)).Return(&apispb.ClientMessage{
					Id: "TODO",
					Content: &apispb.ClientMessage_HttpResponse{
						HttpResponse: &apispb.HttpResponse{
							Status:  200,
							Body:    []byte("Test"),
							Headers: map[string]*apispb.HeaderValue{},
						},
					},
				}, nil)

				By("Generating a HTTP trigger")
				request, _ := http.NewRequest("POST", fmt.Sprintf("%s/x-nitric-api/test/test/", gatewayUrl), bytes.NewReader(payload[:]))
				_, err := http.DefaultClient.Do(request)

				By("Not returning an error")
				Expect(err).To(BeNil())

				// By("Having the provided path")
				// Expect(capturedRequest.GetHttp().Path).To((Equal("/test/")))

				// By("Having expected headers")
				// Expect(capturedRequest.GetHttp().Headers["X-Nitric-Request-Id"].Value[0]).To((Equal("1234")))
				// Expect(capturedRequest.GetHttp().Headers["X-Nitric-Payload-Type"].Value[0]).To((Equal("Test Payload")))
				// Expect(capturedRequest.GetHttp().Headers["User-Agent"].Value[0]).To((Equal("Test")))

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
				request, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/x-nitric-topic/test", gatewayUrl, testEvtToken), bytes.NewReader(requestBody))
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

			mockManager := mock_topics.NewMockSubscriptionRequestHandler(ctrl)
			gatewayOptions.TopicsListenerPlugin = mockManager

			It("Should successfully handle the notification", func() {
				payload := map[string]interface{}{
					"testing": "test",
				}

				structPayload, _ := structpb.NewStruct(payload)

				messagePayload := &topicspb.TopicMessage{
					Content: &topicspb.TopicMessage_StructPayload{
						StructPayload: structPayload,
					},
				}

				messagePayloadBytes, _ := proto.Marshal(messagePayload)

				testTopic := "test"

				// By("Returning the expected worker")
				// pool.EXPECT().GetWorker(gomock.Any()).AnyTimes().Return(mockHandler, nil)

				mockRequest := &topicspb.ServerMessage{
					Content: &topicspb.ServerMessage_MessageRequest{
						MessageRequest: &topicspb.MessageRequest{
							TopicName: testTopic,
							Message:   messagePayload,
						},
					},
				}

				By("Handling exactly 1 request")
				mockManager.EXPECT().HandleRequest(test.ProtoEq(mockRequest)).Return(&topicspb.ClientMessage{
					Content: &topicspb.ClientMessage_MessageResponse{
						MessageResponse: &topicspb.MessageResponse{
							Success: true,
						},
					},
				}, nil)

				testID := "1234"
				evt := []eventgrid.Event{
					{
						ID:    &testID,
						Topic: &testTopic,
						Data:  messagePayloadBytes,
					},
				}

				requestBody, err := json.Marshal(evt)
				Expect(err).To(BeNil())
				request, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/x-nitric-topic/test", gatewayUrl, testEvtToken), bytes.NewReader(requestBody))
				Expect(err).To(BeNil())
				request.Header.Add("aeg-event-type", "Notification")
				_, _ = http.DefaultClient.Do(request)
			})
		})
	})
})
