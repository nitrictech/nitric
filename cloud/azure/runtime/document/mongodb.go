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
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/nitrictech/nitric/cloud/azure/runtime/env"
	document "github.com/nitrictech/nitric/core/pkg/decorators/documents"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	documentpb "github.com/nitrictech/nitric/core/pkg/proto/documents/v1"
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

// Mapping to mongo operators, startsWith will be handled within the function
var mongoOperatorMap = map[string]string{
	"<":  "$lt",
	"<=": "$lte",
	"==": "$eq",
	"!=": "$ne",
	">=": "$gte",
	">":  "$gt",
}

type MongoDocService struct {
	client *mongo.Client
	db     *mongo.Database
}

var _ documentpb.DocumentsServer = &MongoDocService{}

// Get a document from the mongo collection
func (s *MongoDocService) Get(ctx context.Context, req *documentpb.DocumentGetRequest) (*documentpb.DocumentGetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.Get")

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	col := s.getCollection(req.Key)
	docRef := bson.M{primaryKeyAttr: req.Key.Id}

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

	return &documentpb.DocumentGetResponse{
		Document: &documentpb.Document{
			Key:     req.Key,
			Content: content,
		},
	}, nil
}

// Set a document in the mongo collection
func (s *MongoDocService) Set(ctx context.Context, req *documentpb.DocumentSetRequest) (*documentpb.DocumentSetResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.Set")

	if err := document.ValidateKey(req.Key); err != nil {
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

	filter := bson.M{primaryKeyAttr: req.Key.Id}

	update := bson.D{{Key: "$set", Value: value}}

	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error updating value",
			err,
		)
	}

	// add references
	if req.Key.Collection.Parent != nil {
		err := s.updateChildReferences(ctx, req.Key, coll.Name(), "$addToSet")
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error updating child references",
				err,
			)
		}
	}

	return &documentpb.DocumentSetResponse{}, nil
}

// Delete a document from the mongo collection
func (s *MongoDocService) Delete(ctx context.Context, req *documentpb.DocumentDeleteRequest) (*documentpb.DocumentDeleteResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.Delete")

	if err := document.ValidateKey(req.Key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	coll := s.getCollection(req.Key)

	filter := bson.M{primaryKeyAttr: req.Key.Id}

	opts := options.FindOneAndDelete().SetProjection(bson.M{childrenAttr: 1, primaryKeyAttr: 0})

	var deletedDocument map[string]interface{}

	// Delete document
	if err := coll.FindOneAndDelete(ctx, filter, opts).Decode(&deletedDocument); err != nil {
		return nil, newErr(
			codes.Internal,
			"error deleting document",
			err,
		)
	}

	// Delete all the child collection documents
	if deletedDocument[childrenAttr] != nil {
		children := deletedDocument[childrenAttr].(primitive.A)

		for _, v := range children {
			colName := v.(string)
			childCol := s.db.Collection(colName)
			_, err := childCol.DeleteMany(ctx, bson.D{{Key: parentKeyAttr, Value: req.Key.Id}})
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"error deleting child collection document(s)",
					err,
				)
			}
		}
	}

	// clean references if none left
	if req.Key.Collection.Parent != nil {
		err := s.updateChildReferences(ctx, req.Key, coll.Name(), "$pull")
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error removing child references",
				err,
			)
		}
	}

	return &documentpb.DocumentDeleteResponse{}, nil
}

func (s *MongoDocService) getCursor(ctx context.Context, collection *documentpb.Collection, expressions []*documentpb.Expression, limit int32, pagingToken map[string]string) (cursor *mongo.Cursor, orderBy string, err error) {
	coll := s.getCollection(&documentpb.Key{Collection: collection})

	query := bson.M{}

	opts := options.Find()

	opts.SetProjection(bson.M{childrenAttr: 0})

	if limit > 0 {
		opts.SetLimit(int64(limit))

		if len(pagingToken) > 0 {
			opts.SetSort(bson.D{{Key: primaryKeyAttr, Value: 1}})

			if tokens, ok := pagingToken["pagingTokens"]; ok {
				var vals []interface{}
				for _, v := range strings.Split(tokens, "|") {
					vals = append(vals, v)
				}

				query[primaryKeyAttr] = bson.D{{Key: "$gt", Value: vals[len(vals)-1]}}
			}
		}
	}

	if collection.Parent != nil && collection.Parent.Id != "" {
		query[parentKeyAttr] = collection.Parent.Id
	}

	for _, exp := range expressions {
		expOperand := exp.Operand
		if exp.Operator == "startsWith" {
			expVal := fmt.Sprintf("%v", exp.Value)
			endRangeValue := document.GetEndRangeValue(expVal)

			startsWith := bson.D{
				{Key: s.getOperator(">="), Value: expVal},
				{Key: s.getOperator("<"), Value: endRangeValue},
			}

			query[expOperand] = startsWith
		} else {
			query[expOperand] = bson.D{
				{Key: s.getOperator(exp.Operator), Value: exp.Value},
			}
		}

		if exp.Operator != "==" && limit > 0 && orderBy == "" {
			opts.SetSort(bson.D{{Key: expOperand, Value: 1}})
			orderBy = expOperand
		}
	}

	cursor, err = coll.Find(ctx, query, opts)

	return
}

