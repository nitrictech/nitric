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
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/smithy-go"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"

	provider_mocks "github.com/nitrictech/nitric/cloud/aws/mocks/provider"
	sfn_mock "github.com/nitrictech/nitric/cloud/aws/mocks/sfn"
	sns_mock "github.com/nitrictech/nitric/cloud/aws/mocks/sns"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	sns_service "github.com/nitrictech/nitric/cloud/aws/runtime/topic"
	eventpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

var _ = Describe("Sns", func() {
	Context("Publish", func() {
		When("Publishing to an available topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsResourceProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock, nil)
			payload, _ := structpb.NewStruct(map[string]interface{}{"Test": "test"})

			message := &eventpb.TopicMessage{
				Content: &eventpb.TopicMessage_StructPayload{
					StructPayload: payload,
				},
			}

			data, _ := proto.Marshal(message)
			stringData := string(data)

			It("Should publish without error", func() {
				By("Retrieving a list of topics")
				awsMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Topic).Return(map[string]resource.ResolvedResource{
					"test": {ARN: "arn:test"},
				}, nil)

				By("Publishing the message to the topic")
				snsMock.EXPECT().Publish(gomock.Any(), &sns.PublishInput{
					MessageAttributes: map[string]types.MessageAttributeValue{},
					TopicArn:          aws.String("arn:test"),
					Message:           aws.String(stringData),
				})

				_, err := eventsClient.Publish(context.TODO(), &eventpb.TopicPublishRequest{
					TopicName: "test",
					Delay:     nil,
					Message:   message,
				})

				Expect(err).To(BeNil())
			})
		})

		When("Publishing to a non-existent topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsResourceProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock, nil)

			payload, _ := structpb.NewStruct(map[string]interface{}{"Test": "test"})

			It("Should return an error", func() {
				By("Returning no topics")
				awsMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Topic).Return(map[string]resource.ResolvedResource{}, nil)

				_, err := eventsClient.Publish(context.TODO(), &eventpb.TopicPublishRequest{
					TopicName: "test",
					Message: &eventpb.TopicMessage{
						Content: &eventpb.TopicMessage_StructPayload{
							StructPayload: payload,
						},
					},
				})

				Expect(err.Error()).To(ContainSubstring("could not resolve topic ARN from topic name"))
			})
		})

		When("Publishing with access to the topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsResourceProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock, nil)

			payload, _ := structpb.NewStruct(map[string]interface{}{"Test": "test"})

			It("Should return an error", func() {
				By("Returning no topics")
				opErr := &smithy.OperationError{
					ServiceID: "SNS",
					Err:       errors.New("AuthorizationError"),
				}

				awsMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_Topic).Return(nil, opErr)

				_, err := eventsClient.Publish(context.TODO(), &eventpb.TopicPublishRequest{
					TopicName: "test",
					Delay:     durationpb.New(time.Second * 0),
					Message: &eventpb.TopicMessage{
						Content: &eventpb.TopicMessage_StructPayload{
							StructPayload: payload,
						},
					},
				})

				Expect(err.Error()).To(ContainSubstring("unable to publish to topic"))
			})
		})
	})

	Context("Delayed Publish", func() {
		When("Publishing to an available topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsResourceProvider(ctrl)
			sfnMock := sfn_mock.NewMockSFNAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, nil, sfnMock)
			payload, _ := structpb.NewStruct(map[string]interface{}{"Test": "test"})

			message := &eventpb.TopicMessage{
				Content: &eventpb.TopicMessage_StructPayload{
					StructPayload: payload,
				},
			}

			marshalledMessage, _ := proto.Marshal(message)
			stringData := string(marshalledMessage)

			It("Should publish without error", func() {
				By("Retrieving a list of topics")
				awsMock.EXPECT().GetResources(gomock.Any(), resource.AwsResource_StateMachine).Return(map[string]resource.ResolvedResource{
					"test": {ARN: "arn:test"},
				}, nil)

				input, _ := json.Marshal(map[string]interface{}{
					"seconds": 1,
					"message": stringData,
				})

				By("Publishing the message to the topic")
				sfnMock.EXPECT().StartExecution(gomock.Any(), &sfn.StartExecutionInput{
					StateMachineArn: aws.String("arn:test"),
					TraceHeader:     aws.String(""),
					Input:           aws.String(string(input)),
				})

				_, err := eventsClient.Publish(context.TODO(), &eventpb.TopicPublishRequest{
					TopicName: "test",
					Delay:     durationpb.New(time.Second * 1),
					Message: &eventpb.TopicMessage{
						Content: &eventpb.TopicMessage_StructPayload{
							StructPayload: payload,
						},
					},
				})

				Expect(err).To(BeNil())
			})
		})
	})
})
