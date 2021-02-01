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
			evts := []sdk.NitricEvent{evt}
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
			evts := []sdk.NitricEvent{evt}
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

	Context("Pop", func() {
		When("The queue is empty", func() {
			evtsBytes, _ := json.Marshal([]sdk.NitricEvent{})
			mockStorageDriver := mocks.NewMockStorageDriver(&mocks.MockStorageDriverOptions{
				StoredItems: map[string][]byte{
					"/nitric/queues/test": evtsBytes,
				},
			})
			queuePlugin, _ := queue_plugin.NewWithStorageDriver(mockStorageDriver)

			It("Should return an empty slice of queue items", func() {
				depth := uint32(10)
				items, err := queuePlugin.Pop(sdk.PopOptions{
					QueueName: "test",
					Depth: &depth,
				})
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning an empty slice")
				Expect(items).To(HaveLen(0))
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
			evts := []sdk.NitricEvent{evt}
			evtsBytes, _ := json.Marshal(evts)
			mockStorageDriver := mocks.NewMockStorageDriver(&mocks.MockStorageDriverOptions{
				StoredItems: map[string][]byte{
					"/nitric/queues/test": evtsBytes,
				},
			})
			queuePlugin, _ := queue_plugin.NewWithStorageDriver(mockStorageDriver)

			It("Should append to the existing queue", func() {
				depth := uint32(10)
				items, err := queuePlugin.Pop(sdk.PopOptions{
					QueueName: "test",
					Depth: &depth,
				})
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning 1 item")
				Expect(items).To(HaveLen(1))

				var messages []sdk.NitricEvent
				bytes := mockStorageDriver.GetStoredItems()["/nitric/queues/test"]
				json.Unmarshal(bytes, &messages)
				By("Having no remaining messages on the Queue")
				Expect(messages).To(HaveLen(0))
			})
		})

		When("The queue depth is 15", func() {
			evt := sdk.NitricEvent{
				RequestId:   "1234",
				PayloadType: "test-payload",
				Payload: map[string]interface{}{
					"Test": "Test",
				},
			}
			evts := []sdk.NitricEvent{}

			// Add 15 items to the queue
			for i := 0; i < 15; i++ {
				evts = append(evts, evt)
			}

			evtsBytes, _ := json.Marshal(evts)
			mockStorageDriver := mocks.NewMockStorageDriver(&mocks.MockStorageDriverOptions{
				StoredItems: map[string][]byte{
					"/nitric/queues/test": evtsBytes,
				},
			})
			queuePlugin, _ := queue_plugin.NewWithStorageDriver(mockStorageDriver)

			When("Requested depth is 10", func() {
				It("Should return 10 items", func() {
					depth := uint32(10)
					items, err := queuePlugin.Pop(sdk.PopOptions{
						QueueName: "test",
						Depth: &depth,
					})
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning 10 item")
					Expect(items).To(HaveLen(10))

					var messages []sdk.NitricEvent
					bytes := mockStorageDriver.GetStoredItems()["/nitric/queues/test"]
					json.Unmarshal(bytes, &messages)
					By("Having 5 remaining messages on the Queue")
					Expect(messages).To(HaveLen(5))
				})
			})
		})
	})
})
