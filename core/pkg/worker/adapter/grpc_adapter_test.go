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

package adapter

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_nitric "github.com/nitrictech/nitric/core/mocks/nitric"
	mock_sync "github.com/nitrictech/nitric/core/mocks/sync"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

var _ = Describe("grpcAdapter", func() {
	Context("newTicket", func() {
		When("calling newTicket", func() {
			ctrl := gomock.NewController(GinkgoT())
			lck := mock_sync.NewMockLocker(ctrl)
			wkr := &GrpcAdapter{
				responseQueueLock: lck,
				responseQueue:     make(map[string]chan *v1.TriggerResponse),
			}

			It("should thread safely add a new ticket to the queue", func() {
				By("Locking the response queue lock")
				lck.EXPECT().Lock().Times(1)

				By("Unlocking the response queue lock")
				lck.EXPECT().Unlock().Times(1)

				By("returning a new response UUID")
				id, ch := wkr.newTicket()
				_, err := uuid.Parse(id)
				Expect(err).ShouldNot(HaveOccurred())

				By("returning a new trigger response channel")
				Expect(ch).ToNot(BeNil())

				By("storing the new ticket in the response queue")
				Expect(wkr.responseQueue[id]).To(Equal(ch))

				ctrl.Finish()
			})
		})
	})

	Context("resolveTicket", func() {
		When("resolving a ticket id that does not exist", func() {
			ctrl := gomock.NewController(GinkgoT())
			lck := mock_sync.NewMockLocker(ctrl)
			wkr := &GrpcAdapter{
				responseQueueLock: lck,
				responseQueue:     make(map[string]chan *v1.TriggerResponse),
			}

			It("should thread safely add a new ticket to the queue", func() {
				By("Locking the response queue lock")
				lck.EXPECT().Lock().Times(1)

				By("Unlocking the response queue lock")
				lck.EXPECT().Unlock().Times(1)

				ch, err := wkr.resolveTicket("fake-ticket")

				By("returning an error")
				Expect(err).Should(HaveOccurred())

				By("returning a nil channel")
				Expect(ch).To(BeNil())

				ctrl.Finish()
			})
		})

		When("resolving a ticket id that exists", func() {
			ctrl := gomock.NewController(GinkgoT())
			lck := mock_sync.NewMockLocker(ctrl)
			ch := make(chan *v1.TriggerResponse)
			wkr := &GrpcAdapter{
				responseQueueLock: lck,
				responseQueue: map[string]chan *v1.TriggerResponse{
					"test": ch,
				},
			}

			It("should thread safely add a new ticket to the queue", func() {
				By("Locking the response queue lock")
				lck.EXPECT().Lock().Times(1)

				By("Unlocking the response queue lock")
				lck.EXPECT().Unlock().Times(1)

				rch, err := wkr.resolveTicket("test")

				By("not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("returning a the channel")
				Expect(rch).To(Equal(ch))

				By("removing the ticket from the response queue")
				Expect(wkr.responseQueue["test"]).To(BeNil())

				ctrl.Finish()
			})
		})
	})

	Context("send", func() {
		// TODO: possible remove?
	})

	Context("Listen", func() {
		When("receiving valid trigger responses", func() {
			ctrl := gomock.NewController(GinkgoT())
			stream := mock_nitric.NewMockFaasService_TriggerStreamServer(ctrl)
			errChan := make(chan error)
			respChan := make(chan *v1.TriggerResponse)
			mockResp := &v1.TriggerResponse{}
			wkr := &GrpcAdapter{
				responseQueueLock: &sync.Mutex{},
				responseQueue: map[string]chan *v1.TriggerResponse{
					"test": respChan,
				},
				stream: stream,
			}

			It("should process the response", func() {
				By("Sending a client message then closing the channel")
				gomock.InOrder(
					stream.EXPECT().Recv().Return(&v1.ClientMessage{
						Id: "test",
						Content: &v1.ClientMessage_TriggerResponse{
							TriggerResponse: mockResp,
						},
					}, nil),
					stream.EXPECT().Recv().Return(nil, io.EOF),
				)

				go func() {
					wkr.Start(errChan)
				}()

				resp := <-respChan
				err := <-errChan

				By("receiving the send response")
				Expect(resp).To(Equal(mockResp))

				By("receiving the channel close")
				Expect(err).To(Equal(io.EOF))

				ctrl.Finish()
			})
		})

		When("receiving invalid trigger responses", func() {
			ctrl := gomock.NewController(GinkgoT())
			stream := mock_nitric.NewMockFaasService_TriggerStreamServer(ctrl)
			errChan := make(chan error)
			mockResp := &v1.TriggerResponse{}
			wkr := &GrpcAdapter{
				responseQueueLock: &sync.Mutex{},
				responseQueue: map[string]chan *v1.TriggerResponse{
					"test": nil,
				},
				stream: stream,
			}

			It("should exit with error", func() {
				By("Sending a client message then closing the channel")
				gomock.InOrder(
					stream.EXPECT().Recv().Return(&v1.ClientMessage{
						Id: "bad-id",
						Content: &v1.ClientMessage_TriggerResponse{
							TriggerResponse: mockResp,
						},
					}, nil),
				)

				go func() {
					wkr.Start(errChan)
				}()

				err := <-errChan

				By("receiving error")
				Expect(err).Should(HaveOccurred())

				ctrl.Finish()
			})
		})

		When("receiving an error", func() {
			ctrl := gomock.NewController(GinkgoT())
			stream := mock_nitric.NewMockFaasService_TriggerStreamServer(ctrl)
			mockErr := fmt.Errorf("mock error")

			wkr := &GrpcAdapter{
				responseQueueLock: &sync.Mutex{},
				responseQueue:     make(map[string]chan *v1.TriggerResponse),
				stream:            stream,
			}

			It("should exit", func() {
				errChan := make(chan error)

				go func() {
					wkr.Start(errChan)
				}()

				By("the client sending io.EOF")
				stream.EXPECT().Recv().Return(nil, mockErr)

				err := <-errChan

				By("passing through the error")
				Expect(err).To(Equal(mockErr))

				ctrl.Finish()
			})
		})

		When("receiving io.EOF", func() {
			ctrl := gomock.NewController(GinkgoT())
			stream := mock_nitric.NewMockFaasService_TriggerStreamServer(ctrl)

			wkr := &GrpcAdapter{
				responseQueueLock: &sync.Mutex{},
				responseQueue:     make(map[string]chan *v1.TriggerResponse),
				stream:            stream,
			}

			It("should exit", func() {
				errChan := make(chan error)

				go func() {
					wkr.Start(errChan)
				}()

				By("the client sending io.EOF")
				stream.EXPECT().Recv().Return(nil, io.EOF)

				err := <-errChan

				By("passing through the error")
				Expect(err).To(Equal(io.EOF))

				ctrl.Finish()
			})
		})
	})

	Context("HandleHttpRequest", func() {
		When("the worker connection responds with an error", func() {
			ctrl := gomock.NewController(GinkgoT())
			stream := mock_nitric.NewMockFaasService_TriggerStreamServer(ctrl)
			mockErr := fmt.Errorf("mock error")
			wkr := &GrpcAdapter{
				responseQueueLock: &sync.Mutex{},
				responseQueue:     make(map[string]chan *v1.TriggerResponse),
				stream:            stream,
			}

			It("should return an error", func() {
				By("gRPC returning an error")
				stream.EXPECT().Send(gomock.Any()).Return(mockErr)

				By("returning the error")
				_, err := wkr.HandleTrigger(context.TODO(), &v1.TriggerRequest{})
				Expect(err).To(Equal(mockErr))
			})
		})

		PWhen("the worker successfully responds", func() {
			// TODO
		})
	})

	Context("HandleEvent", func() {
		When("the worker connection responds with an error", func() {
			ctrl := gomock.NewController(GinkgoT())
			stream := mock_nitric.NewMockFaasService_TriggerStreamServer(ctrl)
			mockErr := fmt.Errorf("mock error")
			wkr := &GrpcAdapter{
				responseQueueLock: &sync.Mutex{},
				responseQueue:     make(map[string]chan *v1.TriggerResponse),
				stream:            stream,
			}

			It("should return an error", func() {
				By("gRPC returning an error")
				stream.EXPECT().Send(gomock.Any()).Return(mockErr)

				By("returning the error")
				_, err := wkr.HandleTrigger(context.TODO(), &v1.TriggerRequest{})
				Expect(err).To(Equal(mockErr))
			})
		})

		PWhen("the worker successfully responds", func() {
			// TODO
		})
	})
})
