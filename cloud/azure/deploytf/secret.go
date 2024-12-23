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
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// // Secret - Deploy a Secret
func (a *NitricAzureTerraformProvider) Secret(stack cdktf.TerraformStack, name string, config *deploymentspb.Secret) error {
	// Secrets in Azure Key Vaults are unique resources created during deployment, so we don't need to do anything here.
	// Instead, if at least one secret is requested a Key Vault will be created for the stack.
	// Policies are also created which restrict access to the Key Vault and the named secrets inside.
	return nil
}
