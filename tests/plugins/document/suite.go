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
	"fmt"

	"github.com/nitric-dev/membrane/pkg/plugins/document"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Simple 'users' collection test data

var UserKey1 = document.Key{
	Collection: &document.Collection{Name: "users"},
	Id:         "jsmith@server.com",
}
var UserItem1 = map[string]interface{}{
	"firstName": "John",
	"lastName":  "Smith",
	"email":     "jsmith@server.com",
	"country":   "US",
	"age":       "30",
}
var UserKey2 = document.Key{
	Collection: &document.Collection{Name: "users"},
	Id:         "j.smithers@yahoo.com",
}
var UserItem2 = map[string]interface{}{
	"firstName": "Johnson",
	"lastName":  "Smithers",
	"email":     "j.smithers@yahoo.com",
	"country":   "AU",
	"age":       "40",
}
var UserKey3 = document.Key{
	Collection: &document.Collection{Name: "users"},
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

var CustomersKey = document.Key{
	Collection: &document.Collection{Name: "customers"},
}

var CustomersColl = document.Collection{Name: "customers"}

type Customer struct {
	Key     document.Key
	Content map[string]interface{}
	Orders  []Order
}

type Order struct {
	Key     document.Key
	Content map[string]interface{}
}

var Customer1 = Customer{
	Key: document.Key{
		Collection: &document.Collection{Name: "customers"},
		Id:         "1000",
	},
	Content: map[string]interface{}{
		"testName":  "CustomerItem1",
		"firstName": "John",
		"lastName":  "Smith",
		"email":     "jsmith@server.com",
		"country":   "AU",
		"age":       "40",
	},
	Orders: []Order{
		{
			Key: document.Key{
				Collection: &document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Collection: &document.Collection{Name: "customers"},
						Id:         "1000",
					},
				},
				Id: "501",
			},
			Content: map[string]interface{}{
				"testName": "OrderItem1",
				"sku":      "ABC-501",
				"type":     "bike/mountain",
				"number":   "1",
				"price":    "14.95",
			},
		},
		{
			Key: document.Key{
				Collection: &document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Collection: &document.Collection{Name: "customers"},
						Id:         "1000",
					},
				},
				Id: "502",
			},
			Content: map[string]interface{}{
				"testName": "OrderItem2",
				"sku":      "ABC-502",
				"type":     "bike/road",
				"number":   "2",
				"price":    "19.95",
			},
		},
		{
			Key: document.Key{
				Collection: &document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Collection: &document.Collection{Name: "customers"},
						Id:         "1000",
					},
				},
				Id: "503",
			},
			Content: map[string]interface{}{
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
	Key: document.Key{
		Collection: &document.Collection{Name: "customers"},
		Id:         "2000",
	},
	Content: map[string]interface{}{
		"testName":  "CustomerItem2",
		"firstName": "David",
		"lastName":  "Adams",
		"email":     "dadams@server.com",
		"country":   "US",
		"age":       "20",
	},
	Orders: []Order{
		{
			Key: document.Key{
				Collection: &document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Collection: &document.Collection{Name: "customers"},
						Id:         "2000",
					},
				},
				Id: "504",
			},
			Content: map[string]interface{}{
				"testName": "OrderItem4",
				"sku":      "ABC-504",
				"type":     "bike/hybrid",
				"number":   "1",
				"price":    "229.95",
			},
		},
		{
			Key: document.Key{
				Collection: &document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Collection: &document.Collection{Name: "customers"},
						Id:         "2000",
					},
				},
				Id: "505",
			},
			Content: map[string]interface{}{
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
	Key     document.Key
	Content map[string]interface{}
}

var Items = []Item{
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "01"},
		Content: map[string]interface{}{"letter": "A"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "02"},
		Content: map[string]interface{}{"letter": "B"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "03"},
		Content: map[string]interface{}{"letter": "C"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "04"},
		Content: map[string]interface{}{"letter": "D"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "05"},
		Content: map[string]interface{}{"letter": "E"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "06"},
		Content: map[string]interface{}{"letter": "F"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "07"},
		Content: map[string]interface{}{"letter": "G"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "08"},
		Content: map[string]interface{}{"letter": "H"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "09"},
		Content: map[string]interface{}{"letter": "I"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "10"},
		Content: map[string]interface{}{"letter": "J"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "11"},
		Content: map[string]interface{}{"letter": "K"},
	},
	{
		Key:     document.Key{Collection: &document.Collection{Name: "items"}, Id: "12"},
		Content: map[string]interface{}{"letter": "L"},
	},
}

