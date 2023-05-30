// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	mock_resource "github.com/nitrictech/nitric/core/mocks/resource"
	"github.com/nitrictech/nitric/core/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/plugins/resource"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GRPC Resources", func() {
	Context("Declare", func() {
		When("plugin not registered", func() {
			rs := &grpc.ResourcesServiceServer{}
			resp, err := rs.Declare(context.Background(), &v1.ResourceDeclareRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Resource plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockRS := mock_resource.NewMockResourceService(g)

			mockRS.EXPECT().Declare(gomock.Any(), &v1.ResourceDeclareRequest{}).Return(nil)

			_, err := grpc.NewResourcesServiceServer(grpc.WithResourcePlugin(mockRS)).Declare(context.Background(), &v1.ResourceDeclareRequest{})

			It("Should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})

	Context("Details", func() {
		When("plugin not registered", func() {
			rs := &grpc.ResourcesServiceServer{}
			resp, err := rs.Details(context.Background(), &v1.ResourceDetailsRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Resource plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is invalid", func() {
			g := gomock.NewController(GinkgoT())
			mockRS := mock_resource.NewMockResourceService(g)

			_, err := grpc.NewResourcesServiceServer(grpc.WithResourcePlugin(mockRS)).Details(context.Background(), &v1.ResourceDetailsRequest{
				Resource: &v1.Resource{
					Type: v1.ResourceType_Bucket,
				},
			})

			It("Should fail", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("request is valid", func() {
			defer GinkgoRecover()
			g := gomock.NewController(GinkgoT())
			mockRS := mock_resource.NewMockResourceService(g)

			mockRS.EXPECT().Details(gomock.Any(), resource.ResourceType_Api, "test").Return(&resource.DetailsResponse[any]{
				Id:       "test",
				Provider: "mock",
				Service:  "apigateway",
				Detail: resource.ApiDetails{
					URL: "http://mock.url/",
				},
			}, nil)

			_, err := grpc.NewResourcesServiceServer(grpc.WithResourcePlugin(mockRS)).Details(context.Background(), &v1.ResourceDetailsRequest{
				Resource: &v1.Resource{
					Type: v1.ResourceType_Api,
					Name: "test",
				},
			})

			It("Should succeed", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
