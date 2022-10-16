// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc/codes"

	pb "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/secret"
	"github.com/nitrictech/nitric/pkg/span"
)

// GRPC Interface for registered Nitric Secret Plugins
type SecretServer struct {
	pb.UnimplementedSecretServiceServer
	secretPlugin secret.SecretService
}

func (s *SecretServer) checkPluginRegistered() error {
	if s.secretPlugin == nil {
		return NewPluginNotRegisteredError("Secret")
	}

	return nil
}

func (s *SecretServer) Put(ctx context.Context, req *pb.SecretPutRequest) (*pb.SecretPutResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "SecretService.Put", err)
	}

	span := span.FromContext(ctx, "secret-"+req.GetSecret().GetName())

	span.SetAttributes(
		semconv.CodeFunctionKey.String("Secret.Put"),
		attribute.Key("faas.secret.name").String(req.GetSecret().GetName()),
		attribute.Key("faas.secret.operation").String("put"),
	)

	defer span.End()

	if r, err := s.secretPlugin.Put(&secret.Secret{
		Name: req.GetSecret().GetName(),
	}, req.GetValue()); err == nil {
		span.SetAttributes(attribute.Key("faas.secret.version").String(r.SecretVersion.Version))

		return &pb.SecretPutResponse{
			SecretVersion: &pb.SecretVersion{
				Secret: &pb.Secret{
					Name: r.SecretVersion.Secret.Name,
				},
				Version: r.SecretVersion.Version,
			},
		}, nil
	} else {
		span.RecordError(err)

		return nil, NewGrpcError("SecretService.Put", err)
	}
}

func (s *SecretServer) Access(ctx context.Context, req *pb.SecretAccessRequest) (*pb.SecretAccessResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "SecretService.Access", err)
	}

	span := span.FromContext(ctx, "secret-"+req.GetSecretVersion().GetSecret().GetName())

	span.SetAttributes(
		semconv.CodeFunctionKey.String("Secret.Access"),
		attribute.Key("faas.secret.name").String(req.GetSecretVersion().GetSecret().GetName()),
		attribute.Key("faas.secret.operation").String("access"),
		attribute.Key("faas.secret.version").String(req.GetSecretVersion().GetVersion()),
	)

	defer span.End()

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
		span.RecordError(err)

		return nil, NewGrpcError("SecretService.Access", err)
	}
}

func NewSecretServer(secretPlugin secret.SecretService) pb.SecretServiceServer {
	return &SecretServer{
		secretPlugin: secretPlugin,
	}
}
