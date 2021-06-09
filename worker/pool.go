package worker

type WorkerPool interface {
	// WaitForActiveWorkers - A blocking method
	WaitForActiveWorkers(timeout int) error
	GetWorker() (Worker, error)
}
