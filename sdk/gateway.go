package sdk

import (
	"fmt"

	"github.com/nitric-dev/membrane/handler"
)

// SourceType enum
type SourceType int

const (
	Subscription SourceType = iota
	Request
	Custom
)

func (e SourceType) String() string {
	return []string{"SUBSCRIPTION", "REQUEST", "CUSTOM"}[e]
}

type NitricContext struct {
	RequestId   string
	PayloadType string
	Source      string
	SourceType  SourceType
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
	Start(handler handler.SourceHandler) error
}

type UnimplementedGatewayPlugin struct {
	GatewayService
}

func (*UnimplementedGatewayPlugin) Start(_ handler.SourceHandler) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
