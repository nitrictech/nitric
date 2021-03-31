package sns_service_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	plugin "github.com/nitric-dev/membrane/plugins/aws/eventing/sns"
	"github.com/nitric-dev/membrane/sdk"
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
			eventingClient, _ := plugin.NewWithClient(&MockSNSClient{
				availableTopics: []*sns.Topic{&sns.Topic{TopicArn: aws.String("test")}},
			})

			It("Should return the available topics", func() {
				topics, err := eventingClient.ListTopics()

				Expect(err).To(BeNil())
				Expect(topics).To(ContainElements("test"))
			})
		})
	})

	Context("Publish", func() {
		When("Publishing to an available topic", func() {
			eventingClient, _ := plugin.NewWithClient(&MockSNSClient{
				availableTopics: []*sns.Topic{&sns.Topic{TopicArn: aws.String("test")}},
			})
			payload := map[string]interface{}{"Test": "test"}

			It("Should publish without error", func() {
				err := eventingClient.Publish("test", &sdk.NitricEvent{
					ID:          "testing",
					PayloadType: "Test Payload",
					Payload:     payload,
				})

				Expect(err).To(BeNil())
			})
		})

		When("Publishing to a non-existent topic", func() {
			eventingClient, _ := plugin.NewWithClient(&MockSNSClient{
				availableTopics: []*sns.Topic{},
			})

			payload := map[string]interface{}{"Test": "test"}

			It("Should return an error", func() {
				err := eventingClient.Publish("test", &sdk.NitricEvent{
					ID:          "testing",
					PayloadType: "Test Payload",
					Payload:     payload,
				})

				Expect(err.Error()).To(ContainSubstring("Unable to find topic"))
			})
		})
	})
})
