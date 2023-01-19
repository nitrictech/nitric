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

package worker

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock "github.com/nitrictech/nitric/core/mocks/adapter"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

var _ = Describe("SubscriptionWorker", func() {
	Context("Http", func() {
		httpTrigger := &v1.TriggerRequest{
			Context: &v1.TriggerRequest_Http{
				Http: &v1.HttpTriggerContext{},
			},
		}
		subWrkr := &SubscriptionWorker{}

		When("calling HandlesHttpRequest", func() {
			It("should return false", func() {
				Expect(subWrkr.HandlesTrigger(httpTrigger)).To(BeFalse())
			})
		})

		When("calling HandleHttpRequest", func() {
			It("should return an error", func() {
				_, err := subWrkr.HandleTrigger(context.TODO(), httpTrigger)
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Event", func() {
		When("calling HandlesEvent with the wrong topic", func() {
			subWrkr := &SubscriptionWorker{
				topic: "bad",
			}

			It("should return false", func() {
				Expect(subWrkr.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Topic{Topic: &v1.TopicTriggerContext{Topic: "test"}},
				})).To(BeFalse())
			})
		})

		When("calling HandlesEvent with the correct topic", func() {
			subWrkr := &SubscriptionWorker{
				topic: "test",
			}

			It("should return true", func() {
				Expect(subWrkr.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Topic{Topic: &v1.TopicTriggerContext{Topic: "test"}},
				})).To(BeTrue())
			})
		})

		When("calling HandleEvent", func() {
			It("should call the base grpc workers HandleEvent", func() {
				ctrl := gomock.NewController(GinkgoT())
				hndlr := mock.NewMockAdapter(ctrl)

				By("calling the base grpc handler HandleEvent method")
				hndlr.EXPECT().HandleTrigger(gomock.Any(), gomock.Any()).Times(1)

				subWrkr := &SubscriptionWorker{
					topic:   "test",
					Adapter: hndlr,
				}

				_, err := subWrkr.HandleTrigger(context.TODO(), &v1.TriggerRequest{
					Context: &v1.TriggerRequest_Topic{
						Topic: &v1.TopicTriggerContext{
							Topic: "test",
						},
					},
				})

				Expect(err).ShouldNot(HaveOccurred())
				ctrl.Finish()
			})
		})
	})
})
