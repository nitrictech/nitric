package queue_plugin_test

import (
	"encoding/json"

	"github.com/nitric-dev/membrane/plugins/dev/mocks"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	queue_plugin "github.com/nitric-dev/membrane/plugins/dev/queue"
)

var _ = Describe("Queue", func() {
	Context("Push", func() {
		When("The queue is empty", func() {
			mockStorageDriver := mocks.NewMockStorageDriver(&mocks.MockStorageDriverOptions{})
			queuePlugin, _ := queue_plugin.NewWithStorageDriver(mockStorageDriver)
			evt := sdk.NitricEvent{
				RequestId:   "1234",
				PayloadType: "test-payload",
				Payload: map[string]interface{}{
					"Test": "Test",
				},
			}
			evts := []*sdk.NitricEvent{&evt}
			evtsBytes, _ := json.Marshal([]sdk.NitricEvent{evt})
			It("Should store the events in the queue", func() {
				resp, err := queuePlugin.Push("test", evts)
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning No failed messages")
				Expect(resp.FailedMessages).To(BeEmpty())

				By("Storing the sent message, in the given queue")
				Expect(mockStorageDriver.GetStoredItems()["/nitric/queues/test"]).ToNot(BeNil())

				By("Storing the content of the given message")
				Expect(mockStorageDriver.GetStoredItems()["/nitric/queues/test"]).To(BeEquivalentTo(evtsBytes))
			})
		})

		When("The queue is not empty", func() {
			evt := sdk.NitricEvent{
				RequestId:   "1234",
				PayloadType: "test-payload",
				Payload: map[string]interface{}{
					"Test": "Test",
				},
			}
			evts := []*sdk.NitricEvent{&evt}
			evtsBytes, _ := json.Marshal([]sdk.NitricEvent{evt})
			mockStorageDriver := mocks.NewMockStorageDriver(&mocks.MockStorageDriverOptions{
				StoredItems: map[string][]byte{
					"/nitric/queues/test": evtsBytes,
				},
			})
			queuePlugin, _ := queue_plugin.NewWithStorageDriver(mockStorageDriver)

			It("Should append to the existing queue", func() {
				resp, err := queuePlugin.Push("test", evts)
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Having no Failed Messages")
				Expect(resp.FailedMessages).To(BeEmpty())

				By("Storing the sent message, in the given queue")
				Expect(mockStorageDriver.GetStoredItems()["/nitric/queues/test"]).ToNot(BeNil())

				var messages []sdk.NitricEvent
				bytes := mockStorageDriver.GetStoredItems()["/nitric/queues/test"]
				json.Unmarshal(bytes, &messages)
				By("Having 2 messages on the Queue")
				Expect(messages).To(HaveLen(2))
			})
		})
	})
})
