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

	"github.com/nitrictech/nitric/cloud/azure/common"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	commonresources "github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/core/pkg/logger"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pkg/errors"
	apimanagement "github.com/pulumi/pulumi-azure-native-sdk/apimanagement/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	cdn "github.com/pulumi/pulumi-azure-native-sdk/cdn/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/containerinstance/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/dbforpostgresql/v2"
	eventgrid "github.com/pulumi/pulumi-azure-native-sdk/eventgrid/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/keyvault"
	"github.com/pulumi/pulumi-azure-native-sdk/network/v2"
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

	WebsiteStorageAccounts map[string]*storage.StorageAccount
	WebsiteContainers      map[string]*storage.StorageAccountStaticWebsite
	Endpoint               *cdn.AFDEndpoint
	websiteFileMd5Outputs  pulumi.Array

	AzureConfig *common.AzureConfig

	ClientConfig *authorization.GetClientConfigResult

	ResourceGroup  *resources.ResourceGroup
	KeyVault       *keyvault.Vault
	StorageAccount *storage.StorageAccount

	ContainerEnv *ContainerEnv

	Apis                  map[string]ApiResources
	HttpProxies           map[string]ApiResources
	Buckets               map[string]*storage.BlobContainer
	UserDelegationKeyRole *authorization.RoleDefinition

	Queues map[string]*storage.Queue

	Principals map[resourcespb.ResourceType]map[string]*ServicePrincipal

	ContainerApps map[string]*ContainerApp
	Topics        map[string]*eventgrid.Topic

	KeyValueStores map[string]*storage.Table

	SqlMigrations    map[string]*containerinstance.ContainerGroup
	DatabaseServer   *dbforpostgresql.Server
	DbMasterPassword *random.RandomPassword
	VirtualNetwork   *network.VirtualNetwork

	DatabaseSubnet       *network.Subnet
	InfrastructureSubnet *network.Subnet
	ContainerGroupSubnet *network.Subnet

	Roles *Roles
	provider.NitricDefaultOrder
}

var _ provider.NitricPulumiProvider = (*NitricAzurePulumiProvider)(nil)

const (
	pulumiAzureNativeVersion = "2.40.0"
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
		return status.Error(codes.InvalidArgument, err.Error())
	}

	a.AzureConfig, err = common.ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

