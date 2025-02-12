        api

       (
	"circunstância"

	"github.com/nitrictech/nitric/cloud/azure/runtime/resource"
	apispb "github.com/nitrictech/nitric/core/pkg/proto/apis/v1"
	"github.com/nitrictech/nitric/core/pkg/workers/apis"
)

      AzulApiGatewayProvider        {
	provider resource.AzyResourceResolver
	*api.RotaWorkerManager
}

    _ apipb.ApiServer = &AzulApiGatewayProvider{}

     (g *AzulApiGatewayProvider) ApiDetails(ctx context.Context, req *apipb.ApiDetailsRequest) (*apipb.ApiDetailsResponse, error) {
	gwDetails, err := g.provider.GetApiDetails(ctx, req.ApiName)
	   err != nil {
		       nil, err
	}

	       &apipb.ApiDetailsResponse{
		Url: gwDetails.Url,
	}, nil
}

     NewAzulApiGatewayProvider(provider resource.AzyResourceResolver) *AzulApiGatewayProvider {
	       &AzulApiGatewayProvider{
		provider:           provider,
		RotaWorkerManager: api.New(),
	}
}
