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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/mocks/provider"
	mock_s3iface "github.com/nitrictech/nitric/mocks/s3"
	s3_service "github.com/nitrictech/nitric/pkg/plugins/storage/s3"
	"github.com/nitrictech/nitric/pkg/providers/aws/core"
)

var _ = Describe("S3", func() {
	When("Write", func() {
		When("Given the S3 backend is available", func() {
			When("Creating an object in an existing bucket", func() {
				testPayload := []byte("Test")
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_s3iface.NewMockS3API(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)

				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorage)
				It("Should successfully store the object", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(core.AwsResource_Bucket).Return(map[string]string{
						"my-bucket": "arn:aws:s3:::my-bucket",
					}, nil)

					By("writing the item")
					mockStorage.EXPECT().PutObject(gomock.Any()).Return(&s3.PutObjectOutput{}, nil)

					err := storagePlugin.Write("my-bucket", "test-item", testPayload)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("Creating an object in a non-existent bucket", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockStorage := mock_s3iface.NewMockS3API(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorage)
				It("Should fail to store the item", func() {
					By("the bucket not existing")
					mockProvider.EXPECT().GetResources(core.AwsResource_Bucket).Return(map[string]string{}, nil)

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
					ctrl := gomock.NewController(GinkgoT())
					mockStorage := mock_s3iface.NewMockS3API(ctrl)
					mockProvider := mock_provider.NewMockAwsProvider(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorage)

					It("Should successfully retrieve the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(core.AwsResource_Bucket).Return(map[string]string{
							"test-bucket": "arn:aws:s3:::test-bucket",
						}, nil)

						By("the object existing")
						mockStorage.EXPECT().GetObject(&s3.GetObjectInput{
							Bucket: aws.String("test-bucket"),
							Key:    aws.String("test-key"),
						}).Return(&s3.GetObjectOutput{
							Body: io.NopCloser(bytes.NewReader([]byte("Test"))),
						}, nil)

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
					ctrl := gomock.NewController(GinkgoT())
					mockStorage := mock_s3iface.NewMockS3API(ctrl)
					mockProvider := mock_provider.NewMockAwsProvider(ctrl)
					storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorage)

					It("Should successfully delete the object", func() {
						By("the bucket existing")
						mockProvider.EXPECT().GetResources(core.AwsResource_Bucket).Return(map[string]string{
							"test-bucket": "arn:aws:s3:::test-bucket",
						}, nil)

						By("successfully deleting the object")
						mockStorage.EXPECT().DeleteObject(gomock.Any()).Return(&s3.DeleteObjectOutput{}, nil)

						err := storagePlugin.Delete("test-bucket", "test-key")
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
			storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient)

			When("A URL is requested for a known operation", func() {
				It("Should successfully generate the URL", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(core.AwsResource_Bucket).Return(map[string]string{
						"test-bucket": "arn:aws:s3:::test-bucket-aaa111",
					}, nil)

					presign := 0
					mockStorageClient.EXPECT().PutObjectRequest(&s3.PutObjectInput{
						Bucket: aws.String("test-bucket-aaa111"), // the real bucket name should be provided here, not the nitric name
						Key:    aws.String("test-key"),
					}).Times(1).Return(&request.Request{
						Operation: &request.Operation{
							Name:       "",
							HTTPMethod: "",
							HTTPPath:   "",
							Paginator:  nil,
							// Unfortunately, PutObjectRequest returns a struct, instead of an interface,
							// so we can't really mock it. However, if this BeforePresignFn returns an error
							// it currently prevents the rest of the presign call and returns a blank url string.
							// this is good enough to perform basic testing.
							BeforePresignFn: func(r *request.Request) error {
								presign += 1
								return fmt.Errorf("test error")
							},
						},
						HTTPRequest: &http.Request{Host: "", URL: &url.URL{
							Scheme:      "",
							Opaque:      "",
							User:        nil,
							Host:        "aws.example.com",
							Path:        "",
							RawPath:     "",
							ForceQuery:  false,
							RawQuery:    "",
							Fragment:    "",
							RawFragment: "",
						}},
					}, nil)

					url, err := storagePlugin.PreSignUrl("test-bucket", "test-key", 1, uint32(60))
					By("Returning an error")
					// We always get an error due to inability to replace the Request with a mock
					Expect(err).Should(HaveOccurred())

					By("Returning a blank url")
					// always blank - it's the best we can do without a real mock.
					Expect(url).To(Equal(""))
				})
			})
		})
	})

	When("ListFiles", func() {
		When("The bucket exists", func() {
			When("The s3 backend is available", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				mockStorageClient := mock_s3iface.NewMockS3API(ctrl)
				storagePlugin, _ := s3_service.NewWithClient(mockProvider, mockStorageClient)

				It("should list the files contained in the bucket", func() {
					By("the bucket existing")
					mockProvider.EXPECT().GetResources(core.AwsResource_Bucket).Return(map[string]string{
						"test-bucket": "arn:aws:s3:::test-bucket-aaa111",
					}, nil)

					By("s3 returning files")
					mockStorageClient.EXPECT().ListObjects(&s3.ListObjectsInput{
						Bucket: aws.String("test-bucket-aaa111"),
					}).Return(&s3.ListObjectsOutput{
						Contents: []*s3.Object{{
							Key: aws.String("test"),
						}},
					}, nil)

					files, err := storagePlugin.ListFiles("test-bucket")

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
