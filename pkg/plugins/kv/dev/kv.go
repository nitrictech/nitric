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
	scribble "github.com/nanobox-io/golang-scribble"
	utils2 "github.com/nitric-dev/membrane/pkg/utils"
	"github.com/nitric-dev/membrane/pkg/sdk"
)

type DevKVService struct {
	sdk.UnimplementedKeyValuePlugin
	db ScribbleIface
}

func (s *DevKVService) Get(collection string, key string) (map[string]interface{}, error) {
	value := make(map[string]interface{})
	err := s.db.Read(collection, key, &value)

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (s *DevKVService) Put(collection string, key string, value map[string]interface{}) error {
	return s.db.Write(collection, key, value)
}

func (s *DevKVService) Delete(collection string, key string) error {
	error := s.db.Delete(collection, key)

	if error != nil {
		return error
	}

	return nil
}

// New - Create a new dev KV plugin
func New() (sdk.KeyValueService, error) {
	dbDir := utils2.GetEnv("LOCAL_DB_DIR", "/nitric/")
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
