package schema

type ServiceIntent struct {
	Env       map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Container Container         `json:"container" yaml:"container" jsonschema:"oneof_required=container"`

	Dev *Dev `json:"dev,omitempty" yaml:"dev,omitempty"`

	// Only used for schema generation, will always be nil. Do not use or remove.
	ServiceSchemaOnlyHackType string `json:"type" yaml:"-" jsonschema:"type,enum=service"`
	// TODO: should sub-type be sub_type?
	ServiceSchemaOnlyHackSubType string `json:"sub-type,omitempty" yaml:"-,omitempty" jsonschema:"sub-type"`
}

type Dev struct {
	// The script the start the service (because running it locally is orders of magnitude faster than building the containers)
	Script string `json:"script" yaml:"script"`
	// Watch  []string
}

// Runtime represents a union of all possible runtime types
type Container struct {
	Docker *Docker      `json:"docker,omitempty" yaml:"docker,omitempty" jsonschema:"oneof_required=docker"`
	Image  *DockerImage `json:"image,omitempty" yaml:"image,omitempty" jsonschema:"oneof_required=image"`
}

// DockerFileRuntime represents a runtime that uses a Dockerfile
type Docker struct {
	Dockerfile string            `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`
	Context    string            `json:"context,omitempty" yaml:"context,omitempty"`
	Args       map[string]string `json:"args,omitempty" yaml:"args,omitempty"`
}

type DockerImage struct {
	ID string `json:"id,omitempty" yaml:"id,omitempty" jsonschema:"required"`
}
