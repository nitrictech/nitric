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
				Expect(err.Error()).To(Equal("Secret name must not be blank"))
			})
		})

		When("given an invalid secret name", func() {
			err := ValidateSecretName("@test-secret")

			It("should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})

			It("error should indicate a secret name matching the valid pattern be provided", func() {
				Expect(err.Error()).To(Equal("Secret name must match pattern: ^\\w+(-\\w+)*$"))
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
