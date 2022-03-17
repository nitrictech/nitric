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

package eventgrid_service_test

import (
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_eventgrid "github.com/nitrictech/nitric/mocks/mock_event_grid"
	mock_provider "github.com/nitrictech/nitric/mocks/provider"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	eventgrid_service "github.com/nitrictech/nitric/pkg/plugins/events/eventgrid"
	"github.com/nitrictech/nitric/pkg/providers/azure/core"
)

const mockRegion = "local1-test"
const mockTopicName = "Test"
const mockResourceName = "test-abcdef"

var getTopicResourcesResponse = map[string]core.AzGenericResource{
	mockTopicName: {
		Name:     mockResourceName,
		Location: mockRegion,
	},
}

var _ = Describe("Event Grid Plugin", func() {
	When("Listing Available Topics", func() {
		When("There are no topics available", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			mockProvider := mock_provider.NewMockAzProvider(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(mockProvider, eventgridClient)

			It("Should return an empty list of topics", func() {
				By("provider returning no available topics")
				mockProvider.EXPECT().GetResources(core.AzResource_Topic).Return(map[string]core.AzGenericResource{}, nil)

				topics, err := eventgridPlugin.ListTopics()
				Expect(err).To(BeNil())
				Expect(topics).To(BeEmpty())
				ctrl.Finish()
			})
		})

		When("There are topics available", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			mockProvider := mock_provider.NewMockAzProvider(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(mockProvider, eventgridClient)

			It("Should return all available topics", func() {
				By("provider returning a topic")
				mockProvider.EXPECT().GetResources(core.AzResource_Topic).Return(getTopicResourcesResponse, nil)

				topics, err := eventgridPlugin.ListTopics()
				Expect(err).To(BeNil())
				Expect(topics).To(ContainElement("Test"))

				ctrl.Finish()
			})
		})
	})

	When("Publishing Messages", func() {
		event := &events.NitricEvent{
			ID:          "Test",
			PayloadType: "Test",
			Payload: map[string]interface{}{
				"Test": "Test",
			},
		}

		When("To a topic that does not exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			mockProvider := mock_provider.NewMockAzProvider(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(mockProvider, eventgridClient)

			It("should return an error", func() {
				By("provider returning no topics")
				mockProvider.EXPECT().GetResources(core.AzResource_Topic).Return(map[string]core.AzGenericResource{}, nil)

				err := eventgridPlugin.Publish("Test", event)
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
					&http.Response{
						StatusCode: 403,
					},
				}, nil).Times(1)

				By("get resources returning topics")
				mockProvider.EXPECT().GetResources(core.AzResource_Topic).Return(getTopicResourcesResponse, nil)

				err := eventgridPlugin.Publish("Test", event)
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
				mockProvider.EXPECT().GetResources(core.AzResource_Topic).Return(getTopicResourcesResponse, nil)
				By("the eventgrid client publishing to the returned topic")
				eventgridClient.EXPECT().PublishEvents(
					gomock.Any(),
					fmt.Sprintf("%s.%s-1.eventgrid.azure.net", getTopicResourcesResponse["Test"].Name, getTopicResourcesResponse["Test"].Location),
					gomock.Any(),
				).Return(autorest.Response{
					&http.Response{
						StatusCode: 202,
					},
				}, nil).Times(1)

				err := eventgridPlugin.Publish("Test", event)
				Expect(err).ShouldNot(HaveOccurred())

				ctrl.Finish()
			})
		})
	})
})
