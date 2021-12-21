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

	"github.com/nitrictech/nitric/pkg/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/nitrictech/nitric/interfaces/nitric/v1"
)

type FaasServer struct {
	pb.UnimplementedFaasServiceServer
	pool worker.WorkerPool
}

// Starts a new stream
// A reference to this stream will be passed on to a new worker instance
// This represents a new server that is ready to begin processing
func (s *FaasServer) TriggerStream(stream pb.FaasService_TriggerStreamServer) error {
	cm, err := stream.Recv()
	ir := cm.GetInitRequest()

	if ir == nil {
		// SHUT IT DOWN!!!!
		return status.Error(codes.FailedPrecondition, "first message must be InitRequest")
	}

	var wrkr worker.GrpcWorker = nil

	if route := ir.GetRoute(); route != nil {
		// Create a new route worker
		wrkr = worker.NewRouteWorker(stream, &worker.RouteWorkerOptions{
			Path:    route.Path,
			Methods: route.Methods,
		})
	} else if schedule := ir.GetSchedule(); schedule != nil {
		wrkr = worker.NewScheduleWorker(stream, &worker.ScheduleWorkerOptions{
			Description: schedule.GetKey(),
		})
	} else {
		// XXX: Catch all worker type
		wrkr = worker.NewFaasWorker(stream)
	}

	// Add it to our new pool
	if err := s.pool.AddWorker(wrkr); err != nil {
		// Worker could not be added
		// Cancel the stream by returning an error
		// This should cause the spawned child process to exit
		return err
	}

	// We're good to go
	errchan := make(chan error)

	// Start the worker
	go wrkr.Listen(errchan)

	// block here on error returned from the worker
	err = <-errchan
	fmt.Println("FaaS stream closed, removing worker")

	// Worker is done so we can remove it from the pool
	s.pool.RemoveWorker(wrkr)

	return err
}

func NewFaasServer(workerPool worker.WorkerPool) *FaasServer {
	return &FaasServer{
		pool: workerPool,
	}
}
