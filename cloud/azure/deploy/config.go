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

type AzureConfig struct {
	Refresh                                 bool
	Org                                     string `mapstructure:"org"`
	AdminEmail                              string `mapstructure:"adminemail"`
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

	if azureConfig.Config == nil {
		azureConfig.Config = map[string]*AzureConfigItem{}
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
