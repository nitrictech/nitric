package terraform

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/aws/jsii-runtime-go"
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
	app    cdktf.App
	stack  cdktf.TerraformStack
	engine *TerraformEngine

	terraformResources      map[string]cdktf.TerraformHclModule
	terraformInfraResources map[string]cdktf.TerraformHclModule
	terraformVariables      map[string]cdktf.TerraformVariable
}

type SpecReference struct {
	// var/infra/etc
	Source string
	// simple key for var or path for infra e.g. vpc.arn
	Path []string
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

func (tf *TerraformDeployment) resolveTokensForModule(resource *ResourceBlueprint, module cdktf.TerraformHclModule) error {
	for property, value := range resource.Properties {
		module.Set(jsii.String(property), value)

		token, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid token format")
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
		} else if specRef.Source == "var" {
			// TODO: Remove dynamic variable creation, instead reference from spec (add a variables definition to the platform spec)
			tfVariable, ok := tf.terraformVariables[specRef.Path[0]]
			if !ok {
				tf.terraformVariables[specRef.Path[0]] = cdktf.NewTerraformVariable(tf.stack, jsii.String(specRef.Path[0]), &cdktf.TerraformVariableConfig{})
				tfVariable = tf.terraformVariables[specRef.Path[0]]
			}

			// Create a new terraform variable
			module.Set(jsii.String(property), tfVariable.Value())
		}
	}

	return nil
}

func NewTerraformDeployment(engine *TerraformEngine, stackName string) *TerraformDeployment {
	app := cdktf.NewApp(&cdktf.AppConfig{})

	return &TerraformDeployment{
		app:                     app,
		stack:                   cdktf.NewTerraformStack(app, jsii.String(stackName)),
		engine:                  engine,
		terraformResources:      map[string]cdktf.TerraformHclModule{},
		terraformInfraResources: map[string]cdktf.TerraformHclModule{},
		terraformVariables:      map[string]cdktf.TerraformVariable{},
	}
}

func (e *TerraformEngine) resolvePluginsForService(servicePlugin *PluginManifest) (*plugin.PluginDefintion, error) {
	// TODO: Map platform resource plugins to the service plugin
	return &plugin.PluginDefintion{
		Service: plugin.GoPlugin{
			Alias:  "svcPlugin",
			Name:   "default",
			Import: servicePlugin.Runtime.GoModule,
		},
	}, nil
}

func (e *TerraformDeployment) resolveService(name string, spec *app_spec_schema.ServiceIntent, resourceSpec *ServiceBlueprint, plug *PluginManifest) (interface{}, error) {
	// Map the nitric variable
	var nitricVar interface{} = nil
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
			"args":     map[string]interface{}{"PLUGIN_DEFINITION": jsii.String(string(pluginManifestBytes))},
		}
	} else if spec.Container.Docker != nil {
		imageVars = &map[string]interface{}{
			"build_context": jsii.String(spec.Container.Docker.Context),
			"dockerfile":    jsii.String(spec.Container.Docker.Dockerfile),
			"tag":           jsii.String(name),
			"args":          map[string]interface{}{"PLUGIN_DEFINITION": jsii.String(string(pluginManifestBytes))},
		}
	}

	imageModule := cdktf.NewTerraformHclModule(e.stack, jsii.Sprintf("%s_image", name), &cdktf.TerraformHclModuleConfig{
		Source:    jsii.String("github.com/nitrictech/nitric//engines/terraform/modules/image?depth=1&ref=next"),
		Variables: imageVars,
	})

	identityModules := map[string]cdktf.TerraformHclModule{}
	identityModuleOutputs := map[string]interface{}{}
	for identityName, id := range resourceSpec.Identities {
		identityPlugin, err := e.engine.repository.GetPlugin(id.PluginId)
		if err != nil {
			return nil, err
		}

		identityModules[identityName] = cdktf.NewTerraformHclModule(e.stack, jsii.Sprintf("%s_%s_role", name, identityName), &cdktf.TerraformHclModuleConfig{
			Source:    jsii.String(identityPlugin.Deployment.Terraform),
			Variables: &id.Properties,
		})

		identityModuleOutputs[identityName] = identityModules[identityName].Get(jsii.String("nitric"))
	}

	// Create this services identities

	nitricVar = &NitricServiceVariables{
		NitricVariables: NitricVariables{
			Name: jsii.String(name),
		},
		ImageId:    imageModule.GetString(jsii.String("image_id")),
		Identities: &identityModuleOutputs,
		Env:        &spec.Env,
	}

	return nitricVar, nil

}

// Apply the engine to the target environment
func (e *TerraformEngine) Apply(appSpec *app_spec_schema.Application) error {
	tfDeployment := NewTerraformDeployment(e, appSpec.Name)

	// Create a terraform variable to establish the root context for application builds
	// this will be prepended to the path of any internal docker builds
	// tfDeployment.terraformVariables["build_root"] = cdktf.NewTerraformVariable(tfDeployment.stack, jsii.String("build_root"), &cdktf.TerraformVariableConfig{
	// 	Type: jsii.String("string"),
	// })

	// Resolve resource modules
	for intentName, resourceIntent := range appSpec.ResourceIntents {
		var nitricVar interface{} = nil
		var plug *PluginManifest = nil
		if resourceIntent.Type == "service" {
			spec, err := e.platform.GetServiceBlueprint(resourceIntent.SubType)
			if err != nil {
				return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.Type, resourceIntent.SubType, err)
			}
			plug, err = e.repository.GetPlugin(spec.PluginId)
			if err != nil {
				return fmt.Errorf("could not find plugin %s", spec.PluginId)
			}

			nitricVar, err = tfDeployment.resolveService(intentName, resourceIntent.ServiceIntent, spec, plug)
			if err != nil {
				return err
			}
		} else {
			spec, err := e.platform.GetResourceBlueprint(resourceIntent.Type, resourceIntent.SubType)
			if err != nil {
				return fmt.Errorf("could not find platform type for %s.%s: %w", resourceIntent.Type, resourceIntent.SubType, err)
			}
			plug, err = e.repository.GetPlugin(spec.PluginId)
			if err != nil {
				return fmt.Errorf("could not find plugin %s", spec.PluginId)
			}

		}

		tfDeployment.terraformResources[intentName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(intentName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plug.Deployment.Terraform),
			Variables: &map[string]interface{}{
				"nitric": nitricVar,
			},
		})
	}

	// Resolve infra modules
	for infraName, infra := range e.platform.InfraSpecs {
		plugin, err := e.repository.GetPlugin(infra.PluginId)
		if err != nil {
			return err
		}

		tfDeployment.terraformInfraResources[infraName] = cdktf.NewTerraformHclModule(tfDeployment.stack, jsii.String(infraName), &cdktf.TerraformHclModuleConfig{
			// TODO: This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plugin.Deployment.Terraform),
		})
	}

	// Resolve resource tokens
	for resourceName, resource := range appSpec.ResourceIntents {
		resourceSpec, err := e.platform.GetResourceBlueprint(resource.Type, resource.SubType)
		if err != nil {
			return err
		}

		tfDeployment.resolveTokensForModule(resourceSpec, tfDeployment.terraformResources[resourceName])
	}

	// Resolve infra tokens
	for infraName, infra := range e.platform.InfraSpecs {
		tfDeployment.resolveTokensForModule(infra, tfDeployment.terraformInfraResources[infraName])
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
