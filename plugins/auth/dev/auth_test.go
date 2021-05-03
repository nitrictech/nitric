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

package auth_service_test

import (
	mocks "github.com/nitric-dev/membrane/mocks/scribble"
	auth_plugin "github.com/nitric-dev/membrane/plugins/auth/dev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("Auth", func() {
	mockScribble := mocks.NewMockScribble()
	authPlugin := auth_plugin.NewWithDriver(mockScribble)

	AfterEach(func() {
		mockScribble.ClearStore()
	})

	Context("Create", func() {
		When("The user does not already exist", func() {
			It("Should successfully create the user", func() {
				err := authPlugin.Create("test", "test", "test@test.com", "test")

				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Storing the user")
				testUser, ok := mockScribble.GetCollection("auth_test")["test"].(map[string]interface{})

				Expect(ok).To(Equal(true))
				Expect(testUser["id"]).To(Equal("test"))
				Expect(testUser["email"]).To(Equal("test@test.com"))
				Expect(bcrypt.CompareHashAndPassword([]byte(testUser["pwdHashAndSalt"].(string)), []byte("test"))).ShouldNot(HaveOccurred())
			})
		})

		When("A user with the same id already exists", func() {
			BeforeEach(func() {
				// Setup the existing user...
				mockPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
				mockScribble.SetCollection("auth_test", map[string]interface{}{
					"test": map[string]interface{}{
						"id":             "test",
						"email":          "test@test.com",
						"pwdHashAndSalt": mockPassword,
					},
				})
			})

			It("Should return an error", func() {
				err := authPlugin.Create("test", "test", "test2@test.com", "test")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("with id"))
			})
		})

		When("A user with the same email already exists", func() {
			BeforeEach(func() {
				// Setup the existing user...
				mockPassword, _ := bcrypt.GenerateFromPassword([]byte("test"), bcrypt.DefaultCost)
				mockScribble.SetCollection("auth_test", map[string]interface{}{
					"test": map[string]interface{}{
						"id":             "test1",
						"email":          "test2@test.com",
						"pwdHashAndSalt": mockPassword,
					},
				})
			})

			It("Should return an error", func() {
				err := authPlugin.Create("test", "test", "test2@test.com", "test")

				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("with email"))
			})
		})
	})
})
