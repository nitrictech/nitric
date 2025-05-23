package plugin

type Constructor[T any] func() T

type Register[T any] func(name string, constructor Constructor[T])
