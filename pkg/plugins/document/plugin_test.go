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
package document_test

import (
	"github.com/nitric-dev/membrane/pkg/adapters/grpc"
	"github.com/nitric-dev/membrane/pkg/plugins/document"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Document Plugin", func() {

	When("Collection.String", func() {
		It("should print collection", func() {
			col := &document.Collection{Name: "customer"}
			log := grpc.LogArg(col)
			Expect(log).To(BeEquivalentTo("{Name: customer, Parent: <nil>}"))
		})
	})

	When("Key.String", func() {
		It("should print subcollection key", func() {
			key := document.Key{
				Collection: &document.Collection{
					Name: "orders",
					Parent: &document.Key{
						Collection: &document.Collection{Name: "customers"},
						Id:         "user@mail.com",
					},
				},
				Id: "501",
			}
			log := grpc.LogArg(key)
			Expect(log).To(BeEquivalentTo("{Collection: {Name: orders, Parent: {Collection: {Name: customers, Parent: <nil>}, Id: user@mail.com}}, Id: 501}"))
		})
	})
})
