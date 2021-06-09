package worker

import (
	"fmt"
	"sync"

	"github.com/nitric-dev/membrane/handler"
)

type HttpWorkerPool struct {
	maxWorkers int
	workerLock sync.Mutex
	workers    []*FaasWorker
}

// Ensure workers implement the trigger handler interface
func (s *HttpWorkerPool) GetTriggerHandler() (handler.TriggerHandler, error) {
	s.workerLock.Lock()
	defer s.workerLock.Unlock()

	if len(s.workers) > 0 {
		return s.workers[0], nil
	} else {
		return nil, fmt.Errorf("No available workers in this pool!")
	}

	return s
}

func (s *HttpWorkerPool) WaitForActiveWorkers(timeout int) error {
	// Dial the child port to see if it's open and ready...
	maxWaitTime := time.Duration(timeout) * time.Second
	// Longer poll times, e.g. 200 milliseconds results in slow lambda cold starts (15s+)
	pollInterval := time.Duration(15) * time.Millisecond

	var waitedTime = time.Duration(0)
	for {
		if s.GetWorkerCount() >= 1 {
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				return fmt.Errorf("Unable to dial child server, does it expose a http server at: %s?", s.childAddress)
			}
		}
	}

	return nil
}

func (s *HttpWorkerPool) getActiveWorkers() int {
	s.workerLock.Lock()
	defer s.workerLock.Unlock()

	return len(s.workers)
}

func (s *HttpWorkerPool) AddWorker(address string) error {
	s.workerLock.Lock()
	defer s.workerLock.Unlock()
	length := len(s.workers)
	if length < s.maxWorkers {
		s.workers[length] = newHttpWorker(address)
		return nil
	}

	return fmt.Errorf("Unable to add worker, Worker pool limit reached!")
}

func NewHttpWorkerPool() *HttpWorkerPool {
	return &HttpWorkerPool{
		maxWorkers: 1,
		workerLock: sync.Mutex{},
		workers: make([]*FaasWorker, 1)
	}
}
