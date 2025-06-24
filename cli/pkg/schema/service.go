package schema

type ServiceIntent struct {
	Resource
	Env       map[string]string `json:"env,omitempty" yaml:"env,omitempty"`
	Container Container         `json:"container" yaml:"container" jsonschema:"oneof_required=container"`

	Dev *Dev `json:"dev,omitempty" yaml:"dev,omitempty"`

	Triggers map[string]*ServiceTrigger `json:"triggers,omitempty" yaml:"triggers,omitempty"`
}

func (s *ServiceIntent) GetType() string {
	return "service"
}

type ServiceTrigger struct {
	Schedule *Schedule `json:"schedule,omitempty" yaml:"schedule,omitempty" jsonschema:"oneof_required=schedule"`
	// TODO: Add additional trigger types
	// Topic    *TopicTrigger  `json:"topic,omitempty" yaml:"topic,omitempty" jsonschema:"oneof_required=topic"`
	// Bucket   *BucketTrigger `json:"bucket,omitempty" yaml:"bucket,omitempty" jsonschema:"oneof_required=bucket"`

	Path string `json:"path" yaml:"path" jsonschema:"required"`
}

type Schedule struct {
	CronExpression string `json:"cron_expression,omitempty" yaml:"cron_expression,omitempty" jsonschema:"oneof_required=cron_expression"`
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
