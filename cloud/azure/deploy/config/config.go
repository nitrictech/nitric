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

package config

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy/config"
)

type RawConfig = config.AbstractConfig[*RawConfigItem]

type RawConfigItem struct {
	Extras    map[string]any `mapstructure:",remain"`
	Telemetry int
}

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

type AzureConfig = config.AbstractConfig[*AzureConfigItem]

var defaultContainerAppsConfig = &AzureContainerAppsConfig{
	Cpu:         0.25,
	Memory:      0.5,
	MinReplicas: 0,
	MaxReplicas: 10,
}

var defaultAzureConfigItem = AzureConfigItem{
	Telemetry: 0,
}

// Return GcpConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (*AzureConfig, error) {
	rawConfig := RawConfig{}
	err := mapstructure.Decode(attributes, &rawConfig)
	if err != nil {
		return nil, err
	}

	for configName, configVal := range rawConfig.Config {
		if configVal == nil {
			return nil, fmt.Errorf("configuration key %s should not be empty", configName)
		}

		if len(configVal.Extras) > 1 {
			return nil, fmt.Errorf("config items should not contain more than one runtime config")
		}
	}

	gcpConfig := &AzureConfig{}
	err = mapstructure.Decode(attributes, gcpConfig)
	if err != nil {
		return nil, err
	}

	if gcpConfig.Config == nil {
		gcpConfig.Config = map[string]*AzureConfigItem{}
	}

	// if no default then set provider level defaults
	if _, hasDefault := gcpConfig.Config["default"]; !hasDefault {
		gcpConfig.Config["default"] = &defaultAzureConfigItem
		gcpConfig.Config["default"].ContainerApps = defaultContainerAppsConfig
	}

	for configName, configVal := range gcpConfig.Config {
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

		gcpConfig.Config[configName] = configVal
	}

	return gcpConfig, nil
}
