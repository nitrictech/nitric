package pulumix

import (
	"fmt"
	"regexp"

	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

var (
	// URN
	childResourceUrn  = regexp.MustCompile(`(?:^|.+)nitric:(\w+)\$`)
	nitricResourceUrn = regexp.MustCompile(`(?:^|.+)nitric:(\w+)::([\w-]+)$`)
)

// PulumiUrn - Generate a standard Nitric Pulumi URN from a resource identifier
func PulumiUrn(nitricType resourcespb.ResourceType) string {
	return fmt.Sprintf("nitric:%s", resourcespb.ResourceType_name[int32(nitricType)])
}

// IsNitricResource - Checks if the Pulumi resource has a parent that is a Nitric Resource Type.
func IsNitricChildResource(pulumiUrn string) bool {
	return childResourceUrn.MatchString(pulumiUrn)
}

func IsNitricParentResource(pulumiUrn string) bool {
	return nitricResourceUrn.MatchString(pulumiUrn)
}

func NitricResourceIdFromPulumiUrn(pulumiUrn string) *resourcespb.ResourceIdentifier {
	urnGroups := nitricResourceUrn.FindStringSubmatch(pulumiUrn)
	if len(urnGroups) != 3 {
		return nil
	}

	nitricType := urnGroups[1]
	nitricName := urnGroups[2]

	resourceType := resourcespb.ResourceType_value[nitricType]

	return &resourcespb.ResourceIdentifier{
		Name: nitricName,
		Type: resourcespb.ResourceType(resourceType),
	}
}

// NitricResource - A logical Pulumi resources that represents a Nitric resource
// used to group concrete provider resources used to fulfill nitric resource deployments.
type NitricResource struct {
	pulumi.ResourceState
	Name string
	Type resourcespb.ResourceType
}

func ParentResourceFromResourceId(ctx *pulumi.Context, id *resourcespb.ResourceIdentifier) (pulumi.Resource, error) {
	res := &NitricResource{Name: id.Name, Type: id.Type}

	pulumiUrn := PulumiUrn(id.Type)

	err := ctx.RegisterComponentResource(pulumiUrn, id.Name, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
