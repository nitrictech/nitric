package terraform

import (
	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func (td *TerraformDeployment) createVariablesForIntent(intentName string, spec *ResourceBlueprint) {
	for varName, variable := range spec.Variables {
		if td.instancedTerraformVariables[intentName] == nil {
			td.instancedTerraformVariables[intentName] = make(map[string]cdktf.TerraformVariable)
		}

		td.instancedTerraformVariables[intentName][varName] = cdktf.NewTerraformVariable(td.stack, jsii.Sprintf("%s_%s", intentName, varName), &cdktf.TerraformVariableConfig{
			Description: jsii.String(variable.Description),
			Type:        jsii.String(variable.Type),
			Default:     variable.Default,
		})
	}
}

func (td *TerraformDeployment) createPlatformVariables() {
	for varName, variableSpec := range td.engine.platform.Variables {
		td.terraformVariables[varName] = cdktf.NewTerraformVariable(td.stack, jsii.String(varName), &cdktf.TerraformVariableConfig{
			Description: jsii.String(variableSpec.Description),
			Default:     variableSpec.Default,
			Type:        jsii.String(variableSpec.Type),
		})
	}
}