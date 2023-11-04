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
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/nitrictech/nitric/cloud/common/cors"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	apimanagement "github.com/pulumi/pulumi-azure-native-sdk/apimanagement/v20201201"
	"github.com/pulumi/pulumi-azure-native-sdk/managedidentity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	commonutils "github.com/nitrictech/nitric/cloud/common/deploy/utils"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type AzureApiManagementArgs struct {
	StackID           string
	ResourceGroupName pulumi.StringInput
	OrgName           pulumi.StringInput
	AdminEmail        pulumi.StringInput
	OpenAPISpec       *openapi3.T
	Apps              map[string]*exec.ContainerApp
	ManagedIdentity   *managedidentity.UserAssignedIdentity
	Cors              *v1.ApiCorsDefinition
}

type AzureApiManagement struct {
	pulumi.ResourceState

	Name    string
	Api     *apimanagement.Api
	Service *apimanagement.ApiManagementService
}

const policyTemplate = `<policies><inbound><base /><set-backend-service base-url="https://%s" />%s<authentication-managed-identity resource="%s" client-id="%s" /><set-header name="X-Forwarded-Authorization" exists-action="override"><value>@(context.Request.Headers.GetValueOrDefault("Authorization",""))</value></set-header></inbound><backend><base /></backend><outbound><base /></outbound><on-error><base /></on-error></policies>`

const corsTemplate = `
<policies><inbound><base />
<cors allow-credentials="{{.AllowCredentials}}">
    {{if .AllowOrigins}}
    <allowed-origins>
        {{range .AllowOrigins}}<origin>{{.}}</origin>{{end}}
    </allowed-origins>
    {{end}}
    {{if .AllowMethods}}
    <allowed-methods preflight-result-max-age="{{.MaxAge}}">
        {{range .AllowMethods}}<method>{{.}}</method>{{end}}
    </allowed-methods>
    {{end}}
    {{if .AllowHeaders}}
    <allowed-headers>
        {{range .AllowHeaders}}<header>{{.}}</header>{{end}}
    </allowed-headers>
    {{end}}
    {{if .ExposeHeaders}}
    <expose-headers>
        {{range .ExposeHeaders}}<header>{{.}}</header>{{end}}
    </expose-headers>
    {{end}}
</cors>
</inbound><backend><base /></backend><outbound><base /></outbound><on-error><base /></on-error></policies>
`

const jwtTemplate = `<validate-jwt header-name="Authorization" failed-validation-httpcode="401" failed-validation-error-message="Unauthorized. Access token is missing or invalid." require-expiration-time="false">  
<openid-config url="%s.well-known/openid-configuration" />  
<required-claims>  
	<claim name="aud" match="any" separator=",">  
		<value>%s</value>  
	</claim>  
</required-claims>  
</validate-jwt>
`

func marshalOpenAPISpec(spec *openapi3.T) ([]byte, error) {
	sec := spec.Security
	spec.Security = openapi3.SecurityRequirements{}

	b, err := spec.MarshalJSON()

	spec.Security = sec

	return b, err
}

type securityDefinition struct {
	Issuer    string
	Audiences []string
}

func setSecurityRequirements(secReq *openapi3.SecurityRequirements, secDef map[string]securityDefinition) []string {
	jwtTemplates := []string{}

	for _, sec := range *secReq {
		for sn := range sec {
			if sd, ok := secDef[sn]; ok {
				jwtTemplates = append(jwtTemplates, fmt.Sprintf(jwtTemplate, sd.Issuer, strings.Join(sd.Audiences, ",")))
			}
		}
	}

	return jwtTemplates
}

