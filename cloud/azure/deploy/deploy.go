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

	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	commonresources "github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/core/pkg/logger"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pkg/errors"
	apimanagement "github.com/pulumi/pulumi-azure-native-sdk/apimanagement/v20201201"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi-azure-native-sdk/eventgrid"
	"github.com/pulumi/pulumi-azure-native-sdk/keyvault"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ApiResources struct {
	ApiManagementService *apimanagement.ApiManagementService
	Api                  *apimanagement.Api
}

type NitricAzurePulumiProvider struct {
	*deploy.CommonStackDetails

	StackId   string
	resources []*pulumix.NitricPulumiResource[any]

	AzureConfig *AzureConfig

	ClientConfig *authorization.GetClientConfigResult

	ResourceGroup  *resources.ResourceGroup
	KeyVault       *keyvault.Vault
	StorageAccount *storage.StorageAccount

	ContainerEnv *ContainerEnv

	Apis        map[string]ApiResources
	HttpProxies map[string]ApiResources
	Buckets     map[string]*storage.BlobContainer

	Queues map[string]*storage.Queue

	Principals map[resourcespb.ResourceType]map[string]*ServicePrincipal

	ContainerApps map[string]*ContainerApp
	Topics        map[string]*eventgrid.Topic

	KeyValueStores map[string]*storage.Table

	Roles *Roles
	provider.NitricDefaultOrder
}

var _ provider.NitricPulumiProvider = (*NitricAzurePulumiProvider)(nil)

const (
	pulumiAzureNativeVersion = "1.95.0"
	pulumiAzureVersion       = "5.52.0"
)

