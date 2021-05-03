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

package kv_service_test

import (
	mocks "github.com/nitric-dev/membrane/mocks/scribble"
	kv_plugin "github.com/nitric-dev/membrane/plugins/kv/dev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("KV", func() {
	mockDbDriver := mocks.NewMockScribble()
	kvPlugin, _ := kv_plugin.NewWithDB(mockDbDriver)

	AfterEach(func() {
		mockDbDriver.ClearStore()
	})

	Context("Put", func() {
		It("Should successfully store the document", func() {
			testItem := map[string]interface{}{
				"Test": "Test",
			}
			err := kvPlugin.Put("Test", "Test", testItem)

			Expect(err).ShouldNot(HaveOccurred())
			item := mockDbDriver.GetCollection("Test")["Test"]
			Expect(item).To(BeEquivalentTo(testItem))
		})
	})

	Context("Get", func() {
		item := map[string]interface{}{
			"Test": "Test",
		}

		When("the key exists", func() {
			BeforeEach(func() {
				mockDbDriver.SetCollection("Test", map[string]interface{}{
					"Test": item,
				})
			})

			It("should return the stored item", func() {
				gotItem, err := kvPlugin.Get("Test", "Test")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(gotItem).To(BeEquivalentTo(item))
			})
		})

		When("the key does not exist", func() {
			It("should return an error", func() {
				gotItem, err := kvPlugin.Get("Test", "Test")

				Expect(err).Should(HaveOccurred())
				Expect(gotItem).To(BeNil())
			})
		})
	})

	Context("Delete", func() {
		item1 := map[string]interface{}{
			"Test": "Test",
		}

		When("it exists", func() {
			BeforeEach(func() {
				mockDbDriver.SetCollection("Test", map[string]interface{}{
					"Test": item1,
				})
			})

			It("should delete successfully", func() {
				err := kvPlugin.Delete("Test", "Test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(mockDbDriver.GetCollection("Test")["Test"]).To(BeNil())
			})
		})

		When("it does not exist", func() {
			It("should cause en error", func() {
				err := kvPlugin.Delete("Test", "Test")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
