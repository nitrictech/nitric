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
package test

import (
	"fmt"

	"github.com/nitric-dev/membrane/sdk"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Simple 'users' collection test data

var UserKey1 = sdk.Key{
	Collection: "users",
	Id:         "jsmith@server.com",
}
var UserItem1 = map[string]interface{}{
	"firstName": "John",
	"lastName":  "Smith",
	"email":     "jsmith@server.com",
	"country":   "US",
	"age":       "30",
}
var UserKey2 = sdk.Key{
	Collection: "users",
	Id:         "j.smithers@yahoo.com",
}
var UserItem2 = map[string]interface{}{
	"firstName": "Johnson",
	"lastName":  "Smithers",
	"email":     "j.smithers@yahoo.com",
	"country":   "AU",
	"age":       "40",
}
var UserKey3 = sdk.Key{
	Collection: "users",
	Id:         "pdavis@server.com",
}
var UserItem3 = map[string]interface{}{
	"firstName": "Paul",
	"lastName":  "Davis",
	"email":     "pdavis@server.com",
	"country":   "US",
	"age":       "50",
}

// Single Table Design 'customers' collection test data

var CustomersKey = sdk.Key{
	Collection: "customers",
}

type Customer struct {
	Key    sdk.Key
	Data   map[string]interface{}
	Orders []Order
}

type Order struct {
	Key  sdk.Key
	Data map[string]interface{}
}

var Customer1 = Customer{
	Key: sdk.Key{
		Collection: "customers",
		Id:         "1000",
	},
	Data: map[string]interface{}{
		"testName":  "CustomerItem1",
		"firstName": "John",
		"lastName":  "Smith",
		"email":     "jsmith@server.com",
		"country":   "AU",
		"age":       "40",
	},
	Orders: []Order{
		{
			Key: sdk.Key{
				Collection: "orders",
				Id:         "501",
			},
			Data: map[string]interface{}{
				"testName": "OrderItem1",
				"sku":      "ABC-501",
				"type":     "bike/mountain",
				"number":   "1",
				"price":    "14.95",
			},
		},
		{
			Key: sdk.Key{
				Collection: "orders",
				Id:         "502",
			},
			Data: map[string]interface{}{
				"testName": "OrderItem2",
				"sku":      "ABC-502",
				"type":     "bike/road",
				"number":   "2",
				"price":    "19.95",
			},
		},
		{
			Key: sdk.Key{
				Collection: "orders",
				Id:         "503",
			},
			Data: map[string]interface{}{
				"testName": "OrderItem3",
				"sku":      "ABC-503",
				"type":     "scooter/electric",
				"number":   "3",
				"price":    "124.95",
			},
		},
	},
}

var Customer2 = Customer{
	Key: sdk.Key{
		Collection: "customers",
		Id:         "2000",
	},
	Data: map[string]interface{}{
		"testName":  "CustomerItem2",
		"firstName": "David",
		"lastName":  "Adams",
		"email":     "dadams@server.com",
		"country":   "US",
		"age":       "20",
	},
	Orders: []Order{
		{
			Key: sdk.Key{
				Collection: "orders",
				Id:         "504",
			},
			Data: map[string]interface{}{
				"testName": "OrderItem4",
				"sku":      "ABC-504",
				"type":     "bike/hybrid",
				"number":   "1",
				"price":    "229.95",
			},
		},
		{
			Key: sdk.Key{
				Collection: "orders",
				Id:         "505",
			},
			Data: map[string]interface{}{
				"testName": "OrderItem5",
				"sku":      "ABC-505",
				"type":     "scooter/manual",
				"number":   "2",
				"price":    "9.95",
			},
		},
	},
}

type Item struct {
	Key  sdk.Key
	Data map[string]interface{}
}

var Items = []Item{
	{Key: sdk.Key{Collection: "items", Id: "01"}, Data: map[string]interface{}{"letter": "A"}},
	{Key: sdk.Key{Collection: "items", Id: "02"}, Data: map[string]interface{}{"letter": "B"}},
	{Key: sdk.Key{Collection: "items", Id: "03"}, Data: map[string]interface{}{"letter": "C"}},
	{Key: sdk.Key{Collection: "items", Id: "04"}, Data: map[string]interface{}{"letter": "D"}},
	{Key: sdk.Key{Collection: "items", Id: "05"}, Data: map[string]interface{}{"letter": "E"}},
	{Key: sdk.Key{Collection: "items", Id: "06"}, Data: map[string]interface{}{"letter": "F"}},
	{Key: sdk.Key{Collection: "items", Id: "07"}, Data: map[string]interface{}{"letter": "G"}},
	{Key: sdk.Key{Collection: "items", Id: "08"}, Data: map[string]interface{}{"letter": "H"}},
	{Key: sdk.Key{Collection: "items", Id: "09"}, Data: map[string]interface{}{"letter": "I"}},
	{Key: sdk.Key{Collection: "items", Id: "10"}, Data: map[string]interface{}{"letter": "J"}},
	{Key: sdk.Key{Collection: "items", Id: "11"}, Data: map[string]interface{}{"letter": "K"}},
	{Key: sdk.Key{Collection: "items", Id: "12"}, Data: map[string]interface{}{"letter": "L"}},
}

var ParentItemsKey = sdk.Key{
	Collection: "parentItems",
	Id:         "1",
}

// Test Data Loading Functions ------------------------------------------------

func LoadUsersData(kvPlugin sdk.DocumentService) {
	kvPlugin.Set(UserKey1, nil, UserItem1)
	kvPlugin.Set(UserKey2, nil, UserItem2)
	kvPlugin.Set(UserKey3, nil, UserItem3)
}

func LoadCustomersData(docPlugin sdk.DocumentService) {
	docPlugin.Set(Customer1.Key, nil, Customer1.Data)
	docPlugin.Set(Customer1.Key, &Customer1.Orders[0].Key, Customer1.Orders[0].Data)
	docPlugin.Set(Customer1.Key, &Customer1.Orders[1].Key, Customer1.Orders[1].Data)
	docPlugin.Set(Customer1.Key, &Customer1.Orders[2].Key, Customer1.Orders[2].Data)

	docPlugin.Set(Customer2.Key, nil, Customer2.Data)
	docPlugin.Set(Customer2.Key, &Customer2.Orders[0].Key, Customer2.Orders[0].Data)
	docPlugin.Set(Customer2.Key, &Customer2.Orders[1].Key, Customer2.Orders[1].Data)
}

func LoadItemsData(docPlugin sdk.DocumentService) {
	for _, item := range Items {
		docPlugin.Set(item.Key, nil, item.Data)
		docPlugin.Set(ParentItemsKey, &item.Key, item.Data)
	}
}

// Unit Test Functions --------------------------------------------------------

func GetTests(docPlugin sdk.DocumentService) {
	Context("Get", func() {
		When("Blank key.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Id: "1"}
				_, err := docPlugin.Get(key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Unknown key.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "unknown", Id: "1"}
				_, err := docPlugin.Get(key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank key.Id", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users"}
				_, err := docPlugin.Get(key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank parentKey.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users", Id: "1"}
				parentKey := sdk.Key{Id: "1"}
				_, err := docPlugin.Get(key, &parentKey)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank parentKey.Id", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users", Id: "1"}
				parentKey := sdk.Key{Collection: "users"}
				_, err := docPlugin.Get(key, &parentKey)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Get", func() {
			It("Should get item successfully", func() {
				docPlugin.Set(UserKey1, nil, UserItem1)

				doc, err := docPlugin.Get(UserKey1, nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(UserItem1["email"]))
			})
		})
		When("Valid Compound Key Get", func() {
			It("Should store item successfully", func() {
				docPlugin.Set(Customer1.Key, &Customer1.Orders[0].Key, Customer1.Orders[0].Data)

				doc, err := docPlugin.Get(Customer1.Key, &Customer1.Orders[0].Key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(Customer1.Orders[0].Data))
			})
		})
	})
}

