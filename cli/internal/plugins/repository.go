package plugins

import (
	"errors"
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/engines/terraform"
)

type PluginRepository struct {
	apiClient *api.NitricApiClient
}

func (r *PluginRepository) GetResourcePlugin(team, libname, version, name string) (*terraform.ResourcePluginManifest, error) {
	pluginManifest, err := r.apiClient.GetPluginManifest(team, libname, version, name)
	if errors.Is(err, api.ErrNotFound) {
		return nil, fmt.Errorf("plugin %s/%s/%s@%s not found", team, libname, name, version)
	} else if err != nil {
		return nil, err
	}

	resourcePluginManifest, ok := pluginManifest.(*terraform.ResourcePluginManifest)
	if !ok {
		return nil, fmt.Errorf("encountered malformed manifest for plugin %s/%s/%s@%s: %v", team, libname, name, version, err)
	}

	return resourcePluginManifest, nil
}

func (r *PluginRepository) GetIdentityPlugin(team, libname, version, name string) (*terraform.IdentityPluginManifest, error) {
	pluginManifest, err := r.apiClient.GetPluginManifest(team, libname, version, name)
	if errors.Is(err, api.ErrNotFound) {
		return nil, fmt.Errorf("plugin %s/%s/%s@%s not found", team, libname, name, version)
	} else if err != nil {
		return nil, err
	}

	identityPluginManifest, ok := pluginManifest.(*terraform.IdentityPluginManifest)
	if !ok {
		return nil, fmt.Errorf("encountered malformed manifest for plugin %s/%s/%s@%s: %v", team, libname, name, version, err)
	}

	return identityPluginManifest, nil
}

func NewPluginRepository(apiClient *api.NitricApiClient) *PluginRepository {
	return &PluginRepository{
		apiClient: apiClient,
	}
}
