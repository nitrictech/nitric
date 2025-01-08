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
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/api"
	commonutils "github.com/nitrictech/nitric/cloud/common/deploy/utils"
	"github.com/nitrictech/nitric/core/pkg/logger"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type securityDefinition struct {
	Issuer    string   `json:"issuer"`
	Audiences []string `json:"audiences"`
}

const policyTemplate = `<policies><inbound><base /><set-backend-service base-url="https://%s" />%s<authentication-managed-identity resource="%s" client-id="%s" /><set-header name="X-Forwarded-Authorization" exists-action="override"><value>@(context.Request.Headers.GetValueOrDefault("Authorization",""))</value></set-header></inbound><backend><base /></backend><outbound><base /></outbound><on-error><base /></on-error></policies>`

const jwtTemplate = `<validate-jwt header-name="Authorization" failed-validation-httpcode="401" failed-validation-error-message="Unauthorized. Access token is missing or invalid." require-expiration-time="false">  
<openid-config url="%s.well-known/openid-configuration" />  
<required-claims>  
	<claim name="aud" match="any" separator=",">  
		<value>%s</value>  
	</claim>  
</required-claims>  
</validate-jwt>
`

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

func (n *NitricAzureTerraformProvider) Api(stack cdktf.TerraformStack, name string, config *deploymentspb.Api) error {
	// parse over the provided open api spec to translate so we can send it to the module
	additionalApiConfig := n.AzureConfig.Apis[name]

	openapiDoc := &openapi3.T{}
	err := openapiDoc.UnmarshalJSON([]byte(config.GetOpenapi()))
	if err != nil {
		return fmt.Errorf("invalid document supplied for api: %s", name)
	}

	if len(openapiDoc.Paths) < 1 {
		logger.Warnf("skipping deployment of API %s, no routes defined", name)
		return nil
	}

	description := jsii.Sprintf("Nitric API Gateway for %s", n.Stack.StackNameOutput())
	if additionalApiConfig != nil && additionalApiConfig.Description != "" {
		description = jsii.String(additionalApiConfig.Description)
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

	policyTemplates := map[string]*string{}

	// for _, pathItem := range openapiDoc.Paths {
	// 	for _, op := range pathItem.Operations() {
	// 		if v, ok := op.Extensions["x-nitric-target"]; ok {
	// 			var jwtTemplates []string

	// 			// Apply top level security
	// 			if openapiDoc.Security != nil {
	// 				jwtTemplates = setSecurityRequirements(&openapiDoc.Security, secDef)
	// 			}

	// 			// Override with path security
	// 			if op.Security != nil {
	// 				jwtTemplates = setSecurityRequirements(op.Security, secDef)
	// 			}

	// 			jwtTemplateString := strings.Join(jwtTemplates, "\n")
	// 			target := ""

	// 			targetMap, isMap := v.(map[string]interface{})
	// 			if !isMap {
	// 				return fmt.Errorf("operation: %s has malformed x-nitric-target annotation", op.OperationID)
	// 			}

	// 			target, isString := targetMap["name"].(string)
	// 			if !isString {
	// 				return fmt.Errorf("operation: %s has malformed x-nitric-target annotation", op.OperationID)
	// 			}

	// 			app, ok := p.Services[target]
	// 			if !ok {
	// 				return fmt.Errorf("Unable to find container app for service: %s", target)
	// 			}

	// 			policyTemplates[op.OperationID] = fmt.Sprintf(policyTemplate, fmt.Sprintf("%s%s%s", app.App.LatestRevisionFqdn, "/x-nitric-api/", name), jwtTemplateString, app.Sp.ClientID, p.ContainerEnv.ManagedUser.ClientId)

	// 			if err != nil {
	// 				return errors.WithMessage(err, "NewApiOperationPolicy "+op.OperationID)
	// 			}
	// 		} else {
	// 			return fmt.Errorf("operation: %s missing x-nitric-target annotation", op.OperationID)
	// 		}
	// 	}
	// }

	n.Apis[name] = api.NewApi(stack, jsii.String(name), &api.ApiConfig{
		PublisherName:            jsii.String(n.AzureConfig.Org),
		PublisherEmail:           jsii.String(n.AzureConfig.AdminEmail),
		Location:                 jsii.String(n.Region),
		ResourceGroupName:        n.Stack.ResourceGroupNameOutput(),
		AppIdentity:              n.Stack.AppIdentityOutput(),
		Description:              description,
		OperationPolicyTemplates: &policyTemplates,

		// No need to transform the openapi spec, we can just pass it directly
		// We provide a seperate array for the creation of operation policies for the API
		OpenapiSpec: jsii.String(config.GetOpenapi()),
	})

	// For all paths

	return fmt.Errorf("not implemented")
}
