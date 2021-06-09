package grpc

import (
	pb "github.com/nitric-dev/membrane/interfaces/nitric/v1"
	"github.com/nitric-dev/membrane/worker"
)

type FaasServer struct {
	pb.UnimplementedFaasServer
	// srv  pb.Faas_TriggerStreamServer
	pool *worker.FaasWorkerPool
}

// Starts a new stream
// A reference to this stream will be passed on to a new worker instance
// This represents a new server that is ready to begin processing
func (s *FaasServer) TriggerStream(srv pb.Faas_TriggerStreamServer) error {
	if err := s.pool.AddWorker(srv); err != nil {
		// return an error here...
		// TODO: Return proper grpc error with status here...
		return err
	}

	return nil
}

func NewFaasServer(workerPool *worker.FaasWorkerPool) *FaasServer {
	return &FaasServer{
		pool: workerPool,
	}
}
