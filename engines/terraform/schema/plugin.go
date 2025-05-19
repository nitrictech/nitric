package schema

type TerraformPluginManifest struct {
	Name       string                    `json:"name"`
	Version    string                    `json:"version"`
	Deployment TerraformDeploymentModule `json:"deployment"`
	Runtime    RuntimeModule             `json:"runtime"`
	Inputs     []TerraformPluginInput    `json:"inputs"`
	Outputs    []TerraformPluginOutput   `json:"outputs"`
}

type TerraformDeploymentModule struct {
	Terraform string `json:"terraform"`
}

type RuntimeModule struct {
	GoModule string `json:"go_module"`
}

type TerraformPluginInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"default"`
	Required    bool   `json:"required"`
}

type TerraformPluginOutput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
