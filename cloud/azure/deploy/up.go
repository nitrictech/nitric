// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploy

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/authorization"
	"github.com/pulumi/pulumi-azure-native-sdk/keyvault"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-azure-native-sdk/storage"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nitrictech/nitric/cloud/azure/deploy/api"
	"github.com/nitrictech/nitric/cloud/azure/deploy/bucket"
	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/queue"
	"github.com/nitrictech/nitric/cloud/azure/deploy/topic"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	azureStorage "github.com/pulumi/pulumi-azure-native-sdk/storage"
)

type UpStreamMessageWriter struct {
	stream deploy.DeployService_UpServer
}

func (s *UpStreamMessageWriter) Write(bytes []byte) (int, error) {
	err := s.stream.Send(&deploy.DeployUpEvent{
		Content: &deploy.DeployUpEvent_Message{
			Message: &deploy.DeployEventMessage{
				Message: string(bytes),
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}

// Up - Deploy requested infrastructure for a stack
func (d *DeployServer) Up(request *deploy.DeployUpRequest, stream deploy.DeployService_UpServer) error {
	details, err := getStackDetailsFromAttributes(request.Attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	pulumiStack, err := auto.UpsertStackInlineSource(context.TODO(), details.Stack, details.Project, func(ctx *pulumi.Context) error {
		// Calculate unique stackID
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
		stackID := pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.ID())

		clientConfig, err := authorization.GetClientConfig(ctx)
		if err != nil {
			return err
		}

		rg, err := resources.NewResourceGroup(ctx, utils.ResourceName(ctx, "", utils.ResourceGroupRT), &resources.ResourceGroupArgs{
			Location: pulumi.String(details.Region),
			Tags:     common.Tags(ctx, stackID, ctx.Stack()),
		})
		if err != nil {
			return errors.WithMessage(err, "resource group create")
		}

		// Get Execution units
		executionUnits := lo.Filter[*deploy.Resource](request.Spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetExecutionUnit() != nil
		})

		// Get Queues
		queues := lo.Filter[*deploy.Resource](request.Spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetQueue() != nil
		})

		// Get Buckets
		buckets := lo.Filter[*deploy.Resource](request.Spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetBucket() != nil
		})

		// Get Topics
		topics := lo.Filter[*deploy.Resource](request.Spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetTopic() != nil
		})

		// Get Topics
		schedules := lo.Filter[*deploy.Resource](request.Spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetSchedule() != nil
		})

		apis := lo.Filter[*deploy.Resource](request.Spec.Resources, func(item *deploy.Resource, index int) bool {
			return item.GetApi() != nil
		})

		envMap := map[string]string{}
		contEnvArgs := &exec.ContainerEnvArgs{
			ResourceGroupName: rg.Name,
			Location:          rg.Location,
			EnvMap:            envMap,
			StackID:           stackID,
		}

		// Create a stack level keyvault if secrets are enabled
		// At the moment secrets have no config level setting
		kvName := utils.ResourceName(ctx, "", utils.KeyVaultRT)

		kv, err := keyvault.NewVault(ctx, kvName, &keyvault.VaultArgs{
			Location:          rg.Location,
			ResourceGroupName: rg.Name,
			Properties: &keyvault.VaultPropertiesArgs{
				EnableSoftDelete:        pulumi.Bool(false),
				EnableRbacAuthorization: pulumi.Bool(true),
				Sku: &keyvault.SkuArgs{
					Family: pulumi.String("A"),
					Name:   keyvault.SkuNameStandard,
				},
				TenantId: pulumi.String(clientConfig.TenantId),
			},
			Tags: common.Tags(ctx, stackID, kvName),
		})
		if err != nil {
			return err
		}

		contEnvArgs.KVaultName = kv.Name

		// Create a storage account if buckets or queues are required
		var storageAccount *azureStorage.StorageAccount
		if len(buckets) > 0 || len(queues) > 0 {
			accName := utils.ResourceName(ctx, details.Stack, utils.StorageAccountRT)
			storageAccount, err = azureStorage.NewStorageAccount(ctx, accName, &storage.StorageAccountArgs{
				AccessTier:        azureStorage.AccessTierHot,
				ResourceGroupName: rg.Name,
				Kind:              pulumi.String("StorageV2"),
				Sku: azureStorage.SkuArgs{
					Name: pulumi.String(storage.SkuName_Standard_LRS),
				},
				Tags: common.Tags(ctx, stackID, accName),
			})
			if err != nil {
				return err
			}

			contEnvArgs.StorageAccountBlobEndpoint = storageAccount.PrimaryEndpoints.Blob()
			contEnvArgs.StorageAccountQueueEndpoint = storageAccount.PrimaryEndpoints.Queue()
		}

		// For each bucket create a new bucket
		for _, b := range buckets {
			_, err := bucket.NewAzureStorageBucket(ctx, b.Name, &bucket.AzureStorageBucketArgs{
				StackID:       stackID,
				Account:       storageAccount,
				ResourceGroup: rg,
			})
			if err != nil {
				return err
			}
		}

		// For each queue create a new queue
		for _, q := range queues {
			_, err := queue.NewAzureStorageQueue(ctx, q.Name, &queue.AzureStorageQueueArgs{
				StackID:       stackID,
				Account:       storageAccount,
				ResourceGroup: rg,
			})
			if err != nil {
				return err
			}
		}

		deployedTopics := map[string]*topic.AzureEventGridTopic{}

		var contEnv *exec.ContainerEnv

		apps := map[string]*exec.ContainerApp{}

		if len(executionUnits) > 0 {
			contEnv, err = exec.NewContainerEnv(ctx, "containerEnv", contEnvArgs)
			if err != nil {
				return errors.WithMessage(err, "containerApps")
			}

			for _, eu := range executionUnits {
				repositoryUrl := pulumi.Sprintf("%s/%s-%s-%s", contEnv.Registry.LoginServer, details.Project, eu.Name, "azure")

				image, err := image.NewImage(ctx, eu.Name, &image.ImageArgs{
					SourceImage:   eu.GetExecutionUnit().GetImage().GetUri(),
					RepositoryUrl: repositoryUrl,
					Username:      contEnv.RegistryUser.Elem(),
					Password:      contEnv.RegistryPass.Elem(),
					Server:        contEnv.Registry.LoginServer,
					Runtime:       runtime,
				}, pulumi.Parent(contEnv))
				if err != nil {
					return err
				}

				apps[eu.Name], err = exec.NewContainerApp(ctx, eu.Name, &exec.ContainerAppArgs{
					ResourceGroupName: rg.Name,
					Location:          pulumi.String(details.Region),
					SubscriptionID:    pulumi.String(clientConfig.SubscriptionId),
					Registry:          contEnv.Registry,
					RegistryUser:      contEnv.RegistryUser,
					RegistryPass:      contEnv.RegistryPass,
					ManagedEnv:        contEnv.ManagedEnv,
					ImageUri:          image.URI(),
					Env:               contEnv.Env,
					ExecutionUnit:     eu.GetExecutionUnit(),
					ManagedIdentityID: contEnv.ManagedUser.ClientId,
				}, pulumi.Parent(contEnv))
				if err != nil {
					return err
				}
			}
		}

		for _, t := range topics {
			deployedTopics[t.Name], err = topic.NewAzureEventGridTopic(ctx, utils.ResourceName(ctx, t.Name, utils.EventGridRT), &topic.AzureEventGridTopicArgs{
				ResourceGroup: rg,
				StackID:       stackID,
			})
			if err != nil {
				return err
			}

			for _, s := range t.GetTopic().Subscriptions {
				_, err = topic.NewAzureEventGridTopicSubscription(ctx, utils.ResourceName(ctx, fmt.Sprintf("%s-%s", t.Name, s.GetExecutionUnit()), utils.EventSubscriptionRT), &topic.AzureEventGridTopicSubscriptionArgs{
					Topic:  deployedTopics[t.Name],
					Target: apps[s.GetExecutionUnit()],
				})
				if err != nil {
					return err
				}
			}
		}

		if len(schedules) > 0 {
			// TODO: Add schedule support
			// NOTE: Currently CRONTAB support is required, we either need to revisit the design of
			// our scheduled expressions or implement a workaround or request a feature.
			_ = ctx.Log.Warn("Schedules are not currently supported for Azure deployments", &pulumi.LogArgs{})
		}

		for _, a := range apis {
			if a.GetApi().GetOpenapi() == "" {
				return fmt.Errorf("azure provider can only deploy OpenAPI specs")
			}

			doc := &openapi3.T{}
			err := doc.UnmarshalJSON([]byte(a.GetApi().GetOpenapi()))
			if err != nil {
				return fmt.Errorf("invalid document suppled for api: %s", a.Name)
			}

			_, err = api.NewAzureApiManagement(ctx, a.Name, &api.AzureApiManagementArgs{
				ResourceGroupName: rg.Name,
				OrgName:           pulumi.String(details.Org),
				AdminEmail:        pulumi.String(details.AdminEmail),
				OpenAPISpec:       doc,
				Apps:              apps,
				ManagedIdentity:   contEnv.ManagedUser,
				StackID:           stackID,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	_ = pulumiStack.SetConfig(context.TODO(), "azure-native:location", auto.ConfigValue{Value: details.Region})
	_ = pulumiStack.SetConfig(context.TODO(), "azure:location", auto.ConfigValue{Value: details.Region})

	messageWriter := &UpStreamMessageWriter{
		stream: stream,
	}

	_, err = pulumiStack.Up(context.TODO(), optup.ProgressStreams(messageWriter))
	// Run the program
	// _, err = pulumiStack.Up(context.TODO(), optup.ProgressStreams(messageWriter))
	if err != nil {
		return err
	}

	return nil
}
