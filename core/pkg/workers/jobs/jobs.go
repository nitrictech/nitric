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

package jobs

import (
	"fmt"
	"sync"

	"github.com/nitrictech/nitric/core/pkg/help"
	batchpb "github.com/nitrictech/nitric/core/pkg/proto/batch/v1"
	"github.com/nitrictech/nitric/core/pkg/workers"
)

type JobName = string

type WorkerConnection = workers.WorkerRequestBroker[*batchpb.ServerMessage, *batchpb.ClientMessage]

type JobRequestHandler interface {
	batchpb.JobServer
	HandleJobRequest(request *batchpb.ServerMessage) (*batchpb.ClientMessage, error)
	WorkerCount() int
}

var _ JobRequestHandler = (*JobManager)(nil)

type JobManager struct {
	handlers map[JobName]*WorkerConnection
	lock     sync.RWMutex
}

func (s *JobManager) registerHandler(handler *WorkerConnection, registrationRequest *batchpb.RegistrationRequest) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	jobName := registrationRequest.GetJobName()

	if _, exists := s.handlers[jobName]; exists {
		// Only one job handler allowed per job
		return fmt.Errorf("job handler already registered for job: %s", jobName)
	}

	s.handlers[jobName] = handler

	return nil
}

func (s *JobManager) unregisterHandler(jobName string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.handlers, jobName)
}

func (s *JobManager) WorkerCount() int {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return len(s.handlers)
}

// Subscribe allows the local nitric server to register new subscribers
//
//	called by Nitric applications wishing to subscribe to a topic.
//	the stream establishes communication between that subscriber and this server
func (s *JobManager) HandleJob(handlerConnectionStream batchpb.Job_HandleJobServer) error {
	// The client MUST send registration request on the stream as its first message
	initialRequest, err := handlerConnectionStream.Recv()
	if err != nil {
		return err
	}

	registrationRequest := initialRequest.GetRegistrationRequest()
	if registrationRequest == nil {
		return fmt.Errorf("request received from unregistered subscriber, initial request must be a registration request. %s", help.BugInNitricHelpText())
	}

	handler := workers.NewWorkerRequestBroker[*batchpb.ServerMessage, *batchpb.ClientMessage](handlerConnectionStream)
	if err := s.registerHandler(handler, registrationRequest); err != nil {
		return err
	}

	defer s.unregisterHandler(registrationRequest.GetJobName())

	// send acknowledgement of registration
	err = handlerConnectionStream.Send(&batchpb.ServerMessage{
		Content: &batchpb.ServerMessage_RegistrationResponse{
			RegistrationResponse: &batchpb.RegistrationResponse{},
		},
	})
	if err != nil {
		return err
	}

	err = handler.Run()
	if err != nil {
		return fmt.Errorf("job handler connection broker encountered and error: %w", err)
	}

	return nil
}

// findMatchingSubscriber returns the subscribers for a given topic, or an error if none are found
func (s *JobManager) findMatchingHandler(jobName string) (*WorkerConnection, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	worker, ok := s.handlers[jobName]
	if !ok {
		return nil, fmt.Errorf("no worker registered for job: %s", jobName)
	}

	return worker, nil
}

func (s *JobManager) HandleJobRequest(request *batchpb.ServerMessage) (*batchpb.ClientMessage, error) {
	if request.Id == "" {
		request.Id = workers.GenerateUniqueId()
	}

	messageRequest := request.GetJobRequest()
	if messageRequest == nil {
		return nil, fmt.Errorf("invalid request, expected job request. %s", help.BugInNitricHelpText())
	}
	jobName := messageRequest.GetJobName()

	handler, err := s.findMatchingHandler(jobName)
	if err != nil {
		return nil, err
	}

	resp, err := handler.Send(request)
	if err != nil {
		return nil, err
	}
	return *resp, nil
}

func New() *JobManager {
	return &JobManager{
		handlers: make(map[string]*WorkerConnection),
		lock:     sync.RWMutex{},
	}
}
