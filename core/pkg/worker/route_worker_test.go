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

var _ = Describe("RouteWorker", func() {
	Context("Http", func() {
		rWrkr := &RouteWorker{
			methods: []string{"GET"},
			path:    "/test/:param",
		}

		When("calling HandlesTrigger with bad path", func() {
			It("should return false", func() {
				Expect(rWrkr.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "GET",
							Path:   "/test/",
						},
					},
				})).To(BeFalse())
			})
		})

		When("calling HandlesTrigger with bad method", func() {
			It("should return false", func() {
				Expect(rWrkr.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "POST",
							Path:   "/test/test",
						},
					},
				})).To(BeFalse())
			})
		})

		When("calling HandlesHttpRequest with matching path and method", func() {
			It("should return true", func() {
				Expect(rWrkr.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "GET",
							Path:   "/test/test",
						},
					},
				})).To(BeTrue())
			})
		})

		When("calling HandleHttpRequest", func() {
			It("should call the base grpc workers HandleEvent with augmented trigger", func() {
				ctrl := gomock.NewController(GinkgoT())
				hndlr := mock.NewMockAdapter(ctrl)

				By("calling the base grpc handler HandleEvent method")
				hndlr.EXPECT().HandleTrigger(context.TODO(), &v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "GET",
							Path:   "/test/name",
							PathParams: map[string]string{
								"param": "name",
							},
						},
					},
				}).Times(1)

				subWrkr := &RouteWorker{
					methods: []string{"GET"},
					path:    "/test/:param",
					Adapter: hndlr,
				}

				_, err := subWrkr.HandleTrigger(context.TODO(), &v1.TriggerRequest{
					Context: &v1.TriggerRequest_Http{
						Http: &v1.HttpTriggerContext{
							Method: "GET",
							Path:   "/test/name",
						},
					},
				})

				Expect(err).ShouldNot(HaveOccurred())
				ctrl.Finish()
			})
		})
	})

	Context("Event", func() {
		eventTrigger := &v1.TriggerRequest{
			Context: &v1.TriggerRequest_Topic{
				Topic: &v1.TopicTriggerContext{},
			},
		}
		When("calling HandlesTrigger wth an Event", func() {
			rWrkr := &RouteWorker{}

			It("should return false", func() {
				Expect(rWrkr.HandlesTrigger(eventTrigger)).To(BeFalse())
			})
		})

		When("calling HandleEvent", func() {
			subWrkr := &RouteWorker{}

			It("should return an error", func() {
				_, err := subWrkr.HandleTrigger(context.TODO(), eventTrigger)

				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
