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

package kv_service

import (
	"fmt"

	scribble "github.com/nanobox-io/golang-scribble"
	"github.com/nitric-dev/membrane/plugins/kv"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
)

type DevKVService struct {
	sdk.UnimplementedKeyValuePlugin
	db ScribbleIface
}

func (s *DevKVService) Get(collection string, key map[string]interface{}) (map[string]interface{}, error) {
	dbKey := fmt.Sprintf("%v", key["key"])
	value := make(map[string]interface{})
	err := s.db.Read(collection, dbKey, &value)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *DevKVService) Put(collection string, key map[string]interface{}, value map[string]interface{}) error {
	dbKey := fmt.Sprintf("%v", key["key"])
	return s.db.Write(collection, dbKey, value)
}

func (s *DevKVService) Delete(collection string, key map[string]interface{}) error {
	dbKey := fmt.Sprintf("%v", key["key"])
	error := s.db.Delete(collection, dbKey)

	if error != nil {
		return error
	}

	return nil
}

func (p *DevKVService) Query(collection string, expressions []sdk.QueryExpression, limit int) ([]map[string]interface{}, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
	collection, error := kv.GetCollection(collection)
	if error != nil {
		return nil, error
	}
	error = kv.ValidateExpressions(expressions)
	if error != nil {
		return nil, error
	}

	s.db.
}

// New - Create a new dev KV plugin
func New() (sdk.KeyValueService, error) {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", "/nitric/")
	db, err := scribble.New(dbDir, nil)

	if err != nil {
		return nil, err
	}

	return &DevKVService{
		db: db,
	}, nil
}

func NewWithDB(db ScribbleIface) (sdk.KeyValueService, error) {
	return &DevKVService{
		db: db,
	}, nil
}
