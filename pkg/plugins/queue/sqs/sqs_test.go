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

package sqs_service_test

import (
	"encoding/json"
	"fmt"

	"github.com/nitric-dev/membrane/pkg/plugins/eventing"
	sqs_service "github.com/nitric-dev/membrane/pkg/plugins/queue/sqs"
	mocks_sqs "github.com/nitric-dev/membrane/tests/mocks/sqs"

	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sqs", func() {
	// Tests for the BatchPush method
	Context("BatchPush", func() {
		When("Publishing to a queue that exists", func() {
			sqsMock := mocks_sqs.NewMockSqs(&mocks_sqs.MockSqsOptions{
				Queues: []string{"test"},
			})
			plugin := sqs_service.NewWithClient(sqsMock)

			It("Should publish the message", func() {
				_, err := plugin.SendBatch("test", []queue.NitricTask{
					{
						ID:          "1234",
						PayloadType: "test-payload",
						Payload: map[string]interface{}{
							"Test": "Test",
						},
					},
				})

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("Publishing to a queue that doesn't exist", func() {
			sqsMock := mocks_sqs.NewMockSqs(&mocks_sqs.MockSqsOptions{})
			plugin := sqs_service.NewWithClient(sqsMock)

			It("Should fail to publish the message", func() {
				_, err := plugin.SendBatch("test", []queue.NitricTask{
					{
						ID:          "1234",
						PayloadType: "test-payload",
						Payload: map[string]interface{}{
							"Test": "Test",
						},
					},
				})

				Expect(err).Should(HaveOccurred())
			})
		})
	})

	// Tests for the Pop method
	Context("Pop", func() {
		When("Popping from a queue that exists", func() {
			When("There is a message on the queue", func() {
				mockId := "mockmessageid"
				mockReceiptHandle := "mockreceipthandle"
				jsonBytes, _ := json.Marshal(eventing.NitricEvent{
					ID:          "mockrequestid",
					PayloadType: "mockpayloadtype",
					Payload:     map[string]interface{}{},
				})
				mockEventJson := string(jsonBytes)

				sqsMock := mocks_sqs.NewMockSqs(&mocks_sqs.MockSqsOptions{
					Queues: []string{"mock-queue"},
					Messages: map[string][]*mocks_sqs.Message{
						"mock-queue": {
							{
								Id:            &mockId,
								ReceiptHandle: &mockReceiptHandle,
								Body:          &mockEventJson,
							},
						},
					},
				})
				plugin := sqs_service.NewWithClient(sqsMock)

				depth := uint32(10)

				It("Should pop the message", func() {
					msg, err := plugin.Receive(queue.ReceiveOptions{
						QueueName: "mock-queue",
						Depth:     &depth,
					})

					Expect(msg).To(HaveLen(1))

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			When("There are no messages on the queue", func() {
				sqsMock := mocks_sqs.NewMockSqs(&mocks_sqs.MockSqsOptions{
					Queues: []string{"mock-queue"},
					Messages: map[string][]*mocks_sqs.Message{
						// Queue with empty message slice
						"mock-queue": make([]*mocks_sqs.Message, 0),
					},
				})
				plugin := sqs_service.NewWithClient(sqsMock)
				depth := uint32(10)

				It("Should pop the message", func() {
					msg, err := plugin.Receive(queue.ReceiveOptions{
						QueueName: "mock-queue",
						Depth:     &depth,
					})

					Expect(len(msg)).To(Equal(0))

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})

		//When("Popping from a queue that doesn't exist", func() {
		When("Popping from a queue that doesn't exist", func() {
			sqsMock := mocks_sqs.NewMockSqs(&mocks_sqs.MockSqsOptions{
				Queues: []string{},
			})
			plugin := sqs_service.NewWithClient(sqsMock)

			depth := uint32(10)

			It("Should return an error", func() {
				_, err := plugin.Receive(queue.ReceiveOptions{
					QueueName: "non-existent-queue",
					Depth:     &depth,
				})

				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	// Tests for the Complete method
	Context("Complete", func() {
		When("The message is successfully deleted from SQS", func() {
			// No errors set on mock, 'complete' won't return an error.
			sqsMock := mocks_sqs.NewMockSqs(&mocks_sqs.MockSqsOptions{
				Queues: []string{"test-queue"},
			})
			plugin := sqs_service.NewWithClient(sqsMock)

			It("Should not return an error", func() {
				err := plugin.Complete("test-queue", "test-id")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("The message fails to delete from SQS", func() {
			// No errors set on mock, 'complete' won't return an error.
			sqsMock := mocks_sqs.NewMockSqs(&mocks_sqs.MockSqsOptions{
				CompleteError: fmt.Errorf("mock complete error"),
			})
			plugin := sqs_service.NewWithClient(sqsMock)

			It("Should return an error", func() {
				err := plugin.Complete("test-queue", "test-id")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
