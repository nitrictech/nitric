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

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/containerinstance/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/dbforpostgresql/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAzurePulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{a.DatabaseServer})}

	_, err := dbforpostgresql.NewDatabase(ctx, name, &dbforpostgresql.DatabaseArgs{
		DatabaseName:      pulumi.String(name),
		ResourceGroupName: a.ResourceGroup.Name,
		ServerName:        a.DatabaseServer.Name,
		Charset:           pulumi.String("utf8"),
		Collation:         pulumi.String("en_US.utf8"),
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("unable to create nitric database %s: failed to create Azure SQL Database", name))
	}

	if config.GetImageUri() != "" {
		repositoryUrl := pulumi.Sprintf("%s/%s", a.ContainerEnv.Registry.LoginServer, config.GetImageUri())

		inspect, err := image.CommandFromImageInspect(config.GetImageUri(), " ")
		if err != nil {
			return err
		}

		image, err := image.NewLocalImage(ctx, name, &image.LocalImageArgs{
			SourceImage:   config.GetImageUri(),
			SourceImageID: inspect.ID,
			RepositoryUrl: repositoryUrl,
		}, pulumi.Provider(a.ContainerEnv.DockerProvider))
		if err != nil {
			return err
		}

		containerGroupName := fmt.Sprintf("%s-migration-group", name)

		databaseUrl := pulumi.Sprintf("postgres://%s:%s@%s:%s/%s", "nitric", a.DbMasterPassword.Result, a.DatabaseServer.FullyQualifiedDomainName, "5432", name)

		a.SqlMigrations[name], err = containerinstance.NewContainerGroup(ctx, containerGroupName, &containerinstance.ContainerGroupArgs{
			ContainerGroupName: pulumi.String(containerGroupName),
			Containers: containerinstance.ContainerArray{
				&containerinstance.ContainerArgs{
					Image: pulumi.Sprintf("%s/%s-migrations:latest", a.ContainerEnv.Registry.LoginServer, image.Name),
					Name:  pulumi.Sprintf("%s-migration", name),
					Resources: &containerinstance.ResourceRequirementsArgs{
						Requests: &containerinstance.ResourceRequestsArgs{
							Cpu:        pulumi.Float64(1),
							MemoryInGB: pulumi.Float64(1),
						},
					},
					EnvironmentVariables: containerinstance.EnvironmentVariableArray{
						containerinstance.EnvironmentVariableArgs{
							Name:        pulumi.String("DB_URL"),
							SecureValue: databaseUrl,
						},
						containerinstance.EnvironmentVariableArgs{
							Name:  pulumi.String("NITRIC_DB_NAME"),
							Value: pulumi.String(name),
						},
					},
				},
			},
			Location:          pulumi.String(a.Region),
			OsType:            pulumi.String(containerinstance.OperatingSystemTypesLinux),
			ResourceGroupName: a.ResourceGroup.Name,
			RestartPolicy:     pulumi.String(containerinstance.ContainerGroupRestartPolicyNever),
			Sku:               pulumi.String(containerinstance.ContainerGroupSkuStandard),
			SubnetIds: &containerinstance.ContainerGroupSubnetIdArray{
				containerinstance.ContainerGroupSubnetIdArgs{
					Id:   a.ContainerGroupSubnet.ID(),
					Name: a.ContainerGroupSubnet.Name,
				},
			},
			ImageRegistryCredentials: &containerinstance.ImageRegistryCredentialArray{
				&containerinstance.ImageRegistryCredentialArgs{
					Username: a.ContainerEnv.RegistryUser.Elem(),
					Password: a.ContainerEnv.RegistryPass.Elem(),
					Server:   a.ContainerEnv.Registry.LoginServer,
				},
			},
		}, opts...)
		if err != nil {
			return err
		}
	}

	return nil
}
