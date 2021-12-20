package worker

type GrpcWorker interface {
	Worker
	Listen(chan error)
}
