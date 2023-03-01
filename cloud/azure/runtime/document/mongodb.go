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
	"io"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/nitrictech/nitric/core/pkg/plugins/document"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors"
	"github.com/nitrictech/nitric/core/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/core/pkg/utils"
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
	document.UnimplementedDocumentPlugin
}

func (s *MongoDocService) Get(ctx context.Context, key *document.Key) (*document.Document, error) {
	newErr := errors.ErrorsWithScope(
		"MongoDocService.Get",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	col := s.getCollection(key)
	docRef := bson.M{primaryKeyAttr: key.Id}

	var value map[string]interface{}

	opts := options.FindOne()

	// Remove meta data ids and child colls
	opts.SetProjection(bson.M{primaryKeyAttr: 0, parentKeyAttr: 0, childrenAttr: 0})

	err := col.FindOne(ctx, docRef, opts).Decode(&value)
	if err != nil {
		code := codes.Internal
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, newErr(
				codes.NotFound,
				"document not found",
				err,
			)
		}

		return nil, newErr(
			code,
			"unable to retrieve value",
			err,
		)
	}

	return &document.Document{
		Key:     key,
		Content: value,
	}, nil
}

func (s *MongoDocService) Set(ctx context.Context, key *document.Key, value map[string]interface{}) error {
	newErr := errors.ErrorsWithScope(
		"MongoDocService.Set",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	if value == nil {
		return newErr(
			codes.InvalidArgument,
			"provide non-nil value",
			nil,
		)
	}

	coll := s.getCollection(key)

	value = mapKeys(key, value)

	opts := options.Update().SetUpsert(true)

	filter := bson.M{primaryKeyAttr: key.Id}

	update := bson.D{{Key: "$set", Value: value}}

	_, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return newErr(
			codes.Internal,
			"error updating value",
			err,
		)
	}

	// add references
	if key.Collection.Parent != nil {
		err := s.updateChildReferences(ctx, key, coll.Name(), "$addToSet")
		if err != nil {
			return newErr(
				codes.Internal,
				"error updating child references",
				err,
			)
		}
	}

	return nil
}

func (s *MongoDocService) Delete(ctx context.Context, key *document.Key) error {
	newErr := errors.ErrorsWithScope(
		"MongoDocService.Delete",
		map[string]interface{}{
			"key": key,
		},
	)

	if err := document.ValidateKey(key); err != nil {
		return newErr(
			codes.InvalidArgument,
			"invalid key",
			err,
		)
	}

	coll := s.getCollection(key)

	filter := bson.M{primaryKeyAttr: key.Id}

	opts := options.FindOneAndDelete().SetProjection(bson.M{childrenAttr: 1, primaryKeyAttr: 0})

	var deletedDocument map[string]interface{}

	// Delete document
	if err := coll.FindOneAndDelete(ctx, filter, opts).Decode(&deletedDocument); err != nil {
		return newErr(
			codes.Internal,
			"error deleting value",
			err,
		)
	}

	// Delete all the child collection documents
	if deletedDocument[childrenAttr] != nil {
		children := deletedDocument[childrenAttr].(primitive.A)

		for _, v := range children {
			colName := v.(string)
			childCol := s.db.Collection(colName)
			_, err := childCol.DeleteMany(ctx, bson.D{{Key: parentKeyAttr, Value: key.Id}})
			if err != nil {
				return newErr(
					codes.Internal,
					"error deleting child collection value",
					err,
				)
			}
		}
	}

	// clean references if none left
	if key.Collection.Parent != nil {
		err := s.updateChildReferences(ctx, key, coll.Name(), "$pull")
		if err != nil {
			return newErr(
				codes.Internal,
				"error removing child references",
				err,
			)
		}
	}

	return nil
}

