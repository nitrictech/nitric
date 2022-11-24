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

	mock_events "github.com/nitrictech/nitric/mocks/plugins/events"
	"github.com/nitrictech/nitric/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/pkg/api/nitric/v1"
)

var _ = Describe("Event Service gRPC Adapter", func() {
	Context("Publish", func() {
		When("No request id is provided", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockService := mock_events.NewMockEventService(ctrl)
			eventServer := grpc.NewEventServiceServer(mockService)

			It("should successfully publish the event", func() {
				By("Calling the provided service")
				mockService.EXPECT().Publish(gomock.Any(), "test-topic", 0, gomock.Any()).Return(nil).Times(1)

				response, err := eventServer.Publish(context.Background(), &v1.EventPublishRequest{
					Topic: "test-topic",
					Event: &v1.NitricEvent{
						Id:          "",
						PayloadType: "",
						Payload:     nil,
					},
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Autogenerating the event ID")
				Expect(response.Id).ToNot(BeEmpty())
			})
		})

		When("A request id is provided", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockService := mock_events.NewMockEventService(ctrl)
			eventServer := grpc.NewEventServiceServer(mockService)

			It("Should successfully be handled", func() {
				By("Calling the provided service")
				mockService.EXPECT().Publish(context.TODO(), "test-topic", 0, gomock.Any()).Return(nil).Times(1)

				response, err := eventServer.Publish(context.Background(), &v1.EventPublishRequest{
					Topic: "test-topic",
					Event: &v1.NitricEvent{
						Id:          "test-id",
						PayloadType: "",
						Payload:     nil,
					},
				})

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning the provided id")
				Expect(response.Id).To(Equal("test-id"))
			})
		})
	})
})
