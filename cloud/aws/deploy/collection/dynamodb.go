package collection

import (
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type DynamodbCollection struct {
	pulumi.ResourceState

	Table *dynamodb.Table
	Name  string
}

type DynamodbCollectionArgs struct {
	StackID    pulumi.StringInput
	Collection *v1.Collection
}

func NewDynamodbCollection(ctx *pulumi.Context, name string, args *DynamodbCollectionArgs, opts ...pulumi.ResourceOption) (*DynamodbCollection, error) {
	res := &DynamodbCollection{Name: name}

	err := ctx.RegisterComponentResource("nitric:collection:Dynamodb", name, res, opts...)
	if err != nil {
		return nil, err
	}

	res.Table, err = dynamodb.NewTable(ctx, name, &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("_pk"),
				Type: pulumi.String("S"),
			},
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("_sk"),
				Type: pulumi.String("S"),
			},
		},
		HashKey:     pulumi.String("_pk"),
		RangeKey:    pulumi.String("_sk"),
		BillingMode: pulumi.String("PAY_PER_REQUEST"),
		Tags:        tags.Tags(ctx, args.StackID, name),
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}
