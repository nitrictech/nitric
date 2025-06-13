package schema

type SubscriptionIntent struct {
	SubscriptionSchemaOnlyHackType    string `json:"type" jsonschema:"type,enum=subscription"`
	SubscriptionSchemaOnlyHackSubType string `json:"sub-type,omitempty" yaml:"-,omitempty" jsonschema:"sub-type"`

	Source string `json:"source" yaml:"source"`
	Target string `json:"target" yaml:"target"`
	Path   string `json:"path" yaml:"path"`
}
