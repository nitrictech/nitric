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

var _ = Describe("Paths", func() {
	Context("PercentFromIntString", func() {
		When("Calling PercentFromIntString more a non-number", func() {
			_, err := utils.PercentFromIntString("testing")

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("Calling PercentFromIntString with number gte 100", func() {
			num, _ := utils.PercentFromIntString("100")

			It("should return 1", func() {
				Expect(num).To(Equal(float64(1)))
			})
		})

		When("Calling PercentFromIntString with number lte 0", func() {
			num, _ := utils.PercentFromIntString("0")

			It("should return 1", func() {
				Expect(num).To(Equal(float64(0)))
			})
		})

		When("Calling a number gt 0 and lt 100", func() {
			num, _ := utils.PercentFromIntString("97")

			It("should return 0.97", func() {
				Expect(num).To(Equal(0.97))
			})
		})
	})
})
