// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAwsPulumiProvider) Websocket(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Websocket) error {
	defaultTarget := a.lambdas[config.MessageTarget.GetService()]
	connectTarget := a.lambdas[config.ConnectTarget.GetService()]
	disconnectTarget := a.lambdas[config.DisconnectTarget.GetService()]

	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	websocketApi, err := apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		ProtocolType: pulumi.String("WEBSOCKET"),
		Tags:         pulumi.ToStringMap(tags.Tags(a.stackId, name, resources.Websocket)),
		// TODO: We won't actually be using this, but it is required.
		// Instead we'll be using the $default route
		RouteSelectionExpression: pulumi.String("$request.body.action"),
	}, opts...)
	if err != nil {
		return err
	}

	a.websockets[name] = websocketApi

	// Create the API integrations
	integrationDefault, err := apigatewayv2.NewIntegration(ctx, fmt.Sprintf("%s-default-integration", name), &apigatewayv2.IntegrationArgs{
		ApiId:           websocketApi.ID(),
		IntegrationType: pulumi.String("AWS_PROXY"),
		IntegrationUri:  defaultTarget.Arn,
	}, opts...)
	if err != nil {
		return err
	}

	_, err = awslambda.NewPermission(ctx, fmt.Sprintf("%s-default-permission", name), &awslambda.PermissionArgs{
		Function:  defaultTarget.Name,
		Action:    pulumi.String("lambda:InvokeFunction"),
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("%s/*/*", websocketApi.ExecutionArn),
	}, opts...)
	if err != nil {
		return err
	}

	// check if the function name is different if not assign to default
	integrationConnect := integrationDefault
	if connectTarget != defaultTarget {
		integrationConnect, err = apigatewayv2.NewIntegration(ctx, fmt.Sprintf("%s-connect-integration", name), &apigatewayv2.IntegrationArgs{
			ApiId:           websocketApi.ID(),
			IntegrationType: pulumi.String("AWS_PROXY"),
			IntegrationUri:  connectTarget.Arn,
		}, opts...)
		if err != nil {
			return err
		}

		_, err = awslambda.NewPermission(ctx, fmt.Sprintf("%s-connect-permission", name), &awslambda.PermissionArgs{
			Function:  defaultTarget.Name,
			Action:    pulumi.String("lambda:InvokeFunction"),
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("%s/*/*", websocketApi.ExecutionArn),
		}, opts...)
		if err != nil {
			return err
		}
	}

	// check if the function name is different if not assign to default
	integrationDisconnect := integrationDefault
	if disconnectTarget != defaultTarget {
		integrationDisconnect, err = apigatewayv2.NewIntegration(ctx, fmt.Sprintf("%s-disconnect-integration", name), &apigatewayv2.IntegrationArgs{
			ApiId:           websocketApi.ID(),
			IntegrationType: pulumi.String("AWS_PROXY"),
			IntegrationUri:  disconnectTarget.Arn,
		}, opts...)
		if err != nil {
			return err
		}

		_, err = awslambda.NewPermission(ctx, fmt.Sprintf("%s-disconnect-permission", name), &awslambda.PermissionArgs{
			Function:  defaultTarget.Name,
			Action:    pulumi.String("lambda:InvokeFunction"),
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("%s/*/*", websocketApi.ExecutionArn),
		}, opts...)
		if err != nil {
			return err
		}
	}

	// Create the routes for the websocket handler
	// The default message route
	_, err = apigatewayv2.NewRoute(ctx, fmt.Sprintf("%s-default-route", name), &apigatewayv2.RouteArgs{
		ApiId:    websocketApi.ID(),
		RouteKey: pulumi.String("$default"),
		Target:   pulumi.Sprintf("integrations/%s", integrationDefault.ID()),
	}, opts...)
	if err != nil {
		return err
	}

	// The client connection route
	_, err = apigatewayv2.NewRoute(ctx, fmt.Sprintf("%s-connect-route", name), &apigatewayv2.RouteArgs{
		ApiId:    websocketApi.ID(),
		RouteKey: pulumi.String("$connect"),
		Target:   pulumi.Sprintf("integrations/%s", integrationConnect.ID()),
	}, opts...)
	if err != nil {
		return err
	}

	// the client disconnection route
	_, err = apigatewayv2.NewRoute(ctx, fmt.Sprintf("%s-disconnect-route", name), &apigatewayv2.RouteArgs{
		ApiId:    websocketApi.ID(),
		RouteKey: pulumi.String("$disconnect"),
		Target:   pulumi.Sprintf("integrations/%s", integrationDisconnect.ID()),
	}, opts...)
	if err != nil {
		return err
	}

	_, err = apigatewayv2.NewStage(ctx, name+"DefaultStage", &apigatewayv2.StageArgs{
		AutoDeploy: pulumi.BoolPtr(true),
		Name:       pulumi.String(common.DefaultWsStageName),
		ApiId:      websocketApi.ID(),
		Tags:       pulumi.ToStringMap(tags.Tags(a.stackId, name+"DefaultStage", resources.Websocket)),
	}, opts...)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return nil
}
