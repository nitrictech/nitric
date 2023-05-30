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

var _ = Describe("ScheduleWorker", func() {
	Context("HandlesTrigger", func() {
		scheduleWorker := &ScheduleWorker{
			key: "my-test-worker",
		}

		When("calling HandlesTrigger without a TopicContext", func() {
			It("should return false", func() {
				handles := scheduleWorker.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{},
					},
				})

				Expect(handles).To(BeFalse())
			})
		})

		When("calling HandlesTrigger with a TopicContext but the wrong topic", func() {
			It("should return false", func() {
				handles := scheduleWorker.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Topic{
						Topic: &v1.TopicTriggerContext{
							Topic: "bad-topic",
						},
					},
				})
				Expect(handles).To(BeFalse())
			})
		})

		When("calling HandlesTrigger with a TopicContext with the correct topic", func() {
			It("should return true", func() {
				handles := scheduleWorker.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Topic{
						Topic: &v1.TopicTriggerContext{
							Topic: "my-test-worker",
						},
					},
				})
				Expect(handles).To(BeTrue())
			})
		})
	})

	Context("HandleTrigger", func() {
		When("calling HandleTrigger without a topic context", func() {
			scheduleWorker := &ScheduleWorker{
				key: "my-test-worker",
			}

			_, err := scheduleWorker.HandleTrigger(context.TODO(), &v1.TriggerRequest{})

			It("Should return an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("calling HandleTrigger with a topic that does match", func() {
			It("should call the adapter with the Trigger", func() {
				ctrl := gomock.NewController(GinkgoT())
				hndlr := mock.NewMockAdapter(ctrl)
				trigger := &v1.TriggerRequest{
					Context: &v1.TriggerRequest_Topic{
						Topic: &v1.TopicTriggerContext{
							Topic: "test",
						},
					},
				}

				By("calling the base grpc handler HandleEvent method")
				hndlr.EXPECT().HandleTrigger(gomock.Any(), trigger).Times(1)

				subWrkr := &SubscriptionWorker{
					topic:   "test",
					Adapter: hndlr,
				}

				_, err := subWrkr.HandleTrigger(context.TODO(), trigger)

				Expect(err).ShouldNot(HaveOccurred())
				ctrl.Finish()
			})
		})
	})
})
