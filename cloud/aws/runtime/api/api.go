package api

import (
	"context"

	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
)

type AwsApiGatewayProvider struct {
	provider *resource.AwsResourceService

	*apis.RouteWorkerManager
}

var _ apispb.ApiServer = &AwsApiGatewayProvider{}

func (a *AwsApiGatewayProvider) Details(ctx context.Context, req *apispb.ApiDetailsRequest) (*apispb.ApiDetailsResponse, error) {
	gwDetails, err := a.provider.GetAWSApiGatewayDetails(ctx, &resourcespb.ResourceIdentifier{
		Type: resourcespb.ResourceType_Api,
		Name: req.ApiName,
	})
	if err != nil {
		return nil, err
	}

	return &apispb.ApiDetailsResponse{
		Url: gwDetails.Url,
	}, nil
}

func NewAwsApiGatewayProvider(provider *resource.AwsResourceService) *AwsApiGatewayProvider {
	return &AwsApiGatewayProvider{
		provider:           provider,
		RouteWorkerManager: apis.New(),
	}
}
