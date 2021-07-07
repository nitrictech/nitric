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
package document_test

import (
	"errors"
	"os"
	"sort"

	doc "github.com/nitric-dev/membrane/plugins/document"
	"github.com/nitric-dev/membrane/sdk"
	"github.com/nitric-dev/membrane/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Function Test Cases

var _ = Describe("Document Plugin", func() {

	os.Setenv(utils.NITRIC_HOME, "test/")
	os.Setenv(utils.NITRIC_YAML, "nitric.yaml")

	When("ValidateCollection", func() {
		When("Blank collection", func() {
			It("should return error", func() {
				err := doc.ValidateCollection("", "")
				Expect(err).To(BeEquivalentTo(errors.New("provide non-blank collection")))
			})
		})
		When("Missing collection", func() {
			It("should return error", func() {
				err := doc.ValidateCollection("uknown", "")
				Expect(err.Error()).To(BeEquivalentTo("nitric-website collections: uknown: not found"))
			})
		})
		When("Valid collection", func() {
			It("should return nil", func() {
				err := doc.ValidateCollection("users", "")
				Expect(err).To(BeNil())
			})
		})
		When("Missing subcollection", func() {
			It("should return nil", func() {
				err := doc.ValidateCollection("customers", "unknown")
				Expect(err).To(BeNil())
				// TODO: review subcollection validation in future
				// Expect(err.Error()).To(BeEquivalentTo("nitric-website collections: customers: sub-collection: unknown: not found"))
			})
		})
		// TODO: review subcollection validation in future
		// When("Valid subcollection", func() {
		// 	It("should return nil", func() {
		// 		err := doc.ValidateCollection("customers", "addresses")
		// 		Expect(err).To(BeNil())
		// 	})
		// })
	})

	When("ValidateKeys", func() {
		When("Blank key.Collection", func() {
			It("should return error", func() {
				err := doc.ValidateKeys(sdk.Key{}, nil)
				Expect(err).To(BeEquivalentTo(errors.New("provide non-blank key.Collection")))
			})
		})
		When("Blank key.Id", func() {
			It("should return error", func() {
				err := doc.ValidateKeys(sdk.Key{Collection: "users"}, nil)
				Expect(err).To(BeEquivalentTo(errors.New("provide non-blank key.Id")))
			})
		})
		When("Unknown key.Collection", func() {
			It("should return error", func() {
				err := doc.ValidateKeys(sdk.Key{Collection: "unknown", Id: "123"}, nil)
				Expect(err).To(BeEquivalentTo(errors.New("nitric-website collections: unknown: not found")))
			})
		})
		// TODO: review subcollection validate in the future
		// When("Unknown subKey.Collection", func() {
		// 	It("should return error", func() {
		// 		err := doc.ValidateKeys(sdk.Key{Collection: "users", Id: "123"}, &sdk.Key{Collection: "unknown"})
		// 		Expect(err).To(BeEquivalentTo(errors.New("nitric-website collections: users: sub-collection: unknown: not found")))
		// 	})
		// })
	})

	When("GetValueEndCode", func() {
		It("should get next value", func() {
			endCode := doc.GetEndRangeValue("Customer#")
			Expect(endCode).NotTo(BeNil())
			Expect(endCode).To(BeEquivalentTo("Customer$"))
		})
	})

	When("ExpsSort", func() {
		When("order is sorted", func() {
			It("Should not change order", func() {
				exps := []sdk.QueryExpression{
					{Operand: "A", Operator: "==", Value: "1"},
					{Operand: "B", Operator: "==", Value: "2"},
					{Operand: "C", Operator: "==", Value: "3"},
				}
				sort.Sort(doc.ExpsSort(exps))
				Expect(exps[0].Operand).To(BeEquivalentTo("A"))
				Expect(exps[1].Operand).To(BeEquivalentTo("B"))
				Expect(exps[2].Operand).To(BeEquivalentTo("C"))
			})
		})
		When("not order not sorted", func() {
			It("Should not change order", func() {
				exps := []sdk.QueryExpression{
					{Operand: "C", Operator: "==", Value: "3"},
					{Operand: "A", Operator: "==", Value: "1"},
					{Operand: "B", Operator: "==", Value: "2"},
				}
				sort.Sort(doc.ExpsSort(exps))
				Expect(exps[0].Operand).To(BeEquivalentTo("A"))
				Expect(exps[1].Operand).To(BeEquivalentTo("B"))
				Expect(exps[2].Operand).To(BeEquivalentTo("C"))
			})
		})
		When("not order not sorted", func() {
			It("Should not change order", func() {
				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "==", Value: "3"},
					{Operand: "number", Operator: ">=", Value: "1"},
					{Operand: "number", Operator: "<=", Value: "2"},
				}
				sort.Sort(doc.ExpsSort(exps))
				Expect(exps[0].Operator).To(BeEquivalentTo(">="))
				Expect(exps[1].Operator).To(BeEquivalentTo("=="))
				Expect(exps[2].Operator).To(BeEquivalentTo("<="))
			})
		})
	})

	When("ValidateExpression", func() {
		When("expression is valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "123"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("expressions empty", func() {
			It("should be valid", func() {
				err := doc.ValidateExpressions([]sdk.QueryExpression{})
				Expect(err).To(BeNil())
			})
		})
		When("operand is nil", func() {
			It("should return error", func() {
				err := doc.ValidateExpressions(nil)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operand not found", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "", Operator: "==", Value: "123"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operator is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "", Value: "123"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("value is blank", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: ""},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "=", Value: "123"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("operation is not valid", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startWith", Value: "Order#"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
			})
		})
		When("inequality query against muliple operations", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startsWith", Value: "Order#"},
					{Operand: "number", Operator: ">", Value: "1"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("inequality expressions on multiple properties not supported with Firestore:"))
			})
		})
		When("valid range filter expression", func() {
			It("expression is valid", func() {
				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: ">=", Value: "1"},
					{Operand: "number", Operator: "<=", Value: "2"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("valid range filter expression in reverse order", func() {
			It("expression is valid", func() {
				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "<=", Value: "1"},
					{Operand: "number", Operator: ">=", Value: "2"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).To(BeNil())
			})
		})
		When("invalid valid range filter expression", func() {
			It("should return error", func() {
				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: ">", Value: "1"},
					{Operand: "number", Operator: "<=", Value: "2"},
				}
				err := doc.ValidateExpressions(exps)
				Expect(err).ToNot(BeNil())
				Expect(err.Error()).To(HavePrefix("range expression not supported with DynamoDB"))
			})
		})
	})
})
