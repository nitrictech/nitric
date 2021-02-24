package sns_service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type SnsEventService struct {
	sdk.UnimplementedEventingPlugin
	client snsiface.SNSAPI
}

// Retrieve the topicArn for a given named nitric topic
func (s *SnsEventService) getTopicArnFromName(name *string) (*string, error) {
	topicsOutput, error := s.client.ListTopics(&sns.ListTopicsInput{})

	if error != nil {
		return nil, fmt.Errorf("There was an error retrieving SNS topics: %v", error)
	}

	for _, t := range topicsOutput.Topics {
		if strings.Contains(*t.TopicArn, *name) {
			return t.TopicArn, nil
		}
	}

	return nil, fmt.Errorf("Unable to find topic with name: %s", *name)
}

// Publish to a given topic
func (s *SnsEventService) Publish(topic string, event *sdk.NitricEvent) error {
	data, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("Payload marshalling error: %v", err)
	}

	topicArn, err := s.getTopicArnFromName(&topic)

	if err != nil {
		return fmt.Errorf("There was an error resolving the topic ARN for topic: %s, %v", topic, err)
	}

	message := string(data)

	publishInput := &sns.PublishInput{
		TopicArn: topicArn,
		Message:  &message,
		// MessageStructure: json is for an AWS specific JSON format,
		// which sends different messages to different subscription types. Don't use it.
		// MessageStructure: aws.String("json"),
	}

	_, err = s.client.Publish(publishInput)

	if err != nil {
		return fmt.Errorf("Error publishing message: %v", err)
	}

	return nil
}

func (s *SnsEventService) ListTopics() ([]string, error) {
	topicsOutput, error := s.client.ListTopics(&sns.ListTopicsInput{})

	if error != nil {
		return nil, fmt.Errorf("There was an error retrieving SNS topics: %v", error)
	}

	var topics []string
	for _, t := range topicsOutput.Topics {
		// TODO: Extract topic name from ARN
		topics = append(topics, *t.TopicArn)
	}

	return topics, nil
}

// Create new DynamoDB documents server
// XXX: No External Args for function atm (currently the plugin loader does not pass any argument information)
func New() (sdk.EventService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("Error creating new AWS session %v", sessionError)
	}

	snsClient := sns.New(sess)

	return &SnsEventService{
		client: snsClient,
	}, nil
}

func NewWithClient(client snsiface.SNSAPI) (sdk.EventService, error) {
	return &SnsEventService{
		client: client,
	}, nil
}
