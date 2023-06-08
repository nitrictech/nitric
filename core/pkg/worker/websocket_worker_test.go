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
	mock "github.com/nitrictech/nitric/core/mocks/adapter"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WebsocketWorker", func() {
	Context("HandlesTrigger", func() {
		websockWorker := NewWebsocketWorker(nil, &WebsocketWorkerOptions{
			Socket: "test",
			Event:  v1.WebsocketEvent_Connect,
		})

		When("calling HandlesTrigger without a websocket context", func() {
			It("should return false", func() {
				Expect(websockWorker.HandlesTrigger(&v1.TriggerRequest{})).To(BeFalse())
			})
		})

		When("calling HandlesTrigger with a websocket context but non-matching event", func() {
			It("should return false", func() {
				Expect(websockWorker.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Websocket{
						Websocket: &v1.WebsocketTriggerContext{
							Socket: "test",
							Event:  v1.WebsocketEvent_Disconnect,
						},
					},
				})).To(BeFalse())
			})
		})

		When("calling HandlesTrigger with a mismatched socket name", func() {
			It("should return true", func() {
				Expect(websockWorker.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Websocket{
						Websocket: &v1.WebsocketTriggerContext{
							Socket: "bad-socket",
							Event:  v1.WebsocketEvent_Connect,
						},
					},
				})).To(BeFalse())
			})
		})

		When("calling HandlesTrigger with a valid websocket context", func() {
			It("should return true", func() {
				Expect(websockWorker.HandlesTrigger(&v1.TriggerRequest{
					Context: &v1.TriggerRequest_Websocket{
						Websocket: &v1.WebsocketTriggerContext{
							Socket: "test",
							Event:  v1.WebsocketEvent_Connect,
						},
					},
				})).To(BeTrue())
			})
		})
	})

	Context("HandleTrigger", func() {
		When("calling HandleTrigger without a websocket context", func() {
			websockWorker := NewWebsocketWorker(nil, &WebsocketWorkerOptions{
				Socket: "test",
				Event:  v1.WebsocketEvent_Connect,
			})

			It("should return an error", func() {
				_, err := websockWorker.HandleTrigger(context.TODO(), &v1.TriggerRequest{})
				Expect(err).Should(HaveOccurred())
			})
		})

		When("calling HandleTrigger with a websocket context", func() {
			ctrl := gomock.NewController(GinkgoT())
			hndlr := mock.NewMockAdapter(ctrl)
			websockWorker := NewWebsocketWorker(hndlr, &WebsocketWorkerOptions{
				Socket: "test",
				Event:  v1.WebsocketEvent_Connect,
			})

			trigger := &v1.TriggerRequest{
				Context: &v1.TriggerRequest_Websocket{
					Websocket: &v1.WebsocketTriggerContext{
						Socket: "test",
						Event:  v1.WebsocketEvent_Connect,
					},
				},
			}

			It("should not return an error", func() {
				By("calling it's adapter")
				hndlr.EXPECT().HandleTrigger(gomock.Any(), trigger).Times(1).Return(&v1.TriggerResponse{
					Context: &v1.TriggerResponse_Websocket{
						Websocket: &v1.WebsocketResponseContext{
							Success: true,
						},
					},
				}, nil)

				By("successfully handling the trigger")
				_, err := websockWorker.HandleTrigger(context.TODO(), trigger)
				Expect(err).ShouldNot(HaveOccurred())

				ctrl.Finish()
			})
		})
	})
})
