package terraform

import (
	"embed"
	"fmt"
	"io/fs"

	"gopkg.in/yaml.v3"
)

//go:embed plugins/**/manifest.yaml
var manifestFs embed.FS

var resourceManifests = map[string]*ResourcePluginManifest{}
var identityManifests = map[string]*IdentityPluginManifest{}
var allManifests = map[string]*PluginManifest{}

// Read all manifests and build a map of name -> manifest
func init() {
	// Walk filesystem and read each manifest.yaml file
	fs.WalkDir(manifestFs, ".", func(path string, d fs.DirEntry, err error) error {
		// fmt.Println(path)

		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		manifestBytes, err := fs.ReadFile(manifestFs, path)
		if err != nil {
			return err
		}

		var manifest PluginManifest

		err = yaml.Unmarshal(manifestBytes, &manifest)
		if err != nil {
			return err
		}

		if _, ok := allManifests[manifest.Name]; ok {
			return fmt.Errorf("duplicate plugin detected %s", manifest.Name)
		}

		allManifests[manifest.Name] = &manifest

		if manifest.Type == "identity" {
			var identityManifest IdentityPluginManifest

			err = yaml.Unmarshal(manifestBytes, &identityManifest)
			if err != nil {
				return err
			}

			fmt.Println(identityManifest.Name)

			identityManifests[identityManifest.Name] = &identityManifest
		} else {
			var resourceManifest ResourcePluginManifest

			err := yaml.Unmarshal(manifestBytes, &resourceManifest)
			if err != nil {
				return err
			}

			resourceManifests[resourceManifest.Name] = &resourceManifest
		}

		return nil
	})
}

type NitricTerraformPluginRepository struct {
}

func pluginType(name string) (string, error) {
	man, ok := allManifests[name]
	if !ok {
		return "", fmt.Errorf("plugin %s not found", name)
	}

	return man.Type, nil
}

func (r *NitricTerraformPluginRepository) GetResourcePlugin(name string) (*ResourcePluginManifest, error) {
	pluginType, err := pluginType(name)
	if err != nil {
		return nil, fmt.Errorf("resource plugin %s not found", name)
	}

	if pluginType == "identity" {
		return nil, fmt.Errorf("plugin %s is of type %s and cannot be used as a resource", name, pluginType)
	}

	manifest, ok := resourceManifests[name]
	if !ok {
		return nil, fmt.Errorf("resource plugin %s not found", name)
	}

	return manifest, nil
}

func (r *NitricTerraformPluginRepository) GetIdentityPlugin(name string) (*IdentityPluginManifest, error) {
	manifest, ok := identityManifests[name]
	if !ok {
		if _, ok := resourceManifests[name]; ok {
			return nil, fmt.Errorf("%s is a resource plugin and cannot be used as an identity plugin", name)
		}

		return nil, fmt.Errorf("identity plugin %s not found", name)
	}

	return manifest, nil
}

func NewNitricTerraformPluginRepository() *NitricTerraformPluginRepository {
	return &NitricTerraformPluginRepository{}
}
