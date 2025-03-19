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
	"slices"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn"
	"github.com/samber/lo"
)

type ApiGateway struct {
	GatewayURL *string `json:"gateway_url"`
}

// function to create a new cdn
func (n *NitricAzureTerraformProvider) NewCdn(tfstack cdktf.TerraformStack) cdn.Cdn {
	apiGateways := make(map[string]ApiGateway)

	sortedApiKeys := lo.Keys(n.Apis)
	slices.Sort(sortedApiKeys)

	for _, apiName := range sortedApiKeys {
		api := n.Apis[apiName]
		apiGateways[apiName] = ApiGateway{
			GatewayURL: api.ApiGatewayUrlOutput(),
		}
	}

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

	return cdn.NewCdn(tfstack, jsii.String("cdn"), &cdn.CdnConfig{
		StackName:                    n.Stack.StackNameOutput(),
		StorageAccountPrimaryWebHost: n.Stack.StorageAccountWebHostOutput(),
		ResourceGroupName:            n.Stack.ResourceGroupNameOutput(),
		CdnPurgePaths:                filesToPurgeMap,
		Apis:                         apiGateways,
		DependsOn:                    &dependsOn,
	})
}
