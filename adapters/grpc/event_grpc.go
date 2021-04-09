package grpc

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPC Interface for registered Nitric Eventing Plugins
type EventServer struct {
	pb.UnimplementedEventServer
	eventPlugin sdk.EventService
}

func (s *EventServer) checkPluginRegistered() (bool, error) {
	if s.eventPlugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Event plugin not registered")
	}
	return true, nil
}

func (s *EventServer) Publish(ctx context.Context, req *pb.EventPublishRequest) (*pb.EventPublishResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// auto generate an ID if we did not receive one
		var ID = req.GetEvent().GetId()
		if ID == "" {
			ID = uuid.New().String()
		}

		event := &sdk.NitricEvent{
			ID:          ID,
			PayloadType: req.GetEvent().GetPayloadType(),
			Payload:     req.GetEvent().GetPayload().AsMap(),
		}
		if err := s.eventPlugin.Publish(req.GetTopic(), event); err == nil {
			return &pb.EventPublishResponse{
				Id: ID,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NewEventServer(eventingPlugin sdk.EventService) pb.EventServer {
	return &EventServer{
		eventPlugin: eventingPlugin,
	}
}

type TopicServer struct {
	pb.UnimplementedTopicServer
	eventPlugin sdk.EventService
}

func (s *TopicServer) checkPluginRegistered() (bool, error) {
	if s.eventPlugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Event plugin not registered")
	}

	return true, nil
}

func (s *TopicServer) List(context.Context, *pb.TopicListRequest) (*pb.TopicListResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {

		if res, err := s.eventPlugin.ListTopics(); err == nil {
			topics := make([]*pb.NitricTopic, len(res))
			for i, topicName := range res {
				topics[i] = &pb.NitricTopic{
					Name: topicName,
				}
			}

			return &pb.TopicListResponse{
				Topics: topics,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NewTopicServer(eventService sdk.EventService) pb.TopicServer {
	// The external topic/event interfaces are separate. Internally, they're fulfilled together,
	// so the event plugin is all that's needed for both the Event and Topic servers currently.
	return &TopicServer{
		eventPlugin: eventService,
	}
}
