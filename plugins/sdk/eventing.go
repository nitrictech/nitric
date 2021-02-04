package sdk

import "fmt"

type EventService interface {
	GetTopics() ([]string, error)
	Publish(topic string, event *NitricEvent) error
}

type UnimplementedEventingPlugin struct {
	EventService
}

func (*UnimplementedEventingPlugin) GetTopics() ([]string, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedEventingPlugin) Publish(topic string, event *NitricEvent) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
