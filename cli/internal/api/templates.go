package api

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *NitricApiClient) GetTemplates() (*ListTemplatesResponse, error) {
	response, err := c.get("/api/templates")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get templates: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var templates ListTemplatesResponse
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templates: %v, body: %s", err, string(body))
	}

	return &templates, nil
}
