package grpc

import (
	"context"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPC Interface for registered Nitric Storage Plugins
type QueueServer struct {
	pb.UnimplementedQueueServer
	plugin sdk.QueueService
}

func (s *QueueServer) checkPluginRegistered() (bool, error) {
	if s.plugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Queue plugin not registered")
	}

	return true, nil
}

func (s *QueueServer) Send(ctx context.Context, req *pb.QueueSendRequest) (*pb.QueueSendResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		task := req.GetTask()

		nitricTask := sdk.NitricTask{
			ID:          task.GetId(),
			PayloadType: task.GetPayloadType(),
			Payload:     task.GetPayload().AsMap(),
		}

		if err := s.plugin.Send(req.GetQueue(), nitricTask); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}

	// Success
	return &pb.QueueSendResponse{}, nil
}

func (s *QueueServer) SendBatch(ctx context.Context, req *pb.QueueSendBatchRequest) (*pb.QueueSendBatchResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// Translate tasks
		tasks := make([]sdk.NitricTask, len(req.GetTasks()))
		for i, task := range req.GetTasks() {
			tasks[i] = sdk.NitricTask{
				ID:          task.GetId(),
				PayloadType: task.GetPayloadType(),
				Payload:     task.GetPayload().AsMap(),
			}
		}

		if resp, err := s.plugin.SendBatch(req.GetQueue(), tasks); err == nil {
			failedTasks := make([]*pb.FailedTask, len(resp.FailedMessages))
			for i, fmsg := range resp.FailedMessages {
				st, _ := structpb.NewStruct(fmsg.Task.Payload)
				failedTasks[i] = &pb.FailedTask{
					Message: fmsg.Message,
					Task: &pb.NitricTask{
						Id:          fmsg.Task.ID,
						PayloadType: fmsg.Task.PayloadType,
						Payload:     st,
					},
				}
			}
			return &pb.QueueSendBatchResponse{
				FailedTasks: failedTasks,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *QueueServer) Receive(ctx context.Context, req *pb.QueueReceiveRequest) (*pb.QueueReceiveResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// Convert gRPC request to plugin params
		depth := uint32(req.GetDepth())
		popOptions := sdk.ReceiveOptions{
			QueueName: req.GetQueue(),
			Depth:     &depth,
		}

		// Perform the Queue Pop operation
		tasks, err := s.plugin.Receive(popOptions)
		if err != nil {
			return nil, err
		}

		// Convert the NitricTasks to the gRPC type
		grpcTasks := []*pb.NitricTask{}
		for _, task := range tasks {
			st, _ := structpb.NewStruct(task.Payload)
			grpcTasks = append(grpcTasks, &pb.NitricTask{
				Id:          task.ID,
				Payload:     st,
				LeaseId:     task.LeaseID,
				PayloadType: task.PayloadType,
			})
		}

		// Return the tasks
		res := pb.QueueReceiveResponse{
			Tasks: grpcTasks,
		}
		return &res, nil
	} else {
		return nil, err
	}
}

func (s *QueueServer) Complete(ctx context.Context, req *pb.QueueCompleteRequest) (*pb.QueueCompleteResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		// Convert gRPC request to plugin params
		queueName := req.GetQueue()
		leaseId := req.GetLeaseId()

		// Perform the Queue Complete operation
		err := s.plugin.Complete(queueName, leaseId)
		if err != nil {
			return nil, err
		}

		// Return a successful response
		return &pb.QueueCompleteResponse{}, nil
	} else {
		return nil, err
	}
}

func NewQueueServer(plugin sdk.QueueService) pb.QueueServer {
	return &QueueServer{
		plugin: plugin,
	}
}
