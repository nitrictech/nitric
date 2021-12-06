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

	"github.com/nitrictech/nitric/pkg/adapters/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	v1 "github.com/nitrictech/nitric/interfaces/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/events"
)

type MockEventService struct {
	PublishError error
	PublishTopic string
	PublishEvent *events.NitricEvent

	TopicList      []string
	TopicListError error
}

func (m *MockEventService) Publish(topic string, event *events.NitricEvent) error {
	m.PublishTopic = topic
	m.PublishEvent = event
	return m.PublishError
}

func (m *MockEventService) ListTopics() ([]string, error) {
	return m.TopicList, m.TopicListError
}

var _ = Describe("Event Service gRPC Adapter", func() {
	Context("Publish", func() {
		When("No request id is provided", func() {
			mockService := &MockEventService{
				PublishError:   nil,
				TopicListError: nil,
			}

			eventServer := grpc.NewEventServiceServer(mockService)
			response, err := eventServer.Publish(context.Background(), &v1.EventPublishRequest{
				Topic: "test-topic",
				Event: &v1.NitricEvent{
					Id:          "",
					PayloadType: "",
					Payload:     nil,
				},
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should pass the generated id to the implementing service plugin", func() {
				Expect(mockService.PublishEvent.ID).ToNot(BeEmpty())
			})

			It("Should return the generated ID", func() {
				Expect(response.Id).ToNot(BeEmpty())
			})
		})

		When("A request id is provided", func() {
			mockService := &MockEventService{
				PublishError:   nil,
				TopicListError: nil,
			}

			eventServer := grpc.NewEventServiceServer(mockService)
			response, err := eventServer.Publish(context.Background(), &v1.EventPublishRequest{
				Topic: "test-topic",
				Event: &v1.NitricEvent{
					Id:          "test-id",
					PayloadType: "",
					Payload:     nil,
				},
			})

			It("Should not return an error", func() {
				Expect(err).To(BeNil())
			})

			It("Should pass the provided id to the implementing service plugin", func() {
				Expect(mockService.PublishEvent.ID).To(Equal("test-id"))
			})

			It("Should return the provided ID", func() {
				Expect(response.Id).To(Equal("test-id"))
			})
		})
	})
})
