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

	"github.com/nitrictech/nitric/core/pkg/plugins/document"
	"github.com/nitrictech/nitric/core/pkg/utils"
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
	Reviews []Review
}

type Order struct {
	Key     document.Key
	Content map[string]interface{}
}

type Review struct {
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
	Reviews: []Review{
		{
			Key: document.Key{
				Collection: &document.Collection{
					Name: "reviews",
					Parent: &document.Key{
						Collection: &document.Collection{Name: "customers"},
						Id:         "1000",
					},
				},
				Id: "300",
			},
			Content: map[string]interface{}{
				"title": "Good review",
				"stars": "5",
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
	utils.Must(docPlugin.Set(context.TODO(), &UserKey1, UserItem1))
	utils.Must(docPlugin.Set(context.TODO(), &UserKey2, UserItem2))
	utils.Must(docPlugin.Set(context.TODO(), &UserKey3, UserItem3))
}

func LoadCustomersData(docPlugin document.DocumentService) {
	utils.Must(docPlugin.Set(context.TODO(), &Customer1.Key, Customer1.Content))
	utils.Must(docPlugin.Set(context.TODO(), &Customer1.Orders[0].Key, Customer1.Orders[0].Content))
	utils.Must(docPlugin.Set(context.TODO(), &Customer1.Orders[1].Key, Customer1.Orders[1].Content))
	utils.Must(docPlugin.Set(context.TODO(), &Customer1.Orders[2].Key, Customer1.Orders[2].Content))

	utils.Must(docPlugin.Set(context.TODO(), &Customer2.Key, Customer2.Content))
	utils.Must(docPlugin.Set(context.TODO(), &Customer2.Orders[0].Key, Customer2.Orders[0].Content))
	utils.Must(docPlugin.Set(context.TODO(), &Customer2.Orders[1].Key, Customer2.Orders[1].Content))
}

func LoadItemsData(docPlugin document.DocumentService) {
	for _, item := range Items {
		utils.Must(docPlugin.Set(context.TODO(), &item.Key, item.Content))

		key := document.Key{
			Collection: &ChildItemsCollection,
			Id:         item.Key.Id,
		}
		utils.Must(docPlugin.Set(context.TODO(), &key, item.Content))
	}
}

// Unit Test Functions --------------------------------------------------------
