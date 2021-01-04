package sdk

import "fmt"

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

type GatewayHandler func(*NitricRequest) *NitricResponse

type GatewayPlugin interface {
	// Start the Gateway
	// This method should block
	Start(handler GatewayHandler) error
}

type UnimplementedGatewayPlugin struct {
	GatewayPlugin
}

func (*UnimplementedGatewayPlugin) Start(_ GatewayHandler) error {
	return fmt.Errorf("UNIMPLEMENTED")
}
