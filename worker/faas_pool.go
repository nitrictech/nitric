package worker

import (
	"fmt"
	"sync"

	"github.com/nitric-dev/membrane/handler"
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
)


type FaasWorkerPool struct {
	maxWorkers int
	workerLock sync.Mutex
	workers []*FaasWorker
}

// Ensure workers implement the trigger handler interface
func (s *FaasWorkerPool) GetTriggerHandler() (handler.TriggerHandler, error) {
	s.workerLock.Lock()
	defer s.workerLock.Unlock()

	if len(s.workers) > 0 {
		return s.workers[0], nil
	} else {
		return nil, fmt.Errorf("No available workers in this pool!")
	}

	return s
}

// Synchronously wait for at least one active worker
func (s *FaasWorkerPool) WaitForActiveWorker(timeout int) error {
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

func (s *FaasWorkerPool) getWorkerCount() int {
	s.workerLock.Lock()
	defer s.workerLock.Unlock()
	return len(s.workers)
}

// Add a New FaaS worker to this pool
func (s *FaasWorkerPool) AddWorker(stream pb.Faas_TriggerStreamServer) error {
	s.workerLock.Lock()
	defer s.workerLock.Unlock()
	workerCount := len(workers)

	// Ensure we haven't reached the maximum number of workers
	if workerCount > maxWorkers {
		return fmt.Errorf("Max worker capacity reached! Cannot add more workers!")
	}

	// Add a new worker to this pool
	workers[workerCount] = newFaasWorker(stream)
}

func NewFaasWorkerPool() *FaasWorkerPool {
	return &FaasWorkerPool{
		// Only need one at the moment, but leaving open to future proofing
		maxWorkers: 1,
		// Pre-allocate this for efficiency
		workers: make(*FaasWorker, 1)
	}
}
