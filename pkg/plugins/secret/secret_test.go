package secret

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Unimplemented Secret Plugin Tests", func() {
	uisp := &UnimplementedSecretPlugin{}

	Context("Put", func() {
		When("Calling Put on UnimplementedSecretPlugin", func() {
			_, err := uisp.Put(nil, nil)

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})

	Context("Access", func() {
		When("Calling Access on UnimplementedSecretPlugin", func() {
			_, err := uisp.Access(nil)

			It("should return an unimplemented error", func() {
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("UNIMPLEMENTED"))
			})
		})
	})
})
