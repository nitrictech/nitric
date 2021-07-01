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

	ds_plugin "github.com/nitric-dev/membrane/plugins/document/boltdb"
	test "github.com/nitric-dev/membrane/plugins/document/test"
	"github.com/nitric-dev/membrane/utils"
	"github.com/onsi/ginkgo"
)

var _ = ginkgo.Describe("Bolt", func() {

	os.Setenv(utils.NITRIC_HOME, "../test/")
	os.Setenv(utils.NITRIC_YAML, "nitric.yaml")

	docPlugin, err := ds_plugin.New()
	if err != nil {
		panic(err)
	}

	ginkgo.BeforeSuite(func() {
		test.LoadItemsData(docPlugin)
	})

	ginkgo.AfterSuite(func() {
		err = os.RemoveAll(ds_plugin.DEFAULT_DIR)
		if err == nil {
			os.Remove(ds_plugin.DEFAULT_DIR)
			os.Remove("nitric/")
		}
	})

	test.GetTests(docPlugin)
	test.SetTests(docPlugin)
	test.DeleteTests(docPlugin)
	test.QueryTests(docPlugin)
})
