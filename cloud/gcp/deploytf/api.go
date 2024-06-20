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
	"net/http"
	"strings"

	"github.com/aws/jsii-runtime-go"
	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/common/deploy/utils"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/service"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

type nameUrlPair struct {
	name           string
	invokeUrl      *string
	timeoutSeconds *float64
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

func (n *NitricGcpTerraformProvider) Api(stack cdktf.TerraformStack, name string, config *deploymentspb.Api) error {
	if config.GetOpenapi() == "" {
		return fmt.Errorf("gcp provider can only deploy OpenAPI specs")
	}

	openapiDoc := &openapi3.T{}
	err := openapiDoc.UnmarshalJSON([]byte(config.GetOpenapi()))
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
	services := map[string]service.Service{}

	for _, pi := range v2doc.Paths {
		for _, m := range []string{http.MethodGet, http.MethodPatch, http.MethodDelete, http.MethodPost, http.MethodPut} {
			if pi.GetOperation(m) == nil {
				continue
			}

			name, ok := keepOperation(pi.GetOperation(m).Extensions)
			if !ok {
				return fmt.Errorf("found operation missing nitric target property: %+v", pi.GetOperation(m).Extensions)
			}

			if _, ok := n.Services[name]; !ok {
				return fmt.Errorf("unable to find target service %s in %+v", name, n.Services)
			}

			services[name] = n.Services[name]

			break
		}
	}

	nameUrlPairs := make([]nameUrlPair, 0, len(services))

	// collect name arn pairs for output iteration
	for k, v := range services {
		nameUrlPairs = append(nameUrlPairs, nameUrlPair{
			name:           k,
			invokeUrl:      v.ServiceEndpointOutput(),
			timeoutSeconds: v.TimeoutSeconds(),
		})
	}

	naps := make(map[string]*string)
	timeouts := make(map[string]*float64)

	for _, p := range nameUrlPairs {
		naps[p.name] = p.invokeUrl
		timeouts[p.name] = p.timeoutSeconds
	}

	for k, p := range v2doc.Paths {
		p.Get = gcpOperation(name, p.Get, naps, timeouts)
		p.Post = gcpOperation(name, p.Post, naps, timeouts)
		p.Patch = gcpOperation(name, p.Patch, naps, timeouts)
		p.Put = gcpOperation(name, p.Put, naps, timeouts)
		p.Delete = gcpOperation(name, p.Delete, naps, timeouts)
		p.Options = gcpOperation(name, p.Options, naps, timeouts)
		v2doc.Paths[k] = p
	}

	b, err := v2doc.MarshalJSON()
	if err != nil {
		return err
	}

	serviceNames := map[string]*string{}
	for k, v := range services {
		serviceNames[k] = v.ServiceNameOutput()
	}

	dependableServices := []cdktf.ITerraformDependable{}
	for _, v := range services {
		dependableServices = append(dependableServices, v)
	}

	n.Apis[name] = api.NewApi(stack, jsii.Sprintf("api_%s", name), &api.ApiConfig{
		Name:           jsii.String(name),
		OpenapiSpec:    jsii.String(string(b)),
		TargetServices: &serviceNames,
		StackId:        n.Stack.StackIdOutput(),
		DependsOn:      &dependableServices,
	})

	return nil
}

func gcpOperation(apiName string, op *openapi2.Operation, urls map[string]*string, timeouts map[string]*float64) *openapi2.Operation {
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

	op.Extensions["x-google-backend"] = map[string]any{
		// Append the name of the target origin api gateway to the target address
		"address":          fmt.Sprintf("%s/x-nitric-api/%s", *urls[name], apiName),
		"path_translation": "APPEND_PATH_TO_ADDRESS",
		"deadline":         timeouts[name],
	}

	return op
}
