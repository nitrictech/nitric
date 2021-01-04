package mocks

import (
	"context"
	"strconv"

	pb "google.golang.org/genproto/googleapis/pubsub/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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

// Creates the given topic with the given name. See the [resource name rules](
// https://cloud.google.com/pubsub/docs/admin#resource_names).
func (m *MockPubSubPublisherServer) CreateTopic(context.Context, *pb.Topic) (*pb.Topic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateTopic not implemented")
}

// Updates an existing topic. Note that certain properties of a
// topic are not modifiable.
func (m *MockPubSubPublisherServer) UpdateTopic(context.Context, *pb.UpdateTopicRequest) (*pb.Topic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateTopic not implemented")
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

// Gets the configuration of a topic.
func (m *MockPubSubPublisherServer) GetTopic(context.Context, *pb.GetTopicRequest) (*pb.Topic, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTopic not implemented")
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

// Lists the names of the attached subscriptions on this topic.
func (m *MockPubSubPublisherServer) ListTopicSubscriptions(context.Context, *pb.ListTopicSubscriptionsRequest) (*pb.ListTopicSubscriptionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTopicSubscriptions not implemented")
}

// Lists the names of the snapshots on this topic. Snapshots are used in
// [Seek](https://cloud.google.com/pubsub/docs/replay-overview) operations,
// which allow you to manage message acknowledgments in bulk. That is, you can
// set the acknowledgment state of messages in an existing subscription to the
// state captured by a snapshot.
func (m *MockPubSubPublisherServer) ListTopicSnapshots(context.Context, *pb.ListTopicSnapshotsRequest) (*pb.ListTopicSnapshotsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListTopicSnapshots not implemented")
}

// Deletes the topic with the given name. Returns `NOT_FOUND` if the topic
// does not exist. After a topic is deleted, a new topic may be created with
// the same name; this is an entirely new topic with none of the old
// configuration or subscriptions. Existing subscriptions to this topic are
// not deleted, but their `topic` field is set to `_deleted-topic_`.
func (m *MockPubSubPublisherServer) DeleteTopic(context.Context, *pb.DeleteTopicRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteTopic not implemented")
}

// Detaches a subscription from this topic. All messages retained in the
// subscription are dropped. Subsequent `Pull` and `StreamingPull` requests
// will return FAILED_PRECONDITION. If the subscription is a push
// subscription, pushes to the endpoint will stop.
func (m *MockPubSubPublisherServer) DetachSubscription(context.Context, *pb.DetachSubscriptionRequest) (*pb.DetachSubscriptionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DetachSubscription not implemented")
}

func NewPubsubPublisherServer(topics []string) *MockPubSubPublisherServer {
	return &MockPubSubPublisherServer{
		topics:            topics,
		publishedMessages: make(map[string][]*pb.PubsubMessage),
	}
}

// We likely will not need this for now
// as we aren't supplying a mechanism for using it
// type MockPubSubSubscriberServer struct {
// 	pb.UnimplementedSubscriberServer
// }
