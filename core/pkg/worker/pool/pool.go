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

package pool

import (
	"fmt"
	"sync"
	"time"

	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/core/pkg/worker"
)

type WorkerPool interface {
	// WaitForMinimumWorkers - A blocking method
	WaitForMinimumWorkers(timeout int) error
	GetWorkerCount() int
	GetWorker(*GetWorkerOptions) (worker.Worker, error)
	GetWorkers(*GetWorkerOptions) []worker.Worker
	AddWorker(worker.Worker) error
	RemoveWorker(worker.Worker) error
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
	workers    []worker.Worker
	poolErr    chan error
}

func (p *ProcessPool) GetWorkerCount() int {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()
	return len(p.workers)
}

func prepend(slice []worker.Worker, elems ...worker.Worker) []worker.Worker {
	return append(elems, slice...)
}

// return route workers
func (p *ProcessPool) getHttpWorkers() []worker.Worker {
	hws := make([]worker.Worker, 0)

	for _, w := range p.workers {
		switch w.(type) {
		case *worker.ScheduleWorker:
			break
		case *worker.SubscriptionWorker:
			break
		case *worker.RouteWorker:
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
func (p *ProcessPool) getEventWorkers() []worker.Worker {
	hws := make([]worker.Worker, 0)

	for _, w := range p.workers {
		switch w.(type) {
		case *worker.RouteWorker:
			// Ignore route workers
			break
		case *worker.ScheduleWorker:
			hws = prepend(hws, w)
		case *worker.SubscriptionWorker:
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
	waitUntil := time.Now().Add(time.Duration(timeout) * time.Second)
	ticker := time.NewTicker(time.Duration(5) * time.Millisecond)

	// stop the ticker on exit
	defer ticker.Stop()

	for {
		if p.GetWorkerCount() >= p.minWorkers {
			break
		}

		// wait for the next tick
		time := <-ticker.C

		if time.After(waitUntil) {
			return fmt.Errorf("available workers below required minimum of %d, %d available, timed out waiting for more workers", p.minWorkers, p.GetWorkerCount())
		}
	}

	return nil
}

type GetWorkerOptions struct {
	Trigger *v1.TriggerRequest
	Filter  func(w worker.Worker) bool
}

func filterWorkers(ws []worker.Worker, f func(w worker.Worker) bool) []worker.Worker {
	newWs := make([]worker.Worker, 0)
	for _, w := range ws {
		if f(w) {
			newWs = append(newWs, w)
		}
	}

	return newWs
}

// GetWorkers - return a slice of all workers matching the input options.
// useful for retrieving a list of all topic subscribers (for example)
func (p *ProcessPool) GetWorkers(opts *GetWorkerOptions) []worker.Worker {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	workers := append([]worker.Worker{}, p.workers...)

	if opts.Trigger != nil {
		workers = filterWorkers(workers, func(w worker.Worker) bool {
			return w.HandlesTrigger(opts.Trigger)
		})
	}

	if opts.Filter != nil {
		workers = filterWorkers(workers, opts.Filter)
	}

	return workers
}

// GetWorker - Retrieves a worker from this pool
func (p *ProcessPool) GetWorker(opts *GetWorkerOptions) (worker.Worker, error) {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	ws := p.workers

	if opts.Trigger.GetHttp() != nil {
		// prioritise http workers
		ws = p.getHttpWorkers()
	}

	if opts.Trigger.GetTopic() != nil {
		// prioritise event workers
		ws = p.getEventWorkers()
	}

	// fallback to all workers (don't prioritise order based on trigger type)
	if opts.Filter != nil {
		ws = filterWorkers(ws, opts.Filter)
	}

	for _, w := range ws {
		if w.HandlesTrigger(opts.Trigger) {
			return w, nil
		}
	}

	return nil, fmt.Errorf("no valid workers available")
}

// RemoveWorker - Removes the given worker from this pool
func (p *ProcessPool) RemoveWorker(wrkr worker.Worker) error {
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
func (p *ProcessPool) AddWorker(wrkr worker.Worker) error {
	p.workerLock.Lock()
	defer p.workerLock.Unlock()

	workerCount := len(p.workers)

	// Ensure we haven't reached the maximum number of workers
	if workerCount > p.maxWorkers {
		return fmt.Errorf("max worker capacity reached! cannot add more workers")
	}

	p.workers = append(p.workers, wrkr)

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
		workers:    make([]worker.Worker, 0),
		poolErr:    make(chan error),
	}
}
