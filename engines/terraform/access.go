package terraform

import "github.com/samber/lo"

type ResourceType string

const (
	Database ResourceType = "database"
	Bucket   ResourceType = "bucket"
)

const AllAccess = "all"

var accessTypes = map[ResourceType][]string{
	Database: {"query", "mutate"},
	Bucket:   {"read", "write", "delete"},
}

// ExpandActions expands 'all' in an actions array whilst maintaining deduplication
func ExpandActions(actions []string, resourceType ResourceType) []string {
	var expanded []string

	for _, action := range actions {
		if action == AllAccess {
			expanded = append(expanded, accessTypes[resourceType]...)
		} else {
			expanded = append(expanded, action)
		}
	}

	return lo.Uniq(expanded)
}
