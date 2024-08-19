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

package resource

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/samber/lo"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsArn "github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewayv2"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/ifaces/apigatewayv2iface"
	"github.com/nitrictech/nitric/cloud/aws/ifaces/resourcegroupstaggingapiiface"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	resourcepb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
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
	AwsResource_Unknown      AwsResource = "unknown"
)

// Map of resources for which 'details' can be requested.
var resourceDetailsTypeMap = map[resourcepb.ResourceType]AwsResource{
	resourcepb.ResourceType_Api:       AwsResource_Api,
	resourcepb.ResourceType_Websocket: AwsResource_Api,
}

type ResolvedResource struct {
	ARN string
}

// Aws core utility provider
type AwsResourceService struct {
	stackID   string
	cacheLock sync.Mutex
	client    resourcegroupstaggingapiiface.ResourceGroupsTaggingAPIAPI
	apiClient apigatewayv2iface.ApiGatewayV2API
	cache     map[AwsResource]map[string]ResolvedResource
}

type AwsResourceResolver interface {
	GetApiGatewayById(context.Context, string) (*apigatewayv2.GetApiOutput, error)
	GetResources(context.Context, AwsResource) (map[string]ResolvedResource, error)
}

var (
	_ AwsResourceResolver        = &AwsResourceService{}
	_ resourcepb.ResourcesServer = &AwsResourceService{}
)

func (a *AwsResourceService) Declare(ctx context.Context, req *resourcepb.ResourceDeclareRequest) (*resourcepb.ResourceDeclareResponse, error) {
	return &resourcepb.ResourceDeclareResponse{}, nil
}

type AWSApiGatewayDetails struct {
	Url string
}

// GetAWSApiGatewayDetails - Get the details for an AWS API Gateway resource related to a Nitric API or Websocket
func (a *AwsResourceService) GetAWSApiGatewayDetails(ctx context.Context, identifier *resourcespb.ResourceIdentifier) (*AWSApiGatewayDetails, error) {
	resourceName := identifier.Name
	resourceType := identifier.Type

	if resourceType != resourcepb.ResourceType_Api && resourceType != resourcepb.ResourceType_Websocket {
		return nil, fmt.Errorf("resource type %s is not an API Gateway", resourceType)
	}

	rt, ok := resourceDetailsTypeMap[resourceType]
	if !ok {
		return nil, fmt.Errorf("unhandled resource type: %s", resourceType)
	}

	// Get resource references (arns) for the resource type
	resources, err := a.GetResources(ctx, rt)
	if err != nil {
		return nil, err
	}

	arn, ok := resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("unable to find resource %s for name: %s", resourceType, resourceName)
	}

	// split arn to find the apiId
	arnParts := strings.Split(arn.ARN, "/")
	apiId := arnParts[len(arnParts)-1]
	// Get api detail
	api, err := a.GetApiGatewayById(ctx, apiId)
	if err != nil {
		return nil, err
	}

	return &AWSApiGatewayDetails{
		Url: *api.ApiEndpoint,
	}, nil
}

func (a *AwsResourceService) GetApiGatewayById(ctx context.Context, apiId string) (*apigatewayv2.GetApiOutput, error) {
	return a.apiClient.GetApi(context.TODO(), &apigatewayv2.GetApiInput{
		ApiId: aws.String(apiId),
	})
}

func resourceTypeFromArn(arn string) (string, error) {
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
		return AwsResource_Unknown, nil
	}
}

// populate the resource cache
func (a *AwsResourceService) populateCache(ctx context.Context) error {
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()
	if a.cache == nil {
		a.cache = make(map[string]map[string]ResolvedResource)

		resourceNameKey := tags.GetResourceNameKey(a.stackID)

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
							return err
						}

						if a.cache[typ] == nil {
							a.cache[typ] = map[string]ResolvedResource{}
						}

						// Check the value doesn't already exist
						if _, ok := a.cache[typ][*t.Value]; ok {
							// Clear the cache to avoid partial data and allow for a retry if a manual fix is applied
							a.cache = nil
							return fmt.Errorf("unable to uniquely identify %s resource, multiple resources found with matching name: %s: ARNs: %s, %s", typ, *t.Value, a.cache[typ][*t.Value].ARN, *tm.ResourceARN)
						}

						a.cache[typ][*t.Value] = ResolvedResource{ARN: *tm.ResourceARN}
						break
					}
				}
			}
		}

		if len(a.cache[AwsResource_Unknown]) > 0 {
			fmt.Printf("resource cache contains unknown/unsupported resources, tagged with the following names: [%s]\n", strings.Join(lo.Keys(a.cache[AwsResource_Unknown]), ", "))
		}
	}

	return nil
}

type ResourceResolutionError struct {
	Msg   string
	Cause error
}

func (e *ResourceResolutionError) Error() string {
	return e.Msg + ": " + e.Cause.Error()
}

func (e *ResourceResolutionError) Unwrap() error {
	return e.Cause
}

func (a *AwsResourceService) GetResources(ctx context.Context, typ AwsResource) (map[string]ResolvedResource, error) {
	if err := a.populateCache(ctx); err != nil {
		return nil, &ResourceResolutionError{
			Msg:   "error populating resource cache",
			Cause: err,
		}
	}

	return a.cache[typ], nil
}

func New() (*AwsResourceService, error) {
	awsRegion := env.AWS_REGION.String()
	stackID := commonenv.NITRIC_STACK_ID.String()

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

	return &AwsResourceService{
		stackID:   stackID,
		client:    client,
		cacheLock: sync.Mutex{},
		apiClient: apiClient,
	}, nil
}
