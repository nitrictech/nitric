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
	"os"
	"strings"
	"time"

	"github.com/asdine/storm"
	"go.etcd.io/bbolt"

	"github.com/nitrictech/nitric/pkg/plugins/errors"
	"github.com/nitrictech/nitric/pkg/plugins/errors/codes"
	"github.com/nitrictech/nitric/pkg/plugins/storage"
	"github.com/nitrictech/nitric/pkg/utils"
)

const DEV_SUB_DIRECTORY = "./buckets/"

type BoltStorageService struct {
	storage.UnimplementedStoragePlugin
	dbDir string
}

type Object struct {
	Key  string `storm:"id"`
	Data []byte
}

// Write - will create a new item or overwrite an existing item in storage
func (s *BoltStorageService) Write(bucket string, key string, object []byte) error {
	newErr := errors.ErrorsWithScope(
		"BoltStorageService.Write",
		map[string]interface{}{
			"bucket":     bucket,
			"key":        key,
			"object.len": len(object),
		},
	)

	if object == nil {
		return newErr(
			codes.InvalidArgument,
			"provide non-nil object",
			nil,
		)
	}
	if len(object) == 0 {
		return newErr(
			codes.InvalidArgument,
			"provide non-empty object",
			nil,
		)
	}

	db, err := s.createDb(bucket)
	if err != nil {
		return newErr(
			codes.FailedPrecondition,
			"createDb error",
			err,
		)
	}
	defer db.Close()

	obj := Object{
		Key:  key,
		Data: object,
	}

	err = db.Save(&obj)
	if err != nil {
		return newErr(
			codes.Internal,
			"error storing object",
			err,
		)
	}

	return nil
}

// Read - reads an item from Storage
func (s *BoltStorageService) Read(bucket string, key string) ([]byte, error) {
	newErr := errors.ErrorsWithScope(
		"BoltStorageService.Read",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	db, err := s.createDb(bucket)
	if err != nil {
		return nil, newErr(
			codes.FailedPrecondition,
			"createDb error",
			err,
		)
	}
	defer db.Close()

	var obj = Object{}
	err = db.One("Key", key, &obj)
	if err != nil {
		// TODO: not found OK
		return nil, newErr(
			codes.Internal,
			"failed to retrieve key",
			err,
		)
	}

	return obj.Data, nil
}

// Delete - deletes an item from Storage
func (s *BoltStorageService) Delete(bucket string, key string) error {
	newErr := errors.ErrorsWithScope(
		"BoltStorageService.Delete",
		map[string]interface{}{
			"bucket": bucket,
			"key":    key,
		},
	)

	db, err := s.createDb(bucket)
	if err != nil {
		return newErr(
			codes.FailedPrecondition,
			"createDb error",
			err,
		)
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
func New() (storage.StorageService, error) {
	dbDir := utils.GetEnv("LOCAL_BLOB_DIR", utils.GetRelativeDevPath(DEV_SUB_DIRECTORY))

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
