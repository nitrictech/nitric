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
	"io"

	"google.golang.org/grpc/codes"

	pb "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/document"
	"github.com/nitrictech/protoutils"
)

// DocumentServiceServer - GRPC Interface for registered Nitric Document Plugin
type DocumentServiceServer struct {
	pb.UnimplementedDocumentServiceServer
	// TODO: Support multiple plugin registrations
	// Just need to settle on a way of addressing them on calls
	documentPlugin document.DocumentService
}

func (s *DocumentServiceServer) checkPluginRegistered() error {
	if s.documentPlugin == nil {
		return NewPluginNotRegisteredError("Document")
	}

	return nil
}

func (s *DocumentServiceServer) Get(ctx context.Context, req *pb.DocumentGetRequest) (*pb.DocumentGetResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "DocumentService.Get", err)
	}

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
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "DocumentService.Set", err)
	}

	key := keyFromWire(req.Key)

	err := s.documentPlugin.Set(key, req.GetContent().AsMap())
	if err != nil {
		return nil, NewGrpcError("DocumentService.Set", err)
	}

	return &pb.DocumentSetResponse{}, nil
}

func (s *DocumentServiceServer) Delete(ctx context.Context, req *pb.DocumentDeleteRequest) (*pb.DocumentDeleteResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "DocumentService.Delete", err)
	}

	key := keyFromWire(req.Key)

	err := s.documentPlugin.Delete(key)
	if err != nil {
		return nil, NewGrpcError("DocumentService.Delete", err)
	}

	return &pb.DocumentDeleteResponse{}, nil
}

func (s *DocumentServiceServer) Query(ctx context.Context, req *pb.DocumentQueryRequest) (*pb.DocumentQueryResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "DocumentService.Query", err)
	}

	collection := collectionFromWire(req.Collection)
	expressions := expressionsFromWire(req.GetExpressions())

	limit := int(req.GetLimit())
	pagingMap := req.GetPagingToken()

	qr, err := s.documentPlugin.Query(collection, expressions, limit, pagingMap)
	if err != nil {
		return nil, NewGrpcError("DocumentService.Query", err)
	}

	pbDocuments := make([]*pb.Document, 0, len(qr.Documents))
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

func (s *DocumentServiceServer) QueryStream(req *pb.DocumentQueryStreamRequest, srv pb.DocumentService_QueryStreamServer) error {
	if err := s.checkPluginRegistered(); err != nil {
		return err
	}

	col := collectionFromWire(req.Collection)
	expressions := expressionsFromWire(req.Expressions)

	next := s.documentPlugin.QueryStream(col, expressions, int(req.Limit))

	for doc, err := next(); err != io.EOF; doc, err = next() {
		if err != nil {
			return NewGrpcError("DocumentService.QueryStream", err)
		}

		if d, docErr := documentToWire(doc); docErr != nil {
			return NewGrpcError("DocumentService.QueryStream", err)
		} else {
			srv.Send(&pb.DocumentQueryStreamResponse{
				Document: d,
			})
		}
	}

	return nil
}

func NewDocumentServer(docPlugin document.DocumentService) pb.DocumentServiceServer {
	return &DocumentServiceServer{
		documentPlugin: docPlugin,
	}
}

func documentToWire(doc *document.Document) (*pb.Document, error) {
	valStruct, err := protoutils.NewStruct(doc.Content)
	if err != nil {
		return nil, err
	}

	return &pb.Document{
		Content: valStruct,
		Key:     keyToWire(doc.Key),
	}, nil
}

// keyFromWire - returns an Membrane SDK Document Key from the protobuf wire representation
// recursively calls collectionFromWire for the document's collection and parents if present
func keyFromWire(key *pb.Key) *document.Key {
	if key == nil {
		return nil
	}

	sdkKey := &document.Key{
		Collection: collectionFromWire(key.GetCollection()),
		Id:         key.GetId(),
	}

	return sdkKey
}

// keyToWire - translates an SDK key to a gRPC key
func keyToWire(key *document.Key) *pb.Key {
	return &pb.Key{
		Id:         key.Id,
		Collection: collectionToWire(key.Collection),
	}
}

// collectionToWire - translates a SDK collection to a gRPC collection
func collectionToWire(col *document.Collection) *pb.Collection {
	if col.Parent != nil {
		return &pb.Collection{
			Name:   col.Name,
			Parent: keyToWire(col.Parent),
		}
	} else {
		return &pb.Collection{
			Name: col.Name,
		}
	}
}

// collectionFromWire - returns an Membrane SDK Document Collection from the protobuf wire representation
// recursively calls keyFromWire if the collection is a sub-collection under another key
func collectionFromWire(coll *pb.Collection) *document.Collection {
	if coll == nil {
		return nil
	}

	if coll.Parent == nil {
		return &document.Collection{
			Name: coll.Name,
		}
	} else {
		return &document.Collection{
			Name:   coll.Name,
			Parent: keyFromWire(coll.Parent),
		}
	}
}

func expressionsFromWire(exps []*pb.Expression) []document.QueryExpression {
	expressions := make([]document.QueryExpression, len(exps))
	for i, exp := range exps {
		expressions[i] = document.QueryExpression{
			Operand:  exp.GetOperand(),
			Operator: exp.GetOperator(),
			Value:    toExpValue(exp.GetValue()),
		}
	}

	return expressions
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
