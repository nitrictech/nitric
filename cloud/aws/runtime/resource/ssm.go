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
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/runtime/env"
	commonenv "github.com/nitrictech/nitric/cloud/common/runtime/env"
	resourcepb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
)

// Aws core utility provider
type AwsSSMResourceResolver struct {
	stackID   string
	cacheLock sync.Mutex
	// Bot, give me the ssm  client here
	client *ssm.Client

	cache *common.ResourceIndex
}

var _ AwsResourceResolver = &AwsSSMResourceResolver{}

// GetAWSApiGatewayDetails - Get the details for an AWS API Gateway resource related to a Nitric API or Websocket
func (a *AwsSSMResourceResolver) GetAWSApiGatewayDetails(ctx context.Context, identifier *resourcespb.ResourceIdentifier) (*AWSApiGatewayDetails, error) {
	resourceName := identifier.Name
	resourceType := identifier.Type

	endpoint := ""
	switch resourceType {
	case resourcepb.ResourceType_Api:
		endpoint = a.cache.Apis[resourceName].Endpoint
	case resourcepb.ResourceType_Websocket:
		endpoint = a.cache.Websockets[resourceName].Endpoint
	}

	return &AWSApiGatewayDetails{
		Url: endpoint,
	}, nil
}

type ApiGatewayDetails struct {
	Name        string
	Type        string
	ApiEndpoint string
}

func (a *AwsSSMResourceResolver) GetApiGatewayById(ctx context.Context, apiId string) (*ApiGatewayDetails, error) {
	err := a.populateCache(ctx)
	if err != nil {
		return nil, err
	}
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()

	for name, api := range a.cache.Apis {
		if strings.HasSuffix(api.Arn, apiId) {
			return &ApiGatewayDetails{
				Name: name,
				Type: "api",
			}, nil
		}
	}

	for name, api := range a.cache.HttpProxies {
		if strings.HasSuffix(api.Arn, apiId) {
			return &ApiGatewayDetails{
				Name: name,
				Type: "http-proxy",
			}, nil
		}
	}

	for name, api := range a.cache.Websockets {
		if strings.HasSuffix(api.Arn, apiId) {
			return &ApiGatewayDetails{
				Name: name,
				Type: "websocket",
			}, nil
		}
	}

	return nil, fmt.Errorf("api gateway not found")
}

// populate the resource cache
func (a *AwsSSMResourceResolver) populateCache(ctx context.Context) error {
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()

	if a.cache != nil {
		return nil
	}

	response, err := a.client.GetParameter(ctx, &ssm.GetParameterInput{
		Name: aws.String(fmt.Sprintf("/nitric/%s/resource-index", a.stackID)),
	})
	if err != nil {
		return err
	}

	val := response.Parameter.Value
	if val == nil {
		return fmt.Errorf("resource index not found")
	}

	var index common.ResourceIndex

	err = json.Unmarshal([]byte(*val), &index)
	if err != nil {
		return fmt.Errorf("error unmarshalling resource index: %w", err)
	}

	a.cache = &index

	return nil
}

func (a *AwsSSMResourceResolver) GetResources(ctx context.Context, typ AwsResource) (map[string]ResolvedResource, error) {
	if err := a.populateCache(ctx); err != nil {
		return nil, &ResourceResolutionError{
			Msg:   "error populating resource cache",
			Cause: err,
		}
	}

	resolvedResources := map[string]ResolvedResource{}

	switch typ {
	case AwsResource_Api:
		for name, api := range a.cache.Apis {
			resolvedResources[name] = ResolvedResource{
				ARN: api.Arn,
			}
		}
	case AwsResource_StateMachine:
		for name, topic := range a.cache.Topics {
			resolvedResources[name] = ResolvedResource{
				ARN: topic.StateMachineArn,
			}
		}
	case AwsResource_Topic:
		for name, topic := range a.cache.Topics {
			resolvedResources[name] = ResolvedResource{
				ARN: topic.Arn,
			}
		}
	case AwsResource_Collection:
		for name, collection := range a.cache.KvStores {
			resolvedResources[name] = ResolvedResource{
				ARN: collection,
			}
		}
	case AwsResource_Bucket:
		for name, bucket := range a.cache.Buckets {
			resolvedResources[name] = ResolvedResource{
				ARN: bucket,
			}
		}
	case AwsResource_Secret:
		for name, secret := range a.cache.Secrets {
			resolvedResources[name] = ResolvedResource{
				ARN: secret,
			}
		}
	}

	return resolvedResources, nil
}

func NewSSMResourceResolver() (*AwsSSMResourceResolver, error) {
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

	client := ssm.NewFromConfig(cfg)

	return &AwsSSMResourceResolver{
		stackID:   stackID,
		client:    client,
		cacheLock: sync.Mutex{},
	}, nil
}
