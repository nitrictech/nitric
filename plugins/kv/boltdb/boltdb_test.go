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

package boltdb_service_test

import (
	"os"

	kv_plugin "github.com/nitric-dev/membrane/plugins/kv/boltdb"
	data "github.com/nitric-dev/membrane/plugins/kv/test"
	"github.com/nitric-dev/membrane/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KV", func() {
	kvPlugin, err := kv_plugin.New()
	if err != nil {
		panic(err)
	}

	AfterSuite(func() {
		err = os.RemoveAll(kv_plugin.DEFAULT_DIR)
		if err == nil {
			os.Remove(kv_plugin.DEFAULT_DIR)
			os.Remove("nitric/")
		}
	})

	Context("Put", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("", data.UserKey1, data.UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("users", nil, data.UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil item map", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("users", data.UserKey1, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid New Put", func() {
			It("Should store new item successfully", func() {
				err := kvPlugin.Put("users", data.UserKey1, data.UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem1["email"]))
			})
		})
		When("Valid Update Put", func() {
			It("Should update existing item successfully", func() {
				err := kvPlugin.Put("users", data.UserKey1, data.UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem1["email"]))

				err = kvPlugin.Put("users", data.UserKey1, data.UserItem2)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err = kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem2["email"]))
			})
		})
		When("Valid Compound Key Put", func() {
			It("Should store item successfully", func() {
				err := kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("application", data.OrderKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(data.OrderItem1))
			})
		})
	})

	Context("Get", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				_, err := kvPlugin.Get("", data.UserKey1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key", func() {
			It("Should return error", func() {
				_, err := kvPlugin.Get("users", nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Get", func() {
			It("Should get item successfully", func() {
				kvPlugin.Put("users", data.UserKey1, data.UserItem1)

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(data.UserItem1["email"]))
			})
		})
		When("Valid Compound Key Get", func() {
			It("Should store item successfully", func() {
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)

				doc, err := kvPlugin.Get("application", data.OrderKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(data.OrderItem1))
			})
		})
	})

	Context("Delete", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				err := kvPlugin.Delete("", data.UserKey1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key", func() {
			It("Should return error", func() {
				err := kvPlugin.Delete("collection", nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Delete", func() {
			It("Should delete item successfully", func() {
				kvPlugin.Put("users", data.UserKey1, data.UserItem1)

				err := kvPlugin.Delete("users", data.UserKey1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("users", data.UserKey1)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Compound Key Delete", func() {
			It("Should delete item successfully", func() {
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)

				err := kvPlugin.Delete("application", data.OrderKey1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("application", data.OrderKey1)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("Query", func() {
		When("Blank collection argument", func() {
			It("Should return an error", func() {
				vals, err := kvPlugin.Query("", nil, 0)
				Expect(vals).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil key argument", func() {
			It("Should return an error", func() {
				vals, err := kvPlugin.Query("users", nil, 0)
				Expect(vals).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Empty database collection", func() {
			It("Should return empty list", func() {
				vals, err := kvPlugin.Query("users", []sdk.QueryExpression{}, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(0))
			})
		})
		When("Filter users collection", func() {
			It("Should return 1 item", func() {
				kvPlugin.Put("users", data.UserKey1, data.UserItem1)
				kvPlugin.Put("users", data.UserKey2, data.UserItem2)
				kvPlugin.Put("users", data.UserKey3, data.UserItem3)
				exps := []sdk.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
					{Operand: "age", Operator: ">", Value: "40"},
				}
				vals, err := kvPlugin.Query("users", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(1))
				Expect(vals[0]["email"]).To(BeEquivalentTo(data.UserItem3["email"]))
			})
		})
		When("Empty query", func() {
			It("Should return all items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				vals, err := kvPlugin.Query("application", []sdk.QueryExpression{}, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(5))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[3]).To(BeEquivalentTo(data.OrderItem3))
				Expect(vals[4]).To(BeEquivalentTo(data.ProductItem))
			})
		})
		When("Empty limit query", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				vals, err := kvPlugin.Query("application", []sdk.QueryExpression{}, 3)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem2))
			})
		})
		When("PK and SK equality query", func() {
			It("Should return specified item", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "==", Value: "Customer#1000"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(1))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
			})
		})
		When("PK equality query", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(4))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[3]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality limit query", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
				}
				vals, err := kvPlugin.Query("application", exps, 3)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem2))
			})
		})
		When("PK equality and SK startsWith", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "startsWith", Value: "Order#"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))
				Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality and SK >", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: ">", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(2))
				Expect(vals[0]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality and SK >=", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: ">=", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(3))
				Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem2))
				Expect(vals[2]).To(BeEquivalentTo(data.OrderItem3))
			})
		})
		When("PK equality and SK <", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "<", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(1))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
			})
		})
		When("PK equality and SK <=", func() {
			It("Should return specified items", func() {
				kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
				kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
				kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
				kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
				kvPlugin.Put("application", data.ProductKey, data.ProductItem)

				exps := []sdk.QueryExpression{
					{Operand: "pk", Operator: "==", Value: "Customer#1000"},
					{Operand: "sk", Operator: "<=", Value: "Order#501"},
				}
				vals, err := kvPlugin.Query("application", exps, 0)
				Expect(vals).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(vals).To(HaveLen(2))
				Expect(vals[0]).To(BeEquivalentTo(data.CustomerItem))
				Expect(vals[1]).To(BeEquivalentTo(data.OrderItem1))
			})
		})
		// Firestore: cant support multiple property inequality operators
		// When("PK equality and SK startsWith and filter", func() {
		// 	It("Should return specified items", func() {
		// 		kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
		// 		kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
		// 		kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
		// 		kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
		// 		kvPlugin.Put("application", data.ProductKey, data.ProductItem)

		// 		exps := []sdk.QueryExpression{
		// 			{Operand: "pk", Operator: "==", Value: "Customer#1000"},
		// 			{Operand: "sk", Operator: "startsWith", Value: "Order#"},
		// 			{Operand: "number", Operator: ">", Value: "1"},
		// 			{Operand: "price", Operator: "<", Value: "20"},
		// 		}
		// 		vals, err := kvPlugin.Query("application", exps, 0)
		// 		Expect(vals).ToNot(BeNil())
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(vals).To(HaveLen(1))
		// 		Expect(vals[0]).To(BeEquivalentTo(data.OrderItem2))
		// 	})
		// })
		// When("PK equality and SK startsWith and between filter", func() {
		// 	It("Should return specified items", func() {
		// 		kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
		// 		kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
		// 		kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
		// 		kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
		// 		kvPlugin.Put("application", data.ProductKey, data.ProductItem)

		// 		exps := []sdk.QueryExpression{
		// 			{Operand: "pk", Operator: "==", Value: "Customer#1000"},
		// 			{Operand: "sk", Operator: "startsWith", Value: "Order#"},
		// 			{Operand: "number", Operator: ">=", Value: "0"},
		// 			{Operand: "number", Operator: "<=", Value: "1"},
		// 		}
		// 		vals, err := kvPlugin.Query("application", exps, 0)
		// 		Expect(vals).ToNot(BeNil())
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(vals).To(HaveLen(1))
		// 		Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
		// 	})
		// })
		// When("PK equality and SK startsWith and between filters with reversed order", func() {
		// 	It("Should return specified items", func() {
		// 		kvPlugin.Put("application", data.CustomerKey, data.CustomerItem)
		// 		kvPlugin.Put("application", data.OrderKey1, data.OrderItem1)
		// 		kvPlugin.Put("application", data.OrderKey2, data.OrderItem2)
		// 		kvPlugin.Put("application", data.OrderKey3, data.OrderItem3)
		// 		kvPlugin.Put("application", data.ProductKey, data.ProductItem)

		// 		exps := []sdk.QueryExpression{
		// 			{Operand: "pk", Operator: "==", Value: "Customer#1000"},
		// 			{Operand: "sk", Operator: "startsWith", Value: "Order#"},
		// 			{Operand: "number", Operator: "<=", Value: "1"},
		// 			{Operand: "number", Operator: ">=", Value: "0"},
		// 		}
		// 		vals, err := kvPlugin.Query("application", exps, 0)
		// 		Expect(vals).ToNot(BeNil())
		// 		Expect(err).ShouldNot(HaveOccurred())
		// 		Expect(vals).To(HaveLen(1))
		// 		Expect(vals[0]).To(BeEquivalentTo(data.OrderItem1))
		// 	})
		// })
	})

})
