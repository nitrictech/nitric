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
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/iterator"

	storage_mock "github.com/nitrictech/nitric/mocks/gcp_storage"
	plugin "github.com/nitrictech/nitric/pkg/plugins/storage"
	storage_service "github.com/nitrictech/nitric/pkg/plugins/storage/storage"
)

var _ = Describe("Storage", func() {
	Context("Write", func() {
		When("GCloud Storage Backend is available", func() {
			When("Writing to a bucket that exists", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				mockBucket := storage_mock.NewMockBucketHandle(ctrl)
				mockObject := storage_mock.NewMockObjectHandle(ctrl)
				mockWriter := storage_mock.NewMockWriter(ctrl)
				mockStorageServer, _ := storage_service.NewWithClient(mockStorageClient)
				testPayload := []byte("Test")

				It("Should store the item", func() {
					By("The bucket existing")
					gomock.InOrder(
						mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
							Labels: map[string]string{
								"x-nitric-name": "my-bucket",
							},
							Name: "my-bucket-1234",
						}, nil),
						mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
					)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
					mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

					By("The object reference being correct")
					mockBucket.EXPECT().Object("test-file").Return(mockObject)

					By("The writer being called on the object handle")
					mockObject.EXPECT().NewWriter(gomock.Any()).Return(mockWriter)

					By("The bytes being written")
					mockWriter.EXPECT().Write(testPayload).Times(1)
					mockWriter.EXPECT().Close().Times(1)

					err := mockStorageServer.Write(context.TODO(), "my-bucket", "test-file", testPayload)

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					ctrl.Finish()
				})
			})

			When("Writing to a Bucket that does not exist", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				mockStorageServer, _ := storage_service.NewWithClient(mockStorageClient)

				testPayload := []byte("Test")

				It("Should fail to store the item", func() {
					By("The bucket not existing")
					mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)

					err := mockStorageServer.Write(context.TODO(), "my-bucket", "test-file", testPayload)

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					ctrl.Finish()
				})
			})
		})
	})

	Context("Read", func() {
		When("The Google Cloud Storage Backend is available", func() {
			When("The bucket exists", func() {
				When("The item exists", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
					mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
					mockBucket := storage_mock.NewMockBucketHandle(ctrl)
					mockObject := storage_mock.NewMockObjectHandle(ctrl)
					mockReader := storage_mock.NewMockReader(ctrl)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("Should retrieve the item", func() {
						By("the bucket existing")
						gomock.InOrder(
							mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
								Labels: map[string]string{
									"x-nitric-name": "test-bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being valid")
						mockBucket.EXPECT().Object("test-key").Return(mockObject)
						mockObject.EXPECT().NewReader(gomock.Any()).Return(mockReader, nil)

						By("the object reader being called")
						mockReader.EXPECT().Read(gomock.Any()).Return(0, io.EOF)
						mockReader.EXPECT().Close().Times(1)

						item, err := storagePlugin.Read(context.TODO(), "test-bucket", "test-key")

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning read content")
						Expect(item).To(HaveLen(0))
					})
				})

				When("The item doesn't exist", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
					mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
					mockBucket := storage_mock.NewMockBucketHandle(ctrl)
					mockObject := storage_mock.NewMockObjectHandle(ctrl)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("Should return an error", func() {
						By("the bucket existing")
						gomock.InOrder(
							mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
								Labels: map[string]string{
									"x-nitric-name": "test-bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being invalid")
						mockBucket.EXPECT().Object("test-key").Return(mockObject)
						mockObject.EXPECT().NewReader(gomock.Any()).Return(nil, fmt.Errorf("mock-error"))

						item, err := storagePlugin.Read(context.TODO(), "test-bucket", "test-key")

						By("Returning an error")
						Expect(err).Should(HaveOccurred())

						By("Not returning the item")
						Expect(item).To(BeNil())
					})
				})
			})

			When("The bucket doesn't exist", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("Should return an error", func() {
					By("The bucket not existing")
					mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)

					item, err := storagePlugin.Read(context.TODO(), "test-bucket", "test-key")

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
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
					mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
					mockBucket := storage_mock.NewMockBucketHandle(ctrl)
					mockObject := storage_mock.NewMockObjectHandle(ctrl)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("Should retrieve the item", func() {
						By("the bucket existing")
						gomock.InOrder(
							mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
								Labels: map[string]string{
									"x-nitric-name": "test-bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being valid")
						mockBucket.EXPECT().Object("test-key").Return(mockObject)

						By("the object delete being called")
						mockObject.EXPECT().Delete(gomock.Any()).Return(nil)

						err := storagePlugin.Delete(context.TODO(), "test-bucket", "test-key")

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})

				When("The item doesn't exist", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
					mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
					mockBucket := storage_mock.NewMockBucketHandle(ctrl)
					mockObject := storage_mock.NewMockObjectHandle(ctrl)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("Should return an error", func() {
						By("the bucket existing")
						gomock.InOrder(
							mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
								Labels: map[string]string{
									"x-nitric-name": "test-bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being invalid")
						mockBucket.EXPECT().Object("test-key").Return(mockObject)
						mockObject.EXPECT().Delete(gomock.Any()).Return(fmt.Errorf("mock-error"))

						err := storagePlugin.Delete(context.TODO(), "test-bucket", "test-key")

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
					})
				})
			})

			When("The bucket doesn't exist", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("Should return an error", func() {
					By("The bucket not existing")
					mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)

					err := storagePlugin.Delete(context.TODO(), "test-bucket", "test-key")

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})

	Context("SignedUrl", func() {
		When("The bucket exists", func() {
			When("The item exists", func() {
				When("requesting a read url", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
					mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
					mockBucket := storage_mock.NewMockBucketHandle(ctrl)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("should return a readable url", func() {
						By("the bucket existing")
						gomock.InOrder(
							mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
								Labels: map[string]string{
									"x-nitric-name": "test-bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being valid")
						mockBucket.EXPECT().SignedURL("test-key", gomock.Any()).Return("http://example.com", nil)

						url, err := storagePlugin.PreSignUrl(context.TODO(), "test-bucket", "test-key", plugin.READ, 6000)

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("The url being returned")
						Expect(url).To(Equal("http://example.com"))

						ctrl.Finish()
					})
				})

				When("requesting a write url", func() {
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
					mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
					mockBucket := storage_mock.NewMockBucketHandle(ctrl)
					storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

					It("should return a writable url", func() {
						By("the bucket existing")
						gomock.InOrder(
							mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
								Labels: map[string]string{
									"x-nitric-name": "test-bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being valid")
						mockBucket.EXPECT().SignedURL("test-key", gomock.Any()).Return("http://example.com", nil)

						url, err := storagePlugin.PreSignUrl(context.TODO(), "test-bucket", "test-key", plugin.WRITE, 6000)

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("The url being returned")
						Expect(url).To(Equal("http://example.com"))

						ctrl.Finish()
					})
				})
			})

			When("The item doesn't exist", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				mockBucket := storage_mock.NewMockBucketHandle(ctrl)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("Should return an error", func() {
					By("the bucket existing")
					gomock.InOrder(
						mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
							Labels: map[string]string{
								"x-nitric-name": "test-bucket",
							},
							Name: "my-bucket-1234",
						}, nil),
						mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
					)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
					mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

					By("the object reference being invalid")
					mockBucket.EXPECT().SignedURL("test-key", gomock.Any()).Return("", fmt.Errorf("mock-error"))

					url, err := storagePlugin.PreSignUrl(context.TODO(), "test-bucket", "test-key", plugin.READ, 6000)

					By("Returning an error")
					Expect(err).Should(HaveOccurred())

					By("Returning a blank url")
					Expect(url).To(HaveLen(0))
				})
			})
		})

		When("The bucket doesn't exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
			mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
			storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

			It("Should return an error", func() {
				By("The bucket not existing")
				mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
				mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done)

				url, err := storagePlugin.PreSignUrl(context.TODO(), "test-bucket", "test-key", plugin.READ, 60)
				By("Returning an error")
				Expect(err).ShouldNot(BeNil())

				By("Returning an empty url")
				Expect(url).Should(BeEmpty())
			})
		})
	})

	Context("ListFiles", func() {
		When("The bucket exists", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
			mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
			mockObjectIterator := storage_mock.NewMockObjectIterator(ctrl)
			mockBucket := storage_mock.NewMockBucketHandle(ctrl)
			storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

			It("Should return the files on the bucket", func() {
				By("the bucket existing")
				gomock.InOrder(
					mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
						Labels: map[string]string{
							"x-nitric-name": "test-bucket",
						},
						Name: "my-bucket-1234",
					}, nil),
					mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
				)
				mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
				mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

				By("the bucket containing files")
				mockBucket.EXPECT().Objects(gomock.Any(), gomock.Any()).Return(mockObjectIterator)
				gomock.InOrder(
					mockObjectIterator.EXPECT().Next().Return(&storage.ObjectAttrs{
						Name: "test-file",
					}, nil),
					mockObjectIterator.EXPECT().Next().Return(nil, iterator.Done),
				)

				files, err := storagePlugin.ListFiles(context.TODO(), "test-bucket")

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning a single file")
				Expect(files).To(HaveLen(1))

				By("The file having the returned name")
				Expect(files[0].Key).To(Equal("test-file"))
			})
		})

		When("The bucket does not exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
			mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
			storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

			It("Should return an error", func() {
				By("the bucket not existing")
				mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done)
				mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)

				files, err := storagePlugin.ListFiles(context.TODO(), "test-bucket")

				By("returning nil files")
				Expect(files).To(BeNil())

				By("returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
