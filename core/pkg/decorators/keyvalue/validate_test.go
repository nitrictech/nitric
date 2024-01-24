// Copyright 2021 Nitric Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package keyvalue_test

import (
	document "github.com/nitrictech/nitric/core/pkg/decorators/keyvalue"
	keyvaluepb "github.com/nitrictech/nitric/core/pkg/proto/keyvalue/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Function Test Cases

var _ = Describe("Document Plugin", func() {
	When("ValidateKey", func() {
		When("Nil key", func() {
			It("should return error", func() {
				err := document.ValidateValueRef(nil)
				Expect(err.Error()).To(ContainSubstring("provide non-nil key"))
			})
		})
		When("Blank key.Collection", func() {
			It("should return error", func() {
				err := document.ValidateValueRef(&keyvaluepb.Key{})
				Expect(err.Error()).To(ContainSubstring("provide non-blank key.Id"))
			})
		})
		When("Blank key.Id", func() {
			It("should return error", func() {
				key := &keyvaluepb.Key{
					Store: "users",
				}
				err := document.ValidateValueRef(key)
				Expect(err.Error()).To(ContainSubstring("provide non-blank key.Id"))
			})
		})
	})
})
