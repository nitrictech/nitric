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
	"github.com/nitric-dev/membrane/pkg/sdk"
	"golang.org/x/oauth2/google"
)

type FirestoreKVService struct {
	client *firestore.Client
	sdk.UnimplementedKeyValuePlugin
}

func (s *FirestoreKVService) Get(collection string, key string) (map[string]interface{}, error) {
	value, error := s.client.Collection(collection).Doc(key).Get(context.TODO())

	if error != nil {
		return nil, fmt.Errorf("Error retrieving value: %v", error)
	}

	return value.Data(), nil
}

func (s *FirestoreKVService) Put(collection string, key string, value map[string]interface{}) error {
	_, err := s.client.Collection(collection).Doc(key).Set(context.TODO(), value)

	if err != nil {
		return fmt.Errorf("Error updating value: %v", err)
	}

	return nil
}

func (s *FirestoreKVService) Delete(collection string, key string) error {
	_, error := s.client.Collection(collection).Doc(key).Delete(context.TODO())

	if error != nil {
		return fmt.Errorf("Error deleting value: %v", error)
	}

	return nil
}

func New() (sdk.KeyValueService, error) {
	ctx := context.Background()

	credentials, credentialsError := google.FindDefaultCredentials(ctx, pubsub.ScopeCloudPlatform)
	if credentialsError != nil {
		return nil, fmt.Errorf("GCP credentials error: %v", credentialsError)
	}

	client, clientError := firestore.NewClient(ctx, credentials.ProjectID)
	if clientError != nil {
		return nil, fmt.Errorf("firestore client error: %v", clientError)
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
