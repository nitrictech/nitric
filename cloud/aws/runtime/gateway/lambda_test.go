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
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/cloud/aws/mocks/provider"
	"github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	lambda_service "github.com/nitrictech/nitric/cloud/aws/runtime/gateway"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	mock_apis "github.com/nitrictech/nitric/core/mocks/workers/apis"
	mock_storage "github.com/nitrictech/nitric/core/mocks/workers/storage"
	mock_topics "github.com/nitrictech/nitric/core/mocks/workers/topics"
	mock_websockets "github.com/nitrictech/nitric/core/mocks/workers/websockets"
	"github.com/nitrictech/nitric/core/pkg/env"
	coreGateway "github.com/nitrictech/nitric/core/pkg/gateway"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	storagepb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
	ep "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	websocketspb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
)

type MockLambdaRuntime struct {
	lambda_service.LambdaRuntimeHandler
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
		_, err = typedFunc(context.TODO(), evt)

		if err != nil {
			Expect(err).Should(HaveOccurred())
		} else {
			Expect(err).To(BeNil())
		}
	}
}

type protoMatcher struct {
	expected proto.Message
}

func (m *protoMatcher) Matches(x interface{}) bool {
	expectedBytes, _ := proto.Marshal(m.expected)
	actualBytes, _ := proto.Marshal(x.(proto.Message))

	return bytes.Equal(expectedBytes, actualBytes)
}

func (m *protoMatcher) String() string {
	return fmt.Sprintf("equivalent to %+v", m.expected)
}

func EqProto(expected proto.Message) gomock.Matcher {
	return &protoMatcher{expected: expected}
}

