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
	"fmt"
	"os"

	"cloud.google.com/go/storage"
	storage_service "github.com/nitric-dev/membrane/pkg/plugins/storage/storage"
	mock_gcp_storage "github.com/nitric-dev/membrane/tests/mocks/gcp_storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2/jwt"
)

type TestUtils struct{}

func (u *TestUtils) JWTConfigFromJSON(jsonKey []byte) (*jwt.Config, error) {
	return &jwt.Config{
		Email:      "test@test.com",
		PrivateKey: []byte("iamsuperprivateandalsoakey!"),
	}, nil
}

func (u *TestUtils) SignedURL(bucket string, key string, opts *storage.SignedURLOptions) (string, error) {
	return fmt.Sprintf("https://presignedurl/%s/", key), nil
}

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

	Context("PreSignedURL", func() {
		When("The Google Cloud Storage Backend is available", func() {
			storageMap := make(map[string]map[string][]byte)
			mockStorageClient := mock_gcp_storage.NewStorageClient([]string{"test-bucket"}, &storageMap)
			storagePlugin, _ := storage_service.NewWithUtils(mockStorageClient, &TestUtils{})
			When("A service account with sufficient permissions is set up", func() {
				os.Setenv("SERVICE_ACCOUNT_LOCATION", "./service_account.json")
				os.Create("service_account.json")
				When("The item exists", func() {
					It("Should return the presigned url", func() {
						url, err := storagePlugin.PreSignUrl("bucket-name", "image-name.jpg", 0, 100)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(url).Should(Equal("https://presignedurl/image-name.jpg/"))
					})
				})
				When("The item does not exist", func() {
					It("Should still return the presigned url", func() {
						url, err := storagePlugin.PreSignUrl("bucket-name", "image-that-doesnt-exist.jpg", 1, 100)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(url).Should(Equal("https://presignedurl/image-that-doesnt-exist.jpg/"))
					})
				})
				When("Supplying an empty bucket name", func() {
					It("Should return an error", func() {
						_, err := storagePlugin.PreSignUrl("", "image-name", 0, 300)
						Expect(err).Should(HaveOccurred())
					})
				})
				When("Supplying an empty image name", func() {
					It("Should return an error", func() {
						_, err := storagePlugin.PreSignUrl("bucket-name", "", 0, 300)
						Expect(err).Should(HaveOccurred())
					})
				})
				When("Supplying a non 0 or 1 operation", func() {
					It("Should default to a read operation", func() {
						url, err := storagePlugin.PreSignUrl("bucket-name", "image-name.jpg", 2, 300)
						Expect(err).ShouldNot(HaveOccurred())
						Expect(url).Should(Equal("https://presignedurl/image-name.jpg/"))
					})
				})
				When("Supplying an expiry over 7 days (604800 seconds)", func() {
					It("Should return an error", func() {
						_, err := storagePlugin.PreSignUrl("bucket-name", "image-name.jpg", 0, 604801)
						Expect(err).Should(HaveOccurred())
					})
				})
			})
		})
	})
})
