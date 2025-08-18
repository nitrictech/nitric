package schema

import "github.com/samber/lo"

type ResourceType string

const (
	Database ResourceType = "database"
	Bucket   ResourceType = "bucket"
)

const allAccess = "all"

var validActions = map[ResourceType][]string{
	Database: {"query", "mutate"},
	Bucket:   {"read", "write", "delete"},
}

func GetValidActions(resourceType ResourceType) []string {
	actions := validActions[resourceType]
	if actions == nil {
		return []string{}
	}

	return append(actions, allAccess)
}

// ExpandActions expands 'all' in an actions array whilst maintaining deduplication
func ExpandActions(actions []string, resourceType ResourceType) []string {
	expanded := []string{}

	for _, action := range actions {
		if action == allAccess {
			expanded = append(expanded, validActions[resourceType]...)
		} else {
			expanded = append(expanded, action)
		}
	}

	return lo.Uniq(expanded)
}

// ValidateActions ensures that all the actions are valid
func ValidateActions(actions []string, resourceType ResourceType) ([]string, bool) {
	invalidActions := []string{}

	for _, action := range actions {
		if !lo.Contains(GetValidActions(resourceType), action) {
			invalidActions = append(invalidActions, action)
		}
	}

	return invalidActions, len(invalidActions) == 0
}
