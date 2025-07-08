package terraform

type PluginManifest struct {
	Name       string                  `json:"name" yaml:"name"`
	Icon       string                  `json:"icon" yaml:"icon"`
	Deployment DeploymentModule        `json:"deployment" yaml:"deployment"`
	Type       string                  `json:"type" yaml:"type"`
	Runtime    RuntimeModule           `json:"runtime" yaml:"runtime"`
	Inputs     map[string]PluginInput  `json:"inputs" yaml:"inputs"`
	Outputs    map[string]PluginOutput `json:"outputs" yaml:"outputs"`
}

type ResourcePluginManifest struct {
	PluginManifest     `json:",inline" yaml:",inline"`
	RequiredIdentities []string `json:"required_identities" yaml:"required_identities"`
	Capabilities       []string `json:"capabilities" yaml:"capabilities"`
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
	Description string `json:"description" yaml:"description"`
	Type        string `json:"default" yaml:"default"`
	Required    bool   `json:"required" yaml:"required"`
}

type PluginOutput struct {
	Description string `json:"description" yaml:"description"`
}

type PluginRepository interface {
	GetResourcePlugin(team, libname, version, name string) (*ResourcePluginManifest, error)
	GetIdentityPlugin(team, libname, version, name string) (*IdentityPluginManifest, error)
}
