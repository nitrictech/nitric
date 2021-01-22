package mocks

import (
	"context"
	"fmt"
	"strconv"

	"github.com/nitric-dev/membrane/plugins/gcp/ifaces"
	"google.golang.org/api/iterator"
	pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockPubSubPublisherServer struct {
	pb.UnimplementedPublisherServer
	topics            []string
	publishedMessages map[string][]*pb.PubsubMessage
}

func (m *MockPubSubPublisherServer) SetTopics(topics []string) {
	m.topics = topics
	m.ClearMessages()
}

func (m *MockPubSubPublisherServer) ClearMessages() {
	m.publishedMessages = make(map[string][]*pb.PubsubMessage)
}

func (m *MockPubSubPublisherServer) GetMessages() map[string][]*pb.PubsubMessage {
	return m.publishedMessages
}

// Adds one or more messages to the topic. Returns `NOT_FOUND` if the topic
// does not exist.
func (m *MockPubSubPublisherServer) Publish(ctx context.Context, req *pb.PublishRequest) (*pb.PublishResponse, error) {
	topic := req.GetTopic()

	var discoveredTopic string = ""
	for _, t := range m.topics {
		if topic == t {
			discoveredTopic = t
			break
		}
	}

	if discoveredTopic == "" {
		return nil, status.Errorf(codes.NotFound, "Topic does not exist")
	}

	var messageIds = make([]string, 0)
	for _, message := range req.GetMessages() {
		if _, ok := m.publishedMessages[discoveredTopic]; !ok {
			m.publishedMessages[discoveredTopic] = make([]*pb.PubsubMessage, 0)
		}

		messageIds = append(messageIds, strconv.FormatInt(int64(len(m.publishedMessages[discoveredTopic])), 10))
		m.publishedMessages[discoveredTopic] = append(m.publishedMessages[discoveredTopic], message)
	}

	return &pb.PublishResponse{
		MessageIds: messageIds,
	}, nil
}

// Lists matching topics.
func (m *MockPubSubPublisherServer) ListTopics(ctx context.Context, req *pb.ListTopicsRequest) (*pb.ListTopicsResponse, error) {
	topics := make([]*pb.Topic, 0)

	for _, topic := range m.topics {
		topics = append(topics, &pb.Topic{
			Name: topic,
		})
	}

	// List the available topics
	return &pb.ListTopicsResponse{
		Topics: topics,
	}, nil
}

func NewPubsubPublisherServer(topics []string) *MockPubSubPublisherServer {
	return &MockPubSubPublisherServer{
		topics:            topics,
		publishedMessages: make(map[string][]*pb.PubsubMessage),
	}
}

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
		return &MockPubsubTopic{
			c:    s.c,
			name: s.c.topics[s.idx],
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
