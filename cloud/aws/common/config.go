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
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
	"github.com/nitrictech/nitric/cloud/common/deploy/config"
)

type AwsApiConfig struct {
	Description string
	Domains     []string
}

type AwsCdnConfig struct {
	SkipCacheInvalidation bool `mapstructure:"skip-cache-invalidation"`
}

type AwsImports struct {
	// A map of nitric names to ARNs
	Secrets map[string]string
	Buckets map[string]string
}

type EcsLaunchTemplate struct {
	BlockDeviceMappings []struct {
		DeviceName string `mapstructure:"device-name,omitempty"`
		Ebs        struct {
			DeleteOnTermination string `mapstructure:"delete-on-termination,omitempty"`
			VolumeSize          int    `mapstructure:"volume-size,omitempty"`
			VolumeType          string `mapstructure:"volume-type,omitempty"`
		} `mapstructure:"ebs,omitempty"`
	} `mapstructure:"block-device-mappings,omitempty"`
}

type BatchComputeEnvConfig struct {
	MinCpus        int                `mapstructure:"min-cpus"`
	MaxCpus        int                `mapstructure:"max-cpus"`
	InstanceTypes  []string           `mapstructure:"instance-types"`
	LaunchTemplate *EcsLaunchTemplate `mapstructure:"launch-template,omitempty"`
}

type AuroraRdsClusterConfig struct {
	MinCapacity           float64 `mapstructure:"min-capacity"`
	MaxCapacity           float64 `mapstructure:"max-capacity"`
	SecondsUntilAutoPause *int    `mapstructure:"seconds-until-auto-pause"`
}

type AwsConfig struct {
	ScheduleTimezone                      string `mapstructure:"schedule-timezone,omitempty"`
	Import                                AwsImports
	Refresh                               bool
	Apis                                  map[string]*AwsApiConfig
	Cdn                                   *AwsCdnConfig           `mapstructure:"cdn,omitempty"`
	BatchComputeEnvConfig                 *BatchComputeEnvConfig  `mapstructure:"batch-compute-env,omitempty"`
	AuroraRdsClusterConfig                *AuroraRdsClusterConfig `mapstructure:"aurora-rds-cluster,omitempty"`
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
	EphemeralStorage      int                 `mapstructure:"ephemeral-storage"`
	ProvisionedConcurreny int                 `mapstructure:"provisioned-concurrency"`
	Vpc                   *AwsLambdaVpcConfig `mapstructure:"vpc,omitempty"`
}

var defaultLambdaConfig = &AwsLambdaConfig{
	Memory:                128,
	Timeout:               15,
	EphemeralStorage:      512,
	ProvisionedConcurreny: 0,
}

var defaultBatchComputeEnvConfig = &BatchComputeEnvConfig{
	MinCpus:        0,
	MaxCpus:        32,
	InstanceTypes:  []string{"optimal"},
	LaunchTemplate: nil,
}

var defaultAuroraRdsClusterConfig = &AuroraRdsClusterConfig{
	MinCapacity: 0.5,
	MaxCapacity: 1,
}

var defaultCdnConfig = &AwsCdnConfig{
	SkipCacheInvalidation: false,
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
		awsConfig.Apis = map[string]*AwsApiConfig{}
	}

	if awsConfig.Config == nil {
		awsConfig.Config = map[string]*AwsConfigItem{}
	}

	if awsConfig.Cdn == nil {
		awsConfig.Cdn = defaultCdnConfig
	}

	// if no default then set provider level defaults
	if _, hasDefault := awsConfig.Config["default"]; !hasDefault {
		awsConfig.Config["default"] = &defaultAwsConfigItem
		awsConfig.Config["default"].Lambda = defaultLambdaConfig
	}

	if awsConfig.BatchComputeEnvConfig == nil {
		awsConfig.BatchComputeEnvConfig = defaultBatchComputeEnvConfig
	}

	// merge in default values
	err = mergo.Merge(awsConfig.BatchComputeEnvConfig, defaultBatchComputeEnvConfig)
	if err != nil {
		return nil, err
	}

	if awsConfig.AuroraRdsClusterConfig == nil {
		awsConfig.AuroraRdsClusterConfig = defaultAuroraRdsClusterConfig
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
