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

package secrets_manager_secret_service

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	secretsmanager "github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/golang/mock/gomock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_provider "github.com/nitrictech/nitric/mocks/provider"
	mocks "github.com/nitrictech/nitric/mocks/secrets_manager"
	"github.com/nitrictech/nitric/pkg/plugins/secret"
	"github.com/nitrictech/nitric/pkg/providers/aws/core"
)

var _ = Describe("Secrets Manager Plugin", func() {
	testARN := "arn:partition:service:region:account-id:resource-id"
	testVersionID := "yVBWEvgpNjpcCxddXyj9kTefaUpVD999"
	testSecret := secret.Secret{
		Name: "Test",
	}

	testSecretVal := []byte("test")

	When("Put", func() {
		When("Given the Secrets Manager backend is available", func() {
			When("Putting a Secret to an existing secret", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				secretPlugin := &secretsManagerSecretService{
					provider: mockProvider,
					client:   mockSecretClient,
				}
				It("Should successfully store a secret", func() {
					defer ctrl.Finish()

					By("The secret container existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Secret).Return(
						map[string]string{
							"Test": testARN,
						}, nil,
					)

					By("The put operation succeeding")
					mockSecretClient.EXPECT().PutSecretValue(gomock.Any(),
						gomock.AssignableToTypeOf(&secretsmanager.PutSecretValueInput{}),
					).Return(&secretsmanager.PutSecretValueOutput{
						ARN:       aws.String(testARN),
						Name:      aws.String("Test"),
						VersionId: aws.String(testVersionID),
					}, nil).Times(1)

					response, err := secretPlugin.Put(context.TODO(), &testSecret, testSecretVal)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
					By("Returning a response with a version id")
					Expect(response.SecretVersion.Version).To(Equal(testVersionID))
				})
			})
			When("Putting a secret to a non-existent secret", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				secretPlugin := &secretsManagerSecretService{
					provider: mockProvider,
				}
				It("Should return error", func() {
					defer ctrl.Finish()

					By("the secret not existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Secret).Return(map[string]string{}, nil)

					_, err := secretPlugin.Put(context.TODO(), &testSecret, testSecretVal)
					By("returning an error")
					Expect(err).Should(HaveOccurred())
				})
			})
			When("Putting an empty secret", func() {
				secretPlugin := &secretsManagerSecretService{}

				It("Should return an error", func() {
					emptySecret := &secret.Secret{}
					response, err := secretPlugin.Put(context.TODO(), emptySecret, testSecretVal)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					By("Returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
			When("Putting a nil secret", func() {
				secretPlugin := &secretsManagerSecretService{}

				It("Should return an error", func() {
					response, err := secretPlugin.Put(context.TODO(), nil, testSecretVal)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					By("Returning a nil response")
					Expect(response).Should(BeNil())
				})
			})

			When("Putting a secret with a nil value", func() {
				secretPlugin := &secretsManagerSecretService{}

				It("Should return an error", func() {
					response, err := secretPlugin.Put(context.TODO(), &testSecret, nil)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					By("Returning a nil response")
					Expect(response).Should(BeNil())
				})
			})

			When("AWS SecretsManager.PutSecretValue returns an AWS error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client:   mockSecretClient,
					provider: mockProvider,
				}
				It("Should pass through the error", func() {
					defer ctrl.Finish()

					By("The secret existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Secret).Return(map[string]string{
						"Test": testARN,
					}, nil)

					mockSecretClient.EXPECT().PutSecretValue(
						gomock.Any(),
						gomock.Any(),
					).Return(nil, errors.New("aws error")).Times(1)

					response, err := secretPlugin.Put(context.TODO(), &testSecret, testSecretVal)
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("aws error"))
					By("returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
		})
	})
	When("Get", func() {
		When("Given the Secrets Manager backend is available", func() {
			When("The secret exists", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client:   mockSecretClient,
					provider: mockProvider,
				}
				It("Should return the existing secret", func() {
					defer ctrl.Finish()

					By("the secret existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Secret).Return(map[string]string{
						"Test": testARN,
					}, nil)

					mockSecretClient.EXPECT().GetSecretValue(gomock.Any(),
						&secretsmanager.GetSecretValueInput{
							SecretId:  aws.String(testARN),
							VersionId: aws.String("Version-Id"),
						},
					).Return(&secretsmanager.GetSecretValueOutput{
						ARN:          aws.String(testARN),
						Name:         aws.String("Test"),
						VersionId:    aws.String(testVersionID),
						SecretBinary: testSecretVal,
					}, nil).Times(1)

					response, err := secretPlugin.Access(context.TODO(), &secret.SecretVersion{
						Secret: &secret.Secret{
							Name: "Test",
						},
						Version: "Version-Id",
					})
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning a response with the secret")
					Expect(response).ShouldNot(BeNil())
					Expect(response.SecretVersion.Secret.Name).Should(Equal("Test"))
					Expect(response.SecretVersion.Version).Should(Equal("yVBWEvgpNjpcCxddXyj9kTefaUpVD999")) // Didn't return anything
					Expect(response.Value).Should(Equal(testSecretVal))
				})
			})
			When("The secret doesn't exist", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				secretPlugin := &secretsManagerSecretService{
					provider: mockProvider,
				}
				It("Should return a nil secret", func() {
					defer ctrl.Finish()

					By("The secret not existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Secret).Return(map[string]string{}, nil)

					response, err := secretPlugin.Access(context.TODO(), &secret.SecretVersion{
						Secret: &secret.Secret{
							Name: "test-id",
						},
						Version: "test-version-id",
					})
					By("Returning an error")
					Expect(err).Should(HaveOccurred())

					By("Returning a response with the secret")
					Expect(response).Should(BeNil())
				})
			})
			When("Getting the latest secret", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				mockProvider := mock_provider.NewMockAwsProvider(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client:   mockSecretClient,
					provider: mockProvider,
				}
				It("Should return the latest secret", func() {
					defer ctrl.Finish()

					By("The secret already existing")
					mockProvider.EXPECT().GetResources(gomock.Any(), core.AwsResource_Secret).Return(map[string]string{
						"test-id": testARN,
					}, nil)

					mockSecretClient.EXPECT().GetSecretValue(gomock.Any(),
						&secretsmanager.GetSecretValueInput{
							SecretId: aws.String(testARN),
						},
					).Return(&secretsmanager.GetSecretValueOutput{
						ARN:          aws.String(testARN),
						Name:         aws.String("Test"),
						VersionId:    aws.String(testVersionID),
						SecretBinary: testSecretVal,
					}, nil).Times(1)

					response, err := secretPlugin.Access(context.TODO(), &secret.SecretVersion{
						Secret: &secret.Secret{
							Name: "test-id",
						},
						Version: "latest",
					})
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())

					By("Returning a response with the secret")
					Expect(response).ShouldNot(BeNil())
					Expect(response.SecretVersion.Secret.Name).Should(Equal("test-id"))
					Expect(response.SecretVersion.Version).Should(Equal("yVBWEvgpNjpcCxddXyj9kTefaUpVD999")) // Didn't return anything
					Expect(response.Value).Should(Equal(testSecretVal))
				})
			})
			When("An empty id is provided", func() {
				secretPlugin := &secretsManagerSecretService{}

				response, err := secretPlugin.Access(context.TODO(), &secret.SecretVersion{
					Secret:  &secret.Secret{},
					Version: "test-version-id",
				})
				It("Should not return a secret", func() {
					By("Not returning an error")
					Expect(err).Should(HaveOccurred())

					By("Returning a response with the secret")
					Expect(response).Should(BeNil())
				})
			})
			When("An empty versionId is provided", func() {
				secretPlugin := &secretsManagerSecretService{}

				It("Should not return a secret", func() {
					response, err := secretPlugin.Access(context.TODO(), &secret.SecretVersion{
						Secret: &secret.Secret{
							Name: "test-id",
						},
						Version: "",
					})
					By("Not returning an error")
					Expect(err).Should(HaveOccurred())

					By("Returning a response with the secret")
					Expect(response).Should(BeNil())
				})
			})
		})
	})
})
