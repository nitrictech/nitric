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

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/structpb"

	mock_queue "github.com/nitrictech/nitric/mocks/queue"
	"github.com/nitrictech/nitric/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/queue"
)

var _ = Describe("GRPC Queue", func() {
	Context("Send", func() {
		When("plugin not registered", func() {
			ss := &grpc.QueueServiceServer{}
			resp, err := ss.Send(context.Background(), &v1.QueueSendRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Queue plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_queue.NewMockQueueService(g)
			resp, err := grpc.NewQueueServiceServer(mockSS).Send(context.Background(), &v1.QueueSendRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid QueueSendRequest.Queue: value does not match regex pattern"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_queue.NewMockQueueService(g)

			mockSS.EXPECT().Send(gomock.Any(), "job", queue.NitricTask{
				ID:          "tsk",
				PayloadType: "thing",
				Payload: map[string]interface{}{
					"x": "y",
				},
			}).Return(nil)

			resp, err := grpc.NewQueueServiceServer(mockSS).Send(context.Background(), &v1.QueueSendRequest{
				Queue: "job",
				Task: &v1.NitricTask{
					Id:          "tsk",
					PayloadType: "thing",
					Payload: &structpb.Struct{Fields: map[string]*structpb.Value{
						"x": structpb.NewStringValue("y"),
					}},
				},
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(resp.String()).To(Equal(""))
			})
		})
	})

	Context("Receive", func() {
		When("plugin not registered", func() {
			ss := &grpc.QueueServiceServer{}
			resp, err := ss.Receive(context.Background(), &v1.QueueReceiveRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Queue plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_queue.NewMockQueueService(g)
			resp, err := grpc.NewQueueServiceServer(mockSS).Receive(context.Background(), &v1.QueueReceiveRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid QueueReceiveRequest.Queue: value does not match regex pattern"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_queue.NewMockQueueService(g)

			one := uint32(1)
			mockSS.EXPECT().Receive(gomock.Any(), queue.ReceiveOptions{
				QueueName: "job",
				Depth:     &one,
			}).Return([]queue.NitricTask{
				{
					ID:          "tsk",
					PayloadType: "food",
					Payload: map[string]interface{}{
						"ff": "88",
					},
				},
			}, nil)

			resp, err := grpc.NewQueueServiceServer(mockSS).Receive(context.Background(), &v1.QueueReceiveRequest{
				Queue: "job",
				Depth: int32(1),
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(resp.Tasks[0].Id).To(Equal("tsk"))
				Expect(resp.Tasks[0].PayloadType).To(Equal("food"))
			})
		})
	})

	Context("Complete", func() {
		When("plugin not registered", func() {
			ss := &grpc.QueueServiceServer{}
			resp, err := ss.Complete(context.Background(), &v1.QueueCompleteRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Queue plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_queue.NewMockQueueService(g)
			resp, err := grpc.NewQueueServiceServer(mockSS).Complete(context.Background(), &v1.QueueCompleteRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid QueueCompleteRequest.Queue: value does not match regex pattern"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_queue.NewMockQueueService(g)

			mockSS.EXPECT().Complete(gomock.Any(), "job", "45").Return(nil)

			resp, err := grpc.NewQueueServiceServer(mockSS).Complete(context.Background(), &v1.QueueCompleteRequest{
				Queue:   "job",
				LeaseId: "45",
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(resp.String()).To(Equal(""))
			})
		})
	})
})
