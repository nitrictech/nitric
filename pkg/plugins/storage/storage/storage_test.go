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

package storage_service_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	storage_service "github.com/nitrictech/nitric/pkg/plugins/storage/storage"
	mock_gcp_storage "github.com/nitrictech/nitric/tests/mocks/gcp_storage"
)

var _ = Describe("Storage", func() {
	Context("Write", func() {
		When("GCloud Storage Backend is available", func() {
			When("Writing to a bucket that exists", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mock_gcp_storage.NewStorageClient([]string{"my-bucket"}, &storage)
				mockStorageServer, _ := storage_service.NewWithClient(mockStorageClient)
				testPayload := []byte("Test")

				It("Should store the item", func() {
					err := mockStorageServer.Write("my-bucket", "test-file", testPayload)

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Storing the sent item under the given key")
					Expect(storage["my-bucket"]["test-file"]).To(BeEquivalentTo(testPayload))
				})
			})

			When("Writing to a Bucket that does not exist", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mock_gcp_storage.NewStorageClient([]string{}, &storage)
				mockStorageServer, _ := storage_service.NewWithClient(mockStorageClient)
				testPayload := []byte("Test")

				It("Should fail to store the item", func() {
					err := mockStorageServer.Write("my-bucket", "test-file", testPayload)

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})

	Context("Read", func() {
		When("The Google Cloud Storage Backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
					storage := make(map[string]map[string][]byte)
					storage["test-bucket"] = make(map[string][]byte)
					storage["test-bucket"]["test-key"] = []byte("Test")
					mockStorageClient := mock_gcp_storage.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("Should retrieve the item", func() {
						item, err := storagePlugin.Read("test-bucket", "test-key")

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning the item")
						Expect(item).To(Equal([]byte("Test")))
					})
				})

				When("The item doesn't exist", func() {
					storage := make(map[string]map[string][]byte)
					mockStorageClient := mock_gcp_storage.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("Should return an error", func() {
						item, err := storagePlugin.Read("test-bucket", "test-key")

						By("Returning an error")
						Expect(err).Should(HaveOccurred())

						By("Not returning the item")
						Expect(item).To(BeNil())
					})
				})
			})

			When("The bucket doesn't exist", func() {
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mock_gcp_storage.NewStorageClient([]string{}, &storage)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("Should return an error", func() {
					item, err := storagePlugin.Read("test-bucket", "test-key")

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
					mockStorageClient := mock_gcp_storage.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

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
					mockStorageClient := mock_gcp_storage.NewStorageClient([]string{"test-bucket"}, &storage)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

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
				mockStorageClient := mock_gcp_storage.NewStorageClient([]string{}, &storage)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("Should return an error", func() {
					err := storagePlugin.Delete("test-bucket", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
})
