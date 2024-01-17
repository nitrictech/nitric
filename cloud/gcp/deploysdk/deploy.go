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

	_ "embed"

	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/apigateway"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricGcpPulumiProvider struct {
	stackId       string
	projectName   string
	stackName     string
	fullStackName string

	config *GcpConfig
	region string

	apis             map[string]*apigateway.Api
	apiGateways      map[string]*apigateway.Gateway
	cloudRunServices map[string]*cloudrun.Service
	// ecrAuthToken        *ecr.GetAuthorizationTokenResult
	// lambdas             map[string]*lambda.Function
	// lambdaRoles         map[string]*iam.Role
	// httpProxies         map[string]*apigatewayv2.Api
	// apis                map[string]*apigatewayv2.Api
	// secrets             map[string]*secretsmanager.Secret
	// buckets             map[string]*s3.Bucket
	// bucketNotifications map[string]*s3.BucketNotification
	// topics              map[string]*topic
	// collections         map[string]*dynamodb.Table
	// websockets          map[string]*apigatewayv2.Api

	provider.NitricDefaultOrder

	// resourceTaggingClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	// lambdaClient          lambdaiface.LambdaAPI
}

// Embeds the runtime directly into the deploytime binary
// This way the versions will always match as they're always built and versioned together (as a single artifact)
// This should also help with docker build speeds as the runtime has already been "downloaded"
//
//go:embed runtime-aws
var runtime []byte

var _ provider.NitricPulumiProvider = (*NitricGcpPulumiProvider)(nil)

const pulumiGcpVersion = "6.6.0"

func (a *NitricGcpPulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"gcp:location":   auto.ConfigValue{Value: a.region},
		"gcp:version":    auto.ConfigValue{Value: pulumiAwsVersion},
		"docker:version": auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricGcpPulumiProvider) Init(attributes map[string]interface{}) error {
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

func (a *NitricGcpPulumiProvider) Post(ctx *pulumi.Context) error {
	return nil
}

func (a *NitricGcpPulumiProvider) Pre(ctx *pulumi.Context) error {
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

	return nil
}

func NewNitricGcpProvider() *NitricGcpPulumiProvider {

	return &NitricGcpPulumiProvider{}
}
