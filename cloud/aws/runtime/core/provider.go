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
	"sync"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsArn "github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/apigatewayv2iface"
	"github.com/nitrictech/nitric/cloud/aws/ifaces/resourcegroupstaggingapiiface"
	"github.com/nitrictech/nitric/core/pkg/plugins/resource"
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

var resourceTypeMap = map[resource.ResourceType]AwsResource{
	resource.ResourceType_Api:       AwsResource_Api,
	resource.ResourceType_Websocket: AwsResource_Api,
}

type AwsProvider interface {
	resource.ResourceService

	// GetResources API operation for AWS Provider.
	// Returns requested aws resources for the given resource type
	GetResources(context.Context, AwsResource) (map[string]string, error)
	GetApiGatewayById(context.Context, string) (*apigatewayv2.GetApiOutput, error)
}

// Aws core utility provider
type awsProviderImpl struct {
	stack     string
	cacheLock sync.Mutex
	client    resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	apiClient apigatewayv2iface.ApiGatewayV2API
	cache     map[AwsResource]map[string]string
}

var _ AwsProvider = &awsProviderImpl{}

func (a *awsProviderImpl) Declare(ctx context.Context, req resource.ResourceDeclareRequest) error {
	return nil
}

func (a *awsProviderImpl) Details(ctx context.Context, typ resource.ResourceType, name string) (*resource.DetailsResponse[any], error) {
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

	details := &resource.DetailsResponse[any]{
		Id:       arn,
		Provider: "aws",
	}

	switch rt {
	case AwsResource_Api:
		// split arn to find the apiId
		arnParts := strings.Split(arn, "/")
		apiId := arnParts[len(arnParts)-1]
		// Get api detail
		api, err := a.GetApiGatewayById(ctx, apiId)
		if err != nil {
			return nil, err
		}

		details.Service = "ApiGateway"
		if typ == resource.ResourceType_Api {
			details.Detail = resource.ApiDetails{
				URL: *api.ApiEndpoint,
			}
		} else {
			details.Detail = resource.WebsocketDetails{
				URL: fmt.Sprintf("%s/$default", *api.ApiEndpoint),
			}
		}

		return details, nil
	default:
		return nil, fmt.Errorf("unimplemented resource type")
	}
}

func (a *awsProviderImpl) GetApiGatewayById(ctx context.Context, apiId string) (*apigatewayv2.GetApiOutput, error) {
	return a.apiClient.GetApi(context.TODO(), &apigatewayv2.GetApiInput{
		ApiId: aws.String(apiId),
	})
}

func resourceTypeFromArn(arn string) (resource.ResourceType, error) {
	if !awsArn.IsARN(arn) {
		return "", fmt.Errorf("invalid ARN provided")
	}

	parsedArn, err := awsArn.Parse(arn)
	if err != nil {
		return "", err
	}

	switch parsedArn.Service {
	case "s3":
		return AwsResource_Bucket, nil
	case "sns":
		return AwsResource_Topic, nil
	case "sqs":
		return AwsResource_Queue, nil
	case "apigateway":
		return AwsResource_Api, nil
	case "states":
		return AwsResource_StateMachine, nil
	case "secretsmanager":
		return AwsResource_Secret, nil
	case "events":
		return AwsResource_EventRule, nil
	case "dynamodb":
		return AwsResource_Collection, nil
	default:
		return "", fmt.Errorf("invalid resource type")
	}
}

// populate the resource cache
func (a *awsProviderImpl) populateCache(ctx context.Context) error {
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()
	if a.cache == nil {
		a.cache = make(map[string]map[string]string)

		resourceNameKey := tags.GetResourceNameKey(a.stack)

		tagFilters := []types.TagFilter{{
			Key: aws.String(resourceNameKey),
		}}

		paginator := resourcegroupstaggingapi.NewGetResourcesPaginator(a.client, &resourcegroupstaggingapi.GetResourcesInput{
			TagFilters: tagFilters,
			ResourceTypeFilters: []string{
				AwsResource_Api,
				AwsResource_StateMachine,
				AwsResource_Topic,
				AwsResource_Collection,
				AwsResource_Queue,
				AwsResource_Bucket,
				AwsResource_Secret,
				AwsResource_EventRule,
			},
			ResourcesPerPage: aws.Int32(100),
		})

		for paginator.HasMorePages() {
			out, err := paginator.NextPage(ctx)
			if err != nil {
				fmt.Println("failed to retrieve resources:", err)
				return err
			}

			for _, tm := range out.ResourceTagMappingList {
				for _, t := range tm.Tags {
					if *t.Key == resourceNameKey {
						// Get the resource type from the ARN
						typ, err := resourceTypeFromArn(*tm.ResourceARN)
						if err != nil {
							fmt.Printf("unable to identify resource: %s\n", *tm.ResourceARN)
							break
						}

						if a.cache[typ] == nil {
							a.cache[typ] = map[string]string{}
						}

						a.cache[typ][*t.Value] = *tm.ResourceARN

						break
					}
				}
			}
		}
	}

	return nil
}

func (a *awsProviderImpl) GetResources(ctx context.Context, typ AwsResource) (map[string]string, error) {
	if err := a.populateCache(ctx); err != nil {
		return nil, fmt.Errorf("error populating resource cache")
	}

	return a.cache[typ], nil
}

func New() (AwsProvider, error) {
	awsRegion := utils.GetEnv("AWS_REGION", "us-east-1")
	stack := utils.GetEnv("NITRIC_STACK", "")

	cfg, sessionError := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(awsRegion),
		config.WithRetryMode(aws.RetryModeAdaptive),
		config.WithRetryMaxAttempts(10),
	)
	if sessionError != nil {
		return nil, fmt.Errorf("error creating new AWS session %w", sessionError)
	}

	otelaws.AppendMiddlewares(&cfg.APIOptions)

	apiClient := apigatewayv2.NewFromConfig(cfg)
	client := resourcegroupstaggingapi.NewFromConfig(cfg)

	return &awsProviderImpl{
		stack:     stack,
		client:    client,
		cacheLock: sync.Mutex{},
		apiClient: apiClient,
	}, nil
}
