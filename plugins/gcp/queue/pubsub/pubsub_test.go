package pubsub_queue_service_test

import (
	"encoding/json"
	"fmt"

	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"github.com/nitric-dev/membrane/plugins/gcp/mocks"
	pubsub_queue_plugin "github.com/nitric-dev/membrane/plugins/gcp/queue/pubsub"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
)

var _ = Describe("Pubsub", func() {
	Context("BatchPush", func() {

		When("Publishing to a queue that exists", func() {
			mockPubsubClient := mocks.NewMockPubsubClient(
				mocks.MockPubsubOptions{
					Topics: []string{"test"},
				})
			queuePlugin := pubsub_queue_plugin.NewWithClient(mockPubsubClient)

			It("Should queue the Nitric Event", func() {
				resp, err := queuePlugin.SendBatch("test", []sdk.NitricTask{{
					ID:   "1234",
					PayloadType: "test-payload",
					Payload: map[string]interface{}{
						"Test": "Test",
					},
				}})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning no failed messages")
				Expect(resp.FailedMessages).To(HaveLen(0))

				By("The queue containing the published message")
				Expect(mockPubsubClient.PublishedMessages["test"]).To(HaveLen(1))
			})
		})

		When("Publishing to a queue that does not exist", func() {
			mockPubsubClient := mocks.NewMockPubsubClient(mocks.MockPubsubOptions{
				Topics:   []string{},
			})
			queuePlugin := pubsub_queue_plugin.NewWithClient(mockPubsubClient)
			It("Should return the messages that failed to publish", func() {
				_, err := queuePlugin.SendBatch("test", []sdk.NitricTask{{
					ID:   "1234",
					PayloadType: "test-payload",
					Payload: map[string]interface{}{
						"Test": "Test",
					},
				}})

				// It should still attempt
				By("Returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Pop", func() {

		When("Popping from a queue that exists", func() {
			When("There is a message on the queue", func() {
				mockId := "mockmessageid"
				mockReceiptHandle := "mockreceipthandle"
				jsonBytes, _ := json.Marshal(sdk.NitricTask{
					ID:   "mockrequestid",
					PayloadType: "mockpayloadtype",
					Payload:     map[string]interface{}{},
				})

				var mockMessage ifaces.Message = mocks.MockPubsubMessage{
					Id:        mockId,
					AckId:     mockReceiptHandle,
					DataBytes: jsonBytes,
				}

				mockPubsubClient := mocks.NewMockPubsubClient(mocks.MockPubsubOptions{
					Topics: []string{"mock-queue"},
					Messages: map[string][]ifaces.Message{
						"mock-queue": {
							mockMessage,
						},
					},
				})
				queuePlugin := pubsub_queue_plugin.NewWithClients(mockPubsubClient, func(ctx context.Context, opts ...option.ClientOption) (ifaces.SubscriberClient, error) {
					return mocks.MockBaseClient{
						Messages: map[string][]ifaces.Message{
							"mock-queue": {
								mockMessage,
							},
						},
					}, nil
				})

				It("Should pop the message", func() {
					items, err := queuePlugin.Receive(sdk.ReceiveOptions{
						QueueName: "mock-queue",
						Depth:     nil,
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning no failed messages")
					Expect(len(items)).To(Equal(1))

					//By("Removing the popped messages")
				})
			})
			When("There are no messages on the queue", func() {

			})
		})

		When("Publishing to a queue that does not exist", func() {
			//mockPubsubClient := mocks.NewMockPubsubClient([]string{})
			//queuePlugin := pubsub_queue_plugin.NewWithClient(mockPubsubClient)
			//It("Should return the messages that failed to publish", func() {
			//	_, err := queuePlugin.BatchPush("test", []sdk.NitricTask{{
			//		ID:   "1234",
			//		PayloadType: "test-payload",
			//		Payload: map[string]interface{}{
			//			"Test": "Test",
			//		},
			//	}})
			//
			//	// It should still attempt
			//	By("Returning an error")
			//	Expect(err).Should(HaveOccurred())
			//})
		})
	})

	Context("Complete", func() {
		When("Pubsub acknowledge request succeeds", func() {
			mockPubsubClient := mocks.NewMockPubsubClient(mocks.MockPubsubOptions{
				Topics: []string{"mock-queue"},
			})
			queuePlugin := pubsub_queue_plugin.NewWithClients(mockPubsubClient, func(ctx context.Context, opts ...option.ClientOption) (ifaces.SubscriberClient, error) {
				return mocks.MockBaseClient{}, nil
			})

			It("Should not return an error", func() {
				err := queuePlugin.Complete("mock-queue", "test-id")

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
		When("Pubsub acknowledge request errors", func() {
			mockPubsubClient := mocks.NewMockPubsubClient(mocks.MockPubsubOptions{
				Topics: []string{"mock-queue"},
			})
			queuePlugin := pubsub_queue_plugin.NewWithClients(mockPubsubClient, func(ctx context.Context, opts ...option.ClientOption) (ifaces.SubscriberClient, error) {
				return mocks.MockBaseClient{
					CompleteError: fmt.Errorf("mock complete error"),
				}, nil
			})

			It("Should return an error", func() {
				err := queuePlugin.Complete("mock-queue", "test-id")

				By("Not returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
