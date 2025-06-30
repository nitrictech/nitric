package api

import "time"

// TODO: We need a better mechanism to sync the models with the server.

type TemplateResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Source      string    `json:"source"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ListTemplatesResponse struct {
	Templates []*TemplateResponse `json:"templates"`
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
