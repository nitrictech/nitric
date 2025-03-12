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

package deploy

import (
	"fmt"
	"strings"

	"github.com/nitrictech/nitric/cloud/azure/deploy/embeds"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/core/pkg/logger"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	apimanagement "github.com/pulumi/pulumi-azure-native-sdk/apimanagement/v2"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	commonutils "github.com/nitrictech/nitric/cloud/common/deploy/utils"
)

// const policyTemplate = `<policies><inbound><base /><set-backend-service base-url="https://%s" />%s<authentication-managed-identity resource="%s" client-id="%s" /><set-header name="X-Forwarded-Authorization" exists-action="override"><value>@(context.Request.Headers.GetValueOrDefault("Authorization",""))</value></set-header></inbound><backend><base /></backend><outbound><base /></outbound><on-error><base /></on-error></policies>`

// const jwtTemplate = `<validate-jwt header-name="Authorization" failed-validation-httpcode="401" failed-validation-error-message="Unauthorized. Access token is missing or invalid." require-expiration-time="false">
// <openid-config url="%s.well-known/openid-configuration" />
// <required-claims>
// 	<claim name="aud" match="any" separator=",">
// 		<value>%s</value>
// 	</claim>
// </required-claims>
// </validate-jwt>
// `

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

func setSecurityRequirements(secReq *openapi3.SecurityRequirements, secDef map[string]securityDefinition) ([]string, error) {
	jwtTemplates := []string{}

	for _, sec := range *secReq {
		for sn := range sec {
			if sd, ok := secDef[sn]; ok {
				jwtTemplate, err := embeds.GetApiJwtTemplate(embeds.JwtTemplateArgs{
					OidcUri:       sd.Issuer,
					RequiredClaim: strings.Join(sd.Audiences, ","),
				})
				if err != nil {
					return nil, err
				}
				jwtTemplates = append(jwtTemplates, jwtTemplate)
			}
		}
	}

	return jwtTemplates, nil
}