// Query documents from the mongo collection with pagination
func (s *MongoDocService) Query(ctx context.Context, req *documentpb.DocumentQueryRequest) (*documentpb.DocumentQueryResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.Query")

	if colErr, expErr := document.ValidateQueryCollection(req.Collection), document.ValidateExpressions(req.Expressions); colErr != nil || expErr != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid arguments",
			fmt.Errorf("collection: %w, expressions %w", colErr, expErr),
		)
	}

	queryResult := &documentpb.DocumentQueryResponse{
		Documents: make([]*documentpb.Document, 0),
	}

	cursor, orderBy, err := s.getCursor(ctx, req.Collection, req.Expressions, req.Limit, req.PagingToken)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"error creating mongo find",
			err,
		)
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		sdkDoc, err := mongoDocToDocument(req.Collection, cursor)
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error decoding mongo document",
				err,
			)
		}

		queryResult.Documents = append(queryResult.Documents, sdkDoc)

		// If query limit configured determine continue tokens
		if req.Limit > 0 && len(queryResult.Documents) == int(req.Limit) {
			tokens := ""
			if orderBy != "" {
				tokens = fmt.Sprintf("%v", sdkDoc.Content.AsMap()[orderBy]) + "|"
			}
			tokens += sdkDoc.Key.Id

			queryResult.PagingToken = map[string]string{
				"pagingTokens": tokens,
			}
		}
	}

	return queryResult, nil
}

// QueryStream queries documents from the mongo collection as a stream
func (s *MongoDocService) QueryStream(req *documentpb.DocumentQueryStreamRequest, stream documentpb.Documents_QueryStreamServer) error {
	newErr := grpc_errors.ErrorsWithScope("MongoDocService.QueryStream")

	colErr := document.ValidateQueryCollection(req.Collection)
	expErr := document.ValidateExpressions(req.Expressions)

	if colErr != nil || expErr != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid arguments",
			fmt.Errorf("collection error: %w, expression error: %w", colErr, expErr),
		)
	}

	cursor, _, cursorErr := s.getCursor(stream.Context(), req.Collection, req.Expressions, req.Limit, nil)

	if cursorErr != nil {
		return cursorErr
	}

	for cursor.Next(stream.Context()) {
		// return the next document
		doc, err := mongoDocToDocument(req.Collection, cursor)
		if err != nil {
			return newErr(
				codes.Internal,
				"error decoding mongo document",
				err,
			)
		}

		err = stream.Send(&documentpb.DocumentQueryStreamResponse{
			Document: doc,
		})

		if err != nil {
			return err
		}
	}

	if err := cursor.Close(stream.Context()); err != nil {
		return newErr(
			codes.Internal,
			"mongo cursor close error",
			cursor.Err(),
		)
	}

	if cursor.Err() != nil {
		return newErr(
			codes.Internal,
			"mongo cursor error",
			cursor.Err(),
		)
	} else {
		return nil
	}
}

func mongoDocToDocument(coll *documentpb.Collection, cursor *mongo.Cursor) (*documentpb.Document, error) {
	var docSnap map[string]interface{}

	if err := cursor.Decode(&docSnap); err != nil {
		return nil, err
	}

	id := docSnap[primaryKeyAttr].(string)

	// remove id from content
	delete(docSnap, primaryKeyAttr)

	contentStruct, err := structpb.NewStruct(docSnap)
	if err != nil {
		return nil, err
	}

	sdkDoc := &documentpb.Document{
		Content: contentStruct,
		Key: &documentpb.Key{
			Collection: coll,
			Id:         id,
		},
	}

	if docSnap[parentKeyAttr] != nil {
		parentId := docSnap[parentKeyAttr].(string)

		sdkDoc.Key.Collection = &documentpb.Collection{
			Name: coll.Name,
			Parent: &documentpb.Key{
				Collection: coll.Parent.Collection,
				Id:         parentId,
			},
		}

		delete(docSnap, parentKeyAttr)
	}

	return sdkDoc, nil
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

func mapKeys(key *documentpb.Key, source map[string]interface{}) map[string]interface{} {
	// Copy map
	newMap := make(map[string]interface{})

	for key, value := range source {
		newMap[key] = value
	}

	parentKey := key.Collection.Parent

	newMap[primaryKeyAttr] = key.Id

	if parentKey != nil {
		newMap[parentKeyAttr] = parentKey.Id
	}

	return newMap
}

func (s *MongoDocService) updateChildReferences(ctx context.Context, key *documentpb.Key, subCollectionName string, action string) error {
	parentColl := s.getCollection(key.Collection.Parent)
	filter := bson.M{primaryKeyAttr: key.Collection.Parent.Id}
	referenceMeta := bson.M{childrenAttr: subCollectionName}
	update := bson.D{{Key: action, Value: referenceMeta}}

	opts := options.Update().SetUpsert(true)
	_, err := parentColl.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoDocService) getCollection(key *documentpb.Key) *mongo.Collection {
	collectionNames := []string{}
	parentKey := key.Collection.Parent

	for parentKey != nil {
		collectionNames = append(collectionNames, parentKey.Collection.Name)
		parentKey = parentKey.Collection.Parent
	}

	collectionNames = append(collectionNames, key.Collection.Name)

	return s.db.Collection(strings.Join(collectionNames, "."))
}

func (s *MongoDocService) getOperator(operator string) string {
	return mongoOperatorMap[operator]
}
