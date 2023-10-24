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
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy/config"
)

type AwsImports struct {
	// A map of nitric names to ARNs
	Secrets map[string]string
	Buckets map[string]string
}

type AwsConfig struct {
	ScheduleTimezone                      string `mapstructure:"schedule-timezone,omitempty"`
	Import                                AwsImports
	Refresh                               bool
	Apis                                  map[string]*ApiConfig
	config.AbstractConfig[*AwsConfigItem] `mapstructure:"config,squash"`
}

type AwsConfigItem struct {
	Lambda    *AwsLambdaConfig `mapstructure:",omitempty"`
	Telemetry int
}

type AwsLambdaVpcConfig struct {
	SubnetIds        []string `mapstructure:"subnet-ids"`
	SecurityGroupIds []string `mapstructure:"security-group-ids"`
}

type AwsLambdaConfig struct {
	Memory                int
	Timeout               int
	ProvisionedConcurreny int                 `mapstructure:"provisioned-concurrency"`
	Vpc                   *AwsLambdaVpcConfig `mapstructure:"vpc,omitempty"`
}

var defaultLambdaConfig = &AwsLambdaConfig{
	Memory:                128,
	Timeout:               15,
	ProvisionedConcurreny: 0,
}

var defaultAwsConfigItem = AwsConfigItem{
	Telemetry: 0,
}

// Return AwsConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (*AwsConfig, error) {
	// get config attributes
	err := config.ValidateRawConfigKeys(attributes, []string{"lambda"})
	if err != nil {
		return nil, err
	}

	awsConfig := &AwsConfig{}
	err = mapstructure.Decode(attributes, awsConfig)
	if err != nil {
		return nil, err
	}

	// Default timezone if not specified
	if awsConfig.ScheduleTimezone == "" {
		// default to UTC
		awsConfig.ScheduleTimezone = "UTC"
	}

	if awsConfig.Apis == nil {
		awsConfig.Apis = map[string]*ApiConfig{}
	}

	if awsConfig.Config == nil {
		awsConfig.Config = map[string]*AwsConfigItem{}
	}

	// if no default then set provider level defaults
	if _, hasDefault := awsConfig.Config["default"]; !hasDefault {
		awsConfig.Config["default"] = &defaultAwsConfigItem
		awsConfig.Config["default"].Lambda = defaultLambdaConfig
	}

	for configName, configVal := range awsConfig.Config {
		// Add omitted values from default configs where needed.
		err := mergo.Merge(configVal, defaultAwsConfigItem)
		if err != nil {
			return nil, err
		}

		if configVal.Lambda == nil { // check if no runtime config provided, default to Lambda.
			configVal.Lambda = defaultLambdaConfig
		} else {
			err := mergo.Merge(configVal.Lambda, defaultLambdaConfig)
			if err != nil {
				return nil, err
			}
		}

		awsConfig.Config[configName] = configVal
	}

	return awsConfig, nil
}
