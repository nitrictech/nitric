package schema

type SubscriptionResource struct {
	SubscriptionSchemaOnlyHackType string `json:"type" jsonschema:"type,enum=subscription"`

	Source string `json:"source"`
	Target string `json:"target"`
	Path   string `json:"path"`
}
