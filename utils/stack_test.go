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

	"github.com/nitric-dev/membrane/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Function Test Cases
const YAML_VALID = "test/nitric-valid.yaml"

var _ = Describe("Utils", func() {
	defer GinkgoRecover()

	Context("NewStack", func() {
		When("Valid stack definition", func() {
			It("Should return new stack", func() {
				stack, err := utils.NewStack(YAML_VALID)
				Expect(stack).ToNot(BeNil())
				Expect(err).To(BeNil())
				Expect(stack.Collections).To(HaveLen(2))
			})
		})
	})
	Context("HasCollection", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("valid selection", func() {
			It("Should return true", func() {
				Expect(stack.HasCollection("users")).To(BeTrue())
			})
		})
		When("invalid valid selection", func() {
			It("Should return false", func() {
				Expect(stack.HasCollection("unknown")).To(BeFalse())
			})
		})
	})
	Context("HasSubCollection", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("invalid selection", func() {
			It("Should return false", func() {
				Expect(stack.HasSubCollection("unknown", "orders")).To(BeFalse())
			})
		})
		When("invalid selection", func() {
			It("Should return false", func() {
				Expect(stack.HasSubCollection("customers", "unknown")).To(BeFalse())
			})
		})
		When("valid selection", func() {
			It("Should return true", func() {
				Expect(stack.HasSubCollection("customers", "orders")).To(BeTrue())
			})
		})
	})
	Context("SubCollectionNames", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("valid selection", func() {
			It("Should return true", func() {
				names, err := stack.SubCollectionNames("customers")
				Expect(names).To(HaveLen(3))
				Expect(names[0]).To(BeEquivalentTo("addresses"))
				Expect(names[1]).To(BeEquivalentTo("orders"))
				Expect(names[2]).To(BeEquivalentTo("payments"))
				Expect(err).To(BeNil())
			})
		})
		When("invalid selection", func() {
			It("Should return error", func() {
				names, err := stack.SubCollectionNames("unknown")
				Expect(names).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Context("NewStackDefault", func() {
		When("default not available", func() {
			It("Should return error", func() {
				stack, err := utils.NewStackDefault()
				Expect(stack).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
		When("environment variable configured", func() {
			It("Should return new stack", func() {
				os.Setenv(utils.NITRIC_HOME, "test/")
				os.Setenv(utils.NITRIC_YAML, "nitric-valid.yaml")
				stack, err := utils.NewStackDefault()
				Expect(stack).ToNot(BeNil())
				Expect(err).To(BeNil())
				Expect(stack.Collections).To(HaveLen(2))
			})
		})
	})
})
