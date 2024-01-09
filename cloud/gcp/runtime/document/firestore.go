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

package document

import (
	"context"
	"fmt"
	"strings"

	"errors"

	document "github.com/nitrictech/nitric/core/pkg/decorators/documents"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/documents/v1"

	"google.golang.org/grpc/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/status"
)

const pagingTokens = "pagingTokens"

type FirestoreDocService struct {
	client *firestore.Client
}

var _ v1.DocumentsServer = &FirestoreDocService{}

func (s *FirestoreDocService) Get(ctx context.Context, req *v1.DocumentGetRequest) (*v1.DocumentGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Get")

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	doc := s.getDocRef(req.Key)

	value, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == grpcCodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this collection?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"unable to retrieve value",
			err,
		)
	}

	documentContent, err := structpb.NewStruct(value.Data())
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error converting returned document to struct",
			err,
		)
	}

	return &v1.DocumentGetResponse{
		Document: &v1.Document{
			Key:     req.Key,
			Content: documentContent,
		},
	}, nil
}

func (s *FirestoreDocService) Set(ctx context.Context, req *v1.DocumentSetRequest) (*v1.DocumentSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Set")

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if req.Content == nil {
		return nil, newErr(
			codes.InvalidArgument,
			"provide non-nil value",
			nil,
		)
	}

	doc := s.getDocRef(req.Key)

	if _, err := doc.Set(ctx, req.Content.AsMap()); err != nil {
		if status.Code(err) == grpcCodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this collection?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"error updating value",
			err,
		)
	}

	return &v1.DocumentSetResponse{}, nil
}

func (s *FirestoreDocService) Delete(ctx context.Context, req *v1.DocumentDeleteRequest) (*v1.DocumentDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Delete")

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	doc := s.getDocRef(req.Key)

	// Delete any sub collection documents
	collsIter := doc.Collections(ctx)
	for subCol, err := collsIter.Next(); !errors.Is(err, iterator.Done); subCol, err = collsIter.Next() {
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error deleting value",
				err,
			)
		}

		// Loop over sub collection documents, performing batch deletes
		// up to Firestore's maximum batch size
		const maxBatchSize = 20
		for {
			docsIter := subCol.Limit(maxBatchSize).Documents(ctx)
			numDeleted := 0

			batch := s.client.Batch()
			for subDoc, err := docsIter.Next(); !errors.Is(err, iterator.Done); subDoc, err = docsIter.Next() {
				if err != nil {
					return nil, newErr(codes.Internal, "error deleting records", err)
				}

				batch.Delete(subDoc.Ref)
				numDeleted++
			}

			// If no more to delete, completed
			if numDeleted == 0 {
				break
			}

			_, err := batch.Commit(ctx)
			if err != nil {
				return nil, newErr(codes.Internal, "error deleting records", err)
			}
		}
	}

	// Delete document
	if _, err := doc.Delete(ctx); err != nil {
		if status.Code(err) == grpcCodes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this collection?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"error deleting value",
			err,
		)
	}

	return &v1.DocumentDeleteResponse{}, nil
}

func (s *FirestoreDocService) buildQuery(collection *v1.Collection, expressions []*v1.Expression, limit int32) (query firestore.Query, orderBy string) {
	// Select correct root collection to perform query on
	query = s.getQueryRoot(collection)

	for _, exp := range expressions {
		expOperand := exp.Operand
		if exp.Operator == "startsWith" {
			expVal := fmt.Sprintf("%v", exp.Value)
			endRangeValue := document.GetEndRangeValue(expVal)
			query = query.Where(expOperand, ">=", exp.Value).Where(expOperand, "<", endRangeValue)
		} else {
			query = query.Where(expOperand, exp.Operator, exp.Value)
		}

		if exp.Operator != "==" && limit > 0 && orderBy == "" {
			query = query.OrderBy(expOperand, firestore.Asc)
			orderBy = expOperand
		}
	}

	if limit > 0 {
		query = query.Limit(int(limit))
	}

	return
}

