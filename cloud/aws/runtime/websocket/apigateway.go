// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package websocket

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	websocketpb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
)

type ApiGatewayWebsocketService struct {
	provider *resource.AwsResourceService
	clients  map[string]*apigatewaymanagementapi.Client
}

var _ websocketpb.WebsocketServer = &ApiGatewayWebsocketService{}

func (a *ApiGatewayWebsocketService) getClientForSocket(socket string) (*apigatewaymanagementapi.Client, error) {
	awsRegion := env.AWS_REGION.String()

	if client, ok := a.clients[socket]; ok {
		return client, nil
	}

	details, err := a.Details(context.TODO(), &websocketpb.WebsocketDetailsRequest{
		SocketName: socket,
	})
	if err != nil {
		return nil, err
	}

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	callbackUrl := strings.Replace(details.Url, "wss", "https", 1)
	callbackUrl = callbackUrl + "/$default"

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	a.clients[socket] = apigatewaymanagementapi.NewFromConfig(cfg, apigatewaymanagementapi.WithEndpointResolver(apigatewaymanagementapi.EndpointResolverFromURL(
		callbackUrl,
	)))

	return a.clients[socket], nil
}

func (a *ApiGatewayWebsocketService) Details(ctx context.Context, req *websocketpb.WebsocketDetailsRequest) (*websocketpb.WebsocketDetailsResponse, error) {
	gwDetails, err := a.provider.GetAWSApiGatewayDetails(ctx, &resourcespb.ResourceIdentifier{
		Type: resourcespb.ResourceType_Websocket,
		Name: req.SocketName,
	})
	if err != nil {
		return nil, err
	}

	return &websocketpb.WebsocketDetailsResponse{
		Url: gwDetails.Url,
	}, nil
}

func (a *ApiGatewayWebsocketService) Send(ctx context.Context, req *websocketpb.WebsocketSendRequest) (*websocketpb.WebsocketSendResponse, error) {
	client, err := a.getClientForSocket(req.SocketName)
	if err != nil {
		return nil, err
	}

	_, err = client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(req.ConnectionId),
		Data:         req.Data,
	})

	if err != nil {
		return nil, err
	}

	return &websocketpb.WebsocketSendResponse{}, nil
}

func (a *ApiGatewayWebsocketService) Close(ctx context.Context, req *websocketpb.WebsocketCloseRequest) (*websocketpb.WebsocketCloseResponse, error) {
	client, err := a.getClientForSocket(req.SocketName)
	if err != nil {
		return nil, err
	}

	_, err = client.DeleteConnection(ctx, &apigatewaymanagementapi.DeleteConnectionInput{
		ConnectionId: aws.String(req.ConnectionId),
	})
	if err != nil {
		return nil, err
	}

	return &websocketpb.WebsocketCloseResponse{}, nil
}

func NewAwsApiGatewayWebsocket(provider *resource.AwsResourceService) (*ApiGatewayWebsocketService, error) {
	return &ApiGatewayWebsocketService{
		provider: provider,
		clients:  make(map[string]*apigatewaymanagementapi.Client),
	}, nil
}
