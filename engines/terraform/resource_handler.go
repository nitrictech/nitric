package terraform

import (
	"fmt"
	"strings"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	app_spec_schema "github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/server/plugin"
)

func (e *TerraformEngine) resolvePluginsForService(servicePlugin *ResourcePluginManifest) (*plugin.PluginDefinition, error) {
	gets := []string{}

	// Check if Runtime is nil to prevent panic
	if servicePlugin.Runtime == nil {
		return nil, fmt.Errorf("service plugin %s has no runtime configuration", servicePlugin.Name)
	}

	pluginDef := &plugin.PluginDefinition{
		Service: plugin.GoPlugin{
			Alias:  "svcPlugin",
			Name:   "default",
			Import: strings.Split(servicePlugin.Runtime.GoModule, "@")[0],
		},
	}
	gets = append(gets, servicePlugin.Runtime.GoModule)

	storagePlugins, err := e.GetPluginManifestsForType("bucket")
	if err != nil {
		return nil, err
	}

	for name, plug := range storagePlugins {
		pluginDef.Storage = append(pluginDef.Storage, plugin.GoPlugin{
			Alias:  fmt.Sprintf("storage_%s", name),
			Name:   name,
			Import: strings.Split(plug.Runtime.GoModule, "@")[0],
		})
		gets = append(gets, plug.Runtime.GoModule)
	}

	pluginDef.Gets = gets

	return pluginDef, nil
}

func (td *TerraformDeployment) processServiceIdentities(appSpec *app_spec_schema.Application) (map[string]*NitricServiceVariables, error) {
	serviceInputs := map[string]*NitricServiceVariables{}

	for intentName, serviceIntent := range appSpec.ServiceIntents {
		spec, err := td.engine.platform.GetServiceBlueprint(serviceIntent.GetSubType())
		if err != nil {
			return nil, fmt.Errorf("could not find platform type for %s.%s: %w", serviceIntent.GetType(), serviceIntent.GetSubType(), err)
		}
		plug, err := td.engine.resolvePlugin(spec.ResourceBlueprint)
		if err != nil {
			return nil, err
		}

		nitricVar, err := td.resolveService(intentName, serviceIntent, spec, plug)
		if err != nil {
			return nil, err
		}

		serviceInputs[intentName] = nitricVar
	}

	return serviceInputs, nil
}

func (td *TerraformDeployment) processServiceResources(appSpec *app_spec_schema.Application, serviceInputs map[string]*NitricServiceVariables, serviceEnvs map[string][]interface{}) error {
	for intentName, serviceIntent := range appSpec.ServiceIntents {
		spec, err := td.engine.platform.GetResourceBlueprint(serviceIntent.GetType(), serviceIntent.GetSubType())
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", serviceIntent.GetType(), serviceIntent.GetSubType(), err)
		}
		plug, err := td.engine.resolvePlugin(spec)
		if err != nil {
			return err
		}

		nitricVar := serviceInputs[intentName]
		origEnv := nitricVar.Env

		mergedEnv := serviceEnvs[intentName]
		allEnv := append(mergedEnv, origEnv)
		nitricVar.Env = cdktf.Fn_Merge(&allEnv)

		td.createVariablesForIntent(intentName, spec)

		td.terraformResources[intentName] = cdktf.NewTerraformHclModule(td.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})
	}

	return nil
}

