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
	"fmt"

	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/worker"
)

type FaasServer struct {
	pb.UnimplementedFaasServer
	// srv  pb.Faas_TriggerStreamServer
	pool worker.WorkerPool
}

// Starts a new stream
// A reference to this stream will be passed on to a new worker instance
// This represents a new server that is ready to begin processing
func (s *FaasServer) TriggerStream(stream pb.Faas_TriggerStreamServer) error {
	// Create a new worker
	worker := worker.NewFaasWorker(stream)

	// Add it to our new pool
	if err := s.pool.AddWorker(worker); err != nil {
		// Worker could not be added
		// Cancel the stream by returning an error
		// This should cause the spawned child process to exit
		return err
	}

	// We're good to go
	errchan := make(chan error)

	// Start the worker
	go worker.Listen(errchan)

	// block here on error returned from the worker
	err := <-errchan
	fmt.Println("FaaS stream closed, removing worker")

	// Worker is done so we can remove it from the pool
	s.pool.RemoveWorker(worker)

	return err
}

func NewFaasServer(workerPool worker.WorkerPool) *FaasServer {
	return &FaasServer{
		pool: workerPool,
	}
}
