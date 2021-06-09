package worker

import "github.com/nitric-dev/membrane/handler"

type WorkerPool interface {
	// WaitForActiveWorkers - A blocking method
	WaitForActiveWorkers(timeout int) error
	GetTriggerHandler() (handler.TriggerHandler, error)
}
