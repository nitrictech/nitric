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

package grpc_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_storage "github.com/nitrictech/nitric/mocks/storage"
	"github.com/nitrictech/nitric/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/storage"
)

var _ = Describe("GRPC Storage", func() {
	Context("Write", func() {
		When("plugin not registered", func() {
			ss := &grpc.StorageServiceServer{}
			resp, err := ss.Write(context.Background(), &v1.StorageWriteRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Storage plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)
			resp, err := grpc.NewStorageServiceServer(mockSS).Write(context.Background(), &v1.StorageWriteRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid StorageWriteRequest.BucketName"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)

			val := []byte("hush")
			mockSS.EXPECT().Write("bucky", "key", val)

			resp, err := grpc.NewStorageServiceServer(mockSS).Write(context.Background(), &v1.StorageWriteRequest{
				BucketName: "bucky",
				Key:        "key",
				Body:       val,
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(resp.String()).To(Equal(""))
			})
		})
	})

	Context("Read", func() {
		When("plugin not registered", func() {
			ss := &grpc.StorageServiceServer{}
			resp, err := ss.Read(context.Background(), &v1.StorageReadRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Storage plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)
			resp, err := grpc.NewStorageServiceServer(mockSS).Read(context.Background(), &v1.StorageReadRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid StorageReadRequest.BucketName"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)

			val := []byte("hush")
			mockSS.EXPECT().Read("bucky", "key").Return(val, nil)

			resp, err := grpc.NewStorageServiceServer(mockSS).Read(context.Background(), &v1.StorageReadRequest{
				BucketName: "bucky",
				Key:        "key",
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(string(resp.Body)).To(Equal("hush"))
			})
		})
	})

	Context("Delete", func() {
		When("plugin not registered", func() {
			ss := &grpc.StorageServiceServer{}
			resp, err := ss.Delete(context.Background(), &v1.StorageDeleteRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Storage plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)
			resp, err := grpc.NewStorageServiceServer(mockSS).Delete(context.Background(), &v1.StorageDeleteRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid StorageDeleteRequest.BucketName"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)

			mockSS.EXPECT().Delete("bucky", "key").Return(nil)

			_, err := grpc.NewStorageServiceServer(mockSS).Delete(context.Background(), &v1.StorageDeleteRequest{
				BucketName: "bucky",
				Key:        "key",
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("PreSignURL", func() {
		When("plugin not registered", func() {
			ss := &grpc.StorageServiceServer{}
			resp, err := ss.PreSignUrl(context.Background(), &v1.StoragePreSignUrlRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Storage plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid - bucket", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)
			resp, err := grpc.NewStorageServiceServer(mockSS).PreSignUrl(context.Background(), &v1.StoragePreSignUrlRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid StoragePreSignUrlRequest.BucketName"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid - key", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)
			resp, err := grpc.NewStorageServiceServer(mockSS).PreSignUrl(context.Background(), &v1.StoragePreSignUrlRequest{
				BucketName: "bucky",
			})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid StoragePreSignUrlRequest.Key: value length must be at least 1 runes"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)

			val := "example.com/this/and/that"
			mockSS.EXPECT().PreSignUrl("bucky", "key", storage.READ, uint32(0)).Return(val, nil)

			resp, err := grpc.NewStorageServiceServer(mockSS).PreSignUrl(context.Background(), &v1.StoragePreSignUrlRequest{
				BucketName: "bucky",
				Key:        "key",
				Operation:  v1.StoragePreSignUrlRequest_READ,
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(resp.Url).To(Equal(val))
			})
		})
	})

	Context("List", func() {
		When("plugin not registered", func() {
			ss := &grpc.StorageServiceServer{}
			resp, err := ss.ListFiles(context.Background(), &v1.StorageListFilesRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Storage plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)
			resp, err := grpc.NewStorageServiceServer(mockSS).ListFiles(context.Background(), &v1.StorageListFilesRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid StorageListFilesRequest.BucketName"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_storage.NewMockStorageService(g)

			mockSS.EXPECT().ListFiles("bucky").Return([]*storage.FileInfo{}, nil)

			_, err := grpc.NewStorageServiceServer(mockSS).ListFiles(context.Background(), &v1.StorageListFilesRequest{
				BucketName: "bucky",
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
			})
		})
	})
})
