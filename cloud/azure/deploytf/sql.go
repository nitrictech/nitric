// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploytf

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/sql"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (n *NitricAzureTerraformProvider) SqlDatabase(stack cdktf.TerraformStack, name string, config *deploymentspb.SqlDatabase) error {
	sqlConfig := &sql.SqlConfig{
		Name:                       jsii.String(name),
		StackName:                  n.Stack.StackNameOutput(),
		ResourceGroupName:          n.Stack.ResourceGroupNameOutput(),
		Location:                   jsii.String(n.Region),
		ImageRegistryServer:        n.Stack.RegistryLoginServerOutput(),
		ImageRegistryUsername:      n.Stack.RegistryUsernameOutput(),
		ImageRegistryPassword:      n.Stack.RegistryPasswordOutput(),
		ServerId:                   n.Stack.DatabaseServerIdOutput(),
		DatabaseServerFqdn:         n.Stack.DatabaseServerFqdnOutput(),
		DatabaseMasterPassword:     n.Stack.DatabaseMasterPasswordOutput(),
		MigrationContainerSubnetId: n.Stack.SubnetIdOutput(),
		DependsOn:                  &[]cdktf.ITerraformDependable{n.Stack},
	}

	if config.GetImageUri() != "" {
		inspect, err := image.CommandFromImageInspect(config.GetImageUri(), " ")
		if err != nil {
			return err
		}

		sqlConfig.MigrationImageUri = &inspect.ID
	}

	n.Databases[name] = sql.NewSql(stack, jsii.String(name), sqlConfig)

	return nil
}
