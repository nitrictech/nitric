package sdk

import (
	"fmt"
	"strings"
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
	Push(queue string, events []NitricEvent) (*PushResponse, error)
	Pop(options PopOptions) ([]NitricQueueItem, error)
}

type PopOptions struct {
	// Nitric name for the queue.
	//
	// The Nitric name will match the AWS SQS Queue name.
	//
	// queueName is a required field
	QueueName string `type:"string" required:"true"`

	// Max depth of queue messages to pop.
	//
	// If nil or 0, defaults to depth 1.
	Depth *uint32 `type:"int" required:"false"`
}

func (p *PopOptions) Validate() error  {
	// Validation
	var invalidParams []string
	if p.QueueName == "" {
		invalidParams = append(invalidParams, fmt.Errorf("queueName param must not be blank").Error())
	}
	if len(invalidParams) > 0 {
		return fmt.Errorf( "invalid params: %s", strings.Join(invalidParams, "\n"))
	}

	// Defaults
	// Set depth to 1 by default
	if p.Depth == nil {
		p.Depth = new(uint32)
		*p.Depth = 1
	}	else if *p.Depth < 1 {
		*p.Depth = uint32(1)
	}
	return nil
}

// UnimplementedQueuePlugin - A Default interface, that provide implementations of QueuePlugin methods that
// Flag the method as unimplemented
type UnimplementedQueuePlugin struct {
	QueuePlugin
}

// Ensure UnimplementedQueuePlugin conforms to QueuePlugin interface
var _ QueuePlugin = (*UnimplementedQueuePlugin)(nil)

// Push - Unimplemented Stub for the UnimplementedQueuePlugin
func (*UnimplementedQueuePlugin) Push(queue string, events []NitricEvent) (*PushResponse, error) {
	return nil, fmt.Errorf("UNIMPLEMENTED")
}

// Pop - Unimplemented Stub for the UnimplementedQueuePlugin
//func (*UnimplementedQueuePlugin) Pop(options PopOptions) ([]NitricQueueItem, error) {
//	return nil, fmt.Errorf("UNIMPLEMENTED")
//}
