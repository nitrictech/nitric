package utils_test

import (
	"os"

	"github.com/nitric-dev/membrane/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Function Test Cases
const YAML_INVALID_1 = "test/nitric-invalid-1.yaml"
const YAML_INVALID_2 = "test/nitric-invalid-2.yaml"
const YAML_INVALID_3 = "test/nitric-invalid-3.yaml"
const YAML_VALID = "test/nitric-valid.yaml"

var _ = Describe("Utils", func() {
	defer GinkgoRecover()

	Context("NewStack", func() {
		When("Valid stack definition", func() {
			It("Should return new stack", func() {
				stack, err := utils.NewStack(YAML_VALID)
				Expect(stack).ToNot(BeNil())
				Expect(err).To(BeNil())
				Expect(stack.Collections).To(HaveLen(4))
			})
		})
		When("Configured default 'orders' collection", func() {
			It("Should have default attributes and indexes", func() {
				stack, _ := utils.NewStack(YAML_VALID)

				indexes, err := stack.CollectionIndexes("orders")
				Expect(indexes).To(BeEquivalentTo([]string{"key"}))
				Expect(err).To(BeNil())

				indexes, err = stack.CollectionIndexesComposite("orders")
				Expect(indexes).To(BeNil())
				Expect(err).ToNot(BeNil())

				key, err := stack.CollectionIndexesUnique("orders")
				Expect(key).To(BeEquivalentTo("key"))
				Expect(err).To(BeNil())

				attributes, err := stack.CollectionAttributes("orders")
				Expect(attributes).To(BeEquivalentTo([]string{"key", "value"}))
				Expect(err).To(BeNil())
			})
		})
		When("Invalid stack definition", func() {
			It("Should return error", func() {
				stack, err := utils.NewStack(YAML_INVALID_1)
				Expect(stack).To(BeNil())
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(BeEquivalentTo("nitric-website collections: application: indexes: composite: requires 2 values [pk]"))
			})
		})
		When("Invalid stack definition", func() {
			It("Should return error", func() {
				stack, err := utils.NewStack(YAML_INVALID_2)
				Expect(stack).To(BeNil())
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(BeEquivalentTo("nitric-website collections: application: indexes: composite: not defined"))
			})
		})
		When("Invalid stack definition", func() {
			It("Should return error", func() {
				stack, err := utils.NewStack(YAML_INVALID_3)
				Expect(stack).To(BeNil())
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(BeEquivalentTo("nitric-website collections: users: indexes: unique: unknown has no matching collection attribute"))
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
	Context("CollectionAttributes", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("valid selection", func() {
			It("Should return true", func() {
				names, err := stack.CollectionAttributes("users")
				Expect(names).To(HaveLen(2))
				Expect(err).To(BeNil())
			})
		})
		When("invalid selection", func() {
			It("Should return error", func() {
				names, err := stack.CollectionAttributes("unknown")
				Expect(names).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Context("CollectionFilterAttributes", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("valid selection", func() {
			It("Should return true", func() {
				names, err := stack.CollectionFilterAttributes("users")
				Expect(names).To(HaveLen(0))
				Expect(err).To(BeNil())
			})
		})
		When("valid selection", func() {
			It("Should return true", func() {
				names, err := stack.CollectionFilterAttributes("application")
				Expect(err).To(BeNil())
				Expect(names).To(BeEquivalentTo([]string{"created"}))
			})
		})
		When("invalid selection", func() {
			It("Should return error", func() {
				names, err := stack.CollectionFilterAttributes("unknown")
				Expect(names).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Context("CollectionIndexes", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("valid unique selection", func() {
			It("Should return key", func() {
				names, err := stack.CollectionIndexes("users")
				Expect(names).To(HaveLen(1))
				Expect(err).To(BeNil())
			})
		})
		When("valid composite selection", func() {
			It("Should return key", func() {
				names, err := stack.CollectionIndexes("application")
				Expect(names).To(HaveLen(2))
				Expect(err).To(BeNil())
			})
		})
		When("invalid valid selection", func() {
			It("Should return key", func() {
				names, err := stack.CollectionIndexes("unknown")
				Expect(names).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Context("CollectionIndexesUnique", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("valid selection", func() {
			It("Should return key", func() {
				name, err := stack.CollectionIndexesUnique("users")
				Expect(name).To(BeEquivalentTo("key"))
				Expect(err).To(BeNil())
			})
		})
		When("invalid valid selection", func() {
			It("Should return key", func() {
				name, err := stack.CollectionIndexesUnique("unknown")
				Expect(name).To(BeEquivalentTo(""))
				Expect(err).ToNot(BeNil())
			})
		})
	})
	Context("CollectionIndexesComposite", func() {
		stack, _ := utils.NewStack(YAML_VALID)
		When("valid selection", func() {
			It("Should return key", func() {
				names, err := stack.CollectionIndexesComposite("application")
				Expect(names).To(BeEquivalentTo([]string{"pk", "sk"}))
				Expect(err).To(BeNil())
			})
		})
		When("valid selection", func() {
			It("Should return key", func() {
				names, err := stack.CollectionIndexesComposite("users")
				Expect(names).To(BeNil())
				Expect(err).ToNot(BeNil())
			})
		})
		When("invalid valid selection", func() {
			It("Should return key", func() {
				names, err := stack.CollectionIndexesComposite("unknown")
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
				Expect(stack.Collections).To(HaveLen(4))
			})
		})
	})
})
