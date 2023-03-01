package tags

import "github.com/pulumi/pulumi/sdk/v3/go/pulumi"

func Tags(ctx *pulumi.Context, stackID pulumi.StringInput, name string) pulumi.StringMap {
	return pulumi.StringMap{
		"x-nitric-project":    pulumi.String(ctx.Project()),
		"x-nitric-stack":      stackID,
		"x-nitric-stack-name": pulumi.String(ctx.Stack()),
		"x-nitric-name":       pulumi.String(name),
	}
}
