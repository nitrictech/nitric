package document_suite

import (
	"github.com/nitric-dev/membrane/pkg/plugins/document"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
