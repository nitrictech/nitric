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
	"reflect"
	"sort"

	"gopkg.in/yaml.v2"
)

// NITRIC_HOME directory environment variable name, default value:  .nitric/
const NITRIC_HOME = "NITRIC_HOME"

// NITRIC_YAML filename environment variable name, default value: nitric.yaml
const NITRIC_YAML = "NITRIC_YAML"

// Provides a nitric stack collections definition
type Collection struct {
	Attributes map[string]interface{} `yaml:"attributes"`
	Indexes    map[string]interface{} `yaml:"indexes"`
}

// Provides a nitric stack definition
type NitricStack struct {
	Name        string                `yaml:"name"`
	Collections map[string]Collection `yaml:"collections"`
}

// Return true if the named collection is defined in the stack
func (s NitricStack) HasCollection(name string) bool {
	_, found := s.Collections[name]
	return found
}

// Return the collection index names
func (s NitricStack) CollectionIndexes(name string) ([]string, error) {
	if !s.HasCollection(name) {
		return nil, fmt.Errorf("%v collections: %v: not found", s.Name, name)
	}

	for idxName, value := range s.Collections[name].Indexes {
		if idxName == "unique" {
			return []string{fmt.Sprintf("%v", value)}, nil
		} else if idxName == "composite" {
			return marshalInterfaceSlice(value), nil
		}
	}

	return nil, fmt.Errorf("%v collections: %v: indexes: not found", s.Name, name)
}

// Return the unique index attributue name for the specified collection
func (s NitricStack) CollectionIndexesUnique(name string) (string, error) {
	if !s.HasCollection(name) {
		return "", fmt.Errorf("%v collections: %v: not found", s.Name, name)
	}

	for idxName, value := range s.Collections[name].Indexes {
		if idxName == "unique" {
			return fmt.Sprintf("%v", value), nil
		}
	}

	return "", fmt.Errorf("%v collections: %v: indexes: unique: not found", s.Name, name)
}

// Return the composite index attributue names for the specified collection
func (s NitricStack) CollectionIndexesComposite(name string) ([]string, error) {
	if !s.HasCollection(name) {
		return nil, fmt.Errorf("%v collections: %v: not found", s.Name, name)
	}

	for idxName, value := range s.Collections[name].Indexes {
		if idxName == "composite" {
			return marshalInterfaceSlice(value), nil
		}
	}

	return nil, fmt.Errorf("%v collections: %v: indexes: composite: not found", s.Name, name)
}

// Return the colleciton attribute names
func (s NitricStack) CollectionAttributes(name string) ([]string, error) {
	if !s.HasCollection(name) {
		return nil, fmt.Errorf("%v collections: %v: not found", s.Name, name)
	}

	var names []string

	for key := range s.Collections[name].Attributes {
		names = append(names, key)
	}

	sort.Strings(names)

	return names, nil
}

// Return the collection attribute names
func (s NitricStack) CollectionFilterAttributes(colName string) ([]string, error) {
	names, err := s.CollectionAttributes(colName)
	if err != nil {
		return nil, err
	}

	// Remove value attribute which is not filtered
	if i := IndexOf(names, "value"); i != -1 {
		names = Remove(names, i)
	}

	// If has unique indexes remove any value, e.g. "key"
	index, err := s.CollectionIndexesUnique(colName)
	if err == nil {
		if i := IndexOf(names, index); i != -1 {
			names = Remove(names, i)
		}
	}

	// If has composite indexes any values, e.g. ["pk" "sk"]
	indexes, err := s.CollectionIndexesComposite(colName)
	if err == nil {
		for _, index := range indexes {
			if i := IndexOf(names, index); i != -1 {
				names = Remove(names, i)
			}
		}
	}

	return names, nil
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

	// Configure default collection attributes and indexes
	for colName, collection := range stack.Collections {
		if collection.Attributes == nil && collection.Indexes == nil {
			collection.Attributes = make(map[string]interface{})
			collection.Attributes["key"] = "string"
			collection.Attributes["value"] = "string"
			collection.Indexes = make(map[string]interface{})
			collection.Indexes["unique"] = "key"
			stack.Collections[colName] = collection
		}
	}

	err = validateCollections(stack)
	if err != nil {
		return nil, err
	}

	return &stack, nil
}

func hasAttribute(name interface{}, attributes map[string]interface{}) bool {
	for attName := range attributes {
		if attName == name {
			return true
		}
	}
	return false
}

func marshalInterfaceSlice(value interface{}) []string {
	s := reflect.ValueOf(value)
	if s.Kind() != reflect.Slice {
		// Should not occur
		panic("marshalInterfaceSlice() given a non-slice type")
	}

	valueSlice := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		valueSlice[i] = s.Index(i).Interface()
	}

	strings := make([]string, len(valueSlice))
	for i, v := range valueSlice {
		strings[i] = fmt.Sprint(v)
	}

	return strings
}

func validateCollections(stack NitricStack) error {

	// Validate collections
	for colName, collection := range stack.Collections {

		// Ensure indexes defined
		if len(collection.Attributes) > 0 && len(collection.Indexes) == 0 {
			return fmt.Errorf("%s collections: %v: has no indexes:", stack.Name, colName)
		}
		// Ensure attributes defined
		if len(collection.Attributes) == 0 && len(collection.Indexes) > 0 {
			return fmt.Errorf("%s collections: %v: has no attributes:", stack.Name, colName)
		}
		if len(collection.Attributes) == 1 {
			return fmt.Errorf("%s collections: %v: requires 2 or more attributes", stack.Name, colName)
		}
		if len(collection.Attributes) > 0 && !hasAttribute("value", collection.Attributes) {
			return fmt.Errorf("%s collections: %v: requires a value: attribute", stack.Name, colName)
		}

		for attName, value := range collection.Attributes {
			if value != "string" {
				return fmt.Errorf("%s collections: %v: attributes: %v: %v is not supported, use string", stack.Name, colName, attName, value)
			}
		}

		// Ensure index names are valid
		for idxName, value := range collection.Indexes {
			if idxName != "unique" && idxName != "composite" {
				return fmt.Errorf("%s collections: %v: indexes: %v: is invalid, use unique: or composite:", stack.Name, colName, idxName)
			}
			// Ensure index is a defined attribute
			valueKind := reflect.ValueOf(value).Kind()

			if valueKind == reflect.String {
				if idxName == "composite" {
					return fmt.Errorf("%s collections: %v: indexes: %v: requires 2 values", stack.Name, colName, idxName)
				}
				if !hasAttribute(value, collection.Attributes) {
					return fmt.Errorf("%s collections: %v: indexes: %v: %v has no matching collection attribute", stack.Name, colName, idxName, value)
				}

			} else if valueKind == reflect.Slice {
				if idxName == "unique" {
					return fmt.Errorf("%s collections: %v: indexes: %v: does not support composite values %v", stack.Name, colName, idxName, value)
				}

				values := marshalInterfaceSlice(value)

				if len(values) != 2 {
					return fmt.Errorf("%s collections: %v: indexes: %v: requires 2 values %v", stack.Name, colName, idxName, value)
				}

				for _, v := range values {
					if !hasAttribute(v, collection.Attributes) {
						return fmt.Errorf("%s collections: %v: indexes: %v: %v has no matching collection attribute", stack.Name, colName, idxName, value)
					}
				}

			} else if value == nil {
				return fmt.Errorf("%s collections: %v: indexes: %v: not defined", stack.Name, colName, idxName)

			} else {
				return fmt.Errorf("%s collections: %v: indexes: %v: %T invalid", stack.Name, colName, idxName, value)
			}
		}
	}

	return nil
}
