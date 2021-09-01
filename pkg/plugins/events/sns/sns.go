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

package sns_service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nitric-dev/membrane/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/events"
)

type SnsEventService struct {
	events.UnimplementedeventsPlugin
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
func (s *SnsEventService) Publish(topic string, event *events.NitricEvent) error {
	newErr := errors.ErrorsWithScope(
		"SnsEventService.Publish",
		fmt.Sprintf("topic=%s", topic),
	)

	data, err := json.Marshal(event)

	if err != nil {
		return newErr(
			codes.Internal,
			"error marshalling event payload",
			err,
		)
	}

	topicArn, err := s.getTopicArnFromName(&topic)

	if err != nil {
		return newErr(
			codes.NotFound,
			"could not find topic",
			err,
		)
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
		return newErr(
			codes.Internal,
			"unable to publish message",
			err,
		)
	}

	return nil
}

func (s *SnsEventService) ListTopics() ([]string, error) {
	newErr := errors.ErrorsWithScope("SnsEventService.ListTopics")

	topicsOutput, err := s.client.ListTopics(&sns.ListTopicsInput{})

	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error retrieving topics",
			err,
		)
	}

	var topics []string
	for _, t := range topicsOutput.Topics {
		// TODO: Extract topic name from ARN
		topics = append(topics, *t.TopicArn)
	}

	return topics, nil
}

// New - Create new SNS event service plugin
func New() (events.EventService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %v", sessionError)
	}

	snsClient := sns.New(sess)

	return &SnsEventService{
		client: snsClient,
	}, nil
}

func NewWithClient(client snsiface.SNSAPI) (events.EventService, error) {
	return &SnsEventService{
		client: client,
	}, nil
}
