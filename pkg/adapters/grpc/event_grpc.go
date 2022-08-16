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

package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"

	pb "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/events"
)

// GRPC Interface for registered Nitric events Plugins
type EventServiceServer struct {
	pb.UnimplementedEventServiceServer
	eventPlugin events.EventService
}

func (s *EventServiceServer) checkPluginRegistered() error {
	if s.eventPlugin == nil {
		return NewPluginNotRegisteredError("Event")
	}

	return nil
}

func (s *EventServiceServer) Publish(ctx context.Context, req *pb.EventPublishRequest) (*pb.EventPublishResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "EventService.Publish", err)
	}

	// auto generate an ID if we did not receive one
	ID := req.GetEvent().GetId()
	if ID == "" {
		ID = uuid.New().String()
	}

	event := &events.NitricEvent{
		ID:          ID,
		PayloadType: req.GetEvent().GetPayloadType(),
		Payload:     req.GetEvent().GetPayload().AsMap(),
	}
	if err := s.eventPlugin.Publish(req.GetTopic(), event); err == nil {
		return &pb.EventPublishResponse{
			Id: ID,
		}, nil
	} else {
		return nil, NewGrpcError("EventService.Publish", err)
	}
}

func NewEventServiceServer(eventsPlugin events.EventService) pb.EventServiceServer {
	return &EventServiceServer{
		eventPlugin: eventsPlugin,
	}
}

type TopicServiceServer struct {
	pb.UnimplementedTopicServiceServer
	eventPlugin events.EventService
}

func (s *TopicServiceServer) checkPluginRegistered() error {
	if s.eventPlugin == nil {
		return NewPluginNotRegisteredError("Event")
	}

	return nil
}

func (s *TopicServiceServer) List(ctx context.Context, req *pb.TopicListRequest) (*pb.TopicListResponse, error) {
	if err := s.checkPluginRegistered(); err != nil {
		return nil, err
	}

	if err := req.ValidateAll(); err != nil {
		return nil, newGrpcErrorWithCode(codes.InvalidArgument, "TopicService.List", err)
	}

	if res, err := s.eventPlugin.ListTopics(); err == nil {
		topics := make([]*pb.NitricTopic, len(res))
		for i, topicName := range res {
			topics[i] = &pb.NitricTopic{
				Name: topicName,
			}
		}

		return &pb.TopicListResponse{
			Topics: topics,
		}, nil
	} else {
		return nil, NewGrpcError("TopicService.List", err)
	}
}

func NewTopicServiceServer(eventService events.EventService) pb.TopicServiceServer {
	// The external topic/event interfaces are separate. Internally, they're fulfilled together,
	// so the event plugin is all that's needed for both the Event and Topic servers currently.
	return &TopicServiceServer{
		eventPlugin: eventService,
	}
}
