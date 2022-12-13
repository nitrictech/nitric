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
	"io"
	"strings"

	"github.com/nitrictech/nitric/core/pkg/plugins/document"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"

	grpcCodes "google.golang.org/grpc/codes"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/status"
)

const pagingTokens = "pagingTokens"

type FirestoreDocService struct {
	client *firestore.Client
	document.UnimplementedDocumentPlugin
}

func (s *FirestoreDocService) Get(ctx context.Context, key *document.Key) (*document.Document, error) {
	newErr := errors.ErrorsWithScope(
		"FirestoreDocService.Get",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	doc := s.getDocRef(key)

	value, err := doc.Get(ctx)
	if err != nil {
		code := codes.Internal
		if status.Code(err) == grpcCodes.NotFound {
			code = codes.NotFound
		}

		return nil, newErr(
			code,
			"unable to retrieve value",
			err,
		)
	}

	return &document.Document{
		Key:     key,
		Content: value.Data(),
	}, nil
}

func (s *FirestoreDocService) Set(ctx context.Context, key *document.Key, value map[string]interface{}) error {
	newErr := errors.ErrorsWithScope(
		"FirestoreDocService.Set",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if value == nil {
		return newErr(
			codes.InvalidArgument,
			"provide non-nil value",
			nil,
		)
	}

	doc := s.getDocRef(key)

	if _, err := doc.Set(ctx, value); err != nil {
		return newErr(
			codes.Internal,
			"error updating value",
			err,
		)
	}

	return nil
}

func (s *FirestoreDocService) Delete(ctx context.Context, key *document.Key) error {
	newErr := errors.ErrorsWithScope(
		"FirestoreDocService.Delete",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	doc := s.getDocRef(key)

	// Delete any sub collection documents
	collsIter := doc.Collections(ctx)
	for subCol, err := collsIter.Next(); !errors.Is(err, iterator.Done); subCol, err = collsIter.Next() {
		if err != nil {
			return newErr(
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
					return newErr(codes.Internal, "error deleting records", err)
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
				return newErr(codes.Internal, "error deleting records", err)
			}
		}
	}

	// Delete document
	if _, err := doc.Delete(ctx); err != nil {
		return newErr(
			codes.Internal,
			"error deleting value",
			err,
		)
	}

	return nil
}

func (s *FirestoreDocService) buildQuery(collection *document.Collection, expressions []document.QueryExpression, limit int) (query firestore.Query, orderBy string) {
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
		query = query.Limit(limit)
	}

	return
}

func (s *FirestoreDocService) Query(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
	newErr := errors.ErrorsWithScope(
		"FirestoreDocService.Query",
		map[string]interface{}{
			"collection": collection,
		},
	)

	if err := document.ValidateQueryCollection(collection); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if err := document.ValidateExpressions(expressions); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid expressions",
			err,
		)
	}

	queryResult := &document.QueryResult{
		Documents: make([]document.Document, 0),
	}

	// Select correct root collection to perform query on
	query, orderBy := s.buildQuery(collection, expressions, limit)

	if len(pagingToken) > 0 {
		query = query.OrderBy(firestore.DocumentID, firestore.Asc)

		if tokens, ok := pagingToken[pagingTokens]; ok {
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

		sdkDoc := docSnpToDocument(collection, docSnp)
		queryResult.Documents = append(queryResult.Documents, sdkDoc)

		// If query limit configured determine continue tokens
		if limit > 0 && len(queryResult.Documents) == limit {
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

func (s *FirestoreDocService) QueryStream(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int) document.DocumentIterator {
	newErr := errors.ErrorsWithScope(
		"FirestoreDocService.QueryStream",
		map[string]interface{}{
			"collection": collection,
		},
	)

	colErr := document.ValidateQueryCollection(collection)
	expErr := document.ValidateExpressions(expressions)

	if colErr != nil || expErr != nil {
		// Return an error only iterator
		return func() (*document.Document, error) {
			return nil, newErr(
				codes.InvalidArgument,
				"invalid arguments",
				fmt.Errorf("collection error:%w, expression error: %v", colErr, expErr),
			)
		}
	}

	query, _ := s.buildQuery(collection, expressions, limit)

	iter := query.Documents(ctx)

	return func() (*document.Document, error) {
		docSnp, err := iter.Next()
		if err != nil {
			if errors.Is(err, iterator.Done) {
				return nil, io.EOF
			}

			return nil, newErr(
				codes.Internal,
				"error querying value",
				err,
			)
		}

		sdkDoc := docSnpToDocument(collection, docSnp)

		return &sdkDoc, nil
	}
}

func docSnpToDocument(col *document.Collection, snp *firestore.DocumentSnapshot) document.Document {
	sdkDoc := document.Document{
		Content: snp.Data(),
		Key: &document.Key{
			Collection: col,
			Id:         snp.Ref.ID,
		},
	}

	if p := snp.Ref.Parent.Parent; p != nil {
		sdkDoc.Key.Collection = &document.Collection{
			Name: col.Name,
			Parent: &document.Key{
				Collection: col.Parent.Collection,
				Id:         p.ID,
			},
		}
	}

	return sdkDoc
}

func New() (document.DocumentService, error) {
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

func NewWithClient(client *firestore.Client) (document.DocumentService, error) {
	return &FirestoreDocService{
		client: client,
	}, nil
}

func (s *FirestoreDocService) getDocRef(key *document.Key) *firestore.DocumentRef {
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

func (s *FirestoreDocService) getQueryRoot(collection *document.Collection) firestore.Query {
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
