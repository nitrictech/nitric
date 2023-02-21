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

package events

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/propagation"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/sfniface"
	"github.com/nitrictech/nitric/cloud/aws/ifaces/snsiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/core/pkg/plugins/events"
	utils2 "github.com/nitrictech/nitric/core/pkg/utils"
)

type SnsEventService struct {
	events.UnimplementedeventsPlugin
	client    snsiface.SNSAPI
	sfnClient sfniface.SFNAPI
	provider  core.AwsProvider
}

func (s *SnsEventService) getTopics(ctx context.Context) (map[string]string, error) {
	return s.provider.GetResources(ctx, core.AwsResource_Topic)
}

func (s *SnsEventService) getStateMachines(ctx context.Context) (map[string]string, error) {
	return s.provider.GetResources(ctx, core.AwsResource_StateMachine)
}

func (s *SnsEventService) publish(ctx context.Context, topic string, message string) error {
	topics, err := s.getTopics(ctx)
	if err != nil {
		return fmt.Errorf("error finding topics: %w", err)
	}

	topicArn, ok := topics[topic]

	if !ok {
		return fmt.Errorf("could not find topic")
	}

	mc := propagation.MapCarrier{}
	xray.Propagator{}.Inject(ctx, mc)

	attrs := map[string]types.MessageAttributeValue{}
	for k, v := range mc {
		attrs[k] = types.MessageAttributeValue{DataType: aws.String("String"), StringValue: &v}
	}

	publishInput := &sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  &message,
		// MessageStructure: json is for an AWS specific JSON format,
		// which sends different messages to different subscription types. Don't use it.
		// MessageStructure: aws.String("json"),
		MessageAttributes: attrs,
	}

	_, err = s.client.Publish(ctx, publishInput)

	if err != nil {
		return fmt.Errorf("unable to publish message: %w", err)
	}

	return nil
}

func (s *SnsEventService) publishDelayed(ctx context.Context, topic string, delay int, message string) error {
	sfns, err := s.getStateMachines(ctx)
	if err != nil {
		return fmt.Errorf("error gettings state machines: %w", err)
	}

	sfnArn, ok := sfns[topic]
	if !ok {
		return fmt.Errorf("error finding state machine for topic %s: %w", topic, err)
	}

	mc := propagation.MapCarrier{}
	xray.Propagator{}.Inject(ctx, mc)

	_, err = s.sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(sfnArn),
		TraceHeader:     aws.String(mc[xray.Propagator{}.Fields()[0]]),
		Input: aws.String(fmt.Sprintf(`{
			"seconds": %d,
			"message": %s
		}`, delay, message)),
	})
	if err != nil {
		return fmt.Errorf("error starting state machine execution: %w", err)
	}

	return nil
}

// Publish to a given topic
func (s *SnsEventService) Publish(ctx context.Context, topic string, delay int, event *events.NitricEvent) error {
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
		err = s.publishDelayed(ctx, topic, delay, message)
	} else {
		err = s.publish(ctx, topic, message)
	}

	if err != nil {
		return newErr(codes.Internal, "error publishing message", err)
	}

	return nil
}

func (s *SnsEventService) ListTopics(ctx context.Context) ([]string, error) {
	newErr := errors.ErrorsWithScope("SnsEventService.ListTopics", nil)

	topics, err := s.getTopics(ctx)
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

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	snsClient := sns.NewFromConfig(cfg)
	sfnClient := sfn.NewFromConfig(cfg)

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