func (s *MongoDocService) getCursor(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (cursor *mongo.Cursor, orderBy string, err error) {
	coll := s.getCollection(&document.Key{Collection: collection})

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

func (s *MongoDocService) Query(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
	newErr := errors.ErrorsWithScope(
		"MongoDocService.Query",
		map[string]interface{}{
			"collection": collection,
		},
	)

	if colErr, expErr := document.ValidateQueryCollection(collection), document.ValidateExpressions(expressions); colErr != nil || expErr != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"invalid arguments",
			fmt.Errorf("collection: %w, expressions %w", colErr, expErr),
		)
	}

	queryResult := &document.QueryResult{
		Documents: make([]document.Document, 0),
	}

	cursor, orderBy, err := s.getCursor(ctx, collection, expressions, limit, pagingToken)
	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"error creating mongo find",
			err,
		)
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		sdkDoc, err := mongoDocToDocument(collection, cursor)
		if err != nil {
			return nil, newErr(
				codes.Internal,
				"error decoding mongo document",
				err,
			)
		}

		queryResult.Documents = append(queryResult.Documents, *sdkDoc)

		// If query limit configured determine continue tokens
		if limit > 0 && len(queryResult.Documents) == limit {
			tokens := ""
			if orderBy != "" {
				tokens = fmt.Sprintf("%v", sdkDoc.Content[orderBy]) + "|"
			}
			tokens += sdkDoc.Key.Id

			queryResult.PagingToken = map[string]string{
				"pagingTokens": tokens,
			}
		}
	}

	return queryResult, nil
}

func (s *MongoDocService) QueryStream(ctx context.Context, collection *document.Collection, expressions []document.QueryExpression, limit int) document.DocumentIterator {
	newErr := errors.ErrorsWithScope(
		"MongoDocService.QueryStream",
		map[string]interface{}{
			"collection": collection,
		},
	)

	colErr := document.ValidateQueryCollection(collection)
	expErr := document.ValidateExpressions(expressions)

	if colErr != nil || expErr != nil {
		// Return an error only iterator
		return func() (*document.Document, error) {
			return nil, newErr(
				codes.InvalidArgument,
				"invalid arguments",
				fmt.Errorf("collection error: %w, expression error: %w", colErr, expErr),
			)
		}
	}

	cursor, _, cursorErr := s.getCursor(ctx, collection, expressions, limit, nil)

	return func() (*document.Document, error) {
		if cursorErr != nil {
			return nil, cursorErr
		}

		if cursor.Next(ctx) {
			// return the next document
			doc, err := mongoDocToDocument(collection, cursor)
			if err != nil {
				return nil, newErr(
					codes.Internal,
					"error decoding mongo document",
					err,
				)
			}

			return doc, nil
		} else {
			// there was an error
			// Close the cursor
			cursor.Close(ctx)

			if cursor.Err() != nil {
				return nil, newErr(
					codes.Internal,
					"mongo cursor error",
					cursor.Err(),
				)
			} else {
				return nil, io.EOF
			}
		}
	}
}

func mongoDocToDocument(coll *document.Collection, cursor *mongo.Cursor) (*document.Document, error) {
	var docSnap map[string]interface{}

	if err := cursor.Decode(&docSnap); err != nil {
		return nil, err
	}

	id := docSnap[primaryKeyAttr].(string)

	// remove id from content
	delete(docSnap, primaryKeyAttr)

	sdkDoc := document.Document{
		Content: docSnap,
		Key: &document.Key{
			Collection: coll,
			Id:         id,
		},
	}

	if docSnap[parentKeyAttr] != nil {
		parentId := docSnap[parentKeyAttr].(string)

		sdkDoc.Key.Collection = &document.Collection{
			Name: coll.Name,
			Parent: &document.Key{
				Collection: coll.Parent.Collection,
				Id:         parentId,
			},
		}

		delete(docSnap, parentKeyAttr)
	}

	return &sdkDoc, nil
}

func New() (document.DocumentService, error) {
	mongoDBConnectionString := utils.GetEnv(mongoDBConnectionStringEnvVarName, "")

	if mongoDBConnectionString == "" {
		return nil, fmt.Errorf("MongoDB missing environment variable: %v", mongoDBConnectionStringEnvVarName)
	}

	database := utils.GetEnv(mongoDBDatabaseEnvVarName, "")

	if database == "" {
		return nil, fmt.Errorf("MongoDB missing environment variable: %v", mongoDBDatabaseEnvVarName)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mongoDBSetDirect := utils.GetEnv(mongoDBSetDirectEnvVarName, "true")

	clientOptions := options.Client().ApplyURI(mongoDBConnectionString).SetDirect(mongoDBSetDirect == "true")

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

func NewWithClient(client *mongo.Client, database string) document.DocumentService {
	db := client.Database(database)

	return &MongoDocService{
		client: client,
		db:     db,
	}
}

func mapKeys(key *document.Key, source map[string]interface{}) map[string]interface{} {
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

func (s *MongoDocService) updateChildReferences(ctx context.Context, key *document.Key, subCollectionName string, action string) error {
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

func (s *MongoDocService) getCollection(key *document.Key) *mongo.Collection {
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
