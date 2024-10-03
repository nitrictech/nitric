// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	// Minimum of 3 AZs required for consistent cluster deployments
	azCount := 3

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

	a.RdsSecurityGroup, err = awsec2.NewSecurityGroup(ctx, "nitric-db-sg", &awsec2.SecurityGroupArgs{
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

func allVpcSubnetIds(vpc *ec2.Vpc) pulumi.StringArrayOutput {
	return pulumi.All(vpc.PrivateSubnetIds, vpc.PublicSubnetIds).ApplyT(func(args []interface{}) []string {
		subnets := []string{}
		privateSubnets := args[0].([]string)
		publicSubnets := args[1].([]string)

		subnets = append(subnets, privateSubnets...)
		subnets = append(subnets, publicSubnets...)

		return subnets
	}).(pulumi.StringArrayOutput)
}

// Apply rules to VPC security groups allowing them to mutually communicate
func (a *NitricAwsPulumiProvider) applyVpcRules(ctx *pulumi.Context) error {
	// Allow the database security group to communicate with the Batch security group
	if a.RdsSecurityGroup != nil && a.BatchSecurityGroup != nil {
		_, err := awsec2.NewSecurityGroupRule(ctx, "nitric-db-sg-to-batch-sg", &awsec2.SecurityGroupRuleArgs{
			SecurityGroupId:       a.RdsSecurityGroup.ID(),
			SourceSecurityGroupId: a.BatchSecurityGroup.ID(),
			FromPort:              pulumi.Int(0),
			ToPort:                pulumi.Int(0),
			Protocol:              pulumi.String("-1"),
			Type:                  pulumi.String("ingress"),
		})
		if err != nil {
			return err
		}
	}

	return nil
}
