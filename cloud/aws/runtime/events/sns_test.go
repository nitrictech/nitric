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
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	provider_mocks "github.com/nitrictech/nitric/cloud/aws/mocks/provider"
	sfn_mock "github.com/nitrictech/nitric/cloud/aws/mocks/sfn"
	sns_mock "github.com/nitrictech/nitric/cloud/aws/mocks/sns"
	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	sns_service "github.com/nitrictech/nitric/cloud/aws/runtime/events"
	"github.com/nitrictech/nitric/core/pkg/plugins/events"
)

var _ = Describe("Sns", func() {
	Context("Get Topics", func() {
		When("There are available topics", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock, nil)

			It("Should return the available topics", func() {
				By("Topics being available")
				awsMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Topic).Return(map[string]string{
					"test": "arn:test",
				}, nil)

				topics, err := eventsClient.ListTopics(context.TODO())

				Expect(err).To(BeNil())
				Expect(topics).To(ContainElements("test"))

				ctrl.Finish()
			})
		})
	})

	Context("Publish", func() {
		When("Publishing to an available topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock, nil)
			payload := map[string]interface{}{"Test": "test"}
			testEvent := &events.NitricEvent{
				ID:          "testing",
				PayloadType: "Test Payload",
				Payload:     payload,
			}

			data, _ := json.Marshal(testEvent)
			stringData := string(data)

			It("Should publish without error", func() {
				By("Retrieving a list of topics")
				awsMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Topic).Return(map[string]string{
					"test": "arn:test",
				}, nil)

				By("Publishing the message to the topic")
				snsMock.EXPECT().Publish(gomock.Any(), &sns.PublishInput{
					MessageAttributes: map[string]types.MessageAttributeValue{},
					TopicArn:          aws.String("arn:test"),
					Message:           aws.String(stringData),
				})

				err := eventsClient.Publish(context.TODO(), "test", 0, testEvent)

				Expect(err).To(BeNil())
			})
		})

		When("Publishing to a non-existent topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock, nil)

			payload := map[string]interface{}{"Test": "test"}

			It("Should return an error", func() {
				By("Returning no topics")
				awsMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_Topic).Return(map[string]string{}, nil)

				err := eventsClient.Publish(context.TODO(), "test", 0, &events.NitricEvent{
					ID:          "testing",
					PayloadType: "Test Payload",
					Payload:     payload,
				})

				Expect(err.Error()).To(ContainSubstring("could not find topic"))
			})
		})
	})

	Context("Delayed Publish", func() {
		When("Publishing to an available topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsProvider(ctrl)
			sfnMock := sfn_mock.NewMockSFNAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, nil, sfnMock)
			payload := map[string]interface{}{"Test": "test"}
			testEvent := &events.NitricEvent{
				ID:          "testing",
				PayloadType: "Test Payload",
				Payload:     payload,
			}

			data, _ := json.Marshal(testEvent)
			stringData := string(data)

			It("Should publish without error", func() {
				By("Retrieving a list of topics")
				awsMock.EXPECT().GetResources(gomock.Any(), core.AwsResource_StateMachine).Return(map[string]string{
					"test": "arn:test",
				}, nil)

				By("Publishing the message to the topic")
				sfnMock.EXPECT().StartExecution(gomock.Any(), &sfn.StartExecutionInput{
					StateMachineArn: aws.String("arn:test"),
					TraceHeader:     aws.String(""),
					Input: aws.String(fmt.Sprintf(`{
			"seconds": 1,
			"message": %s
		}`, stringData)),
				})

				err := eventsClient.Publish(context.TODO(), "test", 1, testEvent)

				Expect(err).To(BeNil())
			})
		})
	})
})
