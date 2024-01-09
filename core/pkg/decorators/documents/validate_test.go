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

	document "github.com/nitrictech/nitric/core/pkg/decorators/documents"
	documentpb "github.com/nitrictech/nitric/core/pkg/proto/documents/v1"

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
				err := document.ValidateKey(&documentpb.Key{})
				Expect(err.Error()).To(ContainSubstring("provide non-blank key.Id"))
			})
		})
		When("Blank key.Id", func() {
			It("should return error", func() {
				key := &documentpb.Key{
					Collection: &documentpb.Collection{Name: "users"},
				}
				err := document.ValidateKey(key)
				Expect(err.Error()).To(ContainSubstring("provide non-blank key.Id"))
			})
		})
		When("Blank key.Collection.Parent.Collection.Name", func() {
			It("should return error", func() {
				key := &documentpb.Key{
					Collection: &documentpb.Collection{Name: "users", Parent: &documentpb.Key{}},
					Id:         "123",
				}
				err := document.ValidateKey(key)
				Expect(err.Error()).To(ContainSubstring("invalid parent for collection users, provide non-blank key.Id"))
			})
		})
		When("Blank key.Collection.Parent.Id", func() {
			It("should return error", func() {
				key := &documentpb.Key{
					Collection: &documentpb.Collection{
						Name:   "orders",
						Parent: &documentpb.Key{Collection: &documentpb.Collection{Name: "customers"}},
					},
					Id: "123",
				}
				err := document.ValidateKey(key)
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
				err := document.ValidateQueryCollection(&documentpb.Collection{})
				Expect(err.Error()).To(ContainSubstring("provide non-blank collection.Name"))
			})
		})
		When("Blank key.Id", func() {
			It("should return nil", func() {
				coll := &documentpb.Collection{Name: "users"}
				err := document.ValidateQueryCollection(coll)
				Expect(err).To(BeNil())
			})
		})
		When("Blank key.Collection.Parent.Collection.Name", func() {
			It("should return error", func() {
				coll := &documentpb.Collection{
					Name: "users",
					Parent: &documentpb.Key{
						Id:         "test-key",
						Collection: &documentpb.Collection{},
					},
				}
				err := document.ValidateQueryCollection(coll)
				Expect(err.Error()).To(ContainSubstring("provide non-blank collection.Name"))
			})
		})
		When("Blank collection.Parent.Id", func() {
			It("should return nil", func() {
				coll := &documentpb.Collection{
					Name: "orders",
					Parent: &documentpb.Key{
						Id:         "test-key",
						Collection: &documentpb.Collection{Name: "customers"},
					},
				}
				err := document.ValidateQueryCollection(coll)
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
				exps := []*documentpb.Expression{
					{Operand: "A", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 1},
					}},
					{Operand: "B", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 2},
					}},
					{Operand: "C", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 3},
					}},
				}
				sort.Sort(document.ExpsSort(exps))
				Expect(exps[0].Operand).To(BeEquivalentTo("A"))
				Expect(exps[1].Operand).To(BeEquivalentTo("B"))
				Expect(exps[2].Operand).To(BeEquivalentTo("C"))
			})
		})
		When("not order not sorted", func() {
			It("Should not change order", func() {
				exps := []*documentpb.Expression{
					{Operand: "A", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 1},
					}},
					{Operand: "B", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 2},
					}},
					{Operand: "C", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 3},
					}},
				}
				sort.Sort(document.ExpsSort(exps))
				Expect(exps[0].Operand).To(BeEquivalentTo("A"))
				Expect(exps[1].Operand).To(BeEquivalentTo("B"))
				Expect(exps[2].Operand).To(BeEquivalentTo("C"))
			})
		})
		When("not order not sorted", func() {
			It("Should not change order", func() {
				exps := []*documentpb.Expression{
					{Operand: "number", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 1},
					}},
					{Operand: "number", Operator: ">=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 1},
					}},
					{Operand: "number", Operator: "<=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 2},
					}},
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
				exps := []*documentpb.Expression{
					{Operand: "ok", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "123"},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("expressions empty", func() {
			It("should be valid", func() {
				err := document.ValidateExpressions([]*documentpb.Expression{})
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
				exps := []*documentpb.Expression{
					{Operand: "", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "123"},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operator is blank", func() {
			It("should return error", func() {
				exps := []*documentpb.Expression{
					{Operand: "pk", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "123"},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("value is blank", func() {
			It("should return error", func() {
				exps := []*documentpb.Expression{
					{Operand: "", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: ""},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []*documentpb.Expression{
					{Operand: "", Operator: "=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "123"},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []*documentpb.Expression{
					{Operand: "pk", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "Customer#1000"},
					}},
					{Operand: "sk", Operator: "startWith", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "Customer#1000"},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("inequality query against multiple operations", func() {
			It("should return error", func() {
				exps := []*documentpb.Expression{
					{Operand: "pk", Operator: "==", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "Customer#1000"},
					}},
					{Operand: "sk", Operator: "startsWith", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "Order#"},
					}},
					{Operand: "number", Operator: "startsWith", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_StringValue{StringValue: "1"},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("inequality expressions on multiple properties are not supported:"))
			})
		})
		When("valid range filter expression", func() {
			It("expression is valid", func() {
				exps := []*documentpb.Expression{
					{Operand: "number", Operator: ">=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 1},
					}},
					{Operand: "number", Operator: "<=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 2},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("valid range filter expression in reverse order", func() {
			It("expression is valid", func() {
				exps := []*documentpb.Expression{
					{Operand: "number", Operator: "<=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 1},
					}},
					{Operand: "number", Operator: ">=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 2},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("invalid valid range filter expression", func() {
			It("should return error", func() {
				exps := []*documentpb.Expression{
					{Operand: "number", Operator: ">", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 1},
					}},
					{Operand: "number", Operator: "<=", Value: &documentpb.ExpressionValue{
						Kind: &documentpb.ExpressionValue_IntValue{IntValue: 2},
					}},
				}
				err := document.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("range expression combination not supported (use operators >= and <=) :"))
			})
		})
	})
})
