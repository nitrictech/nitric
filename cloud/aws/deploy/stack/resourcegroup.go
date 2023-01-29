package stack

import (
	"encoding/json"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/resourcegroups"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AwsResourceGroup struct {
	pulumi.ResourceState
	Name          string
	ResourceGroup *resourcegroups.Group
}

type AwsResourceGroupArgs struct {
	StackID pulumi.StringInput
}

func NewAwsResourceGroup(ctx *pulumi.Context, name string, args *AwsResourceGroupArgs, opts ...pulumi.ResourceOption) (*AwsResourceGroup, error) {
	res := &AwsResourceGroup{Name: name}

	err := ctx.RegisterComponentResource("nitric:stack:AwsResourceGroup", name, res, opts...)
	if err != nil {
		return nil, err
	}

	rgQueryJSON := args.StackID.ToStringOutput().ApplyT(func(sid string) (string, error) {
		b, err := json.Marshal(map[string]interface{}{
			"ResourceTypeFilters": []string{"AWS::AllSupported"},
			"TagFilters": []interface{}{
				map[string]interface{}{
					"Key":    "x-nitric-stack",
					"Values": []string{sid},
				},
			},
		})
		if err != nil {
			return "", err
		}

		return string(b), nil
	}).(pulumi.StringOutput)

	res.ResourceGroup, err = resourcegroups.NewGroup(ctx, ctx.Stack(), &resourcegroups.GroupArgs{
		Description: pulumi.Sprintf("Nitric RG for stack %s", res.Name),
		ResourceQuery: &resourcegroups.GroupResourceQueryArgs{
			Query: rgQueryJSON,
		},
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
