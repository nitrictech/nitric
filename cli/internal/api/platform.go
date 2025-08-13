package api

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/nitrictech/nitric/engines/terraform"
)

// FIXME: Because of the difference in fields between identity and resource plugins we need to return an interface
func (c *NitricApiClient) GetPlatform(team, name string, revision int) (*terraform.PlatformSpec, error) {
	response, err := c.get(fmt.Sprintf("/api/teams/%s/platforms/%s/revisions/%d", team, name, revision), true)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		if response.StatusCode == 404 {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("received non 200 response from %s plugin details endpoint: %d", version.ProductName, response.StatusCode)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from %s plugin details endpoint: %v", version.ProductName, err)
	}

	var platformRevision GetPlatformRevisionResponse
	err = json.Unmarshal(body, &platformRevision)
	if err != nil {
		return nil, fmt.Errorf("unexpected response from %s plugin details endpoint: %v", version.ProductName, err)
	}
	return &platformRevision.Revision.Content, nil
}
