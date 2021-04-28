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
	"github.com/nitric-dev/membrane/sdk"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPC Interface for registered Nitric KV Plugin
type KeyValueServer struct {
	pb.UnimplementedKeyValueServer
	// TODO: Support multiple plugin registrations
	// Just need to settle on a way of addressing them on calls
	kvPlugin sdk.KeyValueService
}

func (s *KeyValueServer) checkPluginRegistered() (bool, error) {
	if s.kvPlugin == nil {
		return false, status.Errorf(codes.Unimplemented, "KeyValue plugin not registered")
	}

	return true, nil
}

func (s *KeyValueServer) Put(ctx context.Context, req *pb.KeyValuePutRequest) (*pb.KeyValuePutResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.kvPlugin.Put(req.GetCollection(), req.GetKey(), req.GetValue().AsMap()); err == nil {
			return &pb.KeyValuePutResponse{}, nil
		} else {
			// Case: Failed to create the key
			// TODO: Translate from internal KeyValue Plugin Error
			return nil, err
		}
	} else {
		// Case: Plugin was not registered
		return nil, err
	}
}

func (s *KeyValueServer) Get(ctx context.Context, req *pb.KeyValueGetRequest) (*pb.KeyValueGetResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if val, err := s.kvPlugin.Get(req.GetCollection(), req.GetKey()); err == nil {
			if valStruct, err := structpb.NewStruct(val); err == nil {
				return &pb.KeyValueGetResponse{
					Value: valStruct,
				}, nil
			} else {
				// Case: Failed to create PB struct from stored value
				// TODO: Translate from internal KeyValue Plugin Error
				return nil, err
			}
		} else {
			// Case: There was an error retrieving the keyvalue
			// TODO: Translate from internal KeyValue Plugin Error
			return nil, err
		}
	} else {
		// Case: The keyvalue plugin was not registered
		// TODO: Translate from internal KeyValue Plugin Error
		return nil, err
	}
}

func (s *KeyValueServer) Delete(ctx context.Context, req *pb.KeyValueDeleteRequest) (*pb.KeyValueDeleteResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		if err := s.kvPlugin.Delete(req.GetCollection(), req.GetKey()); err == nil {
			return &pb.KeyValueDeleteResponse{}, nil
		} else {
			// Case: Failed to create the keyvalue
			// TODO: Translate from internal KeyValue Plugin Error
			return nil, err
		}
	} else {
		// Case: Plugin was not registered
		return nil, err
	}
}

func NewKeyValueServer(kvPlugin sdk.KeyValueService) pb.KeyValueServer {
	return &KeyValueServer{
		kvPlugin: kvPlugin,
	}
}
