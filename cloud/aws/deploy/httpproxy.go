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
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAwsPulumiProvider) Http(ctx *pulumi.Context, parent pulumi.Resource, name string, http *deploymentspb.Http) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	targetLambda := a.Lambdas[http.Target.GetService()]

	httpProxyGatewayTags := common.Tags(a.StackId, name, resources.HttpProxy)

	doc := targetLambda.InvokeArn.ApplyT(func(invokeArn string) (string, error) {
		spec := newApiSpec(name, invokeArn, httpProxyGatewayTags)

		b, err := spec.MarshalJSON()
		if err != nil {
			return "", err
		}

		return string(b), nil
	}).(pulumi.StringOutput)

	a.HttpProxies[name], err = apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		Body:           doc,
		ProtocolType:   pulumi.String("HTTP"),
		Tags:           pulumi.ToStringMap(httpProxyGatewayTags),
		FailOnWarnings: pulumi.Bool(true),
	}, opts...)
	if err != nil {
		return err
	}

	_, err = apigatewayv2.NewStage(ctx, name+"DefaultStage", &apigatewayv2.StageArgs{
		AutoDeploy: pulumi.BoolPtr(true),
		Name:       pulumi.String("$default"),
		ApiId:      a.HttpProxies[name].ID(),
	}, opts...)
	if err != nil {
		return err
	}

	// Generate lambda permissions enabling the API Gateway to invoke the function it targets
	_, err = awslambda.NewPermission(ctx, name+http.Target.GetService(), &awslambda.PermissionArgs{
		Function:  targetLambda.Name,
		Action:    pulumi.String("lambda:InvokeFunction"),
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: pulumi.Sprintf("%s/*/*/*", a.HttpProxies[name].ExecutionArn),
	}, opts...)
	if err != nil {
		return err
	}

	return nil
}

func newApiSpec(name, invokeArn string, tags map[string]string) *openapi3.T {
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

	for n, v := range tags {
		doc.Tags = append(doc.Tags, &openapi3.Tag{
			Name:       n,
			Extensions: map[string]interface{}{"x-amazon-apigateway-tag-value": v},
		})
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
