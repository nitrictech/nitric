package terraform

import (
	"embed"
	"fmt"
	"io/fs"

	"gopkg.in/yaml.v3"
)

//go:embed plugins/**/manifest.yaml
var manifestFs embed.FS

var manifests = map[string]*PluginManifest{}

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

		manifests[manifest.Name] = &manifest

		return nil
	})
}

type NitricTerraformPluginRepository struct {
}

func (r *NitricTerraformPluginRepository) GetPlugin(name string) (*PluginManifest, error) {
	manifest, ok := manifests[name]
	if !ok {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	return manifest, nil
}

func NewNitricTerraformPluginRepository() *NitricTerraformPluginRepository {
	return &NitricTerraformPluginRepository{}
}
