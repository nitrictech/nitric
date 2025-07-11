package terraform

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aws/jsii-runtime-go"
	random "github.com/cdktf/cdktf-provider-random-go/random/v11/provider"
	"github.com/cdktf/cdktf-provider-random-go/random/v11/stringresource"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	app_spec_schema "github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines"
	"github.com/nitrictech/nitric/server/plugin"
)

type TerraformEngine struct {
	platform   *PlatformSpec
	repository PluginRepository
	output     io.Writer

	outputDir string
}

type TerraformDeployment struct {
	app     cdktf.App
	stack   cdktf.TerraformStack
	stackId stringresource.StringResource
	engine  *TerraformEngine

	serviceIdentities map[string]map[string]interface{}

	terraformResources          map[string]cdktf.TerraformHclModule
	terraformInfraResources     map[string]cdktf.TerraformHclModule
	terraformVariables          map[string]cdktf.TerraformVariable
	instancedTerraformVariables map[string]map[string]cdktf.TerraformVariable
}

type SpecReference struct {
	// var/infra/etc
	Source string
	// simple key for var or path for infra e.g. vpc.arn
	Path []string
}

func (e *TerraformEngine) resolveIdentityPlugin(blueprint *ResourceBlueprint) (*IdentityPluginManifest, error) {
	pluginRef, err := blueprint.ResolvePlugin(e.platform)
	if err != nil {
		return nil, err
	}
	return e.repository.GetIdentityPlugin(pluginRef.Library.Team, pluginRef.Library.Name, pluginRef.Library.Version, pluginRef.Name)
}

func (e *TerraformEngine) resolvePlugin(blueprint *ResourceBlueprint) (*ResourcePluginManifest, error) {
	pluginRef, err := blueprint.ResolvePlugin(e.platform)
	if err != nil {
		return nil, err
	}
	return e.repository.GetResourcePlugin(pluginRef.Library.Team, pluginRef.Library.Name, pluginRef.Library.Version, pluginRef.Name)
}

func (e *TerraformEngine) GetPluginManifestsForType(typ string) (map[string]*ResourcePluginManifest, error) {
	manifests := map[string]*ResourcePluginManifest{}

	blueprints, err := e.platform.GetResourceBlueprintsForType(typ)
	if err != nil {
		return nil, err
	}

	for blueprintIntent, blueprint := range blueprints {
		plug, err := e.resolvePlugin(blueprint)
		if err != nil {
			return nil, err
		}
		manifests[blueprintIntent] = plug
	}

	return manifests, nil
}

func SpecReferenceFromToken(token string) (*SpecReference, error) {
	contents, ok := extractTokenContents(token)
	if !ok {
		return nil, fmt.Errorf("invalid token format")
	}

	parts := strings.Split(contents, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid token format")
	}

	return &SpecReference{
		Source: parts[0],
		Path:   parts[1:],
	}, nil
}

func (tf *TerraformDeployment) resolveDependencies(resource *ResourceBlueprint, module cdktf.TerraformHclModule) error {
	if len(resource.DependsOn) == 0 {
		return nil
	}

	dependsOnResources := []*string{}
	for _, dependsOn := range resource.DependsOn {
		specRef, err := SpecReferenceFromToken(dependsOn)
		if err != nil {
			return err
		}

		if specRef.Source != "infra" {
			return fmt.Errorf("depends_on can only reference infra resources")
		}

		if _, ok := tf.terraformInfraResources[specRef.Path[0]]; !ok {
			// Infra resource has not already been resolved so this dependency won't exist at all
			continue
		}

		moduleId := fmt.Sprintf("module.%s", *tf.terraformInfraResources[specRef.Path[0]].Node().Id())
		dependsOnResources = append(dependsOnResources, jsii.String(moduleId))
	}
	module.SetDependsOn(&dependsOnResources)
	return nil
}

func (tf *TerraformDeployment) resolveInfraResource(infraName string) (cdktf.TerraformHclModule, error) {
	resource, ok := tf.engine.platform.InfraSpecs[infraName]
	if !ok {
		return nil, fmt.Errorf("infra resource %s not found", infraName)
	}

	if _, ok := tf.terraformInfraResources[infraName]; !ok {
		pluginRef, err := tf.engine.resolvePlugin(resource)
		if err != nil {
			return nil, err
		}

		tf.terraformInfraResources[infraName] = cdktf.NewTerraformHclModule(tf.stack, jsii.String(infraName), &cdktf.TerraformHclModuleConfig{
			Source: jsii.String(pluginRef.Deployment.Terraform),
		})
	}

	return tf.terraformInfraResources[infraName], nil
}

