package terraform

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"

	app_spec_schema "github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines"
)

type TerraformEngine struct {
	platform   *PlatformSpec
	repository PluginRepository

	outputDir string
}

type SpecReference struct {
	Source string
	Path   []string
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

// Apply the engine to the target environment
func (e *TerraformEngine) Apply(appSpec *app_spec_schema.Application) (string, error) {
	tfDeployment := NewTerraformDeployment(e, appSpec.Name)

	// Create platform variables
	tfDeployment.createPlatformVariables()

	// Process service identities first (needed by other resources)
	serviceInputs, err := tfDeployment.processServiceIdentities(appSpec)
	if err != nil {
		return "", err
	}

	// Process each resource type and collect environment variables they export to services
	allServiceEnvs := map[string][]interface{}{}

	// Process bucket resources
	if len(appSpec.BucketIntents) > 0 {
		bucketEnvs, err := tfDeployment.processBucketResources(appSpec)
		if err != nil {
			return "", err
		}
		// Merge bucket environment variable exports
		for serviceName, envs := range bucketEnvs {
			allServiceEnvs[serviceName] = append(allServiceEnvs[serviceName], envs...)
		}
	}

	// Process database resources
	if len(appSpec.DatabaseIntents) > 0 {
		databaseEnvs, err := tfDeployment.processDatabaseResources(appSpec)
		if err != nil {
			return "", err
		}
		// Merge database environment variable exports
		for serviceName, envs := range databaseEnvs {
			allServiceEnvs[serviceName] = append(allServiceEnvs[serviceName], envs...)
		}
	}

	// Process service resources (needs environment variables from other resources)
	if len(appSpec.ServiceIntents) > 0 {
		err = tfDeployment.processServiceResources(appSpec, serviceInputs, allServiceEnvs)
		if err != nil {
			return "", err
		}
	}

	// Process entrypoint resources
	if len(appSpec.EntrypointIntents) > 0 {
		err = tfDeployment.processEntrypointResources(appSpec)
		if err != nil {
			return "", err
		}
	}

	// Resolve resource tokens for all created resources
	resourceIntents := appSpec.GetResourceIntents()
	for resourceName, resourceIntent := range resourceIntents {
		resourceSpec, err := e.platform.GetResourceBlueprint(resourceIntent.GetType(), resourceIntent.GetSubType())
		if err != nil {
			return "", err
		}

		err = tfDeployment.resolveTokensForModule(resourceName, resourceSpec, tfDeployment.terraformResources[resourceName])
		if err != nil {
			return "", err
		}
	}

	// Resolve infra tokens
	for infraName, infra := range tfDeployment.terraformInfraResources {
		infraSpec, ok := e.platform.InfraSpecs[infraName]
		if !ok {
			return "", fmt.Errorf("infra resource %s not found", infraName)
		}
		err := tfDeployment.resolveTokensForModule(infraName, infraSpec, infra)
		if err != nil {
			return "", err
		}
	}

	// Resolve dependencies for all created modules
	for resourceName, resource := range resourceIntents {
		resourceSpec, err := e.platform.GetResourceBlueprint(resource.GetType(), resource.GetSubType())
		if err != nil {
			return "", err
		}

		err = tfDeployment.resolveDependencies(resourceSpec, tfDeployment.terraformResources[resourceName])
		if err != nil {
			return "", err
		}
	}

	// Resolve dependencies for all created infrastructure
	for infraName, infra := range tfDeployment.terraformInfraResources {
		infraSpec, ok := e.platform.InfraSpecs[infraName]
		if !ok {
			return "", fmt.Errorf("infra resource %s not found", infraName)
		}

		err := tfDeployment.resolveDependencies(infraSpec, infra)
		if err != nil {
			return "", err
		}
	}

	tfDeployment.Synth()

	return filepath.Join(e.outputDir, "stacks", appSpec.Name), nil
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

func NewFromFile(platformFile io.Reader, opts ...terraformEngineOption) *TerraformEngine {
	platform := &PlatformSpec{}

	json.NewDecoder(platformFile).Decode(platform)

	return New(platform, opts...)
}

func New(platformSpec *PlatformSpec, opts ...terraformEngineOption) *TerraformEngine {
	engine := &TerraformEngine{
		platform:  platformSpec,
		outputDir: "terraform",
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine
}
