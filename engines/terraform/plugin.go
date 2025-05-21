package terraform

type PluginManifest struct {
	Name       string           `json:"name"`
	Version    string           `json:"version"`
	Deployment DeploymentModule `json:"deployment"`
	Runtime    RuntimeModule    `json:"runtime"`
	Inputs     []PluginInput    `json:"inputs"`
	Outputs    []PluginOutput   `json:"outputs"`
}

type DeploymentModule struct {
	Terraform string `json:"terraform"`
}

type RuntimeModule struct {
	GoModule string `json:"go_module"`
}

type PluginInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"default"`
	Required    bool   `json:"required"`
}

type PluginOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PluginRepository interface {
	GetPlugin(name string) (*PluginManifest, error)
}
