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
	"google.golang.org/protobuf/types/known/structpb"
)

// DocumentServiceServer - GRPC Interface for registered Nitric Document Plugin
type DocumentServiceServer struct {
	pb.UnimplementedDocumentServiceServer
	// TODO: Support multiple plugin registrations
	// Just need to settle on a way of addressing them on calls
	documentPlugin sdk.DocumentService
}

func (s *DocumentServiceServer) Get(ctx context.Context, req *pb.DocumentGetRequest) (*pb.DocumentGetResponse, error) {
	key := keyFromWire(req.Key)

	doc, err := s.documentPlugin.Get(key)
	if err != nil {
		return nil, NewGrpcError("DocumentService.Get", err)
	}

	pbDoc, err := documentToWire(doc)
	if err != nil {
		return nil, NewGrpcError("DocumentService.Get", err)
	}

	return &pb.DocumentGetResponse{
		Document: pbDoc,
	}, nil
}

func (s *DocumentServiceServer) Set(ctx context.Context, req *pb.DocumentSetRequest) (*pb.DocumentSetResponse, error) {
	key := keyFromWire(req.Key)

	err := s.documentPlugin.Set(key, req.GetContent().AsMap())
	if err != nil {
		return nil, NewGrpcError("DocumentService.Set", err)
	}

	return &pb.DocumentSetResponse{}, nil
}

func (s *DocumentServiceServer) Delete(ctx context.Context, req *pb.DocumentDeleteRequest) (*pb.DocumentDeleteResponse, error) {
	key := keyFromWire(req.Key)

	err := s.documentPlugin.Delete(key)
	if err != nil {
		return nil, NewGrpcError("DocumentService.Delete", err)
	}

	return &pb.DocumentDeleteResponse{}, nil
}

func (s *DocumentServiceServer) Query(ctx context.Context, req *pb.DocumentQueryRequest) (*pb.DocumentQueryResponse, error) {
	collection := collectionFromWire(req.Collection)
	expressions := make([]sdk.QueryExpression, len(req.GetExpressions()))
	for i, exp := range req.GetExpressions() {
		expressions[i] = sdk.QueryExpression{
			Operand:  exp.GetOperand(),
			Operator: exp.GetOperator(),
			Value:    toExpValue(exp.GetValue()),
		}
	}
	limit := int(req.GetLimit())
	pagingMap := req.GetPagingToken()

	qr, err := s.documentPlugin.Query(collection, expressions, limit, pagingMap)
	if err != nil {
		return nil, NewGrpcError("DocumentService.Query", err)
	}

	pbDocuments := make([]*pb.Document, len(qr.Documents))
	for _, doc := range qr.Documents {
		pbDoc, err := documentToWire(&doc)
		if err != nil {
			return nil, NewGrpcError("DocumentService.Query", err)
		}

		pbDocuments = append(pbDocuments, pbDoc)
	}

	return &pb.DocumentQueryResponse{
		Documents:   pbDocuments,
		PagingToken: qr.PagingToken,
	}, nil
}

func NewDocumentServer(docPlugin sdk.DocumentService) pb.DocumentServiceServer {
	return &DocumentServiceServer{
		documentPlugin: docPlugin,
	}
}

func documentToWire(doc *sdk.Document) (*pb.Document, error) {
	valStruct, err := structpb.NewStruct(doc.Content)
	if err != nil {
		return nil, err
	}

	return &pb.Document{
		Content: valStruct,
	}, nil
}

// keyFromWire - returns an Membrane SDK Document Key from the protobuf wire representation
// recursively calls collectionFromWire for the document's collection and parents if present
func keyFromWire(key *pb.Key) *sdk.Key {
	if key == nil {
		return nil
	}

	sdkKey := &sdk.Key{
		Collection: collectionFromWire(key.GetCollection()),
		Id:         key.GetId(),
	}

	return sdkKey
}

// collectionFromWire - returns an Membrane SDK Document Collection from the protobuf wire representation
// recursively calls keyFromWire if the collection is a sub-collection under another key
func collectionFromWire(coll *pb.Collection) *sdk.Collection {
	if coll == nil {
		return nil
	}

	if coll.Parent == nil {
		return &sdk.Collection{
			Name: coll.Name,
		}
	} else {
		return &sdk.Collection{
			Name:   coll.Name,
			Parent: keyFromWire(coll.Parent),
		}
	}
}

func toExpValue(x *pb.ExpressionValue) interface{} {
	if x, ok := x.GetKind().(*pb.ExpressionValue_IntValue); ok {
		return x.IntValue
	}
	if x, ok := x.GetKind().(*pb.ExpressionValue_DoubleValue); ok {
		return x.DoubleValue
	}
	if x, ok := x.GetKind().(*pb.ExpressionValue_StringValue); ok {
		return x.StringValue
	}
	if x, ok := x.GetKind().(*pb.ExpressionValue_BoolValue); ok {
		return x.BoolValue
	}
	return nil
}
