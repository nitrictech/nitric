package plugins

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/internal/api"
	"github.com/nitrictech/nitric/engines/terraform"
)

type PluginRepository struct {
	apiClient *api.NitricApiClient
}

func (r *PluginRepository) GetResourcePlugin(team, libname, version, name string) (*terraform.ResourcePluginManifest, error) {
	pluginManifest, err := r.apiClient.GetPluginManifest(team, libname, version, name)
	if err != nil {
		return nil, err
	}

	resourcePluginManifest, ok := pluginManifest.(*terraform.ResourcePluginManifest)
	if !ok {
		return nil, fmt.Errorf("unexpected resource plugin manifest: %v", err)
	}

	return resourcePluginManifest, nil
}

func (r *PluginRepository) GetIdentityPlugin(team, libname, version, name string) (*terraform.IdentityPluginManifest, error) {
	pluginManifest, err := r.apiClient.GetPluginManifest(team, libname, version, name)
	if err != nil {
		return nil, err
	}

	identityPluginManifest, ok := pluginManifest.(*terraform.IdentityPluginManifest)
	if !ok {
		return nil, fmt.Errorf("unexpected identity plugin manifest: %v", err)
	}

	return identityPluginManifest, nil
}

func NewPluginRepository(apiClient *api.NitricApiClient) *PluginRepository {
	return &PluginRepository{
		apiClient: apiClient,
	}
}
