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

type AbstractItem = any

type AbstractConfig[T AbstractItem] struct {
	Config map[string]T
}

// ConfigFromAttributes - Merges given attributes into a useable config, all types are updated with the provided default config item
func ConfigFromAttributes[T AbstractItem](attributes map[string]interface{}, defaultItem T) (*AbstractConfig[T], error) {
	config := new(AbstractConfig[T])

	err := mapstructure.Decode(attributes, config)

	// deep merge default type
	if err := mergo.Merge(config, &AbstractConfig[T]{Config: map[string]T{"default": defaultItem}}); err != nil {
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
