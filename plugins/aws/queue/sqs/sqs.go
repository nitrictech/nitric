package sqs_plugin

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/nitric-dev/membrane/plugins/sdk"
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

func (s *SQSPlugin) Push(queue string, events []*sdk.NitricEvent) error {
	if url, err := s.getUrlForQueueName(queue); err == nil {
		evts := make([]*sqs.SendMessageBatchRequestEntry, 0)
	
		for _, evt := range events {
			if bytes, err := json.Marshal(evt); err == nil {
				evts = append(evts, &sqs.SendMessageBatchRequestEntry{
					// Share the request ID here...
					Id: &evt.RequestId,
					MessageBody: aws.String(string(bytes)),
				})
			} else {
				return err
			}			
		} 

		out, err := s.client.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries: evts,
			QueueUrl: url
		})
	} else {
		return err
	}
}

func New() (sdk.QueuePlugin, error) {
	client := sqs.New()

	return &SQSPlugin{
		client: client,
	}
}
