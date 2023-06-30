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
	"github.com/google/uuid"
	"github.com/pkg/errors"
	apimanagement "github.com/pulumi/pulumi-azure-native-sdk/apimanagement/v20201201"
	"github.com/pulumi/pulumi-azure-native-sdk/managedidentity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

type AzureHttpProxyArgs struct {
	StackID           pulumi.StringInput
	ResourceGroupName pulumi.StringInput
	OrgName           pulumi.StringInput
	AdminEmail        pulumi.StringInput
	App               *exec.ContainerApp
	ManagedIdentity   *managedidentity.UserAssignedIdentity
}

type AzureHttpProxy struct {
	pulumi.ResourceState

	Name    string
	Api     *apimanagement.Api
	Service *apimanagement.ApiManagementService
}

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

func NewAzureHttpProxy(ctx *pulumi.Context, name string, args *AzureHttpProxyArgs, opts ...pulumi.ResourceOption) (*AzureHttpProxy, error) {
	res := &AzureHttpProxy{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:AzureApiManagement", name, res, opts...)
	if err != nil {
		return nil, err
	}

	managedIdentities := args.ManagedIdentity.ID().ToStringOutput().ApplyT(func(id string) apimanagement.UserIdentityPropertiesMapOutput {
		return apimanagement.UserIdentityPropertiesMap{
			id: nil,
		}.ToUserIdentityPropertiesMapOutput()
	}).(apimanagement.UserIdentityPropertiesMapOutput)

	res.Service, err = apimanagement.NewApiManagementService(ctx, utils.ResourceName(ctx, name, utils.ApiManagementProxyRT), &apimanagement.ApiManagementServiceArgs{
		ResourceGroupName: args.ResourceGroupName,
		PublisherEmail:    args.AdminEmail,
		PublisherName:     args.OrgName,
		Sku: apimanagement.ApiManagementServiceSkuPropertiesArgs{
			Name:     pulumi.String("Consumption"),
			Capacity: pulumi.Int(0),
		},
		Identity: &apimanagement.ApiManagementServiceIdentityArgs{
			Type:                   pulumi.String("UserAssigned"),
			UserAssignedIdentities: managedIdentities,
		},
		Tags: common.Tags(ctx, args.StackID, name),
	})
	if err != nil {
		return nil, err
	}

	spec := newApiSpec(name)

	b, err := spec.MarshalJSON()
	if err != nil {
		return nil, err
	}

	res.Api, err = apimanagement.NewApi(ctx, utils.ResourceName(ctx, name, utils.ApiHttpProxyRT), &apimanagement.ApiArgs{
		DisplayName:          pulumi.Sprintf("%s-api", name),
		Protocols:            apimanagement.ProtocolArray{"https"},
		ApiId:                pulumi.String(name),
		Format:               pulumi.String("openapi+json"),
		Path:                 pulumi.String("/"),
		ResourceGroupName:    args.ResourceGroupName,
		SubscriptionRequired: pulumi.Bool(false),
		ServiceName:          res.Service.Name,
		Value:                pulumi.String(string(b)),
	})
	if err != nil {
		return nil, err
	}

	// this.api.id returns a URL path, which is the incorrect value here.
	//   We instead need the value passed to apiId in the api creation above.
	// However, we want to maintain the pulumi dependency, so we need to keep the 'apply' call.
	apiId := res.Api.ID().ToStringOutput().ApplyT(func(id string) string {
		return name
	}).(pulumi.StringOutput)
	for _, p := range spec.Paths {
		for _, op := range p.Operations() {
			_, err = apimanagement.NewApiOperationPolicy(ctx, utils.ResourceName(ctx, name+"-"+op.OperationID, utils.ApiOperationPolicyRT), &apimanagement.ApiOperationPolicyArgs{
				ResourceGroupName: args.ResourceGroupName,
				ApiId:             apiId,
				ServiceName:       res.Service.Name,
				OperationId:       pulumi.String(op.OperationID),
				PolicyId:          pulumi.String("policy"),
				Format:            pulumi.String("xml"),
				Value:             pulumi.Sprintf(proxyTemplate, args.App.App.LatestRevisionFqdn, args.ManagedIdentity.ClientId, args.ManagedIdentity.ClientId),
			}, pulumi.Parent(res.Api))
			if err != nil {
				return nil, errors.WithMessage(err, "NewApiOperationPolicy proxy")
			}
		}
	}

	ctx.Export("api:"+name, res.Service.GatewayUrl)

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":    pulumi.String(name),
		"service": res.Service,
		"api":     res.Api,
	})
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
				Get:     getOperation(),
				Post:    getOperation(),
				Patch:   getOperation(),
				Put:     getOperation(),
				Delete:  getOperation(),
				Options: getOperation(),
			},
		},
	}

	return doc
}

func getOperation() *openapi3.Operation {
	defaultDescription := "default description"

	return &openapi3.Operation{
		OperationID: uuid.NewString(),
		Responses: openapi3.Responses{
			"default": &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: &defaultDescription,
				},
			},
		},
	}
}
