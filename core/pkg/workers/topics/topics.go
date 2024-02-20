// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package topics

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/nitrictech/nitric/core/pkg/help"
	topicspb "github.com/nitrictech/nitric/core/pkg/proto/topics/v1"
	"github.com/nitrictech/nitric/core/pkg/workers"
	"golang.org/x/sync/errgroup"
)

type TopicName = string

type WorkerConnection = workers.WorkerRequestBroker[*topicspb.ServerMessage, *topicspb.ClientMessage]

type SubscriptionRequestHandler interface {
	topicspb.SubscriberServer
	HandleRequest(request *topicspb.ServerMessage) (*topicspb.ClientMessage, error)
	WorkerCount() int
}

type SubscriberManager struct {
	subscriberMap map[TopicName][]*WorkerConnection
	lock          sync.RWMutex
}

func (s *SubscriberManager) registerSubscriber(subscriber *WorkerConnection, registrationRequest *topicspb.RegistrationRequest) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	topicName := registrationRequest.GetTopicName()

	if _, exists := s.subscriberMap[topicName]; !exists {
		s.subscriberMap[topicName] = make([]*WorkerConnection, 0, 1)
	}

	s.subscriberMap[topicName] = append(s.subscriberMap[topicName], subscriber)

	return nil
}

func (s *SubscriberManager) unregisterSubscriber(topicName string, subscriber *WorkerConnection) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.subscriberMap[topicName] = slices.DeleteFunc[[]*WorkerConnection](s.subscriberMap[topicName], func(wc *WorkerConnection) bool {
		return wc == subscriber
	})

	if len(s.subscriberMap[topicName]) == 0 {
		delete(s.subscriberMap, topicName)
	}
}

func (s *SubscriberManager) WorkerCount() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	var count int
	for _, subscribers := range s.subscriberMap {
		count += len(subscribers)
	}

	return count
}

// Subscribe allows the local nitric server to register new subscribers
//
//	called by Nitric applications wishing to subscribe to a topic.
//	the stream establishes communication between that subscriber and this server
func (s *SubscriberManager) Subscribe(subscriberConnectionStream topicspb.Subscriber_SubscribeServer) error {
	// The client MUST send registration request on the stream as its first message
	initialRequest, err := subscriberConnectionStream.Recv()
	if err != nil {
		return err
	}

	registrationRequest := initialRequest.GetRegistrationRequest()
	if registrationRequest == nil {
		return fmt.Errorf("request received from unregistered subscriber, initial request must be a registration request. %s", help.BugInNitricHelpText())
	}

	subscriber := workers.NewWorkerRequestBroker[*topicspb.ServerMessage, *topicspb.ClientMessage](subscriberConnectionStream)
	if err := s.registerSubscriber(subscriber, registrationRequest); err != nil {
		return err
	}

	defer s.unregisterSubscriber(registrationRequest.GetTopicName(), subscriber)

	// send acknowledgement of registration
	err = subscriberConnectionStream.Send(&topicspb.ServerMessage{
		Content: &topicspb.ServerMessage_RegistrationResponse{
			RegistrationResponse: &topicspb.RegistrationResponse{},
		},
	})
	if err != nil {
		return err
	}

	err = subscriber.Run()
	if err != nil {
		return fmt.Errorf("subscriber connection broker encountered and error: %w", err)
	}

	return nil
}

// findMatchingSubscriber returns the subscribers for a given topic, or an error if none are found
func (s *SubscriberManager) findMatchingSubscriber(topicName string) ([]*WorkerConnection, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	workers, ok := s.subscriberMap[topicName]
	if !ok || len(workers) == 0 {
		return nil, fmt.Errorf("no workers registered for topic subscription: %s", topicName)
	}

	return workers, nil
}

// ForwardRequestToSubscribers forwards an event to all subscribers for a given topic
// returns a slice of errors encountered while forwarding the event
func ForwardRequestToSubscribers(subscribers []*WorkerConnection, request *topicspb.ServerMessage) (bool, error) {
	success := true
	errs, _ := errgroup.WithContext(context.Background())

	for _, subscriber := range subscribers {
		footLongSub := subscriber
		errs.Go(func() error {
			resp, err := footLongSub.Send(request)
			if err != nil {
				return err
			} else if !(*resp).GetMessageResponse().GetSuccess() {
				success = false
			}
			return nil
		})
	}

	err := errs.Wait()
	if err != nil {
		return false, fmt.Errorf("errors occurred handling subscription %v", errs)
	}

	return success, nil
}

func (s *SubscriberManager) HandleRequest(request *topicspb.ServerMessage) (*topicspb.ClientMessage, error) {
	if request.Id == "" {
		request.Id = workers.GenerateUniqueId()
	}

	messageRequest := request.GetMessageRequest()
	if messageRequest == nil {
		return nil, fmt.Errorf("invalid request, expected message request. %s", help.BugInNitricHelpText())
	}
	topicName := messageRequest.GetTopicName()

	subscribers, err := s.findMatchingSubscriber(topicName)
	if err != nil {
		return nil, err
	}

	success, err := ForwardRequestToSubscribers(subscribers, request)
	if err != nil {
		return nil, err
	}

	// Compute units support multiple subscribers for a single topic.
	//
	// This makes topics are a unique case where the worker's response isn't directly returned.
	// instead we return a response indicating success or failure of all subscribers.
	//
	// Since subscription workers should always be idempotent, this ensures a failure in any subscriber
	// will trigger a retry of the event processing.
	return &topicspb.ClientMessage{
		Content: &topicspb.ClientMessage_MessageResponse{
			MessageResponse: &topicspb.MessageResponse{
				Success: success,
			},
		},
	}, nil
}

func New() *SubscriberManager {
	return &SubscriberManager{
		subscriberMap: make(map[string][]*WorkerConnection),
		lock:          sync.RWMutex{},
	}
}
