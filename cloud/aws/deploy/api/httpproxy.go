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
	"github.com/getkin/kin-openapi/openapi3"
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

	doc := args.LambdaFunction.Function.InvokeArn.ApplyT(func(invokeArn string) (string, error) {
		spec := newApiSpec(name, invokeArn)

		// augment the api specs with security definitions where available
		b, err := spec.MarshalJSON()
		if err != nil {
			return "", err
		}

		return string(b), nil
	}).(pulumi.StringOutput)

	res.Api, err = apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		Body:           doc,
		ProtocolType:   pulumi.String("HTTP"),
		Tags:           common.Tags(ctx, args.StackID, name),
		FailOnWarnings: pulumi.Bool(true),
	}, opts...)
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

	// Generate lambda permissions enabling the API Gateway to invoke the functions it targets
	_, err = awslambda.NewPermission(ctx, name+args.LambdaFunction.Name, &awslambda.PermissionArgs{
		Function:  args.LambdaFunction.Function.Name,
		Action:    pulumi.String("lambda:InvokeFunction"),
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("%s/*/*/*", res.Api.ExecutionArn),
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

func newApiSpec(name, invokeArn string) *openapi3.T {
	doc := &openapi3.T{
		Info: &openapi3.Info{
			Title:   name,
			Version: "v1",
		},
		OpenAPI: "3.0.1",
		Components: &openapi3.Components{
			SecuritySchemes: make(openapi3.SecuritySchemes),
		},
		Paths: openapi3.Paths{
			"/{proxy+}": &openapi3.PathItem{
				Get:     getOperation(invokeArn, "get"),
				Post:    getOperation(invokeArn, "post"),
				Patch:   getOperation(invokeArn, "patch"),
				Put:     getOperation(invokeArn, "put"),
				Delete:  getOperation(invokeArn, "delete"),
				Options: getOperation(invokeArn, "options"),
			},
		},
	}

	return doc
}

func getOperation(invokeArn string, operationId string) *openapi3.Operation {
	defaultDescription := "default description"

	return &openapi3.Operation{
		OperationID: operationId,
		Responses: openapi3.Responses{
			"default": &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: &defaultDescription,
				},
			},
		},
		Extensions: map[string]interface{}{
			"x-amazon-apigateway-integration": map[string]string{
				"type":                 "aws_proxy",
				"httpMethod":           "POST",
				"payloadFormatVersion": "2.0",
				"uri":                  invokeArn,
			},
		},
		Parameters: openapi3.Parameters{
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:              "path",
					Name:            "proxy+",
					AllowEmptyValue: true,
				},
			},
		},
	}
}
