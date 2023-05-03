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
	"github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	lambda_service "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	mock_pool "github.com/nitrictech/nitric/core/mocks/pool"
	mock_worker "github.com/nitrictech/nitric/core/mocks/worker"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	ep "github.com/nitrictech/nitric/core/pkg/plugins/events"
)

type MockLambdaRuntime struct {
	lambda_service.LambdaRuntimeHandler
	// FIXME: Make this a union array of stuff to send....
	eventQueue []interface{}
}

func (m *MockLambdaRuntime) Start(handler interface{}) {
	// cast the function type to what we know it will be
	typedFunc := handler.(func(ctx context.Context, data gateway.Event) (interface{}, error))
	for _, event := range m.eventQueue {
		bytes, _ := json.Marshal(event)
		evt := gateway.Event{}

		err := json.Unmarshal(bytes, &evt)
		Expect(err).To(BeNil())

		// Unmarshal the thing into the event type we expect...
		// TODO: Do something with out results here...
		_, err = typedFunc(context.TODO(), evt)
		Expect(err).To(BeNil())
	}
}

var _ = Describe("Lambda", func() {
	Context("Http Events", func() {
		When("Sending a compliant HTTP Event", func() {
			ctrl := gomock.NewController(GinkgoT())
			pool := mock_pool.NewMockWorkerPool(ctrl)

			mockHandler := mock_worker.NewMockWorker(ctrl)
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
					RouteKey:       "non-null",
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
				By("Returning the worker")
				pool.EXPECT().GetWorker(gomock.Any()).Return(mockHandler, nil)

				By("Handling all request types")
				mockHandler.EXPECT().HandlesTrigger(gomock.Any()).Return(true)

				By("Handling a single HTTP request")
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), &v1.TriggerRequest{
					Data: []byte("Test Payload"),
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "GET",
							Path:   "/test/test",
							Headers: map[string]*v1.HeaderValue{
								"User-Agent":            {Value: []string{"Test"}},
								"x-nitric-payload-type": {Value: []string{"TestPayload"}},
								"x-nitric-request-id":   {Value: []string{"test-request-id"}},
								"Content-Type":          {Value: []string{"text/plain"}},
								"Cookie":                {Value: []string{"test1=testcookie1", "test2=testcookie2"}},
							},
							QueryParams: map[string]*v1.QueryValue{
								"key":  {Value: []string{"test", "test2"}},
								"key2": {Value: []string{"test1"}},
							},
						},
					},
				}).Return(&v1.TriggerResponse{
					Data: []byte("success"),
					Context: &v1.TriggerResponse_Http{
						Http: &v1.HttpResponseContext{},
					},
				}, nil)

				err := client.Start(pool)
				Expect(err).To(BeNil())
			})
		})
	})

	Context("SNS Events", func() {
		When("The Lambda Gateway receives SNS events", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsProvider(ctrl)

			pool := mock_pool.NewMockWorkerPool(ctrl)
			mockHandler := mock_worker.NewMockWorker(ctrl)

			topicName := "MyTopic"
			eventPayload := map[string]interface{}{
				"test": "test",
			}

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
								TopicArn: fmt.Sprintf("arn:aws:sns:us-east-1:12345678910:arn:%s", topicName),
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
					"MyTopic": "arn:aws:sns:us-east-1:12345678910:arn:MyTopic",
				}, nil)

				By("Returning the worker")
				pool.EXPECT().GetWorker(gomock.Any()).Return(mockHandler, nil)

				By("Handling all request types")
				mockHandler.EXPECT().HandlesTrigger(gomock.Any()).Return(true)

				By("Handling a single event")
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), &v1.TriggerRequest{
					Data: messageBytes,
					Context: &v1.TriggerRequest_Topic{
						Topic: &v1.TopicTriggerContext{
							Topic: "MyTopic",
						},
					},
				})
				// This function will block which means we don't need to wait on processing,
				// the function will unblock once processing has finished, this is due to our mock
				// handler only looping once over each request
				err := client.Start(pool)
				Expect(err).To(BeNil())
			})
		})
	})

	Context("S3 Events", func() {
		When("The Lambda Gateway receives S3 Put events", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsProvider(ctrl)

			pool := mock_pool.NewMockWorkerPool(ctrl)
			mockHandler := mock_worker.NewMockWorker(ctrl)

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.S3Event{
					Records: []events.S3EventRecord{
						{
							EventVersion: "",
							EventSource:  "aws:s3",
							EventName:    "ObjectCreated:Put",
							S3: events.S3Entity{
								Bucket: events.S3Bucket{
									Name: "images",
									Arn:  "arn:aws:sns:us-east-1:12345678910:arn:images",
								},
								Object: events.S3Object{
									Key: "cat.png",
								},
							},
							ResponseElements: map[string]string{},
						},
					},
				}},
			}

			client, err := lambda_service.NewWithRuntime(mockProvider, runtime.Start)
			Expect(err).To(BeNil())

			// This function will block which means we don't need to wait on processing,
			// the function will unblock once processing has finished, this is due to our mock
			// handler only looping once over each request
			It("The gateway should translate into a standard NitricRequest", func() {
				By("Returning the worker")
				mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{
					"images": "arn:aws:sns:us-east-1:12345678910:arn:images",
				}, nil)
				pool.EXPECT().GetWorker(gomock.Any()).Return(mockHandler, nil)

				By("Handling all request types")
				mockHandler.EXPECT().HandlesTrigger(gomock.Any()).Return(true)

				By("Handling a single Notification request")
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), &v1.TriggerRequest{
					Context: &v1.TriggerRequest_Notification{
						Notification: &v1.NotificationTriggerContext{
							Source: "images",
							Notification: &v1.NotificationTriggerContext_Bucket{
								Bucket: &v1.BucketNotification{
									Key:  "cat.png",
									Type: v1.BucketNotificationType_Created,
								},
							},
						},
					},
				}).Return(&v1.TriggerResponse{
					Data: []byte("success"),
					Context: &v1.TriggerResponse_Notification{
						Notification: &v1.NotificationResponseContext{
							Success: true,
						},
					},
				}, nil)

				err := client.Start(pool)
				Expect(err).To(BeNil())
			})
		})
		When("The Lambda Gateway receives S3 Delete events", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsProvider(ctrl)

			pool := mock_pool.NewMockWorkerPool(ctrl)
			mockHandler := mock_worker.NewMockWorker(ctrl)

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.S3Event{
					Records: []events.S3EventRecord{
						{
							EventVersion: "",
							EventSource:  "aws:s3",
							EventName:    "ObjectRemoved:Delete",
							S3: events.S3Entity{
								Bucket: events.S3Bucket{
									Name: "images",
									Arn:  "arn:aws:sns:us-east-1:12345678910:arn:images",
								},
								Object: events.S3Object{
									Key: "cat.png",
								},
							},
							ResponseElements: map[string]string{},
						},
					},
				}},
			}

			client, err := lambda_service.NewWithRuntime(mockProvider, runtime.Start)
			Expect(err).To(BeNil())

			// This function will block which means we don't need to wait on processing,
			// the function will unblock once processing has finished, this is due to our mock
			// handler only looping once over each request
			It("The gateway should translate into a standard NitricRequest", func() {
				By("Returning the worker")
				mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{
					"images": "arn:aws:sns:us-east-1:12345678910:arn:images",
				}, nil)
				pool.EXPECT().GetWorker(gomock.Any()).Return(mockHandler, nil)

				By("Handling all request types")
				mockHandler.EXPECT().HandlesTrigger(gomock.Any()).Return(true)

				By("Handling a single Notification request")
				mockHandler.EXPECT().HandleTrigger(gomock.Any(), &v1.TriggerRequest{
					Context: &v1.TriggerRequest_Notification{
						Notification: &v1.NotificationTriggerContext{
							Source: "images",
							Notification: &v1.NotificationTriggerContext_Bucket{
								Bucket: &v1.BucketNotification{
									Key:  "cat.png",
									Type: v1.BucketNotificationType_Deleted,
								},
							},
						},
					},
				}).Return(&v1.TriggerResponse{
					Data: []byte("success"),
					Context: &v1.TriggerResponse_Notification{
						Notification: &v1.NotificationResponseContext{
							Success: true,
						},
					},
				}, nil)

				err := client.Start(pool)
				Expect(err).To(BeNil())
			})
		})
	})
})
