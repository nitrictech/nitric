package firestore_plugin

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

type FirestorePlugin struct {
	client *firestore.Client
	sdk.UnimplementedDocumentsPlugin
}

func (s *FirestorePlugin) CreateDocument(collection string, key string, document map[string]interface{}) error {
	// Create a new document is firestore
	if key == "" {
		return fmt.Errorf("Key autogeneration unimplemented, please provide non-blank key")
	}

	_, error := s.client.Collection(collection).Doc(key).Create(context.TODO(), document)

	if error != nil {
		return fmt.Errorf("Error creating new document: %v", error)
	}

	return nil
}

func (s *FirestorePlugin) GetDocument(collection string, key string) (map[string]interface{}, error) {
	document, error := s.client.Collection(collection).Doc(key).Get(context.TODO())

	if error != nil {
		return nil, fmt.Errorf("Error retrieving document: %v", error)
	}

	return document.Data(), nil
}

func (s *FirestorePlugin) UpdateDocument(collection string, key string, document map[string]interface{}) error {
	docRef := s.client.Collection(collection).Doc(key)

	_, err := docRef.Get(context.TODO())

	if err == nil {
		_, err := s.client.Collection(collection).Doc(key).Set(context.TODO(), document)

		if err != nil {
			return fmt.Errorf("Error updating document: %v", err)
		}
	} else {
		return fmt.Errorf("Document does not exist: %v", err)
	}

	return nil
}

func (s *FirestorePlugin) DeleteDocument(collection string, key string) error {
	_, error := s.client.Collection(collection).Doc(key).Delete(context.TODO())

	if error != nil {
		return fmt.Errorf("Error deleting document: %v", error)
	}

	return nil
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

	return &FirestorePlugin{
		client: client,
	}, nil
}

func NewWithClient(client *firestore.Client) (sdk.DocumentService, error) {
	return &FirestorePlugin{
		client: client,
	}, nil
}
