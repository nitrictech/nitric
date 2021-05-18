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
package kv_test

import (
	"errors"

	"github.com/nitric-dev/membrane/plugins/kv"
	"github.com/nitric-dev/membrane/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KeyValue Plugin", func() {

	When("GetCollection", func() {
		When("Blank collection", func() {
			It("should return error", func() {
				col, err := kv.GetCollection("")
				Expect(err).To(BeEquivalentTo(errors.New("provide non-blank collection")))
				Expect(col).To(BeEquivalentTo(""))
			})
		})
		When("Valid collection", func() {
			It("should return collection", func() {
				col, err := kv.GetCollection("collection")
				Expect(err).To(BeNil())
				Expect(col).To(BeEquivalentTo("collection"))
			})
		})
	})

	When("GetKeyValue", func() {
		When("Nil key", func() {
			It("should return error", func() {
				key, err := kv.GetKeyValue(nil)
				Expect(err).To(BeEquivalentTo(errors.New("provide non-nil key")))
				Expect(key).To(BeEquivalentTo(""))
			})
		})
		When("Empty key", func() {
			It("should return error", func() {
				keyMap := map[string]interface{}{}
				key, err := kv.GetKeyValue(keyMap)
				Expect(err).To(BeEquivalentTo(errors.New("provide non-empty key")))
				Expect(key).To(BeEquivalentTo(""))
			})
		})
		When("Single key", func() {
			It("should return single key value", func() {
				keyMap := map[string]interface{}{
					"key": "user@server.com",
				}
				key, err := kv.GetKeyValue(keyMap)
				Expect(err).To(BeNil())
				Expect(key).To(BeEquivalentTo("user@server.com"))
			})
		})
		When("Multi key", func() {
			It("should return appended key values", func() {
				keyMap := map[string]interface{}{
					"pk": "Customer#123",
					"sk": "Order#456",
				}
				key, err := kv.GetKeyValue(keyMap)
				Expect(err).To(BeNil())
				Expect(key).To(BeEquivalentTo("Customer#123_Order#456"))
			})
		})
	})

	When("ValidateExpression", func() {
		When("expression is valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "==", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("expressions empty", func() {
			It("should be valid", func() {
				err := kv.ValidateExpressions([]sdk.QueryExpression{})
				Expect(err).To(BeNil())
			})
		})
		When("operand is nil", func() {
			It("should return error", func() {
				err := kv.ValidateExpressions(nil)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operand is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "", Operator: "==", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operator is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("value is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "==", Value: ""},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "Pk", Operator: "=", Value: "123"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startWith", Value: "Order#"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("primary key operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: ">", Value: "Customer#1000"},
				}
				err := kv.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
