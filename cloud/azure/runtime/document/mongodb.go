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
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/nitrictech/nitric/cloud/azure/runtime/env"
	"github.com/nitrictech/nitric/core/pkg/decorators/keyvalue"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	keyvaluepb "github.com/nitrictech/nitric/core/pkg/proto/keyvalue/v1"
)

const (
	// environment variables
	mongoDBConnectionStringEnvVarName = "MONGODB_CONNECTION_STRING"
	mongoDBDatabaseEnvVarName         = "MONGODB_DATABASE"
	mongoDBSetDirectEnvVarName        = "MONGODB_DIRECT"

	primaryKeyAttr = "_id"
	parentKeyAttr  = "_parent_id"
	childrenAttr   = "_child_colls"
)

type MongoDocService struct {
	client *mongo.Client
	db     *mongo.Database
}

var _ keyvaluepb.KeyValueServer = &MongoDocService{}

// Get a document from the mongo collection
func (s *MongoDocService) Get(ctx context.Context, req *keyvaluepb.KeyValueGetRequest) (*keyvaluepb.KeyValueGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.Get")

	if err := keyvalue.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	col := s.getCollection(req.Key)
	docRef := bson.M{primaryKeyAttr: req.Key.Key}

	var value map[string]interface{}

	opts := options.FindOne()

	// Remove meta data IDs and child collections
	opts.SetProjection(bson.M{primaryKeyAttr: 0, parentKeyAttr: 0, childrenAttr: 0})

	err := col.FindOne(ctx, docRef, opts).Decode(&value)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, newErr(
				codes.NotFound,
				"document not found",
				err,
			)
		}

		return nil, newErr(
			codes.Unknown,
			"error getting document",
			err,
		)
	}

	content, err := structpb.NewStruct(value)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error creating struct from document content",
			err,
		)
	}

	return &keyvaluepb.KeyValueGetResponse{
		Value: &keyvaluepb.Value{
			Key:     req.Key,
			Content: content,
		},
	}, nil
}

// Set a document in the mongo collection
func (s *MongoDocService) Set(ctx context.Context, req *keyvaluepb.KeyValueSetRequest) (*keyvaluepb.KeyValueSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.Set")

	if err := keyvalue.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if req.Content == nil {
		return nil, newErr(
			codes.InvalidArgument,
			"document value cannot be empty",
			nil,
		)
	}

	coll := s.getCollection(req.Key)

	mapContent := req.Content.AsMap()

	value := mapKeys(req.Key, mapContent)

	opts := options.Update().SetUpsert(true)

	filter := bson.M{primaryKeyAttr: req.Key.Key}

	update := bson.D{{Key: "$set", Value: value}}

	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error updating value",
			err,
		)
	}

	return &keyvaluepb.KeyValueSetResponse{}, nil
}

// Delete a document from the mongo collection
func (s *MongoDocService) Delete(ctx context.Context, req *keyvaluepb.KeyValueDeleteRequest) (*keyvaluepb.KeyValueDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.Delete")

	if err := keyvalue.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	coll := s.getCollection(req.Key)

	filter := bson.M{primaryKeyAttr: req.Key.Key}

	opts := options.FindOneAndDelete().SetProjection(bson.M{childrenAttr: 1, primaryKeyAttr: 0})

	var deletedKeyValue map[string]interface{}

	// Delete document
	if err := coll.FindOneAndDelete(ctx, filter, opts).Decode(&deletedKeyValue); err != nil {
		return nil, newErr(
			codes.Internal,
			"error deleting document",
			err,
		)
	}

	return &keyvaluepb.KeyValueDeleteResponse{}, nil
}

func New() (*MongoDocService, error) {
	mongoDBConnectionString := env.MONGODB_CONNECTION_STRING.String()
	if mongoDBConnectionString == "" {
		return nil, fmt.Errorf("MongoDB missing environment variable: %v", mongoDBConnectionStringEnvVarName)
	}

	database := env.MONGODB_DATABASE.String()
	if database == "" {
		return nil, fmt.Errorf("MongoDB missing environment variable: %v", mongoDBDatabaseEnvVarName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongoDBSetDirect, err := env.MONGODB_DIRECT.Bool()
	if err != nil {
		return nil, err
	}

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(mongoDBSetDirect)

	client, clientError := mongo.NewClient(clientOptions)

	if clientError != nil {
		return nil, fmt.Errorf("mongodb error creating client: %w", clientError)
	}

	connectError := client.Connect(ctx)

	if connectError != nil {
		return nil, fmt.Errorf("mongodb unable to initialize connection: %w", connectError)
	}

	db := client.Database(database)

	return &MongoDocService{
		client: client,
		db:     db,
	}, nil
}

func NewWithClient(client *mongo.Client, database string) *MongoDocService {
	db := client.Database(database)

	return &MongoDocService{
		client: client,
		db:     db,
	}
}

func mapKeys(key *keyvaluepb.Key, source map[string]interface{}) map[string]interface{} {
	// Copy map
	newMap := make(map[string]interface{})

	for key, value := range source {
		newMap[key] = value
	}

	newMap[primaryKeyAttr] = key.Key

	return newMap
}

func (s *MongoDocService) getCollection(key *keyvaluepb.Key) *mongo.Collection {
	return s.db.Collection(key.Store)
}
