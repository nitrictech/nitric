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

package mock_firestore

import (
	"context"
	"fmt"
	"strings"

	pb "google.golang.org/genproto/googleapis/firestore/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Mock Firestore implementation for plugin testing
type MockFirestoreServer struct {
	pb.UnimplementedFirestoreServer
	Store map[string]map[string]map[string]*pb.Value
}

// NOTE: On handling paths..., the collection and resourceId are the last two components
// The full resource path of the document. A document "doc-1" in collection
// "coll-1" would be: "projects/P/databases/D/documents/coll-1/doc-1".
func (m *MockFirestoreServer) GetDocument(ctx context.Context, req *pb.GetDocumentRequest) (*pb.Document, error) {
	path := req.GetName()
	parts := strings.Split(path, "/")
	collection := parts[len(parts)-2]
	key := parts[len(parts)-1]

	if values, ok := m.Store[collection][key]; ok {
		return &pb.Document{
			Name:   req.GetName(),
			Fields: values,
		}, nil
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Document %s not found", path))
}

func (m *MockFirestoreServer) ListDocuments(context.Context, *pb.ListDocumentsRequest) (*pb.ListDocumentsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDocuments not implemented")
}

func (m *MockFirestoreServer) UpdateDocument(context.Context, *pb.UpdateDocumentRequest) (*pb.Document, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func (m *MockFirestoreServer) DeleteDocument(ctx context.Context, req *pb.DeleteDocumentRequest) (*emptypb.Empty, error) {
	path := req.GetName()
	parts := strings.Split(path, "/")
	collection := parts[len(parts)-2]
	key := parts[len(parts)-1]

	if _, ok := m.Store[collection][key]; ok {
		// Clear the reference
		m.Store[collection][key] = nil
		return &emptypb.Empty{}, nil
	}

	return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Document %s not found", path))
}

func (m *MockFirestoreServer) ClearStore() {
	// Nuke the previous references...
	m.Store = make(map[string]map[string]map[string]*pb.Value)
}

func (m *MockFirestoreServer) BatchGetDocuments(req *pb.BatchGetDocumentsRequest, stream pb.Firestore_BatchGetDocumentsServer) error {
	for _, docName := range req.GetDocuments() {
		parts := strings.Split(docName, "/")
		collection := parts[len(parts)-2]
		key := parts[len(parts)-1]

		var err error
		if doc, ok := m.Store[collection][key]; ok {
			currentTime := timestamppb.Now()
			err = stream.Send(&pb.BatchGetDocumentsResponse{
				Result: &pb.BatchGetDocumentsResponse_Found{
					Found: &pb.Document{
						Name:       docName,
						Fields:     doc,
						CreateTime: currentTime,
						UpdateTime: currentTime,
					},
				},
			})

			if err != nil {
				return err
			}
		} else {
			err = stream.Send(&pb.BatchGetDocumentsResponse{
				Result: &pb.BatchGetDocumentsResponse_Missing{
					Missing: docName,
				},
			})
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MockFirestoreServer) BeginTransaction(context.Context, *pb.BeginTransactionRequest) (*pb.BeginTransactionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BeginTransaction not implemented")
}

func (m *MockFirestoreServer) Commit(ctx context.Context, req *pb.CommitRequest) (*pb.CommitResponse, error) {
	// Return a fake commit response...
	for _, write := range req.GetWrites() {
		if op, ok := write.GetOperation().(*pb.Write_Update); ok {
			document := op.Update
			path := document.GetName()
			parts := strings.Split(path, "/")
			collection := parts[len(parts)-2]
			key := parts[len(parts)-1]

			if _, ok := m.Store[collection]; !ok {
				m.Store[collection] = make(map[string]map[string]*pb.Value)
			}

			_, docExists := m.Store[collection][key]

			// If the existing state of the document does not match our write precondition
			// Then we need to throw an error
			if condition, ok := write.GetCurrentDocument().GetConditionType().(*pb.Precondition_Exists); ok {
				if condition.Exists != docExists {
					var returnCode codes.Code
					if ok {
						returnCode = codes.AlreadyExists
					} else {
						returnCode = codes.NotFound
					}

					return nil, status.Errorf(returnCode, fmt.Sprintf("Item: %s/%s", collection, key))
				}
			}

			m.Store[collection][key] = document.GetFields()
		} else if op, ok := write.GetOperation().(*pb.Write_Delete); ok {
			path := op.Delete
			parts := strings.Split(path, "/")
			collection := parts[len(parts)-2]
			key := parts[len(parts)-1]

			if _, ok := m.Store[collection]; !ok {
				return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Collection: %s, does not exist", collection))
			}

			if _, ok := m.Store[collection][key]; !ok {
				return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Item: %s/%s, does not exist", collection, key))
			}

			// Let's delete it...
			m.Store[collection][key] = nil
		}
	}

	return &pb.CommitResponse{
		WriteResults: []*pb.WriteResult{
			{
				TransformResults: make([]*pb.Value, 0),
			},
		},
	}, nil
}

func (m *MockFirestoreServer) Rollback(context.Context, *pb.RollbackRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Rollback not implemented")
}

func (m *MockFirestoreServer) RunQuery(*pb.RunQueryRequest, pb.Firestore_RunQueryServer) error {
	return status.Errorf(codes.Unimplemented, "method RunQuery not implemented")
}

func (m *MockFirestoreServer) PartitionQuery(context.Context, *pb.PartitionQueryRequest) (*pb.PartitionQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PartitionQuery not implemented")
}

func (m *MockFirestoreServer) Write(pb.Firestore_WriteServer) error {
	return status.Errorf(codes.Unimplemented, "method Write not implemented")
}

func (m *MockFirestoreServer) Listen(pb.Firestore_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}

func (m *MockFirestoreServer) ListCollectionIds(context.Context, *pb.ListCollectionIdsRequest) (*pb.ListCollectionIdsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCollectionIds not implemented")
}

func (m *MockFirestoreServer) BatchWrite(context.Context, *pb.BatchWriteRequest) (*pb.BatchWriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchWrite not implemented")
}

func (m *MockFirestoreServer) CreateDocument(ctx context.Context, req *pb.CreateDocumentRequest) (*pb.Document, error) {
	collection := req.CollectionId
	key := req.DocumentId
	document := req.Document.Fields

	if _, ok := m.Store[collection][key]; ok {
		return nil, status.Errorf(codes.AlreadyExists, fmt.Sprintf("Item: %s/%s already exists", collection, key))
	}

	m.Store[collection][key] = document

	// TODO: Probably need to set the document Name here post creation...
	return req.Document, nil
}
