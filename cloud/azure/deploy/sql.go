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

	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-azure-native-sdk/dbforpostgresql/v2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// func createAzureContainerRepository(ctx *pulumi.Context, name string, region string) (*containerregistry.Registry, error) {
// 	acr, err := containerregistry.NewRegistry(ctx, name, &containerregistry.RegistryArgs{
// 		Location: pulumi.String(region),
// 	})
// }

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

	return nil
}
