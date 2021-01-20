package sqs_plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type SQSPlugin struct {
	sdk.UnimplementedQueuePlugin
	client sqsiface.SQSAPI
}

// Get the URL for a given queue name
func (s *SQSPlugin) getUrlForQueueName(queue string) (*string, error) {
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

func (s *SQSPlugin) Push(queue string, events []*sdk.NitricEvent) (*sdk.PushResponse, error) {
	if url, err := s.getUrlForQueueName(queue); err == nil {
		evts := make([]*sqs.SendMessageBatchRequestEntry, 0)

		for _, evt := range events {
			if bytes, err := json.Marshal(evt); err == nil {
				evts = append(evts, &sqs.SendMessageBatchRequestEntry{
					// Share the request ID here...
					Id:          &evt.RequestId,
					MessageBody: aws.String(string(bytes)),
				})
			} else {
				// TODO: Do we want to just mark this one as having errored?
				return nil, err
			}
		}

		// TODO: Get Succeeded/Failed Messages
		if out, err := s.client.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries:  evts,
			QueueUrl: url,
		}); err == nil {
			// process out Failed messages to return to the user...
			failedEvents := make([]*sdk.NitricEvent, 0)
			for _, failed := range out.Failed {
				for _, e := range events {
					if e.RequestId == *failed.Id {
						failedEvents = append(failedEvents, e)
						// continue outer loop
						break
					}
				}
			}

			return &sdk.PushResponse{
				FailedMessages: failedEvents,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func New() (sdk.QueuePlugin, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	sess, sessionError := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if sessionError != nil {
		return nil, fmt.Errorf("Error creating new AWS session %v", sessionError)
	}

	client := sqs.New(sess)

	return &SQSPlugin{
		client: client,
	}
}
