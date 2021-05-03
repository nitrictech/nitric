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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type SQSQueueService struct {
	sdk.UnimplementedQueuePlugin
	client sqsiface.SQSAPI
}

// Get the URL for a given queue name
func (s *SQSQueueService) getUrlForQueueName(queue string) (*string, error) {
	// TODO: Need to be able to guarantee same accound deployment in this case
	// In this case it would be preferred to use this method
	// s.client.GetQueueUrl(&sqs.GetQueueUrlInput{})
	if out, err := s.client.ListQueues(&sqs.ListQueuesInput{}); err == nil {
		for _, url := range out.QueueUrls {
			if strings.HasSuffix(*url, queue) {
				return url, nil
			}
		}
	} else {
		return nil, fmt.Errorf("An Unexpected error occurred: %s", err)
	}

	return nil, fmt.Errorf("Could not find Queue: %s", queue)
}

func (s *SQSQueueService) Send(queue string, task sdk.NitricTask) error {
	tasks := []sdk.NitricTask{task}
	_, err := s.SendBatch(queue, tasks)
	return err
}

func (s *SQSQueueService) SendBatch(queue string, tasks []sdk.NitricTask) (*sdk.SendBatchResponse, error) {
	if url, err := s.getUrlForQueueName(queue); err == nil {
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
				return nil, err
			}
		}

		if out, err := s.client.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries:  entries,
			QueueUrl: url,
		}); err == nil {
			// process out Failed messages to return to the user...
			failedTasks := make([]*sdk.FailedTask, 0)
			for _, failed := range out.Failed {
				for _, e := range tasks {
					if e.ID == *failed.Id {
						failedTasks = append(failedTasks, &sdk.FailedTask{
							Task:    &e,
							Message: *failed.Message,
						})
						// continue processing failed messages
						break
					}
				}
			}

			return &sdk.SendBatchResponse{
				FailedTasks: failedTasks,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *SQSQueueService) Receive(options sdk.ReceiveOptions) ([]sdk.NitricTask, error) {
	err := options.Validate()
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("failed to retrieve messages: %s", err)
		}

		if len(res.Messages) == 0 {
			return []sdk.NitricTask{}, nil
		}

		var tasks []sdk.NitricTask
		for _, m := range res.Messages {
			var nitricTask sdk.NitricTask
			bodyBytes := []byte(*m.Body)
			err := json.Unmarshal(bodyBytes, &nitricTask)
			if err != nil {
				// TODO: append error to error list and Nack the message.
			}

			tasks = append(tasks, sdk.NitricTask{
				ID:          nitricTask.ID,
				Payload:     nitricTask.Payload,
				PayloadType: nitricTask.PayloadType,
				LeaseID:     *m.ReceiptHandle,
			})
		}

		return tasks, nil

	} else {
		return nil, err
	}
}

// Completes a previously popped queue item
func (s *SQSQueueService) Complete(queue string, leaseId string) error {
	if url, err := s.getUrlForQueueName(queue); err == nil {
		req := sqs.DeleteMessageInput{
			QueueUrl:      url,
			ReceiptHandle: aws.String(leaseId),
		}

		_, err := s.client.DeleteMessage(&req)
		if err != nil {
			return fmt.Errorf("failed to complete item: %s", err)
		}

		return nil

	} else {
		return err
	}
}

func New() (sdk.QueueService, error) {
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

func NewWithClient(client sqsiface.SQSAPI) sdk.QueueService {
	return &SQSQueueService{
		client: client,
	}
}
