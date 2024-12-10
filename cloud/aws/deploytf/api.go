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
	"os"
	"path/filepath"

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

	op.Extensions["x-amazon-apigateway-integration"] = map[string]string{
		"type":                 "aws_proxy",
		"httpMethod":           "POST",
		"payloadFormatVersion": "2.0",
		"uri":                  fmt.Sprintf("${%s}", name),
	}

	return op
}

func (n *NitricAwsTerraformProvider) Api(stack cdktf.TerraformStack, name string, config *deploymentspb.Api) error {
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
	for _, apiPath := range openapiDoc.Paths {
		for _, pathOperation := range apiPath.Operations() {
			if apiNitricTarget, ok := pathOperation.Extensions["x-nitric-target"]; ok {
				if targetMap, isMap := apiNitricTarget.(map[string]any); isMap {
					serviceName, ok := targetMap["name"].(string)
					if !ok {
						return fmt.Errorf("missing or invalid 'name' field in x-nitric-target for path %s on API %s", pathOperation.OperationID, name)
					}

					nitricService, ok := n.Services[serviceName]
					if !ok {
						return fmt.Errorf("service %s is registered for path %s on API %s, but that service does not exist in the project", serviceName, pathOperation.OperationID, name)
					}

					nitricServiceTargets[serviceName] = nitricService
				}
			}
		}
	}

	nameArnPairs := map[string]*string{}
	targetNames := map[string]*string{}

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
		Name:       "x-nitric-${stack_id}-name",
		Extensions: map[string]interface{}{"x-amazon-apigateway-tag-value": name},
	}, {
		Name:       "x-nitric-${stack_id}-type",
		Extensions: map[string]interface{}{"x-amazon-apigateway-tag-value": "api"},
	}}

	b, err := json.MarshalIndent(openapiDoc, "", "  ")
	if err != nil {
		return err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("./.nitric/%s.spec.json", name))
	if err != nil {
		return err
	}

	// Write out the spec to the .nitric tmp directory
	err = os.WriteFile(absPath, b, 0o600)
	if err != nil {
		return err
	}

	// Create a terraform asset that references the spec file
	asset := cdktf.NewTerraformAsset(stack, jsii.Sprintf("api_%s_spec", name), &cdktf.TerraformAssetConfig{
		Path:      jsii.String(absPath),
		AssetHash: jsii.String("nitric-api-spec"),
		Type:      cdktf.AssetType_FILE,
	})

	nameArnPairs["stack_id"] = n.Stack.StackIdOutput()

	templateFile := cdktf.Fn_Templatefile(asset.Path(), nameArnPairs)

	domains := []string{}
	if n.AwsConfig != nil && n.AwsConfig.Apis != nil && n.AwsConfig.Apis[name] != nil {
		domains = n.AwsConfig.Apis[name].Domains
	}

	n.Apis[name] = api.NewApi(stack, jsii.Sprintf("api_%s", name), &api.ApiConfig{
		Name:                  jsii.String(name),
		Spec:                  cdktf.Token_AsString(templateFile, &cdktf.EncodingOptions{}),
		TargetLambdaFunctions: &targetNames,
		Domains:               jsii.Strings(domains...),
		StackId:               n.Stack.StackIdOutput(),
	})

	return nil
}
