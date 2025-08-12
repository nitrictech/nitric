package terraform

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/aws/jsii-runtime-go"
	random "github.com/cdktf/cdktf-provider-random-go/random/v11/provider"
	"github.com/cdktf/cdktf-provider-random-go/random/v11/stringresource"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	app_spec_schema "github.com/nitrictech/nitric/cli/pkg/schema"
)

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

func NewTerraformDeployment(engine *TerraformEngine, stackName string) *TerraformDeployment {
	app := cdktf.NewApp(&cdktf.AppConfig{
		Outdir: jsii.String(engine.outputDir),
	})
	stack := cdktf.NewTerraformStack(app, jsii.String(stackName))

	NewNilTerraformBackend(app, jsii.String("nil_backend"))

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

func (td *TerraformDeployment) resolveInfraResource(infraName string) (cdktf.TerraformHclModule, error) {
	resource, ok := td.engine.platform.InfraSpecs[infraName]
	if !ok {
		return nil, fmt.Errorf("infra resource %s not found", infraName)
	}

	if _, ok := td.terraformInfraResources[infraName]; !ok {
		pluginRef, err := td.engine.resolvePlugin(resource)
		if err != nil {
			return nil, err
		}

		td.terraformInfraResources[infraName] = cdktf.NewTerraformHclModule(td.stack, jsii.String(infraName), &cdktf.TerraformHclModuleConfig{
			Source: jsii.String(pluginRef.Deployment.Terraform),
		})
	}

	return td.terraformInfraResources[infraName], nil
}

func (td *TerraformDeployment) resolveEntrypointNitricVar(name string, appSpec *app_spec_schema.Application, spec *app_spec_schema.EntrypointIntent) (interface{}, error) {
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

		hclTarget, ok := td.terraformResources[route.TargetName]
		if !ok {
			return nil, fmt.Errorf("target %s not found", route.TargetName)
		}

		domainNameNitricVar := hclTarget.Get(jsii.String("nitric.domain_name"))
		idNitricVar := hclTarget.Get(jsii.String("nitric.id"))
		resourcesNitricVar := hclTarget.Get(jsii.String("nitric.exports.resources"))

		origins[route.TargetName] = map[string]interface{}{
			"path":        jsii.String(path),
			"base_path":   jsii.String(route.BasePath),
			"type":        jsii.String(intentTargetType),
			"id":          idNitricVar,
			"domain_name": domainNameNitricVar,
			"resources":   resourcesNitricVar,
		}
	}

	nitricVar := map[string]interface{}{
		"name":     jsii.String(name),
		"stack_id": td.stackId.Result(),
		"origins":  origins,
	}

	return nitricVar, nil
}

func (td *TerraformDeployment) resolveService(name string, spec *app_spec_schema.ServiceIntent, resourceSpec *ServiceBlueprint, plug *ResourcePluginManifest) (*NitricServiceVariables, error) {
	var imageVars *map[string]interface{} = nil

	pluginManifest, err := td.engine.resolvePluginsForService(plug)
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

	imageModule := cdktf.NewTerraformHclModule(td.stack, jsii.Sprintf("%s_image", name), &cdktf.TerraformHclModuleConfig{
		Source:    jsii.String(imageModuleRef),
		Variables: imageVars,
	})

	identityModuleOutputs := map[string]interface{}{}
	for _, id := range resourceSpec.Identities {
		identityPlugin, err := td.engine.resolveIdentityPlugin(&id)
		if err != nil {
			return nil, err
		}

		idModule := cdktf.NewTerraformHclModule(td.stack, jsii.Sprintf("%s_%s_role", name, identityPlugin.Name), &cdktf.TerraformHclModuleConfig{
			Source:    jsii.String(identityPlugin.Deployment.Terraform),
			Variables: &id.Properties,
		})

		idModule.Set(jsii.String("nitric"), map[string]interface{}{
			"name":     jsii.String(name),
			"stack_id": td.stackId.Result(),
		})

		identityModuleOutputs[identityPlugin.IdentityType] = idModule.Get(jsii.String("nitric"))
	}

	for _, requiredIdentity := range plug.RequiredIdentities {
		providedIdentities := slices.Collect(maps.Keys(identityModuleOutputs))
		if ok := slices.Contains(providedIdentities, requiredIdentity); !ok {
			return nil, fmt.Errorf("service %s is missing identity %s, required by plugin %s, provided identities were %s", name, requiredIdentity, plug.Name, providedIdentities)
		}
	}

	nitricVar := &NitricServiceVariables{
		NitricVariables: NitricVariables{
			Name: jsii.String(name),
		},
		Schedules:  &schedules,
		ImageId:    imageModule.GetString(jsii.String("image_id")),
		Identities: &identityModuleOutputs,
		StackId:    td.stackId.Result(),
		Env:        &spec.Env,
	}

	td.serviceIdentities[name] = identityModuleOutputs

	return nitricVar, nil
}

func (td *TerraformDeployment) Synth() {
	td.app.Synth()
}
