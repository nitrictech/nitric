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

package sns_service_test

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	provider_mocks "github.com/nitrictech/nitric/mocks/provider"
	sns_mock "github.com/nitrictech/nitric/mocks/sns"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	sns_service "github.com/nitrictech/nitric/pkg/plugins/events/sns"
	"github.com/nitrictech/nitric/pkg/providers/aws/core"
)

var _ = Describe("Sns", func() {
	Context("Get Topics", func() {
		When("There are available topics", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock)

			It("Should return the available topics", func() {
				By("Topics being available")
				awsMock.EXPECT().GetResources(core.AwsResource_Topic).Return(map[string]string{
					"test": "arn:test",
				}, nil)

				topics, err := eventsClient.ListTopics()

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

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock)
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
				awsMock.EXPECT().GetResources(core.AwsResource_Topic).Return(map[string]string{
					"test": "arn:test",
				}, nil)

				By("Publishing the message to the topic")
				snsMock.EXPECT().Publish(&sns.PublishInput{
					TopicArn: aws.String("arn:test"),
					Message:  aws.String(stringData),
				})

				err := eventsClient.Publish("test", testEvent)

				Expect(err).To(BeNil())
			})
		})

		When("Publishing to a non-existent topic", func() {
			ctrl := gomock.NewController(GinkgoT())
			awsMock := provider_mocks.NewMockAwsProvider(ctrl)
			snsMock := sns_mock.NewMockSNSAPI(ctrl)

			eventsClient, _ := sns_service.NewWithClient(awsMock, snsMock)

			payload := map[string]interface{}{"Test": "test"}

			It("Should return an error", func() {
				By("Returning no topics")
				awsMock.EXPECT().GetResources(core.AwsResource_Topic).Return(map[string]string{}, nil)

				err := eventsClient.Publish("test", &events.NitricEvent{
					ID:          "testing",
					PayloadType: "Test Payload",
					Payload:     payload,
				})

				Expect(err.Error()).To(ContainSubstring("could not find topic"))
			})
		})
	})
})
