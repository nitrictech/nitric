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

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/containerinstance/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/dbforpostgresql/v2"
	"github.com/pulumi/pulumi-azure-native-sdk/network/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAzurePulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {
	var err error
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	db, err := dbforpostgresql.NewDatabase(ctx, name, &dbforpostgresql.DatabaseArgs{
		DatabaseName:      pulumi.String(name),
		ResourceGroupName: a.ResourceGroup.Name,
		ServerName:        a.DatabaseServer.Name,
		Charset:           pulumi.String("utf8"),
		Collation:         pulumi.String("en_US.utf8"),
	}, opts...)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("unable to create nitric database %s: failed to create Azure SQL Database", name))
	}

	a.SqlDatabases[name] = db

	if config.GetImageUri() != "" {
		imageUriSplit := strings.Split(config.GetImageUri(), "/")
		imageName := imageUriSplit[len(imageUriSplit)-1]

		repositoryUrl := pulumi.Sprintf("%s/%s", a.ContainerEnv.Registry.LoginServer, imageName)

		image, err := image.NewLocalImage(ctx, name, &image.LocalImageArgs{
			SourceImage:   config.GetImageUri(),
			RepositoryUrl: repositoryUrl,
			Username:      a.ContainerEnv.RegistryUser.Elem(),
			Password:      a.ContainerEnv.RegistryPass.Elem(),
			Server:        a.ContainerEnv.Registry.LoginServer,
		})
		if err != nil {
			return err
		}

		containerGroupSubnet, err := network.NewSubnet(ctx, "container-group-subnet", &network.SubnetArgs{
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
		})
		if err != nil {
			return errors.WithMessage(err, "creating container group subnet")
		}

		containerGroupName := fmt.Sprintf("%s-migration-group", name)

		_, err = containerinstance.NewContainerGroup(ctx, containerGroupName, &containerinstance.ContainerGroupArgs{
			ContainerGroupName: pulumi.String(containerGroupName),
			Containers: containerinstance.ContainerArray{
				&containerinstance.ContainerArgs{
					Image: image.URI(),
					Name:  pulumi.Sprintf("%s-migration", name),
					Resources: &containerinstance.ResourceRequirementsArgs{
						Requests: &containerinstance.ResourceRequestsArgs{
							Cpu:        pulumi.Float64(1),
							MemoryInGB: pulumi.Float64(1),
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
					Id:   containerGroupSubnet.ID(),
					Name: a.VirtualNetwork.Name,
				},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}
