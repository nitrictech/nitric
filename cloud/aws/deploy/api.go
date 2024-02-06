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
	"encoding/json"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func awsOperation(op *openapi3.Operation, funcs map[string]string) *openapi3.Operation {
	if op == nil {
		return nil
	}

	name := ""

	if v, ok := op.Extensions["x-nitric-target"]; ok {
		targetMap, isMap := v.(map[string]any)
		if isMap {
			name = targetMap["name"].(string)
		}
	}

	if name == "" {
		return nil
	}

	if _, ok := funcs[name]; !ok {
		return nil
	}

	arn := funcs[name]

	op.Extensions["x-amazon-apigateway-integration"] = map[string]string{
		"type":                 "aws_proxy",
		"httpMethod":           "POST",
		"payloadFormatVersion": "2.0",
		"uri":                  arn,
	}

	return op
}

type nameArnPair struct {
	name      string
	invokeArn string
}

func (a *NitricAwsPulumiProvider) Api(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Api) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	nameArnPairs := make([]interface{}, 0, len(a.lambdas))

	if config.GetOpenapi() == "" {
		return fmt.Errorf("aws provider can only deploy OpenAPI specs")
	}

	openapiDoc := &openapi3.T{}
	err := openapiDoc.UnmarshalJSON([]byte(config.GetOpenapi()))
	if err != nil {
		return fmt.Errorf("invalid document supplied for api: %s", name)
	}

	// augment open api spec with AWS specific security extensions
	if openapiDoc.Components.SecuritySchemes != nil {
		// Start translating to AWS centric security schemes
		for _, scheme := range openapiDoc.Components.SecuritySchemes {
			// implement OpenIDConnect security
			if scheme.Value.Type == "openIdConnect" {
				// We need to extract audience values as well
				// lets use an extension to store these with the document
				audiences := scheme.Value.Extensions["x-nitric-audiences"]

				// Augment extensions with aws specific extensions
				scheme.Value.Extensions["x-amazon-apigateway-authorizer"] = map[string]interface{}{
					"type": "jwt",
					"jwtConfiguration": map[string]interface{}{
						"audience": audiences,
					},
					"identitySource": "$request.header.Authorization",
				}
			} else {
				return fmt.Errorf("unsupported security definition supplied")
			}
		}
	}

	nitricServiceTargets := map[string]*lambda.Function{}
	for _, p := range openapiDoc.Paths {
		for _, op := range p.Operations() {
			if v, ok := op.Extensions["x-nitric-target"]; ok {
				if targetMap, isMap := v.(map[string]any); isMap {
					serviceName := targetMap["name"].(string)
					nitricServiceTargets[serviceName] = a.lambdas[serviceName]
				}
			}
		}
	}

	// collect name arn pairs for output iteration
	for k, v := range nitricServiceTargets {
		nameArnPairs = append(nameArnPairs, pulumi.All(k, v.InvokeArn).ApplyT(func(args []interface{}) nameArnPair {
			name := args[0].(string)
			arn := args[1].(string)

			return nameArnPair{
				name:      name,
				invokeArn: arn,
			}
		}))
	}

	apiGatewayTags := tags.Tags(a.stackId, name, resources.API)

	doc := pulumi.All(nameArnPairs...).ApplyT(func(pairs []interface{}) (string, error) {
		naps := make(map[string]string)

		for _, p := range pairs {
			if pair, ok := p.(nameArnPair); ok {
				naps[pair.name] = pair.invokeArn
			} else {
				// XXX: Should not occur
				return "", fmt.Errorf("invalid data %T %v", p, p)
			}
		}

		for k, p := range openapiDoc.Paths {
			p.Get = awsOperation(p.Get, naps)
			p.Post = awsOperation(p.Post, naps)
			p.Patch = awsOperation(p.Patch, naps)
			p.Put = awsOperation(p.Put, naps)
			p.Delete = awsOperation(p.Delete, naps)
			p.Options = awsOperation(p.Options, naps)
			openapiDoc.Paths[k] = p
		}

		// Important: AWS will use these on the first deployment, but not subsequent updates
		// subsequent updates use the tags provided to Pulumi below.
		for n, v := range apiGatewayTags {
			openapiDoc.Tags = append(openapiDoc.Tags, &openapi3.Tag{
				Name:       n,
				Extensions: map[string]interface{}{"x-amazon-apigateway-tag-value": v},
			})
		}

		// augment the api specs with security definitions where available
		b, err := json.Marshal(openapiDoc)
		if err != nil {
			return "", err
		}

		return string(b), nil
	}).(pulumi.StringOutput)

	a.apis[name], err = apigatewayv2.NewApi(ctx, name, &apigatewayv2.ApiArgs{
		Body: doc,
		// Name fixed to title in the spec, if these mismatch the name will change on the second deployment.
		Name:           pulumi.String(openapiDoc.Info.Title),
		ProtocolType:   pulumi.String("HTTP"),
		Tags:           pulumi.ToStringMap(apiGatewayTags),
		FailOnWarnings: pulumi.Bool(true),
	}, opts...)
	if err != nil {
		return err
	}

	apiStage, err := apigatewayv2.NewStage(ctx, name+"DefaultStage", &apigatewayv2.StageArgs{
		AutoDeploy: pulumi.BoolPtr(true),
		Name:       pulumi.String("$default"),
		ApiId:      a.apis[name].ID(),
		// Tags:       pulumi.ToStringMap(common.Tags(args.StackID, name+"DefaultStage", resources.API)),
	}, opts...)
	if err != nil {
		return err
	}

	// Generate permissions enabling the API Gateway to invoke the functions it targets
	for fName, fun := range nitricServiceTargets {
		_, err = lambda.NewPermission(ctx, name+fName, &lambda.PermissionArgs{
			Function:  fun.Name,
			Action:    pulumi.String("lambda:InvokeFunction"),
			Principal: pulumi.String("apigateway.amazonaws.com"),
			SourceArn: pulumi.Sprintf("%s/*/*/*", a.apis[name].ExecutionArn),
		}, opts...)
		if err != nil {
			return err
		}
	}

	endPoint := a.apis[name].ApiEndpoint.ApplyT(func(ep string) string {
		return ep
	}).(pulumi.StringInput)

	if a.config.Apis[name] != nil {
		// For each specified domain name
		for _, domainName := range a.config.Apis[name].Domains {
			_, err := newDomainName(ctx, name, domainNameArgs{
				domainName: domainName,
				api:        a.apis[name],
				stage:      apiStage,
			})
			if err != nil {
				return err
			}
		}
	}

	ctx.Export("api:"+name, endPoint)

	return nil
}
