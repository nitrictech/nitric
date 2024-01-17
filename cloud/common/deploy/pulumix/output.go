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
	"fmt"
	"strings"

	deploy "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pterm/pterm"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func PulumiOutputsToResult(outputs auto.OutputMap) *deploy.DeploymentUpEvent {
	apis := auto.OutputMap{}
	for k, v := range outputs {
		// Don't output secrets
		if v.Secret {
			continue
		}

		if strings.HasPrefix(k, "api:") {
			apis[strings.TrimPrefix(k, "api:")] = v
		}
	}

	rows := [][]string{{"API", "Endpoint"}}
	for k, v := range apis {
		rows = append(rows, []string{k, fmt.Sprint(v.Value)})
	}
	table, _ := pterm.DefaultTable.WithData(rows).Srender()

	return &deploy.DeploymentUpEvent{
		Content: &deploy.DeploymentUpEvent_Result{
			Result: &deploy.UpResult{
				Success: true,
				Details: "\n" + table + "\n",
			},
		},
	}
}
