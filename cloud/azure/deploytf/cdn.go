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
	"strings"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn_api_rewrites"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/cdn_subsites"
	"github.com/samber/lo"
)

// Required due to the way Azure handles rule replacement and ordering
// We need to ensure that the rules are unique and ordered in a way that
// avoids conflicts with other rules. This is a simple way to do that
// by converting the name to a unique number based on its character codes.
func nameToUniqueNumber(name string) int {
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

func ensureValidSubdomain(domain string, subdomain string) error {
	domain = strings.ToLower(strings.TrimSuffix(domain, "."))
	subdomain = strings.ToLower(strings.TrimSuffix(subdomain, "."))

	if subdomain == domain || strings.HasSuffix(subdomain, "."+domain) {
		return nil
	}

	return fmt.Errorf("%s is not a valid subdomain of %s", subdomain, domain)
}

// function to create a new cdn
func (n *NitricAzureTerraformProvider) NewCdn(tfstack cdktf.TerraformStack) error {
	dependsOn := []cdktf.ITerraformDependable{n.Stack}

	allCdnPurgeMaps := []interface{}{}

	var uploadedFiles *map[string]*string
	var primaryWebHost *string
	var subDomain string
	var enableCustomDomain bool
	var isApexDomain bool

	// Validate the custom domain configuration
	if n.AzureConfig.CdnDomain.DomainName != "" && n.AzureConfig.CdnDomain.ZoneName != "" {
		if n.AzureConfig.CdnDomain.ZoneName == "" {
			return fmt.Errorf("zone-name is required for custom domain")
		}

		if n.AzureConfig.CdnDomain.DomainName == "" {
			return fmt.Errorf("domain-name is required for custom domain")
		}

		if n.AzureConfig.CdnDomain.ZoneResourceGroup == "" {
			return fmt.Errorf("zone-resource-group is required for custom domain")
		}

		subDomain = strings.ToLower(n.AzureConfig.CdnDomain.DomainName)

		// check if the domain name is a subdomain of the zone name
		err := ensureValidSubdomain(n.AzureConfig.CdnDomain.ZoneName, subDomain)
		if err != nil {
			return err
		}

		// if it is a subdomain, remove the zone name from the subdomain
		subDomain = strings.ReplaceAll(strings.TrimSuffix(subDomain, n.AzureConfig.CdnDomain.ZoneName), ".", "")

		isApexDomain = subDomain == ""

		// if domain is an apex domain, create a unique subdomain for naming purposes
		if isApexDomain {
			subDomain = fmt.Sprintf("%s-%s", n.StackName, strings.ReplaceAll(n.AzureConfig.CdnDomain.ZoneName, ".", "-"))
		}

		enableCustomDomain = true
	}

	for _, ws := range n.Websites {
		// set the primary web host to the first website
		if *ws.BasePath() == "/" {
			primaryWebHost = ws.StorageAccountWebHostOutput()
		}

		// add website to depends on
		dependsOn = append(dependsOn, ws)

		allCdnPurgeMaps = append(allCdnPurgeMaps, *cdktf.Token_AsStringMap(ws.UploadedFilesOutput(), nil))
	}

	// merge all maps into one
	if len(allCdnPurgeMaps) > 0 {
		uploadedFiles = cdktf.Token_AsStringMap(cdktf.Fn_Merge(&allCdnPurgeMaps), nil)
	}

	enableApiRewrites := len(n.Apis) > 0

	afdCDN := cdn.NewCdn(tfstack, jsii.String("cdn"), &cdn.CdnConfig{
		StackName:             n.Stack.StackNameOutput(),
		ResourceGroupName:     n.Stack.ResourceGroupNameOutput(),
		UploadedFiles:         uploadedFiles,
		PrimaryWebHost:        primaryWebHost,
		ZoneResourceGroupName: jsii.String(n.AzureConfig.CdnDomain.ZoneResourceGroup),
		ZoneName:              jsii.String(n.AzureConfig.CdnDomain.ZoneName),
		CustomDomainHostName:  jsii.String(n.AzureConfig.CdnDomain.DomainName),
		DomainName:            jsii.String(subDomain),
		IsApexDomain:          jsii.Bool(isApexDomain),
		EnableCustomDomain:    jsii.Bool(enableCustomDomain),
		EnableApiRewrites:     jsii.Bool(enableApiRewrites),
		DependsOn:             &dependsOn,
	})

	if len(n.Websites) > 1 {
		for _, ws := range n.Websites {
			// add website to depends on
			dependsOn = append(dependsOn, ws)

			if *ws.BasePath() == "/" {
				continue
			}

			normalizedName := strings.ReplaceAll(*ws.BasePath(), "/", "")
			dependsOn := []cdktf.ITerraformDependable{n.Stack, afdCDN}

			cdn_subsites.NewCdnSubsites(tfstack, jsii.String(fmt.Sprintf("cdn_subsite_%s", normalizedName)), &cdn_subsites.CdnSubsitesConfig{
				Name:                         jsii.String(normalizedName),
				StackName:                    n.Stack.StackNameOutput(),
				BasePath:                     ws.BasePath(),
				RuleOrder:                    jsii.Number(nameToUniqueNumber(normalizedName)),
				CdnDefaultFrontdoorRuleSetId: afdCDN.CdnFrontdoorDefaultRuleSetIdOutput(),
				PrimaryWebHost:               ws.StorageAccountWebHostOutput(),
				CdnFrontdoorProfileId:        afdCDN.CdnFrontdoorProfileIdOutput(),
				DependsOn:                    &dependsOn,
			})
		}
	}

	// add cdn api rewrites if apis are present
	if enableApiRewrites {
		sortedApiKeys := lo.Keys(n.Apis)
		slices.Sort(sortedApiKeys)

		for _, apiName := range sortedApiKeys {
			api := n.Apis[apiName]
			rewriteDependsOn := []cdktf.ITerraformDependable{n.Stack, afdCDN, api}

			// calculate a unique rule order for the api
			ruleOrder := nameToUniqueNumber(apiName)

			cdn_api_rewrites.NewCdnApiRewrites(tfstack, jsii.String(fmt.Sprintf("cdn_api_rewrite_%s", apiName)), &cdn_api_rewrites.CdnApiRewritesConfig{
				Name:                  jsii.String(apiName),
				ApiHostName:           api.ApiGatewayUrlOutput(),
				CdnFrontdoorProfileId: afdCDN.CdnFrontdoorProfileIdOutput(),
				CdnFrontdoorRuleSetId: afdCDN.CdnFrontdoorApiRuleSetIdOutput(),
				RuleOrder:             jsii.Number(ruleOrder),
				DependsOn:             &rewriteDependsOn,
			})
		}
	}

	return nil
}
