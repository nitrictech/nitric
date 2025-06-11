package terraform

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"maps"
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

func (e *TerraformEngine) GetPluginManifestsForType(typ string) (map[string]*ResourcePluginManifest, error) {
	manifests := map[string]*ResourcePluginManifest{}

	blueprints, err := e.platform.GetResourceBlueprintsForType(typ)
	if err != nil {
		return nil, err
	}

	for blueprintIntent, blueprint := range blueprints {
		plug, err := e.repository.GetResourcePlugin(blueprint.PluginId)
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

func (tf *TerraformDeployment) resolveTokensForModule(intentName string, resource *ResourceBlueprint, module cdktf.TerraformHclModule) error {
	for property, value := range resource.Properties {
		module.Set(jsii.String(property), value)

		// TODO: Recursive logic for mapping tokens
		token, ok := value.(string)
		if !ok {
			// Skip as its already been set
			continue
		}

		specRef, err := SpecReferenceFromToken(token)
		if err != nil {
			continue
		}

		if specRef.Source == "infra" {
			refName := specRef.Path[0]
			propertyName := specRef.Path[1]
			// map the variable output to the infra resource
			refProperty := tf.terraformInfraResources[refName].Get(jsii.String(propertyName))

			module.Set(jsii.String(property), refProperty)
		} else if specRef.Source == "self" {
			tfVariable, ok := tf.instancedTerraformVariables[intentName][specRef.Path[0]]
			if !ok {
				return fmt.Errorf("Variable %s does not exist for provided blueprint", specRef.Path[0])
			}

			module.Set(jsii.String(property), tfVariable.Value())
		} else if specRef.Source == "var" {
			tfVariable, ok := tf.terraformVariables[specRef.Path[0]]
			if !ok {
				return fmt.Errorf("Variable %s does not exist for this platform", specRef.Path[0])
			}

			// Create a new terraform variable
			module.Set(jsii.String(property), tfVariable.Value())
		}
	}

	return nil
}

func NewTerraformDeployment(engine *TerraformEngine, stackName string) *TerraformDeployment {
	app := cdktf.NewApp(&cdktf.AppConfig{})
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

func (e *TerraformEngine) resolvePluginsForService(servicePlugin *ResourcePluginManifest) (*plugin.PluginDefintion, error) {
	// TODO: Map platform resource plugins to the service plugin
	pluginDef := &plugin.PluginDefintion{
		Service: plugin.GoPlugin{
			Alias:  "svcPlugin",
			Name:   "default",
			Import: servicePlugin.Runtime.GoModule,
		},
	}

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
			Import: plug.Runtime.GoModule,
		})
	}

	return pluginDef, nil
}

var entrypointOriginTypes = []string{"service", "bucket"}

