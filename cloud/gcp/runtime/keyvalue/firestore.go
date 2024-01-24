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

	"github.com/nitrictech/nitric/core/pkg/decorators/keyvalue"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/keyvalue/v1"

	"google.golang.org/grpc/codes"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"golang.org/x/oauth2/google"
	"google.golang.org/grpc/status"
)

type FirestoreDocService struct {
	client *firestore.Client
}

var _ v1.KeyValueServer = &FirestoreDocService{}

func (s *FirestoreDocService) Get(ctx context.Context, req *v1.KeyValueGetRequest) (*v1.KeyValueGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Get")

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

	return &v1.KeyValueGetResponse{
		Value: &v1.Value{
			Ref:     req.Ref,
			Content: documentContent,
		},
	}, nil
}

func (s *FirestoreDocService) Set(ctx context.Context, req *v1.KeyValueSetRequest) (*v1.KeyValueSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Set")

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

	return &v1.KeyValueSetResponse{}, nil
}

func (s *FirestoreDocService) Delete(ctx context.Context, req *v1.KeyValueDeleteRequest) (*v1.KeyValueDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("FirestoreDocService.Delete")

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

	return &v1.KeyValueDeleteResponse{}, nil
}

func New() (v1.KeyValueServer, error) {
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

func NewWithClient(client *firestore.Client) (v1.KeyValueServer, error) {
	return &FirestoreDocService{
		client: client,
	}, nil
}

func (s *FirestoreDocService) getDocRef(ref *v1.ValueRef) *firestore.DocumentRef {
	return s.client.Collection(ref.Store).Doc(ref.Key)
}
