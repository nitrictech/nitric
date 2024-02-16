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

package queue

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/proto"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/sqsiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"

	queuespb "github.com/nitrictech/nitric/core/pkg/proto/queues/v1"
)

type SQSQueueService struct {
	provider resource.AwsResourceProvider
	client   sqsiface.SQSAPI
}

var _ queuespb.QueuesServer = &SQSQueueService{}

// Get the URL for a given queue name
func (s *SQSQueueService) getUrlForQueueName(ctx context.Context, queue string) (*string, error) {
	queues, err := s.provider.GetResources(ctx, resource.AwsResource_Queue)
	if err != nil {
		return nil, fmt.Errorf("error retrieving queue list: %w", err)
	}

	queueArn, ok := queues[queue]

	if !ok {
		return nil, fmt.Errorf("arn for queue %s could not be determined", queue)
	}

	arnParts := strings.Split(queueArn.ARN, ":")
	accountId := arnParts[4]
	queueName := arnParts[5]

	out, err := s.client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName:              aws.String(queueName),
		QueueOwnerAWSAccountId: aws.String(accountId),
	})
	if err != nil {
		return nil, fmt.Errorf("encountered an error retrieving the queue list: %w", err)
	}

	return out.QueueUrl, nil
}

func isSQSAccessDeniedErr(err error) bool {
	var opErr *smithy.OperationError
	if errors.As(err, &opErr) {
		return opErr.Service() == "SQS" && strings.Contains(opErr.Unwrap().Error(), "AccessDenied")
	}
	return false
}

func (s *SQSQueueService) Enqueue(ctx context.Context, req *queuespb.QueueEnqueueRequest) (*queuespb.QueueEnqueueResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SQSQueueService.Enqueue")

	requestIdMap := map[string]*queuespb.QueueMessage{}

	if url, err := s.getUrlForQueueName(ctx, req.QueueName); err == nil {
		entries := make([]types.SendMessageBatchRequestEntry, 0)

		for _, sendTaskReq := range req.Messages {
			t := sendTaskReq

			// generate a unique Id for each task
			id := uuid.New()
			requestIdMap[id.String()] = t

			if bytes, err := proto.Marshal(t); err == nil {
				msgString := base64.StdEncoding.EncodeToString(bytes)

				entries = append(entries, types.SendMessageBatchRequestEntry{
					Id:          aws.String(id.String()),
					MessageBody: aws.String(msgString),
				})
			} else {
				return nil, newErr(
					codes.Internal,
					"error marshalling task to JSON",
					err,
				)
			}
		}

		if out, err := s.client.SendMessageBatch(ctx, &sqs.SendMessageBatchInput{
			Entries:  entries,
			QueueUrl: url,
		}); err == nil {
			// process out Failed messages to return to the user...
			failedTasks := make([]*queuespb.FailedEnqueueMessage, 0, len(out.Failed))
			for _, failed := range out.Failed {
				for id, e := range requestIdMap {
					if id == *failed.Id {
						failedTasks = append(failedTasks, &queuespb.FailedEnqueueMessage{
							Message: e,
							Details: *failed.Message,
						})
						// continue processing failed messages
						break
					}
				}
			}

			return &queuespb.QueueEnqueueResponse{
				FailedMessages: failedTasks,
			}, nil
		} else {
			if isSQSAccessDeniedErr(err) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to send tasks to queue, have you requested access to this queue?",
					err,
				)
			}

			return nil, newErr(
				codes.Internal,
				"error sending tasks",
				err,
			)
		}
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to find queue",
			err,
		)
	}
}

func (s *SQSQueueService) Dequeue(ctx context.Context, req *queuespb.QueueDequeueRequest) (*queuespb.QueueDequeueResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SQSQueueService.Dequeue")

	if url, err := s.getUrlForQueueName(ctx, req.QueueName); err == nil {
		req := sqs.ReceiveMessageInput{
			MaxNumberOfMessages: req.Depth,
			MessageAttributeNames: []string{
				string(types.QueueAttributeNameAll),
			},
			QueueUrl: url,
			// VisibilityTimeout:       nil,
			// WaitTimeSeconds:         nil,
		}

		res, err := s.client.ReceiveMessage(ctx, &req)
		if err != nil {
			if isSQSAccessDeniedErr(err) {
				return nil, newErr(
					codes.PermissionDenied,
					"unable to receive task(s) from queue, have you requested access to this queue?",
					err,
				)
			}

			return nil, newErr(
				codes.Internal,
				"failed to retrieve message",
				err,
			)
		}

		tasks := make([]*queuespb.ReceivedMessage, 0, len(res.Messages))
		for _, m := range res.Messages {
			var queueMessage queuespb.QueueMessage

			msgBytes, err := base64.StdEncoding.DecodeString(*m.Body)
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"failed unmarshalling body",
					err,
				)
			}

			err = proto.Unmarshal(msgBytes, &queueMessage)
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"failed unmarshalling body",
					err,
				)
			}

			tasks = append(tasks, &queuespb.ReceivedMessage{
				LeaseId: *m.ReceiptHandle,
				Message: &queueMessage,
			})
		}

		return &queuespb.QueueDequeueResponse{
			Messages: tasks,
		}, nil
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to find queue",
			err,
		)
	}
}

// Completes a previously popped queue item
func (s *SQSQueueService) Complete(ctx context.Context, req *queuespb.QueueCompleteRequest) (*queuespb.QueueCompleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("SQSQueueService.Complete")

	if url, err := s.getUrlForQueueName(ctx, req.QueueName); err == nil {
		req := sqs.DeleteMessageInput{
			QueueUrl:      url,
			ReceiptHandle: aws.String(req.LeaseId),
		}

		if _, err := s.client.DeleteMessage(ctx, &req); err != nil {
			return nil, newErr(
				codes.Internal,
				"failed to dequeue task",
				err,
			)
		}

		return &queuespb.QueueCompleteResponse{}, nil
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to find queue",
			err,
		)
	}
}

func New(provider resource.AwsResourceProvider) (queuespb.QueuesServer, error) {
	awsRegion := env.AWS_REGION.String()

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	client := sqs.NewFromConfig(cfg)

	return &SQSQueueService{
		client:   client,
		provider: provider,
	}, nil
}

func NewWithClient(provider resource.AwsResourceProvider, client sqsiface.SQSAPI) queuespb.QueuesServer {
	return &SQSQueueService{
		client:   client,
		provider: provider,
	}
}
