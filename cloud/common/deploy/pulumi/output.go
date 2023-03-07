package pulumi

import (
	"fmt"
	"strings"

	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pterm/pterm"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
)

func PulumiOutputsToResult(outputs auto.OutputMap) *deploy.DeployUpEvent {
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

	return &deploy.DeployUpEvent{
		Content: &deploy.DeployUpEvent_Result{
			Result: &deploy.DeployUpEventResult{
				Success: true,
				Result: &deploy.UpResult{
					Content: &deploy.UpResult_StringResult{
						StringResult: "\n" + table + "\n",
					},
				},
			},
		},
	}
}
