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

package pulumix

import (
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NitricPulumiResource[T any] struct {
	Id     *resourcespb.ResourceIdentifier
	Config T
}

type NitricPulumiServiceResource = NitricPulumiResource[*NitricPulumiServiceConfig]

type NitricPulumiServiceConfig struct {
	*deploymentspb.Service

	// Allow for pulumi strings to be used as service environment variables
	env pulumi.StringMap
}

func (n *NitricPulumiServiceConfig) SetEnv(key string, value pulumi.StringInput) {
	if n.env == nil {
		n.env = pulumi.StringMap{}
	}

	n.env[key] = value
}

func (n *NitricPulumiServiceConfig) Env() pulumi.StringMap {
	envMap := pulumi.StringMap{}

	for k, v := range n.Service.GetEnv() {
		envMap[k] = pulumi.String(v)
	}

	for k, v := range n.env {
		envMap[k] = v
	}

	return envMap
}
