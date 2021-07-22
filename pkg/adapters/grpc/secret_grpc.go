package grpc

import (
	"context"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPC Interface for registered Nitric Secret Plugins
type SecretServer struct {
	pb.UnimplementedSecretServiceServer
	secretPlugin secret.SecretService
}

func (s *SecretServer) checkPluginRegistered() error {
	if s.secretPlugin == nil {
		return status.Errorf(codes.Unimplemented, "Secret plugin not registered")
	}

	return nil
}

func (s *SecretServer) Put(ctx context.Context, req *pb.SecretPutRequest) (*pb.SecretPutResponse, error) {
	if err := s.checkPluginRegistered(); err == nil {
		if r, err := s.secretPlugin.Put(&secret.Secret{}); err == nil {
			return &pb.SecretPutResponse{
				Name:      r.Id,
				VersionId: r.VersionId,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *SecretServer) Get(ctx context.Context, req *pb.SecretGetRequest) (*pb.SecretGetResponse, error) {
	if err := s.checkPluginRegistered(); err == nil {
		if s, err := s.secretPlugin.Get(req.GetName(), req.GetVersionId()); err == nil {
			return &pb.SecretGetResponse{
				Secret: &pb.Secret{
					Name:  s.Name,
					Value: s.Value,
				},
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NewSecretServer(secretPlugin secret.SecretService) pb.SecretServiceServer {
	return &SecretServer{
		secretPlugin: secretPlugin,
	}
}
