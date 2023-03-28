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
	"github.com/nitrictech/nitric/cloud/common/deploy/config"
)

type AzureConfigItem struct {
	ContainerApps AzureContainerAppsConfig
	Telemetry     int
	Target        string
}

type AzureContainerAppsConfig struct {
	Cpu         float64
	Memory      float64
	MinReplicas int `mapstructure:"min-replicas"`
	MaxReplicas int `mapstructure:"max-replicas"`
}

type AzureConfig = config.AbstractConfig[AzureConfigItem]

var defaultAzureConfigItem = AzureConfigItem{
	ContainerApps: AzureContainerAppsConfig{
		Cpu:         0.5,
		Memory:      0.5,
		MinReplicas: 0,
		MaxReplicas: 10,
	},
	Telemetry: 0,
	Target:    "containerapps",
}

// Return AzureConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (*AzureConfig, error) {
	// Use common ConfigFromAttributes
	return config.ConfigFromAttributes(attributes, defaultAzureConfigItem)
}