func (a *NitricAzurePulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"azure-native:location": auto.ConfigValue{Value: a.Region},
		"azure:location":        auto.ConfigValue{Value: a.Region},
		"azure-native:version":  auto.ConfigValue{Value: pulumiAzureNativeVersion},
		"azure:version":         auto.ConfigValue{Value: pulumiAzureVersion},
		"docker:version":        auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricAzurePulumiProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	a.AzureConfig, err = ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

func createKeyVault(ctx *pulumi.Context, group *resources.ResourceGroup, tenantId string, tags map[string]string) (*keyvault.Vault, error) {
	// Create a stack level keyvault if secrets are enabled
	// At the moment secrets have no config level setting
	vaultName := ResourceName(ctx, "", KeyVaultRT)

	keyVault, err := keyvault.NewVault(ctx, vaultName, &keyvault.VaultArgs{
		Location:          group.Location,
		ResourceGroupName: group.Name,
		Properties: &keyvault.VaultPropertiesArgs{
			EnableSoftDelete:        pulumi.Bool(false),
			EnableRbacAuthorization: pulumi.Bool(true),
			Sku: &keyvault.SkuArgs{
				Family: pulumi.String("A"),
				Name:   keyvault.SkuNameStandard,
			},
			TenantId: pulumi.String(tenantId),
		},
		Tags: pulumi.ToStringMap(tags),
	})
	if err != nil {
		return nil, err
	}

	return keyVault, nil
}

func createStorageAccount(ctx *pulumi.Context, group *resources.ResourceGroup, tags map[string]string) (*storage.StorageAccount, error) {
	accName := ResourceName(ctx, "", StorageAccountRT)
	storageAccount, err := storage.NewStorageAccount(ctx, accName, &storage.StorageAccountArgs{
		AccessTier:        storage.AccessTierHot,
		ResourceGroupName: group.Name,
		Kind:              pulumi.String("StorageV2"),
		Sku: storage.SkuArgs{
			Name: pulumi.String(storage.SkuName_Standard_LRS),
		},
		Tags: pulumi.ToStringMap(tags),
	})
	if err != nil {
		return nil, err
	}

	return storageAccount, nil
}

func hasResourceType(resources []*pulumix.NitricPulumiResource[any], resourceType resourcespb.ResourceType) bool {
	for _, r := range resources {
		if r.Id.GetType() == resourceType {
			return true
		}
	}

	return false
}

func (a *NitricAzurePulumiProvider) Pre(ctx *pulumi.Context, nitricResources []*pulumix.NitricPulumiResource[any]) error {
	a.resources = nitricResources

	// make our random stackId
	stackRandId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-stack-name", ctx.Stack()), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Upper:   pulumi.Bool(false),
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

	a.StackId = <-stackIdChan

	a.ClientConfig, err = authorization.GetClientConfig(ctx)
	if err != nil {
		return err
	}

	a.ResourceGroup, err = resources.NewResourceGroup(ctx, ResourceName(ctx, "", ResourceGroupRT), &resources.ResourceGroupArgs{
		Location: pulumi.String(a.Region),
		Tags:     pulumi.ToStringMap(tags.Tags(a.StackId, ctx.Stack(), commonresources.Stack)),
	})
	if err != nil {
		return errors.WithMessage(err, "resource group create")
	}

	// Create a key vault if secrets are required.
	// Unlike AWS and GCP which have centralized secrets management, Azure allows for multiple key vaults.
	// This means we need to create a keyvault for each stack.
	if hasResourceType(nitricResources, resourcespb.ResourceType_Secret) {
		logger.Info("Stack declares one or more secrets, creating stack level Azure Key Vault")
		a.KeyVault, err = createKeyVault(ctx, a.ResourceGroup, a.ClientConfig.TenantId, tags.Tags(a.StackId, ctx.Stack(), commonresources.Stack))
		if err != nil {
			return errors.WithMessage(err, "keyvault create")
		}
	}

	hasBuckets := hasResourceType(nitricResources, resourcespb.ResourceType_Bucket)
	hasKvStores := hasResourceType(nitricResources, resourcespb.ResourceType_KeyValueStore)
	hasQueues := hasResourceType(nitricResources, resourcespb.ResourceType_Queue)

	// Create a storage account if buckets, kv stores or queues are required.
	// Unlike AWS and GCP which have centralized storage management, Azure allows for multiple storage accounts.
	// This means we need to create a storage account for each stack, before buckets can be created.
	if hasBuckets || hasKvStores || hasQueues {
		logger.Info("Stack declares bucket(s), key/value store(s) or queue(s), creating stack level Azure Storage Account")
		a.StorageAccount, err = createStorageAccount(ctx, a.ResourceGroup, tags.Tags(a.StackId, ctx.Stack(), commonresources.Stack))
		if err != nil {
			return errors.WithMessage(err, "storage account create")
		}
	}

	a.ContainerEnv, err = a.newContainerEnv(ctx, a.StackId, map[string]string{})
	if err != nil {
		return err
	}

	// Greedily create all the roles for consistency. Could be reduced to required roles only in future.
	a.Roles, err = CreateRoles(ctx, a.StackId, a.ClientConfig.SubscriptionId, a.ResourceGroup.Name)
	if err != nil {
		return err
	}

	return nil
}

func (a *NitricAzurePulumiProvider) Post(ctx *pulumi.Context) error {
	return nil
}

func (a *NitricAzurePulumiProvider) Result(ctx *pulumi.Context) (pulumi.StringOutput, error) {
	outputs := []interface{}{}

	// Add APIs outputs
	if len(a.Apis) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.Apis {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", apiName, api.ApiManagementService.GatewayUrl))
		}
	}

	// Add HTTP Proxy outputs
	if len(a.HttpProxies) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("HTTP Proxies:\n──────────────"))
		for proxyName, proxy := range a.HttpProxies {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", proxyName, proxy.ApiManagementService.GatewayUrl))
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

func NewNitricAzurePulumiProvider() *NitricAzurePulumiProvider {
	principalsMap := map[resourcespb.ResourceType]map[string]*ServicePrincipal{}

	principalsMap[resourcespb.ResourceType_Service] = map[string]*ServicePrincipal{}

	return &NitricAzurePulumiProvider{
		Apis:           map[string]ApiResources{},
		HttpProxies:    map[string]ApiResources{},
		Buckets:        make(map[string]*storage.BlobContainer),
		Queues:         make(map[string]*storage.Queue),
		ContainerApps:  map[string]*ContainerApp{},
		Topics:         map[string]*eventgrid.Topic{},
		Principals:     principalsMap,
		KeyValueStores: map[string]*storage.Table{},
	}
}
