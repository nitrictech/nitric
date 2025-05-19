package terraform

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	coreschema "github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines"
	"github.com/nitrictech/nitric/engines/terraform/schema"
)

type TerraformEngine struct {
	platform schema.TerraformPlatform
}

func resolvePlugin(pluginName string) (*schema.TerraformPluginManifest, error) {
	return nil, fmt.Errorf("plugin %s not found", pluginName)
}

func resolvePluginName(resource schema.TerraformPlatformResource, subtype string) (string, error) {
	plugin := resource.Plugin
	if subtype != "" {
		if _, ok := resource.Subtypes[subtype]; !ok {
			return "", fmt.Errorf("subtype %s not found", subtype)
		}

		plugin = resource.Subtypes[subtype].Plugin
	}

	// resolve the plugins manifest and locate the deployment module
	return plugin, nil
}

func (e *TerraformEngine) getPlatformResourceForType(resourceType string) (schema.TerraformPlatformResource, error) {
	switch resourceType {
	case "service":
		return e.platform.Services, nil
	case "entrypoint":
		return e.platform.Entrypoints, nil
	}
	return schema.TerraformPlatformResource{}, fmt.Errorf("resource type %s not found", resourceType)
}

// Apply the engine to the target environment
func (e *TerraformEngine) Apply(application *coreschema.Application, environment map[string]interface{}) error {
	app := cdktf.NewApp(&cdktf.AppConfig{})

	stack := cdktf.NewTerraformStack(app, jsii.String(application.Name))
	terraformResources := map[string]cdktf.TerraformHclModule{}
	terraformInfraResources := map[string]cdktf.TerraformHclModule{}

	resolvableResourceProperties := map[string]map[string]interface{}{}
	resolvableInfraProperties := map[string]map[string]interface{}{}

	// 1. Start deploying the platform
	for resourceName, resource := range application.Resources {
		terraformPlatformResource, err := e.getPlatformResourceForType(resource.Type)
		if err != nil {
			return err
		}

		pluginName, err := resolvePluginName(terraformPlatformResource, resource.SubType)
		if err != nil {
			return err
		}

		plugin, err := resolvePlugin(pluginName)
		if err != nil {
			return err
		}

		terraformResources[resourceName] = cdktf.NewTerraformHclModule(stack, jsii.String(resourceName), &cdktf.TerraformHclModuleConfig{
			// This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plugin.Deployment.Terraform),
		})
	}

	// 2. Deploy platform infra resources
	for infraName, infra := range e.platform.Infra {
		// Locate the plugin for the infra from platform
		plugin, err := resolvePlugin(infra.Plugin)
		if err != nil {
			return err
		}

		terraformInfraResources[infraName] = cdktf.NewTerraformHclModule(stack, jsii.String(infraName), &cdktf.TerraformHclModuleConfig{
			// This assumes that the plugin is resolvable as a URI
			Source: jsii.String(plugin.Deployment.Terraform),
		})
	}

	// 3. Map inputs and outputs from the platform and environment to the stack resources

	app.Synth()

	return nil
}

var _ engines.Engine = &TerraformEngine{}

func New(platformFile io.Reader) *TerraformEngine {
	platform := &schema.TerraformPlatform{}

	json.NewDecoder(platformFile).Decode(platform)

	return &TerraformEngine{
		platform: platform,
	}
}
