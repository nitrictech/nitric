package sqs_service_test

import (
	"encoding/json"
	"fmt"

	"github.com/nitric-dev/membrane/plugins/aws/mocks"
	sqs_plugin "github.com/nitric-dev/membrane/plugins/aws/queue/sqs"
	"github.com/nitric-dev/membrane/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sqs", func() {
	// Tests for the BatchPush method
	Context("BatchPush", func() {
		When("Publishing to a queue that exists", func() {
			sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{
				Queues: []string{"test"},
			})
			plugin := sqs_plugin.NewWithClient(sqsMock)

			It("Should publish the message", func() {
				_, err := plugin.SendBatch("test", []sdk.NitricTask{
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
			sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{})
			plugin := sqs_plugin.NewWithClient(sqsMock)

			It("Should fail to publish the message", func() {
				_, err := plugin.SendBatch("test", []sdk.NitricTask{
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
				jsonBytes, _ := json.Marshal(sdk.NitricEvent{
					ID:          "mockrequestid",
					PayloadType: "mockpayloadtype",
					Payload:     map[string]interface{}{},
				})
				mockEventJson := string(jsonBytes)

				sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{
					Queues: []string{"mock-queue"},
					Messages: map[string][]*mocks.Message{
						"mock-queue": {
							{
								Id:            &mockId,
								ReceiptHandle: &mockReceiptHandle,
								Body:          &mockEventJson,
							},
						},
					},
				})
				plugin := sqs_plugin.NewWithClient(sqsMock)

				depth := uint32(10)

				It("Should pop the message", func() {
					msg, err := plugin.Receive(sdk.ReceiveOptions{
						QueueName: "mock-queue",
						Depth:     &depth,
					})

					Expect(msg).To(HaveLen(1))

					Expect(err).ShouldNot(HaveOccurred())
				})
			})
			When("There are no messages on the queue", func() {
				sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{
					Queues: []string{"mock-queue"},
					Messages: map[string][]*mocks.Message{
						// Queue with empty message slice
						"mock-queue": make([]*mocks.Message, 0),
					},
				})
				plugin := sqs_plugin.NewWithClient(sqsMock)
				depth := uint32(10)

				It("Should pop the message", func() {
					msg, err := plugin.Receive(sdk.ReceiveOptions{
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
			sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{
				Queues: []string{},
			})
			plugin := sqs_plugin.NewWithClient(sqsMock)

			depth := uint32(10)

			It("Should return an error", func() {
				_, err := plugin.Receive(sdk.ReceiveOptions{
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
			sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{
				Queues: []string{"test-queue"},
			})
			plugin := sqs_plugin.NewWithClient(sqsMock)

			It("Should not return an error", func() {
				err := plugin.Complete("test-queue", "test-id")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("The message fails to delete from SQS", func() {
			// No errors set on mock, 'complete' won't return an error.
			sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{
				CompleteError: fmt.Errorf("mock complete error"),
			})
			plugin := sqs_plugin.NewWithClient(sqsMock)

			It("Should return an error", func() {
				err := plugin.Complete("test-queue", "test-id")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
