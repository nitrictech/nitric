package documents_plugin_test

import (
	"encoding/json"
	"fmt"

	documents_plugin "github.com/nitric-dev/membrane/plugins/dev/documents"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type MockScribbleDriver struct {
	documents_plugin.ScribbleIface
	store map[string]map[string][]byte
}

func (m *MockScribbleDriver) ensureCollectionExists(collection string) {
	if _, ok := m.store[collection]; !ok {
		m.store[collection] = make(map[string][]byte)
	}
}

func (m *MockScribbleDriver) clearStore() {
	m.store = make(map[string]map[string][]byte)
}

func (m *MockScribbleDriver) Read(collection string, key string, v interface{}) error {
	m.ensureCollectionExists(collection)

	if item, ok := m.store[collection][key]; ok {
		json.Unmarshal(item, &v)

		return nil
	}

	// TODO: This should produce the same error as stat()
	// in the case of a file not existing
	return fmt.Errorf("File does not exist!")
}

func (m *MockScribbleDriver) Write(collection string, key string, v interface{}) error {
	m.ensureCollectionExists(collection)

	bytes, _ := json.Marshal(v)

	m.store[collection][key] = bytes

	return nil
}

func (m *MockScribbleDriver) Delete(collection string, key string) error {
	m.ensureCollectionExists(collection)

	if _, ok := m.store[collection][key]; ok {
		m.store[collection][key] = nil
		return nil
	}

	return fmt.Errorf("File does not exist!")
}

func NewMockScribbleDriver() *MockScribbleDriver {
	return &MockScribbleDriver{
		store: make(map[string]map[string][]byte),
	}
}

var _ = Describe("Documents", func() {
	mockDbDriver := NewMockScribbleDriver()
	documentsPlugin, _ := documents_plugin.NewWithDB(mockDbDriver)

	AfterEach(func() {
		mockDbDriver.clearStore()
	})

	When("Creating a document", func() {
		When("the document doesn't yet exist", func() {
			It("Should successfully store the document", func() {
				err := documentsPlugin.CreateDocument("Test", "Test", map[string]interface{}{
					"Test": "Test",
				})

				Expect(err).To(BeNil())
			})
		})

		When("the document already exists", func() {
			BeforeEach(func() {
				mockDbDriver.store = map[string]map[string][]byte{
					"Test": map[string][]byte{
						"Test": []byte("{ \"Test\": \"Test\" }"),
					},
				}
			})

			It("Should return an error", func() {
				err := documentsPlugin.CreateDocument("Test", "Test", map[string]interface{}{
					"Test": "Test",
				})

				Expect(err).ToNot(BeNil())
			})
		})
	})

	When("Retrieving a document", func() {
		item := map[string]interface{}{
			"Test": "Test",
		}

		itemBytes, _ := json.Marshal(item)

		When("the document exists", func() {
			BeforeEach(func() {
				mockDbDriver.store = map[string]map[string][]byte{
					"Test": map[string][]byte{
						"Test": itemBytes,
					},
				}
			})

			It("should return the stored item", func() {
				gotItem, err := documentsPlugin.GetDocument("Test", "Test")

				Expect(err).To(BeNil())
				Expect(gotItem).To(BeEquivalentTo(item))
			})
		})

		When("the document does not exist", func() {
			It("should return an error", func() {
				gotItem, err := documentsPlugin.GetDocument("Test", "Test")

				Expect(err).ToNot(BeNil())
				Expect(gotItem).To(BeNil())
			})
		})
	})

	When("Updating a document", func() {
		item1 := map[string]interface{}{
			"Test": "Test",
		}
		item1Bytes, _ := json.Marshal(item1)
		item2 := map[string]interface{}{
			"Test": "Test2",
		}
		item2Bytes, _ := json.Marshal(item2)

		When("it exists", func() {
			BeforeEach(func() {
				mockDbDriver.store = map[string]map[string][]byte{
					"Test": map[string][]byte{
						"Test": item1Bytes,
					},
				}
			})

			It("should update successfully", func() {
				err := documentsPlugin.UpdateDocument("Test", "Test", item2)
				Expect(err).To(BeNil())
				Expect(mockDbDriver.store["Test"]["Test"]).To(BeEquivalentTo(item2Bytes))
			})
		})

		When("it does not exist", func() {
			It("should cause an error", func() {
				err := documentsPlugin.UpdateDocument("Test", "Test", item2)
				Expect(err).ToNot(BeNil())
			})
		})
	})

	When("Deleting a document", func() {
		item1 := map[string]interface{}{
			"Test": "Test",
		}
		item1Bytes, _ := json.Marshal(item1)

		When("it exists", func() {
			BeforeEach(func() {
				mockDbDriver.store = map[string]map[string][]byte{
					"Test": map[string][]byte{
						"Test": item1Bytes,
					},
				}
			})

			It("should delete successfully", func() {
				err := documentsPlugin.DeleteDocument("Test", "Test")
				Expect(err).To(BeNil())
				Expect(mockDbDriver.store["Test"]["Test"]).To(BeNil())
			})
		})

		When("it does not exist", func() {
			It("should cause en error", func() {
				err := documentsPlugin.DeleteDocument("Test", "Test")
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
