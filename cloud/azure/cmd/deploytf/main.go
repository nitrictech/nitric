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

package main

import (
	"github.com/nitrictech/nitric/cloud/azure/common/runtime"
	"github.com/nitrictech/nitric/cloud/azure/deploytf"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
)

// Start the deployment server
func main() {
	// os.Setenv("JSII_SILENCE_WARNING_UNTESTED_NODE_VERSION", "true")
	// os.Setenv("SYNTH_HCL_OUTPUT", "true")
	azureStack := deploytf.NewNitricAzureProvider()

	providerServer := provider.NewTerraformProviderServer(azureStack, runtime.NitricAzureRuntime)

	// Start the terraform provider server
	providerServer.Start()
}
