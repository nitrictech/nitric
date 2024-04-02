// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"fmt"
	"strings"

	_ "embed"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/resourcegroups"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/secretsmanager"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAwsPulumiProvider struct {
	stackId     string
	projectName string
	stackName   string

	fullStackName string

	config *AwsConfig
	region string

	ecrAuthToken        *ecr.GetAuthorizationTokenResult
	lambdas             map[string]*lambda.Function
	lambdaRoles         map[string]*iam.Role
	httpProxies         map[string]*apigatewayv2.Api
	apis                map[string]*apigatewayv2.Api
	secrets             map[string]*secretsmanager.Secret
	buckets             map[string]*s3.Bucket
	bucketNotifications map[string]*s3.BucketNotification
	topics              map[string]*topic
	queues              map[string]*sqs.Queue
	websockets          map[string]*apigatewayv2.Api
	keyValueStores      map[string]*dynamodb.Table

	provider.NitricDefaultOrder

	resourceTaggingClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	lambdaClient          lambdaiface.LambdaAPI
}

// Embeds the runtime directly into the deploytime binary
// This way the versions will always match as they're always built and versioned together (as a single artifact)
// This should also help with docker build speeds as the runtime has already been "downloaded"
//
//go:embed runtime-aws
var runtime []byte

var _ provider.NitricPulumiProvider = (*NitricAwsPulumiProvider)(nil)

const pulumiAwsVersion = "6.6.0"

func (a *NitricAwsPulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"aws:region":     auto.ConfigValue{Value: a.region},
		"aws:version":    auto.ConfigValue{Value: pulumiAwsVersion},
		"docker:version": auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricAwsPulumiProvider) Init(attributes map[string]interface{}) error {
	var err error

	region, ok := attributes["region"].(string)
	if !ok {
		return fmt.Errorf("Missing region attribute")
	}

	a.region = region

	a.config, err = ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	var isString bool

	iProject, hasProject := attributes["project"]
	a.projectName, isString = iProject.(string)
	if !hasProject || !isString || a.projectName == "" {
		// need a valid project name
		return fmt.Errorf("project is not set or invalid")
	}

	iStack, hasStack := attributes["stack"]
	a.stackName, isString = iStack.(string)
	if !hasStack || !isString || a.stackName == "" {
		// need a valid stack name
		return fmt.Errorf("stack is not set or invalid")
	}

	// Backwards compatible stack name
	// The existing providers in the CLI
	// Use the combined project and stack name
	a.fullStackName = fmt.Sprintf("%s-%s", a.projectName, a.stackName)

	return nil
}

func (a *NitricAwsPulumiProvider) Pre(ctx *pulumi.Context, resources []*deploymentspb.Resource) error {
	// make our random stackId
	stackRandId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-stack-name", ctx.Stack()), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"stack-name": ctx.Stack(),
		}),
	})
	if err != nil {
		return err
	}

	stackIdChan := make(chan string)
	pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.Result).ApplyT(func(id string) string {
		stackIdChan <- id
		return id
	})

	a.stackId = <-stackIdChan

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	a.resourceTaggingClient = resourcegroupstaggingapi.New(sess)

	a.lambdaClient = lambdaClient.New(sess, &aws.Config{Region: aws.String(a.region)})

	a.ecrAuthToken, err = ecr.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenArgs{})
	if err != nil {
		return err
	}

	// Create AWS Resource groups with our tags
	_, err = resourcegroups.NewGroup(ctx, "stack-resource-group", &resourcegroups.GroupArgs{
		Name:        pulumi.String(a.fullStackName),
		Description: pulumi.Sprintf("All deployed resources for the %s nitric stack", a.fullStackName),
		ResourceQuery: &resourcegroups.GroupResourceQueryArgs{
			Query: pulumi.Sprintf(`{
				"ResourceTypeFilters":["AWS::AllSupported"],
				"TagFilters":[{"Key":"%s"}]
			}`, tags.GetResourceNameKey(a.stackId)),
		},
	})

	return err
}

func (a *NitricAwsPulumiProvider) Post(ctx *pulumi.Context) error {
	return nil
}

func (a *NitricAwsPulumiProvider) Result(ctx *pulumi.Context) (pulumi.StringOutput, error) {
	outputs := []interface{}{}

	// Add APIs outputs
	if len(a.apis) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.apis {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", apiName, api.ApiEndpoint))
		}
	}

	// Add HTTP Proxy outputs
	if len(a.httpProxies) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("HTTP Proxies:\n──────────────"))
		for proxyName, proxy := range a.httpProxies {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", proxyName, proxy.ApiEndpoint))
		}
	}

	// Add Websocket outputs
	if len(a.websockets) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("Websockets:\n──────────────"))
		for wsName, ws := range a.websockets {
			outputs = append(outputs, pulumi.Sprintf("%s: %s/%s", wsName, ws.ApiEndpoint, common.DefaultWsStageName))
		}
	}

	output, ok := pulumi.All(outputs...).ApplyT(func(deets []interface{}) string {
		stringyOutputs := make([]string, len(deets))
		for i, d := range deets {
			stringyOutputs[i] = d.(string)
		}

		return strings.Join(stringyOutputs, "\n")
	}).(pulumi.StringOutput)

	if !ok {
		return pulumi.StringOutput{}, fmt.Errorf("Failed to generate pulumi output")
	}

	return output, nil
}

func NewNitricAwsProvider() *NitricAwsPulumiProvider {
	return &NitricAwsPulumiProvider{
		lambdas:             make(map[string]*lambda.Function),
		lambdaRoles:         make(map[string]*iam.Role),
		apis:                make(map[string]*apigatewayv2.Api),
		httpProxies:         make(map[string]*apigatewayv2.Api),
		secrets:             make(map[string]*secretsmanager.Secret),
		buckets:             make(map[string]*s3.Bucket),
		bucketNotifications: make(map[string]*s3.BucketNotification),
		websockets:          make(map[string]*apigatewayv2.Api),
		topics:              make(map[string]*topic),
		queues:              make(map[string]*sqs.Queue),
		keyValueStores:      make(map[string]*dynamodb.Table),
	}
}
