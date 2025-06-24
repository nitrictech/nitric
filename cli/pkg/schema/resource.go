package schema

type Resource struct {
	SubType string `json:"sub-type,omitempty" yaml:"sub-type,omitempty" jsonschema:"-"`
}
