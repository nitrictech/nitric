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
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	"github.com/nitrictech/nitric/cloud/aws/runtime/resource"
	grpc_errors "github.com/nitrictech/nitric/core/pkg/grpc/errors"
	"github.com/nitrictech/nitric/core/pkg/logger"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	websocketpb "github.com/nitrictech/nitric/core/pkg/proto/websockets/v1"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"
	"google.golang.org/grpc/codes"
)

type ApiGatewayWebsocketService struct {
	provider *resource.AwsResourceService
	clients  map[string]*apigatewaymanagementapi.Client
}

var _ websocketpb.WebsocketServer = &ApiGatewayWebsocketService{}

func (a *ApiGatewayWebsocketService) getClientForSocket(socket string) (*apigatewaymanagementapi.Client, error) {
	awsRegion := env.AWS_REGION.String()

	if client, ok := a.clients[socket]; ok {
		logger.Debug("using existing websocket client found in cache")
		return client, nil
	}

	details, err := a.SocketDetails(context.TODO(), &websocketpb.WebsocketDetailsRequest{
		SocketName: socket,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting websocket details: %w", err)
	}

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session: %w", sessionError)
	}

	// post requests are made to the https endpoint, so the scheme needs to be changed
	callbackUrl := details.Url
	if strings.HasPrefix(details.Url, "wss") {
		callbackUrl = strings.Replace(details.Url, "wss", "https", 1)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	a.clients[socket] = apigatewaymanagementapi.NewFromConfig(cfg, apigatewaymanagementapi.WithEndpointResolver(apigatewaymanagementapi.EndpointResolverFromURL(
		callbackUrl,
	)))

	return a.clients[socket], nil
}

func (a *ApiGatewayWebsocketService) SocketDetails(ctx context.Context, req *websocketpb.WebsocketDetailsRequest) (*websocketpb.WebsocketDetailsResponse, error) {
	gwDetails, err := a.provider.GetAWSApiGatewayDetails(ctx, &resourcespb.ResourceIdentifier{
		Type: resourcespb.ResourceType_Websocket,
		Name: req.SocketName,
	})
	if err != nil {
		return nil, err
	}

	return &websocketpb.WebsocketDetailsResponse{
		Url: fmt.Sprintf("%s/%s", gwDetails.Url, common.DefaultWsStageName),
	}, nil
}

func (a *ApiGatewayWebsocketService) SendMessage(ctx context.Context, req *websocketpb.WebsocketSendRequest) (*websocketpb.WebsocketSendResponse, error) {
	newErr := grpc_errors.ErrorsWithScope("ApiGateway.Websocket.Send")

	client, err := a.getClientForSocket(req.SocketName)
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error getting websocket client",
			err,
		)
	}

	_, err = client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(req.ConnectionId),
		Data:         req.Data,
	})
	if err != nil {
		return nil, newErr(
			codes.Internal,
			"error sending message to websocket",
			err,
		)
	}

	return &websocketpb.WebsocketSendResponse{}, nil
}

func (a *ApiGatewayWebsocketService) CloseConnection(ctx context.Context, req *websocketpb.WebsocketCloseConnectionRequest) (*websocketpb.WebsocketCloseConnectionResponse, error) {
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

	return &websocketpb.WebsocketCloseConnectionResponse{}, nil
}

func NewAwsApiGatewayWebsocket(provider *resource.AwsResourceService) (*ApiGatewayWebsocketService, error) {
	return &ApiGatewayWebsocketService{
		provider: provider,
		clients:  make(map[string]*apigatewaymanagementapi.Client),
	}, nil
}
