package document_suite

import (
	"io"

	"github.com/nitric-dev/membrane/pkg/plugins/document"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func unwrapIter(iter document.DocumentIterator) []*document.Document {
	docs := make([]*document.Document, 0)
	for {
		d, err := iter()

		if err != nil {
			Expect(err).To(Equal(io.EOF))
			break
		}

		docs = append(docs, d)
	}

	return docs
}

func QueryStreamTests(docPlugin document.DocumentService) {
	Context("QueryStream", func() {
		// Validation Tests
		When("Invalid - blank key.Collection.Name", func() {
			It("Should return an iterator that errors", func() {
				iter := docPlugin.QueryStream(&document.Collection{}, []document.QueryExpression{}, 0)
				Expect(iter).ToNot(BeNil())

				_, err := iter()
				Expect(err).Should(HaveOccurred())
				Expect(err).ToNot(Equal(io.EOF))
			})
		})
		When("Invalid - nil expressions argument", func() {
			It("Should return an iterator that errors", func() {
				iter := docPlugin.QueryStream(&document.Collection{Name: "users"}, nil, 0)
				Expect(iter).ToNot(BeNil())

				_, err := iter()
				Expect(err).Should(HaveOccurred())
				Expect(err).ToNot(Equal(io.EOF))
			})
		})
		When("Empty database", func() {
			It("Should return io.EOF", func() {
				iter := docPlugin.QueryStream(&document.Collection{Name: "users"}, []document.QueryExpression{}, 0)

				Expect(iter).ToNot(BeNil())

				doc, err := iter()
				Expect(doc).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(io.EOF))
			})
		})

		// Query Tests
		When("key: {users}, subcol: '', exp: []", func() {
			It("Should return all users", func() {
				LoadUsersData(docPlugin)
				LoadCustomersData(docPlugin)

				iter := docPlugin.QueryStream(&document.Collection{Name: "users"}, []document.QueryExpression{}, 0)

				for {
					d, err := iter()

					if err != nil {
						Expect(err).To(Equal(io.EOF))
						break
					}

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

				iter := docPlugin.QueryStream(&CustomersColl, []document.QueryExpression{}, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(2))
				Expect(docs[0].Content["email"]).To(BeEquivalentTo(Customer1.Content["email"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Key))
				Expect(docs[1].Content["email"]).To(BeEquivalentTo(Customer2.Content["email"]))
				Expect(*docs[1].Key).To(BeEquivalentTo(Customer2.Key))
			})
		})
		When("key: {customers, nil}, subcol: '', exp: [country == US]", func() {
			It("Should return 1 item", func() {
				LoadCustomersData(docPlugin)

				exps := []document.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
				}

				iter := docPlugin.QueryStream(&CustomersColl, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(1))
				Expect(docs[0].Content["email"]).To(BeEquivalentTo(Customer2.Content["email"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer2.Key))
			})
		})
		When("key: {customers, nil}, subcol: '', exp: [country == US, age > 40]", func() {
			It("Should return 0 item", func() {
				LoadCustomersData(docPlugin)

				exps := []document.QueryExpression{
					{Operand: "country", Operator: "==", Value: "US"},
					{Operand: "age", Operator: ">", Value: "40"},
				}

				iter := docPlugin.QueryStream(&CustomersColl, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(0))
			})
		})
		When("key: {customers, key1}, subcol: orders", func() {
			It("Should return 3 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &Customer1.Key,
				}

				iter := docPlugin.QueryStream(&coll, []document.QueryExpression{}, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(3))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(docs[1].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[1].Content["testName"]))
				Expect(*docs[1].Key).To(BeEquivalentTo(Customer1.Orders[1].Key))
				Expect(*docs[1].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(docs[2].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(*docs[2].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*docs[2].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(2))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(docs[1].Content["testName"]).To(BeEquivalentTo(Customer2.Orders[0].Content["testName"]))
				Expect(*docs[1].Key).To(BeEquivalentTo(Customer2.Orders[0].Key))
				Expect(*docs[1].Key.Collection.Parent).To(BeEquivalentTo(Customer2.Key))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(1))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(3))

				for _, d := range docs {
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(2))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[1].Content["testName"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[1].Key))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(docs[1].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(*docs[1].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*docs[1].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
			})
		})
		When("key: {customers, nil}, subcol: orders, exps: [number < 1]", func() {
			It("Should return 0 orders", func() {
				LoadCustomersData(docPlugin)

				coll := document.Collection{
					Name:   "orders",
					Parent: &CustomersKey,
				}
				exps := []document.QueryExpression{
					{Operand: "number", Operator: "<", Value: "1"},
				}

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(0))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(0))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(5))

				for _, d := range docs {
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(3))

				for _, d := range docs {
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(2))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(docs[1].Content["testName"]).To(BeEquivalentTo(Customer2.Orders[0].Content["testName"]))
				Expect(*docs[1].Key).To(BeEquivalentTo(Customer2.Orders[0].Key))
				Expect(*docs[1].Key.Collection.Parent).To(BeEquivalentTo(Customer2.Key))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(1))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[0].Content["testName"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[0].Key))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(2))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(docs[1].Content["testName"]).To(BeEquivalentTo(Customer2.Orders[1].Content["testName"]))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*docs[1].Key.Collection.Parent).To(BeEquivalentTo(Customer2.Key))
				Expect(*docs[1].Key).To(BeEquivalentTo(Customer2.Orders[1].Key))
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

				iter := docPlugin.QueryStream(&coll, exps, 0)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(1))
				Expect(docs[0].Content["testName"]).To(BeEquivalentTo(Customer1.Orders[2].Content["testName"]))
				Expect(*docs[0].Key).To(BeEquivalentTo(Customer1.Orders[2].Key))
				Expect(*docs[0].Key.Collection.Parent).To(BeEquivalentTo(Customer1.Key))
			})
		})
		When("key: {items, nil}, subcol: '', exp: [], limit: 10", func() {
			It("Should return limited results", func() {
				LoadItemsData(docPlugin)

				coll := document.Collection{
					Name: "items",
				}

				iter := docPlugin.QueryStream(&coll, []document.QueryExpression{}, 10)
				docs := unwrapIter(iter)

				Expect(docs).To(HaveLen(10))
			})
		})
	})
}
