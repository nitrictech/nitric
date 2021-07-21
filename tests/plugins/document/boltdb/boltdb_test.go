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

package boltdb_service_test

import (
	"os"
	"strings"
	"time"

	"github.com/asdine/storm"
	boltdb_service "github.com/nitric-dev/membrane/pkg/plugins/document/boltdb"
	"github.com/nitric-dev/membrane/pkg/utils"
	"go.etcd.io/bbolt"

	test "github.com/nitric-dev/membrane/tests/plugins/document"
	. "github.com/onsi/ginkgo"
)

func createDb(colName string) (*storm.DB, error) {
	dbDir := utils.GetEnv("LOCAL_DB_DIR", boltdb_service.DEFAULT_DIR)

	// Check whether file exists
	_, err := os.Stat(dbDir)
	if os.IsNotExist(err) {
		// Make directory if not present
		err := os.MkdirAll(dbDir, 0777)
		if err != nil {
			panic(err)
		}
	}

	dbPath := dbDir + strings.ToLower(colName) + ".db"

	options := storm.BoltOptions(0600, &bbolt.Options{Timeout: 1 * time.Second})

	return storm.Open(dbPath, options)
}

func deleteCollection(collection string) {
	db, err := createDb(collection)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var docs []boltdb_service.BoltDoc
	db.All(&docs)

	for _, childDoc := range docs {
		err = db.DeleteStruct(&childDoc)
		if err != nil {
			panic(err)
		}
	}
}

var _ = Describe("Bolt", func() {

	docPlugin, err := boltdb_service.New()
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		test.LoadItemsData(docPlugin)
	})

	AfterEach(func() {
		deleteCollection("customers")
		deleteCollection("users")
		deleteCollection("items")
		deleteCollection("parentItems")
	})

	AfterSuite(func() {
		err = os.RemoveAll(boltdb_service.DEFAULT_DIR)
		if err == nil {
			os.Remove(boltdb_service.DEFAULT_DIR)
			os.Remove("nitric/")
		}
	})

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
})
