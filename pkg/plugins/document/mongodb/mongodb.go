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

package mongodb_service

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/nitric-dev/membrane/pkg/plugins/document"
	"github.com/nitric-dev/membrane/pkg/plugins/errors"
	"github.com/nitric-dev/membrane/pkg/plugins/errors/codes"
	"github.com/nitric-dev/membrane/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	client  *mongo.Client
	db      *mongo.Database
	context context.Context
	document.UnimplementedDocumentPlugin
}

func (s *MongoDocService) Get(key *document.Key) (*document.Document, error) {
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

	err := col.FindOne(s.context, docRef, opts).Decode(&value)

	if err != nil {
		var code = codes.Internal
		if err == mongo.ErrNoDocuments {
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

func (s *MongoDocService) Set(key *document.Key, value map[string]interface{}) error {
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

	update := bson.D{{"$set", value}}

	_, err := coll.UpdateOne(s.context, filter, update, opts)

	if err != nil {
		return newErr(
			codes.Internal,
			"error updating value",
			err,
		)
	}

	// add references
	if key.Collection.Parent != nil {
		err := s.updateChildReferences(key, coll.Name(), "$addToSet")

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

func (s *MongoDocService) Delete(key *document.Key) error {
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
	if err := coll.FindOneAndDelete(s.context, filter, opts).Decode(&deletedDocument); err != nil {
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
			_, err := childCol.DeleteMany(s.context, bson.D{{parentKeyAttr, key.Id}})

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
		err := s.updateChildReferences(key, coll.Name(), "$pull")

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

func (s *MongoDocService) getCursor(collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (cursor *mongo.Cursor, orderBy string, err error) {
	coll := s.getCollection(&document.Key{Collection: collection})

	query := bson.M{}

	opts := options.Find()

	opts.SetProjection(bson.M{childrenAttr: 0})

	if limit > 0 {
		opts.SetLimit(int64(limit))

		if len(pagingToken) > 0 {
			opts.SetSort(bson.D{{primaryKeyAttr, 1}})

			if tokens, ok := pagingToken["pagingTokens"]; ok {
				var vals []interface{}
				for _, v := range strings.Split(tokens, "|") {
					vals = append(vals, v)
				}

				query[primaryKeyAttr] = bson.D{{"$gt", vals[len(vals)-1]}}
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
				{s.getOperator(">="), expVal},
				{s.getOperator("<"), endRangeValue},
			}

			query[expOperand] = startsWith

		} else {
			query[expOperand] = bson.D{
				{s.getOperator(exp.Operator), exp.Value},
			}
		}

		if exp.Operator != "==" && limit > 0 && orderBy == "" {
			opts.SetSort(bson.D{{expOperand, 1}})
			orderBy = expOperand
		}
	}

	cursor, err = coll.Find(s.context, query, opts)

	return
}

func (s *MongoDocService) Query(collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
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
			fmt.Errorf("collection: %v, expressions%v", colErr, expErr),
		)
	}

	queryResult := &document.QueryResult{
		Documents: make([]document.Document, 0),
	}

	cursor, orderBy, err := s.getCursor(collection, expressions, limit, pagingToken)

	if err != nil {
		return nil, newErr(
			codes.InvalidArgument,
			"error creating mongo find",
			err,
		)
	}

	defer cursor.Close(s.context)
	for cursor.Next(s.context) {
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

func (s *MongoDocService) QueryStream(collection *document.Collection, expressions []document.QueryExpression, limit int) document.DocumentIterator {
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
				fmt.Errorf("collection error:%v, expression error: %v", colErr, expErr),
			)
		}
	}

	cursor, _, cursorErr := s.getCursor(collection, expressions, limit, nil)

	return func() (*document.Document, error) {
		if cursorErr != nil {
			return nil, cursorErr
		}

		if cursor.Next(s.context) {
			// return the next document
			return mongoDocToDocument(collection, cursor)
		} else {
			// there was an error
			// Close the cursor
			cursor.Close(s.context)

			// Examine the cursors error to see if it's exhausted
			if cursor.Err() != nil {
				return nil, cursor.Err()
			} else {
				return nil, io.EOF
			}
		}
	}
}

func mongoDocToDocument(coll *document.Collection, cursor *mongo.Cursor) (*document.Document, error) {
	var docSnap map[string]interface{}

	err := cursor.Decode(&docSnap)

	if err != nil {
		return nil, fmt.Errorf("error decoding mongo document")
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
		return nil, fmt.Errorf("mongodb error creating client: %v", clientError)
	}

	connectError := client.Connect(ctx)

	if connectError != nil {
		return nil, fmt.Errorf("mongodb unable to initialize connection: %v", connectError)
	}

	db := client.Database(database)

	return &MongoDocService{
		client:  client,
		db:      db,
		context: context.Background(),
	}, nil
}

func NewWithClient(client *mongo.Client, database string, ctx context.Context) document.DocumentService {
	db := client.Database(database)

	return &MongoDocService{
		client:  client,
		db:      db,
		context: ctx,
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

func (s *MongoDocService) updateChildReferences(key *document.Key, subCollectionName string, action string) error {
	parentColl := s.getCollection(key.Collection.Parent)
	filter := bson.M{primaryKeyAttr: key.Collection.Parent.Id}
	referenceMeta := bson.M{childrenAttr: subCollectionName}
	update := bson.D{{action, referenceMeta}}

	opts := options.Update().SetUpsert(true)
	_, err := parentColl.UpdateOne(s.context, filter, update, opts)

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
