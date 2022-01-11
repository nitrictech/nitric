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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/pkg/plugins/events"
	sns_service "github.com/nitrictech/nitric/pkg/plugins/events/sns"
)

type MockSNSClient struct {
	snsiface.SNSAPI
	// Available topics
	availableTopics []*sns.Topic
}

func (m *MockSNSClient) ListTopics(input *sns.ListTopicsInput) (*sns.ListTopicsOutput, error) {
	return &sns.ListTopicsOutput{
		Topics: m.availableTopics,
	}, nil
}

func (m *MockSNSClient) Publish(input *sns.PublishInput) (*sns.PublishOutput, error) {
	topicArn := input.TopicArn

	var topic *sns.Topic
	for _, t := range m.availableTopics {
		if topicArn == t.TopicArn {
			topic = t
			break
		}
	}

	if topic == nil {
		return nil, awserr.New(sns.ErrCodeNotFoundException, "Topic not found", fmt.Errorf("Topic does not exist"))
	}

	return &sns.PublishOutput{}, nil
}

var _ = Describe("Sns", func() {

	Context("Get Topics", func() {
		When("There are available topics", func() {
			eventsClient, _ := sns_service.NewWithClient(&MockSNSClient{
				availableTopics: []*sns.Topic{{TopicArn: aws.String("test")}},
			})

			It("Should return the available topics", func() {
				topics, err := eventsClient.ListTopics()

				Expect(err).To(BeNil())
				Expect(topics).To(ContainElements("test"))
			})
		})
	})

	Context("Publish", func() {
		When("Publishing to an available topic", func() {
			eventsClient, _ := sns_service.NewWithClient(&MockSNSClient{
				availableTopics: []*sns.Topic{{TopicArn: aws.String("test")}},
			})
			payload := map[string]interface{}{"Test": "test"}

			It("Should publish without error", func() {
				err := eventsClient.Publish("test", &events.NitricEvent{
					ID:          "testing",
					PayloadType: "Test Payload",
					Payload:     payload,
				})

				Expect(err).To(BeNil())
			})
		})

		When("Publishing to a non-existent topic", func() {
			eventsClient, _ := sns_service.NewWithClient(&MockSNSClient{
				availableTopics: []*sns.Topic{},
			})

			payload := map[string]interface{}{"Test": "test"}

			It("Should return an error", func() {
				err := eventsClient.Publish("test", &events.NitricEvent{
					ID:          "testing",
					PayloadType: "Test Payload",
					Payload:     payload,
				})

				Expect(err.Error()).To(ContainSubstring("Unable to find topic"))
			})
		})
	})
})
