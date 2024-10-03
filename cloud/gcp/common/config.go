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
	"strings"

	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy/config"
)

type GcpConfigItem struct {
	CloudRun  *GcpCloudRunConfig `mapstructure:",omitempty"`
	Telemetry int
}

type GcpDatabaseConfig struct {
	DeletionPolicy string `mapstructure:"deletion-policy"`
}

type GcpImports struct {
	// A map of nitric names to GCP Secret Manager names
	Secrets map[string]string
}

type GcpBatchCompute struct {
	// Accelerator to use for the batch compute resources when GPUs are required
	AcceleratorType string `mapstructure:"accelerator-type"`
}

type GcpCloudRunConfig struct {
	Cpus         float64
	Memory       int
	Timeout      int
	MinInstances int `mapstructure:"min-instances"`
	MaxInstances int `mapstructure:"max-instances"`
	Concurrency  int
}

type GcpApiConfig struct {
	Description string
}

type GcpConfig struct {
	config.AbstractConfig[*GcpConfigItem] `mapstructure:"config,squash"`
	Apis                                  map[string]*GcpApiConfig
	Databases                             map[string]*GcpDatabaseConfig `mapstructure:"databases"`
	Import                                GcpImports
	ScheduleTimezone                      string           `mapstructure:"schedule-timezone"`
	ProjectId                             string           `mapstructure:"gcp-project-id"`
	GcpBatchCompute                       *GcpBatchCompute `mapstructure:"batch-compute"`
	Refresh                               bool
}

var defaultCloudRunConfig = &GcpCloudRunConfig{
	Cpus:         1,
	Memory:       512,
	Timeout:      300,
	MinInstances: 0,
	MaxInstances: 80,
	Concurrency:  300,
}

var defaultGcpBatchCompute = &GcpBatchCompute{
	AcceleratorType: "nvidia-tesla-t4",
}

var defaultGcpConfigItem = GcpConfigItem{
	Telemetry: 0,
}

var defaultSqlDatabaseItem = &GcpDatabaseConfig{
	DeletionPolicy: "DELETE",
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

	if gcpConfig.Apis == nil {
		gcpConfig.Apis = map[string]*GcpApiConfig{}
	}

	if gcpConfig.Config == nil {
		gcpConfig.Config = map[string]*GcpConfigItem{}
	}

	if gcpConfig.GcpBatchCompute == nil {
		gcpConfig.GcpBatchCompute = defaultGcpBatchCompute
	}

	// Add omitted values from default configs where needed.
	err = mergo.Merge(gcpConfig.GcpBatchCompute, defaultGcpBatchCompute)
	if err != nil {
		return nil, err
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

	for databaseName, configVal := range gcpConfig.Databases {
		// Add defaults to database config
		err := mergo.Merge(configVal, defaultSqlDatabaseItem)
		if err != nil {
			return nil, err
		}

		if strings.ToLower(configVal.DeletionPolicy) != "abandon" && strings.ToLower(configVal.DeletionPolicy) != "delete" {
			configVal.DeletionPolicy = "DELETE"
		}

		configVal.DeletionPolicy = strings.ToUpper(configVal.DeletionPolicy)

		gcpConfig.Databases[databaseName] = configVal
	}

	return gcpConfig, nil
}
