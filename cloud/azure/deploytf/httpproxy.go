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
	"fmt"
	"strings"

	"github.com/aws/jsii-runtime-go"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/http_proxy"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

const proxyTemplate = `<policies>
	<inbound>
		<base />
		<set-backend-service base-url="https://%s"/>
		<authentication-managed-identity resource="%s" client-id="%s" />
		<set-header name="X-Forwarded-Authorization" exists-action="override">
			<value>@(context.Request.Headers.GetValueOrDefault("Authorization",""))</value>
		</set-header>
	</inbound>
	<backend>
		<base />
	</backend>
	<outbound>
		<base />
	</outbound>
	<on-error>
		<base />
	</on-error>
</policies>`

func (n *NitricAzureTerraformProvider) Http(stack cdktf.TerraformStack, name string, config *deploymentspb.Http) error {
	operationPolicyTemplate := map[string]*string{}

	spec := newApiSpec(name)

	dependsOnServices := map[string]cdktf.ITerraformDependable{}

	for _, path := range spec.Paths {
		for _, op := range path.Operations() {
			target := config.Target.GetService()
			service, ok := n.Services[target]
			if !ok {
				return fmt.Errorf("unable to find container app for service: %s in http proxy: %s", target, name)
			}

			dependsOnServices[target] = service

			operationPolicyTemplate[op.OperationID] = jsii.Sprintf(proxyTemplate, *service.FqdnOutput(), *service.ClientIdOutput(), *n.Stack.AppIdentityClientIdOutput())
		}
	}

	b, err := spec.MarshalJSON()
	if err != nil {
		return err
	}

	// Limit to 42 as the random stack ID that appends is 8 characters
	cleanName := strings.ReplaceAll(name, "_", "-")
	if len(cleanName) > 42 {
		cleanName = cleanName[:42]
	}

	dependsOn := []cdktf.ITerraformDependable{}
	for _, v := range dependsOnServices {
		dependsOn = append(dependsOn, v)
	}

	n.Proxies[name] = http_proxy.NewHttpProxy(stack, jsii.String(cleanName), &http_proxy.HttpProxyConfig{
		Name:                     jsii.String(cleanName),
		PublisherName:            jsii.String(n.AzureConfig.Org),
		PublisherEmail:           jsii.String(n.AzureConfig.AdminEmail),
		Location:                 jsii.String(n.Region),
		ResourceGroupName:        n.Stack.ResourceGroupNameOutput(),
		AppIdentity:              n.Stack.AppIdentityOutput(),
		Description:              jsii.Sprintf("Nitric HTTP Proxy for %s", *n.Stack.StackIdOutput()),
		OperationPolicyTemplates: &operationPolicyTemplate,
		DependsOn:                &dependsOn,

		// No need to transform the openapi spec, we can just pass it directly
		// We provide a separate array for the creation of operation policies for the API
		OpenapiSpec: jsii.String(string(b)),
		Tags:        n.GetTags(*n.Stack.StackIdOutput(), name, resources.HttpProxy),
	})

	return nil
}

func newApiSpec(name string) *openapi3.T {
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
			"/*": &openapi3.PathItem{
				Get:     getOperation("get"),
				Post:    getOperation("post"),
				Patch:   getOperation("patch"),
				Put:     getOperation("put"),
				Delete:  getOperation("delete"),
				Options: getOperation("options"),
			},
		},
	}

	return doc
}

func getOperation(operationId string) *openapi3.Operation {
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
	}
}
