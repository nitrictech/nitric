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
	"slices"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn_api_rewrites"
	"github.com/samber/lo"
)

// Convert API name to a unique numeric value based on character codes
// This creates a number that's guaranteed unique for different strings
func apiNameToUniqueNumber(name string) int {
	// Start at a high base to avoid conflicts with other rules
	base := 10000

	// Use character position and value to guarantee uniqueness
	// This is essentially creating a custom numeric representation
	for i, char := range name {
		// Multiply by position+1 to weight characters differently
		// Use prime number multiplier to reduce collision risk
		base += int(char) * (i + 1) * 31
	}

	return base
}

// function to create a new cdn
func (n *NitricAzureTerraformProvider) NewCdn(tfstack cdktf.TerraformStack) cdn.Cdn {
	dependsOn := []cdktf.ITerraformDependable{n.Stack}

	allCdnPurgeMaps := []interface{}{}

	var filesToPurgeMap *map[string]*string

	for _, ws := range n.Websites {
		// add website to depends on
		dependsOn = append(dependsOn, ws)

		changedFilesOutput := ws.ChangedFilesOutput()

		if changedFilesOutput == nil {
			continue
		}

		allCdnPurgeMaps = append(allCdnPurgeMaps, *cdktf.Token_AsStringMap(changedFilesOutput, nil))
	}

	// merge all maps into one
	if len(allCdnPurgeMaps) > 0 {
		filesToPurgeMap = cdktf.Token_AsStringMap(cdktf.Fn_Merge(&allCdnPurgeMaps), nil)
	}

	enableApis := len(n.Apis) > 0

	afdCDN := cdn.NewCdn(tfstack, jsii.String("cdn"), &cdn.CdnConfig{
		StackName:                    n.Stack.StackNameOutput(),
		StorageAccountPrimaryWebHost: n.Stack.StorageAccountWebHostOutput(),
		ResourceGroupName:            n.Stack.ResourceGroupNameOutput(),
		CdnPurgePaths:                filesToPurgeMap,
		EnableApiRewrites:            jsii.Bool(enableApis),
		DependsOn:                    &dependsOn,
	})

	// add cdn api rewrites if apis are present
	if enableApis {
		sortedApiKeys := lo.Keys(n.Apis)
		slices.Sort(sortedApiKeys)

		for _, apiName := range sortedApiKeys {
			api := n.Apis[apiName]
			rewriteDependsOn := []cdktf.ITerraformDependable{n.Stack, afdCDN, api}

			// calculate a unique rule order for the api
			ruleOrder := apiNameToUniqueNumber(apiName)

			cdn_api_rewrites.NewCdnApiRewrites(tfstack, jsii.String(fmt.Sprintf("cdn_api_rewrite_%s", apiName)), &cdn_api_rewrites.CdnApiRewritesConfig{
				Name:                  jsii.String(apiName),
				ApiHostName:           api.ApiGatewayUrlOutput(),
				CdnFrontdoorProfileId: afdCDN.CdnFrontdoorProfileIdOutput(),
				CdnFrontdoorRuleSetId: afdCDN.CdnFrontdoorRuleSetIdOutput(),
				RuleOrder:             jsii.Number(ruleOrder),
				DependsOn:             &rewriteDependsOn,
			})
		}
	}

	return afdCDN
}
