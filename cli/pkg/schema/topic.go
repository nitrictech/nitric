package schema

type Topic struct {
	Name     string         `json:"name"`
	Triggers []TopicTrigger `json:"triggers,omitempty"`
}

type TopicTriggerType string

type TopicTrigger struct {
	Target TopicTriggerTarget `json:"target"`
}

type TopicTriggerTargetType string

const (
	TopicTriggerTargetType_Service TopicTriggerTargetType = "service"
)

type TopicTriggerTarget struct {
	TargetType TopicTriggerTargetType `json:"type" jsonschema:"enum=service"`
	TargetName string                 `json:"name"`
	Path       string                 `json:"path"`
}
