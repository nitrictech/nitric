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
	"context"
	"net/http"

	eventgridmgmt "github.com/Azure/azure-sdk-for-go/services/eventgrid/mgmt/2020-06-01/eventgrid"
	"github.com/Azure/go-autorest/autorest"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_eventgrid "github.com/nitrictech/nitric/mocks/mock_event_grid"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	eventgrid_service "github.com/nitrictech/nitric/pkg/plugins/events/eventgrid"
)

var _ = Describe("Event Grid Plugin", func() {
	topicName := "Test"
	topicEndpoint := "https://Test.local1-test.eventgrid.azure.net/api/events"
	topicListResponsePage := eventgridmgmt.NewTopicsListResultPage(
		eventgridmgmt.TopicsListResult{
			Value: &[]eventgridmgmt.Topic{
				{
					Name: &topicName,
					TopicProperties: &eventgridmgmt.TopicProperties{
						Endpoint: &topicEndpoint,
					},
				},
			},
		},
		func(context.Context, eventgridmgmt.TopicsListResult) (eventgridmgmt.TopicsListResult, error) {
			return eventgridmgmt.TopicsListResult{}, nil
		},
	)
	When("Listing Available Topics", func() {
		When("There are no topics available", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			topicClient := mock_eventgrid.NewMockTopicsClientAPI(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(eventgridClient, topicClient)
			Expect(topicClient).NotTo(BeNil())
			topicClient.EXPECT().ListBySubscription(
				gomock.Any(),
				"",
				gomock.Any(),
			).Return(eventgridmgmt.TopicsListResultPage{}, nil).Times(1)

			It("Should return an empty list of topics", func() {
				topics, err := eventgridPlugin.ListTopics()
				Expect(err).To(BeNil())
				Expect(topics).To(BeEmpty())
			})
		})

		When("There are topics available", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			topicClient := mock_eventgrid.NewMockTopicsClientAPI(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(eventgridClient, topicClient)

			topicClient.EXPECT().ListBySubscription(
				gomock.Any(),
				"",
				gomock.Any(),
			).Return(topicListResponsePage, nil).Times(1)

			It("Should return all available topics", func() {
				topics, err := eventgridPlugin.ListTopics()
				Expect(err).To(BeNil())
				Expect(topics).To(ContainElement("Test"))
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
			topicClient := mock_eventgrid.NewMockTopicsClientAPI(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(eventgridClient, topicClient)

			topicClient.EXPECT().ListBySubscription(
				gomock.Any(),
				"",
				gomock.Any(),
			).Return(eventgridmgmt.TopicsListResultPage{}, nil).Times(1)

			It("should return an error", func() {
				err := eventgridPlugin.Publish("Test", event)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("publishing to a topic that is unauthorised", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			topicClient := mock_eventgrid.NewMockTopicsClientAPI(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(eventgridClient, topicClient)

			eventgridClient.EXPECT().PublishEvents(
				gomock.Any(),
				"Test.local1-test.eventgrid.azure.net",
				gomock.Any(),
			).Return(autorest.Response{
				&http.Response{
					StatusCode: 403,
				},
			}, nil).Times(1)
			topicClient.EXPECT().ListBySubscription(
				gomock.Any(),
				"",
				gomock.Any(),
			).Return(topicListResponsePage, nil).Times(1)

			It("should return an error", func() {
				err := eventgridPlugin.Publish("Test", event)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("To a topic that does exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			topicClient := mock_eventgrid.NewMockTopicsClientAPI(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(eventgridClient, topicClient)

			eventgridClient.EXPECT().PublishEvents(
				gomock.Any(),
				"Test.local1-test.eventgrid.azure.net",
				gomock.Any(),
			).Return(autorest.Response{
				&http.Response{
					StatusCode: 202,
				},
			}, nil).Times(1)
			topicClient.EXPECT().ListBySubscription(
				gomock.Any(),
				"",
				gomock.Any(),
			).Return(topicListResponsePage, nil).Times(1)

			It("should successfully publish the message", func() {
				err := eventgridPlugin.Publish("Test", event)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("Providing an empty topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			topicClient := mock_eventgrid.NewMockTopicsClientAPI(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(eventgridClient, topicClient)

			It("should return an error", func() {
				err := eventgridPlugin.Publish("", event)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("provided invalid topic"))
			})
		})

		When("Providing an empty event", func() {
			ctrl := gomock.NewController(GinkgoT())
			eventgridClient := mock_eventgrid.NewMockBaseClientAPI(ctrl)
			topicClient := mock_eventgrid.NewMockTopicsClientAPI(ctrl)
			eventgridPlugin, _ := eventgrid_service.NewWithClient(eventgridClient, topicClient)

			It("should return an error", func() {
				err := eventgridPlugin.Publish("Test", nil)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("provided invalid event"))
			})
		})
	})
})
