package documents_plugin_test

import (
	documents_plugin "github.com/nitric-dev/membrane/plugins/dev/documents"
	"github.com/nitric-dev/membrane/plugins/dev/mocks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Documents", func() {
	mockDbDriver := mocks.NewMockScribble()
	documentsPlugin, _ := documents_plugin.NewWithDB(mockDbDriver)

	AfterEach(func() {
		mockDbDriver.ClearStore()
	})

	Context("CreateDocument", func() {
		When("the document doesn't yet exist", func() {
			It("Should successfully store the document", func() {
				testItem := map[string]interface{}{
					"Test": "Test",
				}
				err := documentsPlugin.CreateDocument("Test", "Test", testItem)

				Expect(err).ShouldNot(HaveOccurred())
				item := mockDbDriver.GetCollection("Test")["Test"]
				Expect(item).To(BeEquivalentTo(testItem))
			})
		})

		When("the document already exists", func() {
			BeforeEach(func() {
				mockDbDriver.SetCollection("Test", map[string]interface{}{
					"Test": map[string]interface{}{
						"Test": "Test",
					},
				})
			})

			It("Should return an error", func() {
				err := documentsPlugin.CreateDocument("Test", "Test", map[string]interface{}{
					"Test": "Test",
				})

				Expect(err).ToNot(BeNil())
			})
		})
	})

	Context("GetDocument", func() {
		item := map[string]interface{}{
			"Test": "Test",
		}

		When("the document exists", func() {
			BeforeEach(func() {
				mockDbDriver.SetCollection("Test", map[string]interface{}{
					"Test": item,
				})
			})

			It("should return the stored item", func() {
				gotItem, err := documentsPlugin.GetDocument("Test", "Test")

				Expect(err).ShouldNot(HaveOccurred())
				Expect(gotItem).To(BeEquivalentTo(item))
			})
		})

		When("the document does not exist", func() {
			It("should return an error", func() {
				gotItem, err := documentsPlugin.GetDocument("Test", "Test")

				Expect(err).Should(HaveOccurred())
				Expect(gotItem).To(BeNil())
			})
		})
	})

	Context("UpdateDocument", func() {
		item1 := map[string]interface{}{
			"Test": "Test",
		}
		item2 := map[string]interface{}{
			"Test": "Test2",
		}

		When("The document already exists", func() {
			BeforeEach(func() {
				mockDbDriver.SetCollection("Test", map[string]interface{}{
					"Test": item1,
				})
			})

			It("should update successfully", func() {
				err := documentsPlugin.UpdateDocument("Test", "Test", item2)
				Expect(err).ShouldNot(HaveOccurred())
				item := mockDbDriver.GetCollection("Test")["Test"]

				Expect(item).To(BeEquivalentTo(item2))
			})
		})

		When("It does not already exist", func() {
			It("should return an error", func() {
				err := documentsPlugin.UpdateDocument("Test", "Test", item2)
				Expect(err).Should(HaveOccurred())
			})
		})
	})

	Context("DeleteDocument", func() {
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
				err := documentsPlugin.DeleteDocument("Test", "Test")
				Expect(err).ShouldNot(HaveOccurred())
				Expect(mockDbDriver.GetCollection("Test")["Test"]).To(BeNil())
			})
		})

		When("it does not exist", func() {
			It("should cause en error", func() {
				err := documentsPlugin.DeleteDocument("Test", "Test")
				Expect(err).Should(HaveOccurred())
			})
		})
	})
})
