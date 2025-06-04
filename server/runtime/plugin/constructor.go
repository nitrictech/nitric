package plugin

type Constructor[T any] func() (T, error)

type Register[T any] func(name string, constructor Constructor[T]) error
