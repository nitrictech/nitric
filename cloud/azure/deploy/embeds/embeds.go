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

package embeds

import (
	_ "embed"
	"strings"
	"text/template"
)

//go:embed api-policy-template.xml
var apiPolicyTemplate string

type ApiPolicyTemplateArgs struct {
	BackendHostName         string
	ExtraPolicies           string
	ManagedIdentityResource string
	ManagedIdentityClientId string
}

func GetApiPolicyTemplate(args ApiPolicyTemplateArgs) (string, error) {
	tmpl, err := template.New("apiPolicyTemplate").Parse(apiPolicyTemplate)
	if err != nil {
		return "", err
	}

	var output strings.Builder

	err = tmpl.Execute(&output, args)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

type JwtTemplateArgs struct {
	OidcUri       string
	RequiredClaim string
}

//go:embed api-jwt-template.xml
var jwtPolicyTemplate string

func GetApiJwtTemplate(args JwtTemplateArgs) (string, error) {
	tmpl, err := template.New("jwtPolicyTemplate").Parse(jwtPolicyTemplate)
	if err != nil {
		return "", err
	}

	var output strings.Builder

	err = tmpl.Execute(&output, args)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}
