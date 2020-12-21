package services

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1/eventing"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPC Interface for registered Nitric Eventing Plugins
type EventingServer struct {
	pb.UnimplementedEventingServer
	eventingPlugin sdk.EventingPlugin
}

func (s *EventingServer) checkPluginRegistered() (bool, error) {
	if s.eventingPlugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Eventing plugin not registered")
	}

	return true, nil
}

func (s *EventingServer) Publish(ctx context.Context, req *pb.PublishRequest) (*empty.Empty, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		event := &sdk.NitricEvent{
			RequestId:   req.GetEvent().GetRequestId(),
			PayloadType: req.GetEvent().GetPayloadType(),
			Payload:     req.GetEvent().GetPayload().AsMap(),
		}
		if err := s.eventingPlugin.Publish(req.GetTopicName(), event); err == nil {
			return &empty.Empty{}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *EventingServer) GetTopics(context.Context, *empty.Empty) (*pb.GetTopicsReply, error) {
	if ok, err := s.checkPluginRegistered(); ok {

		if res, err := s.eventingPlugin.GetTopics(); err == nil {
			return &pb.GetTopicsReply{
				Topics: res,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NewEventingServer(eventingPlugin sdk.EventingPlugin) {
	return &EventingServer{
		eventingPlugin: eventingPlugin,
	}
}
