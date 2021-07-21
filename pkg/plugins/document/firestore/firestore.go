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
	"github.com/nitric-dev/membrane/pkg/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

const pagingTokens = "pagingTokens"

type FirestoreDocService struct {
	client  *firestore.Client
	context context.Context
	sdk.UnimplementedDocumentPlugin
}

func (s *FirestoreDocService) Get(key *sdk.Key) (*sdk.Document, error) {
	err := document.ValidateKey(key)
	if err != nil {
		return nil, err
	}

	doc := s.getDocRef(key)

	value, err := doc.Get(s.context)
	if err != nil {
		return nil, fmt.Errorf("error retrieving value: %v", err)
	}

	return &sdk.Document{
		Content: value.Data(),
	}, nil
}

func (s *FirestoreDocService) Set(key *sdk.Key, content map[string]interface{}, merge bool) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	if content == nil {
		return fmt.Errorf("provide non-nil content")
	}

	doc := s.getDocRef(key)

	if merge {
		_, err = doc.Set(s.context, content, firestore.MergeAll)
		if err != nil {
			return fmt.Errorf("error updating content: %v", err)
		}

	} else {
		_, err = doc.Set(s.context, content)
		if err != nil {
			return fmt.Errorf("error updating content: %v", err)
		}
	}

	return nil
}

func (s *FirestoreDocService) Delete(key *sdk.Key) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	doc := s.getDocRef(key)

	_, err = doc.Delete(s.context)
	if err != nil {
		return fmt.Errorf("error deleting value: %v", err)
	}

	// TODO: delete sub collection records

	return nil
}

func (s *FirestoreDocService) Query(collection *sdk.Collection, expressions []sdk.QueryExpression, limit int, pagingToken map[string]string) (*sdk.QueryResult, error) {
	err := document.ValidateQueryCollection(collection)
	if err != nil {
		return nil, err
	}

	err = document.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	queryResult := &sdk.QueryResult{
		Documents: make([]sdk.Document, 0),
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

	for {
		docSnp, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error querying value: %v", err)
		}
		sdkDoc := sdk.Document{Content: docSnp.Data()}
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

func New() (sdk.DocumentService, error) {
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

func NewWithClient(client *firestore.Client, ctx context.Context) (sdk.DocumentService, error) {
	return &FirestoreDocService{
		client:  client,
		context: ctx,
	}, nil
}

func (s *FirestoreDocService) getDocRef(key *sdk.Key) *firestore.DocumentRef {
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

func (s *FirestoreDocService) getQueryRoot(collection *sdk.Collection) firestore.Query {
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
