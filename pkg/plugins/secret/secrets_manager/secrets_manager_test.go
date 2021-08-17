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
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	secretsmanager "github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/golang/mock/gomock"
	mocks "github.com/nitric-dev/membrane/mocks/mock_secrets_manager"
	"github.com/nitric-dev/membrane/pkg/plugins/secret"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Secrets Manager Plugin", func() {
	var testARN = "arn:partition:service:region:account-id:resource-id"
	var testVersionID = "yVBWEvgpNjpcCxddXyj9kTefaUpVD999"
	testSecret := secret.Secret{
		Name: "Test",
	}

	testSecretVal := []byte("test")

	When("Put", func() {
		When("Given the Secrets Manager backend is available", func() {
			When("Putting a Secret to an existing secret", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client: mockSecretClient,
				}
				It("Should successfully store a secret", func() {
					defer ctrl.Finish()

					mockSecretClient.EXPECT().GetSecretValue(
						&secretsmanager.GetSecretValueInput{
							SecretId: aws.String("Test"),
						},
					).Return(&secretsmanager.GetSecretValueOutput{
						ARN:       aws.String(testARN),
						Name:      aws.String("Test"),
						VersionId: aws.String(testVersionID),
					}, nil).Times(1)
					mockSecretClient.EXPECT().PutSecretValue(
						gomock.AssignableToTypeOf(&secretsmanager.PutSecretValueInput{}),
					).Return(&secretsmanager.PutSecretValueOutput{
						ARN:       aws.String(testARN),
						Name:      aws.String("Test"),
						VersionId: aws.String(testVersionID),
					}, nil).Times(1)

					response, err := secretPlugin.Put(&testSecret, testSecretVal)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
					By("Returning a response with a version id")
					Expect(response.SecretVersion.Version).To(Equal(testVersionID))
				})
			})
			When("Putting a secret to a non-existent secret", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client: mockSecretClient,
				}
				It("Should successfully store a secret", func() {
					defer ctrl.Finish()

					mockSecretClient.EXPECT().GetSecretValue(
						&secretsmanager.GetSecretValueInput{
							SecretId: aws.String("Test"),
						},
					).Return(nil, awserr.New(secretsmanager.ErrCodeResourceNotFoundException, "secret does not exist", nil))
					mockSecretClient.EXPECT().CreateSecret(
						&secretsmanager.CreateSecretInput{
							Name:         aws.String("Test"),
							SecretBinary: testSecretVal,
						},
					).Return(&secretsmanager.CreateSecretOutput{
						ARN:       aws.String(testARN),
						Name:      aws.String("Test"),
						VersionId: aws.String(testVersionID),
					}, nil).Times(1)

					response, err := secretPlugin.Put(&testSecret, testSecretVal)
					By("Not returning an error")
					Expect(err).ShouldNot(HaveOccurred())
					By("Returning a response with a version id")
					Expect(response.SecretVersion.Version).To(Equal(testVersionID))
				})
			})
			When("Putting an empty secret", func() {
				secretPlugin := &secretsManagerSecretService{}

				It("Should return an error", func() {
					var emptySecret = &secret.Secret{}
					response, err := secretPlugin.Put(emptySecret, testSecretVal)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					By("Returning a nil response")
					Expect(response).Should(BeNil())
				})
			})
			When("Putting a nil secret", func() {
				secretPlugin := &secretsManagerSecretService{}

				It("Should return an error", func() {
					response, err := secretPlugin.Put(nil, testSecretVal)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					By("Returning a nil response")
					Expect(response).Should(BeNil())
				})
			})

			When("Putting a secret with a nil value", func() {
				secretPlugin := &secretsManagerSecretService{}

				It("Should return an error", func() {
					response, err := secretPlugin.Put(&testSecret, nil)
					By("Returning an error")
					Expect(err).Should(HaveOccurred())
					By("Returning a nil response")
					Expect(response).Should(BeNil())
				})
			})

			When("AWS SecretsManager.GetSecretValue returns an non NOT_FOUND error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client: mockSecretClient,
				}
				It("Should pass through the error", func() {
					defer ctrl.Finish()

					mockSecretClient.EXPECT().GetSecretValue(
						gomock.Any(),
					).Return(nil, fmt.Errorf("non-aws error")).Times(1)

					response, err := secretPlugin.Put(&testSecret, testSecretVal)
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("failed to retrieve secret container: \n non-aws error"))
					By("returning a nil response")
					Expect(response).Should(BeNil())
				})
			})

			When("AWS SecretsManager.PutSecretValue returns an AWS error", func() {
				ctrl := gomock.NewController(GinkgoT())
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client: mockSecretClient,
				}
				It("Should pass through the error", func() {
					defer ctrl.Finish()

					mockSecretClient.EXPECT().GetSecretValue(
						gomock.Any(),
					).Return(nil, awserr.New(secretsmanager.ErrCodeEncryptionFailure, "aws error", nil)).Times(1)
					response, err := secretPlugin.Put(&testSecret, testSecretVal)
					By("returning an error")
					Expect(err).Should(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("EncryptionFailure: aws error"))
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
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client: mockSecretClient,
				}
				It("Should return the existing secret", func() {
					defer ctrl.Finish()

					mockSecretClient.EXPECT().GetSecretValue(
						&secretsmanager.GetSecretValueInput{
							SecretId:  aws.String("Test"),
							VersionId: aws.String("Version-Id"),
						},
					).Return(&secretsmanager.GetSecretValueOutput{
						ARN:          aws.String(testARN),
						Name:         aws.String("Test"),
						VersionId:    aws.String(testVersionID),
						SecretBinary: testSecretVal,
					}, nil).Times(1)

					response, err := secretPlugin.Access(&secret.SecretVersion{
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
				mockSecretClient := mocks.NewMockSecretsManagerAPI(ctrl)
				secretPlugin := &secretsManagerSecretService{
					client: mockSecretClient,
				}
				It("Should return a nil secret", func() {
					defer ctrl.Finish()

					mockSecretClient.EXPECT().GetSecretValue(
						&secretsmanager.GetSecretValueInput{
							SecretId:  aws.String("test-id"),
							VersionId: aws.String("test-version-id"),
						},
					).Return(nil, fmt.Errorf("failed to get secret")).Times(1)

					response, err := secretPlugin.Access(&secret.SecretVersion{
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
				secretPlugin := &secretsManagerSecretService{
					client: mockSecretClient,
				}
				It("Should return the latest secret", func() {
					defer ctrl.Finish()

					mockSecretClient.EXPECT().GetSecretValue(
						&secretsmanager.GetSecretValueInput{
							SecretId: aws.String("test-id"),
						},
					).Return(&secretsmanager.GetSecretValueOutput{
						ARN:          aws.String(testARN),
						Name:         aws.String("Test"),
						VersionId:    aws.String(testVersionID),
						SecretBinary: testSecretVal,
					}, nil).Times(1)

					response, err := secretPlugin.Access(&secret.SecretVersion{
						Secret: &secret.Secret{
							Name: "test-id",
						},
						Version: "latest",
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
			When("An empty id is provided", func() {
				secretPlugin := &secretsManagerSecretService{}

				response, err := secretPlugin.Access(&secret.SecretVersion{
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
					response, err := secretPlugin.Access(&secret.SecretVersion{
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
