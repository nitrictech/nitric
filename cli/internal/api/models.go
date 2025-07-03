package api

import (
	"fmt"
	"time"
)

// TODO: We need a better mechanism to sync the models with the server.

type Template struct {
	Slug     string   `json:"slug"`
	TeamSlug string   `json:"teamSlug"`
	Versions []string `json:"versions"`
}

func (t *Template) String() string {
	return fmt.Sprintf("%s/%s", t.TeamSlug, t.Slug)
}

type TemplateVersion struct {
	TemplateSlug      string    `json:"templateSlug"`
	TeamSlug          string    `json:"teamSlug"`
	Version           string    `json:"version"`
	TemplateLibraryId string    `json:"templateLibraryId"`
	GitSource         string    `json:"gitSource"`
	Public            bool      `json:"public"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

func (t *TemplateVersion) String() string {
	return fmt.Sprintf("%s/%s@%s", t.TeamSlug, t.TemplateSlug, t.Version)
}
