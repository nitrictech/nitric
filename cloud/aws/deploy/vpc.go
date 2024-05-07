package deploy

import (
	"slices"

	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	pulumiAws "github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func (a *NitricAwsPulumiProvider) vpc(ctx *pulumi.Context) error {
	availabilityZones, err := pulumiAws.GetAvailabilityZones(ctx, &pulumiAws.GetAvailabilityZonesArgs{})
	if err != nil {
		return err
	}
	// Ensure AZ order is deterministic
	slices.Sort(availabilityZones.Names)

	// TODO: Make configurable
	azCount := 2

	a.VpcAzs = availabilityZones.Names[0:azCount]

	a.Vpc, err = ec2.NewVpc(ctx, "nitric-vpc", &ec2.VpcArgs{
		EnableDnsHostnames:    pulumi.Bool(true),
		AvailabilityZoneNames: a.VpcAzs,
		// These are quite expensive to run with (~$1.5/day/gateway)
		// with database compute on top of that
		// Replace with a VPC Endpoint if possible
		NatGateways: &ec2.NatGatewayConfigurationArgs{
			// TODO: Internet access with not be HA for resources on private subnets
			// If we remove this then Lambda instances deployed in this stack will
			// not be able to access external resources
			Strategy: ec2.NatGatewayStrategySingle,
		},
		Tags: pulumi.ToStringMap(tags.Tags(a.StackId, "vpc", "Vpc")),
	})
	if err != nil {
		return err
	}

	a.VpcSecurityGroup, err = awsec2.NewSecurityGroup(ctx, "nitric-db-sg", &awsec2.SecurityGroupArgs{
		VpcId: a.Vpc.VpcId,
		// Allow only incoming postgres SQL connections
		Ingress: awsec2.SecurityGroupIngressArray{
			&awsec2.SecurityGroupIngressArgs{
				FromPort: pulumi.Int(5432),
				ToPort:   pulumi.Int(5432),
				Protocol: pulumi.String("tcp"),
				Self:     pulumi.Bool(true),
			},
		},
		// Allow all outgoing traffic
		// TODO: Harden this
		Egress: awsec2.SecurityGroupEgressArray{
			&awsec2.SecurityGroupEgressArgs{
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				Protocol: pulumi.String("-1"),
				CidrBlocks: pulumi.StringArray{
					pulumi.String("0.0.0.0/0"),
				},
			},
		},
		Tags: pulumi.ToStringMap(tags.Tags(a.StackId, "vpc-database-security-group", "VpcSecurityGroup")),
	})
	if err != nil {
		return err
	}

	return nil
}
