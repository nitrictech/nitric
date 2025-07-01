package api

import "time"

// TODO: We need a better mechanism to sync the models with the server.

type TemplateResponse struct {
	Name             string   `json:"name"`
	OrganizationSlug string   `json:"organizationSlug"`
	Versions         []string `json:"versions"`
}

type TemplateVersionResponse struct {
	ID        int64     `json:"id"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ListTemplateVersionsResponse struct {
	Versions []*TemplateVersionResponse `json:"versions"`
}