func SetTests(docPlugin sdk.DocumentService) {
	Context("Set", func() {
		When("Blank key.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Id: "1"}
				err := docPlugin.Set(key, nil, UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Unknown key.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "unknown", Id: "1"}
				err := docPlugin.Set(key, nil, UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank key.Id", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users"}
				err := docPlugin.Set(key, nil, UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank parentKey.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users", Id: "1"}
				parentKey := sdk.Key{Id: "1"}
				err := docPlugin.Set(key, &parentKey, UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank parentKey.Id", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users", Id: "1"}
				parentKey := sdk.Key{Collection: "users"}
				err := docPlugin.Set(key, &parentKey, UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil item map", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users", Id: "1"}
				err := docPlugin.Set(key, nil, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid New Set", func() {
			It("Should store new item successfully", func() {
				err := docPlugin.Set(UserKey1, nil, UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(UserKey1, nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(UserItem1["email"]))
			})
		})
		When("Valid Update Set", func() {
			It("Should update existing item successfully", func() {
				err := docPlugin.Set(UserKey1, nil, UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(UserKey1, nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(UserItem1["email"]))

				err = docPlugin.Set(UserKey1, nil, UserItem2)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err = docPlugin.Get(UserKey1, nil)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc["email"]).To(BeEquivalentTo(UserItem2["email"]))
			})
		})
		When("Valid Compound Key Set", func() {
			It("Should store item successfully", func() {
				err := docPlugin.Set(Customer1.Key, &Customer1.Orders[0].Key, Customer1.Orders[0].Data)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(Customer1.Key, &Customer1.Orders[0].Key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(Customer1.Orders[0].Data))
			})
		})
	})
}

func DeleteTests(docPlugin sdk.DocumentService) {
	Context("Delete", func() {
		When("Blank key.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Id: "1"}
				err := docPlugin.Delete(key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Unknown key.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "unknown", Id: "1"}
				err := docPlugin.Delete(key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank key.Id", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users"}
				err := docPlugin.Delete(key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank parentKey.Collection", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users", Id: "1"}
				parentKey := sdk.Key{Id: "1"}
				err := docPlugin.Delete(key, &parentKey)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank parentKey.Id", func() {
			It("Should return error", func() {
				key := sdk.Key{Collection: "users", Id: "1"}
				parentKey := sdk.Key{Collection: "users"}
				err := docPlugin.Delete(key, &parentKey)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Delete", func() {
			It("Should delete item successfully", func() {
				docPlugin.Set(UserKey1, nil, UserItem1)

				err := docPlugin.Delete(UserKey1, nil)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(UserKey1, nil)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Compound Key Delete", func() {
			It("Should delete item successfully", func() {
				docPlugin.Set(Customer1.Key, &Customer1.Orders[0].Key, Customer1.Orders[0].Data)

				err := docPlugin.Delete(Customer1.Key, &Customer1.Orders[0].Key)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(Customer1.Key, &Customer1.Orders[0].Key)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
	})
}

func QueryTests(docPlugin sdk.DocumentService) {
	Context("Query", func() {
		When("Invalid - blank key.Collection", func() {
			It("Should return an error", func() {
				result, err := docPlugin.Query(sdk.Key{}, "", []sdk.QueryExpression{}, 0, nil)
				Expect(result).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Invalid - uknown key.Collection", func() {
			It("Should return an error", func() {
				result, err := docPlugin.Query(sdk.Key{Collection: "unknown"}, "", []sdk.QueryExpression{}, 0, nil)
				Expect(result).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Invalid - nil expressions argument", func() {
			It("Should return an error", func() {
				result, err := docPlugin.Query(sdk.Key{Collection: "users"}, "", nil, 0, nil)
				Expect(result).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Invalid - unkown subcol", func() {
			It("Should return an error", func() {
				// TODO: review sub-collection validation
				// result, err := docPlugin.Query(Customer1.Key, "unknown", []sdk.QueryExpression{}, 0, nil)
				// Expect(result).To(BeNil())
				// Expect(err).Should(HaveOccurred())
			})
		})
		When("Empty database", func() {
			It("Should return empty list", func() {
				result, err := docPlugin.Query(sdk.Key{Collection: "users"}, "", []sdk.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(0))
				Expect(result.PagingToken).To(BeNil())
			})
		})
		When("key: {users}, subcol: '', exp: []", func() {
			It("Should return all users", func() {
				LoadUsersData(docPlugin)
				LoadCustomersData(docPlugin)

				result, err := docPlugin.Query(sdk.Key{Collection: "users"}, "", []sdk.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(3))
			})
		})
		When("key: {users, key2}: subcol: '', exp: []", func() {
			It("Should return 1 user", func() {
				LoadUsersData(docPlugin)
				LoadCustomersData(docPlugin)

				result, err := docPlugin.Query(UserKey2, "", []sdk.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(1))
				Expect(result.Data[0]["email"]).To(BeEquivalentTo(UserItem2["email"]))
			})
		})
		When("key: {customers, unknown}", func() {
			It("Should return empty list", func() {
				LoadUsersData(docPlugin)
				LoadCustomersData(docPlugin)

				result, err := docPlugin.Query(sdk.Key{Collection: "users", Id: "unknown"}, "", []sdk.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(0))
			})
		})
		When("key: {customers, nil}, subcol: '', exp: []", func() {
			It("Should return 2 items", func() {
				LoadCustomersData(docPlugin)

				result, err := docPlugin.Query(CustomersKey, "", []sdk.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(2))
				Expect(result.Data[0]["email"]).To(BeEquivalentTo(Customer1.Data["email"]))
				Expect(result.Data[1]["email"]).To(BeEquivalentTo(Customer2.Data["email"]))
				Expect(result.PagingToken).To(BeNil())
			})
		})
		When("key: {customers, nil}, subcol: '', exp: [country == US]", func() {
			It("Should return 1 item", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
				}
				result, err := docPlugin.Query(CustomersKey, "", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(1))
				Expect(result.Data[0]["email"]).To(BeEquivalentTo(Customer2.Data["email"]))
				Expect(result.PagingToken).To(BeNil())
			})
		})
		When("key: {customers, nil}, subcol: '', exp: [country == US, age > 40]", func() {
			It("Should return 0 item", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
					{Operand: "age", Operator: ">", Value: "40"},
				}
				result, err := docPlugin.Query(CustomersKey, "", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(0))
			})
		})
		When("key: {customers, key1}, subcol: orders", func() {
			It("Should return 3 orders", func() {
				LoadCustomersData(docPlugin)

				result, err := docPlugin.Query(Customer1.Key, "orders", []sdk.QueryExpression{}, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(3))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Data["testName"]))
				Expect(result.Data[1]["testName"]).To(BeEquivalentTo(Customer1.Orders[1].Data["testName"]))
				Expect(result.Data[2]["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Data["testName"]))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number == 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "==", Value: "1"},
				}
				result, err := docPlugin.Query(CustomersKey, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(2))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Data["testName"]))
				Expect(result.Data[1]["testName"]).To(BeEquivalentTo(Customer2.Orders[0].Data["testName"]))
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number == 1]", func() {
			It("Should return 1 order", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "==", Value: "1"},
				}
				result, err := docPlugin.Query(Customer1.Key, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(1))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Data["testName"]))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number > 1]", func() {
			It("Should return 3 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: ">", Value: "1"},
				}
				result, err := docPlugin.Query(CustomersKey, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(3))
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number > 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: ">", Value: "1"},
				}
				result, err := docPlugin.Query(Customer1.Key, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(2))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[1].Data["testName"]))
				Expect(result.Data[1]["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Data["testName"]))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number < 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "<", Value: "1"},
				}
				result, err := docPlugin.Query(CustomersKey, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(0))
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number < 1]", func() {
			It("Should return 0 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "<", Value: "1"},
				}
				result, err := docPlugin.Query(Customer1.Key, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(0))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number >= 1]", func() {
			It("Should return 5 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: ">=", Value: "1"},
				}
				result, err := docPlugin.Query(CustomersKey, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(5))
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number >= 1]", func() {
			It("Should return 3 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: ">=", Value: "1"},
				}
				result, err := docPlugin.Query(Customer1.Key, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(3))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number <= 1]", func() {
			It("Should return 2 orders", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "<=", Value: "1"},
				}
				result, err := docPlugin.Query(CustomersKey, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(2))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Data["testName"]))
				Expect(result.Data[1]["testName"]).To(BeEquivalentTo(Customer2.Orders[0].Data["testName"]))
			})
		})
		When("key: {customers, key1}, subcol: orders, exps: [number <= 1]", func() {
			It("Should return 1 order", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "number", Operator: "<=", Value: "1"},
				}
				result, err := docPlugin.Query(Customer1.Key, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(1))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Data["testName"]))
			})
		})
		When("key {customers, nil}, subcol: orders, exps: [type startsWith scooter]", func() {
			It("Should return 2 order", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "type", Operator: "startsWith", Value: "scooter"},
				}
				result, err := docPlugin.Query(CustomersKey, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(2))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Data["testName"]))
				Expect(result.Data[1]["testName"]).To(BeEquivalentTo(Customer2.Orders[1].Data["testName"]))
			})
		})
		When("key {customers, key1}, subcol: orders, exps: [type startsWith bike/road]", func() {
			It("Should return 1 order", func() {
				LoadCustomersData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "type", Operator: "startsWith", Value: "scooter"},
				}
				result, err := docPlugin.Query(Customer1.Key, "orders", exps, 0, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(1))
				Expect(result.Data[0]["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Data["testName"]))
			})
		})
		When("key: {items, nil}, subcol: '', exp: [], limit: 10", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				result, err := docPlugin.Query(sdk.Key{Collection: "items"}, "", []sdk.QueryExpression{}, 10, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(10))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(sdk.Key{Collection: "items"}, "", []sdk.QueryExpression{}, 10, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(2))
				Expect(result.PagingToken).To(BeNil())

				// Ensure values are unique
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}
			})
		})
		When("key: {items, nil}, subcol: '', exps: [letter > D], limit: 4", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "letter", Operator: ">", Value: "D"},
				}
				result, err := docPlugin.Query(sdk.Key{Collection: "items"}, "", exps, 4, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(sdk.Key{Collection: "items"}, "", exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}

				result, err = docPlugin.Query(sdk.Key{Collection: "items"}, "", exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(0))
				Expect(result.PagingToken).To(BeEmpty())
			})
		})
		When("key: {parentItems, 1}, subcol: items, exp: [], limit: 10", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				result, err := docPlugin.Query(ParentItemsKey, "items", []sdk.QueryExpression{}, 10, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(10))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(ParentItemsKey, "items", []sdk.QueryExpression{}, 10, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(2))
				Expect(result.PagingToken).To(BeNil())

				// Ensure values are unique
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}
			})
		})
		When("key: {parentItems, 1}, subcol: items, exps: [letter > D], limit: 4", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				exps := []sdk.QueryExpression{
					{Operand: "letter", Operator: ">", Value: "D"},
				}
				result, err := docPlugin.Query(ParentItemsKey, "items", exps, 4, nil)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				dataMap := make(map[string]string)
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					dataMap[val] = val
				}

				result, err = docPlugin.Query(ParentItemsKey, "items", exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(4))
				Expect(result.PagingToken).ToNot(BeEmpty())

				// Ensure values are unique
				for i := range result.Data {
					val := fmt.Sprintf("%v", result.Data[i]["letter"])
					if _, found := dataMap[val]; found {
						Expect("matching value").ShouldNot(HaveOccurred())
					}
				}

				result, err = docPlugin.Query(ParentItemsKey, "items", exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Data).To(HaveLen(0))
				Expect(result.PagingToken).To(BeEmpty())
			})
		})

	})
}
