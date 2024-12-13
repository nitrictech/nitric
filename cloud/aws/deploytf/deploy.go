// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploytf

import (
	"embed"

	"github.com/aws/jsii-runtime-go"
	ecrauth "github.com/cdktf/cdktf-provider-aws-go/aws/v19/dataawsecrauthorizationtoken"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v19/provider"
	dockerprovider "github.com/cdktf/cdktf-provider-docker-go/docker/v11/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/bucket"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/http_proxy"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/keyvalue"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/queue"
	rds "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/rds"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/schedule"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/secret"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/service"
	tfstack "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/stack"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/topic"
	vpc "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/vpc"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/websocket"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAwsTerraformProvider struct {
	*deploy.CommonStackDetails
	Stack tfstack.Stack

	Vpc vpc.Vpc
	Rds rds.Rds

	AwsConfig      *common.AwsConfig
	Apis           map[string]api.Api
	Buckets        map[string]bucket.Bucket
	Topics         map[string]topic.Topic
	HttpProxies    map[string]http_proxy.HttpProxy
	Schedules      map[string]schedule.Schedule
	Services       map[string]service.Service
	Secrets        map[string]secret.Secret
	Queues         map[string]queue.Queue
	KeyValueStores map[string]keyvalue.Keyvalue
	Websockets     map[string]websocket.Websocket

	provider.NitricDefaultOrder
}

var _ provider.NitricTerraformProvider = (*NitricAwsTerraformProvider)(nil)

func (a *NitricAwsTerraformProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	a.AwsConfig, err = common.ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

// embed the modules directory here
//
//go:embed .nitric/modules/**/*
var modules embed.FS

func (a *NitricAwsTerraformProvider) CdkTfModules() ([]provider.ModuleDirectory, error) {
	return []provider.ModuleDirectory{
		{
			ParentDir: ".nitric/modules",
			Modules:   modules,
		},
	}, nil
}

func (a *NitricAwsTerraformProvider) RequiredProviders() map[string]interface{} {
	return map[string]interface{}{}
}

func (a *NitricAwsTerraformProvider) Pre(stack cdktf.TerraformStack, resources []*deploymentspb.Resource) error {
	tfRegion := cdktf.NewTerraformVariable(stack, jsii.String("region"), &cdktf.TerraformVariableConfig{
		Type:        jsii.String("string"),
		Default:     jsii.String(a.Region),
		Description: jsii.String("The AWS region to deploy resources to"),
	})

	awsprovider.NewAwsProvider(stack, jsii.String("aws"), &awsprovider.AwsProviderConfig{
		Region: tfRegion.StringValue(),
	})

	ecrAuthConfig := ecrauth.NewDataAwsEcrAuthorizationToken(stack, jsii.String("ecr_auth"), &ecrauth.DataAwsEcrAuthorizationTokenConfig{})
	dockerprovider.NewDockerProvider(stack, jsii.String("docker"), &dockerprovider.DockerProviderConfig{
		RegistryAuth: &[]*map[string]interface{}{
			{
				"address":  ecrAuthConfig.ProxyEndpoint(),
				"username": ecrAuthConfig.UserName(),
				"password": ecrAuthConfig.Password(),
			},
		},
	})

	a.Stack = tfstack.NewStack(stack, jsii.String("stack"), &tfstack.StackConfig{})

	databases := lo.Filter(resources, func(item *deploymentspb.Resource, idx int) bool {
		return item.Id.Type == resourcespb.ResourceType_SqlDatabase
	})
	// Create a shared database cluster if we have more than one database
	if len(databases) > 0 {
		a.Vpc = vpc.NewVpc(stack, jsii.String("vpc"), &vpc.VpcConfig{})
		a.Rds = rds.NewRds(stack, jsii.String("rds"), &rds.RdsConfig{
			MinCapacity:      jsii.Number(0.5),
			MaxCapacity:      jsii.Number(1),
			VpcId:            a.Vpc.VpcIdOutput(),
			StackId:          a.Stack.StackIdOutput(),
			PrivateSubnetIds: cdktf.Token_AsList(a.Vpc.PrivateSubnetIdsOutput(), &cdktf.EncodingOptions{}),
		})
	}

	return nil
}

func (a *NitricAwsTerraformProvider) Post(stack cdktf.TerraformStack) error {
	// Get a link to this stacks resource group in the AWS console
	// outputs = append(outputs, pulumi.Sprintf("Deployed Resources:\n──────────────"))
	// outputs = append(outputs, pulumi.Sprintf(, a.Region, urlEncodedRgArn))

	// urlEncodedRgArn := a.ResourceGroup.Arn.ApplyT(func(arn string) string {
	// 	// URL encode the ARN
	// 	return url.QueryEscape(arn)
	// })

	cdktf.NewTerraformOutput(stack, jsii.Sprintf("deployed-resources"), &cdktf.TerraformOutputConfig{
		Value: jsii.Sprintf("https://%s.console.aws.amazon.com/resource-groups/group/%s\n", a.Region, *cdktf.Fn_Urlencode(a.Stack.ResourceGroupArnOutput())),
	})

	// Set terraform outputs
	cdktf.NewTerraformOutput(stack, jsii.Sprintf("stack-output"), &cdktf.TerraformOutputConfig{
		Value: a.Stack,
	})

	// loop over all the resources and create outputs for them
	allEndpoints := map[string]*string{}

	for name, api := range a.Apis {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-api-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     api,
		})
		allEndpoints[name] = api.EndpointOutput()
	}

	if len(allEndpoints) > 0 {
		cdktf.NewTerraformOutput(stack, jsii.String("endpoints"), &cdktf.TerraformOutputConfig{
			Value: allEndpoints,
		})
	}

	for name, bucket := range a.Buckets {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-bucket-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     bucket,
		})
	}

	for name, topic := range a.Topics {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-topic-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     topic,
		})
	}

	for name, schedule := range a.Schedules {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-schedule-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     schedule,
		})
	}

	for name, service := range a.Services {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-service-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     service,
		})
	}

	for name, secret := range a.Secrets {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-secret-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     secret,
		})
	}

	for name, queue := range a.Queues {
		cdktf.NewTerraformOutput(stack, jsii.Sprintf("%s-queue-output", name), &cdktf.TerraformOutputConfig{
			Sensitive: jsii.Bool(true),
			Value:     queue,
		})
	}

	// Give all the Services access to the resource index
	accessRoleNames := []string{}
	for _, service := range a.Services {
		accessRoleNames = append(accessRoleNames, *service.RoleNameOutput())
	}

	return a.ResourcesStore(stack, accessRoleNames)
}

// // Post - Called after all resources have been created, before the Pulumi Context is concluded
// Post(stack cdktf.TerraformStack) error

func NewNitricAwsProvider() *NitricAwsTerraformProvider {
	return &NitricAwsTerraformProvider{
		Apis:           make(map[string]api.Api),
		Buckets:        make(map[string]bucket.Bucket),
		Services:       make(map[string]service.Service),
		Topics:         make(map[string]topic.Topic),
		Schedules:      make(map[string]schedule.Schedule),
		Secrets:        make(map[string]secret.Secret),
		Queues:         make(map[string]queue.Queue),
		KeyValueStores: make(map[string]keyvalue.Keyvalue),
		Websockets:     make(map[string]websocket.Websocket),
	}
}
