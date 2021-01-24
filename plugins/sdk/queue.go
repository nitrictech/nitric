package sdk

import (
	"fmt"
)

type FailedMessage struct {
	Event   *NitricEvent
	Message string
}

type PushResponse struct {
	FailedMessages []*FailedMessage
}

// QueuePlugin - The Nitric plugin interface for cloud native queue services
type QueuePlugin interface {
	// Push - The push method for the Nitric Queue Service
	Push(queue string, events []*NitricEvent) (*PushResponse, error)
}

// UnimplementedQueuePlugin - A Default interface, that provide implementations of QueuePlugin methods that
// Flag the method as unimplemented
type UnimplementedQueuePlugin struct {
	QueuePlugin
}

// Push - Unimplemented Stuv for the UnimplementedQueuePlugin
func (*UnimplementedQueuePlugin) Push(queue string, events []*NitricEvent) (*PushResponse, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}
