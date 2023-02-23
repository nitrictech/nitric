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

package gateway

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/apigateway"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/common/deploy/utils"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
)

type ApiGatewayArgs struct {
	ProjectId       string
	StackID         pulumi.StringInput
	OpenAPISpec     *openapi3.T
	Functions       map[string]*exec.CloudRunner
	SecuritySchemes openapi3.SecuritySchemes
}

type ApiGateway struct {
	pulumi.ResourceState

	Name    string
	Gateway *apigateway.Gateway
	Api     *apigateway.Api
}

type nameUrlPair struct {
	name      string
	invokeUrl string
}


func NewApiGateway(ctx *pulumi.Context, name string, args *ApiGatewayArgs, opts ...pulumi.ResourceOption) (*ApiGateway, error) {
	res := &ApiGateway{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:GcpApiGateway", name, res, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(res))

	// augment document with security definitions
	for sn, sd := range args.OpenAPISpec.Components.SecuritySchemes {
		if sd.Value.Type == "openIdConnect" {
			// We need to extract audience values from the extensions
			// the extension is type of []interface and cannot be converted to []string directly
			audiences, err := utils.GetAudiencesFromExtension(sd.Value.Extensions)
			if err != nil {
				return nil, err
			}

			oidConf, err := utils.GetOpenIdConnectConfig(sd.Value.OpenIdConnectUrl)
			if err != nil {
				return nil, err
			}

			args.OpenAPISpec.Components.SecuritySchemes[sn] = &openapi3.SecuritySchemeRef{
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
			return nil, fmt.Errorf("unsupported security definition supplied")
		}
	}

	v2doc, err := openapi2conv.FromV3(args.OpenAPISpec)
	if err != nil {
		return nil, err
	}

	// Get service targets for IAM binding
	funcs := map[string]*exec.CloudRunner{}

	for _, pi := range v2doc.Paths {
		for _, m := range []string{http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodPost, http.MethodPut} {
			if pi.GetOperation(m) == nil {
				continue
			}

			name, ok := keepOperation(pi.GetOperation(m).Extensions)
			if !ok {
				continue
			}

			if _, ok := args.Functions[name]; !ok {
				continue
			}

			funcs[name] = args.Functions[name]

			break
		}
	}

	nameArnPairs := make([]interface{}, 0, len(args.Functions))

	// collect name arn pairs for output iteration
	for k, v := range args.Functions {
		nameArnPairs = append(nameArnPairs, pulumi.All(k, v.Url).ApplyT(func(args []interface{}) nameUrlPair {
			name := args[0].(string)
			url := args[1].(string)

			return nameUrlPair{
				name:      name,
				invokeUrl: url,
			}
		}))
	}

	// Now we need to create the document provided and interpolate the deployed service targets
	// i.e. their Urls...
	// Replace Nitric API Extensions with google api gateway extensions
	doc := pulumi.All(nameArnPairs...).ApplyT(func(pairs []interface{}) (string, error) {
		naps := make(map[string]string)

		for _, p := range pairs {
			if pair, ok := p.(nameUrlPair); ok {
				naps[pair.name] = pair.invokeUrl
			} else {
				// XXX: Should not occur
				return "", fmt.Errorf("invalid data %T %v", p, p)
			}
		}

		for k, p := range v2doc.Paths {
			p.Get = gcpOperation(p.Get, naps)
			p.Post = gcpOperation(p.Post, naps)
			p.Patch = gcpOperation(p.Patch, naps)
			p.Put = gcpOperation(p.Put, naps)
			p.Delete = gcpOperation(p.Delete, naps)
			p.Options = gcpOperation(p.Options, naps)
			v2doc.Paths[k] = p
		}

		b, err := v2doc.MarshalJSON()
		if err != nil {
			return "", err
		}

		return base64.StdEncoding.EncodeToString(b), nil
	}).(pulumi.StringOutput)

	res.Api, err = apigateway.NewApi(ctx, name, &apigateway.ApiArgs{
		ApiId:  pulumi.String(name),
		Labels: common.Tags(ctx, args.StackID, name),
	}, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "api "+name)
	}

	invoker, err := serviceaccount.NewAccount(ctx, name+"-acct", &serviceaccount.AccountArgs{
		AccountId: pulumi.String(utils.StringTrunc(name, 30-5) + "-acct"),
	}, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "api serviceaccount "+name)
	}

	// Bind that IAM account as a member of all available service targets
	for _, fun := range funcs {
		iamName := fmt.Sprintf("%s-%s-binding", name, fun.Name)

		_, err = cloudrun.NewIamMember(ctx, iamName, &cloudrun.IamMemberArgs{
			Service:  fun.Service.Name,
			Location: fun.Service.Location,
			Member:   pulumi.Sprintf("serviceAccount:%s", invoker.Email),
			Role:     pulumi.String("roles/run.invoker"),
		}, opts...)
		if err != nil {
			return nil, errors.WithMessage(err, "api iamMember "+iamName)
		}
	}

	// Deploy the config
	config, err := apigateway.NewApiConfig(ctx, name+"-config", &apigateway.ApiConfigArgs{
		Project:     pulumi.String(args.ProjectId),
		Api:         res.Api.ApiId,
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
				GoogleServiceAccount: invoker.Email,
			},
		},
		Labels: common.Tags(ctx, args.StackID, name),
	}, append(opts, pulumi.ReplaceOnChanges([]string{"*"}))...)
	if err != nil {
		return nil, errors.WithMessage(err, "api config")
	}

	// Deploy the gateway
	res.Gateway, err = apigateway.NewGateway(ctx, name+"-gateway", &apigateway.GatewayArgs{
		DisplayName: pulumi.String(name + "-gateway"),
		GatewayId:   pulumi.String(name + "-gateway"),
		ApiConfig:   pulumi.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", args.ProjectId, res.Api.ApiId, config.ApiConfigId),
		Labels:      common.Tags(ctx, args.StackID, name),
	}, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "api gateway")
	}

	url := res.Gateway.DefaultHostname.ApplyT(func(hn string) string { return "https://" + hn })
	ctx.Export("api:"+name, url)

	return res, nil
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

func gcpOperation(op *openapi2.Operation, urls map[string]string) *openapi2.Operation {
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
		"address":          urls[name],
		"path_translation": "APPEND_PATH_TO_ADDRESS",
	}

	return op
}
