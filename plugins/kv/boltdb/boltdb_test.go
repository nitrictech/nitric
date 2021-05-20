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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KV", func() {
	kvPlugin, err := kv_plugin.New()
	if err != nil {
		panic(err)
	}

	key := map[string]interface{}{
		"key": "jsmith@server.com",
	}
	testItem := map[string]interface{}{
		"firstName": "John",
		"lastName":  "Smith",
	}
	orderKey := map[string]interface{}{
		"pk": "Customer#1000",
		"sk": "Order#500",
	}
	orderItem := map[string]interface{}{
		"sku":    "ABC-123",
		"number": "1",
		"price":  "13.95",
	}

	AfterEach(func() {
		err := os.RemoveAll(kv_plugin.DEFAULT_DIR)
		if err != nil {
			panic(err)
		}

		_, err = os.Stat(kv_plugin.DEFAULT_DIR)
		if os.IsNotExist(err) {
			// Make diretory if not present
			err := os.Mkdir(kv_plugin.DEFAULT_DIR, 0777)
			if err != nil {
				panic(err)
			}
		}
	})

	AfterSuite(func() {
		err := os.RemoveAll(kv_plugin.DEFAULT_DIR)
		if err == nil {
			os.Remove(kv_plugin.DEFAULT_DIR)
			os.Remove("nitric/")
		}
	})

	Context("Put", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("", key, testItem)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("Nil key", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("collection", nil, testItem)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("Nil item map", func() {
			It("Should return error", func() {
				err := kvPlugin.Put("collection", key, nil)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("Valid Put", func() {
			It("Should store item successfully", func() {
				err := kvPlugin.Put("collection", key, testItem)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("collection", key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(testItem))
			})
		})

		When("Valid Compound Key Put", func() {
			It("Should store item successfully", func() {
				err := kvPlugin.Put("collection", orderKey, orderItem)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("collection", orderKey)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(orderItem))
			})
		})
	})

	Context("Get", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				_, err := kvPlugin.Get("", key)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("Nil key", func() {
			It("Should return error", func() {
				_, err := kvPlugin.Get("collection", nil)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("Valid Get", func() {
			It("Should get item successfully", func() {
				kvPlugin.Put("collection", key, testItem)

				doc, err := kvPlugin.Get("collection", key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(testItem))
			})
		})

		When("Valid Compound Key Get", func() {
			It("Should store item successfully", func() {
				kvPlugin.Put("collection", orderKey, orderItem)

				doc, err := kvPlugin.Get("collection", orderKey)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc).To(BeEquivalentTo(orderItem))
			})
		})
	})

	Context("Delete", func() {
		When("Blank collection", func() {
			It("Should return error", func() {
				err := kvPlugin.Delete("", key)
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
				kvPlugin.Put("collection", key, testItem)

				err := kvPlugin.Delete("collection", key)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("collection", key)
				Expect(doc).To(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		// TODO: missing item delete - discuss behaviour

		When("Valid Compound Key Delete", func() {
			It("Should delete item successfully", func() {
				kvPlugin.Put("collection", orderKey, orderItem)

				err := kvPlugin.Delete("collection", orderKey)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := kvPlugin.Get("collection", orderKey)
				Expect(doc).To(BeNil())
				Expect(err).ShouldNot(HaveOccurred())
			})
		})

		// TODO: missing item delete - discuss behaviour
	})

	// Context("Query", func() {
	// 	item1 := map[string]interface{}{
	// 		"email":     "john.smith@server.com",
	// 		"firstName": "John",
	// 		"lastName":  "Smith",
	// 	}
	// 	item2 := map[string]interface{}{
	// 		"email":     "paul.davis@server.com",
	// 		"firstName": "Paul",
	// 		"lastName":  "Davis",
	// 	}

	// 	BeforeEach(func() {
	// 		mockDbDriver.SetCollection("collection", map[string]interface{}{
	// 			"john.smith@server.com": item1,
	// 			"paul.davis@server.com": item2,
	// 		})
	// 	})

	// 	When("it does not exist", func() {
	// 		It("should cause en error", func() {
	// 			vals, err := kvPlugin.Query("collection", []sdk.QueryExpression{}, 0)
	// 			fmt.Println(vals)
	// 			Expect(err).To(BeNil())
	// 		})
	// 	})
	// })
})
