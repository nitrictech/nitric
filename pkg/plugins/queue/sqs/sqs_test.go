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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/golang/mock/gomock"
	mocks_sqs "github.com/nitrictech/nitric/mocks/sqs"

	"github.com/nitrictech/nitric/pkg/plugins/queue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sqs", func() {
	Context("getUrlForQueueName", func() {
		When("List queues returns an error", func() {
			It("Should fail to publish the message", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				plugin := NewWithClient(sqsMock).(*SQSQueueService)

				By("Calling ListQueues and receiving an error")
				sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(nil, fmt.Errorf("mock-error"))

				_, err := plugin.getUrlForQueueName("test-queue")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Encountered an error retrieving the queue list: mock-error"))
				ctrl.Finish()
			})
		})

		When("No queues exist", func() {
			It("Should fail to publish the message", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				plugin := NewWithClient(sqsMock).(*SQSQueueService)

				By("Calling ListQueues and receiving no queue")
				sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(&sqs.ListQueuesOutput{
					QueueUrls: []*string{},
				}, nil)

				_, err := plugin.getUrlForQueueName("test-queue")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find queue with name: test-queue"))
				ctrl.Finish()
			})
		})

		When("No queue tags match the nitric name", func() {
			It("Should fail to publish the message", func() {
				ctrl := gomock.NewController(GinkgoT())
				sqsMock := mocks_sqs.NewMockSQSAPI(ctrl)
				plugin := NewWithClient(sqsMock).(*SQSQueueService)

				// Name is in the URL, but that's not important.
				queueUrl := aws.String("https://example.com/test-queue")

				By("Calling ListQueues and receiving no queue")
				sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(&sqs.ListQueuesOutput{
					QueueUrls: []*string{queueUrl},
				}, nil)

				By("Calling ListQueueTags with the available queues")
				sqsMock.EXPECT().ListQueueTags(&sqs.ListQueueTagsInput{QueueUrl: queueUrl}).Times(1).Return(&sqs.ListQueueTagsOutput{
					Tags: map[string]*string{
						// The nitric name tag doesn't match the expected queue name
						"x-nitric-name": aws.String("not-test-queue"),
					},
				}, nil)

				By("calling getUrlForQueueName with test-queue")
				_, err := plugin.getUrlForQueueName("test-queue")

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unable to find queue with name: test-queue"))
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
				plugin := NewWithClient(sqsMock)

				queueUrl := aws.String("https://example.com/test-queue")

				By("Calling ListQueues to get the queue name")
				sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(&sqs.ListQueuesOutput{
					QueueUrls: []*string{queueUrl},
				}, nil)

				By("Calling ListQueueTags to get the x-nitric-name")
				sqsMock.EXPECT().ListQueueTags(gomock.Any()).Times(1).Return(&sqs.ListQueueTagsOutput{
					Tags: map[string]*string{
						"x-nitric-name": aws.String("test-queue"),
					},
				}, nil)

				By("Calling SendMessageBatch with the expected batch entries")
				sqsMock.EXPECT().SendMessageBatch(&sqs.SendMessageBatchInput{
					QueueUrl: queueUrl,
					Entries: []*sqs.SendMessageBatchRequestEntry{
						{
							Id:          aws.String("1234"),
							MessageBody: aws.String(`{"id":"1234","payloadType":"test-payload","payload":{"Test":"Test"}}`),
						},
					},
				}).Return(&sqs.SendMessageBatchOutput{}, nil)

				_, err := plugin.SendBatch("test-queue", []queue.NitricTask{
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
					plugin := NewWithClient(sqsMock)

					By("Calling ListQueues and receiving an error")
					sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(nil, fmt.Errorf("mock-error"))

					_, err := plugin.SendBatch("test-queue", []queue.NitricTask{
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
					Expect(err.Error()).To(ContainSubstring("Encountered an error retrieving the queue list: mock-error"))
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
					plugin := NewWithClient(sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling ListQueues to get the queue name")
					sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(&sqs.ListQueuesOutput{
						QueueUrls: []*string{queueUrl},
					}, nil)

					By("Calling ListQueueTags to get the x-nitric-name")
					sqsMock.EXPECT().ListQueueTags(gomock.Any()).Times(1).Return(&sqs.ListQueueTagsOutput{
						Tags: map[string]*string{
							"x-nitric-name": aws.String("mock-queue"),
						},
					}, nil)

					By("Calling ReceiveMessage with the expected inputs")
					sqsMock.EXPECT().ReceiveMessage(&sqs.ReceiveMessageInput{
						MaxNumberOfMessages: aws.Int64(int64(10)),
						MessageAttributeNames: []*string{
							aws.String(sqs.QueueAttributeNameAll),
						},
						QueueUrl: queueUrl,
					}).Times(1).Return(&sqs.ReceiveMessageOutput{
						Messages: []*sqs.Message{
							{
								ReceiptHandle: aws.String("mockreceipthandle"),
								Body:          aws.String(`{"id":"1234","payloadType":"test-payload","payload":{"Test":"Test"}}`),
							},
						},
					}, nil)

					depth := uint32(10)

					By("Returning the task")
					messages, err := plugin.Receive(queue.ReceiveOptions{
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
					plugin := NewWithClient(sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling ListQueues to get the queue name")
					sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(&sqs.ListQueuesOutput{
						QueueUrls: []*string{queueUrl},
					}, nil)

					By("Calling ListQueueTags to get the x-nitric-name")
					sqsMock.EXPECT().ListQueueTags(gomock.Any()).Times(1).Return(&sqs.ListQueueTagsOutput{
						Tags: map[string]*string{
							"x-nitric-name": aws.String("mock-queue"),
						},
					}, nil)

					By("Calling ReceiveMessage with the expected inputs")
					sqsMock.EXPECT().ReceiveMessage(&sqs.ReceiveMessageInput{
						MaxNumberOfMessages: aws.Int64(int64(10)),
						MessageAttributeNames: []*string{
							aws.String(sqs.QueueAttributeNameAll),
						},
						QueueUrl: queueUrl,
					}).Times(1).Return(&sqs.ReceiveMessageOutput{
						Messages: []*sqs.Message{},
					}, nil)

					depth := uint32(10)

					msgs, err := plugin.Receive(queue.ReceiveOptions{
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
					plugin := NewWithClient(sqsMock)

					queueUrl := aws.String("https://example.com/test-queue")

					By("Calling ListQueues to get the queue name")
					sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(&sqs.ListQueuesOutput{
						QueueUrls: []*string{queueUrl},
					}, nil)

					By("Calling ListQueueTags to get the x-nitric-name")
					sqsMock.EXPECT().ListQueueTags(gomock.Any()).Times(1).Return(&sqs.ListQueueTagsOutput{
						Tags: map[string]*string{
							"x-nitric-name": aws.String("test-queue"),
						},
					}, nil)

					By("Calling SQS with the queue url and task lease id")
					sqsMock.EXPECT().DeleteMessage(&sqs.DeleteMessageInput{
						QueueUrl:      queueUrl,
						ReceiptHandle: aws.String("lease-id"),
					}).Times(1).Return(
						&sqs.DeleteMessageOutput{},
						nil,
					)

					err := plugin.Complete("test-queue", "lease-id")

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
					plugin := NewWithClient(sqsMock)

					queueUrl := aws.String("http://example.com/queue")

					By("Calling ListQueues to get the queue name")
					sqsMock.EXPECT().ListQueues(&sqs.ListQueuesInput{}).Times(1).Return(&sqs.ListQueuesOutput{
						QueueUrls: []*string{queueUrl},
					}, nil)

					By("Calling ListQueueTags to get the x-nitric-name")
					sqsMock.EXPECT().ListQueueTags(gomock.Any()).Times(1).Return(&sqs.ListQueueTagsOutput{
						Tags: map[string]*string{
							"x-nitric-name": aws.String("test-queue"),
						},
					}, nil)

					By("Calling SQS with the queue url and task lease id")
					sqsMock.EXPECT().DeleteMessage(&sqs.DeleteMessageInput{
						QueueUrl:      queueUrl,
						ReceiptHandle: aws.String("test-id"),
					}).Return(nil, fmt.Errorf("mock-error"))

					err := plugin.Complete("test-queue", "test-id")

					By("returning the error")
					Expect(err).Should(HaveOccurred())

					ctrl.Finish()
				})
			})
		})
	})
})
