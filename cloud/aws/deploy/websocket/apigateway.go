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

package websocket

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

type AwsWebsocketApiGatewayArgs struct {
	DefaultTarget    *exec.LambdaExecUnit
	ConnectTarget    *exec.LambdaExecUnit
	DisconnectTarget *exec.LambdaExecUnit

	StackID pulumi.StringInput
}

type AwsWebsocketApiGateway struct {
	pulumi.ResourceState

	Name string
	Api  *apigatewayv2.Api
}

func NewAwsWebsocketApiGateway(ctx *pulumi.Context, name string, args *AwsWebsocketApiGatewayArgs, opts ...pulumi.ResourceOption) (*AwsWebsocketApiGateway, error) {
	res := &AwsWebsocketApiGateway{Name: name}

	err := ctx.RegisterComponentResource("nitric:websocket:AwsApiGateway", name, res, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(res))

	res.Api, err = apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		ProtocolType: pulumi.String("WEBSOCKET"),
		Tags:         common.Tags(ctx, args.StackID, name),
		// TODO: We won't actually be using this, but it is required.
		// Instead we'll be using the $default route
		RouteSelectionExpression: pulumi.String("$request.body.action"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	// Create the API integrations
	integrationDefault, err := apigatewayv2.NewIntegration(ctx, fmt.Sprintf("%s-default-integration", name), &apigatewayv2.IntegrationArgs{
		ApiId:           res.Api.ID(),
		IntegrationType: pulumi.String("AWS_PROXY"),
		IntegrationUri:  args.DefaultTarget.Function.Arn,
	})
	if err != nil {
		return nil, err
	}

	_, err = awslambda.NewPermission(ctx, fmt.Sprintf("%s-default-permission", name), &awslambda.PermissionArgs{
		Function:  args.DefaultTarget.Function.Name,
		Action:    pulumi.String("lambda:InvokeFunction"),
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("%s/*/*", res.Api.ExecutionArn),
	}, opts...)
	if err != nil {
		return nil, err
	}

	// check if the function name is different if not assign to default
	integrationConnect := integrationDefault
	if args.ConnectTarget != args.DefaultTarget {
		integrationConnect, err = apigatewayv2.NewIntegration(ctx, fmt.Sprintf("%s-connect-integration", name), &apigatewayv2.IntegrationArgs{
			ApiId:           res.Api.ID(),
			IntegrationType: pulumi.String("AWS_PROXY"),
			IntegrationUri:  args.ConnectTarget.Function.Arn,
		})
		if err != nil {
			return nil, err
		}

		_, err = awslambda.NewPermission(ctx, fmt.Sprintf("%s-connect-permission", name), &awslambda.PermissionArgs{
			Function:  args.DefaultTarget.Function.Name,
			Action:    pulumi.String("lambda:InvokeFunction"),
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("%s/*/*", res.Api.ExecutionArn),
		}, opts...)
		if err != nil {
			return nil, err
		}
	}

	// check if the function name is different if not assign to default
	integrationDisconnect := integrationDefault
	if args.DisconnectTarget != args.DefaultTarget {
		integrationDisconnect, err = apigatewayv2.NewIntegration(ctx, fmt.Sprintf("%s-disconnect-integration", name), &apigatewayv2.IntegrationArgs{
			ApiId:           res.Api.ID(),
			IntegrationType: pulumi.String("AWS_PROXY"),
			IntegrationUri:  args.DisconnectTarget.Function.Arn,
		})
		if err != nil {
			return nil, err
		}

		_, err = awslambda.NewPermission(ctx, fmt.Sprintf("%s-disconnect-permission", name), &awslambda.PermissionArgs{
			Function:  args.DefaultTarget.Function.Name,
			Action:    pulumi.String("lambda:InvokeFunction"),
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("%s/*/*", res.Api.ExecutionArn),
		}, opts...)
		if err != nil {
			return nil, err
		}
	}

	// Create the routes for the websocket handler
	// The default message route
	_, err = apigatewayv2.NewRoute(ctx, fmt.Sprintf("%s-default-route", name), &apigatewayv2.RouteArgs{
		ApiId:    res.Api.ID(),
		RouteKey: pulumi.String("$default"),
		Target:   pulumi.Sprintf("integrations/%s", integrationDefault.ID()),
	})
	if err != nil {
		return nil, err
	}

	// The client connection route
	_, err = apigatewayv2.NewRoute(ctx, fmt.Sprintf("%s-connect-route", name), &apigatewayv2.RouteArgs{
		ApiId:    res.Api.ID(),
		RouteKey: pulumi.String("$connect"),
		Target:   pulumi.Sprintf("integrations/%s", integrationConnect.ID()),
	})
	if err != nil {
		return nil, err
	}

	// the client disconnection route
	_, err = apigatewayv2.NewRoute(ctx, fmt.Sprintf("%s-disconnect-route", name), &apigatewayv2.RouteArgs{
		ApiId:    res.Api.ID(),
		RouteKey: pulumi.String("$disconnect"),
		Target:   pulumi.Sprintf("integrations/%s", integrationDisconnect.ID()),
	})
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewStage(ctx, name+"DefaultStage", &apigatewayv2.StageArgs{
		AutoDeploy: pulumi.BoolPtr(true),
		Name:       pulumi.String("$default"),
		ApiId:      res.Api.ID(),
		Tags:       common.Tags(ctx, args.StackID, name+"DefaultStage"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	endPoint := res.Api.ApiEndpoint.ApplyT(func(ep string) string {
		return ep
	}).(pulumi.StringInput)

	ctx.Export("api:"+name, endPoint)

	return res, nil
}
