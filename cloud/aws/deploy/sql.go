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
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Sqldatabase - Implements PostgresSql database deployments use AWS Aurora
func (a *NitricAwsPulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {

	// Run CREATE DATABASE queries against this cluster for each of the databases we want to deploy here
	// We will run with a IF NOT EXISTS clause to ensure we don't overwrite existing databases
	// a.DatabaseCluster

	return fmt.Errorf("Sql databases are unimplemented in the nitric AWS provider")
}
