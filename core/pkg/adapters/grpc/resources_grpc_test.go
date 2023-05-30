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
