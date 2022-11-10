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

package sqs_service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/mocks/provider"
	mocks_sqs "github.com/nitrictech/nitric/mocks/sqs"
	"github.com/nitrictech/nitric/pkg/plugins/queue"
	"github.com/nitrictech/nitric/pkg/providers/aws/core"
)

var _ = Describe("Sqs", func() {
	Context("getUrlForQueueName", func() {
		When("GetResources returns an error", func() {
			It("Should fail to publish the message", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				providerMock := mock_provider.NewMockAwsProvider(ctrl)
				plugin := NewWithClient(providerMock, sqsMock).(*SQSQueueService)

				By("Calling GetResources and receiving an error")
				providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Times(1).Return(nil, fmt.Errorf("mock-error"))

				_, err := plugin.getUrlForQueueName(context.TODO(), "test-queue")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(Equal("error retrieving queue list"))
				ctrl.Finish()
			})
		})

		When("No queues exist", func() {
			It("Should fail to publish the message", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				providerMock := mock_provider.NewMockAwsProvider(ctrl)
				plugin := NewWithClient(providerMock, sqsMock).(*SQSQueueService)

				By("Calling GetResources and have queue be missing")
				providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Times(1).Return(map[string]string{}, nil)

				_, err := plugin.getUrlForQueueName(context.TODO(), "test-queue")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("queue test-queue does not exist"))
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
				providerMock := mock_provider.NewMockAwsProvider(ctrl)
				plugin := NewWithClient(providerMock, sqsMock)

				queueUrl := aws.String("https://example.com/test-queue")

				By("The queue being available")
				providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Return(map[string]string{
					"test-queue": "arn:aws:sqs:us-east-2:444455556666:test-queue",
				}, nil)

				By("Calling GetQueueUrl to get the queue name")
				sqsMock.EXPECT().GetQueueUrl(gomock.Any(), gomock.Any()).Times(1).Return(&sqs.GetQueueUrlOutput{
					QueueUrl: queueUrl,
				}, nil)

				By("Calling SendMessageBatch with the expected batch entries")
				sqsMock.EXPECT().SendMessageBatch(gomock.Any(), &sqs.SendMessageBatchInput{
					QueueUrl: queueUrl,
					Entries: []types.SendMessageBatchRequestEntry{
						{
							Id:          aws.String("1234"),
							MessageBody: aws.String(`{"id":"1234","payloadType":"test-payload","payload":{"Test":"Test"}}`),
						},
					},
				}).Return(&sqs.SendMessageBatchOutput{}, nil)

				_, err := plugin.SendBatch(context.TODO(), "test-queue", []queue.NitricTask{
					{
						ID:          "1234",
						PayloadType: "test-payload",
						Payload: map[string]interface{}{
							"Test": "Test",
						},
					},
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
				ctrl.Finish()
			})
		})

		When("Publishing to a queue that doesn't exist", func() {
			When("List queues returns an error", func() {
				It("Should fail to publish the message", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					By("provider GetResources returning an error")
					providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Return(nil, fmt.Errorf("mock-error"))

					_, err := plugin.SendBatch(context.TODO(), "test-queue", []queue.NitricTask{
						{
							ID:          "1234",
							PayloadType: "test-payload",
							Payload: map[string]interface{}{
								"Test": "Test",
							},
						},
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("error retrieving queue list"))
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
					providerMock := mock_provider.NewMockAwsProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("calling provider GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Return(map[string]string{
						"mock-queue": "arn:aws:sqs:us-east-2:444455556666:mock-queue",
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
								Body:          aws.String(`{"id":"1234","payloadType":"test-payload","payload":{"Test":"Test"}}`),
							},
						},
					}, nil)

					depth := uint32(10)

					By("Returning the task")
					messages, err := plugin.Receive(context.TODO(), queue.ReceiveOptions{
						QueueName: "mock-queue",
						Depth:     &depth,
					})

					Expect(messages).To(HaveLen(1))
					Expect(messages[0]).To(BeEquivalentTo(queue.NitricTask{
						ID:          "1234",
						PayloadType: "test-payload",
						LeaseID:     "mockreceipthandle",
						Payload: map[string]interface{}{
							"Test": "Test",
						},
					}))
					Expect(err).ShouldNot(HaveOccurred())

					ctrl.Finish()
				})
			})

			When("There are no messages on the queue", func() {
				It("Should receive no messages", func() {
					ctrl := gomock.NewController(GinkgoT())
					sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
					providerMock := mock_provider.NewMockAwsProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Return(map[string]string{
						"mock-queue": "arn:aws:sqs:us-east-2:444455556666:mock-queue",
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

					depth := uint32(10)

					msgs, err := plugin.Receive(context.TODO(), queue.ReceiveOptions{
						QueueName: "mock-queue",
						Depth:     &depth,
					})

					By("Returning an empty array of tasks")
					Expect(msgs).To(HaveLen(0))

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

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
					providerMock := mock_provider.NewMockAwsProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Return(map[string]string{
						"test-queue": "arn:aws:sqs:us-east-2:444455556666:test-queue",
					}, nil)

					By("Calling ListQueueTags to get the x-nitric-name")
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

					err := plugin.Complete(context.TODO(), "test-queue", "lease-id")

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
					providerMock := mock_provider.NewMockAwsProvider(ctrl)
					plugin := NewWithClient(providerMock, sqsMock)

					queueUrl := aws.String("http://example.com/queue")

					By("Calling GetResources to get the queue arn")
					providerMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Queue).Times(1).Return(map[string]string{
						"test-queue": "arn:aws:sqs:us-east-2:444455556666:test-queue",
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

					err := plugin.Complete(context.TODO(), "test-queue", "test-id")

					By("returning the error")
					Expect(err).Should(HaveOccurred())

					ctrl.Finish()
				})
			})
		})
	})
})
