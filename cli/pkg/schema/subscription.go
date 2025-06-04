package schema

type SubscriptionIntent struct {
	SubscriptionSchemaOnlyHackType string `json:"type" jsonschema:"type,enum=subscription"`

	Source string `json:"source" yaml:"source"`
	Target string `json:"target" yaml:"target"`
	Path   string `json:"path" yaml:"path"`
}
