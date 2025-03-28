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

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"

	app "github.com/pulumi/pulumi-azure-native-sdk/app"
	"github.com/pulumi/pulumi-azure-native-sdk/containerregistry"
	"github.com/pulumi/pulumi-azure-native-sdk/managedidentity"
	"github.com/pulumi/pulumi-azure-native-sdk/operationalinsights"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type ContainerEnvArgs struct {
	// ResourceGroupName pulumi.StringInput
	// Location          pulumi.StringInput
	EnvMap map[string]string
	// StackID           string

	// KVaultName                  pulumi.StringInput
	// StorageAccountBlobEndpoint  pulumi.StringInput
	// StorageAccountQueueEndpoint pulumi.StringInput
}

type ContainerEnv struct {
	pulumi.ResourceState

	Name           string
	DockerProvider *docker.Provider
	Registry       *containerregistry.Registry
	RegistryArgs   *docker.RegistryArgs
	ManagedEnv     *app.ManagedEnvironment
	Env            app.EnvironmentVarArray
	ManagedUser    *managedidentity.UserAssignedIdentity
}

func (p *NitricAzurePulumiProvider) newContainerEnv(ctx *pulumi.Context, name string, envMap map[string]string, opts ...pulumi.ResourceOption) (*ContainerEnv, error) {
	res := &ContainerEnv{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitricazure:ContainerEnv", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.ManagedUser, err = managedidentity.NewUserAssignedIdentity(ctx, "managed-identity", &managedidentity.UserAssignedIdentityArgs{
		Location:          p.ResourceGroup.Location,
		ResourceGroupName: p.ResourceGroup.Name,
		ResourceName:      pulumi.Sprintf("managed-identity-%s", p.StackId),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	env := app.EnvironmentVarArray{}

	if p.StorageAccount != nil {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_STORAGE_ACCOUNT_NAME"),
			Value: p.StorageAccount.Name,
		})
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_STORAGE_ACCOUNT_BLOB_ENDPOINT"),
			Value: p.StorageAccount.PrimaryEndpoints.Blob(),
		})
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_STORAGE_ACCOUNT_QUEUE_ENDPOINT"),
			Value: p.StorageAccount.PrimaryEndpoints.Queue(),
		})
	}

	if p.KeyVault != nil {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String("KVAULT_NAME"),
			Value: p.KeyVault.Name,
		})
	}

	env = append(env, app.EnvironmentVarArgs{
		Name:  pulumi.String("NITRIC_HTTP_PROXY_PORT"),
		Value: pulumi.String(fmt.Sprint(3000)),
	})

	for k, v := range envMap {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String(k),
			Value: pulumi.String(v),
		})
	}

	res.Env = env

	res.Registry, err = containerregistry.NewRegistry(ctx, ResourceName(ctx, name, RegistryRT), &containerregistry.RegistryArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		Location:          p.ResourceGroup.Location,
		AdminUserEnabled:  pulumi.BoolPtr(true),
		Sku: containerregistry.SkuArgs{
			Name: pulumi.String("Basic"),
		},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	aw, err := operationalinsights.NewWorkspace(ctx, ResourceName(ctx, name, AnalyticsWorkspaceRT), &operationalinsights.WorkspaceArgs{
		Location:          p.ResourceGroup.Location,
		ResourceGroupName: p.ResourceGroup.Name,
		Sku: &operationalinsights.WorkspaceSkuArgs{
			Name: pulumi.String("PerGB2018"),
		},
		RetentionInDays: pulumi.Int(30),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	sharedKeys := operationalinsights.GetSharedKeysOutput(ctx, operationalinsights.GetSharedKeysOutputArgs{
		ResourceGroupName: p.ResourceGroup.Name,
		WorkspaceName:     aw.Name,
	})

	managementArgs := &app.ManagedEnvironmentArgs{
		Location:          p.ResourceGroup.Location,
		ResourceGroupName: p.ResourceGroup.Name,
		AppLogsConfiguration: app.AppLogsConfigurationArgs{
			Destination: pulumi.String("log-analytics"),
			LogAnalyticsConfiguration: app.LogAnalyticsConfigurationArgs{
				SharedKey:  sharedKeys.PrimarySharedKey(),
				CustomerId: aw.CustomerId,
			},
		},
		Tags: pulumi.ToStringMap(p.GetTags(p.StackId, ctx.Stack()+"Kube", resources.Service)),
	}

	if p.InfrastructureSubnet != nil {
		managementArgs.VnetConfiguration = &app.VnetConfigurationArgs{
			InfrastructureSubnetId: p.InfrastructureSubnet.ID(),
		}
	}

	res.ManagedEnv, err = app.NewManagedEnvironment(ctx, ResourceName(ctx, name, KubeRT), managementArgs, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	creds := pulumi.All(p.ResourceGroup.Name, res.Registry.Name).ApplyT(func(args []interface{}) (*containerregistry.ListRegistryCredentialsResult, error) {
		rgName := args[0].(string)
		regName := args[1].(string)

		return containerregistry.ListRegistryCredentials(ctx, &containerregistry.ListRegistryCredentialsArgs{
			ResourceGroupName: rgName,
			RegistryName:      regName,
		})
	})

	registryUser := creds.ApplyT(func(arg interface{}) *string {
		cred := arg.(*containerregistry.ListRegistryCredentialsResult)
		return cred.Username
	}).(pulumi.StringPtrOutput)

	registryPass := creds.ApplyT(func(arg interface{}) (*string, error) {
		cred := arg.(*containerregistry.ListRegistryCredentialsResult)

		if len(cred.Passwords) == 0 || cred.Passwords[0].Value == nil {
			return nil, fmt.Errorf("cannot retrieve container registry credentials")
		}

		return cred.Passwords[0].Value, nil
	}).(pulumi.StringPtrOutput)

	res.RegistryArgs = &docker.RegistryArgs{
		Server:   res.Registry.LoginServer,
		Username: registryUser,
		Password: registryPass,
	}

	res.DockerProvider, err = docker.NewProvider(ctx, "docker-auth-provider", &docker.ProviderArgs{
		RegistryAuth: &docker.ProviderRegistryAuthArray{
			docker.ProviderRegistryAuthArgs{
				Address:  res.Registry.LoginServer,
				Username: registryUser,
				Password: registryPass,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	err = ctx.RegisterResourceOutputs(res, pulumi.Map{})

	return res, err
}
