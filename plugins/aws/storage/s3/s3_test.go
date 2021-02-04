package s3_service_test

import (
	"github.com/nitric-dev/membrane/plugins/aws/mocks"
	s3Plugin "github.com/nitric-dev/membrane/plugins/aws/storage/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("S3", func() {
	When("Put", func() {
		When("Given the S3 backend is available", func() {
			When("Creating an object in an existing bucket", func() {
				testPayload := []byte("Test")
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mocks.NewStorageClient([]*mocks.MockBucket{
					{
						Name: "my-bucket",
						Tags: map[string]string{
							"x-nitric-name": "my-bucket",
						},
					},
				}, &storage)

				storagePlugin, _ := s3Plugin.NewWithClient(mockStorageClient)
				It("Should successfully store the object", func() {
					err := storagePlugin.Put("my-bucket", "test-item", testPayload)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Storing the item")
					Expect(storage["my-bucket"]["test-item"]).To(BeEquivalentTo(testPayload))
				})
			})

			When("Creating an object in a non-existent bucket", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mocks.NewStorageClient([]*mocks.MockBucket{}, &storage)
				storagePlugin, _ := s3Plugin.NewWithClient(mockStorageClient)
				It("Should fail to store the item", func() {
					err := storagePlugin.Put("my-bucket", "test-item", []byte("Test"))
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
	When("Get", func() {
		When("The S3 backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func () {
					// Setup a mock bucket, with a single item
					storage := make(map[string]map[string][]byte)
					storage["test-bucket"] = make(map[string][]byte)
					storage["test-bucket"]["test-key"] = []byte("Test")
					mockStorageClient := mocks.NewStorageClient([]*mocks.MockBucket{
						{
							Name: "test-bucket",
							Tags: map[string]string{
								"x-nitric-name": "test-bucket",
							},
						},
					}, &storage)
					storagePlugin, _ := s3Plugin.NewWithClient(mockStorageClient)

					It("Should successfully retrieve the object", func() {
						object, err := storagePlugin.Get("test-bucket", "test-key")
						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning the item")
						Expect(object).To(Equal([]byte("Test")))
					})
				})
				When("The item doesn't exist", func () {

				})
			})
			When("The bucket doesn't exist", func() {

			})
		})
	})
	When("Delete", func() {
		When("The S3 backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func () {
					// Setup a mock bucket, with a single item
					storage := make(map[string]map[string][]byte)
					storage["test-bucket"] = make(map[string][]byte)
					storage["test-bucket"]["test-key"] = []byte("Test")
					mockStorageClient := mocks.NewStorageClient([]*mocks.MockBucket{
						{
							Name: "test-bucket",
							Tags: map[string]string{
								"x-nitric-name": "test-bucket",
							},
						},
					}, &storage)
					storagePlugin, _ := s3Plugin.NewWithClient(mockStorageClient)

					It("Should successfully delete the object", func() {
						err := storagePlugin.Delete("test-bucket", "test-key")
						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Deleting the item")
						Expect(storage["test-bucket"]["test-key"]).To(BeNil())
					})
				})
				When("The item doesn't exist", func () {

				})
			})
			When("The bucket doesn't exist", func() {

			})
		})
	})
})
