package grpc

import (
	"context"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserServer - GRPC API for the nitric user plugin
type UserServer struct {
	pb.UnimplementedUserServer
	// TODO: Support multiple plugin registerations
	// Just need to settle on a way of addressing them on calls
	plugin sdk.UserService
}

func (s *UserServer) checkPluginRegistered() (bool, error) {
	if s.plugin == nil {
		return false, status.Errorf(codes.Unimplemented, "User auth plugin not registered")
	}

	return true, nil
}

// Create - Creates a new user
func (s *UserServer) Create(ctx context.Context, req *pb.UserCreateRequest) (*pb.UserCreateResponse, error) {
	if ok, err := s.checkPluginRegistered(); !ok {
		return nil, err
	}

	err := s.plugin.Create(req.GetTenant(), req.GetId(), req.GetEmail(), req.GetPassword())

	if err != nil {
		return nil, err
	}

	return &pb.UserCreateResponse{}, nil
}

// NewUserServer - Returns a new concrete instance of the GRCP implementation for the Nitric User plugin
func NewUserServer(plugin sdk.UserService) pb.UserServer {
	return &UserServer{
		plugin: plugin,
	}
}
