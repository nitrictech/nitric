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

package sqs_service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/plugins/queue"
	"github.com/nitric-dev/membrane/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

type SQSQueueService struct {
	queue.UnimplementedQueuePlugin
	client sqsiface.SQSAPI
}

// Get the URL for a given queue name
func (s *SQSQueueService) getUrlForQueueName(queueName string) (*string, error) {
	// TODO: Need to be able to guarantee same accound deployment in this case
	// In this case it would be preferred to use this method
	// s.client.GetQueueUrl(&sqs.GetQueueUrlInput{})
	if out, err := s.client.ListQueues(&sqs.ListQueuesInput{}); err == nil {
		for _, url := range out.QueueUrls {
			if strings.HasSuffix(*url, queueName) {
				return url, nil
			}
		}
	} else {
		return nil, fmt.Errorf("An Unexpected error occurred: %s", err)
	}

	return nil, fmt.Errorf("Could not find Queue: %s", queueName)
}

func (s *SQSQueueService) Send(queueName string, task queue.NitricTask) error {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Send",
		fmt.Sprintf("queue=%s", queueName),
	)

	tasks := []queue.NitricTask{task}
	if _, err := s.SendBatch(queueName, tasks); err != nil {
		return newErr(
			codes.Internal,
			"failed to send task",
			err,
		)
	}
	return nil
}

func (s *SQSQueueService) SendBatch(queueName string, tasks []queue.NitricTask) (*queue.SendBatchResponse, error) {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.SendBatch",
		fmt.Sprintf("queue=%s", queueName),
	)

	if url, err := s.getUrlForQueueName(queueName); err == nil {
		entries := make([]*sqs.SendMessageBatchRequestEntry, 0)

		for _, task := range tasks {
			if bytes, err := json.Marshal(task); err == nil {
				entries = append(entries, &sqs.SendMessageBatchRequestEntry{
					// Share the request ID here...
					Id:          &task.ID,
					MessageBody: aws.String(string(bytes)),
				})
			} else {
				// TODO: Do we want to just mark this one as having errored?
				return nil, newErr(
					codes.Internal,
					"error marshalling task",
					err,
				)
			}
		}

		if out, err := s.client.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries:  entries,
			QueueUrl: url,
		}); err == nil {
			// process out Failed messages to return to the user...
			failedTasks := make([]*queue.FailedTask, 0)
			for _, failed := range out.Failed {
				for _, e := range tasks {
					if e.ID == *failed.Id {
						failedTasks = append(failedTasks, &queue.FailedTask{
							Task:    &e,
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

func (s *SQSQueueService) Receive(options queue.ReceiveOptions) ([]queue.NitricTask, error) {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Receive",
		fmt.Sprintf("options=%v", options),
	)

	if err := options.Validate(); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid receive options",
			err,
		)
	}

	if url, err := s.getUrlForQueueName(options.QueueName); err == nil {
		req := sqs.ReceiveMessageInput{
			MaxNumberOfMessages: aws.Int64(int64(*options.Depth)),
			MessageAttributeNames: []*string{
				aws.String(sqs.QueueAttributeNameAll),
			},
			QueueUrl: url,
			// TODO: Consider explicit timeout values
			//VisibilityTimeout:       nil,
			//WaitTimeSeconds:         nil,
		}

		res, err := s.client.ReceiveMessage(&req)
		if err != nil {
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
				// TODO: append error to error list and Nack the message.
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
func (s *SQSQueueService) Complete(q string, leaseId string) error {
	newErr := errors.ErrorsWithScope(
		"SQSQueueService.Complete",
		fmt.Sprintf("queue=%s", q),
	)

	if url, err := s.getUrlForQueueName(q); err == nil {
		req := sqs.DeleteMessageInput{
			QueueUrl:      url,
			ReceiptHandle: aws.String(leaseId),
		}

		if _, err := s.client.DeleteMessage(&req); err != nil {
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

func New() (queue.QueueService, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("Error creating new AWS session %v", sessionError)
	}

	client := sqs.New(sess)

	return &SQSQueueService{
		client: client,
	}, nil
}

func NewWithClient(client sqsiface.SQSAPI) queue.QueueService {
	return &SQSQueueService{
		client: client,
	}
}
