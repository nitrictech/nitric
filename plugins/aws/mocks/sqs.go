package mocks

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/nitric-dev/membrane/plugins/sdk"
)

type MockSqsOptions struct {
	Queues []string
}

type Message struct {
	Id    string
	Event sdk.NitricEvent
}

type MockSqs struct {
	sqsiface.SQSAPI
	queues   []string
	messages map[string][]Message
}

func (s *MockSqs) ListQueues(in *sqs.ListQueuesInput) (*sqs.ListQueuesOutput, error) {
	queueUrls := make([]*string, 0)

	for _, queue := range s.queues {
		queueUrls = append(queueUrls, &queue)
	}

	return &sqs.ListQueuesOutput{
		QueueUrls: queueUrls,
	}, nil
}

func (s *MockSqs) SendMessageBatch(in *sqs.SendMessageBatchInput) (*sqs.SendMessageBatchOutput, error) {
	for _, q := range s.queues {
		if *in.QueueUrl == q {
			if s.messages[q] == nil {
				s.messages[q] = make([]Message, 0)
			}

			successfulMessages := make([]*sqs.SendMessageBatchResultEntry, 0)
			failedMessages := make([]*sqs.BatchResultErrorEntry, 0)
			for _, e := range in.Entries {
				var evt sdk.NitricEvent

				json.Unmarshal([]byte(*e.MessageBody), &evt)

				s.messages[q] = append(s.messages[q], Message{
					Id:    *e.Id,
					Event: evt,
				})

				successfulMessages = append(successfulMessages, &sqs.SendMessageBatchResultEntry{
					Id: e.Id,
				})
			}

			// TODO: Add a configurable failure mechanism here...
			return &sqs.SendMessageBatchOutput{
				Successful: successfulMessages,
				Failed:     failedMessages,
			}, nil
		}
	}

	return nil, fmt.Errorf("Queue: %s does not exist", *in.QueueUrl)
}

func NewMockSqs(opts *MockSqsOptions) *MockSqs {
	return &MockSqs{
		queues:   opts.Queues,
		messages: make(map[string][]Message),
	}
}
