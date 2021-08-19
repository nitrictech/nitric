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

package key_vault_secret_service

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/armkeyvault"
	"github.com/golang/mock/gomock"
	mocks "github.com/nitric-dev/membrane/mocks/mock_key_vault"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Key Vault", func() {
	secretName := "secret-name"
	secretVersion := "secret-name/version-name"
	secretVal := []byte("Super Secret Message")
	secretString := string(secretVal)
	mockSecretResponse := armkeyvault.SecretResponse{
		Secret: &armkeyvault.Secret{
			Properties: &armkeyvault.SecretProperties{
				SecretURI:            &secretName,
				SecretURIWithVersion: &secretVersion,
				Value:                &secretString,
			},
		},
	}
	testSecret := &secret.Secret{
		Name: "secret-name",
	}
	testSecretVersion := &secret.SecretVersion{
		Secret:  testSecret,
		Version: secretVersion,
	}

	When("Put", func() {
		When("Given the Key Vault backend is available", func() {
			When("Putting a Secret to an existing secret", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}

				It("Should successfully store a secret", func() {
					// Assert all methods are called at least their number of times
					defer crtl.Finish()

					//Mocking expects
					mockSecretClient.EXPECT().CreateOrUpdate(
						context.Background(),
						secretPlugin.resourceGroupName,
						secretPlugin.vaultName,
						testSecret.Name,
						gomock.Any(),
						gomock.Any(),
					).Return(mockSecretResponse, nil).Times(1)

					response, err := secretPlugin.Put(testSecret, secretVal)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
					By("Returning the service provided version id")
					Expect(response.SecretVersion.Version).To(Equal(secretVersion))
					Expect(response.SecretVersion.Secret.Name).To(Equal(secretName))
				})
			})

			When("Putting a Secret to a non-existent secret", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}
				It("Should successfully store a secret", func() {
					defer crtl.Finish()

					//Mocking expects
					mockSecretClient.EXPECT().CreateOrUpdate(
						context.Background(),
						secretPlugin.resourceGroupName,
						secretPlugin.vaultName,
						testSecret.Name,
						gomock.Any(),
						gomock.Any(),
					).Return(mockSecretResponse, nil).Times(1)

					response, err := secretPlugin.Put(testSecret, secretVal)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
					By("Returning the correct secret")
					Expect(response.SecretVersion.Version).To(Equal(secretVersion))
					Expect(response.SecretVersion.Secret.Name).To(Equal(secretName))
				})
			})

			When("Putting a nil secret", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}

				It("Should invalidate the secret", func() {
					_, err := secretPlugin.Put(nil, secretVal)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})

			When("Putting a secret with an empty name", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}

				It("Should invalidate the secret", func() {
					_, err := secretPlugin.Put(&secret.Secret{Name: ""}, secretVal)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})

			When("Putting a secret with an empty value", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}

				It("Should invalidate the secret", func() {
					_, err := secretPlugin.Put(testSecret, nil)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
		})
	})

	When("Access", func() {
		When("Given the Key Vault backend is available", func() {
			When("The secret store exists", func() {
				When("The secret exists", func() {
					crtl := gomock.NewController(GinkgoT())
					mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
					secretPlugin := &keyVaultSecretService{
						client:            mockSecretClient,
						resourceGroupName: "resource-group-name",
						vaultName:         "vault-name",
					}

					It("Should successfully return a secret", func() {
						defer crtl.Finish()
						//Mocking expects
						mockSecretClient.EXPECT().Get(
							context.Background(),
							secretPlugin.resourceGroupName,
							secretPlugin.vaultName,
							secretVersion,
							gomock.Any(),
						).Return(mockSecretResponse, nil).Times(1)

						response, err := secretPlugin.Access(testSecretVersion)
						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())
						By("Returning the correct secret")
						Expect(response.Value).To(Equal(secretVal))
						Expect(response.SecretVersion).ToNot(BeNil())
						Expect(response.SecretVersion.Version).To(Equal(secretVersion))
						Expect(response.SecretVersion.Secret.Name).To(Equal(secretName))
					})
				})
			})
			When("The secret doesn't exist", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}
				It("Should return an error", func() {
					defer crtl.Finish()

					mockSecretClient.EXPECT().Get(
						context.Background(),
						secretPlugin.resourceGroupName,
						secretPlugin.vaultName,
						secretVersion,
						gomock.Any(),
					).Return(armkeyvault.SecretResponse{}, fmt.Errorf("secret does not exist")).Times(1)

					response, err := secretPlugin.Access(testSecretVersion)
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("secret does not exist"))
					By("returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
			When("An empty secret version is provided", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}
				It("Should return an error", func() {
					defer crtl.Finish()

					response, err := secretPlugin.Access(nil)
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					By("returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
			When("An empty secret is provided", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}
				It("Should return an error", func() {
					defer crtl.Finish()

					response, err := secretPlugin.Access(&secret.SecretVersion{Secret: nil, Version: secretVersion})
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					By("returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
			When("An empty secret name is provided", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}
				It("Should return an error", func() {
					defer crtl.Finish()

					response, err := secretPlugin.Access(&secret.SecretVersion{Secret: &secret.Secret{Name: ""}, Version: secretVersion})
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					By("returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
			When("An empty version is provided", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockKeyVaultClient(crtl)
				secretPlugin := &keyVaultSecretService{
					client:            mockSecretClient,
					resourceGroupName: "resource-group-name",
					vaultName:         "vault-name",
				}
				It("Should return an error", func() {
					defer crtl.Finish()

					response, err := secretPlugin.Access(&secret.SecretVersion{Secret: testSecret, Version: ""})
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					By("returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
		})
	})
})
