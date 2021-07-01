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

	"github.com/Knetic/govaluate"
	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/nitric-dev/membrane/plugins/document"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
	"go.etcd.io/bbolt"
)

const DEFAULT_DIR = "nitric/collections/"

const skipTokenName = "skip"
const idName = "Id"
const partionKeyName = "ParitionKey"
const sortKeyName = "SortKey"

type BoltDocService struct {
	sdk.UnimplementedDocumentPlugin
	dbDir string
}

type BoltDoc struct {
	Id          string `storm:"id"`
	ParitionKey string `storm:"index"`
	SortKey     string `storm:"index"`
	Value       map[string]interface{}
}

func (s *BoltDocService) Get(key sdk.Key, subKey *sdk.Key) (map[string]interface{}, error) {
	err := document.ValidateKeys(key, subKey)
	if err != nil {
		return nil, err
	}

	db, err := s.createdDb(key.Collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	doc := createDoc(key, subKey)

	err = db.One(idName, doc.Id, &doc)

	if err != nil {
		return nil, err
	}

	return doc.Value, nil
}

func (s *BoltDocService) Set(key sdk.Key, subKey *sdk.Key, value map[string]interface{}) error {
	err := document.ValidateKeys(key, subKey)
	if err != nil {
		return err
	}

	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	db, err := s.createdDb(key.Collection)
	if err != nil {
		return err
	}
	defer db.Close()

	doc := createDoc(key, subKey)
	doc.Value = value

	return db.Save(&doc)
}

func (s *BoltDocService) Delete(key sdk.Key, subKey *sdk.Key) error {
	err := document.ValidateKeys(key, subKey)
	if err != nil {
		return err
	}

	db, err := s.createdDb(key.Collection)
	if err != nil {
		return err
	}
	defer db.Close()

	doc := createDoc(key, subKey)

	err = db.DeleteStruct(&doc)

	// TODO: delete sub collection records

	return err
}

func (s *BoltDocService) Query(key sdk.Key, subcollection string, expressions []sdk.QueryExpression, limit int, pagingToken map[string]string) (*sdk.QueryResult, error) {
	err := document.ValidateCollection(key.Collection, subcollection)
	if err != nil {
		return nil, err
	}

	err = document.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	db, err := s.createdDb(key.Collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = document.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	// Build up chain of expression matchers
	matchers := []q.Matcher{}

	// Apply collection/sub-collection filters
	if key.Id != "" {
		matchers = append(matchers, q.Eq(partionKeyName, key.Id))
	}
	if subcollection != "" {
		matchers = append(matchers, q.Gte(sortKeyName, subcollection+"#"))
		matchers = append(matchers, q.Lt(sortKeyName, document.GetEndRangeValue(subcollection+"#")))
	} else {
		matchers = append(matchers, q.Eq(sortKeyName, key.Collection+"#"))
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
		if expStr.Len() > 0 {
			expStr.WriteString(" && ")
		}
		if exp.Operator == "startsWith" {
			expStr.WriteString(exp.Operand + " >= '" + exp.Value + "' && ")
			expStr.WriteString(exp.Operand + " < '" + document.GetEndRangeValue(exp.Value) + "'")

		} else {
			expStr.WriteString(exp.Operand + " " + exp.Operator + " '" + exp.Value + "'")
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
	results := make([]map[string]interface{}, 0)
	scanCount := 0
	for _, doc := range docs {

		if filterExp != nil {
			eval, err := filterExp.Evaluate(doc.Value)
			if err != nil {
				return nil, err
			}
			include := eval.(bool)
			if include {
				results = append(results, doc.Value)
			}

		} else {
			results = append(results, doc.Value)
		}

		scanCount += 1

		// Break if greater than fetch limit
		if limit > 0 && len(results) == limit {
			break
		}
	}

	// Provide paging token to skip previous reads
	var resultPagingToken map[string]string
	if limit > 0 && len(results) == limit {
		resultPagingToken = make(map[string]string)
		resultPagingToken[skipTokenName] = fmt.Sprintf("%v", pagingSkip+scanCount)
	}

	return &sdk.QueryResult{
		Data:        results,
		PagingToken: resultPagingToken,
	}, nil
}

// New - Create a new dev KV plugin
func New() (*BoltDocService, error) {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", DEFAULT_DIR)

	// Check whether file exists
	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		// Make diretory if not present
		err := os.MkdirAll(dbDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	return &BoltDocService{dbDir: dbDir}, nil
}

func (s *BoltDocService) createdDb(collection string) (*storm.DB, error) {
	err := document.ValidateCollection(collection, "")
	if err != nil {
		return nil, err
	}

	dbPath := s.dbDir + strings.ToLower(collection) + ".db"

	options := storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second})
	db, err := storm.Open(dbPath, options)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDoc(key sdk.Key, subKey *sdk.Key) BoltDoc {

	doc := BoltDoc{
		Id:          key.Id,
		ParitionKey: key.Id,
	}
	if subKey != nil {
		doc.Id += "_" + subKey.Id
	}

	// Top level collection item
	if subKey == nil {
		doc.SortKey = key.Collection + "#"

	} else {
		// Sub-collection item
		doc.SortKey = subKey.Collection + "#" + subKey.Id
	}

	return doc
}