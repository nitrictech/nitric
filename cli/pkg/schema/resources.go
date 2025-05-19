package schema

type BucketResource struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	BucketSchemaOnlyHackType string `json:"type" jsonschema:"type,enum=bucket"`
}

type DatabaseResource struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	DatabaseSchemaOnlyHackType string `json:"type" jsonschema:"type,enum=database"`
}

type StateResource struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	StateSchemaOnlyHackType string `json:"type" jsonschema:"type,enum=state"`
}
