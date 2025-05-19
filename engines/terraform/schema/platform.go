package schema

type TerraformPlatform struct {
	Name        string                           `json:"name"`
	Services    TerraformPlatformResource        `json:"services"`
	Entrypoints TerraformPlatformResource        `json:"entrypoints"`
	Infra       map[string]BaseTerraformResource `json:"infra"`
}

type BaseTerraformResource struct {
	Plugin     string            `json:"plugin"`
	Properties map[string]string `json:"properties"` // XXX: May need to be map[string]interface{}
}

type TerraformPlatformResource struct {
	BaseTerraformResource `json:",inline"`
	Subtypes              map[string]BaseTerraformResource `json:"subtypes"`
}
