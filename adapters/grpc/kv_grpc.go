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
	"fmt"
	"reflect"

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
		keyMap, err := toSdkKeyMap(req.GetKey())
		if err != nil {
			return nil, err
		}
		if err := s.kvPlugin.Put(req.GetCollection(), keyMap, req.GetValue().AsMap()); err == nil {
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
		keyMap, err := toSdkKeyMap(req.GetKey())
		if err != nil {
			return nil, err
		}
		if val, err := s.kvPlugin.Get(req.GetCollection(), keyMap); err == nil {
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
		keyMap, err := toSdkKeyMap(req.GetKey())
		if err != nil {
			return nil, err
		}
		if err := s.kvPlugin.Delete(req.GetCollection(), keyMap); err == nil {
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

func (s *KeyValueServer) Query(ctx context.Context, req *pb.KeyValueQueryRequest) (*pb.KeyValueQueryResponse, error) {
	if ok, err := s.checkPluginRegistered(); ok {
		collection := req.GetCollection()
		expressions := make([]sdk.QueryExpression, len(req.GetExpressions()))
		for i, exp := range req.GetExpressions() {
			expressions[i] = sdk.QueryExpression{
				Operand:  exp.GetOperand(),
				Operator: exp.GetOperator(),
				Value:    exp.GetValue(),
			}
		}
		limit := int(req.GetLimit())

		var pagingMap map[string]interface{}
		if req.PagingToken != nil {
			pagingMap, err = toSdkKeyMap(req.PagingToken)
			if err != nil {
				return nil, err
			}
		}

		if qr, err := s.kvPlugin.Query(collection, expressions, limit, pagingMap); err == nil {

			valStructs := make([]*structpb.Struct, len(qr.Data))
			for i, valMap := range qr.Data {
				if valStruct, err := structpb.NewStruct(valMap); err == nil {
					valStructs[i] = valStruct

				} else {
					// Case: Failed to create PB struct from stored value
					// TODO: Translate from internal KeyValue Plugin Error
					return nil, err
				}
			}

			pagingToken, err := toProtoKeyMap(qr.PagingToken)
			if err != nil {
				return nil, err
			}

			return &pb.KeyValueQueryResponse{
				Values:      valStructs,
				PagingToken: pagingToken,
			}, nil

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

func toSdkKeyMap(keyMap map[string]*pb.Key) (map[string]interface{}, error) {
	if keyMap == nil {
		return nil, fmt.Errorf("provide non-nil key")
	}

	sdkMap := make(map[string]interface{})

	for k, v := range keyMap {
		if x, ok := v.GetKey().(*pb.Key_String_); ok {
			sdkMap[k] = x.String_
			break
		}
		if x, ok := v.GetKey().(*pb.Key_Number); ok {
			sdkMap[k] = x.Number
			break
		}
		// Else unsupported type
		return nil, fmt.Errorf("unsupported key type: %v", v)
	}

	return sdkMap, nil
}

func toProtoKeyMap(keyMap map[string]interface{}) (map[string]*pb.Key, error) {
	if keyMap == nil {
		return nil, fmt.Errorf("provide non-nil key")
	}

	protoMap := make(map[string]*pb.Key)

	for k, v := range keyMap {
		valueKind := reflect.ValueOf(v).Kind()
		if valueKind == reflect.String {
			key := pb.Key{
				Key: &pb.Key_String_{String_: v.(string)},
			}
			protoMap[k] = &key

		} else if valueKind == reflect.Int64 {
			key := pb.Key{
				Key: &pb.Key_Number{Number: v.(int64)},
			}
			protoMap[k] = &key

		} else {
			return nil, fmt.Errorf("unsupported new key type %v", v)
		}
	}

	return protoMap, nil
}
