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

package queue

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/cloud/aws/mocks/provider"
	mocks_sqs "github.com/nitrictech/nitric/cloud/aws/mocks/sqs"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	queuepb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
)

var _ = Describe("Sqs", func() {
	testStruct, err := structpb.NewStruct(map[string]interface{}{"Test": "Test"})

	Expect(err).To(BeNil())

	testPayloadBytes, err := proto.Marshal(testStruct)

	Expect(err).To(BeNil())

	testPayloadB64 := base64.StdEncoding.EncodeToString(testPayloadBytes)

	Context("getUrlForQueueName", func() {
		When("GetResources returns an error", func() {
			It("Should fail to publish the message", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
				plugin := NewWithClient(providerMock, sqsMock).(*SQSQueueService)

				By("Calling GetResources and receiving an error")
				providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Times(1).Return(nil, fmt.Errorf("mock-error"))

				_, err := plugin.getUrlForQueueName(context.TODO(), "test-queue")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("error retrieving queue list: mock-error"))
				ctrl.Finish()
			})
		})

		When("No queues exist", func() {
			It("Should fail to publish the message", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
				plugin := NewWithClient(providerMock, sqsMock).(*SQSQueueService)

				By("Calling GetResources and have queue be missing")
				providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Times(1).Return(map[string]resource.ResolvedResource{}, nil)

				_, err := plugin.getUrlForQueueName(context.TODO(), "test-queue")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("arn for queue test-queue could not be determined"))
				ctrl.Finish()
			})
		})
	})

	// Tests for the BatchPush method
	Context("BatchSend", func() {
		When("Sending to a queue that exists", func() {
			It("Should send the task to the queue", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
				plugin := NewWithClient(providerMock, sqsMock)

				queueUrl := aws.String("https://example.com/test-queue")

				By("The queue being available")
				providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Return(map[string]resource.ResolvedResource{
					"test-queue": {
						ARN: "arn:aws:sqs:us-east-2:444455556666:test-queue",
					},
				}, nil)

				By("Calling GetQueueUrl to get the queue name")
				sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Times(1).Return(&sqs.GetQueueUrlOutput{
					QueueUrl: queueUrl,
				}, nil)

				By("Calling SendMessageBatch with the expected batch entries")
				sqsMock.EXPECT().SendMessageBatch(gomock.Any(), gomock.Any()).Return(&sqs.SendMessageBatchOutput{}, nil)

				_, err := plugin.Send(context.TODO(), &queuepb.QueueSendRequestBatch{
					QueueName: "test-queue",
					Requests: []*queuepb.QueueSendRequest{
						{
							Payload: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"Test": structpb.NewStringValue("Test"),
								},
							},
						},
					},
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
				ctrl.Finish()
			})
		})

		When("Permission to access the queue is missing", func() {
			It("Should return an error", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
				plugin := NewWithClient(providerMock, sqsMock)

				queueUrl := aws.String("https://example.com/test-queue")

				By("The queue being available")
				providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Return(map[string]resource.ResolvedResource{
					"test-queue": {
						ARN: "arn:aws:sqs:us-east-2:444455556666:test-queue",
					},
				}, nil)

				By("Calling GetQueueUrl to get the queue name")
				sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Times(1).Return(&sqs.GetQueueUrlOutput{
					QueueUrl: queueUrl,
				}, nil)

				opErr := &smithy.OperationError{
					ServiceID: "SQS",
					Err:       errors.New("AccessDenied"),
				}

				By("Calling SendMessageBatch with the expected batch entries")
				sqsMock.EXPECT().SendMessageBatch(gomock.Any(), gomock.Any()).Return(nil, opErr)

				_, err := plugin.Send(context.TODO(), &queuepb.QueueSendRequestBatch{
					QueueName: "test-queue",
					Requests: []*queuepb.QueueSendRequest{
						{
							Payload: &structpb.Struct{
								Fields: map[string]*structpb.Value{
									"Test": structpb.NewStringValue("Test"),
								},
							},
						},
					},
				})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unable to send task"))
				ctrl.Finish()
			})
		})

		When("Publishing to a queue that doesn't exist", func() {
			When("List queues returns an error", func() {
				It("Should fail to publish the message", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					By("provider GetResources returning an error")
					providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Return(nil, fmt.Errorf("mock-error"))

					_, err := plugin.Send(context.TODO(), &queuepb.QueueSendRequestBatch{
						QueueName: "test-queue",
						Requests: []*queuepb.QueueSendRequest{
							{
								Payload: &structpb.Struct{
									Fields: map[string]*structpb.Value{
										"Test": structpb.NewStringValue("Test"),
									},
								},
							},
						},
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("rpc error: code = NotFound desc = SQSQueueService.SendBatch unable to find queue"))
					ctrl.Finish()
				})
			})
		})
	})

	// Tests for the Receive method
	Context("Receive", func() {
		When("Receive from a queue that exists", func() {
			When("There is a message on the queue", func() {
				It("Should receive the message", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("calling provider GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Return(map[string]resource.ResolvedResource{
						"mock-queue": {
							ARN: "arn:aws:sqs:us-east-2:444455556666:mock-queue",
						},
					}, nil)

					By("calling provider GetQueuUrl")
					sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Return(&sqs.GetQueueUrlOutput{
						QueueUrl: queueUrl,
					}, nil)

					By("Calling ReceiveMessage with the expected inputs")
					sqsMock.EXPECT().ReceiveMessage(gomock.Any(), &sqs.ReceiveMessageInput{
						MaxNumberOfMessages: int32(10),
						MessageAttributeNames: []string{
							string(types.QueueAttributeNameAll),
						},
						QueueUrl: queueUrl,
					}).Times(1).Return(&sqs.ReceiveMessageOutput{
						Messages: []types.Message{
							{
								ReceiptHandle: aws.String("mockreceipthandle"),
								Body:          aws.String(testPayloadB64),
							},
						},
					}, nil)

					By("Returning the task")
					response, err := plugin.Receive(context.TODO(), &queuepb.QueueReceiveRequest{
						QueueName: "mock-queue",
						Depth:     10,
					})

					Expect(err).ShouldNot(HaveOccurred())
					Expect(response.Tasks).To(HaveLen(1))
					Expect(response.Tasks[0].LeaseId).To(BeEquivalentTo("mockreceipthandle"))
					Expect(response.Tasks[0].Payload.AsMap()).To(BeEquivalentTo(testStruct.AsMap()))

					ctrl.Finish()
				})
			})

			When("There are no messages on the queue", func() {
				It("Should receive no messages", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Return(map[string]resource.ResolvedResource{
						"mock-queue": {
							ARN: "arn:aws:sqs:us-east-2:444455556666:mock-queue",
						},
					}, nil)

					By("Calling GetQueueUrl to get the queue url")
					sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Times(1).Return(&sqs.GetQueueUrlOutput{
						QueueUrl: queueUrl,
					}, nil)

					By("Calling ReceiveMessage with the expected inputs")
					sqsMock.EXPECT().ReceiveMessage(gomock.Any(), &sqs.ReceiveMessageInput{
						MaxNumberOfMessages: int32(10),
						MessageAttributeNames: []string{
							string(types.QueueAttributeNameAll),
						},
						QueueUrl: queueUrl,
					}).Times(1).Return(&sqs.ReceiveMessageOutput{
						Messages: []types.Message{},
					}, nil)

					response, err := plugin.Receive(context.TODO(), &queuepb.QueueReceiveRequest{
						QueueName: "mock-queue",
						Depth:     10,
					})

					By("Returning an empty array of tasks")
					Expect(response.Tasks).To(HaveLen(0))

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					ctrl.Finish()
				})
			})

			When("Permission to access the queue is missing", func() {
				It("Should return an error", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Return(map[string]resource.ResolvedResource{
						"mock-queue": {
							ARN: "arn:aws:sqs:us-east-2:444455556666:mock-queue",
						},
					}, nil)

					opErr := &smithy.OperationError{
						ServiceID: "SQS",
						Err:       errors.New("AccessDenied"),
					}

					By("Calling GetQueueUrl to get the queue url")
					sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Times(1).Return(&sqs.GetQueueUrlOutput{
						QueueUrl: queueUrl,
					}, nil)

					By("Calling ReceiveMessage with the expected inputs")
					sqsMock.EXPECT().ReceiveMessage(gomock.Any(), &sqs.ReceiveMessageInput{
						MaxNumberOfMessages: int32(10),
						MessageAttributeNames: []string{
							string(types.QueueAttributeNameAll),
						},
						QueueUrl: queueUrl,
					}).Times(1).Return(nil, opErr)

					_, err := plugin.Receive(context.TODO(), &queuepb.QueueReceiveRequest{
						QueueName: "mock-queue",
						Depth:     10,
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unable to receive task"))

					ctrl.Finish()
				})
			})
		})

		// Tests for the Complete method
		Context("Complete", func() {
			When("The message is successfully deleted from SQS", func() {
				// No errors set on mock, 'complete' won't return an error.
				It("Should successfully delete the task", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Return(map[string]resource.ResolvedResource{
						"test-queue": {
							ARN: "arn:aws:sqs:us-east-2:444455556666:test-queue",
						},
					}, nil)

					By("Calling ListQueueTags to get the stack specific nitric name")
					sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Times(1).Return(&sqs.GetQueueUrlOutput{
						QueueUrl: queueUrl,
					}, nil)

					By("Calling SQS with the queue url and task lease id")
					sqsMock.EXPECT().DeleteMessage(gomock.Any(), &sqs.DeleteMessageInput{
						QueueUrl:      queueUrl,
						ReceiptHandle: aws.String("lease-id"),
					}).Times(1).Return(
						&sqs.DeleteMessageOutput{},
						nil,
					)

					_, err := plugin.Complete(context.TODO(), &queuepb.QueueCompleteRequest{
						QueueName: "test-queue",
						LeaseId:   "lease-id",
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					ctrl.Finish()
				})
			})

			When("The message fails to delete from SQS", func() {
				// No errors set on mock, 'complete' won't return an error.
				It("Return an error", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsResourceProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("http://example.com/queue")

					By("Calling GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Queue).Times(1).Return(map[string]resource.ResolvedResource{
						"test-queue": {
							ARN: "arn:aws:sqs:us-east-2:444455556666:test-queue",
						},
					}, nil)

					By("Calling GetQueueUrl to get the queueurl")
					sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Times(1).Return(&sqs.GetQueueUrlOutput{
						QueueUrl: queueUrl,
					}, nil)

					By("Calling SQS with the queue url and task lease id")
					sqsMock.EXPECT().DeleteMessage(gomock.Any(), &sqs.DeleteMessageInput{
						QueueUrl:      queueUrl,
						ReceiptHandle: aws.String("test-id"),
					}).Return(nil, fmt.Errorf("mock-error"))

					_, err := plugin.Complete(context.TODO(), &queuepb.QueueCompleteRequest{
						QueueName: "test-queue",
						LeaseId:   "test-id",
					})

					By("returning the error")
					Expect(err).Should(HaveOccurred())

					ctrl.Finish()
				})
			})
		})
	})
})