func (p *NitricAzurePulumiProvider) Api(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Api) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	additionalApiConfig := p.AzureConfig.Apis[name]

	openapiDoc := &openapi3.T{}
	err := openapiDoc.UnmarshalJSON([]byte(config.GetOpenapi()))
	if err != nil {
		return fmt.Errorf("invalid document supplied for api: %s", name)
	}

	if len(openapiDoc.Paths) < 1 {
		logger.Warnf("skipping deployment of API %s, no routes defined", name)
		return nil
	}

	managedIdentities := p.ContainerEnv.ManagedUser.ID().ToStringOutput().ApplyT(func(id string) apimanagement.UserIdentityPropertiesMapOutput {
		return apimanagement.UserIdentityPropertiesMap{
			id: nil,
		}.ToUserIdentityPropertiesMapOutput()
	}).(apimanagement.UserIdentityPropertiesMapOutput)

	serviceName := ResourceName(ctx, name, ApiManagementServiceRT)
	managedServiceId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-id", name), &random.RandomStringArgs{
		Length: pulumi.Int(4),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"name": name,
		}),
		Upper:   pulumi.Bool(false),
		Special: pulumi.Bool(false),
	})
	if err != nil {
		return err
	}

	managedServiceName := pulumi.Sprintf("%s%s", serviceName, managedServiceId.Result)

	mgmtService, err := apimanagement.NewApiManagementService(ctx, fmt.Sprintf("%s-mgmt", name), &apimanagement.ApiManagementServiceArgs{
		ServiceName:       managedServiceName,
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
		Tags: pulumi.ToStringMap(common.Tags(p.StackId, name, resources.API)),
	}, opts...)
	if err != nil {
		return err
	}

	displayName := name + "-api"
	if openapiDoc.Info != nil && openapiDoc.Info.Title != "" {
		displayName = openapiDoc.Info.Title
	}

	description := fmt.Sprintf("Nitric API Gateway for %s", p.StackId)

	if additionalApiConfig != nil && additionalApiConfig.Description != "" {
		description = additionalApiConfig.Description
	}

	openapiDoc.Info.Description = description

	b, err := marshalOpenAPISpec(openapiDoc)
	if err != nil {
		return err
	}

	api, err := apimanagement.NewApi(ctx, fmt.Sprintf("%s-api", name), &apimanagement.ApiArgs{
		Description:          pulumi.String(description),
		DisplayName:          pulumi.String(displayName),
		Protocols:            pulumi.StringArray{pulumi.String("https")},
		ApiId:                pulumi.String(name),
		Format:               pulumi.String("openapi+json"),
		Path:                 pulumi.String("/"),
		ResourceGroupName:    p.ResourceGroup.Name,
		SubscriptionRequired: pulumi.Bool(false),
		ServiceName:          mgmtService.Name,
		// No need to transform the original spec, the mapping occurs as part of the operation policies below
		Value: pulumi.String(string(b)),
	}, opts...)
	if err != nil {
		return err
	}

	p.Apis[name] = ApiResources{
		Api:                  api,
		ApiManagementService: mgmtService,
	}

	secDef := map[string]securityDefinition{}

	if openapiDoc.Components.SecuritySchemes != nil {
		// Start translating to Azure centric security schemes
		for apiName, scheme := range openapiDoc.Components.SecuritySchemes {
			// implement OpenIDConnect security
			if scheme.Value.Type == "openIdConnect" {
				// We need to extract audience values as well
				// lets use an extension to store these with the document
				audiences, err := commonutils.GetAudiencesFromExtension(scheme.Value.Extensions)
				if err != nil {
					return err
				}

				oidConf, err := commonutils.GetOpenIdConnectConfig(scheme.Value.OpenIdConnectUrl)
				if err != nil {
					return err
				}

				secDef[apiName] = securityDefinition{
					Audiences: audiences,
					Issuer:    oidConf.Issuer,
				}
			}
		}
	}

	if len(openapiDoc.Paths) < 1 {
		return fmt.Errorf("no paths defined in api: %s", name)
	}

	for _, pathItem := range openapiDoc.Paths {
		for _, op := range pathItem.Operations() {
			if v, ok := op.Extensions["x-nitric-target"]; ok {
				var jwtTemplates []string

				// Apply top level security
				if openapiDoc.Security != nil {
					jwtTemplates, err = setSecurityRequirements(&openapiDoc.Security, secDef)
					if err != nil {
						return err
					}
				}

				// Override with path security
				if op.Security != nil {
					jwtTemplates, err = setSecurityRequirements(op.Security, secDef)
					if err != nil {
						return err
					}
				}

				jwtTemplateString := strings.Join(jwtTemplates, "\n")
				target := ""

				targetMap, isMap := v.(map[string]interface{})
				if !isMap {
					return fmt.Errorf("operation: %s has malformed x-nitric-target annotation", op.OperationID)
				}

				target, isString := targetMap["name"].(string)
				if !isString {
					return fmt.Errorf("operation: %s has malformed x-nitric-target annotation", op.OperationID)
				}

				app, ok := p.ContainerApps[target]
				if !ok {
					return fmt.Errorf("unable to find container app for service: %s", target)
				}

				// this.api.id returns a URL path, which is the incorrect value here.
				//   We instead need the value passed to apiId in the api creation above.
				// However, we want to maintain the pulumi dependency, so we need to keep the 'apply' call.
				apiId := api.ID().ToStringOutput().ApplyT(func(id string) string {
					return name
				}).(pulumi.StringOutput)

				_ = ctx.Log.Info("op policy "+op.OperationID+" , name "+name, &pulumi.LogArgs{Ephemeral: true})

				policyValue := pulumi.All(app.App.LatestRevisionFqdn, app.Sp.ClientID, p.ContainerEnv.ManagedUser.ClientId).ApplyT(func(args []interface{}) (string, error) {
					policyTemplate, err := embeds.GetApiPolicyTemplate(embeds.ApiPolicyTemplateArgs{
						BackendHostName:         args[0].(string),
						ExtraPolicies:           jwtTemplateString,
						ManagedIdentityResource: args[1].(string),
						ManagedIdentityClientId: args[2].(string),
					})
					if err != nil {
						return "", err
					}

					return policyTemplate, nil
				}).(pulumi.StringOutput)

				_, err = apimanagement.NewApiOperationPolicy(ctx, ResourceName(ctx, name+"-"+op.OperationID, ApiOperationPolicyRT), &apimanagement.ApiOperationPolicyArgs{
					ResourceGroupName: p.ResourceGroup.Name,
					ApiId:             apiId,
					ServiceName:       mgmtService.Name,
					OperationId:       pulumi.String(op.OperationID),
					PolicyId:          pulumi.String("policy"),
					Format:            pulumi.String("xml"),
					// Value:             pulumi.Sprintf(policyTemplate, pulumi.Sprintf("%s%s%s", app.App.LatestRevisionFqdn, "/x-nitric-api/", name), jwtTemplateString, app.Sp.ClientID, p.ContainerEnv.ManagedUser.ClientId),
					Value: policyValue,
				}, pulumi.Parent(api))
				if err != nil {
					return errors.WithMessage(err, "NewApiOperationPolicy "+op.OperationID)
				}
			} else {
				return fmt.Errorf("operation: %s missing x-nitric-target annotation", op.OperationID)
			}
		}
	}

	return nil
}
