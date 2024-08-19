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

package topic

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/smithy-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/sfniface"
	"github.com/nitrictech/nitric/cloud/aws/ifaces/snsiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	"github.com/nitrictech/nitric/core/pkg/help"
	topicpb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
)

type SnsEventService struct {
	client    snsiface.SNSAPI
	sfnClient sfniface.SFNAPI
	resolver  resource.AwsResourceResolver
}

var _ topicpb.TopicsServer = &SnsEventService{}

func isSNSAccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "SNS" && strings.Contains(opErr.Unwrap().Error(), "AuthorizationError")
	}
	return false
}

func (s *SnsEventService) getTopics(ctx context.Context) (map[string]resource.ResolvedResource, error) {
	return s.resolver.GetResources(ctx, resource.AwsResource_Topic)
}

func (s *SnsEventService) getStateMachines(ctx context.Context) (map[string]resource.ResolvedResource, error) {
	return s.resolver.GetResources(ctx, resource.AwsResource_StateMachine)
}

func (s *SnsEventService) publish(ctx context.Context, topic string, message string) error {
	topics, err := s.getTopics(ctx)
	if err != nil {
		return fmt.Errorf("error finding topics: %w", err)
	}

	snsTopic, ok := topics[topic]

	if !ok {
		return fmt.Errorf("could not resolve topic ARN from topic name")
	}

	mc := propagation.MapCarrier{}
	xray.Propagator{}.Inject(ctx, mc)

	attrs := map[string]types.MessageAttributeValue{}
	for k, v := range mc {
		attrs[k] = types.MessageAttributeValue{DataType: aws.String("String"), StringValue: aws.String(v)}
	}

	publishInput := &sns.PublishInput{
		TopicArn: aws.String(snsTopic.ARN),
		Message:  &message,
		// MessageStructure: json is for an AWS specific JSON format,
		// which sends different messages to different subscription types. Don't use it.
		// MessageStructure: aws.String("json"),
		MessageAttributes: attrs,
	}

	_, err = s.client.Publish(ctx, publishInput)

	return err
}

func (s *SnsEventService) publishDelayed(ctx context.Context, topic string, delay time.Duration, message string) error {
	stepFunctions, err := s.getStateMachines(ctx)
	if err != nil {
		return fmt.Errorf("error getting state machines: %w", err)
	}

	stepFunction, ok := stepFunctions[topic]
	if !ok {
		return fmt.Errorf("error resolving state machine ARN for topic %s: %w", topic, err)
	}

	mc := propagation.MapCarrier{}
	xray.Propagator{}.Inject(ctx, mc)

	input, err := json.Marshal(map[string]interface{}{
		"seconds": int(delay / time.Second),
		"message": message,
	})
	if err != nil {
		return err
	}

	_, err = s.sfnClient.StartExecution(ctx, &sfn.StartExecutionInput{
		StateMachineArn: aws.String(stepFunction.ARN),
		TraceHeader:     aws.String(mc[xray.Propagator{}.Fields()[0]]),
		Input:           aws.String(string(input)),
	})

	return err
}

// Publish to a given topic
func (s *SnsEventService) Publish(ctx context.Context, req *topicpb.TopicPublishRequest) (*topicpb.TopicPublishResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SnsEventService.Publish")

	messageBytes, err := proto.Marshal(req.Message)
	if err != nil {
		return nil, newErr(
			codes.Unknown,
			fmt.Sprintf("unable to serialize message. %s", help.BugInNitricHelpText()),
			err,
		)
	}
	message := base64.StdEncoding.EncodeToString(messageBytes)

	if req.Delay != nil && req.Delay.AsDuration() > 0 {
		err = s.publishDelayed(ctx, req.TopicName, req.Delay.AsDuration(), message)
	} else {
		err = s.publish(ctx, req.TopicName, message)
	}

	if err != nil {
		if isSNSAccessDeniedErr(err) {
			return nil, newErr(
				codes.PermissionDenied,
				"unable to publish to topic, this may be due to a missing permissions request in your code.",
				err,
			)
		}

		return nil, newErr(codes.Internal, "error publishing message", err)
	}

	return &topicpb.TopicPublishResponse{}, nil
}

// Create new SNS event service plugin
func New(resolver resource.AwsResourceResolver) (*SnsEventService, error) {
	awsRegion := env.AWS_REGION.String()

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
		resolver:  resolver,
	}, nil
}

func NewWithClient(provider resource.AwsResourceResolver, client snsiface.SNSAPI, sfnClient sfniface.SFNAPI) (*SnsEventService, error) {
	return &SnsEventService{
		resolver:  provider,
		client:    client,
		sfnClient: sfnClient,
	}, nil
}
