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

type GcpConfigItem struct {
	CloudRun  *GcpCloudRunConfig `mapstructure:",omitempty"`
	Telemetry int
}

type GcpCloudRunConfig struct {
	Memory       int
	Timeout      int
	MinInstances int `mapstructure:"min-instances"`
	MaxInstances int `mapstructure:"max-instances"`
	Concurrency  int
}

type GcpConfig struct {
	config.AbstractConfig[*GcpConfigItem] `mapstructure:"config,squash"`
	ScheduleTimezone                      string `mapstructure:"schedule-timezone"`
	ProjectId                             string `mapstructure:"project-id"`
	Refresh                               bool
}

var defaultCloudRunConfig = &GcpCloudRunConfig{
	Memory:       512,
	Timeout:      300,
	MinInstances: 0,
	MaxInstances: 80,
	Concurrency:  300,
}

var defaultGcpConfigItem = GcpConfigItem{
	Telemetry: 0,
}

// Return GcpConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (*GcpConfig, error) {
	err := config.ValidateRawConfigKeys(attributes, []string{"cloudrun"})
	if err != nil {
		return nil, err
	}

	gcpConfig := &GcpConfig{}
	err = mapstructure.Decode(attributes, gcpConfig)
	if err != nil {
		return nil, err
	}

	if gcpConfig.ScheduleTimezone == "" {
		gcpConfig.ScheduleTimezone = "UTC"
	}

	if gcpConfig.Config == nil {
		gcpConfig.Config = map[string]*GcpConfigItem{}
	}

	// if no default then set provider level defaults
	if _, hasDefault := gcpConfig.Config["default"]; !hasDefault {
		gcpConfig.Config["default"] = &defaultGcpConfigItem
		gcpConfig.Config["default"].CloudRun = defaultCloudRunConfig
	}

	for configName, configVal := range gcpConfig.Config {
		// Add omitted values from default configs where needed.
		err := mergo.Merge(configVal, defaultGcpConfigItem)
		if err != nil {
			return nil, err
		}

		if configVal.CloudRun == nil { // check if no runtime config provided, default to Lambda.
			configVal.CloudRun = defaultCloudRunConfig
		} else {
			err := mergo.Merge(configVal.CloudRun, defaultCloudRunConfig)
			if err != nil {
				return nil, err
			}
		}

		gcpConfig.Config[configName] = configVal
	}

	return gcpConfig, nil
}
