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
	"github.com/aws/smithy-go"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/sqsiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/core/pkg/plugins/queue"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

const (
	// ErrCodeNoSuchTagSet - AWS API neglects to include a constant for this error code.
	ErrCodeNoSuchTagSet = "NoSuchTagSet"
	ErrCodeAccessDenied = "AccessDenied"
)

type SQSQueueService struct {
	queue.UnimplementedQueuePlugin
	provder core.AwsProvider
	client  sqsiface.SQSAPI
}

// Get the URL for a given queue name
func (s *SQSQueueService) getUrlForQueueName(ctx context.Context, queue string) (*string, error) {
	queues, err := s.provder.GetResources(ctx, core.AwsResource_Queue)
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

func (s *SQSQueueService) Send(ctx context.Context, queueName string, task queue.NitricTask) error {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Send",
		map[string]interface{}{
			"queue": queueName,
			"task":  task,
		},
	)

	tasks := []queue.NitricTask{task}
	if _, err := s.SendBatch(ctx, queueName, tasks); err != nil {
		if isSQSAccessDeniedErr(err) {
			return newErr(
				codes.PermissionDenied,
				"unable to send task to queue, have you requested access to this queue?",
				err,
			)
		}

		return newErr(
			codes.Internal,
			"failed to send task",
			err,
		)
	}
	return nil
}

func (s *SQSQueueService) SendBatch(ctx context.Context, queueName string, tasks []queue.NitricTask) (*queue.SendBatchResponse, error) {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.SendBatch",
		map[string]interface{}{
			"queue":     queueName,
			"tasks.len": len(tasks),
		},
	)

	if url, err := s.getUrlForQueueName(ctx, queueName); err == nil {
		entries := make([]types.SendMessageBatchRequestEntry, 0)

		for _, task := range tasks {
			t := task
			if bytes, err := json.Marshal(t); err == nil {
				entries = append(entries, types.SendMessageBatchRequestEntry{
					// Share the request ID here...
					Id:          &t.ID,
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
			failedTasks := make([]*queue.FailedTask, 0)
			for _, failed := range out.Failed {
				for _, e := range tasks {
					task := e
					if task.ID == *failed.Id {
						failedTasks = append(failedTasks, &queue.FailedTask{
							Task:    &task,
							Message: *failed.Message,
						})
						// continue processing failed messages
						break
					}
				}
			}

			return &queue.SendBatchResponse{
				FailedTasks: failedTasks,
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

func (s *SQSQueueService) Receive(ctx context.Context, options queue.ReceiveOptions) ([]queue.NitricTask, error) {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Receive",
		map[string]interface{}{
			"options": options,
		},
	)

	if err := options.Validate(); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid receive options",
			err,
		)
	}

	if url, err := s.getUrlForQueueName(ctx, options.QueueName); err == nil {
		req := sqs.ReceiveMessageInput{
			MaxNumberOfMessages: int32(*options.Depth),
			MessageAttributeNames: []string{
				string(types.QueueAttributeNameAll),
			},
			QueueUrl: url,
			// TODO: Consider explicit timeout values
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

		if len(res.Messages) == 0 {
			return []queue.NitricTask{}, nil
		}

		var tasks []queue.NitricTask
		for _, m := range res.Messages {
			var nitricTask queue.NitricTask
			bodyBytes := []byte(*m.Body)
			err := json.Unmarshal(bodyBytes, &nitricTask)
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"failed unmarshalling body",
					err,
				)
			}

			tasks = append(tasks, queue.NitricTask{
				ID:          nitricTask.ID,
				Payload:     nitricTask.Payload,
				PayloadType: nitricTask.PayloadType,
				LeaseID:     *m.ReceiptHandle,
			})
		}

		return tasks, nil
	} else {
		return nil, newErr(
			codes.NotFound,
			"unable to find queue",
			err,
		)
	}
}

// Completes a previously popped queue item
func (s *SQSQueueService) Complete(ctx context.Context, q string, leaseId string) error {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Complete",
		map[string]interface{}{
			"queue":   q,
			"leaseId": leaseId,
		},
	)

	if url, err := s.getUrlForQueueName(ctx, q); err == nil {
		req := sqs.DeleteMessageInput{
			QueueUrl:      url,
			ReceiptHandle: aws.String(leaseId),
		}

		if _, err := s.client.DeleteMessage(ctx, &req); err != nil {
			return newErr(
				codes.Internal,
				"failed to dequeue task",
				err,
			)
		}

		return nil
	} else {
		return newErr(
			codes.NotFound,
			"unable to find queue",
			err,
		)
	}
}

func New(provider core.AwsProvider) (queue.QueueService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	client := sqs.NewFromConfig(cfg)

	return &SQSQueueService{
		client:  client,
		provder: provider,
	}, nil
}

func NewWithClient(provider core.AwsProvider, client sqsiface.SQSAPI) queue.QueueService {
	return &SQSQueueService{
		client:  client,
		provder: provider,
	}
}
