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

package common

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy/config"
)

type AzureConfigItem struct {
	ContainerApps *AzureContainerAppsConfig `mapstructure:"containerapps,omitempty"`
	Telemetry     int
}

type AzureContainerAppsConfig struct {
	Cpu         float64
	Memory      float64
	MinReplicas int `mapstructure:"min-replicas"`
	MaxReplicas int `mapstructure:"max-replicas"`
}

type AzureApiConfig struct {
	Description string
}

const (
	// cdnSkuStandard deploys websites using the Standard tier of Azure Front Door (default).
	cdnSkuStandard = "standard"
	// cdnSkuPremium deploys websites using the Premium tier of Azure Front Door.
	cdnSkuPremium = "premium"
)

type CdnDomainConfig struct {
	// Name of the managed zone in azure DNS zones
	ZoneName string `mapstructure:"zone-name"`
	// Resource group of the managed zone in azure DNS zones
	ZoneResourceGroup string `mapstructure:"zone-resource-group"`
	// Name of the domain name or sub domain to set for the CDN
	DomainName string `mapstructure:"domain-name"`
	// If true, CDN cache invalidation will not be run on deployment
	SkipCacheInvalidation bool `mapstructure:"skip-cache-invalidation"`
	// The Azure Front Door tier to deploy, either standard or premium
	Sku string `mapstructure:"sku"`
}

// FrontDoorSkuName returns the name Azure uses for the configured Front Door tier
func (c CdnDomainConfig) FrontDoorSkuName() string {
	if c.Sku == cdnSkuPremium {
		return "Premium_AzureFrontDoor"
	}

	return "Standard_AzureFrontDoor"
}

type AzureConfig struct {
	Refresh                                 bool
	ResourceGroup                           string `mapstructure:"resource-group"`
	Org                                     string `mapstructure:"org"`
	AdminEmail                              string `mapstructure:"adminemail"`
	Apis                                    map[string]*AzureApiConfig
	CdnDomain                               CdnDomainConfig `mapstructure:"cdn"`
	config.AbstractConfig[*AzureConfigItem] `mapstructure:"config,squash"`
}

var defaultContainerAppsConfig = &AzureContainerAppsConfig{
	Cpu:         0.25,
	Memory:      0.5,
	MinReplicas: 0,
	MaxReplicas: 10,
}

var defaultAzureConfigItem = AzureConfigItem{
	Telemetry: 0,
}

// Return AzureConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (*AzureConfig, error) {
	err := config.ValidateRawConfigKeys(attributes, []string{"containerapps"})
	if err != nil {
		return nil, err
	}

	azureConfig := &AzureConfig{}
	err = mapstructure.Decode(attributes, azureConfig)
	if err != nil {
		return nil, err
	}

	if azureConfig.AdminEmail == "" {
		azureConfig.AdminEmail = "unknown@example.com"
	}

	if azureConfig.Org == "" {
		azureConfig.Org = "unknown"
	}

	if azureConfig.Apis == nil {
		azureConfig.Apis = map[string]*AzureApiConfig{}
	}

	if azureConfig.Config == nil {
		azureConfig.Config = map[string]*AzureConfigItem{}
	}

	if azureConfig.CdnDomain.Sku == "" {
		azureConfig.CdnDomain.Sku = cdnSkuStandard
	}

	if azureConfig.CdnDomain.Sku != cdnSkuStandard && azureConfig.CdnDomain.Sku != cdnSkuPremium {
		return nil, fmt.Errorf("invalid cdn sku value %q: must be %q or %q", azureConfig.CdnDomain.Sku, cdnSkuStandard, cdnSkuPremium)
	}

	// if no default then set provider level defaults
	if _, hasDefault := azureConfig.Config["default"]; !hasDefault {
		azureConfig.Config["default"] = &defaultAzureConfigItem
		azureConfig.Config["default"].ContainerApps = defaultContainerAppsConfig
	}

	for configName, configVal := range azureConfig.Config {
		// Add omitted values from default configs where needed.
		err := mergo.Merge(configVal, defaultAzureConfigItem)
		if err != nil {
			return nil, err
		}

		if configVal.ContainerApps == nil { // check if no runtime config provided, default to Lambda.
			configVal.ContainerApps = defaultContainerAppsConfig
		} else {
			err := mergo.Merge(configVal.ContainerApps, defaultContainerAppsConfig)
			if err != nil {
				return nil, err
			}
		}

		azureConfig.Config[configName] = configVal
	}

	return azureConfig, nil
}
