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

type BoltKVService struct {
	sdk.UnimplementedKeyValuePlugin
	dbDir string
}

type BoltDoc struct {
	Key        string `storm:"id"`
	PartionKey string `storm:"index"`
	SortKey    string `storm:"index"`
	Attribute1 string `storm:"index"`
	Attribute2 string `storm:"index"`
	Attribute3 string `storm:"index"`
	Attribute4 string `storm:"index"`
	Attribute5 string `storm:"index"`
	Value      map[string]interface{}
}

// Implement sdk.KeyValueService interface

func (s *BoltKVService) Get(collection string, key map[string]interface{}) (map[string]interface{}, error) {
	db, err := s.createdDb(collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return nil, err
	}

	keyValue := kv.GetKeyValue(key)

	var doc = BoltDoc{}

	err = db.One("Key", keyValue, &doc)

	if err != nil {
		return nil, err
	}

	return doc.Value, nil
}

func (s *BoltKVService) Put(collection string, key map[string]interface{}, value map[string]interface{}) error {
	db, err := s.createdDb(collection)
	if err != nil {
		return err
	}
	defer db.Close()

	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return err
	}
	if value == nil {
		return fmt.Errorf("provide non-nil value")
	}

	doc := BoltDoc{
		Key:   kv.GetKeyValue(key),
		Value: value,
	}

	// Specify any composite indexes for filtering
	compositeKeys, err := kv.Stack.CollectionIndexesComposite(collection)
	if compositeKeys != nil && err == nil {
		doc.PartionKey = fmt.Sprintf("%v", key[compositeKeys[0]])
		doc.SortKey = fmt.Sprintf("%v", key[compositeKeys[1]])
	}

	// Project any collection filter attributes into Doc filter attributes
	filterAttributes, err := kv.Stack.CollectionFilterAttributes(collection)
	if filterAttributes != nil && err == nil {
		for i, name := range filterAttributes {
			valueStr := fmt.Sprintf("%v", value[name])
			switch i {
			case 0:
				doc.Attribute1 = valueStr
			case 1:
				doc.Attribute2 = valueStr
			case 2:
				doc.Attribute3 = valueStr
			case 3:
				doc.Attribute4 = valueStr
			case 4:
				doc.Attribute5 = valueStr
			}
		}
	}

	return db.Save(&doc)
}

func (s *BoltKVService) Delete(collection string, key map[string]interface{}) error {
	db, err := s.createdDb(collection)
	if err != nil {
		return err
	}
	defer db.Close()

	err = kv.ValidateKeyMap(collection, key)
	if err != nil {
		return err
	}

	keyValue := kv.GetKeyValue(key)
	doc := BoltDoc{
		Key: keyValue,
	}

	err = db.DeleteStruct(&doc)

	return err
}

func (s *BoltKVService) Query(collection string, expressions []sdk.QueryExpression, limit int) ([]map[string]interface{}, error) {
	db, err := s.createdDb(collection)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	err = kv.ValidateExpressions(collection, expressions)
	if err != nil {
		return nil, err
	}

	// Build up chain of expression matchers
	matchers := []q.Matcher{}

	uniqueKeyName, _ := kv.Stack.CollectionIndexesUnique(collection)
	compositeKeyNames, _ := kv.Stack.CollectionIndexesComposite(collection)
	filterNames, _ := kv.Stack.CollectionFilterAttributes(collection)

	for _, exp := range expressions {
		operand := ""

		if exp.Operand == uniqueKeyName {
			operand = "Key"
		}
		if index := utils.IndexOf(compositeKeyNames, exp.Operand); index != -1 {
			switch index {
			case 0:
				operand = "PartionKey"
			case 1:
				operand = "SortKey"
			}
		}
		if index := utils.IndexOf(filterNames, exp.Operand); index != -1 {
			switch index {
			case 0:
				operand = "Attribute1"
			case 1:
				operand = "Attribute2"
			case 2:
				operand = "Attribute3"
			case 3:
				operand = "Attribute4"
			case 4:
				operand = "Attribute5"
			}
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
	var docs []BoltDoc

	query.Find(&docs)

	results := make([]map[string]interface{}, 0)
	for _, doc := range docs {
		results = append(results, doc.Value)
	}

	return results, nil
}

// New - Create a new dev KV plugin
func New() (*BoltKVService, error) {
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

	return &BoltKVService{dbDir: dbDir}, nil
}

func (s *BoltKVService) createdDb(collection string) (*storm.DB, error) {
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
