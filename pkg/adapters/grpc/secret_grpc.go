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
		if r, err := s.secretPlugin.Put(&secret.Secret{
			Name: req.GetSecret().GetName(),
		}, req.GetValue()); err == nil {
			return &pb.SecretPutResponse{
				SecretVersion: &pb.SecretVersion{
					Secret: &pb.Secret{
						Name: r.SecretVersion.Secret.Name,
					},
					Version: r.SecretVersion.Version,
				},
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *SecretServer) Access(ctx context.Context, req *pb.SecretAccessRequest) (*pb.SecretAccessResponse, error) {
	if err := s.checkPluginRegistered(); err == nil {
		if s, err := s.secretPlugin.Access(&secret.SecretVersion{
			Secret: &secret.Secret{
				Name: req.GetSecretVersion().GetSecret().GetName(),
			},
			Version: req.GetSecretVersion().GetVersion(),
		}); err == nil {
			return &pb.SecretAccessResponse{
				SecretVersion: &pb.SecretVersion{
					Secret: &pb.Secret{
						Name: s.SecretVersion.Secret.Name,
					},
					Version: s.SecretVersion.Version,
				},
				Value: s.Value,
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
