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

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"

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

type HttpProxyArgs struct {
	ProjectId string
	StackID   string
	Function  *exec.CloudRunner
}

type HttpProxy struct {
	pulumi.ResourceState

	Name    string
	Gateway *apigateway.Gateway
	Api     *apigateway.Api
}

func NewHttpProxy(ctx *pulumi.Context, name string, args *HttpProxyArgs, opts ...pulumi.ResourceOption) (*HttpProxy, error) {
	res := &HttpProxy{Name: name}

	err := ctx.RegisterComponentResource("nitric:api:GcpApiGateway", name, res, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(res))

	resourceLabels := common.Tags(args.StackID, name, resources.HttpProxy)

	res.Api, err = apigateway.NewApi(ctx, name, &apigateway.ApiArgs{
		ApiId:  pulumi.String(name),
		Labels: pulumi.ToStringMap(resourceLabels),
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

	// Bind that IAM account as a member of the function

	iamName := fmt.Sprintf("%s-%s-binding", name, args.Function.Name)

	_, err = cloudrun.NewIamMember(ctx, iamName, &cloudrun.IamMemberArgs{
		Service:  args.Function.Service.Name,
		Location: args.Function.Service.Location,
		Member:   pulumi.Sprintf("serviceAccount:%s", invoker.Email),
		Role:     pulumi.String("roles/run.invoker"),
	}, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "api iamMember "+iamName)
	}

	doc := args.Function.Url.ToStringOutput().ApplyT(func(url string) (string, error) {
		apiDoc := newApiSpec(name, url)

		v2doc, err := openapi2conv.FromV3(apiDoc)
		if err != nil {
			return "", err
		}

		b, err := v2doc.MarshalJSON()
		if err != nil {
			return "", err
		}

		return base64.StdEncoding.EncodeToString(b), nil
	}).(pulumi.StringOutput)

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
		Labels: pulumi.ToStringMap(resourceLabels),
	}, append(opts, pulumi.ReplaceOnChanges([]string{"*"}))...)
	if err != nil {
		return nil, errors.WithMessage(err, "api config")
	}

	// Deploy the gateway
	res.Gateway, err = apigateway.NewGateway(ctx, name+"-gateway", &apigateway.GatewayArgs{
		DisplayName: pulumi.String(name + "-gateway"),
		GatewayId:   pulumi.String(name + "-gateway"),
		ApiConfig:   pulumi.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", args.ProjectId, res.Api.ApiId, config.ApiConfigId),
		Labels:      pulumi.ToStringMap(resourceLabels),
	}, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "api gateway")
	}

	url := res.Gateway.DefaultHostname.ApplyT(func(hn string) string { return "https://" + hn })
	ctx.Export("api:"+name, url)

	return res, nil
}

func newApiSpec(name, functionUrl string) *openapi3.T {
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
			"/**": &openapi3.PathItem{
				Get:     getOperation(functionUrl, "get"),
				Post:    getOperation(functionUrl, "post"),
				Patch:   getOperation(functionUrl, "patch"),
				Put:     getOperation(functionUrl, "put"),
				Delete:  getOperation(functionUrl, "delete"),
				Options: getOperation(functionUrl, "options"),
			},
		},
	}

	return doc
}

func getOperation(functionUrl string, operationId string) *openapi3.Operation {
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
		Extensions: map[string]interface{}{
			"x-google-backend": map[string]string{
				"address":          functionUrl,
				"path_translation": "APPEND_PATH_TO_ADDRESS",
			},
		},
	}
}
