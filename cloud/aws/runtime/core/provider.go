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
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/apigatewayv2iface"
	"github.com/nitrictech/nitric/cloud/aws/ifaces/resourcegroupstaggingapiiface"
	"github.com/nitrictech/nitric/core/pkg/providers/common"
	"github.com/nitrictech/nitric/core/pkg/utils"
)

type AwsResource = string

const (
	AwsResource_Api          AwsResource = "apigateway:apis"
	AwsResource_StateMachine AwsResource = "states:stateMachine"
	AwsResource_Topic        AwsResource = "sns:topic"
	AwsResource_Collection   AwsResource = "dynamodb:table"
	AwsResource_Queue        AwsResource = "sqs:queue"
	AwsResource_Bucket       AwsResource = "s3:bucket"
	AwsResource_Secret       AwsResource = "secretsmanager:secret"
	AwsResource_EventRule    AwsResource = "events:rule"
)

var resourceTypeMap = map[common.ResourceType]AwsResource{
	common.ResourceType_Api: AwsResource_Api,
}

type AwsProvider interface {
	common.ResourceService
	// GetResources API operation for AWS Provider.
	// Returns requested aws resources for the given resource type
	GetResources(context.Context, AwsResource) (map[string]string, error)
}

// Aws core utility provider
type awsProviderImpl struct {
	stack     string
	client    resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	apiClient apigatewayv2iface.ApiGatewayV2API
	cache     map[AwsResource]map[string]string
}

var _ AwsProvider = &awsProviderImpl{}

func (a *awsProviderImpl) Details(ctx context.Context, typ common.ResourceType, name string) (*common.DetailsResponse[any], error) {
	rt, ok := resourceTypeMap[typ]
	if !ok {
		return nil, fmt.Errorf("unhandled resource type: %s", typ)
	}

	// Get resource references (arns) for the resource type
	resources, err := a.GetResources(ctx, rt)
	if err != nil {
		return nil, err
	}

	arn, ok := resources[name]
	if !ok {
		return nil, fmt.Errorf("unable to find resource %s for name: %s", typ, name)
	}

	details := &common.DetailsResponse[any]{
		Id:       arn,
		Provider: "aws",
	}

	switch rt {
	case AwsResource_Api:
		// split arn to find the apiId
		arnParts := strings.Split(arn, "/")
		apiId := arnParts[len(arnParts)-1]
		// Get api detail
		api, err := a.apiClient.GetApi(context.TODO(), &apigatewayv2.GetApiInput{
			ApiId: aws.String(apiId),
		})
		if err != nil {
			return nil, err
		}

		details.Service = "ApiGateway"
		details.Detail = common.ApiDetails{
			URL: *api.ApiEndpoint,
		}

		return details, nil
	default:
		return nil, fmt.Errorf("unimplemented resource type")
	}
}

func (a *awsProviderImpl) GetResources(ctx context.Context, typ AwsResource) (map[string]string, error) {
	if a.cache[typ] == nil {
		resources := make(map[string]string)
		tagFilters := []types.TagFilter{{
			Key: aws.String("x-nitric-name"),
		}}

		if a.stack != "" {
			tagFilters = append(tagFilters, types.TagFilter{
				Key:    aws.String("x-nitric-stack"),
				Values: []string{a.stack},
			})
		}

		out, err := a.client.GetResources(ctx, &resourcegroupstaggingapi.GetResourcesInput{
			ResourceTypeFilters: []string{typ},
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

	cfg, sessionError := config.LoadDefaultConfig(context.TODO(), config.WithRegion(awsRegion))
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	apiClient := apigatewayv2.NewFromConfig(cfg)
	client := resourcegroupstaggingapi.NewFromConfig(cfg)

	return &awsProviderImpl{
		stack:     stack,
		client:    client,
		apiClient: apiClient,
		cache:     make(map[AwsResource]map[string]string),
	}, nil
}
