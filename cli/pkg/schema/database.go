package schema

type DatabaseIntent struct {
	// Only used for schema generation, will always be nil. Do not use or remove.
	DatabaseSchemaOnlyHackType string `json:"type" yaml:"-" jsonschema:"type,enum=database"`
}
