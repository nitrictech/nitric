package sdk

import "fmt"

type EventService interface {
	Publish(topic string, event *NitricEvent) error
	ListTopics() ([]string, error)
}

type UnimplementedEventingPlugin struct {
	EventService
}

func (*UnimplementedEventingPlugin) Publish(topic string, event *NitricEvent) error {
	return fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedEventingPlugin) ListTopics() ([]string, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
