package worker

import (
	"github.com/golang/mock/gomock"
	mock "github.com/nitrictech/nitric/mocks/worker"
	"github.com/nitrictech/nitric/pkg/triggers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("SubscriptionWorker", func() {

	Context("Http", func() {
		subWrkr := &SubscriptionWorker{}

		When("calling HandlesHttpRequest", func() {
			It("should return false", func() {
				Expect(subWrkr.HandlesHttpRequest(&triggers.HttpRequest{})).To(BeFalse())
			})
		})

		When("calling HandleHttpRequest", func() {
			It("should return an error", func() {
				_, err := subWrkr.HandleHttpRequest(&triggers.HttpRequest{})
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Event", func() {
		When("calling HandlesEvent with the wrong topic", func() {
			subWrkr := &SubscriptionWorker{
				topic: "bad",
			}

			It("should return true", func() {
				Expect(subWrkr.HandlesEvent(&triggers.Event{
					Topic: "test",
				})).To(BeFalse())
			})
		})

		When("calling HandlesEvent with the correct topic", func() {
			subWrkr := &SubscriptionWorker{
				topic: "test",
			}

			It("should return true", func() {
				Expect(subWrkr.HandlesEvent(&triggers.Event{
					Topic: "test",
				})).To(BeTrue())
			})
		})

		When("calling HandleEvent", func() {
			It("should call the base grpc workers HandleEvent", func() {
				ctrl := gomock.NewController(GinkgoT())
				mGrpc := mock.NewMockGrpcWorker(ctrl)

				By("calling the base grpc handler HandleEvent method")
				mGrpc.EXPECT().HandleEvent(gomock.Any()).Times(1)

				subWrkr := &SubscriptionWorker{
					topic:      "test",
					GrpcWorker: mGrpc,
				}

				err := subWrkr.HandleEvent(&triggers.Event{})

				Expect(err).ShouldNot(HaveOccurred())
				ctrl.Finish()
			})
		})
	})
})
