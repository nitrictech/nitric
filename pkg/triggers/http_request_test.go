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

package triggers

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Http Request", func() {
	Context("parsePathParams", func() {
		When("Parsing single path parameters with a match", func() {
			params, err := parsePathParams("/user/:userId", "/user/123")

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should contain a userId path parameter", func() {
				Expect(params["userId"]).To(Equal("123"))
			})
		})

		When("Parsing deep path parameters with a match", func() {
			params, err := parsePathParams("/user/:userId/order/:orderId", "/user/123/order/12345")

			It("should not return an error", func() {
				Expect(err).ShouldNot(HaveOccurred())
			})

			It("should contain a userId path parameter", func() {
				Expect(params["userId"]).To(Equal("123"))
			})

			It("should contain an orderId path parameter", func() {
				Expect(params["orderId"]).To(Equal("12345"))
			})
		})

		When("Parsing deep path parameters without a match", func() {
			params, err := parsePathParams("/user/:userId/items/:itemId", "/user/123/order/12345")

			It("should not return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("should return nil parameters", func() {
				Expect(params).To(BeNil())
			})
		})

		When("Parsing parameters with a segment mismatch", func() {
			params, err := parsePathParams("/user/:userId/items/:itemId", "/user/123")

			It("should not return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("should return nil parameters", func() {
				Expect(params).To(BeNil())
			})
		})

		When("Parsing single path parameters without a match", func() {
			params, err := parsePathParams("/user/:userId", "/testing/123")

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("should return nil params", func() {
				Expect(params).To(BeNil())
			})
		})
	})
})
