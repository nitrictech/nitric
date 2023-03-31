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

type GcpConfigItem struct {
	CloudRun  GcpCloudRunConfig
	Telemetry int
	Target    string
}

type GcpCloudRunConfig struct {
	Memory       int
	Timeout      int
	MinInstances int `mapstructure:"min-instances"`
	MaxInstances int `mapstructure:"max-instances"`
	Concurrency  int
}

type GcpConfig = config.AbstractConfig[GcpConfigItem]

var defaultGcpConfigItem = GcpConfigItem{
	CloudRun: GcpCloudRunConfig{
		Memory:       512,
		Timeout:      300,
		MinInstances: 0,
		MaxInstances: 80,
		Concurrency:  300,
	},
	Telemetry: 0,
	Target:    "cloudrun",
}

// Return GcpConfig from stack attributes
func ConfigFromAttributes(attributes map[string]interface{}) (*GcpConfig, error) {
	// Use common ConfigFromAttributes
	return config.ConfigFromAttributes(attributes, defaultGcpConfigItem)
}
