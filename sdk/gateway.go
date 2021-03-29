package sdk

import (
	"fmt"

	"github.com/nitric-dev/membrane/handler"
	"github.com/nitric-dev/membrane/triggers"
)

type NitricContext struct {
	RequestId   string
	PayloadType string
	Trigger     string
	TriggerType triggers.TriggerType
}

// Normalized NitricRequest
type NitricRequest struct {
	Context     *NitricContext
	ContentType string
	Payload     []byte
}

type NitricResponse struct {
	Headers map[string]string
	Status  int
	Body    []byte
}

type GatewayService interface {
	// Start the Gateway
	// This method should block
	Start(handler handler.TriggerHandler) error
}

type UnimplementedGatewayPlugin struct {
	GatewayService
}

func (*UnimplementedGatewayPlugin) Start(_ handler.TriggerHandler) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
