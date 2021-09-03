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

package secret_service_test

import (
	"os"

	"github.com/nitric-dev/membrane/pkg/utils"

	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	secretPlugin "github.com/nitric-dev/membrane/pkg/plugins/secret/dev"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dev Secret Manager", func() {
	AfterSuite(func() {
		// Cleanup default secret directory
		os.RemoveAll(utils.GetDevVolumePath())
	})

	testSecret := secret.Secret{
		Name: "Test",
	}
	testSecretVal := []byte("Super Secret Message")
	When("Put", func() {
		When("Putting a secret to a non-existent secret", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should successfully store a secret", func() {
				response, err := secretPlugin.Put(&testSecret, testSecretVal)
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning a non nil response")
				Expect(response).ShouldNot(BeNil())
			})
		})
		When("Putting a secret to an existing secret", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should succesfully store a secret", func() {
				response, err := secretPlugin.Put(&testSecret, testSecretVal)
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("Returning a non nil response")
				Expect(response).ShouldNot(BeNil())
			})
		})
		When("Putting a secret with an empty name", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should throw an error", func() {
				emptySecretName := &secret.Secret{}
				response, err := secretPlugin.Put(emptySecretName, testSecretVal)
				By("Returning an error")
				Expect(err).Should(HaveOccurred())

				By("Returning a nil response")
				Expect(response).Should(BeNil())
			})
		})
	})
	When("Get", func() {
		When("Getting a secret that exists", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should return the secret", func() {
				putResponse, _ := secretPlugin.Put(&testSecret, testSecretVal)
				response, err := secretPlugin.Access(putResponse.SecretVersion)
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
				By("Returning a response")
				Expect(response.SecretVersion.Secret.Name).Should(Equal(testSecret.Name))
				Expect(response.Value).Should(Equal(testSecretVal))
			})
		})
		When("Getting the latest secret", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should return the latest secret", func() {
				putResponse, _ := secretPlugin.Put(&testSecret, testSecretVal)
				response, err := secretPlugin.Access(&secret.SecretVersion{
					Secret: &secret.Secret{
						Name: testSecret.Name,
					},
					Version: "latest",
				})
				By("Not returning an error")
				Expect(err).ShouldNot(HaveOccurred())
				By("Returning a response")
				Expect(response.SecretVersion.Secret.Name).Should(Equal(testSecret.Name))
				Expect(response.SecretVersion.Version).Should(Equal(putResponse.SecretVersion.Version))
				Expect(response.Value).Should(Equal(testSecretVal))
			})
		})
		When("Getting a secret that doesn't exist", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should return an error", func() {
				response, err := secretPlugin.Access(&secret.SecretVersion{
					Secret: &secret.Secret{
						Name: "test-id",
					},
					Version: "test-version-id",
				})
				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				By("Returning a nil response")
				Expect(response).Should(BeNil())
			})
		})
		When("Getting a secret with an empty id", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should return an error", func() {
				response, err := secretPlugin.Access(&secret.SecretVersion{
					Secret:  &secret.Secret{},
					Version: "test-version-id",
				})
				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				By("Returning a nil response")
				Expect(response).Should(BeNil())
			})
		})
		When("Getting a secret with an empty version id", func() {
			secretPlugin, _ := secretPlugin.New()
			It("Should return an error", func() {
				response, err := secretPlugin.Access(&secret.SecretVersion{
					Secret: &secret.Secret{
						Name: "test-id",
					},
				})
				By("Returning an error")
				Expect(err).Should(HaveOccurred())
				By("Returning a nil response")
				Expect(response).Should(BeNil())
			})
		})
	})
})
