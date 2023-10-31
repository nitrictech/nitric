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

package exec

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"

	app "github.com/pulumi/pulumi-azure-native-sdk/app"
	"github.com/pulumi/pulumi-azure-native-sdk/containerregistry"
	"github.com/pulumi/pulumi-azure-native-sdk/managedidentity"
	"github.com/pulumi/pulumi-azure-native-sdk/operationalinsights"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/azure/deploy/utils"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
)

type ContainerEnvArgs struct {
	ResourceGroupName pulumi.StringInput
	Location          pulumi.StringInput
	EnvMap            map[string]string
	StackID           string

	KVaultName                  pulumi.StringInput
	StorageAccountBlobEndpoint  pulumi.StringInput
	StorageAccountQueueEndpoint pulumi.StringInput
}

type ContainerEnv struct {
	pulumi.ResourceState

	Name         string
	Registry     *containerregistry.Registry
	RegistryUser pulumi.StringPtrOutput
	RegistryPass pulumi.StringPtrOutput
	ManagedEnv   *app.ManagedEnvironment
	Env          app.EnvironmentVarArray
	ManagedUser  *managedidentity.UserAssignedIdentity
}

func NewContainerEnv(ctx *pulumi.Context, name string, args *ContainerEnvArgs, opts ...pulumi.ResourceOption) (*ContainerEnv, error) {
	res := &ContainerEnv{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:func:ContainerEnv", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.ManagedUser, err = managedidentity.NewUserAssignedIdentity(ctx, "managed-identity", &managedidentity.UserAssignedIdentityArgs{
		Location:          args.Location,
		ResourceGroupName: args.ResourceGroupName,
		ResourceName:      pulumi.String("managed-identity"),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	env := app.EnvironmentVarArray{}

	if args.StorageAccountBlobEndpoint != nil {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_STORAGE_ACCOUNT_BLOB_ENDPOINT"),
			Value: args.StorageAccountBlobEndpoint,
		})
	}

	if args.StorageAccountQueueEndpoint != nil {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String("AZURE_STORAGE_ACCOUNT_QUEUE_ENDPOINT"),
			Value: args.StorageAccountQueueEndpoint,
		})
	}

	if args.KVaultName != nil {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String("KVAULT_NAME"),
			Value: args.KVaultName,
		})
	}

	env = append(env, app.EnvironmentVarArgs{
		Name:  pulumi.String("NITRIC_HTTP_PROXY_PORT"),
		Value: pulumi.String(fmt.Sprint(3000)),
	})

	for k, v := range args.EnvMap {
		env = append(env, app.EnvironmentVarArgs{
			Name:  pulumi.String(k),
			Value: pulumi.String(v),
		})
	}

	res.Env = env

	res.Registry, err = containerregistry.NewRegistry(ctx, utils.ResourceName(ctx, name, utils.RegistryRT), &containerregistry.RegistryArgs{
		ResourceGroupName: args.ResourceGroupName,
		Location:          args.Location,
		AdminUserEnabled:  pulumi.BoolPtr(true),
		Sku: containerregistry.SkuArgs{
			Name: pulumi.String("Basic"),
		},
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	aw, err := operationalinsights.NewWorkspace(ctx, utils.ResourceName(ctx, name, utils.AnalyticsWorkspaceRT), &operationalinsights.WorkspaceArgs{
		Location:          args.Location,
		ResourceGroupName: args.ResourceGroupName,
		Sku: &operationalinsights.WorkspaceSkuArgs{
			Name: pulumi.String("PerGB2018"),
		},
		RetentionInDays: pulumi.Int(30),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	sharedKeys := operationalinsights.GetSharedKeysOutput(ctx, operationalinsights.GetSharedKeysOutputArgs{
		ResourceGroupName: args.ResourceGroupName,
		WorkspaceName:     aw.Name,
	})

	res.ManagedEnv, err = app.NewManagedEnvironment(ctx, utils.ResourceName(ctx, name, utils.KubeRT), &app.ManagedEnvironmentArgs{
		Location:          args.Location,
		ResourceGroupName: args.ResourceGroupName,
		AppLogsConfiguration: app.AppLogsConfigurationArgs{
			Destination: pulumi.String("log-analytics"),
			LogAnalyticsConfiguration: app.LogAnalyticsConfigurationArgs{
				SharedKey:  sharedKeys.PrimarySharedKey(),
				CustomerId: aw.CustomerId,
			},
		},
		Tags: pulumi.ToStringMap(common.Tags(args.StackID, ctx.Stack()+"Kube", resources.ExecutionUnit)),
	}, pulumi.Parent(res))
	if err != nil {
		return nil, err
	}

	creds := pulumi.All(args.ResourceGroupName, res.Registry.Name).ApplyT(func(args []interface{}) (*containerregistry.ListRegistryCredentialsResult, error) {
		rgName := args[0].(string)
		regName := args[1].(string)

		return containerregistry.ListRegistryCredentials(ctx, &containerregistry.ListRegistryCredentialsArgs{
			ResourceGroupName: rgName,
			RegistryName:      regName,
		})
	})

	res.RegistryUser = creds.ApplyT(func(arg interface{}) *string {
		cred := arg.(*containerregistry.ListRegistryCredentialsResult)
		return cred.Username
	}).(pulumi.StringPtrOutput)

	res.RegistryPass = creds.ApplyT(func(arg interface{}) (*string, error) {
		cred := arg.(*containerregistry.ListRegistryCredentialsResult)

		if len(cred.Passwords) == 0 || cred.Passwords[0].Value == nil {
			return nil, fmt.Errorf("cannot retrieve container registry credentials")
		}

		return cred.Passwords[0].Value, nil
	}).(pulumi.StringPtrOutput)

	return res, nil
}
