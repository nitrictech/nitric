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

	boltdb_service "github.com/nitric-dev/membrane/pkg/plugins/document/boltdb"
	"github.com/nitric-dev/membrane/pkg/sdk"

	test "github.com/nitric-dev/membrane/tests/plugins/document"
	. "github.com/onsi/ginkgo"
)

func deleteCollection(docPlugin sdk.DocumentService, collection string) {
	// TODO: delete items when query doc.key support available
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
		deleteCollection(docPlugin, "customers")
		deleteCollection(docPlugin, "users")
		deleteCollection(docPlugin, "items")
		deleteCollection(docPlugin, "parentItems")
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
