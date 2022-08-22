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

package mock_pubsub

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"

	ifaces_pubsub "github.com/nitrictech/nitric/pkg/ifaces/pubsub"
)

type MockPubsubClient struct {
	ifaces_pubsub.PubsubClient
	topics                []string
	PublishedMessages     map[string][]ifaces_pubsub.Message
	publishedMessageCount int64
}

type MockPubsubMessage struct {
	Id        string
	AckId     string
	DataBytes []byte
}

const MockProjectID = "mock-project"

var _ ifaces_pubsub.Message = (*MockPubsubMessage)(nil)

func (m MockPubsubMessage) ID() string {
	return m.Id
}

func (m MockPubsubMessage) Data() []byte {
	return m.DataBytes
}

func (m MockPubsubMessage) Attributes() map[string]string {
	return make(map[string]string)
}

func (m MockPubsubMessage) PublishTime() time.Time {
	return time.Now()
}

func (m MockPubsubMessage) Ack() {
}

func (m MockPubsubMessage) Nack() {
}

type MockPubsubOptions struct {
	Topics   []string
	Messages map[string][]ifaces_pubsub.Message
}

func NewMockPubsubClient(opts MockPubsubOptions) *MockPubsubClient {
	if opts.Messages == nil {
		opts.Messages = make(map[string][]ifaces_pubsub.Message)
	}
	publishedMessageCount := 0
	for t := range opts.Messages {
		publishedMessageCount += len(opts.Messages[t])
	}
	return &MockPubsubClient{
		topics:                opts.Topics,
		PublishedMessages:     opts.Messages,
		publishedMessageCount: int64(publishedMessageCount),
	}
}

func (s *MockPubsubClient) Topic(name string) ifaces_pubsub.Topic {
	return &MockPubsubTopic{
		name: name,
		c:    s,
		// subscriptions: nil,
	}
}

func (s *MockPubsubClient) Topics(context.Context) ifaces_pubsub.TopicIterator {
	return &MockTopicIterator{
		c:   s,
		idx: 0,
	}
}

type MockTopicIterator struct {
	ifaces_pubsub.TopicIterator
	c   *MockPubsubClient
	idx int
}

func (s *MockTopicIterator) Next() (ifaces_pubsub.Topic, error) {
	if s.idx < len(s.c.topics) {
		s.idx++
		return &MockPubsubTopic{
			c:    s.c,
			name: s.c.topics[s.idx-1],
		}, nil
	}

	return nil, iterator.Done
}

type MockPubsubTopic struct {
	// ifaces.Topic
	c    *MockPubsubClient
	name string
	// subscriptions []ifaces.Subscription
}

var _ ifaces_pubsub.Topic = (*MockPubsubTopic)(nil)

func (s *MockPubsubTopic) Publish(ctx context.Context, msg ifaces_pubsub.Message) ifaces_pubsub.PublishResult {
	for _, t := range s.c.topics {
		if t == s.name {
			if s.c.PublishedMessages[t] == nil {
				s.c.PublishedMessages[t] = make([]ifaces_pubsub.Message, 0)
			}

			s.c.PublishedMessages[t] = append(s.c.PublishedMessages[t], msg)
			s.c.publishedMessageCount++

			return &MockPublishResult{
				id:  strconv.FormatInt(s.c.publishedMessageCount, 10),
				err: nil,
			}
		}
	}

	return &MockPublishResult{
		id:  "",
		err: fmt.Errorf("Topic %s does not exist", s.name),
	}
}

func (s *MockPubsubTopic) Exists(ctx context.Context) (bool, error) {
	// TODO: Add handle for a backend (non-NotFound error)

	for _, t := range s.c.topics {
		if t == s.name {
			return true, nil
		}
	}

	return false, nil
}

func (s *MockPubsubTopic) Subscriptions(ctx context.Context) ifaces_pubsub.SubscriptionIterator {
	return &MockSubscriptionIterator{
		Subscriptions: []ifaces_pubsub.Subscription{
			// Just setup a single subscription, which is mocked to behave like the default pull subscription.
			MockSubscription{topic: s},
		},
	}
}

func (s *MockPubsubTopic) ID() string {
	return s.name
}

func (s *MockPubsubTopic) String() string {
	return s.name
}

type MockSubscription struct {
	topic ifaces_pubsub.Topic
}

var _ ifaces_pubsub.Subscription = (*MockSubscription)(nil)

func (m MockSubscription) ID() string {
	return fmt.Sprintf("%s-nitricqueue", m.topic.ID())
}

// Return the full unique name for the subscription, in the format "projects/{project}/subscriptions/{subscription}"
func (m MockSubscription) String() string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", MockProjectID, m.ID())
}

type MockSubscriptionIterator struct {
	i             int
	Subscriptions []ifaces_pubsub.Subscription
}

func (s *MockSubscriptionIterator) Next() (ifaces_pubsub.Subscription, error) {
	if len(s.Subscriptions) > s.i {
		s.i++
		return s.Subscriptions[s.i-1], nil
	}
	return nil, iterator.Done
}

type MockPublishResult struct {
	id  string
	err error
}

func (s *MockPublishResult) Get(context.Context) (string, error) {
	if s.err != nil {
		return "", s.err
	}

	return s.id, nil
}

type MockBaseClient struct {
	Messages      map[string][]ifaces_pubsub.Message
	CompleteError error
}

func (m MockBaseClient) Close() error {
	// do nothing, no need to close a mock.
	return nil
}

func (m MockBaseClient) Acknowledge(ctx context.Context, req *pubsubpb.AcknowledgeRequest, opts ...gax.CallOption) error {
	if m.CompleteError != nil {
		return m.CompleteError
	}
	return nil
}

func (m MockBaseClient) Pull(ctx context.Context, req *pubsubpb.PullRequest, opts ...gax.CallOption) (*pubsubpb.PullResponse, error) {
	sub := req.Subscription

	for topicName := range m.Messages {
		topicSubName := fmt.Sprintf("projects/%s/subscriptions/%s-nitricqueue", MockProjectID, topicName)
		if topicSubName == sub {
			mockMessages := m.Messages[topicName]

			if mockMessages == nil || len(mockMessages) < 1 {
				return &pubsubpb.PullResponse{}, nil
			}

			var messages []*pubsubpb.ReceivedMessage

			for i, m := range mockMessages {
				// Only return up to the max number of messages requested.
				if int32(i) >= req.MaxMessages {
					break
				}
				messages = append(messages, &pubsubpb.ReceivedMessage{
					AckId: m.ID(),
					Message: &pubsubpb.PubsubMessage{
						MessageId:  m.ID(),
						Data:       m.Data(),
						Attributes: m.Attributes(),
					},
				})
				mockMessages[i] = nil
			}

			res := &pubsubpb.PullResponse{
				ReceivedMessages: messages,
			}

			return res, nil
		}
	}
	return nil, fmt.Errorf("queue not found")
}
