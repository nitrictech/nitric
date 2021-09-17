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

package secret_manager_secret_service

import (
	"fmt"

	"github.com/golang/mock/gomock"
	mocks "github.com/nitric-dev/membrane/mocks/secret_manager"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ = Describe("Secret Manager", func() {
	var mockSecret = &secretmanagerpb.Secret{
		Name:   "Test",
		Labels: make(map[string]string),
	}
	testSecret := secret.Secret{
		Name: "Test",
	}
	testSecretVal := []byte("Super Secret Message")

	When("Put", func() {
		When("Given the Secret Manager backend is available", func() {
			When("Putting a Secret to an existing secret", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretManagerClient(crtl)
				secretPlugin := &secretManagerSecretService{
					client:    mockSecretClient,
					projectId: "my-project",
				}

				It("Should successfully store a secret", func() {
					// Assert all methods are called at least their number of times
					defer crtl.Finish()
					//Mocking expects
					By("calling SecretManagerService.GetSecret with the expected request")
					mockSecretClient.EXPECT().GetSecret(
						gomock.Any(),
						&secretmanagerpb.GetSecretRequest{
							Name: "projects/my-project/secrets/Test",
						},
					).Return(&secretmanagerpb.Secret{
						Name: "projects/my-project/secrets/Test",
					}, nil).Times(1)

					By("not calling SecretManagerService.CreateSecret")
					mockSecretClient.EXPECT().CreateSecret(
						gomock.Any(),
						gomock.AssignableToTypeOf(&secretmanagerpb.CreateSecretRequest{}),
					).Return(mockSecret, nil).Times(0)

					By("Calling SecretManagerService AddSecretVersion with the expected payload")
					mockSecretClient.EXPECT().AddSecretVersion(
						gomock.Any(),
						&secretmanagerpb.AddSecretVersionRequest{
							Parent: "projects/my-project/secrets/Test",
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte("Super Secret Message"),
							},
						},
					).Return(&secretmanagerpb.SecretVersion{
						Name: "/projects/secrets/Test/versions/1",
					}, nil).Times(1)

					response, err := secretPlugin.Put(&testSecret, testSecretVal)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning the service provided version id")
					Expect(response.SecretVersion.Version).To(Equal("1"))
				})
			})

			When("Putting a Secret to a non-existent secret", func() {
				crtl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretManagerClient(crtl)
				secretPlugin := &secretManagerSecretService{
					client:    mockSecretClient,
					projectId: "my-project",
				}

				It("Should successfully store a secret", func() {
					defer crtl.Finish()

					//Mocking expects
					By("Calling SecretManagerService.GetSecret")
					mockSecretClient.EXPECT().GetSecret(
						gomock.Any(),
						&secretmanagerpb.GetSecretRequest{
							Name: "projects/my-project/secrets/Test",
						},
					).Return(nil, status.Error(codes.NotFound, "secret not found")).Times(1)

					By("Calling SecretManagerService.CreateSecret with the expected payload")
					mockSecretClient.EXPECT().CreateSecret(
						gomock.Any(),
						&secretmanagerpb.CreateSecretRequest{
							Parent:   "projects/my-project",
							SecretId: "Test",
							Secret: &secretmanagerpb.Secret{
								Replication: &secretmanagerpb.Replication{
									Replication: &secretmanagerpb.Replication_Automatic_{
										Automatic: &secretmanagerpb.Replication_Automatic{},
									},
								},
							},
						},
					).Return(&secretmanagerpb.Secret{
						Name: "projects/my-project/secrets/Test",
					}, nil).Times(1)

					By("Calling SecretManagerService.AddSecretVersion")
					mockSecretClient.EXPECT().AddSecretVersion(
						gomock.Any(),
						&secretmanagerpb.AddSecretVersionRequest{
							Parent: "projects/my-project/secrets/Test",
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte("Super Secret Message"),
							},
						},
					).Return(&secretmanagerpb.SecretVersion{
						Name: "/projects/my-project/Test/versions/1",
					}, nil).Times(1)

					response, err := secretPlugin.Put(&testSecret, testSecretVal)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning a response with a version id")
					Expect(response.SecretVersion.Version).To(Equal("1"))
				})
			})

			When("Putting a nil secret", func() {
				secretPlugin := &secretManagerSecretService{
					projectId: "my-project",
				}

				It("Should return an error", func() {
					_, err := secretPlugin.Put(nil, testSecretVal)
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("provide non-nil secret"))
				})
			})

			When("Putting a secret with an empty name", func() {
				secretPlugin := &secretManagerSecretService{
					projectId: "my-project",
				}

				It("Should return an error", func() {
					var emptySecretName = &secret.Secret{}
					_, err := secretPlugin.Put(emptySecretName, testSecretVal)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("provide non-blank secret name"))
				})
			})

			When("Putting a secret with an empty value", func() {
				secretPlugin := &secretManagerSecretService{
					projectId: "my-project",
				}

				It("Should return an error", func() {
					_, err := secretPlugin.Put(&testSecret, nil)

					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("provide non-blank secret value"))
				})
			})
		})
	})

	When("Get", func() {
		When("Given the Secret Manager backend is available", func() {
			When("The secret store exists", func() {
				When("The secret exists", func() {
					crtl := gomock.NewController(GinkgoT())
					mockSecretClient := mocks.NewMockSecretManagerClient(crtl)
					secretPlugin := &secretManagerSecretService{
						client:    mockSecretClient,
						projectId: "my-project",
					}
					It("Should successfully return a secret", func() {
						defer crtl.Finish()
						//Mocking expects
						By("calling SecretManagerService.AccessSecretVersion with the expected payload")
						mockSecretClient.EXPECT().AccessSecretVersion(
							gomock.Any(),
							&secretmanagerpb.AccessSecretVersionRequest{
								Name: "projects/my-project/secrets/test-id/versions/test-version-id",
							},
						).Return(&secretmanagerpb.AccessSecretVersionResponse{
							Name: "/projects/my-project/test-id/versions/test-version-id",
							Payload: &secretmanagerpb.SecretPayload{
								Data: []byte("Super Secret Message"),
							},
						}, nil).Times(1)
						response, err := secretPlugin.Access(&secret.SecretVersion{
							Secret: &secret.Secret{
								Name: "test-id",
							},
							Version: "test-version-id",
						})
						By("Not returning an error")
						Expect(err).ShouldNot(HaveOccurred())

						By("Returning a response with the secret")
						Expect(response).ShouldNot(BeNil())
						Expect(response.SecretVersion.Secret.Name).To(Equal("test-id"))
						Expect(response.SecretVersion.Version).To(Equal("test-version-id"))
						Expect(response.Value).To(Equal([]byte("Super Secret Message")))
					})
				})
				When("The secret doesn't exist", func() {
					crtl := gomock.NewController(GinkgoT())
					mockSecretClient := mocks.NewMockSecretManagerClient(crtl)
					secretPlugin := &secretManagerSecretService{
						client:    mockSecretClient,
						projectId: "my-project",
					}
					It("Should return an error", func() {
						defer crtl.Finish()

						mockSecretClient.EXPECT().AccessSecretVersion(
							gomock.Any(),
							&secretmanagerpb.AccessSecretVersionRequest{
								Name: "projects/my-project/secrets/test-id/versions/test-version-id",
							},
						).Return(nil, fmt.Errorf("failed to access secret")).Times(1)

						response, err := secretPlugin.Access(&secret.SecretVersion{
							Secret: &secret.Secret{
								Name: "test-id",
							},
							Version: "test-version-id",
						})

						By("returning an error")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("failed to access secret"))

						By("returning a nil response")
						Expect(response).Should(BeNil())
					})
				})
				When("An empty name is provided", func() {
					secretPlugin := &secretManagerSecretService{
						projectId: "my-project",
					}

					It("Should return an error", func() {
						response, err := secretPlugin.Access(&secret.SecretVersion{
							Secret: &secret.Secret{
								Name: "",
							},
							Version: "test-version-id",
						})

						By("returning an error")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("provide non-blank name"))

						By("returning a nil response")
						Expect(response).Should(BeNil())
					})
				})
				When("An empty version is provided", func() {
					secretPlugin := &secretManagerSecretService{
						projectId: "my-project",
					}
					It("Should return an error", func() {
						response, err := secretPlugin.Access(&secret.SecretVersion{
							Secret: &secret.Secret{
								Name: "test-id",
							},
							Version: "",
						})

						By("returning an error")
						Expect(err).Should(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("provide non-blank version"))

						By("returning a nil response")
						Expect(response).Should(BeNil())
					})
				})
			})
		})
	})
})
