package triggers

// SourceType enum
type TriggerType int

const (
	TriggerType_Subscription TriggerType = iota
	TriggerType_Request
	TriggerType_Custom
)

func (e TriggerType) String() string {
	return []string{"SUBSCRIPTION", "REQUEST", "CUSTOM"}[e]
}
