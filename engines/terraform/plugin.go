package terraform

type PluginManifest struct {
	Name       string           `json:"name" yaml:"name"`
	Version    string           `json:"version" yaml:"version"`
	Deployment DeploymentModule `json:"deployment" yaml:"deployment"`
	Type       string           `json:"type" yaml:"type"`
	Runtime    RuntimeModule    `json:"runtime" yaml:"runtime"`
	Inputs     []PluginInput    `json:"inputs" yaml:"inputs"`
	Outputs    []PluginOutput   `json:"outputs" yaml:"outputs"`
}

type ResourcePluginManifest struct {
	PluginManifest     `json:",inline" yaml:",inline"`
	RequiredIdentities []string `json:"required_identities" yaml:"required_identities"`
}

type IdentityPluginManifest struct {
	PluginManifest `json:",inline" yaml:",inline"`
	IdentityType   string `json:"identity_type" yaml:"identity_type"`
}

type DeploymentModule struct {
	Terraform string `json:"terraform" yaml:"terraform"`
}

type RuntimeModule struct {
	GoModule string `json:"go_module" yaml:"go_module"`
}

type PluginInput struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Type        string `json:"default" yaml:"default"`
	Required    bool   `json:"required" yaml:"required"`
}

type PluginOutput struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
}

type PluginRepository interface {
	GetResourcePlugin(name string) (*ResourcePluginManifest, error)
	GetIdentityPlugin(name string) (*IdentityPluginManifest, error)
}