var ChildItemsCollection = document.Collection{
	Name: "items",
	Parent: &document.Key{
		Collection: &document.Collection{Name: "parentItems"},
		Id:         "1",
	},
}

// Test Data Loading Functions ------------------------------------------------

func LoadUsersData(docPlugin document.DocumentService) {
	docPlugin.Set(&UserKey1, UserItem1)
	docPlugin.Set(&UserKey2, UserItem2)
	docPlugin.Set(&UserKey3, UserItem3)
}

func LoadCustomersData(docPlugin document.DocumentService) {
	docPlugin.Set(&Customer1.Key, Customer1.Content)
	docPlugin.Set(&Customer1.Orders[0].Key, Customer1.Orders[0].Content)
	docPlugin.Set(&Customer1.Orders[1].Key, Customer1.Orders[1].Content)
	docPlugin.Set(&Customer1.Orders[2].Key, Customer1.Orders[2].Content)

	docPlugin.Set(&Customer2.Key, Customer2.Content)
	docPlugin.Set(&Customer2.Orders[0].Key, Customer2.Orders[0].Content)
	docPlugin.Set(&Customer2.Orders[1].Key, Customer2.Orders[1].Content)
}

func LoadItemsData(docPlugin document.DocumentService) {
	for _, item := range Items {
		docPlugin.Set(&item.Key, item.Content)

		key := document.Key{
			Collection: &ChildItemsCollection,
			Id:         item.Key.Id,
		}
		docPlugin.Set(&key, item.Content)
	}
}

// Unit Test Functions --------------------------------------------------------

