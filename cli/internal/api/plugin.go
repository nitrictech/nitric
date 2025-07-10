package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nitrictech/nitric/engines/terraform"
)

// FIXME: Because of the difference in fields between identity and resource plugins we need to return an interface
func (c *NitricApiClient) GetPluginManifest(team, lib, version, name string) (interface{}, error) {
	response, err := c.get(fmt.Sprintf("/api/plugin_libraries/%s/%s/versions/%s/plugins/%s", team, lib, version, name), true)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nitric auth details endpoint: %v", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("received non 200 response from nitric plugin details endpoint: %d", response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from nitric auth details endpoint: %v", err)
	}

	var pluginManifest terraform.ResourcePluginManifest
	err = json.Unmarshal(body, &pluginManifest)
	if err != nil {
		return nil, fmt.Errorf("unexpected response from nitric plugin details endpoint: %v", err)
	}

	if pluginManifest.Type == "identity" {
		var identityPluginManifest terraform.IdentityPluginManifest
		err = json.Unmarshal(body, &identityPluginManifest)
		if err != nil {
			return nil, fmt.Errorf("unexpected response from nitric plugin details endpoint: %v", err)
		}

		return &identityPluginManifest, nil
	}

	return &pluginManifest, nil
}
