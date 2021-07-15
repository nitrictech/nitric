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
	"github.com/nitric-dev/membrane/pkg/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPC Interface for registered Nitric Storage Plugins
type StorageServer struct {
	pb.UnimplementedStorageServer
	storagePlugin sdk.StorageService
}

// Checks that the storage server is registered and returns gRPC Unimplemented error if not.
func (s *StorageServer) checkPluginRegistered() (bool, error) {
	if s.storagePlugin == nil {
		return false, status.Errorf(codes.Unimplemented, "Storage plugin not registered")
	}

	return true, nil
}

func (s *StorageServer) Write(ctx context.Context, req *pb.StorageWriteRequest) (*pb.StorageWriteResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.storagePlugin.Write(req.GetBucketName(), req.GetKey(), req.GetBody()); err == nil {
			return &pb.StorageWriteResponse{}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *StorageServer) Read(ctx context.Context, req *pb.StorageReadRequest) (*pb.StorageReadResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if object, err := s.storagePlugin.Read(req.GetBucketName(), req.GetKey()); err == nil {
			return &pb.StorageReadResponse{
				Body: object,
			}, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (s *StorageServer) Delete(ctx context.Context, req *pb.StorageDeleteRequest) (*pb.StorageDeleteResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.storagePlugin.Delete(req.GetBucketName(), req.GetKey()); err == nil {
			return &pb.StorageDeleteResponse{}, nil
		} else {
			// TODO: handle specific error codes.
			return nil, err
		}
	} else {
		return nil, err
	}
}

func NewStorageServer(storagePlugin sdk.StorageService) pb.StorageServer {
	return &StorageServer{
		storagePlugin: storagePlugin,
	}
}
