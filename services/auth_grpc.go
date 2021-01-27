package services

import (
	"context"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1/auth"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServer - GRPC API for the nitric auth plugin
type AuthServer struct {
	pb.UnimplementedAuthServer
	// TODO: Support multiple plugin registerations
	// Just need to settle on a way of addressing them on calls
	plugin sdk.AuthPlugin
}

func (s *AuthServer) checkPluginRegistered() (bool, error) {
	if s.plugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Auth plugin not registered")
	}

	return true, nil
}

// CreateUser - Creates a new user
func (s *AuthServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if ok, err := s.checkPluginRegistered(); !ok {
		return nil, err
	}

	err := s.plugin.CreateUser(req.GetTenant(), req.GetId(), req.GetEmail(), req.GetPassword())

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{}, nil
}

// NewAuthServer - Returns a new concrete instance of the GRCP implementation for the Nitric Auth plugin
func NewAuthServer(plugin sdk.AuthPlugin) pb.AuthServer {
	return &AuthServer{
		plugin: plugin,
	}
}
