package pubsub_queue_plugin_test

import (
	"github.com/nitric-dev/membrane/plugins/gcp/mocks"
	pubsub_queue_plugin "github.com/nitric-dev/membrane/plugins/gcp/queue/pubsub"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Pubsub", func() {
	Context("Push", func() {

		When("Publishing to a queue that exists", func() {
			mockPubsubClient := mocks.NewMockPubsubClient([]string{"test"})
			queuePlugin := pubsub_queue_plugin.NewWithClient(mockPubsubClient)

			It("Should queue the Nitric Event", func() {
				resp, err := queuePlugin.Push("test", []*sdk.NitricEvent{&sdk.NitricEvent{
					RequestId:   "1234",
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
			mockPubsubClient := mocks.NewMockPubsubClient([]string{})
			queuePlugin := pubsub_queue_plugin.NewWithClient(mockPubsubClient)
			It("Should return the messages that failed to publish", func() {
				_, err := queuePlugin.Push("test", []*sdk.NitricEvent{&sdk.NitricEvent{
					RequestId:   "1234",
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
})
