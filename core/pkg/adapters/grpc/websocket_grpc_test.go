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

package grpc_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/core/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/websocket/v1"

	mock "github.com/nitrictech/nitric/core/mocks/websocket"
)

var _ = Describe("GRPC Websockets", func() {
	Context("Send", func() {
		When("plugin not registered", func() {
			ws := grpc.NewWebsocketServiceServer(nil)
			resp, err := ws.Send(context.Background(), &v1.WebsocketSendRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Websocket plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("calling send with an invalid request", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockWs := mock.NewMockWebsocketService(ctrl)

			ws := grpc.NewWebsocketServiceServer(mockWs)
			_, err := ws.Send(context.Background(), &v1.WebsocketSendRequest{})

			It("Should report an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("calling send with a valid request", func() {
			When("send returns an error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockWs := mock.NewMockWebsocketService(ctrl)
				data := []byte("test")

				ws := grpc.NewWebsocketServiceServer(mockWs)

				It("Should fail", func() {
					By("calling the provided plugin")
					mockWs.EXPECT().Send(gomock.Any(), "test", "test-connection-id", data).Return(fmt.Errorf("mock-error"))

					_, err := ws.Send(context.Background(), &v1.WebsocketSendRequest{
						Socket:       "test",
						ConnectionId: "test-connection-id",
						Data:         data,
					})

					By("returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})

			When("send returns a response", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockWs := mock.NewMockWebsocketService(ctrl)
				data := []byte("test")

				ws := grpc.NewWebsocketServiceServer(mockWs)

				It("Should succeed", func() {
					By("calling the provided plugin")
					mockWs.EXPECT().Send(gomock.Any(), "test", "test-connection-id", data).Return(nil)

					resp, err := ws.Send(context.Background(), &v1.WebsocketSendRequest{
						Socket:       "test",
						ConnectionId: "test-connection-id",
						Data:         data,
					})

					By("not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("returning a non-nil response")
					Expect(resp).ToNot(BeNil())
				})
			})
		})
	})

	Context("Close", func() {
		When("plugin not registered", func() {
			ws := grpc.NewWebsocketServiceServer(nil)
			resp, err := ws.Close(context.Background(), &v1.WebsocketCloseRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Websocket plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("calling close with an invalid request", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockWs := mock.NewMockWebsocketService(ctrl)

			ws := grpc.NewWebsocketServiceServer(mockWs)
			_, err := ws.Close(context.Background(), &v1.WebsocketCloseRequest{})

			It("Should report an error", func() {
				Expect(err).Should(HaveOccurred())
			})
		})

		When("calling close with a valid request", func() {
			When("send returns an error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockWs := mock.NewMockWebsocketService(ctrl)
				ws := grpc.NewWebsocketServiceServer(mockWs)

				It("Should fail", func() {
					By("calling the provided plugin")
					mockWs.EXPECT().Close(gomock.Any(), "test", "test-connection-id").Return(fmt.Errorf("mock-error"))

					_, err := ws.Close(context.Background(), &v1.WebsocketCloseRequest{
						Socket:       "test",
						ConnectionId: "test-connection-id",
					})

					By("returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})

			When("send returns a response", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockWs := mock.NewMockWebsocketService(ctrl)
				ws := grpc.NewWebsocketServiceServer(mockWs)

				It("Should succeed", func() {
					By("calling the provided plugin")
					mockWs.EXPECT().Close(gomock.Any(), "test", "test-connection-id").Return(nil)

					resp, err := ws.Close(context.Background(), &v1.WebsocketCloseRequest{
						Socket:       "test",
						ConnectionId: "test-connection-id",
					})

					By("not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("returning a non-nil response")
					Expect(resp).ToNot(BeNil())
				})
			})
		})
	})
})
