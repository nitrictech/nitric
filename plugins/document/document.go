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

package document

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
)

const ATTRIB_PK = "_pk"

const ATTRIB_SK = "_sk"

const ROOT_SK = "ROOT#"

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

func Stack() (*utils.NitricStack, error) {
	if stack == nil {
		nitricStack, err := utils.NewStackDefault()
		if err != nil {
			return nil, fmt.Errorf("error loading Nitric stack definition: %v", err)
		}
		stack = nitricStack
	}
	return stack, nil
}

// Validate the collection name and subcollection (optional)
func ValidateCollection(collection string, subcollection string) error {
	if collection == "" {
		return fmt.Errorf("provide non-blank collection")
	}
	stack, err := Stack()
	if err != nil {
		return err
	}
	if !stack.HasCollection(collection) {
		return fmt.Errorf("%v collections: %v: not found", stack.Name, collection)
	}
	// TODO: review subcollection validation
	// if subcollection != "" && !stack.HasSubCollection(collection, subcollection) {
	// 	return fmt.Errorf("%v collections: %v: sub-collection: %v: not found", stack.Name, collection, subcollection)
	// }
	return nil
}

func ValidateKeys(key sdk.Key, subKey *sdk.Key) error {
	stack, err := Stack()
	if err != nil {
		return err
	}

	if key.Collection == "" {
		return fmt.Errorf("provide non-blank key.Collection")
	}
	if !stack.HasCollection(key.Collection) {
		return fmt.Errorf("%v collections: %v: not found", stack.Name, key.Collection)
	}
	if key.Id == "" {
		return fmt.Errorf("provide non-blank key.Id")
	}
	if subKey != nil {
		if subKey.Collection == "" {
			return fmt.Errorf("provide non-blank subKey.Collection")
		}
		// TODO: review subcollection validate in the future
		// if !stack.HasSubCollection(key.Collection, subKey.Collection) {
		// 	return fmt.Errorf("%v collections: %v: sub-collection: %v: not found", stack.Name, key.Collection, subKey.Collection)
		// }
		if subKey.Id == "" {
			return fmt.Errorf("provide non-blank subKey.Id")
		}
	}
	return nil
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
func ValidateExpressions(expressions []sdk.QueryExpression) error {
	if expressions == nil {
		return errors.New("provide non-nil query expressions")
	}

	inequalityProperties := make(map[string]string)

	for _, exp := range expressions {
		if exp.Operand == "" {
			return fmt.Errorf("provide non-blank query expression operand: %v", exp)
		}

		if _, found := validOperators[exp.Operator]; !found {
			return fmt.Errorf("provide valid query expression operator [==, <, >, <=, >=, startsWith]: %v", exp.Operator)
		}
		if exp.Value == "" {
			return fmt.Errorf("provide non-blank query expression value: %v", exp)
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
	if err := hasRangeError(expressions); err != nil {
		return err
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
