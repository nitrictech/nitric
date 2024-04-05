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
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/core/pkg/help"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/apigateway"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/common/deploy/utils"
)

type nameUrlPair struct {
	name      string
	invokeUrl string
}

func (p *NitricGcpPulumiProvider) Api(ctx *pulumi.Context, parent pulumi.Resource, name string, apiConfig *deploymentspb.Api) error {
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	if apiConfig.GetOpenapi() == "" {
		return fmt.Errorf("gcp provider can only deploy OpenAPI specs")
	}

	openapiDoc := &openapi3.T{}
	err := openapiDoc.UnmarshalJSON([]byte(apiConfig.GetOpenapi()))
	if err != nil {
		return fmt.Errorf("invalid document supplied for api: %s", name)
	}

	// augment document with security definitions
	for sn, sd := range openapiDoc.Components.SecuritySchemes {
		if sd.Value.Type == "openIdConnect" {
			// We need to extract audience values from the extensions
			// the extension is type of []interface and cannot be converted to []string directly
			audiences, err := utils.GetAudiencesFromExtension(sd.Value.Extensions)
			if err != nil {
				return err
			}

			oidConf, err := utils.GetOpenIdConnectConfig(sd.Value.OpenIdConnectUrl)
			if err != nil {
				return err
			}

			openapiDoc.Components.SecuritySchemes[sn] = &openapi3.SecuritySchemeRef{
				Value: &openapi3.SecurityScheme{
					Type: "oauth2",
					Flows: &openapi3.OAuthFlows{
						Implicit: &openapi3.OAuthFlow{
							AuthorizationURL: oidConf.AuthEndpoint,
						},
					},
					Extensions: map[string]interface{}{
						"x-google-issuer":    oidConf.Issuer,
						"x-google-jwks_uri":  oidConf.JwksUri,
						"x-google-audiences": strings.Join(audiences, ","),
					},
				},
			}
		} else {
			return fmt.Errorf("unsupported security definition supplied")
		}
	}

	v2doc, err := openapi2conv.FromV3(openapiDoc)
	if err != nil {
		return err
	}

	// Get service targets for IAM binding
	services := p.CloudRunServices

	for _, pi := range v2doc.Paths {
		for _, m := range []string{http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodPost, http.MethodPut} {
			if pi.GetOperation(m) == nil {
				continue
			}

			name, ok := keepOperation(pi.GetOperation(m).Extensions)
			if !ok {
				return fmt.Errorf("found operation missing nitric target property: %+v", pi.GetOperation(m).Extensions)
			}

			if _, ok := p.CloudRunServices[name]; !ok {
				return fmt.Errorf("unable to find target service %s in %+v", name, p.CloudRunServices)
			}

			services[name] = p.CloudRunServices[name]

			break
		}
	}

	nameUrlPairs := make([]interface{}, 0, len(services))

	// collect name arn pairs for output iteration
	for k, v := range services {
		nameUrlPairs = append(nameUrlPairs, pulumi.All(k, v.Url).ApplyT(func(args []interface{}) (nameUrlPair, error) {
			name, nameOk := args[0].(string)
			url, urlOk := args[1].(string)

			if !nameOk || !urlOk {
				return nameUrlPair{}, fmt.Errorf("invalid data %T %v", args, args)
			}

			return nameUrlPair{
				name:      name,
				invokeUrl: url,
			}, nil
		}))
	}

	// Now we need to create the document provided and interpolate the deployed service targets
	// i.e. their Urls...
	// Replace Nitric API Extensions with google api gateway extensions
	doc := pulumi.All(nameUrlPairs...).ApplyT(func(pairs []interface{}) (string, error) {
		naps := make(map[string]string)

		for _, p := range pairs {
			if pair, ok := p.(nameUrlPair); ok {
				naps[pair.name] = pair.invokeUrl
			} else {
				return "", fmt.Errorf("failed to resolve Cloud Run container URL for api %s, invalid name URL pair value %T %v, %s", name, p, p, help.BugInNitricHelpText())
			}
		}

		for k, p := range v2doc.Paths {
			p.Get = gcpOperation(name, p.Get, naps)
			p.Post = gcpOperation(name, p.Post, naps)
			p.Patch = gcpOperation(name, p.Patch, naps)
			p.Put = gcpOperation(name, p.Put, naps)
			p.Delete = gcpOperation(name, p.Delete, naps)
			p.Options = gcpOperation(name, p.Options, naps)
			v2doc.Paths[k] = p
		}

		b, err := v2doc.MarshalJSON()
		if err != nil {
			return "", err
		}

		return base64.StdEncoding.EncodeToString(b), nil
	}).(pulumi.StringOutput)

	resourceLabels := common.Tags(p.StackId, name, resources.API)

	api, err := apigateway.NewApi(ctx, name, &apigateway.ApiArgs{
		ApiId:  pulumi.String(name),
		Labels: pulumi.ToStringMap(resourceLabels),
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return errors.WithMessage(err, "api "+name)
	}

	svcAcct, err := NewServiceAccount(ctx, name+"-api-invoker", &GcpIamServiceAccountArgs{
		AccountId: name + "-api",
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return errors.WithMessage(err, "api serviceaccount "+name)
	}

	// Bind that IAM account as a member of all available service targets
	for _, serv := range services {
		iamName := fmt.Sprintf("%s-%s-binding", name, serv.Name)

		_, err = cloudrun.NewIamMember(ctx, iamName, &cloudrun.IamMemberArgs{
			Service:  serv.Service.Name,
			Location: serv.Service.Location,
			Member:   pulumi.Sprintf("serviceAccount:%s", svcAcct.ServiceAccount.Email),
			Role:     pulumi.String("roles/run.invoker"),
		}, p.WithDefaultResourceOptions(opts...)...)
		if err != nil {
			return errors.WithMessage(err, "api iamMember "+iamName)
		}
	}

	// Deploy the config
	config, err := apigateway.NewApiConfig(ctx, name+"-config", &apigateway.ApiConfigArgs{
		Project:     pulumi.String(p.GcpConfig.ProjectId),
		Api:         api.ApiId,
		DisplayName: pulumi.String(name + "-config"),
		OpenapiDocuments: apigateway.ApiConfigOpenapiDocumentArray{
			apigateway.ApiConfigOpenapiDocumentArgs{
				Document: apigateway.ApiConfigOpenapiDocumentDocumentArgs{
					Path:     pulumi.String("openapi.json"),
					Contents: doc,
				},
			},
		},
		GatewayConfig: apigateway.ApiConfigGatewayConfigArgs{
			BackendConfig: apigateway.ApiConfigGatewayConfigBackendConfigArgs{
				// Add the service account for the invoker here...
				GoogleServiceAccount: svcAcct.ServiceAccount.Email,
			},
		},
		Labels: pulumi.ToStringMap(resourceLabels),
	}, p.WithDefaultResourceOptions(append(opts, pulumi.ReplaceOnChanges([]string{"*"}))...)...)
	if err != nil {
		return errors.WithMessage(err, "api config")
	}

	// Deploy the gateway
	p.ApiGateways[name], err = apigateway.NewGateway(ctx, name+"-gateway", &apigateway.GatewayArgs{
		DisplayName: pulumi.String(name + "-gateway"),
		GatewayId:   pulumi.String(name + "-gateway"),
		ApiConfig:   pulumi.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", p.GcpConfig.ProjectId, api.ApiId, config.ApiConfigId),
		Labels:      pulumi.ToStringMap(resourceLabels),
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return errors.WithMessage(err, "api gateway")
	}

	// url := res.Gateway.DefaultHostname.ApplyT(func(hn string) string { return "https://" + hn })
	// ctx.Export("api:"+name, url)

	return nil
}

func keepOperation(opExt map[string]interface{}) (string, bool) {
	if opExt == nil {
		return "", false
	}

	name := ""

	if v, ok := opExt["x-nitric-target"]; ok {
		targetMap, isMap := v.(map[string]interface{})
		if isMap {
			name, _ = targetMap["name"].(string)
		}
	}

	if name == "" {
		return "", false
	}

	return name, true
}

func gcpOperation(apiName string, op *openapi2.Operation, urls map[string]string) *openapi2.Operation {
	if op == nil {
		return nil
	}

	name, ok := keepOperation(op.Extensions)
	if !ok {
		return nil
	}

	if _, ok := urls[name]; !ok {
		return nil
	}

	if s, ok := op.Extensions["x-nitric-security"]; ok {
		secName, isString := s.(string)

		if isString {
			op.Security = &openapi2.SecurityRequirements{
				{
					secName: {},
				},
			}
		}
	}

	for i, r := range op.Responses {
		if r.Description == "" {
			op.Responses[i].Description = name
		}
	}

	op.Extensions["x-google-backend"] = map[string]string{
		// Append the name of the target origin api gateway to the target address
		"address":          fmt.Sprintf("%s/x-nitric-api/%s", urls[name], apiName),
		"path_translation": "APPEND_PATH_TO_ADDRESS",
	}

	return op
}
