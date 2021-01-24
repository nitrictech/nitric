package sdk

import "fmt"

type EventingPlugin interface {
	GetTopics() ([]string, error)
	Publish(topic string, event *NitricEvent) error
}

type UnimplementedEventingPlugin struct {
	EventingPlugin
}

func (*UnimplementedEventingPlugin) GetTopics() ([]string, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

func (*UnimplementedEventingPlugin) Publish(topic string, event *NitricEvent) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
