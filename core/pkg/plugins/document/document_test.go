// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package document_test

import (
	"sort"

	"github.com/nitrictech/nitric/core/pkg/plugins/document"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Function Test Cases

var _ = Describe("Document Plugin", func() {
	When("ValidateKey", func() {
		When("Nil key", func() {
			It("should return error", func() {
				err := document.ValidateKey(nil)
				Expect(err.Error()).To(ContainSubstring("provide non-nil key"))
			})
		})
		When("Blank key.Collection", func() {
			It("should return error", func() {
				err := document.ValidateKey(&document.Key{})
				Expect(err.Error()).To(ContainSubstring("provide non-blank key.Id"))
			})
		})
		When("Blank key.Id", func() {
			It("should return error", func() {
				key := document.Key{
					Collection: &document.Collection{Name: "users"},
				}
				err := document.ValidateKey(&key)
				Expect(err.Error()).To(ContainSubstring("provide non-blank key.Id"))
			})
		})
		When("Blank key.Collection.Parent.Collection.Name", func() {
			It("should return error", func() {
				key := document.Key{
					Collection: &document.Collection{Name: "users", Parent: &document.Key{}},
					Id:         "123",
				}
				err := document.ValidateKey(&key)
				Expect(err.Error()).To(ContainSubstring("invalid parent for collection users, provide non-blank key.Id"))
			})
		})
		When("Blank key.Collection.Parent.Id", func() {
			It("should return error", func() {
				key := document.Key{
					Collection: &document.Collection{
						Name:   "orders",
						Parent: &document.Key{Collection: &document.Collection{Name: "customers"}},
					},
					Id: "123",
				}
				err := document.ValidateKey(&key)
				Expect(err.Error()).To(ContainSubstring("invalid parent for collection orders, provide non-blank key.Id"))
			})
		})
	})

	When("ValidateQueryCollection", func() {
		When("Nil key", func() {
			It("should return error", func() {
				err := document.ValidateQueryCollection(nil)
				Expect(err.Error()).To(ContainSubstring("provide non-nil collection"))
			})
		})
		When("Blank key.Collection", func() {
			It("should return error", func() {
				err := document.ValidateQueryCollection(&document.Collection{})
				Expect(err.Error()).To(ContainSubstring("provide non-blank collection.Name"))
			})
		})
		When("Blank key.Id", func() {
			It("should return nil", func() {
				coll := document.Collection{Name: "users"}
				err := document.ValidateQueryCollection(&coll)
				Expect(err).To(BeNil())
			})
		})
		When("Blank key.Collection.Parent.Collection.Name", func() {
			It("should return error", func() {
				coll := document.Collection{
					Name: "users",
					Parent: &document.Key{
						Id:         "test-key",
						Collection: &document.Collection{},
					},
				}
				err := document.ValidateQueryCollection(&coll)
				Expect(err.Error()).To(ContainSubstring("provide non-blank collection.Name"))
			})
		})
		When("Blank collection.Parent.Id", func() {
			It("should return nil", func() {
				coll := document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Id:         "test-key",
						Collection: &document.Collection{Name: "customers"},
					},
				}
				err := document.ValidateQueryCollection(&coll)
				Expect(err).To(BeNil())
			})
		})
	})

	When("GetValueEndCode", func() {
		It("should get next value", func() {
			endCode := document.GetEndRangeValue("Customer#")
			Expect(endCode).NotTo(BeNil())
			Expect(endCode).To(BeEquivalentTo("Customer$"))
		})
	})

	When("ExpsSort", func() {
		When("order is sorted", func() {
			It("Should not change order", func() {
				exps := []document.QueryExpression{
					{Operand: "A", Operator: "==", Value: "1"},
					{Operand: "B", Operator: "==", Value: "2"},
					{Operand: "C", Operator: "==", Value: "3"},
				}
				sort.Sort(document.ExpsSort(exps))
				Expect(exps[0].Operand).To(BeEquivalentTo("A"))
				Expect(exps[1].Operand).To(BeEquivalentTo("B"))
				Expect(exps[2].Operand).To(BeEquivalentTo("C"))
			})
		})
		When("not order not sorted", func() {
			It("Should not change order", func() {
				exps := []document.QueryExpression{
					{Operand: "C", Operator: "==", Value: "3"},
					{Operand: "A", Operator: "==", Value: "1"},
					{Operand: "B", Operator: "==", Value: "2"},
				}
				sort.Sort(document.ExpsSort(exps))
				Expect(exps[0].Operand).To(BeEquivalentTo("A"))
				Expect(exps[1].Operand).To(BeEquivalentTo("B"))
				Expect(exps[2].Operand).To(BeEquivalentTo("C"))
			})
		})
		When("not order not sorted", func() {
			It("Should not change order", func() {
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "==", Value: "3"},
					{Operand: "number", Operator: ">=", Value: "1"},
					{Operand: "number", Operator: "<=", Value: "2"},
				}
				sort.Sort(document.ExpsSort(exps))
				Expect(exps[0].Operator).To(BeEquivalentTo(">="))
				Expect(exps[1].Operator).To(BeEquivalentTo("=="))
				Expect(exps[2].Operator).To(BeEquivalentTo("<="))
			})
		})
	})

	When("ValidateExpression", func() {
		When("expression is valid", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "123"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("expressions empty", func() {
			It("should be valid", func() {
				err := document.ValidateExpressions([]document.QueryExpression{})
				Expect(err).To(BeNil())
			})
		})
		When("operand is nil", func() {
			It("should return error", func() {
				err := document.ValidateExpressions(nil)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operand not found", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "", Operator: "==", Value: "123"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operator is blank", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "pk", Operator: "", Value: "123"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("value is blank", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "pk", Operator: "==", Value: ""},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "pk", Operator: "=", Value: "123"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startWith", Value: "Order#"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("inequality query against multiple operations", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startsWith", Value: "Order#"},
					{Operand: "number", Operator: ">", Value: "1"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("inequality expressions on multiple properties are not supported:"))
			})
		})
		When("valid range filter expression", func() {
			It("expression is valid", func() {
				exps := []document.QueryExpression{
					{Operand: "number", Operator: ">=", Value: "1"},
					{Operand: "number", Operator: "<=", Value: "2"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("valid range filter expression in reverse order", func() {
			It("expression is valid", func() {
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "<=", Value: "1"},
					{Operand: "number", Operator: ">=", Value: "2"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("invalid valid range filter expression", func() {
			It("should return error", func() {
				exps := []document.QueryExpression{
					{Operand: "number", Operator: ">", Value: "1"},
					{Operand: "number", Operator: "<=", Value: "2"},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("range expression combination not supported (use operators >= and <=) :"))
			})
		})
	})
})
