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
	"io/fs"

	"github.com/aws/jsii-runtime-go"
	awsprovider "github.com/cdktf/cdktf-provider-aws-go/aws/v10/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/bucket"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/keyvalue"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/queue"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/schedule"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/secret"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/service"
	tfstack "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/stack"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/topic"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/websocket"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAwsTerraformProvider struct {
	*deploy.CommonStackDetails
	Stack tfstack.Stack

	AwsConfig      *common.AwsConfig
	Apis           map[string]api.Api
	Buckets        map[string]bucket.Bucket
	Topics         map[string]topic.Topic
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
		return status.Errorf(codes.InvalidArgument, err.Error())
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

func (a *NitricAwsTerraformProvider) CdkTfModules() (string, fs.FS, error) {
	return ".nitric/modules", modules, nil
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

	a.Stack = tfstack.NewStack(stack, jsii.String("stack"), &tfstack.StackConfig{})

	return nil
}

func (a *NitricAwsTerraformProvider) Post(stack cdktf.TerraformStack) error {
	return nil
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
