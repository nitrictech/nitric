package tags

import (
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
)

func Tags(stackID string, resourceName string, resourceType resources.ResourceType) map[string]string {
	return map[string]string{
		// Locate the unique stack by the presence of the key and the resource by its name
		GetResourceNameKey(stackID): resourceName,
		// Identifies the nitric resource type that led to the creation of the cloud resource.
		GetResourceTypeKey(stackID): string(resourceType),
	}
}

// GetResourceNameKey returns the key used to retrieve a resource's stack specific name from its tags.
func GetResourceNameKey(stackID string) string {
	return fmt.Sprintf("x-nitric-%s-name", stackID)
}

// GetResourceTypeKey returns the key used to retrieve a resource's stack specific type from its tags.
func GetResourceTypeKey(stackID string) string {
	return fmt.Sprintf("x-nitric-%s-type", stackID)
}
