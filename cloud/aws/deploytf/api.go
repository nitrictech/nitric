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

package deploytf

import (
	"encoding/json"
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/service"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func awsOperation(op *openapi3.Operation, funcs map[string]*string) *openapi3.Operation {
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
		"uri":                  *arn,
	}

	return op
}

func (n *NitricAwsTerraformProvider) Api(stack cdktf.TerraformStack, name string, config *deploymentspb.Api) error {
	nameArnPairs := map[string]*string{}
	targetNames := map[string]*string{}

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

	nitricServiceTargets := map[string]service.Service{}
	for _, p := range openapiDoc.Paths {
		for _, op := range p.Operations() {
			if v, ok := op.Extensions["x-nitric-target"]; ok {
				if targetMap, isMap := v.(map[string]any); isMap {
					serviceName := targetMap["name"].(string)
					nitricServiceTargets[serviceName] = n.Services[serviceName]
				}
			}
		}
	}

	// collect name arn pairs for output iteration
	for k, v := range nitricServiceTargets {
		nameArnPairs[k] = v.InvokeArnOutput()
		targetNames[k] = v.LambdaFunctionNameOutput()
	}

	for k, p := range openapiDoc.Paths {
		p.Get = awsOperation(p.Get, nameArnPairs)
		p.Post = awsOperation(p.Post, nameArnPairs)
		p.Patch = awsOperation(p.Patch, nameArnPairs)
		p.Put = awsOperation(p.Put, nameArnPairs)
		p.Delete = awsOperation(p.Delete, nameArnPairs)
		p.Options = awsOperation(p.Options, nameArnPairs)
		openapiDoc.Paths[k] = p
	}

	// TODO: Use common tags method and ensure it works with pointer templating
	openapiDoc.Tags = []*openapi3.Tag{{
		Name:       fmt.Sprintf("x-nitric-%s-name", *n.Stack.StackIdOutput()),
		Extensions: map[string]interface{}{"x-amazon-apigateway-tag-value": name},
	}, {
		Name:       fmt.Sprintf("x-nitric-%s-type", *n.Stack.StackIdOutput()),
		Extensions: map[string]interface{}{"x-amazon-apigateway-tag-value": "api"},
	}}

	b, err := json.Marshal(openapiDoc)
	if err != nil {
		return err
	}

	domains := []string{}
	if n.AwsConfig != nil && n.AwsConfig.Apis != nil && n.AwsConfig.Apis[name] != nil {
		domains = n.AwsConfig.Apis[name].Domains
	}

	n.Apis[name] = api.NewApi(stack, jsii.Sprintf("api_%s", name), &api.ApiConfig{
		Name:                  jsii.String(name),
		Spec:                  jsii.String(string(b)),
		TargetLambdaFunctions: &targetNames,
		Domains:               jsii.Strings(domains...),
		StackId:               n.Stack.StackIdOutput(),
	})

	return nil
}
