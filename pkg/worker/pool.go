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

	"github.com/nitrictech/nitric/pkg/triggers"
)

type WorkerPool interface {
	// WaitForMinimumWorkers - A blocking method
	WaitForMinimumWorkers(timeout int) error
	GetWorkerCount() int
	GetWorker(*GetWorkerOptions) (Worker, error)
	GetWorkers(*GetWorkerOptions) []Worker
	AddWorker(Worker) error
	RemoveWorker(Worker) error
	Monitor() error
}

type ProcessPoolOptions struct {
	MinWorkers int
	MaxWorkers int
}

// ProcessPool - A worker pool that represent co-located processes
type ProcessPool struct {
	minWorkers int
	maxWorkers int
	workerLock sync.Locker
	workers    []Worker
	poolErr    chan error
}

func (p *ProcessPool) GetWorkerCount() int {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()
	return len(p.workers)
}

func prepend(slice []Worker, elems ...Worker) []Worker {
	return append(elems, slice...)
}

// return route workers
func (p *ProcessPool) getHttpWorkers() []Worker {
	hws := make([]Worker, 0)

	for _, w := range p.workers {
		switch w.(type) {
		case *ScheduleWorker:
			break
		case *SubscriptionWorker:
			break
		case *RouteWorker:
			// Prioritise Route Workers
			hws = prepend(hws, w)
			break
		default:
			hws = append(hws, w)
		}
	}

	return hws
}

// return route workers
func (p *ProcessPool) getEventWorkers() []Worker {
	hws := make([]Worker, 0)

	for _, w := range p.workers {
		switch w.(type) {
		case *RouteWorker:
			// Ignore route workers
			break
		case *ScheduleWorker:
			hws = prepend(hws, w)
		case *SubscriptionWorker:
			hws = prepend(hws, w)
		default:
			hws = append(hws, w)
		}
	}

	return hws
}

// GetMinWorkers - return the minimum number of workers for this pool
func (p *ProcessPool) GetMinWorkers() int {
	return p.minWorkers
}

// GetMaxWorkers - return the maximum number of workers for this pool
func (p *ProcessPool) GetMaxWorkers() int {
	return p.maxWorkers
}

// Monitor - Blocks the current thread to supervise this worker pool
func (p *ProcessPool) Monitor() error {
	// Returns a pool error
	// In future we can catch this and attempt to create new workers to recover
	err := <-p.poolErr

	return err
}

// WaitForMinimumWorkers - Waits for the configured minimum number of workers to be available in this pool
func (p *ProcessPool) WaitForMinimumWorkers(timeout int) error {
	maxWaitTime := time.Duration(timeout) * time.Second
	// Longer poll times, e.g. 200 milliseconds results in slow lambda cold starts (15s+)
	pollInterval := time.Duration(15) * time.Millisecond

	var waitedTime = time.Duration(0)
	for {
		if p.GetWorkerCount() >= p.minWorkers {
			break
		} else {
			if waitedTime < maxWaitTime {
				time.Sleep(pollInterval)
				waitedTime += pollInterval
			} else {
				return fmt.Errorf("available workers below required minimum of %d, %d available, timedout waiting for more workers", p.minWorkers, p.GetWorkerCount())
			}
		}
	}

	return nil
}

type GetWorkerOptions struct {
	Http   *triggers.HttpRequest
	Event  *triggers.Event
	Filter func(w Worker) bool
}

func filterWorkers(ws []Worker, f func(w Worker) bool) []Worker {
	newWs := make([]Worker, 0)
	for _, w := range ws {
		if f(w) {
			newWs = append(newWs, w)
		}
	}

	return newWs
}

// GetWorkers - return a slice of all workers matching the input options.
// useful for retrieving a list of all topic subscribers (for example)
func (p *ProcessPool) GetWorkers(opts *GetWorkerOptions) []Worker {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	workers := make([]Worker, 0)

	if opts.Http != nil {
		ws := p.getHttpWorkers()

		for _, w := range ws {
			if w.HandlesHttpRequest(opts.Http) {
				workers = append(workers, w)
			}
		}
	}

	if opts.Event != nil {
		ws := p.getEventWorkers()

		for _, w := range ws {
			if w.HandlesEvent(opts.Event) {
				workers = append(workers, w)
			}
		}
	}

	if opts.Filter != nil {
		workers = filterWorkers(workers, opts.Filter)
	}

	return workers
}

// GetWorker - Retrieves a worker from this pool
func (p *ProcessPool) GetWorker(opts *GetWorkerOptions) (Worker, error) {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	if opts.Http != nil {
		ws := p.getHttpWorkers()

		if opts.Filter != nil {
			ws = filterWorkers(ws, opts.Filter)
		}

		for _, w := range ws {
			if w.HandlesHttpRequest(opts.Http) {
				return w, nil
			}
		}
	}

	if opts.Event != nil {
		ws := p.getEventWorkers()

		if opts.Filter != nil {
			ws = filterWorkers(ws, opts.Filter)
		}

		for _, w := range ws {
			if w.HandlesEvent(opts.Event) {
				return w, nil
			}
		}
	}

	return nil, fmt.Errorf("no valid workers available")
}

// RemoveWorker - Removes the given worker from this pool
func (p *ProcessPool) RemoveWorker(wrkr Worker) error {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	for i, w := range p.workers {
		if wrkr == w {
			p.workers = append(p.workers[:i], p.workers[i+1:]...)
			if len(p.workers) < p.minWorkers {
				p.poolErr <- fmt.Errorf("insufficient workers in pool, need minimum of %d, %d available", p.minWorkers, len(p.workers))
			}

			return nil
		}
	}

	return fmt.Errorf("worker does not exist in this pool")
}

// AddWorker - Adds the given worker to this pool
func (p *ProcessPool) AddWorker(wrkr Worker) error {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	workerCount := len(p.workers)

	// Ensure we haven't reached the maximum number of workers
	if workerCount > p.maxWorkers {
		return fmt.Errorf("max worker capacity reached! cannot add more workers")
	}

	p.workers = append(p.workers, wrkr)

	fmt.Printf("new Worker in pool: %d workers available\n", len(p.workers))

	return nil
}

// NewProcessPool - Creates a new process pool
func NewProcessPool(opts *ProcessPoolOptions) WorkerPool {
	if opts.MaxWorkers < 1 {
		opts.MaxWorkers = 1
	}

	return &ProcessPool{
		minWorkers: opts.MinWorkers,
		maxWorkers: opts.MaxWorkers,
		workerLock: &sync.Mutex{},
		workers:    make([]Worker, 0),
		poolErr:    make(chan error),
	}
}
