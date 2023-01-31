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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/apigateway"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/common/deploy/utils"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
)

type ApiGatewayArgs struct {
	ProjectId           pulumi.StringInput
	StackID             pulumi.StringInput
	OpenAPISpec         *openapi2.T
	Functions           map[string]*exec.CloudRunner
	SecurityDefinitions map[string]*v1.ApiSecurityDefinition
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

type openIdConfig struct {
	JwksUri       string `json:"jwks_uri"`
	TokenEndpoint string `json:"token_endpoint"`
	AuthEndpoint  string `json:"authorization_endpoint"`
}

func getOpenIdConnectConfig(issuer string) (*openIdConfig, error) {
	// append well-known configuration to issuer
	url, err := url.Parse(issuer)
	if err != nil {
		return nil, err
	}

	url.Path = path.Join(url.Path, ".well-known/openid-configuration")

	// get the configuration document
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non 200 status retrieving openid-configuration: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	oidConf := &openIdConfig{}

	if err := json.Unmarshal(body, oidConf); err != nil {
		return nil, err
	}

	return oidConf, nil
}

func NewApiGateway(ctx *pulumi.Context, name string, args *ApiGatewayArgs, opts ...pulumi.ResourceOption) (*ApiGateway, error) {
	res := &ApiGateway{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:GcpApiGateway", name, res, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(res))

	// augment document with security definitions
	for sn, sd := range args.SecurityDefinitions {
		if args.OpenAPISpec.SecurityDefinitions == nil {
			args.OpenAPISpec.SecurityDefinitions = make(map[string]*openapi2.SecurityScheme)
		}

		if sd.GetJwt() != nil {
			oidConf, err := getOpenIdConnectConfig(sd.GetJwt().GetIssuer())
			if err != nil {
				return nil, err
			}

			args.OpenAPISpec.SecurityDefinitions[sn] = &openapi2.SecurityScheme{
				Type:             "oauth2",
				Flow:             "implicit",
				AuthorizationURL: oidConf.AuthEndpoint,
				Extensions: map[string]interface{}{
					"x-google-issuer":    sd.GetJwt().Issuer,
					"x-google-jwks_uri":  oidConf.JwksUri,
					"x-google-audiences": strings.Join(sd.GetJwt().GetAudiences(), ","),
				},
			}
		} else {
			return nil, fmt.Errorf("error deploying gateway: unsupported security definition provided")
		}
	}

	// Get service targets for IAM binding
	funcs := map[string]*exec.CloudRunner{}

	for _, pi := range args.OpenAPISpec.Paths {
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

		for k, p := range args.OpenAPISpec.Paths {
			p.Get = gcpOperation(p.Get, naps)
			p.Post = gcpOperation(p.Post, naps)
			p.Patch = gcpOperation(p.Patch, naps)
			p.Put = gcpOperation(p.Put, naps)
			p.Delete = gcpOperation(p.Delete, naps)
			p.Options = gcpOperation(p.Options, naps)
			args.OpenAPISpec.Paths[k] = p
		}

		b, err := args.OpenAPISpec.MarshalJSON()
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
		Project:     args.ProjectId,
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
		targetMap, isMap := v.(map[string]string)
		if isMap {
			name = targetMap["name"]
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
