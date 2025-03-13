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
	"github.com/nitrictech/nitric/cloud/azure/deploy/embeds"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	apimanagement "github.com/pulumi/pulumi-azure-native-sdk/apimanagement/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/managedidentity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

type AzureHttpProxyArgs struct {
	StackID           string
	ResourceGroupName pulumi.StringInput
	OrgName           pulumi.StringInput
	AdminEmail        pulumi.StringInput
	App               *ContainerApp
	ManagedIdentity   *managedidentity.UserAssignedIdentity
}

type AzureHttpProxy struct {
	pulumi.ResourceState

	Name    string
	Api     *apimanagement.Api
	Service *apimanagement.ApiManagementService
}

// const proxyTemplate = `<policies>
// 	<inbound>
// 		<base />
// 		<set-backend-service base-url="https://%s"/>
// 		<authentication-managed-identity resource="%s" client-id="%s" />
// 		<set-header name="X-Forwarded-Authorization" exists-action="override">
// 			<value>@(context.Request.Headers.GetValueOrDefault("Authorization",""))</value>
// 		</set-header>
// 	</inbound>
// 	<backend>
// 		<base />
// 	</backend>
// 	<outbound>
// 		<base />
// 	</outbound>
// 	<on-error>
// 		<base />
// 	</on-error>
// </policies>`

func (p *NitricAzurePulumiProvider) Http(ctx *pulumi.Context, parent pulumi.Resource, name string, http *deploymentspb.Http) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	managedIdentities := p.ContainerEnv.ManagedUser.ID().ToStringOutput().ApplyT(func(id string) apimanagement.UserIdentityPropertiesMapOutput {
		return apimanagement.UserIdentityPropertiesMap{
			id: nil,
		}.ToUserIdentityPropertiesMapOutput()
	}).(apimanagement.UserIdentityPropertiesMapOutput)

	mgmtService, err := apimanagement.NewApiManagementService(ctx, ResourceName(ctx, name, ApiManagementProxyRT), &apimanagement.ApiManagementServiceArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		PublisherEmail:    pulumi.String(p.AzureConfig.AdminEmail),
		PublisherName:     pulumi.String(p.AzureConfig.Org),
		Sku: apimanagement.ApiManagementServiceSkuPropertiesArgs{
			Name:     pulumi.String("Consumption"),
			Capacity: pulumi.Int(0),
		},
		Identity: &apimanagement.ApiManagementServiceIdentityArgs{
			Type:                   pulumi.String("UserAssigned"),
			UserAssignedIdentities: managedIdentities,
		},
		Tags: pulumi.ToStringMap(common.Tags(p.StackId, name, resources.HttpProxy)),
	}, opts...)
	if err != nil {
		return err
	}

	spec := newApiSpec(name)

	b, err := spec.MarshalJSON()
	if err != nil {
		return err
	}

	apiId := pulumi.String(name)

	proxyApi, err := apimanagement.NewApi(ctx, ResourceName(ctx, name, ApiHttpProxyRT), &apimanagement.ApiArgs{
		DisplayName:          pulumi.Sprintf("%s-api", name),
		Protocols:            pulumi.StringArray{pulumi.String("https")},
		ApiId:                apiId,
		Format:               pulumi.String("openapi+json"),
		Path:                 pulumi.String("/"),
		ResourceGroupName:    p.ResourceGroup.Name,
		SubscriptionRequired: pulumi.Bool(false),
		ServiceName:          mgmtService.Name,
		Value:                pulumi.String(string(b)),
	}, opts...)
	if err != nil {
		return err
	}

	p.HttpProxies[name] = ApiResources{
		Api:                  proxyApi,
		ApiManagementService: mgmtService,
	}

	targetContainerApp := p.ContainerApps[http.GetTarget().GetService()]

	apiPolicy := pulumi.All(targetContainerApp.App.LatestRevisionFqdn, targetContainerApp.Sp.ClientID, p.ContainerEnv.ManagedUser.ClientId).ApplyT(func(args []interface{}) (string, error) {
		backendHostName := args[0].(string)
		servicePrincipalClientId := args[1].(string)
		managedUserClientId := args[2].(string)

		policy, err := embeds.GetApiPolicyTemplate(embeds.ApiPolicyTemplateArgs{
			BackendHostName:         backendHostName,
			ManagedIdentityClientId: servicePrincipalClientId,
			ManagedIdentityResource: managedUserClientId,
		})
		if err != nil {
			return "", err
		}

		return policy, nil
	}).(pulumi.StringOutput)

	for _, path := range spec.Paths {
		for _, op := range path.Operations() {
			_, err = apimanagement.NewApiOperationPolicy(ctx, ResourceName(ctx, name+"-"+op.OperationID, ApiOperationPolicyRT), &apimanagement.ApiOperationPolicyArgs{
				ResourceGroupName: p.ResourceGroup.Name,
				ApiId:             apiId,
				ServiceName:       mgmtService.Name,
				OperationId:       pulumi.String(op.OperationID),
				PolicyId:          pulumi.String("policy"),
				Format:            pulumi.String("xml"),
				Value:             apiPolicy,
			}, pulumi.Parent(proxyApi), pulumi.DependsOn([]pulumi.Resource{proxyApi}))
			if err != nil {
				return errors.WithMessage(err, "NewApiOperationPolicy proxy")
			}
		}
	}

	/*
		res, ctx.RegisterResourceOutputs(res, pulumi.Map{
			"name":    pulumi.String(name),
			"service": res.Service,
			"api":     res.Api,
		})
	*/

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
