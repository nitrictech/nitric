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
)

type AwsConfig struct {
	Config map[string]AwsConfigItem
}

type AwsConfigItem struct {
	Lambda    AwsLambdaConfig
	Telemetry int
	Target    string
}

type AwsLambdaConfig struct {
	Memory                int
	Timeout               int
	ProvisionedConcurreny int `mapstructure:"provisioned-concurrency"`
}

var defaultAwsConfig = &AwsConfig{
	Config: map[string]AwsConfigItem{
		"default": {
			Lambda: AwsLambdaConfig{
				Memory:                128,
				Timeout:               15,
				ProvisionedConcurreny: 0,
			},
			Telemetry: 0,
			Target:    "lambda",
		},
	},
}

// Return AwsConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (AwsConfig, error) {
	config := AwsConfig{}

	err := mapstructure.Decode(attributes, &config)

	// deep merge with defaults
	if err := mergo.Merge(&config, defaultAwsConfig); err != nil {
		return config, err
	}

	// capture default config and have other configs inherit it
	defaultConfig := config.Config["default"]

	// merge each no default key with defaults as well
	for name, val := range config.Config {
		if name == "default" {
			continue
		}

		defaultVal := defaultConfig

		if err := mergo.Merge(&defaultVal, &val, mergo.WithOverride); err != nil {
			return config, err
		}

		config.Config[name] = defaultVal
	}

	return config, err
}
