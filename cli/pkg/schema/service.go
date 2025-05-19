package schema

type ServiceResource struct {
	Port      int               `json:"port"`
	Env       map[string]string `json:"env,omitempty"`
	Container Container         `json:"container" jsonschema:"oneof_required=container"`
	// Only used for schema generation, will always be nil. Do not use or remove.
	ServiceSchemaOnlyHackType string `json:"type" jsonschema:"type,enum=service"`
}

// Runtime represents a union of all possible runtime types
type Container struct {
	Docker *Docker      `json:"docker,omitempty" jsonschema:"oneof_required=docker"`
	Image  *DockerImage `json:"image,omitempty" jsonschema:"oneof_required=image"`
}

// DockerFileRuntime represents a runtime that uses a Dockerfile
type Docker struct {
	Dockerfile string `json:"dockerfile,omitempty"`
	Context    string `json:"context,omitempty"`
}

type DockerImage struct {
	ID string `json:"id,omitempty" jsonschema:"required"`
}
