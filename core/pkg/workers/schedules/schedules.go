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

package schedules

import (
	"fmt"
	"sync"

	schedulespb "github.com/nitrictech/nitric/core/pkg/proto/schedules/v1"
	"github.com/nitrictech/nitric/core/pkg/workers"
)

type ScheduleName = string

type WorkerConnection = workers.WorkerRequestBroker[*schedulespb.ServerMessage, *schedulespb.ClientMessage]

type ScheduleRequestHandler interface {
	schedulespb.SchedulesServer
	HandleRequest(request *schedulespb.ServerMessage) (*schedulespb.ClientMessage, error)
	WorkerCount() int
}

type ScheduleWorkerManager struct {
	workerMap map[ScheduleName]*WorkerConnection
	mutex     sync.RWMutex
}

var _ schedulespb.SchedulesServer = &ScheduleWorkerManager{}

func (s *ScheduleWorkerManager) registerSchedule(scheduleWorker *WorkerConnection, request *schedulespb.RegistrationRequest) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	scheduleName := request.GetScheduleName()

	if _, exists := s.workerMap[scheduleName]; exists {
		return fmt.Errorf("schedule already registered with name: %s", scheduleName)
	}

	s.workerMap[scheduleName] = scheduleWorker

	return nil
}

func (s *ScheduleWorkerManager) Schedule(stream schedulespb.Schedules_ScheduleServer) error {
	initRequest, err := stream.Recv()
	if err != nil {
		return err
	}

	if initRequest.GetRegistrationRequest() == nil {
		return fmt.Errorf("first request must be an init request")
	}

	worker := workers.NewWorkerRequestBroker[*schedulespb.ServerMessage, *schedulespb.ClientMessage](stream)
	if err := s.registerSchedule(worker, initRequest.GetRegistrationRequest()); err != nil {
		return err
	}

	// send acknowledgement of registration
	err = stream.Send(&schedulespb.ServerMessage{
		Content: &schedulespb.ServerMessage_RegistrationResponse{
			RegistrationResponse: &schedulespb.RegistrationResponse{},
		},
	})
	if err != nil {
		return err
	}

	return worker.Run()
}

func (s *ScheduleWorkerManager) HandleRequest(request *schedulespb.ServerMessage) (*schedulespb.ClientMessage, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if request.Id == "" {
		request.Id = workers.GenerateUniqueId()
	}

	worker, ok := s.workerMap[request.GetIntervalRequest().GetScheduleName()]

	if !ok {
		return nil, fmt.Errorf("no worker registered for schedule: %s", request.GetIntervalRequest().GetScheduleName())
	}

	resp, err := worker.Send(request)

	return *resp, err
}

func (s *ScheduleWorkerManager) WorkerCount() int {
	return len(s.workerMap)
}

func New() *ScheduleWorkerManager {
	return &ScheduleWorkerManager{
		workerMap: make(map[string]*WorkerConnection),
		mutex:     sync.RWMutex{},
	}
}
