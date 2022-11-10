package resourcegroupstaggingapiiface

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
)

type ResourceGroupsTaggingAPIAPI interface {
	GetResources(ctx context.Context, params *resourcegroupstaggingapi.GetResourcesInput, optFns ...func(*resourcegroupstaggingapi.Options)) (*resourcegroupstaggingapi.GetResourcesOutput, error)
}
