package document_suite

import (
	"github.com/nitric-dev/membrane/pkg/plugins/document"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
		When("Document Doesn't Exist", func() {
			It("Should return NotFound error", func() {
				key := document.Key{Collection: &document.Collection{Name: "items"}, Id: "not-exist"}
				doc, err := docPlugin.Get(&key)
				Expect(doc).To(BeNil())
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})
		When("Valid Collection Get when there is a Sub Collection", func() {
			It("Should store item successfully", func() {
				docPlugin.Set(&Customer1.Key, Customer1.Content)

				doc, err := docPlugin.Get(&Customer1.Key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Key).To(Equal(&Customer1.Key))
				Expect(doc.Content).To(BeEquivalentTo(Customer1.Content))
			})
		})
	})
}