func (tf *TerraformDeployment) resolveValue(intentName string, value interface{}) (interface{}, error) {
	switch v := value.(type) {
	case string:
		// Handle token resolution for strings
		specRef, err := SpecReferenceFromToken(v)
		if err != nil {
			return v, nil // Return original value if not a token
		}

		if specRef.Source == "infra" {
			refName := specRef.Path[0]
			propertyName := specRef.Path[1]

			infraResource, err := tf.resolveInfraResource(refName)
			if err != nil {
				return nil, err
			}
			// map the variable output to the infra resource
			return infraResource.Get(jsii.String(propertyName)), nil
		} else if specRef.Source == "self" {
			tfVariable, ok := tf.instancedTerraformVariables[intentName][specRef.Path[0]]
			if !ok {
				return nil, fmt.Errorf("Variable %s does not exist for provided blueprint", specRef.Path[0])
			}
			return tfVariable.Value(), nil
		} else if specRef.Source == "var" {
			tfVariable, ok := tf.terraformVariables[specRef.Path[0]]
			if !ok {
				return nil, fmt.Errorf("Variable %s does not exist for this platform", specRef.Path[0])
			}
			return tfVariable.Value(), nil
		}
		return v, nil

	case map[string]interface{}:
		// Recursively process map values
		result := make(map[string]interface{})
		for key, val := range v {
			resolvedVal, err := tf.resolveValue(intentName, val)
			if err != nil {
				return nil, err
			}
			result[key] = resolvedVal
		}
		return result, nil

	case []interface{}:
		// Recursively process slice values
		result := make([]interface{}, len(v))
		for i, val := range v {
			resolvedVal, err := tf.resolveValue(intentName, val)
			if err != nil {
				return nil, err
			}
			result[i] = resolvedVal
		}
		return result, nil

	default:
		// Return primitive values as is
		return v, nil
	}
}

func (tf *TerraformDeployment) resolveTokensForModule(intentName string, resource *ResourceBlueprint, module cdktf.TerraformHclModule) error {
	for property, value := range resource.Properties {
		resolvedValue, err := tf.resolveValue(intentName, value)
		if err != nil {
			return err
		}
		module.Set(jsii.String(property), resolvedValue)
	}

	return nil
}

func NewTerraformDeployment(engine *TerraformEngine, stackName string) *TerraformDeployment {
	app := cdktf.NewApp(&cdktf.AppConfig{
		Outdir: jsii.String(engine.outputDir),
	})
	stack := cdktf.NewTerraformStack(app, jsii.String(stackName))

	random.NewRandomProvider(stack, jsii.String("random"), &random.RandomProviderConfig{})

	stackId := stringresource.NewStringResource(stack, jsii.String("stack_id"), &stringresource.StringResourceConfig{
		Length:  jsii.Number(8),
		Upper:   jsii.Bool(false),
		Lower:   jsii.Bool(true),
		Numeric: jsii.Bool(false),
		Special: jsii.Bool(false),
	})

	return &TerraformDeployment{
		app:                         app,
		stack:                       stack,
		stackId:                     stackId,
		engine:                      engine,
		terraformResources:          map[string]cdktf.TerraformHclModule{},
		terraformInfraResources:     map[string]cdktf.TerraformHclModule{},
		terraformVariables:          map[string]cdktf.TerraformVariable{},
		instancedTerraformVariables: map[string]map[string]cdktf.TerraformVariable{},
		serviceIdentities:           map[string]map[string]interface{}{},
	}
}

