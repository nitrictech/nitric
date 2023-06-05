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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/core/pkg/utils"
)

var _ = Describe("Parse", func() {
	Context("PercentFromIntString", func() {
		When("Calling SplitPath with a path", func() {
			strs := utils.SplitPath("test/testOne")

			It("should return strings split by path", func() {
				Expect(strs).To(HaveLen(2))
				Expect(strs[0]).To(Equal("test"))
				Expect(strs[1]).To(Equal("testOne"))
			})
		})

		When("Calling SplitPath with non-path", func() {
			strs := utils.SplitPath("test")

			It("should return strings split by path", func() {
				Expect(strs).To(HaveLen(1))
				Expect(strs[0]).To(Equal("test"))
			})
		})
	})
})
