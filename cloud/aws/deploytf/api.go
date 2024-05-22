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
	"github.com/samber/lo"
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
	// Api config transformation logic
	// opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	// nameArnPairs := make([]interface{}, 0, len(n.Services))
	nameArnPairs := map[string]*string{}
	targetArns := []*string{}

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
		targetArns = append(targetArns, v.InvokeArnOutput())
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

	b, err := json.Marshal(openapiDoc)
	if err != nil {
		return err
	}

	domains := lo.Ternary(n.AwsConfig != nil && n.AwsConfig.Apis != nil && n.AwsConfig.Apis[name] != nil, n.AwsConfig.Apis[name].Domains, nil)
	if domains == nil {
		domains = []string{}
	}

	n.Apis[name] = api.NewApi(stack, jsii.String(name), &api.ApiConfig{
		Name:             jsii.String(name),
		Spec:             jsii.String(string(b)),
		TargetLambdaArns: &targetArns,
		Domains:          jsii.Strings(domains...),
		StackId:          n.Stack.StackIdOutput(),
	})

	return nil

}
