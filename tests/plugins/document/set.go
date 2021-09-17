package document_suite

import (
	"github.com/nitric-dev/membrane/pkg/plugins/document"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

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
		When("Valid Mutliple Sub Collection Set", func() {
			It("Should store item successfully", func() {
				err := docPlugin.Set(&Customer1.Reviews[0].Key, Customer1.Reviews[0].Content)
				Expect(err).ShouldNot(HaveOccurred())

				doc, err := docPlugin.Get(&Customer1.Reviews[0].Key)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(doc).ToNot(BeNil())
				Expect(doc.Content).To(BeEquivalentTo(Customer1.Reviews[0].Content))
			})
		})
	})
}