func (s *FirestoreDocService) Query(ctx context.Context, req *v1.DocumentQueryRequest) (*v1.DocumentQueryResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Query")

	if err := document.ValidateQueryCollection(req.Collection); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if err := document.ValidateExpressions(req.Expressions); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid expressions",
			err,
		)
	}

	queryResult := &v1.DocumentQueryResponse{
		Documents: make([]*v1.Document, 0),
	}

	// Select correct root collection to perform query on
	query, orderBy := s.buildQuery(req.Collection, req.Expressions, req.Limit)

	if len(req.PagingToken) > 0 {
		query = query.OrderBy(firestore.DocumentID, firestore.Asc)

		if tokens, ok := req.PagingToken[pagingTokens]; ok {
			var vals []interface{}
			for _, v := range strings.Split(tokens, "|") {
				vals = append(vals, v)
			}
			query = query.StartAfter(vals...)
		}
	}

	itr := query.Documents(ctx)
	for docSnp, err := itr.Next(); !errors.Is(err, iterator.Done); docSnp, err = itr.Next() {
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error querying value",
				err,
			)
		}

		sdkDoc, err := firestoreSnapshotToDocument(req.Collection, docSnp)
		if err != nil {
			return nil, err
		}

		queryResult.Documents = append(queryResult.Documents, sdkDoc)

		// If query limit configured determine continue tokens
		if req.Limit > 0 && int32(len(queryResult.Documents)) == req.Limit {
			tokens := ""
			if orderBy != "" {
				tokens = fmt.Sprintf("%v", docSnp.Data()[orderBy]) + "|"
			}
			tokens += docSnp.Ref.ID

			queryResult.PagingToken = map[string]string{
				pagingTokens: tokens,
			}
		}
	}

	return queryResult, nil
}

func (s *FirestoreDocService) QueryStream(req *v1.DocumentQueryStreamRequest, srv v1.Documents_QueryStreamServer) error {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.QueryStream")

	colErr := document.ValidateQueryCollection(req.Collection)
	if colErr != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid arguments",
			fmt.Errorf("collection error: %w", colErr),
		)
	}

	expErr := document.ValidateExpressions(req.Expressions)
	if expErr != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid arguments",
			fmt.Errorf("expression error: %w", expErr),
		)
	}

	query, _ := s.buildQuery(req.Collection, req.Expressions, req.Limit)
	iter := query.Documents(srv.Context())

	for {
		docSnp, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				return nil
			}

			return newErr(
				codes.Internal,
				"error querying value",
				err,
			)
		}

		sdkDoc, err := firestoreSnapshotToDocument(req.Collection, docSnp)
		if err != nil {
			return newErr(codes.Internal, "error serializing firestore document", err)
		}

		if err := srv.Send(&v1.DocumentQueryStreamResponse{
			Document: sdkDoc,
		}); err != nil {
			return newErr(
				codes.Internal,
				"error sending document",
				err,
			)
		}
	}
}

func firestoreSnapshotToDocument(col *v1.Collection, snapshot *firestore.DocumentSnapshot) (*v1.Document, error) {
	documentContent, err := structpb.NewStruct(snapshot.Data())
	if err != nil {
		return nil, err
	}

	doc := &v1.Document{
		Content: documentContent,
		Key: &v1.Key{
			Collection: col,
			Id:         snapshot.Ref.ID,
		},
	}

	if p := snapshot.Ref.Parent.Parent; p != nil {
		doc.Key.Collection = &v1.Collection{
			Name: col.Name,
			Parent: &v1.Key{
				Collection: col.Parent.Collection,
				Id:         p.ID,
			},
		}
	}

	return doc, nil
}

func New() (v1.DocumentsServer, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %w", credentialsError)
	}

	client, clientError := firestore.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("firestore client error: %w", clientError)
	}

	return &FirestoreDocService{
		client: client,
	}, nil
}

func NewWithClient(client *firestore.Client) (v1.DocumentsServer, error) {
	return &FirestoreDocService{
		client: client,
	}, nil
}

func (s *FirestoreDocService) getDocRef(key *v1.Key) *firestore.DocumentRef {
	parentKey := key.Collection.Parent

	if parentKey == nil {
		return s.client.Collection(key.Collection.Name).Doc(key.Id)
	} else {
		return s.client.Collection(parentKey.Collection.Name).
			Doc(parentKey.Id).
			Collection(key.Collection.Name).
			Doc(key.Id)
	}
}

func (s *FirestoreDocService) getQueryRoot(collection *v1.Collection) firestore.Query {
	parentKey := collection.Parent

	if parentKey == nil {
		return s.client.Collection(collection.Name).Offset(0)
	} else {
		if parentKey.Id != "" {
			return s.client.Collection(parentKey.Collection.Name).
				Doc(parentKey.Id).
				Collection(collection.Name).
				Offset(0)
		} else {
			// Note there is a risk of subcollection name collison
			// TODO: future YAML validation could help mitigate this
			return s.client.CollectionGroup(collection.Name).Offset(0)
		}
	}
}
