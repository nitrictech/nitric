package sqs_plugin_test

import (
	"github.com/nitric-dev/membrane/plugins/aws/mocks"
	sqs_plugin "github.com/nitric-dev/membrane/plugins/aws/queue/sqs"
	"github.com/nitric-dev/membrane/plugins/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sqs", func() {
	// Tests for the Push method
	Context("Push", func() {
		When("Publishing to a queue that exists", func() {
			sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{
				Queues: []string{"test"},
			})
			plugin := sqs_plugin.NewWithClient(sqsMock)

			It("Should fail to publish the message", func() {
				_, err := plugin.Push("test", []*sdk.NitricEvent{
					&sdk.NitricEvent{
						RequestId:   "1234",
						PayloadType: "test-payload",
						Payload: map[string]interface{}{
							"Test": "Test",
						},
					},
				})

				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("Publishing to a queue that does not exist", func() {
			sqsMock := mocks.NewMockSqs(&mocks.MockSqsOptions{})
			plugin := sqs_plugin.NewWithClient(sqsMock)

			It("Should fail to publish the message", func() {
				_, err := plugin.Push("test", []*sdk.NitricEvent{
					&sdk.NitricEvent{
						RequestId:   "1234",
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

})
