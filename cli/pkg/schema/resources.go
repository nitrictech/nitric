package schema

type BucketIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	BucketSchemaOnlyHackType string `json:"type" yaml:"-" jsonschema:"type,enum=bucket"`
}

type DatabaseIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	DatabaseSchemaOnlyHackType string `json:"type" yaml:"-" jsonschema:"type,enum=database"`
}

type StateIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	StateSchemaOnlyHackType string `json:"type" yaml:"-" jsonschema:"type,enum=state"`
}
