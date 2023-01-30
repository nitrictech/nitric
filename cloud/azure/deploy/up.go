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
	"github.com/pulumi/pulumi-azure-native-sdk/eventgrid"
	"github.com/pulumi/pulumi-azure-native-sdk/keyvault"
	"github.com/pulumi/pulumi-azure-native-sdk/resources"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nitrictech/nitric/cloud/azure/deploy/api"
	"github.com/nitrictech/nitric/cloud/azure/deploy/exec"
	"github.com/nitrictech/nitric/cloud/azure/deploy/storage"
	"github.com/nitrictech/nitric/cloud/azure/deploy/subscription"
	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
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

		envMap := map[string]string{}
		hasExecUnit := false
		bucketNames := []string{}
		queueNames := []string{}
		execUnits := map[string]*deploy.ExecutionUnit{}

		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_ExecutionUnit:
				execUnits[res.Name] = t.ExecutionUnit
			case *deploy.Resource_Queue:
				queueNames = append(queueNames, res.Name)
			case *deploy.Resource_Bucket:
				bucketNames = append(bucketNames, res.Name)
			}
		}

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

		if len(bucketNames) > 0 || len(queueNames) > 0 {
			sr, err := storage.NewStorageResources(ctx, "storage", &storage.StorageArgs{
				ResourceGroupName: rg.Name,
				StackID:           stackID,
				BucketNames:       bucketNames,
				QueueNames:        queueNames,
			})
			if err != nil {
				return errors.WithMessage(err, "storage create")
			}

			contEnvArgs.StorageAccountBlobEndpoint = sr.Account.PrimaryEndpoints.Blob()
			contEnvArgs.StorageAccountQueueEndpoint = sr.Account.PrimaryEndpoints.Queue()
		}

		topics := map[string]*eventgrid.Topic{}
		execSubscriptions := map[string]map[string]*eventgrid.Topic{}

		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Schedule:
				// TODO: Add schedule support
				// NOTE: Currently CRONTAB support is required, we either need to revisit the design of
				// our scheduled expressions or implement a workaround or request a feature.
				_ = ctx.Log.Warn("Schedules are not currently supported for Azure deployments", &pulumi.LogArgs{})

			case *deploy.Resource_Topic:
				topics[res.Name], err = eventgrid.NewTopic(ctx, utils.ResourceName(ctx, res.Name, utils.EventGridRT), &eventgrid.TopicArgs{
					ResourceGroupName: rg.Name,
					Location:          rg.Location,
					Tags:              common.Tags(ctx, stackID, res.Name),
				})
				if err != nil {
					return errors.WithMessage(err, "eventgrid topic "+res.Name)
				}

				for _, s := range t.Topic.Subscriptions {
					targ := s.Target.(*deploy.SubscriptionTarget_ExecutionUnit)
					subs, ok := execSubscriptions[targ.ExecutionUnit]
					if ok {
						subs[res.Name] = topics[res.Name]
						execSubscriptions[targ.ExecutionUnit] = subs
					} else {
						execSubscriptions[targ.ExecutionUnit] = map[string]*eventgrid.Topic{res.Name: topics[res.Name]}
					}
				}
			}
		}

		var contEnv *exec.ContainerEnv

		apps := map[string]*exec.ContainerApp{}

		if hasExecUnit {
			contEnv, err = exec.NewContainerEnv(ctx, "containerEnv", contEnvArgs)
			if err != nil {
				return errors.WithMessage(err, "containerApps")
			}

			for name, eu := range execUnits {
				repositoryUrl := pulumi.Sprintf("%s/%s-%s-%s", contEnv.Registry.LoginServer, details.Project, name, "azure")

				image, err := image.NewImage(ctx, name, &image.ImageArgs{
					SourceImage:   eu.GetImage().GetUri(),
					RepositoryUrl: repositoryUrl,
					Username:      contEnv.RegistryUser.Elem(),
					Password:      contEnv.RegistryUser.Elem(),
					Server:        contEnv.Registry.LoginServer,
					Runtime:       runtime,
				}, pulumi.Parent(contEnv))
				if err != nil {
					return errors.WithMessage(err, "function image tag "+name)
				}

				apps[name], err = exec.NewContainerApp(ctx, name, &exec.ContainerAppArgs{
					ResourceGroupName: rg.Name,
					Location:          pulumi.String(details.Region),
					SubscriptionID:    pulumi.String(clientConfig.SubscriptionId),
					Registry:          contEnv.Registry,
					RegistryUser:      contEnv.RegistryUser,
					RegistryPass:      contEnv.RegistryPass,
					ManagedEnv:        contEnv.ManagedEnv,
					ImageUri:          image.URI(),
					Env:               contEnv.Env,
					ExecutionUnit:     eu,
					ManagedIdentityID: contEnv.ManagedUser.ClientId,
				}, pulumi.Parent(contEnv))
				if err != nil {
					return err
				}

				_, err = subscription.NewSubscription(ctx, name+"-subscription", &subscription.SubscriptionArgs{
					ResourceGroupName:  rg.Name,
					Sp:                 apps[name].Sp,
					AppName:            name,
					LatestRevisionFqdn: apps[name].App.LatestRevisionFqdn,
					Subscriptions:      execSubscriptions[name],
				}, pulumi.Parent(apps[name]))
				if err != nil {
					return errors.WithMessage(err, "subscriptions")
				}
			}
		}

		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Api:
				if t.Api.GetOpenapi() == "" {
					return fmt.Errorf("azure provider can only deploy OpenAPI specs")
				}

				doc := &openapi3.T{}
				err := doc.UnmarshalJSON([]byte(t.Api.GetOpenapi()))
				if err != nil {
					return fmt.Errorf("invalid document suppled for api: %s", res.Name)
				}

				_, err = api.NewAzureApiManagement(ctx, res.Name, &api.AzureApiManagementArgs{
					ResourceGroupName: rg.Name,
					OrgName:           pulumi.String(details.Org),
					AdminEmail:        pulumi.String(details.AdminEmail),
					OpenAPISpec:       doc,
					Apps:              apps,
					ManagedIdentity:   contEnv.ManagedUser,
					StackID:           stackID,
				})
				if err != nil {
					return errors.WithMessage(err, "gateway "+res.Name)
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	_ = pulumiStack.SetConfig(context.TODO(), "azure:region", auto.ConfigValue{Value: details.Region})

	messageWriter := &UpStreamMessageWriter{
		stream: stream,
	}

	// Run the program
	_, err = pulumiStack.Up(context.TODO(), optup.ProgressStreams(messageWriter))
	if err != nil {
		return err
	}

	return nil
}
