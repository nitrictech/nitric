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
	"fmt"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricGcpTerraformProvider) SqlDatabase(stack cdktf.TerraformStack, name string, config *deploymentspb.SqlDatabase) error {
	// Return gRPC unimplemented error
	return fmt.Errorf("sql databases aren't yet supported by Nitric using the Terraform Google Cloud provider, remove the sql databases from your project and try again or use an alternate provider")
}
