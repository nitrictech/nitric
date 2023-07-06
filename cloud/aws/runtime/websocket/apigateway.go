package websocket

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/nitrictech/nitric/cloud/aws/runtime/core"
	"github.com/nitrictech/nitric/core/pkg/plugins/resource"
	"github.com/nitrictech/nitric/core/pkg/plugins/websocket"
	"github.com/nitrictech/nitric/core/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type ApiGatewayWebsocketService struct {
	websocket.UnimplementedWebsocketService
	provider core.AwsProvider
	clients  map[string]*apigatewaymanagementapi.Client
}

var _ websocket.WebsocketService = &ApiGatewayWebsocketService{}

func (a *ApiGatewayWebsocketService) getClientForSocket(socket string) (*apigatewaymanagementapi.Client, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")

	if client, ok := a.clients[socket]; ok {
		return client, nil
	}

	details, err := a.provider.Details(context.TODO(), resource.ResourceType_Api, socket)
	if err != nil {
		return nil, err
	}

	apiDetails, ok := details.Detail.(resource.ApiDetails)
	if !ok {
		return nil, fmt.Errorf("an error occurred resolving API Gateway details")
	}

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	callbackUrl := strings.Replace(apiDetails.URL, "wss", "https", 1)
	callbackUrl = callbackUrl + "/$default"

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	a.clients[socket] = apigatewaymanagementapi.NewFromConfig(cfg, apigatewaymanagementapi.WithEndpointResolver(apigatewaymanagementapi.EndpointResolverFromURL(
		callbackUrl,
	)))

	return a.clients[socket], nil
}

func (a *ApiGatewayWebsocketService) Send(ctx context.Context, socket string, connectionId string, message []byte) error {
	client, err := a.getClientForSocket(socket)
	if err != nil {
		return err
	}

	_, err = client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionId),
		Data:         message,
	})

	if err != nil {
		return err
	}

	return nil
}

func (a *ApiGatewayWebsocketService) Close(ctx context.Context, socket string, connectionId string) error {
	client, err := a.getClientForSocket(socket)
	if err != nil {
		return err
	}

	_, err = client.DeleteConnection(ctx, &apigatewaymanagementapi.DeleteConnectionInput{
		ConnectionId: aws.String(connectionId),
	})
	if err != nil {
		return err
	}

	return nil
}

func NewAwsApiGatewayWebsocket(provider core.AwsProvider) (*ApiGatewayWebsocketService, error) {
	return &ApiGatewayWebsocketService{
		provider: provider,
		clients:  make(map[string]*apigatewaymanagementapi.Client),
	}, nil
}
