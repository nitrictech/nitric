package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/nitrictech/nitric/engines/terraform"
)

// FIXME: Because of the difference in fields between identity and resource plugins we need to return an interface
func (c *NitricApiClient) GetPluginManifest(team, lib, libVersion, name string) (interface{}, error) {
	response, err := c.get(fmt.Sprintf("/api/plugin_libraries/%s/%s/versions/%s/plugins/%s", team, lib, libVersion, name), true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s plugin details endpoint: %v", version.ProductName, err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("received non 200 response from %s plugin details endpoint: %d", version.ProductName, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s plugin details endpoint: %v", version.ProductName, err)
	}

	var pluginManifest terraform.ResourcePluginManifest
	err = json.Unmarshal(body, &pluginManifest)
	if err != nil {
		return nil, fmt.Errorf("unexpected response from %s plugin details endpoint: %v", version.ProductName, err)
	}

	if pluginManifest.Type == "identity" {
		var identityPluginManifest terraform.IdentityPluginManifest
		err = json.Unmarshal(body, &identityPluginManifest)
		if err != nil {
			return nil, fmt.Errorf("unexpected response from %s plugin details endpoint: %v", version.ProductName, err)
		}

		return &identityPluginManifest, nil
	}

	return &pluginManifest, nil
}
