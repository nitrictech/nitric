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

package grpc_test

import (
	"context"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mock_secret "github.com/nitrictech/nitric/mocks/secret"
	"github.com/nitrictech/nitric/pkg/adapters/grpc"
	v1 "github.com/nitrictech/nitric/pkg/api/nitric/v1"
	"github.com/nitrictech/nitric/pkg/plugins/secret"
)

var _ = Describe("GRPC Secret", func() {
	Context("Put", func() {
		When("plugin not registered", func() {
			ss := &grpc.SecretServer{}
			resp, err := ss.Put(context.Background(), &v1.SecretPutRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Secret plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_secret.NewMockSecretService(g)
			resp, err := grpc.NewSecretServer(mockSS).Put(context.Background(), &v1.SecretPutRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid SecretPutRequest.Secret: value is required"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_secret.NewMockSecretService(g)

			val := []byte("hush")
			mockSS.EXPECT().Put(&secret.Secret{Name: "foo"}, val).Return(&secret.SecretPutResponse{
				SecretVersion: &secret.SecretVersion{
					Secret: &secret.Secret{
						Name: "foo",
					},
					Version: "3",
				},
			}, nil)

			resp, err := grpc.NewSecretServer(mockSS).Put(context.Background(), &v1.SecretPutRequest{
				Secret: &v1.Secret{Name: "foo"},
				Value:  val,
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(resp.SecretVersion.Version).To(Equal("3"))
				Expect(resp.SecretVersion.Secret.Name).To(Equal("foo"))
			})
		})
	})

	Context("Access", func() {
		When("plugin not registered", func() {
			ss := &grpc.SecretServer{}
			resp, err := ss.Access(context.Background(), &v1.SecretAccessRequest{})
			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("Secret plugin not registered"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request not valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_secret.NewMockSecretService(g)
			resp, err := grpc.NewSecretServer(mockSS).Access(context.Background(), &v1.SecretAccessRequest{})

			It("Should report an error", func() {
				Expect(err.Error()).Should(ContainSubstring("invalid SecretAccessRequest.SecretVersion: value is required"))
				Expect(resp).Should(BeNil())
			})
		})

		When("request is valid", func() {
			g := gomock.NewController(GinkgoT())
			mockSS := mock_secret.NewMockSecretService(g)

			mockSS.EXPECT().Access(&secret.SecretVersion{Secret: &secret.Secret{Name: "foo"}, Version: "3"}).Return(&secret.SecretAccessResponse{
				SecretVersion: &secret.SecretVersion{
					Secret: &secret.Secret{
						Name: "foo",
					},
					Version: "3",
				},
				Value: []byte("the value"),
			}, nil)

			resp, err := grpc.NewSecretServer(mockSS).Access(context.Background(), &v1.SecretAccessRequest{
				SecretVersion: &v1.SecretVersion{
					Secret: &v1.Secret{
						Name: "foo",
					},
					Version: "3",
				},
			})

			It("Should succeed", func() {
				Expect(err).Should(BeNil())
				Expect(resp.SecretVersion.Version).To(Equal("3"))
				Expect(resp.SecretVersion.Secret.Name).To(Equal("foo"))
				Expect(resp.Value).To(Equal([]byte("the value")))
			})
		})
	})
})
