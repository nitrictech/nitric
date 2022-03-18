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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi/resourcegroupstaggingapiiface"

	"github.com/nitrictech/nitric/pkg/utils"
)

type AwsResource = string

const (
	AwsResource_Topic      AwsResource = "sns:topic"
	AwsResource_Collection AwsResource = "dynamodb:table"
	AwsResource_Queue      AwsResource = "sqs:queue"
	AwsResource_Bucket     AwsResource = "s3:bucket"
	AwsResource_Secret     AwsResource = "secretsmanager:secret"
)

type AwsProvider interface {
	// GetResources API operation for AWS Provider.
	// Returns requested aws resources for the given resource type
	GetResources(AwsResource) (map[string]string, error)
}

// Aws core utility provider
type awsProviderImpl struct {
	stack  string
	client resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	cache  map[AwsResource]map[string]string
}

var _ AwsProvider = &awsProviderImpl{}

func (a *awsProviderImpl) GetResources(typ AwsResource) (map[string]string, error) {
	if a.cache[typ] == nil {
		resources := make(map[string]string)
		tagFilters := []*resourcegroupstaggingapi.TagFilter{{
			Key: aws.String("x-nitric-name"),
		}}

		if a.stack != "" {
			tagFilters = append(tagFilters, &resourcegroupstaggingapi.TagFilter{
				Key:    aws.String("x-nitric-stack"),
				Values: []*string{aws.String(a.stack)},
			})
		}

		out, err := a.client.GetResources(&resourcegroupstaggingapi.GetResourcesInput{
			ResourceTypeFilters: []*string{aws.String(typ)},
			TagFilters:          tagFilters,
		})

		if err != nil {
			return nil, err
		}

		for _, tm := range out.ResourceTagMappingList {
			for _, t := range tm.Tags {
				if *t.Key == "x-nitric-name" {
					resources[*t.Value] = *tm.ResourceARN
					break
				}
			}
		}

		a.cache[typ] = resources
	}

	return a.cache[typ], nil
}

func New() (AwsProvider, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")
	stack := utils.GetEnv("NITRIC_STACK", "")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})

	if err != nil {
		return nil, err
	}

	client := resourcegroupstaggingapi.New(sess)

	return &awsProviderImpl{
		stack:  stack,
		client: client,
		cache:  make(map[AwsResource]map[string]string),
	}, nil
}
