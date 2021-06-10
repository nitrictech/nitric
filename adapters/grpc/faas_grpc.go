package grpc

import (
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
	return <-errchan
}

func NewFaasServer(workerPool worker.WorkerPool) *FaasServer {
	return &FaasServer{
		pool: workerPool,
	}
}
