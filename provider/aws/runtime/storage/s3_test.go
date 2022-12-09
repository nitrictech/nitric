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
	"bytes"
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/provider/aws/mocks/provider"
	mock_s3iface "github.com/nitrictech/nitric/provider/aws/mocks/s3"
	"github.com/nitrictech/nitric/provider/aws/runtime/core"
	s3_service "github.com/nitrictech/nitric/provider/aws/runtime/storage"
)

var _ = Describe("S3", func() {
	When("Write", func() {
		When("Given the S3 backend is available", func() {
			When("Creating an object in an existing bucket", func() {
				testPayload := []byte("Test")
				ctrl := gomock.NewController(GinkgoT())

				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

				It("Should successfully store the object", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{
						"my-bucket": "arn:aws:s3:::my-bucket",
					}, nil)

					By("writing the item")
					mockStorageClient.EXPECT().PutObject(gomock.Any(), gomock.Any()).Return(&s3.PutObjectOutput{}, nil)

					err := storagePlugin.Write(context.TODO(), "my-bucket", "test-item", testPayload)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("Creating an object in a non-existent bucket", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)
				It("Should fail to store the item", func() {
					By("the bucket not existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{}, nil)

					err := storagePlugin.Write(context.TODO(), "my-bucket", "test-item", []byte("Test"))
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
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
					mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
					mockProvider := mock_provider.NewMockAwsProvider(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

					It("Should successfully retrieve the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{
							"test-bucket": "arn:aws:s3:::test-bucket",
						}, nil)

						By("the object existing")
						mockStorageClient.EXPECT().GetObject(gomock.Any(), &s3.GetObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						}).Return(&s3.GetObjectOutput{
							Body: io.NopCloser(bytes.NewReader([]byte("Test"))),
						}, nil)

						object, err := storagePlugin.Read(context.TODO(), "test-bucket", "test-key")
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
					ctrl := gomock.NewController(GinkgoT())
					mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
					mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
					mockProvider := mock_provider.NewMockAwsProvider(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

					It("Should successfully delete the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{
							"test-bucket": "arn:aws:s3:::test-bucket",
						}, nil)

						By("successfully deleting the object")
						mockStorageClient.EXPECT().DeleteObject(gomock.Any(), gomock.Any()).Return(&s3.DeleteObjectOutput{}, nil)

						err := storagePlugin.Delete(context.TODO(), "test-bucket", "test-key")
						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())
					})
				})
				When("The item doesn't exist", func() {
				})
			})
			When("The bucket doesn't exist", func() {
			})
		})
	})
	When("PreSignUrl", func() {
		When("The bucket exists", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockProvider := mock_provider.NewMockAwsProvider(ctrl)
			mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
			mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
			storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

			When("A URL is requested for a known operation", func() {
				It("Should successfully generate the URL", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{
						"test-bucket": "arn:aws:s3:::test-bucket-aaa111",
					}, nil)

					mockPSStorageClient.EXPECT().PresignPutObject(gomock.Any(), &s3.PutObjectInput{
						Bucket: aws.String("test-bucket-aaa111"), // the real bucket name should be provided here, not the nitric name
						Key:    aws.String("test-key"),
					}, gomock.Any()).Times(1).Return(&v4.PresignedHTTPRequest{
						URL: "aws.example.com",
					}, nil)

					url, err := storagePlugin.PreSignUrl(context.TODO(), "test-bucket", "test-key", 1, uint32(60))
					Expect(err).ShouldNot(HaveOccurred())

					By("Return the correct url")
					// always blank - it's the best we can do without a real mock.
					Expect(url).To(Equal("aws.example.com"))
				})
			})
		})
	})

	When("ListFiles", func() {
		When("The bucket exists", func() {
			When("The s3 backend is available", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				mockPSStorageClient := mock_s3iface.NewMockPreSignAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient, mockPSStorageClient)

				It("should list the files contained in the bucket", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Bucket).Return(map[string]string{
						"test-bucket": "arn:aws:s3:::test-bucket-aaa111",
					}, nil)

					By("s3 returning files")
					mockStorageClient.EXPECT().ListObjectsV2(gomock.Any(), &s3.ListObjectsV2Input{
						Bucket: aws.String("test-bucket-aaa111"),
					}).Return(&s3.ListObjectsV2Output{
						Contents: []types.Object{{
							Key: aws.String("test"),
						}},
					}, nil)

					files, err := storagePlugin.ListFiles(context.TODO(), "test-bucket")

					By("not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("returning the file listing from s3")
					Expect(files).To(HaveLen(1))

					By("having the returned keys")
					Expect(files[0].Key).To(Equal("test"))
				})
			})
		})
	})
})
