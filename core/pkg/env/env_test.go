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

package env

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Env", func() {
	Context("GetEnv", func() {
		When("TEST_ENV is not set", func() {
			env := GetEnv("TEST_ENV", "my-fallback")

			It("should return the fallback string value", func() {
				Expect(env.String()).To(Equal("my-fallback"))
			})
		})

		When("TEST_ENV is a string value", func() {
			os.Setenv("TEST_ENV", "testing")

			env := GetEnv("TEST_ENV", "my-fallback")

			It("should return the set value", func() {
				Expect(env.String()).To(Equal("testing"))
			})

			It("should not return bool values", func() {
				_, err := env.Bool()

				Expect(err).To(HaveOccurred())
			})

			It("should not return int values", func() {
				_, err := env.Int()

				Expect(err).To(HaveOccurred())
			})
		})

		When("TEST_ENV is an integer string", func() {
			os.Setenv("TEST_ENV", "123")

			env := GetEnv("MY_TEST_ENV", "")

			It("should return int values", func() {
				val, err := env.Int()

				By("correctly parsing the int value")
				Expect(val).To(Equal(123))

				By("not returning an error")
				Expect(err).To(BeNil())
			})

			It("should not return bool values", func() {
				_, err := env.Bool()

				Expect(err).To(HaveOccurred())
			})
		})

		When("TEST_ENV is a boolean string", func() {
			os.Setenv("TEST_ENV", "True")

			env := GetEnv("MY_TEST_ENV", "")

			It("should return boolean values", func() {
				val, err := env.Bool()

				By("correctly parsing the bool value")
				Expect(val).To(Equal(true))

				By("not returning an error")
				Expect(err).To(BeNil())
			})

			It("should not return int values", func() {
				_, err := env.Int()

				Expect(err).To(HaveOccurred())
			})
		})
	})
})
