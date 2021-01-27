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

	Context("Get", func() {
		When("The Google Cloud Storage Backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
					storage := make(map[string]map[string][]byte)
					storage["test-bucket"] = make(map[string][]byte)
					storage["test-bucket"]["test-key"] = []byte("Test")
					mockStorageClient := mocks.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_plugin.NewWithClient(mockStorageClient)

					It("Should retrieve the item", func() {
						item, err := storagePlugin.Get("test-bucket", "test-key")

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning the item")
						Expect(item).To(Equal([]byte("Test")))
					})
				})

				When("The item doesn't exist", func() {
					storage := make(map[string]map[string][]byte)
					mockStorageClient := mocks.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_plugin.NewWithClient(mockStorageClient)

					It("Should return an error", func() {
						item, err := storagePlugin.Get("test-bucket", "test-key")

						By("Returning an error")
						Expect(err).Should(HaveOccurred())

						By("Not returning the item")
						Expect(item).To(BeNil())
					})
				})
			})

			When("The bucket doesn't exist", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mocks.NewStorageClient([]string{}, &storage)
				storagePlugin, _ := storage_plugin.NewWithClient(mockStorageClient)

				It("Should return an error", func() {
					item, err := storagePlugin.Get("test-bucket", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())

					By("Not returning the item")
					Expect(item).To(BeNil())
				})
			})
		})
	})

	Context("Delete", func() {
		When("The Google Cloud Storage Backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
					storage := make(map[string]map[string][]byte)
					storage["test-bucket"] = make(map[string][]byte)
					storage["test-bucket"]["test-key"] = []byte("Test")
					mockStorageClient := mocks.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_plugin.NewWithClient(mockStorageClient)

					It("Should delete the item", func() {
						err := storagePlugin.Delete("test-bucket", "test-key")

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Deleting the item")
						Expect(storage["test-bucket"]["test-key"]).To(BeNil())
					})
				})

				When("The item doesn't exist", func() {
					storage := make(map[string]map[string][]byte)
					mockStorageClient := mocks.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_plugin.NewWithClient(mockStorageClient)

					// Since no item existed to begin with, no error is thrown deleting it.
					It("Should not return an error", func() {
						err := storagePlugin.Delete("test-bucket", "test-key")

						By("Not returning an error")
						Expect(err).Should(BeNil())
					})
				})
			})

			When("The bucket doesn't exist", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mocks.NewStorageClient([]string{}, &storage)
				storagePlugin, _ := storage_plugin.NewWithClient(mockStorageClient)

				It("Should return an error", func() {
					err := storagePlugin.Delete("test-bucket", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
})
