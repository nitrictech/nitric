package api

import (
	"context"

	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
)

type AzureApiGatewayProvider struct {
	provider *resource.AzureResourceService
	*apis.RouteWorkerManager
}

var _ apispb.ApiServer = &AzureApiGatewayProvider{}

func (g *AzureApiGatewayProvider) Details(ctx context.Context, req *apispb.ApiDetailsRequest) (*apispb.ApiDetailsResponse, error) {
	gwDetails, err := g.provider.GetApiDetails(ctx, req.ApiName)
	if err != nil {
		return nil, err
	}

	return &apispb.ApiDetailsResponse{
		Url: gwDetails.Url,
	}, nil
}

func NewAzureApiGatewayProvider(provider *resource.AzureResourceService) *AzureApiGatewayProvider {
	return &AzureApiGatewayProvider{
		provider:           provider,
		RouteWorkerManager: apis.New(),
	}
}
