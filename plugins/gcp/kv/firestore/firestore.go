package firestore_service

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"github.com/nitric-dev/membrane/plugins/sdk"
	"golang.org/x/oauth2/google"
)

// type FirestoreClientIface interface {
// 	Collection(string) FirestoreCollectionIface
// }

// type FirestoreCollectionIface interface {
// 	Doc(string) FirestoreDocumentIface
// }

// type FirestoreDocumentIface interface {
// 	Create(context.Context, interface{}) (*firestore.WriteResult, error)
// 	Get(context.Context) (*firestore.DocumentSnapshot, error)
// 	Set(context.Context, interface{}) (*firestore.WriteResult, error)
// }

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
