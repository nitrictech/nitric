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

package secret

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Secret Plugin Utils", func() {
	Context("ValidateSecretName", func() {
		When("given a blank secret name", func() {
			err := ValidateSecretName("")

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("should indicate a non-blank secret name be provided", func() {
				Expect(err.Error()).To(Equal("provide non-blank secret name"))
			})
		})

		When("given an invalid secret name", func() {
			err := ValidateSecretName("@test-secret")

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("error should indicate a secret name matching the valid pattern be provided", func() {
				Expect(err.Error()).To(Equal("provide secret name matching pattern: ^\\w+(-\\w+)*$"))
			})
		})

		Context("Valid Secret Names", func() {
			When("with mixed case", func() {
				err := ValidateSecretName("TeStSeCrEt")

				It("should not return an error", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("with underscores", func() {
				err := ValidateSecretName("_TeSt_SeCrEt")

				It("should not return an error", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})
			})

			When("with hyphens", func() {
				err := ValidateSecretName("TeSt-SeCrEt")

				It("should not return an error", func() {
					Expect(err).ShouldNot(HaveOccurred())
				})
			})
		})
	})
})
