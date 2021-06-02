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
	"strings"

	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
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

var stack *utils.NitricStack

func Stack() utils.NitricStack {
	if stack == nil {
		nitricStack, err := utils.NewStackDefault()
		if err != nil {
			panic(fmt.Sprintf("Error loading Nitric stack definition: %v", err))
		}
		stack = nitricStack
	}
	return *stack
}

// Validate the collection name
func ValidateCollection(collection string) error {
	if collection == "" {
		return fmt.Errorf("provide non-blank collection")
	}
	if !Stack().HasCollection(collection) {
		return fmt.Errorf("%v collections: %v: not found", Stack().Name, collection)
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

	strEndCode := value[len(value)-1:]

	return strFrontCode + string(strEndCode[0]+1)
}

// Validate the provided query expressions
func ValidateExpressions(collection string, expressions []sdk.QueryExpression) error {
	if !Stack().HasCollection(collection) {
		return fmt.Errorf("%v collections: %v: not found", Stack().Name, collection)
	}

	if expressions == nil {
		return errors.New("provide non-nil query expressions")
	}

	attributes, err := Stack().CollectionAttributes(collection)
	if err != nil {
		return err
	}

	inequalityProperties := make(map[string]string)

	for _, exp := range expressions {
		if utils.IndexOf(attributes, exp.Operand) == -1 {
			attributes, _ := Stack().CollectionAttributes(collection)
			return fmt.Errorf("query expression '%v' operand not found in %v collections: %v: attributes: %v", exp.Operand, Stack().Name, collection, attributes)
		}
		if exp.Operand == "value" {
			return fmt.Errorf("query expression 'value' operand is not indexed: %v", exp)
		}
		if _, found := validOperators[exp.Operator]; !found {
			return fmt.Errorf("provide valid query expression operator [==, <, >, <=, >=, startsWith]: %v", exp.Operator)
		}
		if exp.Value == "" {
			return fmt.Errorf("provide non-blank query expression value: %v", exp)
		}

		// Ensure key expressions are valid
		keys, err := Stack().CollectionIndexes(collection)
		if err != nil {
			return err
		}
		if keys[0] == exp.Operand && exp.Operator != "==" {
			return fmt.Errorf("collection: '%v' key '%v' only supports '==' query operator: %v", collection, exp.Operand, exp)
		}

		if exp.Operator != "==" {
			inequalityProperties[exp.Operand] = exp.Operator
		}
	}

	// Firestore inequality compatability check
	if len(inequalityProperties) > 1 {
		msg := ""
		for prop, exp := range inequalityProperties {
			if msg != "" {
				msg += ", "
			}
			msg += prop + " " + exp
		}
		return fmt.Errorf("inequality expressions on multiple properties not supported with Firestore: [ %v ]", msg)
	}

	// DynamoDB range expression compatability check
	if err = hasRangeError(expressions); err != nil {
		return err
	}

	return nil
}

// Validate the provided key map
func ValidateKeyMap(collection string, key map[string]interface{}) error {
	if !Stack().HasCollection(collection) {
		return fmt.Errorf("%v collections: %v: not found", Stack().Name, collection)
	}

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

	// Validate names and values
	indexes, err := Stack().CollectionIndexes(collection)
	if err != nil {
		return err
	}

	for name, value := range key {
		if utils.IndexOf(indexes, name) == -1 {
			return fmt.Errorf("%v collections: %v: indexes: key '%s' not found", Stack().Name, collection, name)
		}

		if value == "" {
			return fmt.Errorf("provide non-blank key value")
		}
	}

	return nil
}

// QueryExpression sorting support with sort.Interface

type ExpsSort []sdk.QueryExpression

func (exps ExpsSort) Len() int {
	return len(exps)
}

// Sort by Operand then Operator then Value
func (exps ExpsSort) Less(i, j int) bool {

	operandCompare := strings.Compare(exps[i].Operand, exps[j].Operand)
	if operandCompare == 0 {

		// Reverse operator comparison for to support range expressions
		operatorCompare := strings.Compare(exps[j].Operator, exps[i].Operator)
		if operatorCompare == 0 {

			return strings.Compare(exps[i].Value, exps[j].Value) < 0

		} else {
			return operatorCompare < 0
		}

	} else {
		return operandCompare < 0
	}
}

func (exps ExpsSort) Swap(i, j int) {
	exps[i], exps[j] = exps[j], exps[i]
}

// DynamoDB only supports query range operands: >= AND <=
// For example: WHERE price >= 20.00 AND price <= 50.0
func hasRangeError(exps []sdk.QueryExpression) error {

	sortedExps := make([]sdk.QueryExpression, len(exps))
	copy(sortedExps, exps)

	sort.Sort(ExpsSort(sortedExps))

	for index, exp := range sortedExps {
		if index < (len(sortedExps) - 1) {
			nextExp := sortedExps[index+1]

			if exp.Operand == nextExp.Operand &&
				((exp.Operator == ">" && nextExp.Operator == "<") ||
					(exp.Operator == ">" && nextExp.Operator == "<=") ||
					(exp.Operator == ">=" && nextExp.Operator == "<")) {

				return fmt.Errorf("range expression not supported with DynamoDB (use operators >= and <=) : %v", exp)
			}
		}
	}

	return nil
}
