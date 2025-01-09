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
	dockerprovider "github.com/cdktf/cdktf-provider-docker-go/docker/v11/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/common"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/keyvalue"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/queue"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/roles"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/service"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/stack"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/topic"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAzureTerraformProvider struct {
	*deploy.CommonStackDetails

	Stack stack.Stack
	Roles roles.Roles

	Apis     map[string]api.Api
	Services map[string]service.Service
	Queues   map[string]queue.Queue
	KvStores map[string]keyvalue.Keyvalue
	Topics   map[string]topic.Topic

	AzureConfig *common.AzureConfig

	provider.NitricDefaultOrder
}

var _ provider.NitricTerraformProvider = (*NitricAzureTerraformProvider)(nil)

func (a *NitricAzureTerraformProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	a.AzureConfig, err = common.ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

// embed the modules directory here
//
//go:embed .nitric/modules/**/*
var modules embed.FS

func (a *NitricAzureTerraformProvider) CdkTfModules() ([]provider.ModuleDirectory, error) {
	return []provider.ModuleDirectory{
		{
			ParentDir: ".nitric/modules",
			Modules:   modules,
		},
	}, nil
}

func (a *NitricAzureTerraformProvider) RequiredProviders() map[string]interface{} {
	return map[string]interface{}{}
}

func (a *NitricAzureTerraformProvider) Pre(tfstack cdktf.TerraformStack, resources []*deploymentspb.Resource) error {

	// If resources contains queues/buckets then we need to enable storage
	_, enableStorage := lo.Find(resources, func(item *deploymentspb.Resource) bool {
		return item.Id.GetType() == resourcespb.ResourceType_Bucket || item.Id.GetType() == resourcespb.ResourceType_Queue
	})

	_, enableDatabase := lo.Find(resources, func(item *deploymentspb.Resource) bool {
		return item.Id.GetType() == resourcespb.ResourceType_SqlDatabase
	})

	_, enableKeyvault := lo.Find(resources, func(item *deploymentspb.Resource) bool {
		return item.Id.GetType() == resourcespb.ResourceType_Secret
	})

	// Deploy the stack - this deploys all pre-requisite environment level resources to support the nitric stack
	a.Stack = stack.NewStack(tfstack, jsii.String("stack"), &stack.StackConfig{
		EnableStorage:  jsii.Bool(enableStorage),
		EnableDatabase: jsii.Bool(enableDatabase),
		EnableKeyvault: jsii.Bool(enableKeyvault),
	})

	a.Roles = roles.NewRoles(tfstack, jsii.String("roles"), &roles.RolesConfig{
		ResourceGroupName: a.Stack.ResourceGroupNameOutput(),
	})

	dockerprovider.NewDockerProvider(tfstack, jsii.String("docker"), &dockerprovider.DockerProviderConfig{
		RegistryAuth: &[]*map[string]interface{}{
			{
				"address":  a.Stack.RegistryLoginServerOutput(),
				"username": a.Stack.RegistryUsernameOutput(),
				"password": a.Stack.RegistryPasswordOutput(),
			},
		},
	})

	return nil
}

func (a *NitricAzureTerraformProvider) Post(stack cdktf.TerraformStack) error {
	return nil
}

func NewNitricAzureProvider() *NitricAzureTerraformProvider {
	return &NitricAzureTerraformProvider{
		Apis:     make(map[string]api.Api),
		Services: make(map[string]service.Service),
		Queues:   make(map[string]queue.Queue),
		Topics:   make(map[string]topic.Topic),
		KvStores: make(map[string]keyvalue.Keyvalue),
	}
}
