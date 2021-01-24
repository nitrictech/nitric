package services

import (
	"context"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1/queue"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPC Interface for registered Nitric Storage Plugins
type QueueServer struct {
	pb.UnimplementedQueueServer
	plugin sdk.QueuePlugin
}

func (s *QueueServer) checkPluginRegistered() (bool, error) {
	if s.plugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Queue plugin not registered")
	}

	return true, nil
}

func (s *QueueServer) Push(ctx context.Context, req *pb.PushRequest) (*pb.PushResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// Translate events
		evts := make([]*sdk.NitricEvent, len(req.GetEvents()))
		for i, evt := range req.GetEvents() {
			evts[i] = &sdk.NitricEvent{
				RequestId:   evt.GetRequestId(),
				PayloadType: evt.GetPayloadType(),
				Payload:     evt.GetPayload().AsMap(),
			}
		}

		if resp, err := s.plugin.Push(req.GetQueue(), evts); err == nil {
			failedMessages := make([]*pb.FailedMessage, len(resp.FailedMessages))
			for i, fmsg := range resp.FailedMessages {
				st, _ := structpb.NewStruct(fmsg.Event.Payload)
				failedMessages[i] = &pb.FailedMessage{
					Message: fmsg.Message,
					Event: &pb.NitricEvent{
						RequestId:   fmsg.Event.RequestId,
						PayloadType: fmsg.Event.PayloadType,
						Payload:     st,
					},
				}
			}
			return &pb.PushResponse{
				FailedMessages: failedMessages,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NewQueueServer(plugin sdk.QueuePlugin) pb.QueueServer {
	return &QueueServer{
		plugin: plugin,
	}
}
