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

package topic_test

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_eventgrid "github.com/nitrictech/nitric/cloud/azure/mocks/mock_event_grid"
	mock_provider "github.com/nitrictech/nitric/cloud/azure/mocks/provider"
	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	eventgrid_service "github.com/nitrictech/nitric/cloud/azure/runtime/topic"
	topicpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

const (
	mockRegion       = "local1-test"
	mockTopicName    = "Test"
	mockResourceName = "test-abcdef"
)

var getTopicResourcesResponse = map[string]resource.AzGenericResource{
	mockTopicName: {
		Name:     mockResourceName,
		Location: mockRegion,
	},
}

var _ = Describe("Event Grid Plugin", func() {
	When("Publishing Messages", func() {
		eventPayload := &topicpb.Message{}

		When("To a topic that does not exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			mockProvider := mock_provider.NewMockAzProvider(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(mockProvider, eventgridClient)

			It("should return an error", func() {
				By("provider returning no topics")
				mockProvider.EXPECT().GetResources(gomock.Any(), resource.AzResource_Topic).Return(map[string]resource.AzGenericResource{}, nil)

				_, err := eventgridPlugin.Publish(context.TODO(), &topicpb.TopicPublishRequest{
					TopicName: "Test",
					Message:   eventPayload,
				})
				Expect(err).Should(HaveOccurred())

				ctrl.Finish()
			})
		})

		When("publishing to a topic that is unauthorised", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			mockProvider := mock_provider.NewMockAzProvider(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(mockProvider, eventgridClient)

			It("should return an error", func() {
				By("publish events returning an unauthorized error")
				eventgridClient.EXPECT().PublishEvents(
					gomock.Any(),
					fmt.Sprintf("%s.%s-1.eventgrid.azure.net", getTopicResourcesResponse["Test"].Name, getTopicResourcesResponse["Test"].Location),
					gomock.Any(),
				).Return(autorest.Response{
					Response: &http.Response{
						StatusCode: 403,
					},
				}, nil).Times(1)

				By("get resources returning topics")
				mockProvider.EXPECT().GetResources(gomock.Any(), resource.AzResource_Topic).Return(getTopicResourcesResponse, nil)

				_, err := eventgridPlugin.Publish(context.TODO(), &topicpb.TopicPublishRequest{
					TopicName: "Test",
					Message:   eventPayload,
				})
				Expect(err).Should(HaveOccurred())

				ctrl.Finish()
			})
		})

		When("To a topic that does exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			mockProvider := mock_provider.NewMockAzProvider(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(mockProvider, eventgridClient)

			It("should successfully publish the message", func() {
				By("the az provider returning topics")
				mockProvider.EXPECT().GetResources(gomock.Any(), resource.AzResource_Topic).Return(getTopicResourcesResponse, nil)
				By("the eventgrid client publishing to the returned topic")
				eventgridClient.EXPECT().PublishEvents(
					gomock.Any(),
					fmt.Sprintf("%s.%s-1.eventgrid.azure.net", getTopicResourcesResponse["Test"].Name, getTopicResourcesResponse["Test"].Location),
					gomock.Any(),
				).Return(autorest.Response{
					Response: &http.Response{
						StatusCode: 202,
					},
				}, nil).Times(1)

				_, err := eventgridPlugin.Publish(context.TODO(), &topicpb.TopicPublishRequest{
					TopicName: "Test",
					Message:   eventPayload,
				})
				Expect(err).ShouldNot(HaveOccurred())

				ctrl.Finish()
			})
		})
	})
})
