package sources

// SourceType enum
type SourceType int

const (
	SourceType_Subscription SourceType = iota
	SourceType_Request
	SourceType_Custom
)

func (e SourceType) String() string {
	return []string{"SUBSCRIPTION", "REQUEST", "CUSTOM"}[e]
}
