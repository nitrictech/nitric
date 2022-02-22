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
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_s3iface "github.com/nitrictech/nitric/mocks/s3"
	s3_service "github.com/nitrictech/nitric/pkg/plugins/storage/s3"
	mock_s3 "github.com/nitrictech/nitric/tests/mocks/s3"
)

var _ = Describe("S3", func() {
	When("Write", func() {
		When("Given the S3 backend is available", func() {
			When("Creating an object in an existing bucket", func() {
				testPayload := []byte("Test")
				storage := make(map[string]map[string][]byte)
				mockStorageClient := mock_s3.NewStorageClient([]*mock_s3.MockBucket{
					{
						Name: "my-bucket",
						Tags: map[string]string{
							"x-nitric-name": "my-bucket",
						},
					},
				}, &storage)

				storagePlugin, _ := s3_service.NewWithClient(mockStorageClient)
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
				mockStorageClient := mock_s3.NewStorageClient([]*mock_s3.MockBucket{}, &storage)
				storagePlugin, _ := s3_service.NewWithClient(mockStorageClient)
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
					mockStorageClient := mock_s3.NewStorageClient([]*mock_s3.MockBucket{
						{
							Name: "test-bucket",
							Tags: map[string]string{
								"x-nitric-name": "test-bucket",
							},
						},
					}, &storage)
					storagePlugin, _ := s3_service.NewWithClient(mockStorageClient)

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
					mockStorageClient := mock_s3.NewStorageClient([]*mock_s3.MockBucket{
						{
							Name: "test-bucket",
							Tags: map[string]string{
								"x-nitric-name": "test-bucket",
							},
						},
					}, &storage)
					storagePlugin, _ := s3_service.NewWithClient(mockStorageClient)

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
	When("PreSignUrl", func() {
		When("The bucket exists", func() {
			// Set up a mock bucket, with a single item
			storage := make(map[string]map[string][]byte)
			storage["test-bucket"] = make(map[string][]byte)
			crtl := gomock.NewController(GinkgoT())
			mockStorageClient := mock_s3iface.NewMockS3API(crtl)
			storagePlugin, _ := s3_service.NewWithClient(mockStorageClient)

			When("A URL is requested for a known operation", func() {
				It("Should successfully generate the URL", func() {

					By("Calling ListBuckets to map the bucket name")
					mockStorageClient.EXPECT().ListBuckets(gomock.Any()).Times(1).Return(&s3.ListBucketsOutput{
						Buckets: []*s3.Bucket{{
							Name: aws.String("test-bucket-aaa111"),
						}},
					}, nil)

					mockStorageClient.EXPECT().GetBucketTagging(gomock.Any()).Times(1).Return(&s3.GetBucketTaggingOutput{TagSet: []*s3.Tag{{
						Key:   aws.String("x-nitric-name"),
						Value: aws.String("test-bucket"),
					}}}, nil)

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
		When("The bucket doesn't exist", func() {

		})
	})
})
