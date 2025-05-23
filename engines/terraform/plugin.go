package terraform

type PluginManifest struct {
	Name       string           `json:"name" yaml:"name"`
	Version    string           `json:"version" yaml:"version"`
	Deployment DeploymentModule `json:"deployment" yaml:"deployment"`
	Runtime    RuntimeModule    `json:"runtime" yaml:"runtime"`
	Inputs     []PluginInput    `json:"inputs" yaml:"inputs"`
	Outputs    []PluginOutput   `json:"outputs" yaml:"outputs"`
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
	GetPlugin(name string) (*PluginManifest, error)
}
