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
	"strings"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"

	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

func (p *NitricGcpPulumiProvider) Http(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Http) error {
	opts := append([]pulumi.ResourceOption{}, pulumi.Parent(parent))

	resourceLabels := common.Tags(p.stackId, name, resources.HttpProxy)

	targetService := p.cloudRunServices[config.Target.GetService()]

	// normalise the name to match the required format
	// ^projects/([a-z0-9-]+)/locations/([a-z0-9-]+)(/([a-z]+)/([a-z0-9-.]+))+$

	normalizedName := strings.Replace(name, "_", "-", -1)

	api, err := apigateway.NewApi(ctx, normalizedName, &apigateway.ApiArgs{
		ApiId:  pulumi.String(normalizedName),
		Labels: pulumi.ToStringMap(resourceLabels),
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, "api "+name)
	}

	doc := targetService.Url.ToStringOutput().ApplyT(func(url string) (string, error) {
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
	apiConfig, err := apigateway.NewApiConfig(ctx, normalizedName+"-config", &apigateway.ApiConfigArgs{
		Project:     pulumi.String(p.config.ProjectId),
		Api:         api.ApiId,
		DisplayName: pulumi.String(normalizedName + "-config"),
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
				GoogleServiceAccount: targetService.Invoker.Email,
			},
		},
		Labels: pulumi.ToStringMap(resourceLabels),
	}, append(opts, pulumi.ReplaceOnChanges([]string{"*"}))...)
	if err != nil {
		return errors.WithMessage(err, "api config")
	}

	// Deploy the gateway
	_, err = apigateway.NewGateway(ctx, normalizedName+"-gateway", &apigateway.GatewayArgs{
		DisplayName: pulumi.String(normalizedName + "-gateway"),
		GatewayId:   pulumi.String(normalizedName + "-gateway"),
		ApiConfig:   pulumi.Sprintf("projects/%s/locations/global/apis/%s/configs/%s", p.config.ProjectId, api.ApiId, apiConfig.ApiConfigId),
		Labels:      pulumi.ToStringMap(resourceLabels),
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, "api gateway")
	}

	return nil
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
