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

package firestore_service

import (
	"context"
	"fmt"
	"strings"

	"github.com/nitric-dev/membrane/pkg/plugins/document"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

const pagingTokens = "pagingTokens"

type FirestoreDocService struct {
	client  *firestore.Client
	context context.Context
	document.UnimplementedDocumentPlugin
}

func (s *FirestoreDocService) Get(key *document.Key) (*document.Document, error) {
	err := document.ValidateKey(key)
	if err != nil {
		return nil, err
	}

	doc := s.getDocRef(key)

	value, err := doc.Get(s.context)
	if err != nil {
		return nil, fmt.Errorf("error retrieving value: %v", err)
	}

	return &document.Document{
		Key:     key,
		Content: value.Data(),
	}, nil
}

func (s *FirestoreDocService) Set(key *document.Key, value map[string]interface{}) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	doc := s.getDocRef(key)

	_, err = doc.Set(s.context, value)
	if err != nil {
		return fmt.Errorf("error updating value: %v", err)
	}

	return nil
}

func (s *FirestoreDocService) Delete(key *document.Key) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	doc := s.getDocRef(key)

	// Delete any sub collection documents
	collsIter := doc.Collections(s.context)
	for subCol, err := collsIter.Next(); err != iterator.Done; subCol, err = collsIter.Next() {
		if err != nil {
			return fmt.Errorf("error deleting value: %v", err)
		}

		// Loop over sub collection documents, performing batch deletes
		// up to Firestore's maximum batch size
		const maxBatchSize = 500
		for {
			docsIter := subCol.Limit(maxBatchSize).Documents(s.context)
			numDeleted := 0

			batch := s.client.Batch()
			for subDoc, err := docsIter.Next(); err != iterator.Done; subDoc, err = docsIter.Next() {
				if err != nil {
					return err
				}

				batch.Delete(subDoc.Ref)
				numDeleted++
			}

			// If no more to delete, completed
			if numDeleted == 0 {
				break
			}

			_, err := batch.Commit(s.context)
			if err != nil {
				return err
			}
		}
	}

	// Delete document
	_, err = doc.Delete(s.context)
	if err != nil {
		return fmt.Errorf("error deleting value: %v", err)
	}

	return nil
}

func (s *FirestoreDocService) Query(collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
	err := document.ValidateQueryCollection(collection)
	if err != nil {
		return nil, err
	}

	err = document.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	queryResult := &document.QueryResult{
		Documents: make([]document.Document, 0),
	}

	// Select correct root collection to perform query on
	query := s.getQueryRoot(collection)

	var orderByAttrib string

	for _, exp := range expressions {
		expOperand := exp.Operand
		if exp.Operator == "startsWith" {
			expVal := fmt.Sprintf("%v", exp.Value)
			endRangeValue := document.GetEndRangeValue(expVal)
			query = query.Where(expOperand, ">=", exp.Value).Where(expOperand, "<", endRangeValue)

		} else {
			query = query.Where(expOperand, exp.Operator, exp.Value)
		}

		if exp.Operator != "==" && limit > 0 && orderByAttrib == "" {
			query = query.OrderBy(expOperand, firestore.Asc)
			orderByAttrib = expOperand
		}
	}

	if limit > 0 {
		query = query.Limit(limit)

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
	}

	itr := query.Documents(s.context)
	for docSnp, err := itr.Next(); err != iterator.Done; docSnp, err = itr.Next() {
		if err != nil {
			return nil, fmt.Errorf("error querying value: %v", err)
		}
		sdkDoc := document.Document{
			Content: docSnp.Data(),
			Key: &document.Key{
				Collection: collection,
				Id:         docSnp.Ref.ID,
			},
		}

		if p := docSnp.Ref.Parent.Parent; p != nil {
			sdkDoc.Key.Collection = &document.Collection{
				Name: collection.Name,
				Parent: &document.Key{
					Collection: collection.Parent.Collection,
					Id:         p.ID,
				},
			}
		}

		queryResult.Documents = append(queryResult.Documents, sdkDoc)

		// If query limit configured determine continue tokens
		if limit > 0 && len(queryResult.Documents) == limit {
			tokens := ""
			if orderByAttrib != "" {
				tokens = fmt.Sprintf("%v", docSnp.Data()[orderByAttrib]) + "|"
			}
			tokens += docSnp.Ref.ID

			queryResult.PagingToken = map[string]string{
				pagingTokens: tokens,
			}
		}
	}

	return queryResult, nil
}

func New() (document.DocumentService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}

	client, clientError := firestore.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("firestore client error: %v", clientError)
	}

	return &FirestoreDocService{
		client:  client,
		context: ctx,
	}, nil
}

func NewWithClient(client *firestore.Client, ctx context.Context) (document.DocumentService, error) {
	return &FirestoreDocService{
		client:  client,
		context: ctx,
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
