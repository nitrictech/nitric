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

package api

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

type AwsHttpProxyArgs struct {
	LambdaFunction *exec.LambdaExecUnit
	StackID        pulumi.StringInput
}

type AwsHttpProxy struct {
	pulumi.ResourceState

	Name string
	Api  *apigatewayv2.Api
}

func NewAwsHttpProxy(ctx *pulumi.Context, name string, args *AwsHttpProxyArgs, opts ...pulumi.ResourceOption) (*AwsHttpProxy, error) {
	res := &AwsHttpProxy{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:AwsApiGateway", name, res, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(res))

	res.Api, err = apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		ProtocolType:   pulumi.String("HTTP"),
		Tags:           common.Tags(ctx, args.StackID, name),
		FailOnWarnings: pulumi.Bool(true),
	}, opts...)
	if err != nil {
		return nil, err
	}

	integrationDefault, err := apigatewayv2.NewIntegration(ctx, fmt.Sprintf("%s-default-integration", name), &apigatewayv2.IntegrationArgs{
		ApiId:           res.Api.ID(),
		IntegrationType: pulumi.String("AWS_PROXY"),
		IntegrationUri:  args.LambdaFunction.Function.Arn,
	})
	if err != nil {
		return nil, err
	}

	_, err = awslambda.NewPermission(ctx, fmt.Sprintf("%s-default-permission", name), &awslambda.PermissionArgs{
		Function:  args.LambdaFunction.Function.Name,
		Action:    pulumi.String("lambda:InvokeFunction"),
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("%s/*/*", res.Api.ExecutionArn),
	}, opts...)
	if err != nil {
		return nil, err
	}

	_, err = apigatewayv2.NewRoute(ctx, name+"DefaultRoute", &apigatewayv2.RouteArgs{
		ApiId:    res.Api.ID(),
		RouteKey: pulumi.String("$default"),
		Target:   pulumi.Sprintf("integrations/%s", integrationDefault.ID()),
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

	endPoint := res.Api.ApiEndpoint.ApplyT(func(ep string) string {
		return ep
	}).(pulumi.StringInput)

	ctx.Export("api:"+name, endPoint)

	return res, nil
}
