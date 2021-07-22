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

package lambda_service_test

import (
	"context"
	"encoding/json"
	"fmt"

	lambda_service "github.com/nitric-dev/membrane/pkg/plugins/gateway/lambda"
	"github.com/nitric-dev/membrane/pkg/triggers"
	"github.com/nitric-dev/membrane/pkg/worker"
	mock_worker "github.com/nitric-dev/membrane/tests/mocks/worker"

	"github.com/aws/aws-lambda-go/events"
	ep "github.com/nitric-dev/membrane/pkg/plugins/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockLambdaRuntime struct {
	lambda_service.LambdaRuntimeHandler
	// FIXME: Make this a union array of stuff to send....
	eventQueue []interface{}
}

func (m *MockLambdaRuntime) Start(handler interface{}) {
	// cast the function type to what we know it will be
	typedFunc := handler.(func(ctx context.Context, event lambda_service.Event) (interface{}, error))
	for _, event := range m.eventQueue {

		bytes, _ := json.Marshal(event)
		evt := lambda_service.Event{}

		json.Unmarshal(bytes, &evt)
		// Unmarshal the thing into the event type we expect...
		// TODO: Do something with out results here...
		_, err := typedFunc(context.TODO(), evt)

		if err != nil {
			// Print the error?
		}
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
	pool.AddWorker(mockHandler)

	AfterEach(func() {
		mockHandler.Reset()
	})

	Context("Http Events", func() {
		When("Sending a compliant HTTP Event", func() {

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.APIGatewayV2HTTPRequest{
					Headers: map[string]string{
						"User-Agent":            "Test",
						"x-nitric-payload-type": "TestPayload",
						"x-nitric-request-id":   "test-request-id",
						"Content-Type":          "text/plain",
					},
					RawPath: "/test/test",
					Body:    "Test Payload",
					RequestContext: events.APIGatewayV2HTTPRequestContext{
						HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
							Method: "GET",
						},
					},
				}},
			}

			client, _ := lambda_service.NewWithRuntime(runtime.Start)

			// This function will block which means we don't need to wait on processing,
			// the function will unblock once processing has finished, this is due to our mock
			// handler only looping once over each request
			It("The gateway should translate into a standard NitricRequest", func() {
				client.Start(pool)

				By("Handling a single HTTP request")
				Expect(len(mockHandler.ReceivedRequests)).To(Equal(1))

				request := mockHandler.ReceivedRequests[0]

				By("Retaining the body")
				Expect(string(request.Body)).To(BeEquivalentTo("Test Payload"))
				By("Retaining the Headers")
				Expect(request.Header["User-Agent"]).To(Equal("Test"))
				Expect(request.Header["x-nitric-payload-type"]).To(Equal("TestPayload"))
				Expect(request.Header["x-nitric-request-id"]).To(Equal("test-request-id"))
				Expect(request.Header["Content-Type"]).To(Equal("text/plain"))
				By("Retaining the method")
				Expect(request.Method).To(Equal("GET"))
				By("Retaining the path")
				Expect(request.Path).To(Equal("/test/test"))
			})
		})
	})

	Context("SNS Events", func() {
		When("The Lambda Gateway receives SNS events", func() {
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

			messageBytes, _ := json.Marshal(&event)

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

			client, _ := lambda_service.NewWithRuntime(runtime.Start)

			It("The gateway should translate into a standard NitricRequest", func() {
				// This function will block which means we don't need to wait on processing,
				// the function will unblock once processing has finished, this is due to our mock
				// handler only looping once over each request
				client.Start(pool)

				By("Handling a single event")
				Expect(len(mockHandler.ReceivedEvents)).To(Equal(1))

				request := mockHandler.ReceivedEvents[0]

				By("Containing the Source Topic")
				Expect(request.Topic).To(Equal("MyTopic"))
			})
		})
	})
})
