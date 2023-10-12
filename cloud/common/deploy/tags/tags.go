package tags

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Tags(ctx *pulumi.Context, stackID string, name string) map[string]string {
	return map[string]string{
		"x-nitric-project":    ctx.Project(),
		"x-nitric-stack":      stackID,
		"x-nitric-stack-name": ctx.Stack(),
		"x-nitric-name":       name,
		// New tag for resource identification
		// Locate the unique stack by the presence of the key and the resource by its name
		GetResourceNameKey(stackID): name,
	}
}

// GetResourceNameKey returns the key used to retrieve a resource's stack specific name from its tags.
func GetResourceNameKey(stackID string) string {
	return fmt.Sprintf("x-nitric-stack-%s", stackID)
}