func (a *NitricAzurePulumiProvider) GetTags(stackID string, resourceName string, resourceType commonresources.ResourceType) map[string]string {
	resourceTags := tags.Tags(stackID, resourceName, resourceType)
	// Augment with common tags
	commonTags := a.CommonStackDetails.Tags

	// merge resourceTags and common tags
	combinedTags := make(map[string]string)
	for k, v := range commonTags {
		combinedTags[k] = v
	}
	for k, v := range resourceTags {
		combinedTags[k] = v
	}

	return combinedTags
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

func (a *NitricAzurePulumiProvider) RequiredProviders() map[string]interface{} {
	return map[string]interface{}{}
}

func (a *NitricAzurePulumiProvider) createDatabaseServer(ctx *pulumi.Context, tags map[string]string) error {
	var err error

	virtualNetworkName := "db-virtual-network"
	a.VirtualNetwork, err = network.NewVirtualNetwork(ctx, virtualNetworkName, &network.VirtualNetworkArgs{
		AddressSpace: &network.AddressSpaceArgs{
			AddressPrefixes: pulumi.StringArray{
				pulumi.String("10.0.0.0/16"),
			},
		},
		FlowTimeoutInMinutes: pulumi.Int(10),
		Location:             pulumi.String(a.Region),
		ResourceGroupName:    a.ResourceGroup.Name,
		VirtualNetworkName:   pulumi.String(virtualNetworkName),
	})
	if err != nil {
		return errors.WithMessage(err, "creating virtual network")
	}

	a.DatabaseSubnet, err = network.NewSubnet(ctx, "db-subnet", &network.SubnetArgs{
		AddressPrefix:      pulumi.String("10.0.0.0/18"),
		ResourceGroupName:  a.ResourceGroup.Name,
		SubnetName:         pulumi.String("db-subnet"),
		VirtualNetworkName: a.VirtualNetwork.Name,
		Delegations: network.DelegationArray{
			network.DelegationArgs{
				Name:        pulumi.String("db-delegation"),
				ServiceName: pulumi.String("Microsoft.DBforPostgreSQL/flexibleServers"),
			},
		},
	})
	if err != nil {
		return errors.WithMessage(err, "creating database subnet")
	}

	a.InfrastructureSubnet, err = network.NewSubnet(ctx, "infrastructure-subnet", &network.SubnetArgs{
		AddressPrefix:      pulumi.String("10.0.64.0/18"),
		ResourceGroupName:  a.ResourceGroup.Name,
		SubnetName:         pulumi.String("infrastructure-subnet"),
		VirtualNetworkName: a.VirtualNetwork.Name,
	}, pulumi.DependsOn([]pulumi.Resource{a.DatabaseSubnet}))
	if err != nil {
		return errors.WithMessage(err, "creating infrastructure subnet")
	}

	a.ContainerGroupSubnet, err = network.NewSubnet(ctx, "container-group-subnet", &network.SubnetArgs{
		AddressPrefix:      pulumi.String("10.0.192.0/18"),
		ResourceGroupName:  a.ResourceGroup.Name,
		SubnetName:         pulumi.String("container-group-subnet"),
		VirtualNetworkName: a.VirtualNetwork.Name,
		Delegations: network.DelegationArray{
			network.DelegationArgs{
				Name:        pulumi.String("container-instance-delegation"),
				ServiceName: pulumi.String("Microsoft.ContainerInstance/containerGroups"),
			},
		},
	}, pulumi.DependsOn([]pulumi.Resource{a.InfrastructureSubnet}))
	if err != nil {
		return errors.WithMessage(err, "creating container group subnet")
	}

	privateDns, err := network.NewPrivateZone(ctx, "db-private-dns", &network.PrivateZoneArgs{
		ResourceGroupName: a.ResourceGroup.Name,
		PrivateZoneName:   pulumi.String("db-private-dns.postgres.database.azure.com"),
		Location:          pulumi.String("global"),
	})
	if err != nil {
		return errors.WithMessage(err, "creating private dns zone")
	}

	vnetLink, err := network.NewVirtualNetworkLink(ctx, "db-private-dns-link", &network.VirtualNetworkLinkArgs{
		Location:            pulumi.String("global"),
		PrivateZoneName:     privateDns.Name,
		RegistrationEnabled: pulumi.Bool(false),
		ResourceGroupName:   a.ResourceGroup.Name,
		VirtualNetwork: &network.SubResourceArgs{
			Id: a.VirtualNetwork.ID(),
		},
		VirtualNetworkLinkName: pulumi.String("db-private-dns-link"),
	})
	if err != nil {
		return err
	}

	dbServerName := ResourceName(ctx, "", DatabaseServerRT)

	// generate a db random password
	a.DbMasterPassword, err = random.NewRandomPassword(ctx, "db-master-password", &random.RandomPasswordArgs{
		Length:  pulumi.Int(16),
		Special: pulumi.Bool(false),
	})
	if err != nil {
		return errors.WithMessage(err, "creating master password")
	}

	a.DatabaseServer, err = dbforpostgresql.NewServer(ctx, dbServerName, &dbforpostgresql.ServerArgs{
		ResourceGroupName:          a.ResourceGroup.Name,
		Location:                   a.ResourceGroup.Location,
		AdministratorLogin:         pulumi.String("nitric"),
		AdministratorLoginPassword: a.DbMasterPassword.Result,
		CreateMode:                 pulumi.String(dbforpostgresql.CreateModeDefault),
		AvailabilityZone:           pulumi.String("1"),
		Version:                    pulumi.String(dbforpostgresql.ServerVersion_14),
		Network: &dbforpostgresql.NetworkArgs{
			DelegatedSubnetResourceId:   a.DatabaseSubnet.ID(),
			PrivateDnsZoneArmResourceId: privateDns.ID(),
		},
		Sku: &dbforpostgresql.SkuArgs{
			Name: pulumi.String("Standard_B1ms"),
			Tier: pulumi.String(dbforpostgresql.SkuTierBurstable),
		},
		HighAvailability: &dbforpostgresql.HighAvailabilityArgs{
			Mode: pulumi.String(dbforpostgresql.HighAvailabilityModeDisabled),
		},
		Storage: &dbforpostgresql.StorageArgs{
			StorageSizeGB: pulumi.Int(32),
		},
		Tags: pulumi.ToStringMap(tags),
	}, pulumi.DependsOn([]pulumi.Resource{a.DatabaseSubnet, privateDns, vnetLink}))
	if err != nil {
		return err
	}

	return nil
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

	if a.AzureConfig.ResourceGroup != "" {
		rgId := fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", a.ClientConfig.SubscriptionId, a.AzureConfig.ResourceGroup)
		a.ResourceGroup, err = resources.GetResourceGroup(ctx, ResourceName(ctx, "", ResourceGroupRT), pulumi.ID(rgId), nil, pulumi.RetainOnDelete(true))
	} else {
		a.ResourceGroup, err = resources.NewResourceGroup(ctx, ResourceName(ctx, "", ResourceGroupRT), &resources.ResourceGroupArgs{
			Location: pulumi.String(a.Region),
			Tags:     pulumi.ToStringMap(a.GetTags(a.StackId, ctx.Stack(), commonresources.Stack)),
		})
	}

	if err != nil {
		return errors.WithMessage(err, "resource group create")
	}

	if hasResourceType(nitricResources, resourcespb.ResourceType_SqlDatabase) {
		logger.Info("Stack declares one or more databases, creating stack level PostgreSQL Database Server")
		err := a.createDatabaseServer(ctx, a.GetTags(a.StackId, ctx.Stack(), commonresources.Stack))
		if err != nil {
			return errors.WithMessage(err, "create azure sql flexible server")
		}
	}

	// Create a key vault if secrets are required.
	// Unlike AWS and GCP which have centralized secrets management, Azure allows for multiple key vaults.
	// This means we need to create a keyvault for each stack.
	if hasResourceType(nitricResources, resourcespb.ResourceType_Secret) {
		logger.Info("Stack declares one or more secrets, creating stack level Azure Key Vault")
		a.KeyVault, err = createKeyVault(ctx, a.ResourceGroup, a.ClientConfig.TenantId, a.GetTags(a.StackId, ctx.Stack(), commonresources.Stack))
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
		a.StorageAccount, err = createStorageAccount(ctx, a.ResourceGroup, a.GetTags(a.StackId, ctx.Stack(), commonresources.Stack))
		if err != nil {
			return errors.WithMessage(err, "storage account create")
		}
	}

	a.ContainerEnv, err = a.newContainerEnv(ctx, a.StackId, map[string]string{})
	if err != nil {
		return err
	}

	// Greedily create all the roles for consistency. Could be reduced to required roles only in future.
	a.Roles, err = a.CreateRoles(ctx, a.StackId, a.ClientConfig.SubscriptionId, a.ResourceGroup.Name)
	if err != nil {
		return err
	}

	return nil
}

func (a *NitricAzurePulumiProvider) Post(ctx *pulumi.Context) error {
	if len(a.WebsiteContainers) > 0 {
		err := a.deployCDN(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *NitricAzurePulumiProvider) Result(ctx *pulumi.Context) (pulumi.StringOutput, error) {
	outputs := []interface{}{}

	// Add Resource Group output as a link to the Azure Portal
	outputs = append(outputs, pulumi.Sprintf("Resource Group:\n──────────────\n%s\n", pulumi.Sprintf("https://portal.azure.com/#@%s/resource/subscriptions/%s/resourceGroups/%s/overview", a.ClientConfig.TenantId, a.ClientConfig.SubscriptionId, a.ResourceGroup.Name)))

	// Add APIs outputs
	if len(a.Apis) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.Apis {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", apiName, api.ApiManagementService.GatewayUrl))
		}
	}

	if a.Endpoint != nil {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("CDN:\n──────────────"))
		outputs = append(outputs, pulumi.Sprintf("https://%s", a.Endpoint.HostName))
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
		return pulumi.StringOutput{}, fmt.Errorf("failed to generate pulumi output")
	}

	return output, nil
}

func NewNitricAzurePulumiProvider() *NitricAzurePulumiProvider {
	principalsMap := map[resourcespb.ResourceType]map[string]*ServicePrincipal{}

	principalsMap[resourcespb.ResourceType_Service] = map[string]*ServicePrincipal{}

	return &NitricAzurePulumiProvider{
		Apis:                   make(map[string]ApiResources),
		HttpProxies:            make(map[string]ApiResources),
		Buckets:                make(map[string]*storage.BlobContainer),
		Queues:                 make(map[string]*storage.Queue),
		ContainerApps:          make(map[string]*ContainerApp),
		Topics:                 make(map[string]*eventgrid.Topic),
		SqlMigrations:          make(map[string]*containerinstance.ContainerGroup),
		Principals:             principalsMap,
		KeyValueStores:         make(map[string]*storage.Table),
		WebsiteStorageAccounts: make(map[string]*storage.StorageAccount),
		WebsiteContainers:      make(map[string]*storage.StorageAccountStaticWebsite),
	}
}
