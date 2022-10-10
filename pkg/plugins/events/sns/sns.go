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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
	"github.com/aws/aws-sdk-go/service/sfn/sfniface"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"

	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/events"
	"github.com/nitrictech/nitric/pkg/providers/aws/core"
	utils2 "github.com/nitrictech/nitric/pkg/utils"
)

type SnsEventService struct {
	events.UnimplementedeventsPlugin
	client    snsiface.SNSAPI
	sfnClient sfniface.SFNAPI
	provider  core.AwsProvider
}

func (s *SnsEventService) getTopics() (map[string]string, error) {
	return s.provider.GetResources(core.AwsResource_Topic)
}

func (s *SnsEventService) getStateMachines() (map[string]string, error) {
	return s.provider.GetResources(core.AwsResource_StateMachine)
}

func (s *SnsEventService) publish(topic string, message string) error {
	topics, err := s.getTopics()
	if err != nil {
		return fmt.Errorf("error finding topics: %v", err)
	}

	topicArn, ok := topics[topic]

	if !ok {
		return fmt.Errorf("could not find topic")
	}

	publishInput := &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  &message,
		// MessageStructure: json is for an AWS specific JSON format,
		// which sends different messages to different subscription types. Don't use it.
		// MessageStructure: aws.String("json"),
	}

	_, err = s.client.Publish(publishInput)

	if err != nil {
		return fmt.Errorf("unable to publish message: %v", err)
	}

	return nil
}

func (s *SnsEventService) publishDelayed(topic string, delay int, message string) error {
	sfns, err := s.getStateMachines()
	if err != nil {
		return fmt.Errorf("error finding state machine: %v", err)
	}

	sfnArn, ok := sfns[topic]
	if !ok {
		return fmt.Errorf("error finding state machine:: %v", err)
	}

	_, err = s.sfnClient.StartExecution(&sfn.StartExecutionInput{
		StateMachineArn: aws.String(sfnArn),
		Input: aws.String(fmt.Sprintf(`{
			"seconds": %d,
			"message": %s
		}`, delay, message)),
	})
	if err != nil {
		return fmt.Errorf("error starting state machine execution: %v", err)
	}

	return nil
}

// Publish to a given topic
func (s *SnsEventService) Publish(topic string, delay int, event *events.NitricEvent) error {
	newErr := errors.ErrorsWithScope(
		"SnsEventService.Publish",
		map[string]interface{}{
			"topic": topic,
			"event": event,
			"delay": delay,
		},
	)

	data, err := json.Marshal(event)
	if err != nil {
		return newErr(
			codes.Internal,
			"error marshalling event payload",
			err,
		)
	}
	message := string(data)

	if delay > 0 {
		err = s.publishDelayed(topic, delay, message)
	} else {
		err = s.publish(topic, message)
	}

	if err != nil {
		return newErr(codes.Internal, "error publishing message", err)
	}

	return nil
}

func (s *SnsEventService) ListTopics() ([]string, error) {
	newErr := errors.ErrorsWithScope("SnsEventService.ListTopics", nil)

	topics, err := s.getTopics()
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error retrieving topics",
			err,
		)
	}

	topicNames := make([]string, 0, len(topics))
	for name := range topics {
		// TODO: Extract topic name from ARN
		topicNames = append(topicNames, name)
	}

	return topicNames, nil
}

// Create new SNS event service plugin
func New(provider core.AwsProvider) (events.EventService, error) {
	awsRegion := utils2.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %v", sessionError)
	}

	snsClient := sns.New(sess)
	sfnClient := sfn.New(sess)

	return &SnsEventService{
		client:    snsClient,
		sfnClient: sfnClient,
		provider:  provider,
	}, nil
}

func NewWithClient(provider core.AwsProvider, client snsiface.SNSAPI, sfnClient sfniface.SFNAPI) (events.EventService, error) {
	return &SnsEventService{
		provider:  provider,
		client:    client,
		sfnClient: sfnClient,
	}, nil
}
