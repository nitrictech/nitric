package deploy

import (
	"context"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeployServer struct {
	app *App
	pb.UnimplementedFaasServiceServer
	pb.UnimplementedResourceServiceServer
}

// Starts a new stream
// The deploy server will collect information from stream InitRequests and
// Immediately terminate the stream
func (s *DeployServer) TriggerStream(stream pb.FaasService_TriggerStreamServer) error {
	cm, err := stream.Recv()

	if err != nil {
		return status.Errorf(codes.Internal, "error reading message from stream: %v", err)
	}

	ir := cm.GetInitRequest()

	if ir == nil {
		// SHUT IT DOWN!!!!
		// The first message must be an init request from the prospective FaaS worker
		return status.Error(codes.FailedPrecondition, "first message must be InitRequest")
	}

	switch w := ir.Worker.(type) {
	case *pb.InitRequest_Api:
		s.app.AddApiHandler(w.Api)
	case *pb.InitRequest_Schedule:
		s.app.AddScheduleHandler(w.Schedule)
	case *pb.InitRequest_Subscription:
		s.app.AddSubscriptionHandler(w.Subscription)
	default:
		// treat as normal function worker
		// XXX: No-op for now. This can be handled exclusively at runtime
	}

	// Close the stream, once we've recieved the InitRequest
	return nil
}

// Declare - Accepts resource declarations and adds them to the Nitric App
func (s *DeployServer) Declare(ctx context.Context, req *pb.ResourceDeclareRequest) (*pb.ResourceDeclareResponse, error) {

	switch req.Resource.Type {
	case pb.ResourceType_Bucket:
		s.app.AddBucket(req.Resource.Name, req.GetBucket())
	case pb.ResourceType_Collection:
		s.app.AddCollection(req.Resource.Name, req.GetCollection())
	case pb.ResourceType_Queue:
		s.app.AddQueue(req.Resource.Name, req.GetQueue())
	case pb.ResourceType_Topic:
		s.app.AddTopic(req.Resource.Name, req.GetTopic())
	}

	return &pb.ResourceDeclareResponse{}, nil
}

// Create a new DeployServer
func New(app *App) *DeployServer {
	return &DeployServer{
		app: app,
	}
}