var _ = Describe("Lambda", func() {
	commonenv.NITRIC_STACK_ID = env.GetEnv("NITRIC_STACK_ID", "test-stack-id")

	Context("Http Events", func() {
		When("Sending a compliant HTTP Event", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockManager := mock_apis.NewMockApiRequestHandler(ctrl)
			// mockHandler := mock_worker.NewMockWorker(ctrl)
			mockProvider := mock_provider.NewMockAwsResourceProvider(ctrl)

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
						APIID: "test-api",
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
				By("The api gateway existing")
				mockProvider.EXPECT().GetApiGatewayById(gomock.Any(), "test-api").Return(&apigatewayv2.GetApiOutput{
					Tags: map[string]string{
						"x-nitric-test-stack-name": "test-api",
					},
				}, nil)

				By("Having at least one worker available")
				mockManager.EXPECT().WorkerCount().Return(1)
				// mockHandler.EXPECT().HandlesTrigger(gomock.Any()).Return(true)

				By("Handling a single HTTP request")
				mockManager.EXPECT().HandleRequest(gomock.Any(), EqProto(&apispb.ServerMessage{
					Content: &apispb.ServerMessage_HttpRequest{
						HttpRequest: &apispb.HttpRequest{
							Method: "GET",
							Path:   "/test/test",
							Headers: map[string]*apispb.HeaderValue{
								"User-Agent":   {Value: []string{"Test"}},
								"Content-Type": {Value: []string{"text/plain"}},
								"Cookie":       {Value: []string{"test1=testcookie1", "test2=testcookie2"}},
							},
							QueryParams: map[string]*apispb.QueryValue{
								"key":  {Value: []string{"test", "test2"}},
								"key2": {Value: []string{"test1"}},
							},
							Body: []byte("Test Payload"),
						},
					},
				})).Return(&apispb.ClientMessage{
					Content: &apispb.ClientMessage_HttpResponse{
						HttpResponse: &apispb.HttpResponse{
							Body: []byte("success"),
						},
					},
				}, nil)

				err := client.Start(&coreGateway.GatewayStartOpts{
					ApiPlugin: mockManager,
				})
				Expect(err).To(BeNil())
			})
		})
	})

	Context("Websocket Events", func() {
		When("Sending a compliant Websocket Event", func() {
			_ = os.Setenv("NITRIC_STACK_ID", "test-stack")
			ctrl := gomock.NewController(GinkgoT())
			mockManager := mock_websockets.NewMockWebsocketRequestHandler(ctrl)

			mockProvider := mock_provider.NewMockAwsResourceProvider(ctrl)

			runtime := MockLambdaRuntime{
				// Setup mock events for our runtime to process...
				eventQueue: []interface{}{&events.APIGatewayWebsocketProxyRequest{
					Headers: map[string]string{
						"User-Agent":   "Test",
						"Content-Type": "text/plain",
					},
					Body: "Test Payload",
					RequestContext: events.APIGatewayWebsocketProxyRequestContext{
						APIID: "test-api",
						// as a connection request
						RouteKey:     "$connect",
						ConnectionID: "testing",
					},
				}},
			}

			client, err := lambda_service.NewWithRuntime(mockProvider, runtime.Start)
			Expect(err).To(BeNil())

			// This function will block which means we don't need to wait on processing,
			// the function will unblock once processing has finished, this is due to our mock
			// handler only looping once over each request
			It("The gateway should translate into a standard NitricRequest", func() {
				// By("Returning the worker")
				// pool.EXPECT().GetWorker(gomock.Any()).Return(mockHandler, nil)

				By("Having at least one worker")
				mockManager.EXPECT().WorkerCount().Return(1)

				By("The websocket gateway existing")
				mockProvider.EXPECT().GetApiGatewayById(gomock.Any(), "test-api").Return(&apigatewayv2.GetApiOutput{
					Tags: map[string]string{
						"x-nitric-test-stack-name": "test-api",
					},
				}, nil)

				By("Handling a single HTTP request")
				mockManager.EXPECT().HandleRequest(&websocketspb.ServerMessage{
					Content: &websocketspb.ServerMessage_WebsocketEventRequest{
						WebsocketEventRequest: &websocketspb.WebsocketEventRequest{
							SocketName:   "test-api",
							ConnectionId: "testing",
							WebsocketEvent: &websocketspb.WebsocketEventRequest_Connection{
								Connection: &websocketspb.WebsocketConnectionEvent{
									QueryParams: map[string]*websocketspb.QueryValue{},
								},
							},
						},
					},
				}).Return(&websocketspb.ClientMessage{
					Id: "TODO",
					Content: &websocketspb.ClientMessage_WebsocketEventResponse{
						WebsocketEventResponse: &websocketspb.WebsocketEventResponse{
							WebsocketResponse: &websocketspb.WebsocketEventResponse_ConnectionResponse{
								ConnectionResponse: &websocketspb.WebsocketConnectionResponse{
									Reject: false,
								},
							},
						},
					},
				}, nil)

				err := client.Start(&coreGateway.GatewayStartOpts{
					WebsocketListenerPlugin: mockManager,
				})
				Expect(err).To(BeNil())
			})
		})
	})

	Context("SNS Events", func() {
		When("The Lambda Gateway receives SNS events", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsResourceProvider(ctrl)

			// pool := mock_pool.NewMockWorkerPool(ctrl)
			mockManager := mock_topics.NewMockSubscriptionRequestHandler(ctrl)
			// mockHandler := mock_worker.NewMockWorker(ctrl)

			topicName := "MyTopic"
			content, _ := structpb.NewStruct(map[string]interface{}{
				"test": "test",
			})

			message := ep.TopicMessage{
				Content: &ep.TopicMessage_StructPayload{
					StructPayload: content,
				},
			}

			messageBytes, err := proto.Marshal(&message)
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
				mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Topic).Return(map[string]resource.ResolvedResource{
					"MyTopic": {
						ARN: "arn:aws:sns:us-east-1:12345678910:arn:MyTopic",
					},
				}, nil)

				By("having at least one worker")
				mockManager.EXPECT().WorkerCount().Return(1)

				By("Handling a single event")
				mockManager.EXPECT().HandleRequest(EqProto(&topicspb.ServerMessage{
					Content: &topicspb.ServerMessage_MessageRequest{
						MessageRequest: &topicspb.MessageRequest{
							TopicName: "MyTopic",
							Message:   &message,
						},
					},
				})).Return(&topicspb.ClientMessage{
					Content: &topicspb.ClientMessage_MessageResponse{
						MessageResponse: &topicspb.MessageResponse{
							Success: true,
						},
					},
				}, nil)

				// This function will block which means we don't need to wait on processing,
				// the function will unblock once processing has finished, this is due to our mock
				// handler only looping once over each request
				err := client.Start(&coreGateway.GatewayStartOpts{
					TopicsListenerPlugin: mockManager,
				})
				Expect(err).To(BeNil())
			})
		})
	})

	Context("S3 Events", func() {
		When("The Lambda Gateway receives S3 Put events", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsResourceProvider(ctrl)

			// pool := mock_pool.NewMockWorkerPool(ctrl)
			mockManager := mock_storage.NewMockBucketRequestHandler(ctrl)
			// mockHandler := mock_worker.NewMockWorker(ctrl)

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
				mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
					"images": {ARN: "arn:aws:sns:us-east-1:12345678910:arn:images"},
				}, nil)
				// pool.EXPECT().GetWorker(gomock.Any()).Return(mockHandler, nil)

				By("Having at least one worker")
				mockManager.EXPECT().WorkerCount().Return(1)
				// mockHandler.EXPECT().HandlesTrigger(gomock.Any()).Return(true)

				By("Handling a single Notification request")
				mockManager.EXPECT().HandleRequest(&storagepb.ServerMessage{
					Content: &storagepb.ServerMessage_BlobEventRequest{
						BlobEventRequest: &storagepb.BlobEventRequest{
							BucketName: "images",
							Event: &storagepb.BlobEventRequest_BlobEvent{
								BlobEvent: &storagepb.BlobEvent{
									Key:  "cat.png",
									Type: storagepb.BlobEventType_Created,
								},
							},
						},
					},
				}).Return(&storagepb.ClientMessage{
					Content: &storagepb.ClientMessage_BlobEventResponse{
						BlobEventResponse: &storagepb.BlobEventResponse{
							Success: true,
						},
					},
				}, nil)

				err := client.Start(&coreGateway.GatewayStartOpts{
					StorageListenerPlugin: mockManager,
				})
				Expect(err).To(BeNil())
			})
		})
		When("The Lambda Gateway receives S3 Delete events", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsResourceProvider(ctrl)

			mockManager := mock_storage.NewMockBucketRequestHandler(ctrl)

			// pool := mock_pool.NewMockWorkerPool(ctrl)
			// mockHandler := mock_worker.NewMockWorker(ctrl)

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
				By("The bucket existing")
				mockProvider.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Bucket).Return(map[string]resource.ResolvedResource{
					"images": {ARN: "arn:aws:sns:us-east-1:12345678910:arn:images"},
				}, nil)

				By("Having at least one worker")
				mockManager.EXPECT().WorkerCount().Return(1)

				By("Handling a single Notification request")
				mockManager.EXPECT().HandleRequest(&storagepb.ServerMessage{
					Content: &storagepb.ServerMessage_BlobEventRequest{
						BlobEventRequest: &storagepb.BlobEventRequest{
							BucketName: "images",
							Event: &storagepb.BlobEventRequest_BlobEvent{
								BlobEvent: &storagepb.BlobEvent{
									Key:  "cat.png",
									Type: storagepb.BlobEventType_Deleted,
								},
							},
						},
					},
				}).Return(&storagepb.ClientMessage{
					Content: &storagepb.ClientMessage_BlobEventResponse{
						BlobEventResponse: &storagepb.BlobEventResponse{
							Success: true,
						},
					},
				}, nil)

				err := client.Start(&coreGateway.GatewayStartOpts{
					StorageListenerPlugin: mockManager,
				})
				Expect(err).To(BeNil())
			})
		})
	})
})
