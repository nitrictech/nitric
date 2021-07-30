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

package boltdb_service

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/nitric-dev/membrane/pkg/plugins/document"
	"github.com/nitric-dev/membrane/pkg/utils"

	"github.com/Knetic/govaluate"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"go.etcd.io/bbolt"
)

const DEFAULT_DIR = "nitric/collections/"

const skipTokenName = "skip"
const idName = "Id"
const partionKeyName = "ParitionKey"
const sortKeyName = "SortKey"

type BoltDocService struct {
	document.UnimplementedDocumentPlugin
	dbDir string
}

type BoltDoc struct {
	Id          string `storm:"id"`
	ParitionKey string `storm:"index"`
	SortKey     string `storm:"index"`
	Value       map[string]interface{}
}

func (d BoltDoc) String() string {
	return fmt.Sprintf("BoltDoc{Id: %v ParitionKey: %v SortKey: %v Value: %v}\n", d.Id, d.ParitionKey, d.SortKey, d.Value)
}

func (s *BoltDocService) Get(key *document.Key) (*document.Document, error) {
	err := document.ValidateKey(key)
	if err != nil {
		return nil, err
	}

	db, err := s.createdDb(*key.Collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	doc := createDoc(key)

	err = db.One(idName, doc.Id, &doc)

	if err != nil {
		return nil, err
	}

	return toSdkDoc(key.Collection, doc), nil
}

func (s *BoltDocService) Set(key *document.Key, content map[string]interface{}) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	if content == nil {
		return fmt.Errorf("provide non-nil content")
	}

	db, err := s.createdDb(*key.Collection)
	if err != nil {
		return err
	}
	defer db.Close()

	doc := createDoc(key)
	doc.Value = content

	return db.Save(&doc)
}