func (td *TerraformDeployment) processBucketResources(appSpec *app_spec_schema.Application) (map[string][]interface{}, error) {
	serviceEnvs := map[string][]interface{}{}

	for intentName, bucketIntent := range appSpec.BucketIntents {
		contentPath := ""
		if bucketIntent != nil {
			contentPath = bucketIntent.ContentPath
		}

		spec, err := td.engine.platform.GetResourceBlueprint(bucketIntent.GetType(), bucketIntent.GetSubType())
		if err != nil {
			return nil, fmt.Errorf("could not find platform type for %s.%s: %w", bucketIntent.GetType(), bucketIntent.GetSubType(), err)
		}
		plug, err := td.engine.resolvePlugin(spec)
		if err != nil {
			return nil, err
		}

		servicesInput := map[string]any{}
		if access, ok := bucketIntent.GetAccess(); ok {
			for serviceName, actions := range access {
				expandedActions := app_spec_schema.ExpandActions(actions, app_spec_schema.Bucket)

				idMap, ok := td.serviceIdentities[serviceName]
				if !ok {
					return nil, fmt.Errorf("could not give access to bucket %s: service %s not found", intentName, serviceName)
				}

				servicesInput[serviceName] = map[string]interface{}{
					"actions":    jsii.Strings(expandedActions...),
					"identities": idMap,
				}
			}
		}

		nitricVar := map[string]any{
			"name":         intentName,
			"stack_id":     td.stackId.Result(),
			"content_path": contentPath,
			"services":     servicesInput,
		}

		td.createVariablesForIntent(intentName, spec)

		td.terraformResources[intentName] = cdktf.NewTerraformHclModule(td.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})

		// Collect environment variables that buckets export to services
		for serviceName := range td.serviceIdentities {
			env := cdktf.Fn_Try(&[]interface{}{td.terraformResources[intentName].Get(jsii.Sprintf("nitric.exports.services.%s.env", serviceName)), map[string]interface{}{}})
			serviceEnvs[serviceName] = append(serviceEnvs[serviceName], env)
		}
	}

	return serviceEnvs, nil
}

func (td *TerraformDeployment) processEntrypointResources(appSpec *app_spec_schema.Application) error {
	for intentName, entrypointIntent := range appSpec.EntrypointIntents {
		spec, err := td.engine.platform.GetResourceBlueprint(entrypointIntent.GetType(), entrypointIntent.GetSubType())
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", entrypointIntent.GetType(), entrypointIntent.GetSubType(), err)
		}
		plug, err := td.engine.resolvePlugin(spec)
		if err != nil {
			return err
		}

		nitricVar, err := td.resolveEntrypointNitricVar(intentName, appSpec, entrypointIntent)
		if err != nil {
			return err
		}

		td.createVariablesForIntent(intentName, spec)

		td.terraformResources[intentName] = cdktf.NewTerraformHclModule(td.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})
	}

	return nil
}

func (td *TerraformDeployment) processDatabaseResources(appSpec *app_spec_schema.Application) (map[string][]interface{}, error) {
	serviceEnvs := map[string][]interface{}{}

	for intentName, databaseIntent := range appSpec.DatabaseIntents {
		spec, err := td.engine.platform.GetResourceBlueprint(databaseIntent.GetType(), databaseIntent.GetSubType())
		if err != nil {
			return nil, fmt.Errorf("could not find platform type for %s.%s: %w", databaseIntent.GetType(), databaseIntent.GetSubType(), err)
		}
		plug, err := td.engine.resolvePlugin(spec)
		if err != nil {
			return nil, err
		}

		servicesInput := map[string]any{}
		if access, ok := databaseIntent.GetAccess(); ok {
			for serviceName, actions := range access {
				expandedActions := app_spec_schema.ExpandActions(actions, app_spec_schema.Database)

				idMap, ok := td.serviceIdentities[serviceName]
				if !ok {
					return nil, fmt.Errorf("could not give access to database %s: service %s not found", intentName, serviceName)
				}

				servicesInput[serviceName] = map[string]interface{}{
					"actions":    jsii.Strings(expandedActions...),
					"identities": idMap,
				}
			}
		}

		nitricVar := map[string]any{
			"name":        intentName,
			"stack_id":    td.stackId.Result(),
			"services":    servicesInput,
			"env_var_key": databaseIntent.EnvVarKey,
		}

		td.createVariablesForIntent(intentName, spec)

		td.terraformResources[intentName] = cdktf.NewTerraformHclModule(td.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})

		// Collect environment variables that databases export to services
		for serviceName := range td.serviceIdentities {
			env := cdktf.Fn_Try(&[]interface{}{td.terraformResources[intentName].Get(jsii.Sprintf("nitric.exports.services.%s.env", serviceName)), map[string]interface{}{}})
			serviceEnvs[serviceName] = append(serviceEnvs[serviceName], env)
		}
	}

	return serviceEnvs, nil
}
