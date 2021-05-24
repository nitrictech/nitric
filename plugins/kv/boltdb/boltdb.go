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
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/asdine/storm/q"
	"github.com/nitric-dev/membrane/plugins/kv"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
	"go.etcd.io/bbolt"
)

const DEFAULT_DIR = "nitric/collections/"

type DevKVService struct {
	sdk.UnimplementedKeyValuePlugin
	dbDir string
}

type Document struct {
	Key        string `storm:"id"`
	PartionKey string `storm:"index"`
	SortKey    string `storm:"index"`
	Attribute1 string `storm:"index"`
	Attribute2 string `storm:"index"`
	Attribute3 string `storm:"index"`
	Value      map[string]interface{}
}

func (d Document) String() string {
	return fmt.Sprintf("{Key:'%v', ParitionKey:'%v', SortKey:'%v', Value:%v}", d.Key, d.PartionKey, d.SortKey, d.Value)
}

func (s *DevKVService) Get(collection string, key map[string]interface{}) (map[string]interface{}, error) {
	db, err := s.createDb(collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = kv.ValidateKeyMap(key)
	if err != nil {
		return nil, err
	}

	keyValue := kv.GetKeyValue(key)

	var doc = Document{}
	err = db.One("Key", keyValue, &doc)
	if err != nil {
		return nil, err
	}

	return doc.Value, nil
}

func (s *DevKVService) Put(collection string, key map[string]interface{}, value map[string]interface{}) error {
	db, err := s.createDb(collection)
	if err != nil {
		return err
	}
	defer db.Close()

	err = kv.ValidateKeyMap(key)
	if err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	keyValue := kv.GetKeyValue(key)
	doc := Document{
		Key:   keyValue,
		Value: value,
	}

	keyValues := kv.GetKeyValues(key)
	if len(keyValues) > 1 {
		doc.PartionKey = keyValues[0]
		doc.SortKey = keyValues[1]
	}

	err = db.Save(&doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *DevKVService) Delete(collection string, key map[string]interface{}) error {
	db, err := s.createDb(collection)
	if err != nil {
		return err
	}
	defer db.Close()

	err = kv.ValidateKeyMap(key)
	if err != nil {
		return err
	}

	keyValue := kv.GetKeyValue(key)
	doc := Document{
		Key: keyValue,
	}

	return db.DeleteStruct(&doc)

	// TODO: discuss delete behaviour for Get and Delete
	// err = db.DeleteStruct(&doc)
	// if err != nil && err.Error() != "not found" {
	// 	return err
	// } else {
	// 	return nil
	// }
}

func (s *DevKVService) Query(collection string, expressions []sdk.QueryExpression, limit int) ([]map[string]interface{}, error) {
	db, err := s.createDb(collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = kv.ValidateExpressions(expressions)
	if err != nil {
		return nil, err
	}

	// Build up chain of expression matchers
	matchers := []q.Matcher{}

	for i, exp := range expressions {
		// Operand is determined by expressions index
		// TODO: replace attribute lookup with nitric stack 'collections' mapping
		operand := "PartionKey"
		if i == 1 {
			operand = "SortKey"
		} else if i == 2 {
			operand = "Attribute1"
		} else if i == 3 {
			operand = "Attribute2"
		} else if i == 4 {
			operand = "Attribute3"
		}

		if exp.Operator == "==" {
			matchers = append(matchers, q.Eq(operand, exp.Value))
		} else if exp.Operator == "<" {
			matchers = append(matchers, q.Lt(operand, exp.Value))
		} else if exp.Operator == ">" {
			matchers = append(matchers, q.Gt(operand, exp.Value))
		} else if exp.Operator == "<=" {
			matchers = append(matchers, q.Lte(operand, exp.Value))
		} else if exp.Operator == ">=" {
			matchers = append(matchers, q.Gte(operand, exp.Value))
		} else if exp.Operator == "startsWith" {
			matchers = append(matchers, q.Gte(operand, exp.Value))
			matchers = append(matchers, q.Lt(operand, kv.GetEndRangeValue(exp.Value)))
		}
	}

	// Create query object
	matcher := q.And(matchers[:]...)
	query := db.Select(matcher)

	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute query
	var docs []Document
	query.Find(&docs)

	results := make([]map[string]interface{}, 0)
	for _, doc := range docs {
		results = append(results, doc.Value)
	}

	return results, nil
}

// New - Create a new dev KV plugin
func New() (sdk.KeyValueService, error) {
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

	return &DevKVService{dbDir: dbDir}, nil
}

func (s *DevKVService) createDb(collection string) (*storm.DB, error) {
	err := kv.ValidateCollection(collection)
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
