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

package worker

import (
	"fmt"
	"sync"
	"time"
)

type WorkerPool interface {
	// WaitForActiveWorkers - A blocking method
	WaitForActiveWorkers(timeout int) error
	GetWorker() (Worker, error)
	AddWorker(Worker) error
}

type ProcessPoolOptions struct {
	MaxWorkers int
}

// A worker pool that represent co-located processes
type ProcessPool struct {
	maxWorkers int
	workerLock sync.Mutex
	workers    []Worker
}

func (p *ProcessPool) getWorkerCount() int {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()
	return len(p.workers)
}

func (p *ProcessPool) WaitForActiveWorkers(timeout int) error {
	maxWaitTime := time.Duration(timeout) * time.Second
	// Longer poll times, e.g. 200 milliseconds results in slow lambda cold starts (15s+)
	pollInterval := time.Duration(15) * time.Millisecond

	var waitedTime = time.Duration(0)
	for {
		if p.getWorkerCount() >= 1 {
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				return fmt.Errorf("No workers available!")
			}
		}
	}

	return nil
}

// TODO: Add policy logic for managing worker assignment
func (p *ProcessPool) GetWorker() (Worker, error) {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	if len(p.workers) > 0 {
		return p.workers[0], nil
	} else {
		return nil, fmt.Errorf("No available workers in this pool!")
	}
}

func (p *ProcessPool) AddWorker(wrkr Worker) error {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	workerCount := len(p.workers)

	// Ensure we haven't reached the maximum number of workers
	if workerCount > p.maxWorkers {
		return fmt.Errorf("Max worker capacity reached! Cannot add more workers!")
	}

	p.workers = append(p.workers, wrkr)

	return nil
}

func NewProcessPool(opts *ProcessPoolOptions) WorkerPool {
	if opts.MaxWorkers < 1 {
		opts.MaxWorkers = 1
	}

	return &ProcessPool{
		maxWorkers: opts.MaxWorkers,
		workerLock: sync.Mutex{},
		workers:    make([]Worker, 0),
	}
}
