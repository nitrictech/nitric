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
	commonresources "github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/core/pkg/logger"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
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

//go:embed runtime-azure
var runtime []byte

type ApiResources struct {
	ApiManagementService *apimanagement.ApiManagementService
	Api                  *apimanagement.Api
}

type NitricAzurePulumiProvider struct {
	stackId       string
	projectName   string
	stackName     string
	fullStackName string
	resources     []*deploymentspb.Resource

	config *AzureConfig
	region string

	clientConfig *authorization.GetClientConfigResult

	resourceGroup  *resources.ResourceGroup
	keyVault       *keyvault.Vault
	storageAccount *storage.StorageAccount

	containerEnv *ContainerEnv

	apis        map[string]ApiResources
	httpProxies map[string]ApiResources
	buckets     map[string]*storage.BlobContainer

	queues map[string]*storage.Queue

	principals map[resourcespb.ResourceType]map[string]*ServicePrincipal

	containerApps map[string]*ContainerApp
	topics        map[string]*eventgrid.Topic

	keyValueStores map[string]*storage.Table

	roles *Roles
	provider.NitricDefaultOrder
}

var _ provider.NitricPulumiProvider = (*NitricAzurePulumiProvider)(nil)

const (
	pulumiAzureNativeVersion = "1.95.0"
	pulumiAzureVersion       = "5.52.0"
)

func (a *NitricAzurePulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"azure-native:location": auto.ConfigValue{Value: a.region},
		"azure:location":        auto.ConfigValue{Value: a.region},
		"azure-native:version":  auto.ConfigValue{Value: pulumiAzureNativeVersion},
		"azure:version":         auto.ConfigValue{Value: pulumiAzureVersion},
		"docker:version":        auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricAzurePulumiProvider) Init(attributes map[string]interface{}) error {
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

func hasResourceType(resources []*deploymentspb.Resource, resourceType resourcespb.ResourceType) bool {
	for _, r := range resources {
		if r.GetId().GetType() == resourceType {
			return true
		}
	}

	return false
}

func (a *NitricAzurePulumiProvider) Pre(ctx *pulumi.Context, nitricResources []*deploymentspb.Resource) error {
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

	a.stackId = <-stackIdChan

	a.clientConfig, err = authorization.GetClientConfig(ctx)
	if err != nil {
		return err
	}

	a.resourceGroup, err = resources.NewResourceGroup(ctx, ResourceName(ctx, "", ResourceGroupRT), &resources.ResourceGroupArgs{
		Location: pulumi.String(a.region),
		Tags:     pulumi.ToStringMap(tags.Tags(a.stackId, ctx.Stack(), commonresources.Stack)),
	})
	if err != nil {
		return errors.WithMessage(err, "resource group create")
	}

	// Create a key vault if secrets are required.
	// Unlike AWS and GCP which have centralized secrets management, Azure allows for multiple key vaults.
	// This means we need to create a keyvault for each stack.
	if hasResourceType(nitricResources, resourcespb.ResourceType_Secret) {
		logger.Info("Stack declares one or more secrets, creating stack level Azure Key Vault")
		a.keyVault, err = createKeyVault(ctx, a.resourceGroup, a.clientConfig.TenantId, tags.Tags(a.stackId, ctx.Stack(), commonresources.Stack))
		if err != nil {
			return errors.WithMessage(err, "keyvault create")
		}
	}

	// Create a storage account if buckets
	// Unlike AWS and GCP which have centralized storage management, Azure allows for multiple storage accounts.
	// This means we need to create a storage account for each stack, before buckets can be created.
	if hasResourceType(nitricResources, resourcespb.ResourceType_Bucket) {
		logger.Info("Stack declares one or more buckets, creating stack level Azure Storage Account")
		a.storageAccount, err = createStorageAccount(ctx, a.resourceGroup, tags.Tags(a.stackId, ctx.Stack(), commonresources.Stack))
		if err != nil {
			return errors.WithMessage(err, "storage account create")
		}
	}

	a.containerEnv, err = a.newContainerEnv(ctx, a.stackId, map[string]string{})
	if err != nil {
		return err
	}

	// Greedily create all the roles for consistency. Could be reduced to required roles only in future.
	a.roles, err = CreateRoles(ctx, a.stackId, a.clientConfig.SubscriptionId, a.resourceGroup.Name)
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
	if len(a.apis) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.apis {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", apiName, api.ApiManagementService.GatewayUrl))
		}
	}

	// Add HTTP Proxy outputs
	if len(a.httpProxies) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("HTTP Proxies:\n──────────────"))
		for proxyName, proxy := range a.httpProxies {
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
		apis:           map[string]ApiResources{},
		httpProxies:    map[string]ApiResources{},
		buckets:        make(map[string]*storage.BlobContainer),
		queues:         make(map[string]*storage.Queue),
		containerApps:  map[string]*ContainerApp{},
		topics:         map[string]*eventgrid.Topic{},
		principals:     principalsMap,
		keyValueStores: map[string]*storage.Table{},
	}
}
