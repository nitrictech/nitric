package schema

type TopicIntent struct {
	TopicSchemaOnlyHackType       string              `json:"type" yaml:"-" jsonschema:"type,enum=topic"`
	TopicAccessSchemaOnlyHackType map[string][]string `json:"access,omitempty" yaml:"access,omitempty"`
	TopicSchemaOnlyHackSubType    string              `json:"sub-type,omitempty" yaml:"-,omitempty" jsonschema:"sub-type"`
}
