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

var validOperators = map[string]bool{
	"==":         true,
	">":          true,
	"<":          true,
	">=":         true,
	"<=":         true,
	"startsWith": true,
}

func GetCollection(collection string) (string, error) {
	if collection == "" {
		return "", fmt.Errorf("provide non-blank collection")
	}

	return collection, nil
}

func GetKeyValue(key map[string]interface{}) (string, error) {
	// Get key
	if key == nil {
		return "", fmt.Errorf("provide non-nil key")
	}
	if len(key) == 0 {
		return "", fmt.Errorf("provide non-empty key")
	}

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
		if key[k] == "" {
			return "", fmt.Errorf("provide non-empty key")
		}
		returnKey += fmt.Sprintf("%v", key[k])
	}

	return returnKey, nil
}

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
