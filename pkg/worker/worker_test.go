package worker

import (
	"github.com/nitrictech/nitric/pkg/triggers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Worker", func() {

	Context("UnimplementedWorker", func() {
		uiWrkr := &UnimplementedWorker{}

		When("calling HandlesEvent", func() {
			It("should return false", func() {
				Expect(uiWrkr.HandlesEvent(&triggers.Event{})).To(BeFalse())
			})
		})

		When("calling HandleEvent", func() {
			It("should return an error", func() {
				err := uiWrkr.HandleEvent(&triggers.Event{})
				Expect(err).Should(HaveOccurred())
			})
		})

		When("calling HandlesHttpRequest", func() {
			It("should return false", func() {
				Expect(uiWrkr.HandlesHttpRequest(&triggers.HttpRequest{})).To(BeFalse())
			})
		})

		When("calling HandleHttpRequest", func() {
			It("should return an error", func() {
				_, err := uiWrkr.HandleHttpRequest(&triggers.HttpRequest{})
				Expect(err).Should(HaveOccurred())
			})
		})

	})

})
