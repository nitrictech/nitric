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
	"strings"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	mock_provider "github.com/nitrictech/nitric/cloud/gcp/mocks/provider"
	cloudrun_service "github.com/nitrictech/nitric/cloud/gcp/runtime/gateway"
	mock_apis "github.com/nitrictech/nitric/core/mocks/workers/apis"
	mock_topics "github.com/nitrictech/nitric/core/mocks/workers/topics"
	"github.com/nitrictech/nitric/core/pkg/gateway"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

const GATEWAY_ADDRESS = "127.0.0.1:9001"

var _ = Describe("Http", func() {
	defer GinkgoRecover()

	ctrl := gomock.NewController(GinkgoT())
	gatewayUrl := fmt.Sprintf("http://%s", GATEWAY_ADDRESS)

	mockApiRequestHandler := mock_apis.NewMockApiRequestHandler(ctrl)
	mockTopicRequestHandler := mock_topics.NewMockSubscriptionRequestHandler(ctrl)

	// Set this to loopback to ensure its not public in our CI/Testing environments
	BeforeSuite(func() {
		os.Setenv("GATEWAY_ADDRESS", GATEWAY_ADDRESS)
	})

	provider := mock_provider.NewMockGcpResourceResolver(ctrl)

	httpPlugin, err := cloudrun_service.New(provider)
	Expect(err).To(BeNil())

	// Run on a non-blocking thread
	go func(gw gateway.GatewayService) {
		defer GinkgoRecover()
		_ = gw.Start(&gateway.GatewayStartOpts{
			ApiPlugin:            mockApiRequestHandler,
			TopicsListenerPlugin: mockTopicRequestHandler,
		})
	}(httpPlugin)

	// Give the gateway time to start, ideally we would block on a channel
	time.Sleep(500 * time.Millisecond)

	When("Invoking the GCP HTTP Gateway", func() {
		When("with a HTTP request", func() {
			It("Should be handled successfully", func() {
				payload := []byte("Test")

				var capturedRequest *apispb.ServerMessage
				var capturedApiName string

				By("Handling exactly 1 request")
				mockApiRequestHandler.EXPECT().HandleRequest(gomock.Any(), gomock.Any()).Times(1).DoAndReturn(func(arg0 interface{}, arg1 interface{}) (*apispb.ClientMessage, error) {
					capturedApiName = arg0.(string)
					capturedRequest = arg1.(*apispb.ServerMessage)
					// apiName string, request *apispb.ServerMessage

					return &apispb.ClientMessage{
						Id: "test",
						Content: &apispb.ClientMessage_HttpResponse{
							HttpResponse: &apispb.HttpResponse{
								Status: 200,
								Body:   []byte("success"),
							},
						},
					}, nil
				})

				request, err := http.NewRequest("POST", fmt.Sprintf("%s/x-nitric-api/%s/test", gatewayUrl, "test-api"), bytes.NewReader(payload))
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

				By("Routing to the correct api")
				Expect(capturedApiName).To(Equal("test-api"))

				By("Trimming the api name off the request path")
				Expect(strings.Contains(capturedRequest.GetHttpRequest().Path, capturedApiName)).To(BeFalse())

				By("Preserving the original requests method")
				Expect(capturedRequest.GetHttpRequest().Method).To(Equal("POST"))

				By("Preserving the original requests path")
				Expect(capturedRequest.GetHttpRequest().Path).To(Equal("test"))

				By("Preserving the original requests body")
				Expect(capturedRequest.GetHttpRequest().Body).To(BeEquivalentTo([]byte("Test")))

				By("Preserving the original requests headers")
				Expect(capturedRequest.GetHttpRequest().Headers["User-Agent"].Value[0]).To(Equal("Test"))
				Expect(capturedRequest.GetHttpRequest().Headers["X-Nitric-Request-Id"].Value[0]).To(Equal("1234"))
				Expect(capturedRequest.GetHttpRequest().Headers["X-Nitric-Payload-Type"].Value[0]).To(Equal("Test Payload"))
				Expect(capturedRequest.GetHttpRequest().Headers["Cookie"].Value[0]).To(Equal("test1=testcookie1; test2=testcookie2"))

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("success"))
			})
		})

		When("From a subscription with a NitricEvent", func() {
			content, _ := structpb.NewStruct(map[string]interface{}{
				"Test": "Test",
			})

			message := topicspb.TopicMessage{
				Content: &topicspb.TopicMessage_StructPayload{
					StructPayload: content,
				},
			}

			messageBytes, err := proto.Marshal(&message)
			Expect(err).To(BeNil())

			b64Event := base64.StdEncoding.EncodeToString(messageBytes)
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
				var capturedRequest *topicspb.ServerMessage

				By("Handling exactly 1 request")
				mockTopicRequestHandler.EXPECT().HandleRequest(gomock.Any()).Times(1).DoAndReturn(func(arg0 interface{}) (*topicspb.ClientMessage, error) {
					capturedRequest = arg0.(*topicspb.ServerMessage)

					return &topicspb.ClientMessage{
						Id: "test",
						Content: &topicspb.ClientMessage_MessageResponse{
							MessageResponse: &topicspb.MessageResponse{
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

				capturedMessageBytes, err := proto.Marshal(capturedRequest.GetMessageRequest().GetMessage())
				Expect(err).To(BeNil())

				By("Passing through the published message data")
				Expect(capturedMessageBytes).To(Equal(messageBytes))

				By("The request returns a successful status")
				Expect(resp.StatusCode).To(Equal(200))

				By("Returning the expected output")
				Expect(string(responseBody)).To(Equal("success"))
			})
		})
	})
})
