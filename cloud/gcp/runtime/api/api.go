package api

import (
	"context"

	"github.com/nitrictech/nitric/cloud/gcp/runtime/resource"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
)

type GcpApiGatewayProvider struct {
	provider *resource.GcpResourceService
	*apis.RouteWorkerManager
}

var _ apispb.ApiServer = &GcpApiGatewayProvider{}

func (g *GcpApiGatewayProvider) Details(ctx context.Context, req *apispb.ApiDetailsRequest) (*apispb.ApiDetailsResponse, error) {
	gwDetails, err := g.provider.GetApiGatewayDetails(ctx, req.ApiName)
	if err != nil {
		return nil, err
	}

	return &apispb.ApiDetailsResponse{
		Url: gwDetails.Url,
	}, nil
}

func NewGcpApiGatewayProvider(provider *resource.GcpResourceService) *GcpApiGatewayProvider {
	return &GcpApiGatewayProvider{
		provider:           provider,
		RouteWorkerManager: apis.New(),
	}
}
