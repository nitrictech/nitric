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

package document_suite

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/nitrictech/nitric/pkg/plugins/document"
)

func QueryTests(docPlugin document.DocumentService) {
	Context("Query", func() {
		When("Invalid - blank key.Collection.Name", func() {
			It("Should return an error", func() {
				result, err := docPlugin.Query(context.TODO(), &document.Collection{}, []document.QueryExpression{}, 0, nil)
				Expect(result).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Invalid - nil expressions argument", func() {
			It("Should return an error", func() {
				result, err := docPlugin.Query(context.TODO(), &document.Collection{Name: "users"}, nil, 0, nil)
				Expect(result).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Empty database", func() {
			It("Should return empty list", func() {
				result, err := docPlugin.Query(context.TODO(), &document.Collection{Name: "users"}, []document.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))
				Expect(result.PagingToken).To(BeNil())
			})
		})
		When("key: {users}, subcol: '', exp: []", func() {
			It("Should return all users", func() {
				LoadUsersData(docPlugin)
				LoadCustomersData(docPlugin)

				result, err := docPlugin.Query(context.TODO(), &document.Collection{Name: "users"}, []document.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(3))

				for _, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					Expect(d.Key.Collection.Name).To(Equal("users"))
					Expect(d.Key.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Parent).To(BeNil())
				}
			})
		})
		When("key: {customers, nil}, subcol: '', exp: []", func() {
			It("Should return 2 items", func() {
				LoadCustomersData(docPlugin)

				result, err := docPlugin.Query(context.TODO(), &CustomersColl, []document.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(2))
				Expect(result.Documents[0].Content["email"]).To(BeEquivalentTo(Customer1.Content["email"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Key))
				Expect(result.Documents[1].Content["email"]).To(BeEquivalentTo(Customer2.Content["email"]))
				Expect(*result.Documents[1].Key).To(BeEquivalentTo(Customer2.Key))
				Expect(result.PagingToken).To(BeNil())
			})
		})
		When("key: {customers, nil}, subcol: '', exp: [country == US]", func() {
			It("Should return 1 item", func() {
				LoadCustomersData(docPlugin)

				exps := []document.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
				}
				result, err := docPlugin.Query(context.TODO(), &CustomersColl, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(1))
				Expect(result.Documents[0].Content["email"]).To(BeEquivalentTo(Customer2.Content["email"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer2.Key))
				Expect(result.PagingToken).To(BeNil())
			})
		})
		When("key: {customers, nil}, subcol: '', exp: [country == US, age > 40]", func() {
			It("Should return 0 item", func() {
				LoadCustomersData(docPlugin)

				exps := []document.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
					{Operand: "age", Operator: ">", Value: "40"},
				}
				result, err := docPlugin.Query(context.TODO(), &CustomersColl, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))
			})
		})
		When("key: {customers, key1}, subcol: orders", func() {
			It("Should return 3 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}
				result, err := docPlugin.Query(context.TODO(), &coll, []document.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(3))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(result.Documents[1].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[1].Content["testName"]))
				Expect(*result.Documents[1].Key).To(BeEquivalentTo(Customer1.Orders[1].Key))
				Expect(*result.Documents[1].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(result.Documents[2].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(*result.Documents[2].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*result.Documents[2].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number == 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &CustomersKey,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "==", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(2))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(result.Documents[1].Content["testName"]).To(BeEquivalentTo(Customer2.Orders[0].Content["testName"]))
				Expect(*result.Documents[1].Key).To(BeEquivalentTo(Customer2.Orders[0].Key))
				Expect(*result.Documents[1].Key.Collection.Parent).To(BeEquivalentTo(Customer2.Key))
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number == 1]", func() {
			It("Should return 1 order", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "==", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(1))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number > 1]", func() {
			It("Should return 3 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &CustomersKey,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: ">", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(3))

				for _, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					Expect(d.Key.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Name).To(Equal("orders"))
					Expect(d.Key.Collection.Parent).ToNot(BeNil())
					Expect(d.Key.Collection.Parent.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Parent.Collection.Name).To(Equal("customers"))
				}
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number > 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: ">", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(2))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[1].Content["testName"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[1].Key))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(result.Documents[1].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(*result.Documents[1].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*result.Documents[1].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number < 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &CustomersKey,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "<", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))

				for _, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					Expect(d.Key.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Name).To(Equal("orders"))
					Expect(d.Key.Collection.Parent).ToNot(BeNil())
					Expect(d.Key.Collection.Parent.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Parent.Collection.Name).To(Equal("customers"))
				}
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number < 1]", func() {
			It("Should return 0 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "<", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number >= 1]", func() {
			It("Should return 5 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &CustomersKey,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: ">=", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(5))

				for _, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					Expect(d.Key.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Name).To(Equal("orders"))
					Expect(d.Key.Collection.Parent).ToNot(BeNil())
					Expect(d.Key.Collection.Parent.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Parent.Collection.Name).To(Equal("customers"))
				}
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number >= 1]", func() {
			It("Should return 3 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: ">=", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(3))

				for _, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					Expect(d.Key.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Name).To(Equal("orders"))
					Expect(d.Key.Collection.Parent).ToNot(BeNil())
					Expect(d.Key.Collection.Parent.Id).ToNot(Equal(""))
					Expect(d.Key.Collection.Parent.Collection.Name).To(Equal("customers"))
				}
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number <= 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &CustomersKey,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "<=", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(2))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(result.Documents[1].Content["testName"]).To(BeEquivalentTo(Customer2.Orders[0].Content["testName"]))
				Expect(*result.Documents[1].Key).To(BeEquivalentTo(Customer2.Orders[0].Key))
				Expect(*result.Documents[1].Key.Collection.Parent).To(BeEquivalentTo(Customer2.Key))
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number <= 1]", func() {
			It("Should return 1 order", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "<=", Value: "1"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(1))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
			})
		})
		When("key {customers, nil}, subcol: orders, exps: [type startsWith scooter]", func() {
			It("Should return 2 order", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &CustomersKey,
				}
				exps := []document.QueryExpression{
					{Operand: "type", Operator: "startsWith", Value: "scooter"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(2))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(result.Documents[1].Content["testName"]).To(BeEquivalentTo(Customer2.Orders[1].Content["testName"]))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*result.Documents[1].Key.Collection.Parent).To(BeEquivalentTo(Customer2.Key))
				Expect(*result.Documents[1].Key).To(BeEquivalentTo(Customer2.Orders[1].Key))
			})
		})
		When("key {customers, key1}, subcol: orders, exps: [type startsWith bike/road]", func() {
			It("Should return 1 order", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}
				exps := []document.QueryExpression{
					{Operand: "type", Operator: "startsWith", Value: "scooter"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(1))
				Expect(result.Documents[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(*result.Documents[0].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*result.Documents[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
			})
		})
		When("key: {items, nil}, subcol: '', exp: [], limit: 10", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				coll := document.Collection{
					Name: "items",
				}
				result, err := docPlugin.Query(context.TODO(), &coll, []document.QueryExpression{}, 10, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(10))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(context.TODO(), &coll, []document.QueryExpression{}, 10, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(2))
				Expect(result.PagingToken).To(BeNil())

				// Ensure values are unique
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}
			})
		})
		When("key: {items, nil}, subcol: '', exps: [letter > D], limit: 4", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				coll := document.Collection{
					Name: "items",
				}
				exps := []document.QueryExpression{
					{Operand: "letter", Operator: ">", Value: "D"},
				}
				result, err := docPlugin.Query(context.TODO(), &coll, exps, 4, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(context.TODO(), &coll, exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}

				result, err = docPlugin.Query(context.TODO(), &coll, exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))
				Expect(result.PagingToken).To(BeEmpty())
			})
		})
		When("key: {parentItems, 1}, subcol: items, exp: [], limit: 10", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				result, err := docPlugin.Query(context.TODO(), &ChildItemsCollection, []document.QueryExpression{}, 10, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(10))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(context.TODO(), &ChildItemsCollection, []document.QueryExpression{}, 10, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(2))
				Expect(result.PagingToken).To(BeNil())

				// Ensure values are unique
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}
			})
		})
		When("key: {parentItems, 1}, subcol: items, exps: [letter > D], limit: 4", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				exps := []document.QueryExpression{
					{Operand: "letter", Operator: ">", Value: "D"},
				}
				result, err := docPlugin.Query(context.TODO(), &ChildItemsCollection, exps, 4, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(context.TODO(), &ChildItemsCollection, exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				for i, d := range result.Documents {
					Expect(d.Key).ToNot(BeNil())
					val := fmt.Sprintf("%v", result.Documents[i].Content["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}

				result, err = docPlugin.Query(context.TODO(), &ChildItemsCollection, exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))
				Expect(result.PagingToken).To(BeEmpty())
			})
		})
	})
}