func GetTests(docPlugin document.DocumentService) {
	Context("Get", func() {
		When("Blank key.Collection.Name", func() {
			It("Should return error", func() {
				key := document.Key{Id: "1"}
				_, err := docPlugin.Get(&key)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank key.Id", func() {
			It("Should return error", func() {
				key := document.Key{Collection: &document.Collection{Name: "users"}}
				_, err := docPlugin.Get(&key)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Get", func() {
			It("Should get item successfully", func() {
				docPlugin.Set(&UserKey1, UserItem1)

				doc, err := docPlugin.Get(&UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Key).To(Equal(&UserKey1))
				Expect(doc.Content["email"]).To(BeEquivalentTo(UserItem1["email"]))
			})
		})
		When("Valid Sub Collection Get", func() {
			It("Should store item successfully", func() {
				docPlugin.Set(&Customer1.Orders[0].Key, Customer1.Orders[0].Content)

				doc, err := docPlugin.Get(&Customer1.Orders[0].Key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Key).To(Equal(&Customer1.Orders[0].Key))
				Expect(doc.Content).To(BeEquivalentTo(Customer1.Orders[0].Content))
			})
		})
	})
}

func SetTests(docPlugin document.DocumentService) {
	Context("Set", func() {
		When("Blank key.Collection.Name", func() {
			It("Should return error", func() {
				key := document.Key{Id: "1"}
				err := docPlugin.Set(&key, UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank key.Id", func() {
			It("Should return error", func() {
				key := document.Key{Collection: &document.Collection{Name: "users"}}
				err := docPlugin.Set(&key, UserItem1)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Nil item map", func() {
			It("Should return error", func() {
				key := document.Key{Collection: &document.Collection{Name: "users"}, Id: "1"}
				err := docPlugin.Set(&key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid New Set", func() {
			It("Should store new item successfully", func() {
				err := docPlugin.Set(&UserKey1, UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(&UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Content["email"]).To(BeEquivalentTo(UserItem1["email"]))
			})
		})
		When("Valid Update Set", func() {
			It("Should update existing item successfully", func() {
				err := docPlugin.Set(&UserKey1, UserItem1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(&UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Content["email"]).To(BeEquivalentTo(UserItem1["email"]))

				err = docPlugin.Set(&UserKey1, UserItem2)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err = docPlugin.Get(&UserKey1)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Content["email"]).To(BeEquivalentTo(UserItem2["email"]))
			})
		})
		When("Valid Sub Collection Set", func() {
			It("Should store item successfully", func() {
				err := docPlugin.Set(&Customer1.Orders[0].Key, Customer1.Orders[0].Content)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(&Customer1.Orders[0].Key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Content).To(BeEquivalentTo(Customer1.Orders[0].Content))
			})
		})
	})
}

func DeleteTests(docPlugin document.DocumentService) {
	Context("Delete", func() {
		When("Blank key.Collection.Name", func() {
			It("Should return error", func() {
				key := document.Key{Id: "1"}
				err := docPlugin.Delete(&key)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Blank key.Id", func() {
			It("Should return error", func() {
				key := document.Key{Collection: &document.Collection{Name: "users"}}
				err := docPlugin.Delete(&key)
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Delete", func() {
			It("Should delete item successfully", func() {
				docPlugin.Set(&UserKey1, UserItem1)

				err := docPlugin.Delete(&UserKey1)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(&UserKey1)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Sub Collection Delete", func() {
			It("Should delete item successfully", func() {
				docPlugin.Set(&Customer1.Orders[0].Key, Customer1.Orders[0].Content)

				err := docPlugin.Delete(&Customer1.Orders[0].Key)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(&Customer1.Orders[0].Key)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Valid Parent and Sub Collection Delete", func() {
			It("Should delete all children", func() {
				LoadCustomersData(docPlugin)

				col := document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Collection: &document.Collection{
							Name: "customers",
						},
					},
				}

				result, err := docPlugin.Query(&col, []document.QueryExpression{}, 0, nil)
				Expect(err).To(BeNil())
				Expect(result.Documents).To(HaveLen(5))

				err = docPlugin.Delete(&Customer1.Key)
				Expect(err).ShouldNot(HaveOccurred())

				err = docPlugin.Delete(&Customer2.Key)
				Expect(err).ShouldNot(HaveOccurred())

				result, err = docPlugin.Query(&col, []document.QueryExpression{}, 0, nil)
				Expect(err).To(BeNil())
				Expect(result.Documents).To(HaveLen(0))
			})
		})
	})
}

func QueryTests(docPlugin document.DocumentService) {
	Context("Query", func() {
		When("Invalid - blank key.Collection.Name", func() {
			It("Should return an error", func() {
				result, err := docPlugin.Query(&document.Collection{}, []document.QueryExpression{}, 0, nil)
				Expect(result).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Invalid - nil expressions argument", func() {
			It("Should return an error", func() {
				result, err := docPlugin.Query(&document.Collection{Name: "users"}, nil, 0, nil)
				Expect(result).To(BeNil())
				Expect(err).Should(HaveOccurred())
			})
		})
		When("Empty database", func() {
			It("Should return empty list", func() {
				result, err := docPlugin.Query(&document.Collection{Name: "users"}, []document.QueryExpression{}, 0, nil)
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

				result, err := docPlugin.Query(&document.Collection{Name: "users"}, []document.QueryExpression{}, 0, nil)
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

				result, err := docPlugin.Query(&CustomersColl, []document.QueryExpression{}, 0, nil)
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
				result, err := docPlugin.Query(&CustomersColl, exps, 0, nil)
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
				result, err := docPlugin.Query(&CustomersColl, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, []document.QueryExpression{}, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, exps, 0, nil)
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
				result, err := docPlugin.Query(&coll, []document.QueryExpression{}, 10, nil)
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

				result, err = docPlugin.Query(&coll, []document.QueryExpression{}, 10, result.PagingToken)
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
				result, err := docPlugin.Query(&coll, exps, 4, nil)
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

				result, err = docPlugin.Query(&coll, exps, 4, result.PagingToken)
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

				result, err = docPlugin.Query(&coll, exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))
				Expect(result.PagingToken).To(BeEmpty())
			})
		})
		When("key: {parentItems, 1}, subcol: items, exp: [], limit: 10", func() {
			It("Should return have multiple pages", func() {
				LoadItemsData(docPlugin)

				result, err := docPlugin.Query(&ChildItemsCollection, []document.QueryExpression{}, 10, nil)
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

				result, err = docPlugin.Query(&ChildItemsCollection, []document.QueryExpression{}, 10, result.PagingToken)
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
				result, err := docPlugin.Query(&ChildItemsCollection, exps, 4, nil)
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

				result, err = docPlugin.Query(&ChildItemsCollection, exps, 4, result.PagingToken)
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

				result, err = docPlugin.Query(&ChildItemsCollection, exps, 4, result.PagingToken)
				Expect(result).ToNot(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
				Expect(result.Documents).To(HaveLen(0))
				Expect(result.PagingToken).To(BeEmpty())
			})
		})
	})
}
