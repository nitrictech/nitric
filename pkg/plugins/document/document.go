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

	"github.com/nitric-dev/membrane/pkg/sdk"
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

// ValidateSubCollectionDepth - returns an error if the provided collection exceeds the maximum supported
// depth for a sub-collection.
func ValidateSubCollectionDepth(collection *sdk.Collection) error {
	coll := collection
	depth := 0
	for coll.Parent != nil {
		depth += 1
		coll = coll.Parent.Collection
	}
	if depth > sdk.MaxSubCollectionDepth {
		return fmt.Errorf(
			"sub-collections only supported to a depth of %d, found depth of %d for collection %s",
			sdk.MaxSubCollectionDepth,
			depth,
			collection.Name,
		)
	}
	return nil
}

// ValidateKey - validates a document key, used for operations on a single document e.g. Get, Set, Delete
func ValidateKey(key *sdk.Key) error {
	if key == nil {
		return fmt.Errorf("provide non-nil key")
	}
	if key.Id == "" {
		return fmt.Errorf("provide non-blank key.Id")
	}
	if key.Collection == nil {
		return fmt.Errorf("provide non-nil key.Collection")
	} else {
		if err := ValidateCollection(key.Collection); err != nil {
			return fmt.Errorf("invalid collection for document key %s, %v", key.Id, err)
		}
	}
	return nil
}

// ValidateCollection - validates a collection key, used for operations on a single document/collection e.g. Get, Set, Delete
func ValidateCollection(collection *sdk.Collection) error {
	if collection == nil {
		return fmt.Errorf("provide non-nil collection")
	}
	if collection.Name == "" {
		return fmt.Errorf("provide non-blank collection.Name")
	}
	if collection.Parent != nil {
		if err := ValidateKey(collection.Parent); err != nil {
			return fmt.Errorf("invalid parent for collection %s, %v", collection.Name, err)
		}
	}

	return ValidateSubCollectionDepth(collection)
}

// ValidateQueryKey - Validates a key used for query operations.
// unique from ValidateKey in that it permits blank key.Id values for wildcard query scenarios.
// e.g. querying values in a sub-collection for all documents in the parent collection.
func ValidateQueryKey(key *sdk.Key) error {
	if key == nil {
		return fmt.Errorf("provide non-nil key")
	}
	if key.Collection == nil {
		return fmt.Errorf("provide non-nil key.Collection")
	} else {
		if err := ValidateQueryCollection(key.Collection); err != nil {
			return fmt.Errorf("invalid collection for document key %s, %v", key.Id, err)
		}
	}
	return nil
}

// ValidateQueryCollection - Validates a collection used for query operations.
// unique from ValidateCollection in that it calls ValidateQueryKey for the collection.Key
func ValidateQueryCollection(collection *sdk.Collection) error {
	if collection == nil {
		return fmt.Errorf("provide non-nil collection")
	}
	if collection.Name == "" {
		return fmt.Errorf("provide non-blank collection.Name")
	}
	if collection.Parent != nil {
		if err := ValidateQueryKey(collection.Parent); err != nil {
			return fmt.Errorf("invalid parent for collection %s, %v", collection.Name, err)
		}
	}
	return ValidateSubCollectionDepth(collection)
}

// GetEndRangeValue - Get end range value to implement "startsWith" expression operator using where clause.
// For example with sdk.Expression("pk", "startsWith", "Customer#") this translates to:
// WHERE pk >= {startRangeValue} AND pk < {endRangeValue}
// WHERE pk >= "Customer#" AND pk < "Customer!"
func GetEndRangeValue(value string) string {
	strFrontCode := value[:len(value)-1]

	strEndCode := value[len(value)-1:]

	return strFrontCode + string(strEndCode[0]+1)
}

// ValidateExpressions - Validate the provided query expressions
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
		// Firestore does not support inequality expressions on multiple properties.
		// Firestore requires composite key to be created at deployment time.
		return fmt.Errorf("inequality expressions on multiple properties are not supported: [ %v ]", msg)
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

// Less - Sort by Operand then Operator then Value
func (exps ExpsSort) Less(i, j int) bool {

	operandCompare := strings.Compare(exps[i].Operand, exps[j].Operand)
	if operandCompare == 0 {

		// Reverse operator comparison for to support range expressions
		operatorCompare := strings.Compare(exps[j].Operator, exps[i].Operator)
		if operatorCompare == 0 {
			iVal := fmt.Sprintf("%v", exps[i].Value)
			jVal := fmt.Sprintf("%v", exps[2].Value)

			return strings.Compare(iVal, jVal) < 0

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

				// Range expression combination not supported with DynamoDB, must use >= and <= which maps to DynamoDB BETWEEN
				return fmt.Errorf("range expression combination not supported (use operators >= and <=) : %v", exp)
			}
		}
	}

	return nil
}
