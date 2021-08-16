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

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/pkg/plugins/storage"
	"google.golang.org/grpc/codes"
)

// GRPC Interface for registered Nitric Storage Plugins
type StorageServiceServer struct {
	pb.UnimplementedStorageServiceServer
	storagePlugin storage.StorageService
}

func (s *StorageServiceServer) checkPluginRegistered() error {
	if s.storagePlugin == nil {
		return NewPluginNotRegisteredError("Storage")
	}

	return nil
}

func (s *StorageServiceServer) Write(ctx context.Context, req *pb.StorageWriteRequest) (*pb.StorageWriteResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "StorageService.Write", err)
	}

	if err := s.storagePlugin.Write(req.GetBucketName(), req.GetKey(), req.GetBody()); err == nil {
		return &pb.StorageWriteResponse{}, nil
	} else {
		return nil, NewGrpcError("StorageService.Write", err)
	}
}

func (s *StorageServiceServer) Read(ctx context.Context, req *pb.StorageReadRequest) (*pb.StorageReadResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "StorageService.Read", err)
	}

	if object, err := s.storagePlugin.Read(req.GetBucketName(), req.GetKey()); err == nil {
		return &pb.StorageReadResponse{
			Body: object,
		}, nil
	} else {
		return nil, NewGrpcError("StorageService.Read", err)
	}
}

func (s *StorageServiceServer) Delete(ctx context.Context, req *pb.StorageDeleteRequest) (*pb.StorageDeleteResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "StorageService.Delete", err)
	}

	if err := s.storagePlugin.Delete(req.GetBucketName(), req.GetKey()); err == nil {
		return &pb.StorageDeleteResponse{}, nil
	} else {
		return nil, NewGrpcError("StorageService.Delete", err)
	}
}

func NewStorageServiceServer(storagePlugin storage.StorageService) pb.StorageServiceServer {
	return &StorageServiceServer{
		storagePlugin: storagePlugin,
	}
}
