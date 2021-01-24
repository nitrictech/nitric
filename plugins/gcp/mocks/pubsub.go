package mocks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"google.golang.org/api/iterator"
)

type MockPubsubClient struct {
	ifaces.PubsubClient
	topics                []string
	PublishedMessages     map[string][]ifaces.Message
	publishedMessageCount int64
}

func NewMockPubsubClient(topics []string) *MockPubsubClient {
	return &MockPubsubClient{
		topics:                topics,
		PublishedMessages:     make(map[string][]ifaces.Message),
		publishedMessageCount: 0,
	}
}

func (s *MockPubsubClient) Topic(name string) ifaces.Topic {
	return &MockPubsubTopic{
		name: name,
		c:    s,
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
	ifaces.Topic
	c    *MockPubsubClient
	name string
}

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

func (s *MockPubsubTopic) ID() string {
	return s.name
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

// We likely will not need this for now
// as we aren't supplying a mechanism for using it
// type MockPubSubSubscriberServer struct {
// 	pb.UnimplementedSubscriberServer
// }