func (e *TerraformDeployment) resolveEntrypointNitricVar(name string, appSpec *app_spec_schema.Application, spec *app_spec_schema.EntrypointIntent) (interface{}, error) {
	origins := map[string]interface{}{}
	for path, route := range spec.Routes {
		intentTarget, ok := appSpec.ResourceIntents[route.TargetName]
		if !ok {
			return nil, fmt.Errorf("target %s not found", route.TargetName)
		}

		hclTarget, ok := e.terraformResources[route.TargetName]
		if !ok {
			return nil, fmt.Errorf("target %s not found", route.TargetName)
		}

		domainNameNitricVar := hclTarget.Get(jsii.String("nitric.domain_name"))
		idNitricVar := hclTarget.Get(jsii.String("nitric.id"))
		resourcesNitricVar := hclTarget.Get(jsii.String("nitric.exports.resources"))

		origins[route.TargetName] = map[string]interface{}{
			"path": jsii.String(path),
			"type": jsii.String(intentTarget.Type),
			"id":   idNitricVar,
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
		identityPlugin, err := e.engine.repository.GetIdentityPlugin(id.PluginId)
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
		ImageId:    imageModule.GetString(jsii.String("image_id")),
		Identities: &identityModuleOutputs,
		StackId:    e.stackId.Result(),
		Env:        &spec.Env,
	}

	e.serviceIdentities[name] = identityModuleOutputs

	return nitricVar, nil
}

func (e *TerraformDeployment) createVariablesForIntent(intentName string, intent app_spec_schema.Resource, spec *ResourceBlueprint) {
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

	// Create a terraform variable to establish the root context for application builds
	// this will be prepended to the path of any internal docker builds
	// tfDeployment.terraformVariables["build_root"] = cdktf.NewTerraformVariable(tfDeployment.stack, jsii.String("build_root"), &cdktf.TerraformVariableConfig{
	// 	Type: jsii.String("string"),
	// })

	// Create platform variables ahead of time
	for varName, variableSpec := range e.platform.Variables {
		tfDeployment.terraformVariables[varName] = cdktf.NewTerraformVariable(tfDeployment.stack, jsii.String(varName), &cdktf.TerraformVariableConfig{
			Description: jsii.String(variableSpec.Description),
			Default:     variableSpec.Default,
			Type:        jsii.String(variableSpec.Type),
		})
	}

	// Resolve infra modules
	for infraName, infra := range e.platform.InfraSpecs {
		plugin, err := e.repository.GetResourcePlugin(infra.PluginId)
		if err != nil {
			return err
		}

		tfDeployment.terraformInfraResources[infraName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(infraName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plugin.Deployment.Terraform),
		})
	}

	// Prepare service inputs and identities
	serviceInputs := map[string]*NitricServiceVariables{}
	for intentName, resourceIntent := range appSpec.ResourceIntents {
		if resourceIntent.Type != "service" {
			continue
		}

		spec, err := e.platform.GetServiceBlueprint(resourceIntent.SubType)
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.Type, resourceIntent.SubType, err)
		}
		plug, err := e.repository.GetResourcePlugin(spec.PluginId)
		if err != nil {
			return fmt.Errorf("could not find plugin %s", spec.PluginId)
		}

		nitricVar, err := tfDeployment.resolveService(intentName, resourceIntent.ServiceIntent, spec, plug)
		if err != nil {
			return err
		}

		serviceInputs[intentName] = nitricVar
	}

	serviceEnvs := map[string][]interface{}{}

	// Prepare non-service/non-entrypoint modules
	for intentName, resourceIntent := range appSpec.ResourceIntents {
		if resourceIntent.Type == "service" || resourceIntent.Type == "entrypoint" {
			continue
		}

		spec, err := e.platform.GetResourceBlueprint(resourceIntent.Type, resourceIntent.SubType)
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.Type, resourceIntent.SubType, err)
		}
		plug, err := e.repository.GetResourcePlugin(spec.PluginId)
		if err != nil {
			return fmt.Errorf("could not find plugin %s", spec.PluginId)
		}

		servicesInput := map[string]any{}
		if access, ok := app_spec_schema.IsAccessible(resourceIntent); ok {
			for serviceName, actions := range access.GetAccess() {
				idMap, ok := tfDeployment.serviceIdentities[serviceName]
				if !ok {
					return fmt.Errorf("service %s not found", serviceName)
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
		tfDeployment.createVariablesForIntent(intentName, resourceIntent, spec)

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
	for intentName, resourceIntent := range appSpec.ResourceIntents {
		if resourceIntent.Type != "service" {
			continue
		}
		spec, err := e.platform.GetResourceBlueprint(resourceIntent.Type, resourceIntent.SubType)
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.Type, resourceIntent.SubType, err)
		}
		plug, err := e.repository.GetResourcePlugin(spec.PluginId)
		if err != nil {
			return fmt.Errorf("could not find plugin %s", spec.PluginId)
		}

		var nitricVar *NitricServiceVariables = serviceInputs[intentName]

		origEnv := nitricVar.Env

		mergedEnv := serviceEnvs[intentName]
		allEnv := append(mergedEnv, origEnv)

		nitricVar.Env = cdktf.Fn_Merge(&allEnv)

		tfDeployment.createVariablesForIntent(intentName, resourceIntent, spec)

		tfDeployment.terraformResources[intentName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})
	}

	// Resolve entrypoint modules
	for intentName, resourceIntent := range appSpec.ResourceIntents {
		if resourceIntent.Type != "entrypoint" {
			continue
		}

		spec, err := e.platform.GetResourceBlueprint(resourceIntent.Type, resourceIntent.SubType)
		if err != nil {
			return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.Type, resourceIntent.SubType, err)
		}
		plug, err := e.repository.GetResourcePlugin(spec.PluginId)
		if err != nil {
			return fmt.Errorf("could not find plugin %s", spec.PluginId)
		}

		nitricVar, err := tfDeployment.resolveEntrypointNitricVar(intentName, appSpec, resourceIntent.EntrypointIntent)
		if err != nil {
			return err
		}

		tfDeployment.createVariablesForIntent(intentName, resourceIntent, spec)

		tfDeployment.terraformResources[intentName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})
	}

	// Resolve resource tokens
	for resourceName, resource := range appSpec.ResourceIntents {
		resourceSpec, err := e.platform.GetResourceBlueprint(resource.Type, resource.SubType)
		if err != nil {
			return err
		}

		err = tfDeployment.resolveTokensForModule(resourceName, resourceSpec, tfDeployment.terraformResources[resourceName])
		if err != nil {
			return err
		}
	}

	// Resolve infra tokens
	for infraName, infra := range e.platform.InfraSpecs {
		// TODO: This is overloading this method as infra-name is not usable in this context as infra cannot resolve `self` tokens
		err := tfDeployment.resolveTokensForModule(infraName, infra, tfDeployment.terraformInfraResources[infraName])
		if err != nil {
			return err
		}
	}

	tfDeployment.app.Synth()

	return nil
}

var _ engines.Engine = &TerraformEngine{}

type terraformEngineOption func(*TerraformEngine)

func WithRepository(repository PluginRepository) terraformEngineOption {
	return func(engine *TerraformEngine) {
		engine.repository = repository
	}
}

func NewFromFile(platformFile io.Reader, opts ...terraformEngineOption) *TerraformEngine {
	platform := &PlatformSpec{}

	json.NewDecoder(platformFile).Decode(platform)

	return New(platform, opts...)
}

func New(platformSpec *PlatformSpec, opts ...terraformEngineOption) *TerraformEngine {
	engine := &TerraformEngine{
		platform: platformSpec,
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine
}
