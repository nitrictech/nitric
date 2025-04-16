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
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/http_proxy"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/keyvalue"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/queue"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/roles"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/service"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/sql"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/stack"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/topic"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/website"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
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
	Proxies   map[string]http_proxy.HttpProxy
	Buckets   map[string]bucket.Bucket
	Services  map[string]service.Service
	Queues    map[string]queue.Queue
	KvStores  map[string]keyvalue.Keyvalue
	Topics    map[string]topic.Topic
	Databases map[string]sql.Sql
	Websites  map[string]website.Website

	EnableWebsites bool

	SubscriptionId string
	AzureConfig    *common.AzureConfig

	dockerProvider dockerprovider.DockerProvider

	provider.NitricDefaultOrder
}

var _ provider.NitricTerraformProvider = (*NitricAzureTerraformProvider)(nil)

func (a *NitricAzureTerraformProvider) GetGlobalTags() *map[string]*string {
	commonTags := a.CommonStackDetails.Tags

	// merge resourceTags and common tags
	combinedTags := make(map[string]*string)
	for k, v := range commonTags {
		combinedTags[k] = jsii.String(v)
	}

	return &combinedTags
}

func (a *NitricAzureTerraformProvider) GetTags(stackID string, resourceName string, resourceType resources.ResourceType) *map[string]*string {
	resourceTags := tags.Tags(stackID, resourceName, resourceType)
	// Augment with common tags
	commonTags := a.CommonStackDetails.Tags

	// merge resourceTags and common tags
	combinedTags := make(map[string]*string)
	for k, v := range commonTags {
		combinedTags[k] = jsii.String(v)
	}
	for k, v := range resourceTags {
		combinedTags[k] = jsii.String(v)
	}

	return &combinedTags
}

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

	// If resources contains queues/buckets/keyvalue then we need to enable storage
	_, enableStorage := lo.Find(resources, func(item *deploymentspb.Resource) bool {
		return item.Id.GetType() == resourcespb.ResourceType_Bucket ||
			item.Id.GetType() == resourcespb.ResourceType_Queue ||
			item.Id.GetType() == resourcespb.ResourceType_KeyValueStore
	})

	_, enableDatabase := lo.Find(resources, func(item *deploymentspb.Resource) bool {
		return item.Id.GetType() == resourcespb.ResourceType_SqlDatabase
	})

	_, enableKeyvault := lo.Find(resources, func(item *deploymentspb.Resource) bool {
		return item.Id.GetType() == resourcespb.ResourceType_Secret
	})

	var resourceGroupImport *string = nil
	if a.AzureConfig.ResourceGroup != "" {
		resourceGroupImport = jsii.String(a.AzureConfig.ResourceGroup)
	}

	enablePrivateEndpoints := a.AzureConfig.Network.PrivateEndpoints
	var subnetId *string = nil
	if a.AzureConfig.Network.SubnetId != "" {
		subnetId = jsii.String(a.AzureConfig.Network.SubnetId)
	}

	var vnetName *string = nil
	if a.AzureConfig.Network.VnetId != "" {
		vnetName = jsii.String(a.AzureConfig.Network.VnetId)
	}

	createDnsZones := a.AzureConfig.Network.CreateDnsZones

	// Deploy the stack - this deploys all pre-requisite environment level resources to support the nitric stack
	a.Stack = stack.NewStack(tfstack, jsii.String("stack"), &stack.StackConfig{
		SubscriptionId:    jsii.String(a.SubscriptionId),
		EnableStorage:     jsii.Bool(enableStorage),
		EnableKeyvault:    jsii.Bool(enableKeyvault),
		EnableDatabase:    jsii.Bool(enableDatabase),
		Location:          jsii.String(a.Region),
		StackName:         jsii.String(a.StackName),
		Tags:              a.GetGlobalTags(),
		ResourceGroupName: resourceGroupImport,
		PrivateEndpoints:  jsii.Bool(enablePrivateEndpoints),
		SubnetId:          subnetId,
		VnetId:            vnetName,
		CreateDnsZones:    jsii.Bool(createDnsZones),
	})

	a.Roles = roles.NewRoles(tfstack, jsii.String("roles"), &roles.RolesConfig{
		ResourceGroupName: a.Stack.ResourceGroupNameOutput(),
		StackName:         a.Stack.StackNameOutput(),
	})

	auths := []dockerprovider.DockerProviderRegistryAuth{
		{
			Address:    a.Stack.RegistryLoginServerOutput(),
			Username:   a.Stack.RegistryUsernameOutput(),
			Password:   a.Stack.RegistryPasswordOutput(),
			ConfigFile: jsii.String(""), // This is unset so username/password takes precedence over the config file
		},
	}

	a.dockerProvider = dockerprovider.NewDockerProvider(tfstack, jsii.String("docker"), &dockerprovider.DockerProviderConfig{
		RegistryAuth: auths,
	})

	return nil
}

func (a *NitricAzureTerraformProvider) Post(stack cdktf.TerraformStack) error {
	// Create a CDN for the stack if we have a website
	if len(a.Websites) > 0 {
		return a.NewCdn(stack)
	}

	return nil
}

func NewNitricAzureProvider() *NitricAzureTerraformProvider {
	return &NitricAzureTerraformProvider{
		Apis:      make(map[string]api.Api),
		Buckets:   make(map[string]bucket.Bucket),
		Services:  make(map[string]service.Service),
		Proxies:   make(map[string]http_proxy.HttpProxy),
		Queues:    make(map[string]queue.Queue),
		Topics:    make(map[string]topic.Topic),
		KvStores:  make(map[string]keyvalue.Keyvalue),
		Databases: make(map[string]sql.Sql),
		Websites:  make(map[string]website.Website),
	}
}