func NewAzureApiManagement(ctx *pulumi.Context, name string, args *AzureApiManagementArgs, opts ...pulumi.ResourceOption) (*AzureApiManagement, error) {
	res := &AzureApiManagement{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:AzureApiManagement", name, res, opts...)
	if err != nil {
		return nil, err
	}

	managedIdentities := args.ManagedIdentity.ID().ToStringOutput().ApplyT(func(id string) apimanagement.UserIdentityPropertiesMapOutput {
		return apimanagement.UserIdentityPropertiesMap{
			id: nil,
		}.ToUserIdentityPropertiesMapOutput()
	}).(apimanagement.UserIdentityPropertiesMapOutput)

	res.Service, err = apimanagement.NewApiManagementService(ctx, utils.ResourceName(ctx, name, utils.ApiManagementServiceRT), &apimanagement.ApiManagementServiceArgs{
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
		Tags: pulumi.ToStringMap(common.Tags(args.StackID, name, resources.API)),
	})
	if err != nil {
		return nil, err
	}

	displayName := name + "-api"
	if args.OpenAPISpec.Info != nil && args.OpenAPISpec.Info.Title != "" {
		displayName = args.OpenAPISpec.Info.Title
	}

	b, err := marshalOpenAPISpec(args.OpenAPISpec)
	if err != nil {
		return nil, err
	}

	res.Api, err = apimanagement.NewApi(ctx, utils.ResourceName(ctx, name, utils.ApiRT), &apimanagement.ApiArgs{
		DisplayName:          pulumi.String(displayName),
		Protocols:            apimanagement.ProtocolArray{"https"},
		ApiId:                pulumi.String(name),
		Format:               pulumi.String("openapi+json"),
		Path:                 pulumi.String("/"),
		ResourceGroupName:    args.ResourceGroupName,
		SubscriptionRequired: pulumi.Bool(false),
		ServiceName:          res.Service.Name,
		// XXX: Do we need to stringify this?
		// Not need to transform the original spec,
		// the mapping occurs as part of the operation policies below
		Value: pulumi.String(string(b)),
	})
	if err != nil {
		return nil, err
	}

	secDef := map[string]securityDefinition{}

	if args.OpenAPISpec.Components.SecuritySchemes != nil {
		// Start translating to AWS centric security schemes
		for apiName, scheme := range args.OpenAPISpec.Components.SecuritySchemes {
			// implement OpenIDConnect security
			if scheme.Value.Type == "openIdConnect" {
				// We need to extract audience values as well
				// lets use an extension to store these with the document
				audiences, err := commonutils.GetAudiencesFromExtension(scheme.Value.Extensions)
				if err != nil {
					return nil, err
				}

				oidConf, err := commonutils.GetOpenIdConnectConfig(scheme.Value.OpenIdConnectUrl)
				if err != nil {
					return nil, err
				}

				secDef[apiName] = securityDefinition{
					Audiences: audiences,
					Issuer:    oidConf.Issuer,
				}
			}
		}
	}

	// this.api.id returns a URL path, which is the incorrect value here.
	// We instead need the value passed to apiId in the api creation above.
	// However, we want to maintain the pulumi dependency, so we need to keep the 'apply' call.
	apiId := res.Api.ID().ToStringOutput().ApplyT(func(id string) string {
		return name
	}).(pulumi.StringOutput)

	for _, pathItem := range args.OpenAPISpec.Paths {
		for _, op := range pathItem.Operations() {
			if v, ok := op.Extensions["x-nitric-target"]; ok {
				var jwtTemplates []string

				// Apply top level security
				if args.OpenAPISpec.Security != nil {
					jwtTemplates = setSecurityRequirements(&args.OpenAPISpec.Security, secDef)
				}

				// Override with path security
				if op.Security != nil {
					jwtTemplates = setSecurityRequirements(op.Security, secDef)
				}

				jwtTemplateString := strings.Join(jwtTemplates, "\n")
				target := ""

				targetMap, isMap := v.(map[string]interface{})
				if !isMap {
					continue
				}

				target, isString := targetMap["name"].(string)
				if !isString {
					continue
				}

				app, ok := args.Apps[target]
				if !ok {
					continue
				}

				_ = ctx.Log.Info("op policy "+op.OperationID+" , name "+name, &pulumi.LogArgs{Ephemeral: true})

				_, err = apimanagement.NewApiOperationPolicy(ctx, utils.ResourceName(ctx, name+"-"+op.OperationID, utils.ApiOperationPolicyRT), &apimanagement.ApiOperationPolicyArgs{
					ResourceGroupName: args.ResourceGroupName,
					ApiId:             apiId,
					ServiceName:       res.Service.Name,
					OperationId:       pulumi.String(op.OperationID),
					PolicyId:          pulumi.String("policy"),
					Format:            pulumi.String("xml"),
					Value:             pulumi.Sprintf(policyTemplate, app.App.LatestRevisionFqdn, jwtTemplateString, args.ManagedIdentity.ClientId, args.ManagedIdentity.ClientId),
				}, pulumi.Parent(res.Api))
				if err != nil {
					return nil, errors.WithMessage(err, "NewApiOperationPolicy "+op.OperationID)
				}
			}
		}
	}

	if args.Cors != nil {
		corsConfig, err := cors.GetCorsConfig(args.Cors)
		if err != nil {
			return nil, err
		}

		var resultBuffer bytes.Buffer
		t := template.Must(template.New("corsTemplate").Parse(corsTemplate))

		err = t.Execute(&resultBuffer, corsConfig)
		if err != nil {
			return nil, err
		}

		corsTemplateResult := resultBuffer.String()

		_, err = apimanagement.NewApiPolicy(ctx, utils.ResourceName(ctx, name+"-cors", utils.ApiOperationPolicyRT), &apimanagement.ApiPolicyArgs{
			ResourceGroupName: args.ResourceGroupName,
			ApiId:             apiId,
			ServiceName:       res.Service.Name,
			PolicyId:          pulumi.String("policy"),
			Format:            pulumi.String("xml"),
			Value:             pulumi.String(corsTemplateResult),
		})
		if err != nil {
			return nil, errors.WithMessage(err, "NewApiPolicy "+name+"-cors")
		}
	}

	ctx.Export("api:"+name, res.Service.GatewayUrl)

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":    pulumi.String(name),
		"service": res.Service,
		"api":     res.Api,
	})
}
