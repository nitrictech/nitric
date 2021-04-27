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
