package schema

type BucketIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	BucketSchemaOnlyHackType string `json:"type" yaml:"-" jsonschema:"type,enum=bucket"`
	Accessible
}

type StateIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	StateSchemaOnlyHackType string `json:"type" yaml:"-" jsonschema:"type,enum=state"`
}
