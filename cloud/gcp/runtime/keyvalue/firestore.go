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

package keyvalue

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nitrictech/nitric/core/pkg/decorators/keyvalue"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/kvstore/v1"

	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/status"
)

type FirestoreDocService struct {
	client *firestore.Client
}

var _ v1.KvStoreServer = &FirestoreDocService{}

func (s *FirestoreDocService) GetKey(ctx context.Context, req *v1.KvStoreGetRequest) (*v1.KvStoreGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.GetKey")

	if err := keyvalue.ValidateValueRef(req.Ref); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	doc := s.getDocRef(req.Ref)

	value, err := doc.Get(ctx)
	if err != nil {
		switch status.Code(err) {
		case codes.NotFound:
			return nil, newErr(
				codes.NotFound,
				fmt.Sprintf("key %s not found in store %s", req.Ref.Key, req.Ref.Store),
				err,
			)
		case codes.PermissionDenied:
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this key value store?",
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

	return &v1.KvStoreGetResponse{
		Value: &v1.Value{
			Ref:     req.Ref,
			Content: documentContent,
		},
	}, nil
}

func (s *FirestoreDocService) SetKey(ctx context.Context, req *v1.KvStoreSetRequest) (*v1.KvStoreSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.SetKey")

	if err := keyvalue.ValidateValueRef(req.Ref); err != nil {
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

	doc := s.getDocRef(req.Ref)

	if _, err := doc.Set(ctx, req.Content.AsMap()); err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this key value store?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"error updating value",
			err,
		)
	}

	return &v1.KvStoreSetResponse{}, nil
}

func (s *FirestoreDocService) DeleteKey(ctx context.Context, req *v1.KvStoreDeleteRequest) (*v1.KvStoreDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.DeleteKey")

	if err := keyvalue.ValidateValueRef(req.Ref); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	doc := s.getDocRef(req.Ref)

	// Delete document
	if _, err := doc.Delete(ctx); err != nil {
		if status.Code(err) == codes.PermissionDenied {
			return nil, newErr(
				codes.PermissionDenied,
				"permission denied, have you requested access to this key value store?",
				err,
			)
		}

		return nil, newErr(
			codes.Internal,
			"error deleting value",
			err,
		)
	}

	return &v1.KvStoreDeleteResponse{}, nil
}

func (s *FirestoreDocService) Keys(req *v1.KvStoreKeysRequest, stream v1.KvStore_KeysServer) error {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Keys")
	storeName := req.GetStore().GetName()

	if storeName == "" {
		return newErr(
			codes.InvalidArgument,
			"store name is required",
			nil,
		)
	}

	iter := s.getCollectionRef(storeName).DocumentRefs(stream.Context())

	for {
		doc, err := iter.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return newErr(
				codes.Internal,
				"error iterating over firestore collection",
				err,
			)
		}

		// Range queries don't appear to be supported when querying based on document ID.
		// e.g. Where(firestore.DocumentID, "<=", req.Prefix)
		// since prefix is a string not a DocumentRef.
		// Instead we filter the results as they're returned
		if !strings.HasPrefix(doc.ID, req.Prefix) {
			continue
		}

		if err := stream.Send(&v1.KvStoreKeysResponse{
			Key: doc.ID,
		}); err != nil {
			return newErr(
				codes.Internal,
				"failed to send response",
				err,
			)
		}
	}

	return nil
}

func New() (v1.KvStoreServer, error) {
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

func NewWithClient(client *firestore.Client) (v1.KvStoreServer, error) {
	return &FirestoreDocService{
		client: client,
	}, nil
}

func (s *FirestoreDocService) getDocRef(ref *v1.ValueRef) *firestore.DocumentRef {
	return s.getCollectionRef(ref.Store).Doc(ref.Key)
}

func (s *FirestoreDocService) getCollectionRef(store string) *firestore.CollectionRef {
	return s.client.Collection(store)
}
