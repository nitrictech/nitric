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

package events_test

import (
	"context"
	"os"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	mock_cloudtasks "github.com/nitrictech/nitric/cloud/gcp/mocks/cloudtasks"
	mock_core "github.com/nitrictech/nitric/cloud/gcp/mocks/provider"
	mock_pubsub "github.com/nitrictech/nitric/cloud/gcp/mocks/pubsub"
	pubsub_service "github.com/nitrictech/nitric/cloud/gcp/runtime/events"
	"github.com/nitrictech/nitric/core/pkg/plugins/events"
)

var _ = Describe("Pubsub Plugin", func() {
	_ = os.Setenv("NITRIC_STACK_ID", "test-stack")

	When("Listing Available Topics", func() {
		When("There are no topics available", func() {
			ctrl := gomock.NewController(GinkgoT())
			pubsubClient := mock_pubsub.NewMockPubsubClient(ctrl)
			mockIterator := mock_pubsub.NewMockTopicIterator(ctrl)
			pubsubPlugin, _ := pubsub_service.NewWithClient(nil, pubsubClient, nil)

			It("Should return an empty list of topics", func() {
				By("Returning an empty topic iterator")
				mockIterator.EXPECT().Next().Return(nil, iterator.Done)
				pubsubClient.EXPECT().Topics(gomock.Any()).Return(mockIterator)

				topics, err := pubsubPlugin.ListTopics(context.TODO())
				Expect(err).To(BeNil())
				Expect(topics).To(BeEmpty())
			})
		})

		When("There are topics available", func() {
			ctrl := gomock.NewController(GinkgoT())
			pubsubClient := mock_pubsub.NewMockPubsubClient(ctrl)
			mockIterator := mock_pubsub.NewMockTopicIterator(ctrl)
			mockTopic := mock_pubsub.NewMockTopic(ctrl)
			pubsubPlugin, _ := pubsub_service.NewWithClient(nil, pubsubClient, nil)

			It("Should return all available topics", func() {
				By("Returning an non-empty topic iterator")
				mockTopic.EXPECT().ID().Return("Test")
				gomock.InOrder(
					mockIterator.EXPECT().Next().Return(mockTopic, nil),
					mockIterator.EXPECT().Next().Return(nil, iterator.Done),
				)
				pubsubClient.EXPECT().Topics(gomock.Any()).Return(mockIterator)

				topics, err := pubsubPlugin.ListTopics(context.TODO())
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

		When("To a topic that does exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			pubsubClient := mock_pubsub.NewMockPubsubClient(ctrl)
			mockTopic := mock_pubsub.NewMockTopic(ctrl)
			mockIterator := mock_pubsub.NewMockTopicIterator(ctrl)
			mockPublishResult := mock_pubsub.NewMockPublishResult(ctrl)
			pubsubPlugin, _ := pubsub_service.NewWithClient(nil, pubsubClient, nil)

			It("should successfully publish the message", func() {
				By("the publish being successful")
				mockPublishResult.EXPECT().Get(gomock.Any()).Return("mock-server", nil)

				By("the topic existing")
				pubsubClient.EXPECT().Topics(gomock.Any()).Return(mockIterator)
				gomock.InOrder(
					mockIterator.EXPECT().Next().Return(mockTopic, nil),
					mockIterator.EXPECT().Next().Return(nil, iterator.Done),
				)
				mockTopic.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(mockPublishResult)

				mockTopic.EXPECT().Labels(gomock.Any()).Return(map[string]string{
					"x-nitric-test-stack-name": "Test",
					"x-nitric-test-stack-type": "topic",
				}, nil)

				err := pubsubPlugin.Publish(context.TODO(), "Test", 0, event)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		When("There are insufficient permissions", func() {
			ctrl := gomock.NewController(GinkgoT())
			pubsubClient := mock_pubsub.NewMockPubsubClient(ctrl)
			mockTopic := mock_pubsub.NewMockTopic(ctrl)
			mockIterator := mock_pubsub.NewMockTopicIterator(ctrl)
			mockPublishResult := mock_pubsub.NewMockPublishResult(ctrl)
			pubsubPlugin, _ := pubsub_service.NewWithClient(nil, pubsubClient, nil)

			It("should return an error", func() {
				By("the topic existing")
				pubsubClient.EXPECT().Topics(gomock.Any()).Return(mockIterator)
				gomock.InOrder(
					mockIterator.EXPECT().Next().Return(mockTopic, nil),
					mockIterator.EXPECT().Next().Return(nil, iterator.Done),
				)

				mockPublishResult.EXPECT().Get(gomock.Any()).Return("", status.Error(codes.PermissionDenied, "insufficient permissions")).Times(1)
				mockTopic.EXPECT().Publish(gomock.Any(), gomock.Any()).Return(mockPublishResult).Times(1)

				mockTopic.EXPECT().Labels(gomock.Any()).Return(map[string]string{
					"x-nitric-test-stack-name": "Test",
					"x-nitric-test-stack-type": "topic",
				}, nil)

				By("an insufficient permissions error is returned")
				err := pubsubPlugin.Publish(context.TODO(), "Test", 0, event)

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("permission denied, have you requested access to this topic?: \n rpc error: code = PermissionDenied desc = insufficient permissions"))
			})
		})
	})

	When("Publishing Delayed Messages", func() {
		event := &events.NitricEvent{
			ID:          "Test",
			PayloadType: "Test",
			Payload: map[string]interface{}{
				"Test": "Test",
			},
		}

		When("To a topic that does exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			pubsubClient := mock_pubsub.NewMockPubsubClient(ctrl)
			cloudtasksClient := mock_cloudtasks.NewMockCloudtasksClient(ctrl)
			mockGcp := mock_core.NewMockGcpProvider(ctrl)
			mockTopic := mock_pubsub.NewMockTopic(ctrl)
			mockIterator := mock_pubsub.NewMockTopicIterator(ctrl)
			pubsubPlugin, _ := pubsub_service.NewWithClient(mockGcp, pubsubClient, cloudtasksClient)

			It("should successfully publish the message", func() {
				By("having a valid service account email")
				mockGcp.EXPECT().GetServiceAccountEmail().Return("test@test.com", nil)

				By("having a valid project id")
				mockGcp.EXPECT().GetProjectID().Return("mock-project-id", nil)

				By("the topic existing")
				pubsubClient.EXPECT().Topics(gomock.Any()).Return(mockIterator)

				gomock.InOrder(
					mockIterator.EXPECT().Next().Return(mockTopic, nil),
					mockIterator.EXPECT().Next().Return(nil, iterator.Done),
				)
				mockTopic.EXPECT().Labels(gomock.Any()).Return(map[string]string{
					"x-nitric-test-stack-name": "Test",
					"x-nitric-test-stack-type": "topic",
				}, nil)
				mockTopic.EXPECT().String().Return("test")

				By("the publish being successful")
				// TODO: We want to validate that create task is called with the correct parameters.
				// This will require a custom gomock matcher, implemented here...
				cloudtasksClient.EXPECT().CreateTask(gomock.Any(), gomock.Any()).Return(nil, nil)

				err := pubsubPlugin.Publish(context.TODO(), "Test", 1, event)
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
