package api

import (
	"encoding/json"
	"fmt"
	"io"
)

func (c *NitricApiClient) GetTemplates() ([]Template, error) {
	response, err := c.get("/api/templates", true)
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

	var templates []Template
	if err := json.Unmarshal(body, &templates); err != nil {
		return nil, fmt.Errorf("failed to unmarshal templates: %v, body: %s", err, string(body))
	}

	return templates, nil
}

// GetTemplate gets a specific template by teamSlug, templateName and version
// version is optional, if it is not provided, the latest version will be returned
func (c *NitricApiClient) GetTemplate(teamSlug string, templateName string, version string) (*TemplateVersion, error) {
	// latest version URL is /api/templates/{teamSlug}/{templateName}
	// specific version URL is /api/templates/{teamSlug}/{templateName}/v/{version}

	templatePath := fmt.Sprintf("/api/templates/%s/%s", teamSlug, templateName)

	if version != "" {
		templatePath = fmt.Sprintf("%s/v/%s", templatePath, version)
	}

	response, err := c.get(templatePath, true)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get template: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var template *TemplateVersion
	if err := json.Unmarshal(body, &template); err != nil {
		return nil, fmt.Errorf("failed to unmarshal template: %v, body: %s", err, string(body))
	}

	return template, nil
}
