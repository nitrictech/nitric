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

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed api-url-rewrite.js
var cloudfront_ApiUrlRewriteFunction string

//go:embed url-rewrite.tmpl.js
var cloudfront_UrlRewriteFunctionName string

func GetApiUrlRewriteFunction() pulumi.StringInput {
	return pulumi.String(cloudfront_ApiUrlRewriteFunction)
}

func GetUrlRewriteFunction(basePath string) (pulumi.StringInput, error) {
	tmpl, err := template.New("rewrite-function").Parse(cloudfront_UrlRewriteFunctionName)
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"BasePath": basePath,
	}

	var output strings.Builder
	err = tmpl.Execute(&output, data)
	if err != nil {
		return nil, err
	}

	return pulumi.String(output.String()), nil
}
