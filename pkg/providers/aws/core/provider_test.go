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

package core

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mocks "github.com/nitrictech/nitric/mocks/resourcetaggingapi"
)

var _ = Describe("AwsProvider", func() {
	When("Calling get resources with filled cache", func() {
		provider := &awsProviderImpl{
			cache: map[string]map[string]string{
				AwsResource_Bucket: {
					"test": "arn:aws:::test",
				},
			},
		}

		It("should return the map of available resources", func() {
			res, err := provider.GetResources(AwsResource_Bucket)

			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(res)).To(Equal(1))
		})
	})

	When("Calling get resources with no cache", func() {
		When("Call to GetResources fails", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockClient := mocks.NewMockResourceGroupsTaggingAPIAPI(ctrl)
			provider := &awsProviderImpl{
				cache:  make(map[string]map[string]string),
				client: mockClient,
			}

			It("should return an error", func() {
				defer ctrl.Finish()

				By("failing to call GetResources")
				mockClient.EXPECT().GetResources(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("mock-error"))

				res, err := provider.GetResources(AwsResource_Topic)

				By("returning an error")
				Expect(err).Should(HaveOccurred())

				By("returning nil resources")
				Expect(res).To(BeNil())
			})
		})

		When("Call to GetResources succeeds", func() {
			ctrl := gomock.NewController(GinkgoT())
			mockClient := mocks.NewMockResourceGroupsTaggingAPIAPI(ctrl)
			provider := &awsProviderImpl{
				cache:  make(map[string]map[string]string),
				client: mockClient,
				stack:  "test-stack",
			}

			It("should return available resources", func() {
				By("failing to call GetResources")
				mockClient.EXPECT().GetResources(gomock.Any(), gomock.Any()).Return(&resourcegroupstaggingapi.GetResourcesOutput{
					ResourceTagMappingList: []types.ResourceTagMapping{{
						ResourceARN: aws.String("arn:aws:::sns:test"),
						Tags: []types.Tag{{
							Key:   aws.String("x-nitric-name"),
							Value: aws.String("test"),
						}, {
							Key:   aws.String("x-nitric-stack"),
							Value: aws.String("test-stack"),
						}},
					}},
				}, nil)

				res, err := provider.GetResources(AwsResource_Topic)

				By("not returning an error")
				Expect(err).ShouldNot(HaveOccurred())

				By("returning resources")
				Expect(len(res)).To(Equal(1))
			})
		})
	})
})
