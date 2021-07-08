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
	"google.golang.org/protobuf/types/known/structpb"
)

// GRPC Interface for registered Nitric Document Plugin
type DocumentServer struct {
	pb.UnimplementedDocumentServer
	// TODO: Support multiple plugin registrations
	// Just need to settle on a way of addressing them on calls
	documentPlugin sdk.DocumentService
}

func (s *DocumentServer) Set(ctx context.Context, req *pb.DocumentSetRequest) (*pb.DocumentSetResponse, error) {
	key := toSdkKey(req.Key)
	subKey := toSdkKey(req.SubKey)

	if err := s.documentPlugin.Set(key, subKey, req.GetValue().AsMap()); err == nil {
		return &pb.DocumentSetResponse{}, nil
	} else {
		// Case: Failed to create the key
		// TODO: Translate from internal Document Service Error
		return nil, err
	}
}

func (s *DocumentServer) Get(ctx context.Context, req *pb.DocumentGetRequest) (*pb.DocumentGetResponse, error) {
	key := toSdkKey(req.Key)
	subKey := toSdkKey(req.SubKey)

	if val, err := s.documentPlugin.Get(key, subKey); err == nil {
		if valStruct, err := structpb.NewStruct(val); err == nil {
			return &pb.DocumentGetResponse{
				Value: valStruct,
			}, nil
		} else {
			// Case: Failed to create PB struct from stored value
			// TODO: Translate from internal Document Plugin Error
			return nil, err
		}
	} else {
		// Case: There was an error retrieving the keyvalue
		// TODO: Translate from internal Document Plugin Error
		return nil, err
	}
}

func (s *DocumentServer) Delete(ctx context.Context, req *pb.DocumentDeleteRequest) (*pb.DocumentDeleteResponse, error) {
	key := toSdkKey(req.Key)
	subKey := toSdkKey(req.SubKey)

	if err := s.documentPlugin.Delete(key, subKey); err == nil {
		return &pb.DocumentDeleteResponse{}, nil
	} else {
		// Case: Failed to create the keyvalue
		// TODO: Translate from internal Document Plugin Error
		return nil, err
	}
}

func (s *DocumentServer) Query(ctx context.Context, req *pb.DocumentQueryRequest) (*pb.DocumentQueryResponse, error) {
	key := toSdkKey(req.Key)
	subcoll := req.GetSubCollection()
	expressions := make([]sdk.QueryExpression, len(req.GetExpressions()))
	for i, exp := range req.GetExpressions() {
		expressions[i] = sdk.QueryExpression{
			Operand:  exp.GetOperand(),
			Operator: exp.GetOperator(),
			Value:    exp.GetValue(),
		}
	}
	limit := int(req.GetLimit())
	pagingMap := req.GetPagingToken()

	if qr, err := s.documentPlugin.Query(key, subcoll, expressions, limit, pagingMap); err == nil {

		valStructs := make([]*structpb.Struct, len(qr.Data))
		for i, valMap := range qr.Data {
			if valStruct, err := structpb.NewStruct(valMap); err == nil {
				valStructs[i] = valStruct

			} else {
				// Case: Failed to create PB struct from stored value
				// TODO: Translate from internal Document Plugin Error
				return nil, err
			}
		}

		return &pb.DocumentQueryResponse{
			Values:      valStructs,
			PagingToken: qr.PagingToken,
		}, nil

	} else {
		// Case: Failed to create the keyvalue
		// TODO: Translate from internal Document Plugin Error
		return nil, err
	}
}

func NewDocumentServer(docPlugin sdk.DocumentService) pb.DocumentServer {
	return &DocumentServer{
		documentPlugin: docPlugin,
	}
}

func toSdkKey(key *pb.Key) *sdk.Key {
	if key != nil {
		return &sdk.Key{
			Collection: key.GetCollection(),
			Id:         key.GetId(),
		}
	}
	return nil
}
