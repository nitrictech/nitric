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
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/smithy-go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/sqsiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	queuepb "github.com/nitrictech/nitric/core/pkg/proto/queue/v1"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

const (
	// ErrCodeNoSuchTagSet - AWS API neglects to include a constant for this error code.
	ErrCodeNoSuchTagSet = "NoSuchTagSet"
	ErrCodeAccessDenied = "AccessDenied"
)

type SQSQueueService struct {
	provider resource.AwsResourceProvider
	client   sqsiface.SQSAPI
}

var _ queuepb.QueueServiceServer = &SQSQueueService{}

// Get the URL for a given queue name
func (s *SQSQueueService) getUrlForQueueName(ctx context.Context, queue string) (*string, error) {
	queues, err := s.provider.GetResources(ctx, resource.AwsResource_Queue)
	if err != nil {
		return nil, fmt.Errorf("error retrieving queue list")
	}

	queueArn, ok := queues[queue]

	if !ok {
		return nil, fmt.Errorf("queue %s does not exist", queue)
	}

	arnParts := strings.Split(queueArn, ":")
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

func (s *SQSQueueService) Send(ctx context.Context, req *queuepb.QueueSendRequestBatch) (*queuepb.QueueSendResponse, error) {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.SendBatch",
		map[string]interface{}{
			"queue":     req.QueueName,
			"tasks.len": len(req.Requests),
		},
	)

	if url, err := s.getUrlForQueueName(ctx, req.QueueName); err == nil {
		entries := make([]types.SendMessageBatchRequestEntry, 0)

		for _, sendTaskReq := range req.Requests {
			t := sendTaskReq
			if bytes, err := json.Marshal(t); err == nil {
				entries = append(entries, types.SendMessageBatchRequestEntry{
					Id:          &t.Id,
					MessageBody: aws.String(string(bytes)),
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
			failedTasks := make([]*queuepb.FailedSendRequest, 0, len(out.Failed))
			for _, failed := range out.Failed {
				for _, e := range req.Requests {
					if e.Id == *failed.Id {
						failedTasks = append(failedTasks, &queuepb.FailedSendRequest{
							Request: e,
							Message: *failed.Message,
						})
						// continue processing failed messages
						break
					}
				}
			}

			return &queuepb.QueueSendResponse{
				FailedRequests: failedTasks,
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

func (s *SQSQueueService) Receive(ctx context.Context, req *queuepb.QueueReceiveRequest) (*queuepb.QueueReceiveResponse, error) {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Receive",
		map[string]interface{}{
			"depth": req.Depth,
			"queue": req.QueueName,
		},
	)

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

		tasks := make([]*queuepb.ReceivedTask, 0, len(res.Messages))
		for _, m := range res.Messages {
			receivedTask := &queuepb.ReceivedTask{
				LeaseId: *m.ReceiptHandle,
			}
			err := json.Unmarshal([]byte(*m.Body), &receivedTask.Payload)
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"failed unmarshalling body",
					err,
				)
			}

			tasks = append(tasks, receivedTask)
		}

		return &queuepb.QueueReceiveResponse{
			Tasks: tasks,
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
func (s *SQSQueueService) Complete(ctx context.Context, req *queuepb.QueueCompleteRequest) (*queuepb.QueueCompleteResponse, error) {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Complete",
		map[string]interface{}{
			"queue":   req.QueueName,
			"leaseId": req.LeaseId,
		},
	)

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

		return &queuepb.QueueCompleteResponse{}, nil
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to find queue",
			err,
		)
	}
}

func New(provider resource.AwsResourceProvider) (queuepb.QueueServiceServer, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

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

func NewWithClient(provider resource.AwsResourceProvider, client sqsiface.SQSAPI) queuepb.QueueServiceServer {
	return &SQSQueueService{
		client:   client,
		provider: provider,
	}
}
