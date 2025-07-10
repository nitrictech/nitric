package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nitrictech/nitric/engines/terraform"
)

type PlatformRevisionResponse struct {
	Content terraform.PlatformSpec `json:"content"`
}

// FIXME: Because of the difference in fields between identity and resource plugins we need to return an interface
func (c *NitricApiClient) GetPlatform(team, name string, revision int) (*terraform.PlatformSpec, error) {
	response, err := c.get(fmt.Sprintf("/api/platforms/%s/%s/revisions/%d", team, name, revision), true)
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

	var platformSpec PlatformRevisionResponse
	err = json.Unmarshal(body, &platformSpec)
	if err != nil {
		return nil, fmt.Errorf("unexpected response from nitric plugin details endpoint: %v", err)
	}
	return &platformSpec.Content, nil
}
