// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// NITRIC_HOME directory environment variable name, default value:  .nitric/
const NITRIC_HOME = "NITRIC_HOME"

// NITRIC_YAML filename environment variable name, default value: nitric.yaml
const NITRIC_YAML = "NITRIC_YAML"

// Provides a nitric stack definition
type NitricStack struct {
	Name        string                 `yaml:"name"`
	Collections map[string]interface{} `yaml:"collections"`
}

// Return true if the named collection is defined in the stack
func (s NitricStack) HasCollection(name string) bool {
	_, found := s.Collections[name]
	return found
}

// Create a Nitric Stack definition with default path
func NewStackDefault() (*NitricStack, error) {
	// Determine path
	filePath := GetEnv("NITRIC_HOME", ".nitric/") + GetEnv("NITRIC_YAML", "nitric.yaml")

	nitricStack, err := NewStack(filePath)
	if err != nil {
		return nil, err
	}
	return nitricStack, nil
}

// Create a new Nitric Stack definition
func NewStack(filename string) (*NitricStack, error) {
	if filename == "" {
		return nil, fmt.Errorf("provide non-blank filename")
	}

	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var stack NitricStack

	err = yaml.Unmarshal(source, &stack)
	if err != nil {
		return nil, err
	}

	return &stack, nil
}
