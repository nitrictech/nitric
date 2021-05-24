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

package kv

import (
	"errors"
	"fmt"
	"sort"

	"github.com/nitric-dev/membrane/sdk"
)

// Map of valid expression operators
var validOperators = map[string]bool{
	"==":         true,
	">":          true,
	"<":          true,
	">=":         true,
	"<=":         true,
	"startsWith": true,
}

// Validate the collection name
func ValidateCollection(collection string) error {
	if collection == "" {
		return fmt.Errorf("provide non-blank collection")
	}
	return nil
}

func GetKeyMap(key string) (map[string]string, error) {
	if key == "" {
		return nil, fmt.Errorf("provide non-blank key")
	}

	return map[string]string{
		"key": key,
	}, nil
}

// Return a single key value, which appends multiple key values when the key map is > 0.
// For example: {"pk": "Customer#1000", "sk": "Order#200"} => "Customer#1000_Order#200"
func GetKeyValue(key map[string]interface{}) string {

	// Create sorted keys
	ks := make([]string, 0, len(key))
	for k := range key {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	returnKey := ""
	for i, k := range ks {
		if i > 0 {
			returnKey += "_"
		}
		returnKey += fmt.Sprintf("%v", key[k])
	}

	return returnKey
}

// Return a sorted list of key map values
func GetKeyValues(key map[string]interface{}) []string {

	// Create sorted keys
	ks := make([]string, 0, len(key))
	for k := range key {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	kv := make([]string, 0, len(key))

	for _, k := range ks {
		kv = append(kv, fmt.Sprintf("%v", key[k]))
	}

	return kv
}

// Get end range value to implement "startsWith" expression operator using where clause.
// For example with sdk.Expression("pk", "startsWith", "Customer#") this translates to:
// WHERE pk >= {startRangeValue} AND pk < {endRangeValue}
// WHERE pk >= "Customer#" AND pk < "Customer!"
func GetEndRangeValue(value string) string {
	strFrontCode := value[:len(value)-1]

	strEndCode := value[len(value)-1 : len(value)]

	return strFrontCode + string(strEndCode[0]+1)
}

// Validate the provided query expressions
func ValidateExpressions(expressions []sdk.QueryExpression) error {
	if expressions == nil {
		return errors.New("provide non-nil expressions")
	}

	for _, exp := range expressions {
		if exp.Operand == "" {
			return fmt.Errorf("provide non-blank expression operand: %v", exp)
		}
		if _, found := validOperators[exp.Operator]; !found {
			return fmt.Errorf("provide valid expression operator [==, <, >, <=, >=, startsWith]: %v", exp.Operator)
		}
		if exp.Value == "" {
			return fmt.Errorf("provide non-blank expression value: %v", exp)
		}
	}
	if len(expressions) > 0 && expressions[0].Operator != "==" {
		return fmt.Errorf("provide identity operator (==) for primary key expressions: %v", expressions[0])
	}

	return nil
}

// Validate the provided key map
func ValidateKeyMap(key map[string]interface{}) error {
	// Get key
	if key == nil {
		return fmt.Errorf("provide non-nil key")
	}
	if len(key) == 0 {
		return fmt.Errorf("provide non-empty key")
	}
	if len(key) > 2 {
		return fmt.Errorf("provide key with 1 or 2 items")
	}

	for _, v := range key {
		if v == "" {
			return fmt.Errorf("provide non-blank key value")
		}
	}

	return nil
}
