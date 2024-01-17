package api

import (
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type NitricApi struct {
	pulumi.ResourceState

	Name string
}

type NitricApiArgs struct {
}

type ApiDeploymentFunc = func() error

func NewNitricApi(ctx *pulumi.Context, name string, args NitricApiArgs, opts ...pulumi.ResourceOption) (*NitricApi, error) {
	res := &NitricApi{Name: name}

	err := ctx.RegisterComponentResource(pulumix.PulumiUrn(resourcespb.ResourceType_Api), name, res, opts...)
	if err != nil {
		return nil, err
	}

	return res, err
}
