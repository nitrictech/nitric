package resource

import (
	"context"

	resourcepb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

type RuntimeResourceService struct {
}

// At declaring resources at runtime is a no-op, as the resource were already declared during deployment.
func (p *RuntimeResourceService) Declare(ctx context.Context, req *resourcepb.ResourceDeclareRequest) (*resourcepb.ResourceDeclareResponse, error) {
	return &resourcepb.ResourceDeclareResponse{}, nil
}
