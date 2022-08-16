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

	boltdb_service "github.com/nitrictech/nitric/pkg/plugins/document/boltdb"
	"github.com/nitrictech/nitric/pkg/utils"

	. "github.com/onsi/ginkgo"

	test "github.com/nitrictech/nitric/tests/plugins/document"
)

var _ = Describe("Bolt", func() {
	docPlugin, err := boltdb_service.New()
	if err != nil {
		panic(err)
	}

	BeforeSuite(func() {
		test.LoadItemsData(docPlugin)
	})

	AfterSuite(func() {
		local_boltdb_path := utils.GetRelativeDevPath(boltdb_service.DEV_SUB_DIRECTORY)
		err = os.RemoveAll(local_boltdb_path)
		if err == nil {
			os.Remove(local_boltdb_path)
			os.Remove(utils.GetDevVolumePath())
		}
	})

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
	test.QueryStreamTests(docPlugin)
})
