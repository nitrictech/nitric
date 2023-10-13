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

package utils_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/core/pkg/utils"
)

var _ = Describe("Env", func() {
	Context("GetEnv", func() {
		When("Calling GetEnv for a non-set env variable", func() {
			env := utils.GetEnv("MY_FAKE_ENV", "my-fallback")

			It("should return the fallback value", func() {
				Expect(env).To(Equal("my-fallback"))
			})
		})

		When("Calling GetEnv for a set env variable", func() {
			os.Setenv("MY_TEST_ENV", "testing")

			env := utils.GetEnv("MY_TEST_ENV", "my-fallback")

			It("should return the set value", func() {
				Expect(env).To(Equal("testing"))
			})
		})
	})
})