func (e *TerraformEngine) resolvePluginsForService(servicePlugin *ResourcePluginManifest) (*plugin.PluginDefinition, error) {
	// TODO: Map platform resource plugins to the service plugin
	gets := []string{}

	pluginDef := &plugin.PluginDefinition{
		Service: plugin.GoPlugin{
			Alias:  "svcPlugin",
			Name:   "default",
			Import: strings.Split(servicePlugin.Runtime.GoModule, "@")[0],
		},
	}
	gets = append(gets, servicePlugin.Runtime.GoModule)

	// FIXME: This add all storage plugins without regard to actually requiring access
	storagePlugins, err := e.GetPluginManifestsForType("bucket")
	if err != nil {
		return nil, err
	}

	// Add storage plugins to the runtime
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

var entrypointOriginTypes = []string{"service", "bucket"}

func (e *TerraformDeployment) resolveEntrypointNitricVar(name string, appSpec *app_spec_schema.Application, spec *app_spec_schema.EntrypointIntent) (interface{}, error) {
	origins := map[string]interface{}{}
	for path, route := range spec.Routes {
		intentTarget, ok := appSpec.GetResourceIntent(route.TargetName)
		if !ok {
			return nil, fmt.Errorf("target %s not found", route.TargetName)
		}

		var intentTargetType string
		switch intentTarget.(type) {
		case *app_spec_schema.ServiceIntent:
			intentTargetType = "service"
		case *app_spec_schema.BucketIntent:
			intentTargetType = "bucket"
		default:
			return nil, fmt.Errorf("target %s is not a service or bucket", route.TargetName)
		}

		hclTarget, ok := e.terraformResources[route.TargetName]
		if !ok {
			return nil, fmt.Errorf("target %s not found", route.TargetName)
		}

		domainNameNitricVar := hclTarget.Get(jsii.String("nitric.domain_name"))
		idNitricVar := hclTarget.Get(jsii.String("nitric.id"))
		resourcesNitricVar := hclTarget.Get(jsii.String("nitric.exports.resources"))

		origins[route.TargetName] = map[string]interface{}{
			"path":      jsii.String(path),
			"base_path": jsii.String(route.BasePath),
			"type":      jsii.String(intentTargetType),
			"id":        idNitricVar,
			// Assume the output var has a http_endpoint property
			"domain_name": domainNameNitricVar,
			"resources":   resourcesNitricVar,
		}
	}

	// Build the origins map
	nitricVar := map[string]interface{}{
		"name":     jsii.String(name),
		"stack_id": e.stackId.Result(),
		"origins":  origins,
	}

	return nitricVar, nil
}

func (e *TerraformDeployment) resolveService(name string, spec *app_spec_schema.ServiceIntent, resourceSpec *ServiceBlueprint, plug *ResourcePluginManifest) (*NitricServiceVariables, error) {
	// Map the nitric variable
	var imageVars *map[string]interface{} = nil

	pluginManifest, err := e.engine.resolvePluginsForService(plug)
	if err != nil {
		return nil, err
	}

	pluginManifestBytes, err := json.Marshal(pluginManifest)
	if err != nil {
		return nil, err
	}

	var schedules map[string]NitricServiceSchedule = nil
	if len(schedules) > 0 && !slices.Contains(plug.Capabilities, "schedules") {
		return nil, fmt.Errorf("service %s has schedules but the plugin %s does not support schedules", name, plug.Name)
	} else {
		schedules = map[string]NitricServiceSchedule{}
	}

	for triggerName, trigger := range spec.Triggers {
		if trigger.Schedule == nil {
			continue
		}

		schedules[triggerName] = NitricServiceSchedule{
			CronExpression: jsii.String(trigger.Schedule.CronExpression),
			Path:           jsii.String(trigger.Path),
		}
	}

	if spec.Container.Image != nil {
		imageVars = &map[string]interface{}{
			"image_id": jsii.String(spec.Container.Image.ID),
			"tag":      jsii.String(name),
			"args":     map[string]*string{"PLUGIN_DEFINITION": jsii.String(string(pluginManifestBytes))},
		}
	} else if spec.Container.Docker != nil {
		// merge args
		args := map[string]*string{"PLUGIN_DEFINITION": jsii.String(string(pluginManifestBytes))}
		for k, v := range spec.Container.Docker.Args {
			args[k] = jsii.String(v)
		}

		imageVars = &map[string]interface{}{
			"build_context": jsii.String(spec.Container.Docker.Context),
			"dockerfile":    jsii.String(spec.Container.Docker.Dockerfile),
			"tag":           jsii.String(name),
			"args":          args,
		}
	}

	imageModule := cdktf.NewTerraformHclModule(e.stack, jsii.Sprintf("%s_image", name), &cdktf.TerraformHclModuleConfig{
		Source:    jsii.String(imageModuleRef),
		Variables: imageVars,
	})

	identityModuleOutputs := map[string]interface{}{}
	for _, id := range resourceSpec.Identities {
		identityPlugin, err := e.engine.resolveIdentityPlugin(&id)
		if err != nil {
			return nil, err
		}

		idModule := cdktf.NewTerraformHclModule(e.stack, jsii.Sprintf("%s_%s_role", name, identityPlugin.Name), &cdktf.TerraformHclModuleConfig{
			Source: jsii.String(identityPlugin.Deployment.Terraform),
			// TODO: Properly resolve tokens here
			Variables: &id.Properties,
		})

		idModule.Set(jsii.String("nitric"), map[string]interface{}{
			"name":     jsii.String(name),
			"stack_id": e.stackId.Result(),
		})

		identityModuleOutputs[identityPlugin.IdentityType] = idModule.Get(jsii.String("nitric"))
	}

	for _, requiredIdentity := range plug.RequiredIdentities {
		providedIdentities := slices.Collect(maps.Keys(identityModuleOutputs))
		if ok := slices.Contains(providedIdentities, requiredIdentity); !ok {
			return nil, fmt.Errorf("service %s is missing identity %s, required by plugin %s, provided identities were %s", name, requiredIdentity, plug.Name, providedIdentities)
		}
	}

	// Create this services identities
	nitricVar := &NitricServiceVariables{
		NitricVariables: NitricVariables{
			Name: jsii.String(name),
		},
		Schedules:  &schedules,
		ImageId:    imageModule.GetString(jsii.String("image_id")),
		Identities: &identityModuleOutputs,
		StackId:    e.stackId.Result(),
		Env:        &spec.Env,
	}

	e.serviceIdentities[name] = identityModuleOutputs

	return nitricVar, nil
}

func (e *TerraformDeployment) createVariablesForIntent(intentName string, spec *ResourceBlueprint) {
	for varName, variable := range spec.Variables {
		if e.instancedTerraformVariables[intentName] == nil {
			e.instancedTerraformVariables[intentName] = make(map[string]cdktf.TerraformVariable)
		}

		e.instancedTerraformVariables[intentName][varName] = cdktf.NewTerraformVariable(e.stack, jsii.Sprintf("%s_%s", intentName, varName), &cdktf.TerraformVariableConfig{
			Description: jsii.String(variable.Description),
			Type:        jsii.String(variable.Type),
			// TODO: Possibly resolve a token?
			Default: variable.Default,
		})
	}
}

// Apply the engine to the target environment
func (e *TerraformEngine) Apply(appSpec *app_spec_schema.Application) error {
	tfDeployment := NewTerraformDeployment(e, appSpec.Name)

	// Create platform variables ahead of time
	for varName, variableSpec := range e.platform.Variables {
		tfDeployment.terraformVariables[varName] = cdktf.NewTerraformVariable(tfDeployment.stack, jsii.String(varName), &cdktf.TerraformVariableConfig{
			Description: jsii.String(variableSpec.Description),
			Default:     variableSpec.Default,
			Type:        jsii.String(variableSpec.Type),
		})
	}

	// Prepare service inputs and identities
	resourceIntents := appSpec.GetResourceIntents()

	serviceInputs := map[string]*NitricServiceVariables{}
	for intentName, resourceIntent := range resourceIntents {
		var serviceIntent *app_spec_schema.ServiceIntent
		switch resourceIntent.(type) {
		case *app_spec_schema.ServiceIntent:
			serviceIntent = resourceIntent.(*app_spec_schema.ServiceIntent)
		default:
			continue
		}

		spec, err := e.platform.GetServiceBlueprint(serviceIntent.GetSubType())
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", serviceIntent.GetType(), serviceIntent.GetSubType(), err)
		}
		plug, err := e.resolvePlugin(spec.ResourceBlueprint)
		if err != nil {
			return fmt.Errorf("could not find plugin %s: %v", spec.PluginId, err)
		}

		nitricVar, err := tfDeployment.resolveService(intentName, serviceIntent, spec, plug)
		if err != nil {
			return err
		}

		serviceInputs[intentName] = nitricVar
	}

	serviceEnvs := map[string][]interface{}{}

	// Resolve non-service/non-entrypoint/non-bucket modules
	for intentName, resourceIntent := range resourceIntents {
		resourceType := resourceIntent.GetType()

		if resourceType == "service" || resourceType == "entrypoint" || resourceType == "bucket" {
			continue
		}

		spec, err := e.platform.GetResourceBlueprint(resourceIntent.GetType(), resourceIntent.GetSubType())
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.GetType(), resourceIntent.GetSubType(), err)
		}
		plug, err := e.resolvePlugin(spec)
		if err != nil {
			return fmt.Errorf("could not find plugin %s: %v", spec.PluginId, err)
		}

		servicesInput := map[string]any{}
		if access, ok := resourceIntent.GetAccess(); ok {
			for serviceName, actions := range access {
				idMap, ok := tfDeployment.serviceIdentities[serviceName]
				if !ok {
					return fmt.Errorf("could not give access to %s %s: service %s not found", resourceIntent.GetType(), intentName, serviceName)
				}

				servicesInput[serviceName] = map[string]interface{}{
					"actions":    jsii.Strings(actions...),
					"identities": idMap,
				}
			}
		}

		nitricVar := map[string]any{
			"name":     intentName,
			"stack_id": tfDeployment.stackId.Result(),
			"services": servicesInput,
		}

		// Create terraform variables for intent for a spec
		tfDeployment.createVariablesForIntent(intentName, spec)

		tfDeployment.terraformResources[intentName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})

		for serviceName, _ := range serviceInputs {
			env := cdktf.Fn_Try(&[]interface{}{tfDeployment.terraformResources[intentName].Get(jsii.Sprintf("nitric.exports.services.%s.env", serviceName)), map[string]interface{}{}})
			serviceEnvs[serviceName] = append(serviceEnvs[serviceName], env)
		}
	}

	// Resolve bucket modules
	for intentName, bucketIntent := range appSpec.BucketIntents {

		contentPath := ""
		if bucketIntent != nil {
			contentPath = bucketIntent.ContentPath
		}

		spec, err := e.platform.GetResourceBlueprint(bucketIntent.GetType(), bucketIntent.GetSubType())
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", bucketIntent.GetType(), bucketIntent.GetSubType(), err)
		}
		plug, err := e.resolvePlugin(spec)
		if err != nil {
			return fmt.Errorf("could not find plugin %s: %v", spec.PluginId, err)
		}

		servicesInput := map[string]any{}
		if access, ok := bucketIntent.GetAccess(); ok {
			for serviceName, actions := range access {
				idMap, ok := tfDeployment.serviceIdentities[serviceName]
				if !ok {
					return fmt.Errorf("could not give access to bucket %s: service %s not found", intentName, serviceName)
				}

				servicesInput[serviceName] = map[string]interface{}{
					"actions":    jsii.Strings(actions...),
					"identities": idMap,
				}
			}
		}

		nitricVar := map[string]any{
			"name":         intentName,
			"stack_id":     tfDeployment.stackId.Result(),
			"content_path": contentPath,
			"services":     servicesInput,
		}

		// Create terraform variables for intent for a spec
		tfDeployment.createVariablesForIntent(intentName, spec)

		tfDeployment.terraformResources[intentName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})

		for serviceName, _ := range serviceInputs {
			env := cdktf.Fn_Try(&[]interface{}{tfDeployment.terraformResources[intentName].Get(jsii.Sprintf("nitric.exports.services.%s.env", serviceName)), map[string]interface{}{}})
			serviceEnvs[serviceName] = append(serviceEnvs[serviceName], env)
		}
	}

	// Resolve service modules
	for intentName, resourceIntent := range appSpec.ServiceIntents {
		spec, err := e.platform.GetResourceBlueprint(resourceIntent.GetType(), resourceIntent.GetSubType())
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.GetType(), resourceIntent.GetSubType(), err)
		}
		plug, err := e.resolvePlugin(spec)
		if err != nil {
			return fmt.Errorf("could not find plugin %s: %v", spec.PluginId, err)
		}

		var nitricVar *NitricServiceVariables = serviceInputs[intentName]

		origEnv := nitricVar.Env

		mergedEnv := serviceEnvs[intentName]
		allEnv := append(mergedEnv, origEnv)

		nitricVar.Env = cdktf.Fn_Merge(&allEnv)

		tfDeployment.createVariablesForIntent(intentName, spec)

		tfDeployment.terraformResources[intentName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})
	}

	// Resolve entrypoint modules
	for intentName, resourceIntent := range appSpec.EntrypointIntents {
		spec, err := e.platform.GetResourceBlueprint(resourceIntent.GetType(), resourceIntent.GetSubType())
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.GetType(), resourceIntent.GetSubType(), err)
		}
		plug, err := e.resolvePlugin(spec)
		if err != nil {
			return fmt.Errorf("could not find plugin %s: %v", spec.PluginId, err)
		}

		nitricVar, err := tfDeployment.resolveEntrypointNitricVar(intentName, appSpec, resourceIntent)
		if err != nil {
			return err
		}

		tfDeployment.createVariablesForIntent(intentName, spec)

		tfDeployment.terraformResources[intentName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})
	}

	// Resolve resource tokens
	for resourceName, resourceIntent := range resourceIntents {
		resourceSpec, err := e.platform.GetResourceBlueprint(resourceIntent.GetType(), resourceIntent.GetSubType())
		if err != nil {
			return err
		}

		err = tfDeployment.resolveTokensForModule(resourceName, resourceSpec, tfDeployment.terraformResources[resourceName])
		if err != nil {
			return err
		}
	}

	// Resolve infra tokens
	for infraName, infra := range tfDeployment.terraformInfraResources {
		infraSpec, ok := e.platform.InfraSpecs[infraName]
		if !ok {
			return fmt.Errorf("infra resource %s not found", infraName)
		}
		// TODO: This is overloading this method as infra-name is not usable in this context as infra cannot resolve `self` tokens
		err := tfDeployment.resolveTokensForModule(infraName, infraSpec, infra)
		if err != nil {
			return err
		}
	}

	// resolve dependencies for all created modules
	for resourceName, resource := range resourceIntents {
		resourceSpec, err := e.platform.GetResourceBlueprint(resource.GetType(), resource.GetSubType())
		if err != nil {
			return err
		}

		err = tfDeployment.resolveDependencies(resourceSpec, tfDeployment.terraformResources[resourceName])
		if err != nil {
			return err
		}
	}

	// resolve dependencies for all created infrastructure
	for infraName, infra := range tfDeployment.terraformInfraResources {
		infraSpec, ok := e.platform.InfraSpecs[infraName]
		if !ok {
			return fmt.Errorf("infra resource %s not found", infraName)
		}

		err := tfDeployment.resolveDependencies(infraSpec, infra)
		if err != nil {
			return err
		}
	}

	tfDeployment.app.Synth()

	fmt.Fprintf(e.output, "✔ Terraform successfully output to %s\n", filepath.Join(e.outputDir, "stacks", appSpec.Name))

	stackPath := filepath.Join(e.outputDir, "stacks", appSpec.Name)

	fmt.Fprintf(e.output, "\nNext steps:\n")
	fmt.Fprintf(e.output, "1. Run `cd %s && terraform init` to initialize the stack\n", stackPath)
	fmt.Fprintf(e.output, "2. Run `terraform plan` to preview the stack\n")
	fmt.Fprintf(e.output, "3. Run `terraform apply` to deploy the stack\n\n")

	return nil
}

var _ engines.Engine = &TerraformEngine{}

type terraformEngineOption func(*TerraformEngine)

func WithRepository(repository PluginRepository) terraformEngineOption {
	return func(engine *TerraformEngine) {
		engine.repository = repository
	}
}

func WithOutputDir(outputDir string) terraformEngineOption {
	return func(engine *TerraformEngine) {
		engine.outputDir = outputDir
	}
}

func WithOutput(output io.Writer) terraformEngineOption {
	return func(engine *TerraformEngine) {
		engine.output = output
	}
}

func NewFromFile(platformFile io.Reader, opts ...terraformEngineOption) *TerraformEngine {
	platform := &PlatformSpec{}

	json.NewDecoder(platformFile).Decode(platform)

	return New(platform, opts...)
}

func New(platformSpec *PlatformSpec, opts ...terraformEngineOption) *TerraformEngine {
	engine := &TerraformEngine{
		platform:  platformSpec,
		outputDir: "terraform",
		output:    io.Discard,
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine
}
