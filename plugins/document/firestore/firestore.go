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

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/document"
	"github.com/nitric-dev/membrane/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

const pagingTokens = "pagingTokens"

type FirestoreDocService struct {
	client  *firestore.Client
	context context.Context
	sdk.UnimplementedDocumentPlugin
}

func (s *FirestoreDocService) Get(key sdk.Key, subKey *sdk.Key) (map[string]interface{}, error) {
	err := document.ValidateKeys(key, subKey)
	if err != nil {
		return nil, err
	}

	doc := s.client.Collection(key.Collection).Doc(key.Id)

	if subKey != nil {
		doc = doc.Collection(subKey.Collection).Doc(subKey.Id)
	}

	value, err := doc.Get(s.context)
	if err != nil {
		return nil, fmt.Errorf("error retrieving value: %v", err)
	}

	return value.Data(), nil
}

func (s *FirestoreDocService) Set(key sdk.Key, subKey *sdk.Key, value map[string]interface{}) error {
	err := document.ValidateKeys(key, subKey)
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	doc := s.client.Collection(key.Collection).Doc(key.Id)

	if subKey != nil {
		doc = doc.Collection(subKey.Collection).Doc(subKey.Id)
	}

	_, err = doc.Set(s.context, value)
	if err != nil {
		return fmt.Errorf("error updating value: %v", err)
	}

	return nil
}

func (s *FirestoreDocService) Delete(key sdk.Key, subKey *sdk.Key) error {
	err := document.ValidateKeys(key, subKey)
	if err != nil {
		return err
	}

	doc := s.client.Collection(key.Collection).Doc(key.Id)

	if subKey != nil {
		doc = doc.Collection(subKey.Collection).Doc(subKey.Id)
	}

	_, err = doc.Delete(s.context)
	if err != nil {
		return fmt.Errorf("error deleting value: %v", err)
	}

	// TODO: delete sub collection records

	return nil
}

func (s *FirestoreDocService) Query(key sdk.Key, subcollection string, expressions []sdk.QueryExpression, limit int, pagingToken map[string]string) (*sdk.QueryResult, error) {
	err := document.ValidateCollection(key.Collection, subcollection)
	if err != nil {
		return nil, err
	}

	err = document.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	queryResult := &sdk.QueryResult{
		Data: make([]map[string]interface{}, 0),
	}

	collRef := s.client.Collection(key.Collection)

	// Fast path lookup document
	if key.Id != "" && subcollection == "" && len(expressions) == 0 {
		value, err := collRef.Doc(key.Id).Get(s.context)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return queryResult, nil

			} else {
				return nil, fmt.Errorf("error retrieving value: %v", err)
			}
		}

		queryResult.Data = append(queryResult.Data, value.Data())
		return queryResult, nil
	}

	var query firestore.Query

	// Select correct root collection to perform query on
	if key.Id != "" {
		subcollRef := collRef.Doc(key.Id).Collection(subcollection)
		query = subcollRef.Offset(0)

	} else {
		if subcollection != "" {
			query = s.client.CollectionGroup(subcollection).Offset(0)
		} else {
			query = collRef.Offset(0)
		}
	}

	var orderByAttrib string

	for _, exp := range expressions {
		expOperand := exp.Operand
		if exp.Operator == "startsWith" {
			endRangeValue := document.GetEndRangeValue(exp.Value)
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
		queryResult.Data = append(queryResult.Data, docSnp.Data())

		// If query limit configured determine continue tokens
		if limit > 0 && len(queryResult.Data) == limit {
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
