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

package s3_service_test

import (
	mocks "github.com/nitric-dev/membrane/mocks/s3"
	s3Plugin "github.com/nitric-dev/membrane/plugins/storage/s3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("S3", func() {
	When("Write", func() {
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
					err := storagePlugin.Write("my-bucket", "test-item", testPayload)
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
					err := storagePlugin.Write("my-bucket", "test-item", []byte("Test"))
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})
	When("Read", func() {
		When("The S3 backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
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
						object, err := storagePlugin.Read("test-bucket", "test-key")
						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning the item")
						Expect(object).To(Equal([]byte("Test")))
					})
				})
				When("The item doesn't exist", func() {

				})
			})
			When("The bucket doesn't exist", func() {

			})
		})
	})
	When("Delete", func() {
		When("The S3 backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
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
				When("The item doesn't exist", func() {

				})
			})
			When("The bucket doesn't exist", func() {

			})
		})
	})
})
