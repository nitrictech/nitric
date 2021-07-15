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

package boltdb_storage_service

import (
	"fmt"
	utils2 "github.com/nitric-dev/membrane/pkg/utils"
	"os"
	"strings"
	"time"

	"github.com/asdine/storm"
	"github.com/nitric-dev/membrane/pkg/sdk"
	"go.etcd.io/bbolt"
)

const DEFAULT_DIR = "nitric/buckets/"

type BoltStorageService struct {
	sdk.UnimplementedStoragePlugin
	dbDir string
}

type Object struct {
	Key  string `storm:"id"`
	Data []byte
}

// Write - will create a new item or overwrite an existing item in storage
func (s *BoltStorageService) Write(bucket string, key string, object []byte) error {
	if bucket == "" {
		return fmt.Errorf("provide non-blank bucket")
	}
	if key == "" {
		return fmt.Errorf("provide non-blank key")
	}
	if object == nil {
		return fmt.Errorf("provide non-nil object")
	}
	if len(object) == 0 {
		return fmt.Errorf("provide non-empty object")
	}

	db, err := s.createDb(bucket)
	if err != nil {
		return err
	}
	defer db.Close()

	obj := Object{
		Key:  key,
		Data: object,
	}

	err = db.Save(&obj)
	if err != nil {
		return fmt.Errorf("Error storing %s : %v", key, err)
	}

	return nil
}

// Read - reads an item from Storage
func (s *BoltStorageService) Read(bucket string, key string) ([]byte, error) {
	if bucket == "" {
		return nil, fmt.Errorf("provide non-blank bucket")
	}
	if key == "" {
		return nil, fmt.Errorf("provide non-blank key")
	}

	db, err := s.createDb(bucket)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var obj = Object{}
	err = db.One("Key", key, &obj)
	if err != nil {
		// TODO: not found OK
		return nil, err
	}

	return obj.Data, nil
}

// Delete - deletes an item from Storage
func (s *BoltStorageService) Delete(bucket string, key string) error {
	if bucket == "" {
		return fmt.Errorf("provide non-blank bucket")
	}
	if key == "" {
		return fmt.Errorf("provide non-blank key")
	}

	db, err := s.createDb(bucket)
	if err != nil {
		return err
	}
	defer db.Close()

	doc := Object{
		Key: key,
	}
	return db.DeleteStruct(&doc)

	// TODO: discuss delete behaviour for Read and Delete
	// err = db.Delete(collection, keyValue)
	// if err != nil && err.Error() != "not found" {
	// 	return err
	// }
}

// New - Create a new BoltDB Storage plugin
func New() (sdk.StorageService, error) {
	dbDir := utils2.GetEnv("LOCAL_BLOB_DIR", DEFAULT_DIR)

	// Check whether file exists
	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		// Make diretory if not present
		err := os.MkdirAll(dbDir, 0777)
		if err != nil {
			return nil, err
		}
	}

	return &BoltStorageService{dbDir: dbDir}, nil
}

func (s *BoltStorageService) createDb(bucket string) (*storm.DB, error) {
	dbPath := s.dbDir + strings.ToLower(bucket) + ".db"

	options := storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second})
	db, err := storm.Open(dbPath, options)
	if err != nil {
		return nil, err
	}

	return db, nil
}
