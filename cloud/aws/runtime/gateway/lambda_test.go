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
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/cloud/aws/mocks/provider"
	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	lambda_service "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	ep "github.com/nitrictech/nitric/core/pkg/plugins/events"
	"github.com/nitrictech/nitric/core/pkg/triggers"
	"github.com/nitrictech/nitric/core/pkg/worker"
	mock_worker "github.com/nitrictech/nitric/core/tests/mocks/worker"
)

type MockLambdaRuntime struct {
	lambda_service.LambdaRuntimeHandler
	// FIXME: Make this a union array of stuff to send....
	eventQueue []interface{}
}

func (m *MockLambdaRuntime) Start(handler interface{}) {
	// cast the function type to what we know it will be
	typedFunc := handler.(func(ctx context.Context, data map[string]interface{}) (interface{}, error))
	for _, event := range m.eventQueue {
		bytes, _ := json.Marshal(event)
		evt := map[string]interface{}{}

		err := json.Unmarshal(bytes, &evt)
		Expect(err).To(BeNil())

		// Unmarshal the thing into the event type we expect...
		// TODO: Do something with out results here...
		_, err = typedFunc(context.TODO(), evt)
		Expect(err).To(BeNil())
	}
}

var _ = Describe("Lambda", func() {
	pool := worker.NewProcessPool(&worker.ProcessPoolOptions{})

	mockHandler := mock_worker.NewMockWorker(&mock_worker.MockWorkerOptions{
		ReturnHttp: &triggers.HttpResponse{
			Body:       []byte("success"),
			StatusCode: 200,
		},
	})
	err := pool.AddWorker(mockHandler)
	Expect(err).NotTo(HaveOccurred())

	AfterEach(func() {
		mockHandler.Reset()
	})

	Context("Http Events", func() {
		When("Sending a compliant HTTP Event", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsProvider(ctrl)

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.APIGatewayV2HTTPRequest{
					Headers: map[string]string{
						"User-Agent":            "Test",
						"x-nitric-payload-type": "TestPayload",
						"x-nitric-request-id":   "test-request-id",
						"Content-Type":          "text/plain",
					},
					RawPath:        "/test/test",
					RawQueryString: "key=test&key2=test1&key=test2",
					Body:           "Test Payload",
					RequestContext: events.APIGatewayV2HTTPRequestContext{
						HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
							Method: "GET",
						},
					},
					Cookies: []string{"test1=testcookie1", "test2=testcookie2"},
				}},
			}

			client, err := lambda_service.NewWithRuntime(mockProvider, runtime.Start)
			Expect(err).To(BeNil())

			// This function will block which means we don't need to wait on processing,
			// the function will unblock once processing has finished, this is due to our mock
			// handler only looping once over each request
			It("The gateway should translate into a standard NitricRequest", func() {
				err := client.Start(pool)
				Expect(err).To(BeNil())

				By("Handling a single HTTP request")
				Expect(len(mockHandler.ReceivedRequests)).To(Equal(1))

				request := mockHandler.ReceivedRequests[0]

				By("Retaining the body")
				Expect(string(request.Body)).To(BeEquivalentTo("Test Payload"))
				By("Retaining the Headers")
				Expect(request.Header["User-Agent"][0]).To(Equal("Test"))
				Expect(request.Header["x-nitric-payload-type"][0]).To(Equal("TestPayload"))
				Expect(request.Header["x-nitric-request-id"][0]).To(Equal("test-request-id"))
				Expect(request.Header["Content-Type"][0]).To(Equal("text/plain"))
				Expect(request.Header["Cookie"]).To(Equal([]string{"test1=testcookie1", "test2=testcookie2"}))
				By("Retaining the method")
				Expect(request.Method).To(Equal("GET"))
				By("Retaining the path")
				Expect(request.Path).To(Equal("/test/test"))

				By("Retaining the query parameters")
				Expect(request.Query["key"]).To(BeEquivalentTo([]string{"test", "test2"}))
				Expect(request.Query["key2"]).To(BeEquivalentTo([]string{"test1"}))
			})
		})
	})

	Context("SNS Events", func() {
		When("The Lambda Gateway receives SNS events", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsProvider(ctrl)

			topicName := "MyTopic"
			eventPayload := map[string]interface{}{
				"test": "test",
			}
			// eventBytes, _ := json.Marshal(&eventPayload)

			event := ep.NitricEvent{
				ID:          "test-request-id",
				PayloadType: "test-payload",
				Payload:     eventPayload,
			}

			messageBytes, err := json.Marshal(&event)
			Expect(err).To(BeNil())

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.SNSEvent{
					Records: []events.SNSEventRecord{
						{
							EventVersion:         "",
							EventSource:          "aws:sns",
							EventSubscriptionArn: "some:arbitrary:subscription:arn:MySubscription",
							SNS: events.SNSEntity{
								TopicArn: fmt.Sprintf("some:arbitrary:topic:arn:%s", topicName),
								Message:  string(messageBytes),
							},
						},
					},
				}},
			}

			client, err := lambda_service.NewWithRuntime(mockProvider, runtime.Start)
			Expect(err).To(BeNil())

			It("The gateway should translate into a standard NitricRequest", func() {
				By("having the topic available")
				mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Topic).Return(map[string]string{
					"MyTopic": "some:arbitrary:topic:arn:MyTopic",
				}, nil)

				// This function will block which means we don't need to wait on processing,
				// the function will unblock once processing has finished, this is due to our mock
				// handler only looping once over each request
				err := client.Start(pool)
				Expect(err).To(BeNil())

				By("Handling a single event")
				Expect(len(mockHandler.ReceivedEvents)).To(Equal(1))

				request := mockHandler.ReceivedEvents[0]

				By("Containing the Source Topic")
				Expect(request.Topic).To(Equal("MyTopic"))
			})
		})
	})
})