func (s *BoltDocService) Delete(key *document.Key) error {
	err := document.ValidateKey(key)
	if err != nil {
		return err
	}

	db, err := s.createdDb(*key.Collection)
	if err != nil {
		return err
	}
	defer db.Close()

	doc := createDoc(key)

	err = db.DeleteStruct(&doc)
	if err != nil {
		return err
	}

	// Delete sub collection documents
	if key.Collection.Parent == nil {
		childDocs, err := fetchChildDocs(key, db)
		if err != nil {
			return err
		}

		for _, childDoc := range childDocs {
			err = db.DeleteStruct(&childDoc)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *BoltDocService) Query(collection *document.Collection, expressions []document.QueryExpression, limit int, pagingToken map[string]string) (*document.QueryResult, error) {
	err := document.ValidateQueryCollection(collection)
	if err != nil {
		return nil, err
	}

	err = document.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	db, err := s.createdDb(*collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Build up chain of expression matchers
	matchers := []q.Matcher{}

	// Apply collection/sub-collection filters
	parentKey := collection.Parent

	if parentKey == nil {
		matchers = append(matchers, q.Eq(sortKeyName, collection.Name+"#"))

	} else {
		if parentKey.Id != "" {
			matchers = append(matchers, q.Eq(partionKeyName, parentKey.Id))
		}
		matchers = append(matchers, q.Gte(sortKeyName, collection.Name+"#"))
		matchers = append(matchers, q.Lt(sortKeyName, document.GetEndRangeValue(collection.Name+"#")))
	}

	// Create query object
	matcher := q.And(matchers[:]...)
	query := db.Select(matcher)

	var pagingSkip = 0

	// If fetch limit configured skip past previous reads
	if limit > 0 && len(pagingToken) > 0 {
		if val, found := pagingToken[skipTokenName]; found {
			pagingSkip, err = strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid pagingToken: %v", pagingToken)
			}
			query = query.Skip(pagingSkip)
		}
	}

	// Execute query
	var docs []BoltDoc
	query.Find(&docs)

	// Create values map filter expression, for example : country == 'US' && age < '12'
	expStr := strings.Builder{}
	for _, exp := range expressions {
		// TODO: test typing capabilities of library and rewrite expressions based on value type
		expValue := fmt.Sprintf("%v", exp.Value)

		if expStr.Len() > 0 {
			expStr.WriteString(" && ")
		}
		if exp.Operator == "startsWith" {
			expStr.WriteString(exp.Operand + " >= '" + expValue + "' && ")
			expStr.WriteString(exp.Operand + " < '" + document.GetEndRangeValue(expValue) + "'")
		} else {
			if stringValue, ok := exp.Value.(string); ok {
				expValue = fmt.Sprintf("'%s'", stringValue)
			}
			expStr.WriteString(exp.Operand + " " + exp.Operator + " " + expValue)
		}
	}
	var filterExp *govaluate.EvaluableExpression
	if expStr.Len() > 0 {
		filterExp, err = govaluate.NewEvaluableExpression(expStr.String())
		if err != nil {
			return nil, fmt.Errorf("could not create filter expression: %v, error: %v", expStr.String(), err)
		}
	}

	// Process query results, applying value filter expressions and fetch limit
	documents := make([]document.Document, 0)
	scanCount := 0
	for _, doc := range docs {
		scanCount += 1

		if filterExp != nil {
			include, err := filterExp.Evaluate(doc.Value)
			if err != nil || !(include.(bool)) {
				// TODO: determine if skipping failed evaluations is always appropriate.
				// 	errors are usually a datatype mismatch or a missing key/prop on the doc, which is essentially a failed match.
				// Treat a failed or false eval as a mismatch
				continue
			}
		}
		sdkDoc := toSdkDoc(collection, doc)
		documents = append(documents, *sdkDoc)

		// Break if greater than fetch limit
		if limit > 0 && len(documents) == limit {
			break
		}
	}

	// Provide paging token to skip previous reads
	var resultPagingToken map[string]string
	if limit > 0 && len(documents) == limit {
		resultPagingToken = make(map[string]string)
		resultPagingToken[skipTokenName] = fmt.Sprintf("%v", pagingSkip+scanCount)
	}

	return &document.QueryResult{
		Documents:   documents,
		PagingToken: resultPagingToken,
	}, nil
}

// New - Create a new dev KV plugin
func New() (*BoltDocService, error) {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", DEFAULT_DIR)

	// Check whether file exists
	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		// Make directory if not present
		err := os.MkdirAll(dbDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	return &BoltDocService{dbDir: dbDir}, nil
}

func (s *BoltDocService) createdDb(coll document.Collection) (*storm.DB, error) {
	for coll.Parent != nil {
		coll = *coll.Parent.Collection
	}

	dbPath := s.dbDir + strings.ToLower(coll.Name) + ".db"

	options := storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second})
	db, err := storm.Open(dbPath, options)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDoc(key *document.Key) BoltDoc {

	parentKey := key.Collection.Parent

	// Top Level Collection
	if parentKey == nil {
		return BoltDoc{
			Id:          key.Id,
			ParitionKey: key.Id,
			SortKey:     key.Collection.Name + "#",
		}

	} else {
		return BoltDoc{
			Id:          parentKey.Id + "_" + key.Id,
			ParitionKey: parentKey.Id,
			SortKey:     key.Collection.Name + "#" + key.Id,
		}
	}
}

func toSdkDoc(col *document.Collection, doc BoltDoc) *document.Document {
	keys := strings.Split(doc.Id, "_")

	// Translate the boltdb Id into a nitric document key Id
	var id string
	var c *document.Collection
	if len(keys) > 1 {
		// sub document
		id = keys[len(keys)-1]
		c = &document.Collection{
			Name: col.Name,
			Parent: &document.Key{
				Collection: col.Parent.Collection,
				Id:         keys[0],
			},
		}
	} else {
		id = doc.Id
		c = col
	}

	return &document.Document{
		Content: doc.Value,
		Key: &document.Key{
			Collection: c,
			Id:         id,
		},
	}
}

func fetchChildDocs(key *document.Key, db *storm.DB) ([]BoltDoc, error) {
	var childDocs []BoltDoc

	err := db.Find(partionKeyName, key.Id, &childDocs)
	if err != nil {
		if err.Error() == "not found" {
			return childDocs, nil
		} else {
			return nil, err
		}
	}

	return childDocs, nil
}
