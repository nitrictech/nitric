package storage_plugin_test

import (
	"github.com/nitric-dev/membrane/plugins/gcp/mocks"
	storage_plugin "github.com/nitric-dev/membrane/plugins/gcp/storage/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Storage", func() {
	Context("Put", func() {
		When("GCloud Storage Backend is available", func() {
			When("Writing to a bucket that exists", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mocks.NewStorageClient([]string{"my-bucket"}, &storage)
				mockStorageServer, _ := storage_plugin.NewWithClient(mockStorageClient)
				testPayload := []byte("Test")

				It("Should store the item", func() {
					err := mockStorageServer.Put("my-bucket", "test-file", testPayload)

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Storing the sent item under the given key")
					Expect(storage["my-bucket"]["test-file"]).To(BeEquivalentTo(testPayload))
				})
			})

			When("Writing to a Bucket that does not exist", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mocks.NewStorageClient([]string{}, &storage)
				mockStorageServer, _ := storage_plugin.NewWithClient(mockStorageClient)
				testPayload := []byte("Test")

				It("Should fail to store the item", func() {
					err := mockStorageServer.Put("my-bucket", "test-file", testPayload)

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
})
