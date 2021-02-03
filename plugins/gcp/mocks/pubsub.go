package mocks

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/googleapis/gax-go/v2"
	pubsubpb "google.golang.org/genproto/googleapis/pubsub/v1"

	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"google.golang.org/api/iterator"
)

type MockPubsubClient struct {
	ifaces.PubsubClient
	topics                []string
	PublishedMessages     map[string][]ifaces.Message
	publishedMessageCount int64
}

type MockPubsubMessage struct {
	Id        string
	AckId     string
	DataBytes []byte
}

const MockProjectID = "mock-project"

var _ ifaces.Message = (*MockPubsubMessage)(nil)

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
	Messages map[string][]ifaces.Message
}

func NewMockPubsubClient(opts MockPubsubOptions) *MockPubsubClient {
	if opts.Messages == nil {
		opts.Messages = make(map[string][]ifaces.Message)
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

func (s *MockPubsubClient) Topic(name string) ifaces.Topic {
	return &MockPubsubTopic{
		name: name,
		c:    s,
		//subscriptions: nil,
	}
}

func (s *MockPubsubClient) Topics(context.Context) ifaces.TopicIterator {
	return &MockTopicIterator{
		c:   s,
		idx: 0,
	}
}

type MockTopicIterator struct {
	ifaces.TopicIterator
	c   *MockPubsubClient
	idx int
}

func (s *MockTopicIterator) Next() (ifaces.Topic, error) {
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
	//ifaces.Topic
	c    *MockPubsubClient
	name string
	//subscriptions []ifaces.Subscription
}

var _ ifaces.Topic = (*MockPubsubTopic)(nil)

func (s *MockPubsubTopic) Publish(ctx context.Context, msg ifaces.Message) ifaces.PublishResult {
	for _, t := range s.c.topics {
		if t == s.name {
			if s.c.PublishedMessages[t] == nil {
				s.c.PublishedMessages[t] = make([]ifaces.Message, 0)
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

func (s *MockPubsubTopic) Subscriptions(ctx context.Context) ifaces.SubscriptionIterator {
	return &MockSubscriptionIterator{
		Subscriptions: []ifaces.Subscription{
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
	topic ifaces.Topic
}

var _ ifaces.Subscription = (*MockSubscription)(nil)

func (m MockSubscription) ID() string {
	fmt.Println("MOCK SUB ID METHOD")
	return fmt.Sprintf("%s-nitricqueue", m.topic.ID())
}

// Return the full unique name for the subscription, in the format "projects/{project}/subscriptions/{subscription}"
func (m MockSubscription) String() string {
	return fmt.Sprintf("projects/%s/subscriptions/%s", MockProjectID, m.ID())
}

type MockSubscriptionIterator struct {
	i             int
	Subscriptions []ifaces.Subscription
}

func (s *MockSubscriptionIterator) Next() (ifaces.Subscription, error) {
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
	Messages map[string][]ifaces.Message
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

	for topicName, _ := range m.Messages {
		topicSubName := fmt.Sprintf("projects/%s/subscriptions/%s-nitricqueue", MockProjectID, topicName)
		fmt.Println(fmt.Sprintf("Topic: %s, Sub: %s", topicSubName, sub))
		if topicSubName == sub {
			fmt.Println("FOUND QUEUE, ATTEMPTING TO POP MESSAGES")
			mockMessages := m.Messages[topicName]

			if mockMessages == nil || len(mockMessages) < 1 {
				fmt.Println("NO MESSAGES TO POP")
				return &pubsubpb.PullResponse{}, nil
			}
			fmt.Println(fmt.Sprintf("%d MESSAGES TO POP", len(mockMessages)))

			var messages []*pubsubpb.ReceivedMessage

			for i, m := range mockMessages {
				fmt.Println(fmt.Sprintf("RETURNING MESSAGE i:%d", i))
				// Only return up to the max number of messages requested.
				if int32(i) >= req.MaxMessages {
					fmt.Println(fmt.Sprintf("REACHED MAX MESSAGES i:%d, max:%d", i, req.MaxMessages))
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