package api

import (
	"fmt"
	"time"

	"github.com/nitrictech/nitric/engines/terraform"
)

// TODO: We need a better mechanism to sync the models with the server.

// Template DTOs from backend
type TemplateResponse struct {
	Slug     string   `json:"slug"`
	TeamSlug string   `json:"teamSlug"`
	Versions []string `json:"versions"`
}

type ListTemplatesResponse struct {
	Templates []TemplateResponse `json:"templates"`
}

// Legacy aliases for backward compatibility
type Template = TemplateResponse
type GetTemplatesResponse = ListTemplatesResponse

func (t *TemplateResponse) String() string {
	return fmt.Sprintf("%s/%s", t.TeamSlug, t.Slug)
}

type TemplateVersion struct {
	TemplateSlug string `json:"templateSlug"`
	TeamSlug     string `json:"teamSlug"`
	Version      string `json:"version"`
	GitSource    string `json:"gitSource"`
	Public       bool   `json:"public"`
}

type GetTemplateVersionResponse struct {
	Template *TemplateVersion `json:"template"`
}

func (t *TemplateVersion) String() string {
	return fmt.Sprintf("%s/%s@%s", t.TeamSlug, t.TemplateSlug, t.Version)
}

// Platform DTOs from backend
type PlatformResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	TeamSlug    string    `json:"teamSlug"`
	Revisions   []int32   `json:"revisions"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PlatformsResponse struct {
	Platforms []PlatformResponse `json:"platforms"`
}

type CreatePlatformRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type CreatePlatformRevisionRequest = terraform.PlatformSpec

type Platform struct {
	ID          int64     `json:"id"`
	TeamSlug    string    `json:"teamSlug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Revisions   []int32   `json:"revisions,omitempty"`
}

type PlatformRevision struct {
	ID         int64                  `json:"id"`
	PlatformID int64                  `json:"platformId"`
	Revision   int32                  `json:"revision"`
	Content    terraform.PlatformSpec `json:"content"`
	CreatedAt  time.Time              `json:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt"`
}

type CreatePlatformResponse struct {
	Platform *Platform `json:"platform"`
}

type CreatePlatformRevisionResponse struct {
	Revision *PlatformRevision `json:"revision"`
}

type GetPlatformResponse struct {
	Platform *Platform `json:"platform"`
}

type GetPlatformRevisionResponse struct {
	Revision *PlatformRevision `json:"revision"`
}

// Plugin DTOs from backend
type PluginResponse struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Public bool   `json:"public"`
	TeamID int64  `json:"teamId"`
}

type PluginVersionResponse struct {
	ID       int64  `json:"id"`
	PluginID int64  `json:"pluginId"`
	Version  string `json:"version"`
	Manifest string `json:"manifest"`
}

type PluginVersionsResponse struct {
	PluginVersions []PluginVersionResponse `json:"pluginVersions"`
}

type PluginsResponse struct {
	Plugins []PluginResponse `json:"plugins"`
}

type PluginProviderResponse struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Public   bool   `json:"public"`
	TeamSlug string `json:"teamSlug"`
}

type CreatePluginLibraryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`
	Public      bool   `json:"public"`
}

type CreatePluginLibraryVersionRequest struct {
	PluginProviderID int64  `json:"pluginProviderId"`
	Version          string `json:"version"`
}

type PluginProviderVersionResponse struct {
	ID               int64  `json:"id"`
	PluginProviderID int64  `json:"pluginProviderId"`
	Version          string `json:"version"`
}

type PluginSummary struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Type string `json:"type"`
}

type PluginLibrary struct {
	ID          int64     `json:"id"`
	TeamSlug    string    `json:"teamSlug"`
	Name        string    `json:"name"`
	Source      string    `json:"source"`
	Public      bool      `json:"public"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type PluginLibraryWithVersions struct {
	PluginLibrary
	Versions []string `json:"versions"`
}

type PluginLibraryVersion struct {
	PluginLibrary
	Version string          `json:"version"`
	Plugins []PluginSummary `json:"plugins"`
}

type ListPluginLibrariesResponse struct {
	Libraries []PluginLibraryWithVersions `json:"libraries"`
}

type GetPluginLibraryResponse struct {
	Library *PluginLibraryWithVersions `json:"library"`
}

type GetPluginLibraryVersionResponse struct {
	Version *PluginLibraryVersion `json:"version"`
}

type GetPluginManifestResponse struct {
	Manifest map[string]interface{} `json:"manifest"`
}

// Team DTOs from backend
type CreateTeamRequest struct {
	Name string `json:"name"`
}

type CreateTeamResponse struct {
	Team *Team `json:"team"`
}
