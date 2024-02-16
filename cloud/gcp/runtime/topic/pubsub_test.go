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
	"os"
	"time"

	tasks "cloud.google.com/go/cloudtasks/apiv2/cloudtaskspb"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	mock_cloudtasks "github.com/nitrictech/nitric/cloud/gcp/mocks/cloudtasks"
	mock_core "github.com/nitrictech/nitric/cloud/gcp/mocks/provider"
	mock_pubsub "github.com/nitrictech/nitric/cloud/gcp/mocks/pubsub"
	"github.com/nitrictech/nitric/cloud/gcp/runtime/env"
	pubsub_service "github.com/nitrictech/nitric/cloud/gcp/runtime/topic"
	topicpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"github.com/nitrictech/nitric/test"
)

var _ = Describe("Pubsub Plugin", func() {
	_ = os.Setenv("NITRIC_STACK_ID", "test-stack")

	When("Publishing Messages", func() {
		payload, _ := structpb.NewStruct(map[string]interface{}{"Test": "test"})
		message := &topicpb.Message{
			Content: &topicpb.Message_StructPayload{
				StructPayload: payload,
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

				_, err := pubsubPlugin.Publish(context.TODO(), &topicpb.TopicPublishRequest{
					TopicName: "Test",
					Message:   message,
					Delay:     durationpb.New(0),
				})
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
				_, err := pubsubPlugin.Publish(context.TODO(), &topicpb.TopicPublishRequest{
					TopicName: "Test",
					Message:   message,
				})

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("PermissionDenied desc = PubsubEventService.Publish permission denied, have you requested access to this topic?"))
			})
		})
	})

	When("Publishing Delayed Messages", func() {
		payload, _ := structpb.NewStruct(map[string]interface{}{"Test": "test"})
		message := &topicpb.Message{
			Content: &topicpb.Message_StructPayload{
				StructPayload: payload,
			},
		}

		When("To a topic that does exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			pubsubClient := mock_pubsub.NewMockPubsubClient(ctrl)
			cloudtasksClient := mock_cloudtasks.NewMockCloudtasksClient(ctrl)
			mockGcp := mock_core.NewMockGcpResourceProvider(ctrl)
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

				matchBy := test.MatchBy(func(actual *tasks.CreateTaskRequest) (string, bool) {
					if actual.Parent != env.DELAY_QUEUE_NAME.String() {
						fmt.Printf("Parent didn't match\n")
						return fmt.Sprintf("parent to match %s got %s", env.DELAY_QUEUE_NAME.String(), actual.Parent), false
					}

					httpRequest := &tasks.HttpRequest{
						AuthorizationHeader: &tasks.HttpRequest_OauthToken{
							OauthToken: &tasks.OAuthToken{
								ServiceAccountEmail: "test@test.com",
							},
						},
						HttpMethod: tasks.HttpMethod_POST,
						Url:        fmt.Sprintf("https://pubsub.googleapis.com/v1/%s:publish", "test"),
						Body:       []byte("{\"messages\":[{\"attributes\":{\"x-cloud-trace-context\":\"00000000000000000000000000000000/0;o=0\",\"x-nitric-topic\":\"Test\"},\"data\":\"ChAKDgoEVGVzdBIGGgR0ZXN0\"}]}"),
					}

					if !proto.Equal(actual.Task.GetHttpRequest(), httpRequest) {
						return fmt.Sprintf("%+v\nActual: %+v", httpRequest, actual.Task.GetHttpRequest()), false
					}

					return "partial match is valid", true
				})

				cloudtasksClient.EXPECT().CreateTask(gomock.Any(), matchBy).Return(nil, nil)

				_, err := pubsubPlugin.Publish(context.TODO(), &topicpb.TopicPublishRequest{
					TopicName: "Test",
					Message:   message,
					Delay:     durationpb.New(60 * time.Second),
				})
				Expect(err).ShouldNot(HaveOccurred())
			})
		})
	})
})
