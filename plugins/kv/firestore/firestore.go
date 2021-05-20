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

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/kv"
	"github.com/nitric-dev/membrane/sdk"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
)

type FirestoreKVService struct {
	client *firestore.Client
	sdk.UnimplementedKeyValuePlugin
}

func (s *FirestoreKVService) Get(collection string, key map[string]interface{}) (map[string]interface{}, error) {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return nil, err
	}

	keyValue, err := kv.GetKeyValue(key)
	if err != nil {
		return nil, err
	}

	value, err := s.client.Collection(collection).Doc(keyValue).Get(context.TODO())

	if err != nil {
		return nil, fmt.Errorf("Error retrieving value: %v", err)
	}

	return value.Data(), nil
}

func (s *FirestoreKVService) Put(collection string, key map[string]interface{}, value map[string]interface{}) error {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return err
	}

	keyValue, err := kv.GetKeyValue(key)
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	_, err = s.client.Collection(collection).Doc(keyValue).Set(context.TODO(), value)

	if err != nil {
		return fmt.Errorf("Error updating value: %v", err)
	}

	return nil
}

func (s *FirestoreKVService) Delete(collection string, key map[string]interface{}) error {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return err
	}
	keyValue, err := kv.GetKeyValue(key)
	if err != nil {
		return err
	}

	_, err = s.client.Collection(collection).Doc(keyValue).Delete(context.TODO())

	if err != nil {
		return fmt.Errorf("Error deleting value: %v", err)
	}

	return nil
}

func (s *FirestoreKVService) Query(collection string, expressions []sdk.QueryExpression, limit int) ([]map[string]interface{}, error) {
	err := kv.ValidateCollection(collection)
	if err != nil {
		return nil, err
	}
	err = kv.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	query := s.client.Collection(collection).Select("Value")

	for _, exp := range expressions {
		if exp.Operator == "startsWith" {
			endRangeValue := kv.GetEndRangeValue(exp.Value)
			query = query.Where(exp.Operand, ">=", exp.Value).Where(exp.Operand, "<", endRangeValue)

		} else {
			query = query.Where(exp.Operand, exp.Operator, exp.Value)
		}
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	itr := query.Documents(context.TODO())

	results := make([]map[string]interface{}, 0)

	for {
		docSnp, err := itr.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error querying value: %v", err)
		}
		results = append(results, docSnp.Data())
	}

	return results, nil
}

func New() (sdk.KeyValueService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials err: %v", credentialsError)
	}

	client, clientError := firestore.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("firestore client err: %v", clientError)
	}

	return &FirestoreKVService{
		client: client,
	}, nil
}

func NewWithClient(client *firestore.Client) (sdk.KeyValueService, error) {
	return &FirestoreKVService{
		client: client,
	}, nil
}
