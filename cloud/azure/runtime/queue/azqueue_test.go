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
	"fmt"
	"time"

	azqueue "github.com/Azure/azure-storage-queue-go/azqueue"
	"github.com/golang/mock/gomock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_azqueue "github.com/nitrictech/nitric/cloud/azure/mocks/azqueue"
	queuepb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
)

var _ = Describe("Azqueue", func() {
	structPayload, err := structpb.NewStruct(map[string]interface{}{"Test": "Test"})
	Expect(err).To(BeNil())

	testStruct := &queuepb.QueueMessage{
		Content: &queuepb.QueueMessage_StructPayload{
			StructPayload: structPayload,
		},
	}

	testPayloadBytes, err := proto.Marshal(testStruct)
	testB64Payload := base64.StdEncoding.EncodeToString(testPayloadBytes)
	Expect(err).To(BeNil())

	Context("Send", func() {
		When("Azure returns a successfully response", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzqueue := mock_azqueue.NewMockAzqueueServiceUrlIface(crtl)
			mockQueue := mock_azqueue.NewMockAzqueueQueueUrlIface(crtl)
			mockMessages := mock_azqueue.NewMockAzqueueMessageUrlIface(crtl)
			// mockMessageId := mock_azqueue.NewMockAzqueueMessageIdUrlIface(crtl)

			queuePlugin := &AzqueueQueueService{
				client: mockAzqueue,
			}

			It("should successfully send the queue item(s)", func() {
				By("Retrieving the Queue URL for the requested queue")
				mockAzqueue.EXPECT().NewQueueURL("test-queue").Times(2).Return(mockQueue)

				By("Retrieving the Message URL of the requested queue")
				mockQueue.EXPECT().NewMessageURL().Times(2).Return(mockMessages)

				By("Calling Enqueue once on the Message URL with the expected options")
				mockMessages.EXPECT().Enqueue(
					gomock.Any(),
					testB64Payload,
					time.Duration(0),
					time.Duration(0),
				).Times(2).Return(&azqueue.EnqueueMessageResponse{}, nil)

				resp, err := queuePlugin.Enqueue(context.TODO(), &queuepb.QueueEnqueueRequest{
					QueueName: "test-queue",
					Messages: []*queuepb.QueueMessage{
						testStruct,
						testStruct,
					},
				})

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("Not returning failed tasks")
				Expect(len(resp.FailedMessages)).To(Equal(0))

				crtl.Finish()
			})
		})

		When("Failing to send one task", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzqueue := mock_azqueue.NewMockAzqueueServiceUrlIface(crtl)
			mockQueue := mock_azqueue.NewMockAzqueueQueueUrlIface(crtl)
			mockMessages := mock_azqueue.NewMockAzqueueMessageUrlIface(crtl)
			// mockMessageId := mock_azqueue.NewMockAzqueueMessageIdUrlIface(crtl)

			queuePlugin := &AzqueueQueueService{
				client: mockAzqueue,
			}

			It("should successfully send the queue item(s)", func() {
				By("Retrieving the Queue URL for the requested queue")
				mockAzqueue.EXPECT().NewQueueURL("test-queue").Times(2).Return(mockQueue)

				By("Retrieving the Message URL of the requested queue")
				mockQueue.EXPECT().NewMessageURL().Times(2).Return(mockMessages)

				By("Calling Enqueue once on the Message URL with the expected options")
				mockMessages.EXPECT().Enqueue(
					gomock.Any(),
					testB64Payload,
					time.Duration(0),
					time.Duration(0),
				).AnyTimes( /* Using AnyTimes because Times(2) doesn't work for multiple returns */
				).Return(nil, fmt.Errorf("a test error")).Return(&azqueue.EnqueueMessageResponse{}, nil)

				// testStruct, _ := structpb.NewStruct(map[string]interface{}{"testval": "testkey"})

				resp, err := queuePlugin.Enqueue(context.TODO(), &queuepb.QueueEnqueueRequest{
					QueueName: "test-queue",
					Messages: []*queuepb.QueueMessage{
						testStruct,
						testStruct,
					},
				})

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("Not returning failed tasks")
				Expect(len(resp.FailedMessages)).To(Equal(0))

				crtl.Finish()
			})
		})
	})

	Context("Receive", func() {
		When("Azure returns a successfully response", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzqueue := mock_azqueue.NewMockAzqueueServiceUrlIface(crtl)
			mockQueue := mock_azqueue.NewMockAzqueueQueueUrlIface(crtl)
			mockMessages := mock_azqueue.NewMockAzqueueMessageUrlIface(crtl)
			mockDequeueResp := mock_azqueue.NewMockDequeueMessagesResponseIface(crtl)
			// mockMessageId := mock_azqueue.NewMockAzqueueMessageIdUrlIface(crtl)

			queuePlugin := &AzqueueQueueService{
				client: mockAzqueue,
			}

			It("should successfully send the queue item(s)", func() {
				By("Retrieving the Queue URL for the requested queue")
				mockAzqueue.EXPECT().NewQueueURL("test-queue").Times(1).Return(mockQueue)

				By("Retrieving the Message URL of the requested queue")
				mockQueue.EXPECT().NewMessageURL().Times(1).Return(mockMessages)

				By("Calling Dequeue once on the Message URL with the expected options")
				mockMessages.EXPECT().Dequeue(
					gomock.Any(),   // ctx
					int32(1),       // depth
					30*time.Second, // visibility timeout - defaulted to 30 seconds
				).Times(1).Return(mockDequeueResp, nil)

				mockDequeueResp.EXPECT().NumMessages().AnyTimes().Return(int32(1))
				mockDequeueResp.EXPECT().Message(int32(0)).Times(1).Return(&azqueue.DequeuedMessage{
					ID: "testid",
					// InsertionTime:   time.Time{},
					// ExpirationTime:  time.Time{},
					PopReceipt:      "popreceipt",
					NextVisibleTime: time.Time{},
					DequeueCount:    0,
					Text:            testB64Payload,
				})

				resp, err := queuePlugin.Dequeue(context.TODO(), &queuepb.QueueDequeueRequest{
					QueueName: "test-queue",
					Depth:     1,
				})

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				By("Returning the dequeued task")
				Expect(len(resp.Messages)).To(Equal(1))
				Expect(resp.Messages[0].Message.GetStructPayload().AsMap()).To(Equal(map[string]interface{}{"Test": "Test"}))

				crtl.Finish()
			})
		})

		When("Azure returns an error", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzqueue := mock_azqueue.NewMockAzqueueServiceUrlIface(crtl)
			mockQueue := mock_azqueue.NewMockAzqueueQueueUrlIface(crtl)
			mockMessages := mock_azqueue.NewMockAzqueueMessageUrlIface(crtl)
			// mockDequeueResp := mock_azqueue.NewMockDequeueMessagesResponseIface(crtl)
			// mockMessageId := mock_azqueue.NewMockAzqueueMessageIdUrlIface(crtl)

			queuePlugin := &AzqueueQueueService{
				client: mockAzqueue,
			}

			It("should successfully send the queue item(s)", func() {
				By("Retrieving the Queue URL for the requested queue")
				mockAzqueue.EXPECT().NewQueueURL("test-queue").Times(1).Return(mockQueue)

				By("Retrieving the Message URL of the requested queue")
				mockQueue.EXPECT().NewMessageURL().Times(1).Return(mockMessages)

				By("Calling Dequeue once on the Message URL with the expected options")
				mockMessages.EXPECT().Dequeue(
					gomock.Any(),   // ctx
					int32(1),       // depth
					30*time.Second, // visibility timeout - defaulted to 30 seconds
				).Times(1).Return(nil, fmt.Errorf("a test error"))

				_, err := queuePlugin.Dequeue(context.TODO(), &queuepb.QueueDequeueRequest{
					QueueName: "test-queue",
					Depth:     1,
				})

				By("Returning an error")
				Expect(err).To(HaveOccurred())

				crtl.Finish()
			})
		})
	})

	Context("Complete", func() {
		When("Azure returns a successfully response", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzqueue := mock_azqueue.NewMockAzqueueServiceUrlIface(crtl)
			mockQueue := mock_azqueue.NewMockAzqueueQueueUrlIface(crtl)
			mockMessages := mock_azqueue.NewMockAzqueueMessageUrlIface(crtl)
			// mockDequeueResp := mock_azqueue.NewMockDequeueMessagesResponseIface(crtl)
			mockMessageId := mock_azqueue.NewMockAzqueueMessageIdUrlIface(crtl)

			queuePlugin := &AzqueueQueueService{
				client: mockAzqueue,
			}

			It("should successfully send the queue item(s)", func() {
				By("Retrieving the Queue URL for the requested queue")
				mockAzqueue.EXPECT().NewQueueURL("test-queue").Times(1).Return(mockQueue)

				By("Retrieving the Message URL of the requested queue")
				mockQueue.EXPECT().NewMessageURL().Times(1).Return(mockMessages)

				lease := AzureQueueItemLease{
					ID:         "testid",
					PopReceipt: "testreceipt",
				}
				leaseStr, _ := lease.String()

				By("Retrieving the Message ID URL specific to the dequeued task")
				mockMessages.EXPECT().NewMessageIDURL(azqueue.MessageID("testid")).Times(1).Return(mockMessageId)
				mockMessageId.EXPECT().Delete(gomock.Any(), azqueue.PopReceipt(lease.PopReceipt)).Times(1).Return(nil, nil)

				_, err := queuePlugin.Complete(context.TODO(), &queuepb.QueueCompleteRequest{
					QueueName: "test-queue",
					LeaseId:   leaseStr,
				})

				By("Not returning an error")
				Expect(err).ToNot(HaveOccurred())

				crtl.Finish()
			})
		})

		When("Azure returns an error", func() {
			crtl := gomock.NewController(GinkgoT())
			mockAzqueue := mock_azqueue.NewMockAzqueueServiceUrlIface(crtl)
			mockQueue := mock_azqueue.NewMockAzqueueQueueUrlIface(crtl)
			mockMessages := mock_azqueue.NewMockAzqueueMessageUrlIface(crtl)
			// mockDequeueResp := mock_azqueue.NewMockDequeueMessagesResponseIface(crtl)
			mockMessageId := mock_azqueue.NewMockAzqueueMessageIdUrlIface(crtl)

			queuePlugin := &AzqueueQueueService{
				client: mockAzqueue,
			}

			It("should successfully send the queue item(s)", func() {
				By("Retrieving the Queue URL for the requested queue")
				mockAzqueue.EXPECT().NewQueueURL("test-queue").Times(1).Return(mockQueue)

				By("Retrieving the Message URL of the requested queue")
				mockQueue.EXPECT().NewMessageURL().Times(1).Return(mockMessages)

				lease := AzureQueueItemLease{
					ID:         "testid",
					PopReceipt: "testreceipt",
				}
				leaseStr, _ := lease.String()

				By("Retrieving the Message ID URL specific to the dequeued task")
				mockMessages.EXPECT().NewMessageIDURL(azqueue.MessageID("testid")).Times(1).Return(mockMessageId)
				mockMessageId.EXPECT().Delete(gomock.Any(), azqueue.PopReceipt(lease.PopReceipt)).Times(1).Return(nil, fmt.Errorf("a test error"))

				_, err := queuePlugin.Complete(context.TODO(), &queuepb.QueueCompleteRequest{
					QueueName: "test-queue",
					LeaseId:   leaseStr,
				})

				By("Returning an error")
				Expect(err).To(HaveOccurred())

				crtl.Finish()
			})
		})
	})
})
