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

package storage_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/durationpb"

	storage_mock "github.com/nitrictech/nitric/cloud/gcp/mocks/gcp_storage"
	storage_service "github.com/nitrictech/nitric/cloud/gcp/runtime/storage"
	storagePb "github.com/nitrictech/nitric/core/pkg/proto/storage/v1"
)

var _ = Describe("Storage", func() {
	os.Setenv("NITRIC_STACK_ID", "test-stack")

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
								"x-nitric-test-stack-name": "my-bucket",
								"x-nitric-test-stack-type": "bucket",
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

					_, err := mockStorageServer.Write(context.TODO(), &storagePb.StorageWriteRequest{
						BucketName: "my-bucket",
						Key:        "test-file",
						Body:       testPayload,
					})

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

					_, err := mockStorageServer.Write(context.TODO(), &storagePb.StorageWriteRequest{
						BucketName: "my-bucket",
						Key:        "test-file",
						Body:       testPayload,
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					ctrl.Finish()
				})
			})

			When("There are insufficient permissions", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucket := storage_mock.NewMockBucketHandle(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				mockObject := storage_mock.NewMockObjectHandle(ctrl)
				mockWriter := storage_mock.NewMockWriter(ctrl)
				mockStorageServer, _ := storage_service.NewWithClient(mockStorageClient)

				testPayload := []byte("Test")

				It("Should fail to store the item", func() {
					By("The bucket existing")
					gomock.InOrder(
						mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
							Labels: map[string]string{
								"x-nitric-test-stack-name": "my-bucket",
								"x-nitric-test-stack-type": "bucket",
							},
							Name: "my-bucket-1234",
						}, nil),
						mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
					)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
					mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

					By("The object reference being correct")
					mockBucket.EXPECT().Object("test-file").Return(mockObject)

					mockObject.EXPECT().NewWriter(gomock.Any()).Return(mockWriter)
					mockWriter.EXPECT().Write(gomock.Any()).Return(0, &googleapi.Error{Code: 403, Message: "insufficient permissions"})

					_, err := mockStorageServer.Write(context.TODO(), &storagePb.StorageWriteRequest{
						BucketName: "my-bucket",
						Key:        "test-file",
						Body:       testPayload,
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring("googleapi: Error 403: insufficient permissions"))
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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
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

						resp, err := storagePlugin.Read(context.TODO(), &storagePb.StorageReadRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning read content")
						Expect(resp.Body).To(HaveLen(0))
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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
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

						item, err := storagePlugin.Read(context.TODO(), &storagePb.StorageReadRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						By("Returning an error")
						Expect(err).Should(HaveOccurred())

						By("Not returning the item")
						Expect(item).To(BeNil())
					})
				})

				When("There are insufficient permissions", func() {
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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						mockBucket.EXPECT().Object("test-key").Return(mockObject)

						By("Returning a permission denied error")
						mockObject.EXPECT().NewReader(gomock.Any()).Return(nil, &googleapi.Error{Code: 403, Message: "insufficient permissions"})

						item, err := storagePlugin.Read(context.TODO(), &storagePb.StorageReadRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("googleapi: Error 403: insufficient permissions"))
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

					item, err := storagePlugin.Read(context.TODO(), &storagePb.StorageReadRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
					})

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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
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

						_, err := storagePlugin.Delete(context.TODO(), &storagePb.StorageDeleteRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
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

						_, err := storagePlugin.Delete(context.TODO(), &storagePb.StorageDeleteRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
					})
				})

				When("There are insufficient permissions", func() {
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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						mockBucket.EXPECT().Object("test-key").Return(mockObject)

						By("Returning a permission denied error")
						mockObject.EXPECT().Delete(gomock.Any()).Return(&googleapi.Error{Code: 403, Message: "insufficient permissions"})

						_, err := storagePlugin.Delete(context.TODO(), &storagePb.StorageDeleteRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
						})

						By("Returning an error")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).Should(ContainSubstring("googleapi: Error 403: insufficient permissions"))
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

					_, err := storagePlugin.Delete(context.TODO(), &storagePb.StorageDeleteRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
					})

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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being valid")
						mockBucket.EXPECT().SignedURL("test-key", gomock.Any()).Return("http://example.com", nil)

						resp, err := storagePlugin.PreSignUrl(context.TODO(), &storagePb.StoragePreSignUrlRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
							Expiry:     durationpb.New(time.Second * 60),
							Operation:  storagePb.StoragePreSignUrlRequest_READ,
						})

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("The url being returned")
						Expect(resp.Url).To(Equal("http://example.com"))

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
									"x-nitric-test-stack-name": "test-bucket",
									"x-nitric-test-stack-type": "bucket",
								},
								Name: "my-bucket-1234",
							}, nil),
							mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
						)
						mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
						mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

						By("the object reference being valid")
						mockBucket.EXPECT().SignedURL("test-key", gomock.Any()).Return("http://example.com", nil)

						response, err := storagePlugin.PreSignUrl(context.TODO(), &storagePb.StoragePreSignUrlRequest{
							BucketName: "test-bucket",
							Key:        "test-key",
							Operation:  storagePb.StoragePreSignUrlRequest_READ,
							Expiry:     durationpb.New(time.Second * 60),
						})

						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("The url being returned")
						Expect(response.Url).To(Equal("http://example.com"))

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
								"x-nitric-test-stack-name": "test-bucket",
								"x-nitric-test-stack-type": "bucket",
							},
							Name: "my-bucket-1234",
						}, nil),
						mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
					)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
					mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

					By("the object reference being invalid")
					mockBucket.EXPECT().SignedURL("test-key", gomock.Any()).Return("", fmt.Errorf("mock-error"))

					response, err := storagePlugin.PreSignUrl(context.TODO(), &storagePb.StoragePreSignUrlRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
						Operation:  storagePb.StoragePreSignUrlRequest_READ,
						Expiry:     durationpb.New(time.Second * 60),
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())

					By("Returning a blank url")
					Expect(response).To(BeNil())
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

				response, err := storagePlugin.PreSignUrl(context.TODO(), &storagePb.StoragePreSignUrlRequest{
					BucketName: "test-bucket",
					Key:        "test-key",
					Operation:  storagePb.StoragePreSignUrlRequest_READ,
					Expiry:     durationpb.New(time.Second * 60),
				})
				By("Returning an error")
				Expect(err).Should(HaveOccurred())

				By("Returning a nil response")
				Expect(response).Should(BeNil())
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
							"x-nitric-test-stack-name": "test-bucket",
							"x-nitric-test-stack-type": "bucket",
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
						Name: "test/test-file",
					}, nil),
					mockObjectIterator.EXPECT().Next().Return(nil, iterator.Done),
				)

				resp, err := storagePlugin.ListFiles(context.TODO(), &storagePb.StorageListFilesRequest{
					BucketName: "test-bucket",
					Prefix:     "test/",
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning a single file")
				Expect(resp.Files).To(HaveLen(1))

				By("The file having the returned name")
				Expect(resp.Files[0].Key).To(Equal("test/test-file"))
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

				resp, err := storagePlugin.ListFiles(context.TODO(), &storagePb.StorageListFilesRequest{
					BucketName: "test-bucket",
				})

				By("returning nil response")
				Expect(resp).To(BeNil())

				By("returning an error")
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Exists", func() {
		When("The bucket exists", func() {
			When("The item exists", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				mockBucket := storage_mock.NewMockBucketHandle(ctrl)
				mockObject := storage_mock.NewMockObjectHandle(ctrl)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("should return true", func() {
					By("the bucket existing")
					gomock.InOrder(
						mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
							Labels: map[string]string{
								"x-nitric-test-stack-name": "test-bucket",
								"x-nitric-test-stack-type": "bucket",
							},
							Name: "my-bucket-1234",
						}, nil),
						mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
					)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
					mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

					By("the object reference being valid")
					mockBucket.EXPECT().Object("test-key").Return(mockObject)

					By("returning object attributes")
					mockObject.EXPECT().Attrs(context.TODO()).Times(1).Return(&storage.ObjectAttrs{}, nil)

					response, err := storagePlugin.Exists(context.TODO(), &storagePb.StorageExistsRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("The url being returned")
					Expect(response.Exists).To(BeTrue())

					ctrl.Finish()
				})
			})

			When("The item doesn't exist", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				mockBucket := storage_mock.NewMockBucketHandle(ctrl)
				mockObject := storage_mock.NewMockObjectHandle(ctrl)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("should return true", func() {
					By("the bucket existing")
					gomock.InOrder(
						mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
							Labels: map[string]string{
								"x-nitric-test-stack-name": "test-bucket",
								"x-nitric-test-stack-type": "bucket",
							},
							Name: "my-bucket-1234",
						}, nil),
						mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
					)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
					mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

					By("the object reference being valid")
					mockBucket.EXPECT().Object("test-key").Return(mockObject)

					By("returning object attributes")
					mockObject.EXPECT().Attrs(context.TODO()).Times(1).Return(nil, storage.ErrObjectNotExist)

					response, err := storagePlugin.Exists(context.TODO(), &storagePb.StorageExistsRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
					})

					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("The url being returned")
					Expect(response.Exists).To(BeFalse())

					ctrl.Finish()
				})
			})

			When("google cloud returns an unknown error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := storage_mock.NewMockStorageClient(ctrl)
				mockBucketIterator := storage_mock.NewMockBucketIterator(ctrl)
				mockBucket := storage_mock.NewMockBucketHandle(ctrl)
				mockObject := storage_mock.NewMockObjectHandle(ctrl)
				storagePlugin, _ := storage_service.NewWithClient(mockStorageClient)

				It("should return an error", func() {
					By("the bucket existing")
					gomock.InOrder(
						mockBucketIterator.EXPECT().Next().Return(&storage.BucketAttrs{
							Labels: map[string]string{
								"x-nitric-test-stack-name": "test-bucket",
								"x-nitric-test-stack-type": "bucket",
							},
							Name: "my-bucket-1234",
						}, nil),
						mockBucketIterator.EXPECT().Next().Return(nil, iterator.Done),
					)
					mockStorageClient.EXPECT().Buckets(gomock.Any(), gomock.Any()).Return(mockBucketIterator)
					mockStorageClient.EXPECT().Bucket("my-bucket-1234").Return(mockBucket)

					By("the object reference being valid")
					mockBucket.EXPECT().Object("test-key").Return(mockObject)

					// Return the unknown error
					mockObject.EXPECT().Attrs(context.TODO()).Times(1).Return(nil, errors.New("unknown error"))

					response, err := storagePlugin.Exists(context.TODO(), &storagePb.StorageExistsRequest{
						BucketName: "test-bucket",
						Key:        "test-key",
					})

					By("Returning an error")
					Expect(err).Should(HaveOccurred())

					By("The response being nil")
					Expect(response).To(BeNil())

					ctrl.Finish()
				})
			})
		})
	})
})
