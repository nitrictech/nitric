package schema

type BucketIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	BucketSchemaOnlyHackType       string              `json:"type" yaml:"-" jsonschema:"type,enum=bucket"`
	BucketAccessSchemaOnlyHackType map[string][]string `json:"access,omitempty" yaml:"access,omitempty"`
	BucketSchemaOnlyHackSubType    string              `json:"sub-type,omitempty" yaml:"-,omitempty" jsonschema:"sub-type"`
}

type StateIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	StateSchemaOnlyHackType    string `json:"type" yaml:"-" jsonschema:"type,enum=state"`
	StateSchemaOnlyHackSubType string `json:"sub-type,omitempty" yaml:"-,omitempty" jsonschema:"sub-type"`
}
