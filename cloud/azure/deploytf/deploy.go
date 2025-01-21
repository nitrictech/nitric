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
	"fmt"

	"github.com/aws/jsii-runtime-go"
	azadprovider "github.com/cdktf/cdktf-provider-azuread-go/azuread/v13/provider"
	azprovider "github.com/cdktf/cdktf-provider-azurerm-go/azurerm/v13/provider"
	dockerprovider "github.com/cdktf/cdktf-provider-docker-go/docker/v11/provider"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/common"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/api"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/bucket"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/keyvalue"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/queue"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/roles"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/service"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/sql"
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

	Apis      map[string]api.Api
	Buckets   map[string]bucket.Bucket
	Services  map[string]service.Service
	Queues    map[string]queue.Queue
	KvStores  map[string]keyvalue.Keyvalue
	Topics    map[string]topic.Topic
	Databases map[string]sql.Sql

	SubscriptionId string
	AzureConfig    *common.AzureConfig

	dockerProvider dockerprovider.DockerProvider

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

	var ok bool
	a.SubscriptionId, ok = attributes["subscription-id"].(string)
	if !ok {
		return fmt.Errorf("missing subscription-id attribute, required for configuring azure provider")
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
	// Create an azure terraform provider
	// provider

	azprovider.NewAzurermProvider(tfstack, jsii.String("azurerm"), &azprovider.AzurermProviderConfig{
		Features:       &struct{}{},
		SubscriptionId: jsii.String(a.SubscriptionId),
	})

	azadprovider.NewAzureadProvider(tfstack, jsii.String("azuread"), &azadprovider.AzureadProviderConfig{})

	// azprovider.NewAzureProvider

	// azadprovider.NewAzureProvider

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
		Location:       jsii.String(a.Region),
		StackName:      jsii.String(a.StackName),
	})

	a.Roles = roles.NewRoles(tfstack, jsii.String("roles"), &roles.RolesConfig{
		ResourceGroupName: a.Stack.ResourceGroupNameOutput(),
	})

	a.dockerProvider = dockerprovider.NewDockerProvider(tfstack, jsii.String("docker"), &dockerprovider.DockerProviderConfig{
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
		Apis:      make(map[string]api.Api),
		Buckets:   make(map[string]bucket.Bucket),
		Services:  make(map[string]service.Service),
		Queues:    make(map[string]queue.Queue),
		Topics:    make(map[string]topic.Topic),
		KvStores:  make(map[string]keyvalue.Keyvalue),
		Databases: make(map[string]sql.Sql),
	}
}
