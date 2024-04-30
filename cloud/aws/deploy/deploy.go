// Copyright Nitric Pty Ltd.
//
// SPDX-License-Identifier: Apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
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
	"fmt"
	"strings"

	_ "embed"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	lambdaClient "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/nitrictech/nitric/cloud/aws/common"
	"github.com/nitrictech/nitric/cloud/aws/deploy/embeds"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	resourcespb "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	pulumiAws "github.com/pulumi/pulumi-aws/sdk/v5/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/resourcegroups"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/secretsmanager"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/sqs"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/codebuild"
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAwsPulumiProvider struct {
	*deploy.CommonStackDetails

	StackId   string
	AwsConfig *AwsConfig

	Vpc              *ec2.Vpc
	VpcSecurityGroup *awsec2.SecurityGroup
	// A codebuild job for creating the requested databases for a single database cluster
	CreateDatabaseProject *codebuild.Project
	DatabaseCluster       *rds.Cluster
	EcrAuthToken          *ecr.GetAuthorizationTokenResult
	Lambdas               map[string]*lambda.Function
	LambdaRoles           map[string]*iam.Role
	HttpProxies           map[string]*apigatewayv2.Api
	Apis                  map[string]*apigatewayv2.Api
	Secrets               map[string]*secretsmanager.Secret
	Buckets               map[string]*s3.Bucket
	BucketNotifications   map[string]*s3.BucketNotification
	Topics                map[string]*topic
	Queues                map[string]*sqs.Queue
	Websockets            map[string]*apigatewayv2.Api
	KeyValueStores        map[string]*dynamodb.Table

	provider.NitricDefaultOrder

	ResourceTaggingClient *resourcegroupstaggingapi.ResourceGroupsTaggingAPI
	LambdaClient          lambdaiface.LambdaAPI
}

var _ provider.NitricPulumiProvider = (*NitricAwsPulumiProvider)(nil)

const pulumiAwsVersion = "6.6.0"

func (a *NitricAwsPulumiProvider) Config() (auto.ConfigMap, error) {
	return auto.ConfigMap{
		"aws:region":     auto.ConfigValue{Value: a.Region},
		"aws:version":    auto.ConfigValue{Value: pulumiAwsVersion},
		"docker:version": auto.ConfigValue{Value: deploy.PulumiDockerVersion},
	}, nil
}

func (a *NitricAwsPulumiProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	a.AwsConfig, err = ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

func (a *NitricAwsPulumiProvider) Pre(ctx *pulumi.Context, resources []*pulumix.NitricPulumiResource[any]) error {
	// make our random stackId
	stackRandId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-stack-name", ctx.Stack()), &random.RandomStringArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(8),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"stack-name": ctx.Stack(),
		}),
	})
	if err != nil {
		return err
	}

	stackIdChan := make(chan string)
	pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.Result).ApplyT(func(id string) string {
		stackIdChan <- id
		return id
	})

	a.StackId = <-stackIdChan

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	a.ResourceTaggingClient = resourcegroupstaggingapi.New(sess)

	a.LambdaClient = lambdaClient.New(sess, &aws.Config{Region: aws.String(a.Region)})

	a.EcrAuthToken, err = ecr.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenArgs{})
	if err != nil {
		return err
	}

	// Create AWS Resource groups with our tags
	_, err = resourcegroups.NewGroup(ctx, "stack-resource-group", &resourcegroups.GroupArgs{
		Name:        pulumi.String(a.FullStackName),
		Description: pulumi.Sprintf("All deployed resources for the %s nitric stack", a.FullStackName),
		ResourceQuery: &resourcegroups.GroupResourceQueryArgs{
			Query: pulumi.Sprintf(`{
				"ResourceTypeFilters":["AWS::AllSupported"],
				"TagFilters":[{"Key":"%s"}]
			}`, tags.GetResourceNameKey(a.StackId)),
		},
	})

	databases := lo.Filter(resources, func(item *pulumix.NitricPulumiResource[any], idx int) bool {
		return item.Id.Type == resourcespb.ResourceType_SqlDatabase
	})
	// Create a shared database cluster if we have more than one database
	if len(databases) > 0 {
		// deploy a VPC and security groups for this database cluster

		availabilityZones, err := pulumiAws.GetAvailabilityZones(ctx, &pulumiAws.GetAvailabilityZonesArgs{})
		if err != nil {
			return err
		}

		// TODO: Make configurable
		azCount := 1

		selectedAzs := availabilityZones.Names[0:azCount]

		clusterAvailabilityZones := pulumi.StringArray{}

		for _, az := range selectedAzs {
			clusterAvailabilityZones = append(clusterAvailabilityZones, pulumi.String(az))
		}

		// generate a db cluster random password
		dbMasterPassword, err := random.NewRandomPassword(ctx, "db-master-password", &random.RandomPasswordArgs{
			Length:          pulumi.Int(16),
			OverrideSpecial: pulumi.String("!#$%^&*()"),
		})
		if err != nil {
			return err
		}

		a.Vpc, err = ec2.NewVpc(ctx, "nitric-vpc", &ec2.VpcArgs{
			EnableDnsHostnames:    pulumi.Bool(true),
			AvailabilityZoneNames: selectedAzs,
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
		})
		if err != nil {
			return err
		}

		dbSubnetGroup, err := rds.NewSubnetGroup(ctx, "dbsubnetgroup", &rds.SubnetGroupArgs{
			SubnetIds: a.Vpc.PrivateSubnetIds,
		})
		if err != nil {
			return err
		}

		a.DatabaseCluster, err = rds.NewCluster(ctx, "postgresql", &rds.ClusterArgs{
			Engine:        pulumi.String(rds.EngineTypeAuroraPostgresql),
			EngineVersion: pulumi.String("13.14"),
			// TODO: limit number of availability zones
			AvailabilityZones: clusterAvailabilityZones,
			DatabaseName:      pulumi.String("nitric"),
			MasterUsername:    pulumi.String("nitric"),
			MasterPassword:    dbMasterPassword.Result,
			EngineMode:        pulumi.String(rds.EngineModeProvisioned),
			Serverlessv2ScalingConfiguration: &rds.ClusterServerlessv2ScalingConfigurationArgs{
				MaxCapacity: pulumi.Float64(1),
				MinCapacity: pulumi.Float64(0.5),
			},
			// TODO: Validate timezones used here
			// PreferredBackupWindow: pulumi.String("07:00-09:00"),
			VpcSecurityGroupIds: pulumi.StringArray{a.VpcSecurityGroup.ID()},
			DbSubnetGroupName:   dbSubnetGroup.Name,
			SkipFinalSnapshot:   pulumi.Bool(true),
		})
		if err != nil {
			return err
		}

		dbInstance, err := rds.NewClusterInstance(ctx, "example", &rds.ClusterInstanceArgs{
			ClusterIdentifier: a.DatabaseCluster.ID(),
			InstanceClass:     pulumi.String("db.serverless"),
			Engine:            a.DatabaseCluster.Engine,
			EngineVersion:     a.DatabaseCluster.EngineVersion,
			DbSubnetGroupName: a.DatabaseCluster.DbSubnetGroupName,
		})
		if err != nil {
			return err
		}

		role, err := iam.NewRole(ctx, "codeBuildRole", &iam.RoleArgs{
			AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Action": "sts:AssumeRole",
						"Principal": {
							"Service": "codebuild.amazonaws.com"
						},
						"Effect": "Allow",
						"Sid": ""
					}
				]
			}`),
		})
		if err != nil {
			return err
		}

		codebuildManagedPolicies := map[string]iam.ManagedPolicy{
			"codeBuildAdmin": iam.ManagedPolicyAWSCodeBuildAdminAccess,
			"rdsAdmin":       iam.ManagedPolicyAmazonRDSFullAccess,
			"ec2Admin":       iam.ManagedPolicyAmazonEC2FullAccess,
			"cloudWatchLogs": iam.ManagedPolicyCloudWatchLogsFullAccess,
		}

		for name, policy := range codebuildManagedPolicies {
			_, err = iam.NewRolePolicyAttachment(ctx, name+"PolicyAttachment", &iam.RolePolicyAttachmentArgs{
				Role:      role.Name,
				PolicyArn: policy,
			})
			if err != nil {
				return err
			}
		}

		// Attach the AWSCodeBuildDeveloperAccess policy to the role
		_, err = iam.NewRolePolicyAttachment(ctx, "codeBuildPolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: iam.ManagedPolicyAWSCodeBuildAdminAccess,
		})
		if err != nil {
			return err
		}

		// Attach the VPC access policy to the role
		_, err = iam.NewRolePolicyAttachment(ctx, "codeBuildRdsPolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: iam.ManagedPolicyAmazonRDSFullAccess,
		})
		if err != nil {
			return err
		}

		// Attach the VPC access policy to the role
		_, err = iam.NewRolePolicyAttachment(ctx, "codeBuildEc2PolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      role.Name,
			PolicyArn: iam.ManagedPolicyAmazonEC2FullAccess,
		})
		if err != nil {
			return err
		}

		// Use a codebuild project to create the databases within the cluster
		a.CreateDatabaseProject, err = codebuild.NewProject(ctx, "create-nitric-databases", &codebuild.ProjectArgs{
			Artifacts: &codebuild.ProjectArtifactsArgs{
				Type: pulumi.String("NO_ARTIFACTS"),
			},
			Environment: &codebuild.ProjectEnvironmentArgs{
				ComputeType: pulumi.String("BUILD_GENERAL1_SMALL"),
				Image:       pulumi.String("aws/codebuild/amazonlinux2-x86_64-standard:4.0"),
				Type:        pulumi.String("LINUX_CONTAINER"),
				EnvironmentVariables: codebuild.ProjectEnvironmentEnvironmentVariableArray{
					&codebuild.ProjectEnvironmentEnvironmentVariableArgs{
						Name:  pulumi.String("DB_CLUSTER_ENDPOINT"),
						Value: a.DatabaseCluster.Endpoint,
					},
					&codebuild.ProjectEnvironmentEnvironmentVariableArgs{
						Name:  pulumi.String("DB_MASTER_USERNAME"),
						Value: pulumi.String("nitric"),
					},
					&codebuild.ProjectEnvironmentEnvironmentVariableArgs{
						Name:  pulumi.String("DB_MASTER_PASSWORD"),
						Value: dbMasterPassword.Result,
					},
				},
			},
			ServiceRole: role.Arn,
			Source: &codebuild.ProjectSourceArgs{
				Type:      pulumi.String("NO_SOURCE"),
				Buildspec: embeds.GetCodeBuildCreateDatabaseConfig(),
			},
			VpcConfig: &codebuild.ProjectVpcConfigArgs{
				SecurityGroupIds: a.DatabaseCluster.VpcSecurityGroupIds,
				Subnets:          a.Vpc.PrivateSubnetIds,
				VpcId:            a.Vpc.VpcId,
			},
			// Don't deploy the build until after the database cluster is completely ready
		}, pulumi.DependsOn([]pulumi.Resource{a.DatabaseCluster, dbInstance}))
		if err != nil {
			return err
		}
	}

	return err
}

func (a *NitricAwsPulumiProvider) Post(ctx *pulumi.Context) error {
	return nil
}

func (a *NitricAwsPulumiProvider) Result(ctx *pulumi.Context) (pulumi.StringOutput, error) {
	outputs := []interface{}{}

	// Add APIs outputs
	if len(a.Apis) > 0 {
		outputs = append(outputs, pulumi.Sprintf("API Endpoints:\n──────────────"))
		for apiName, api := range a.Apis {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", apiName, api.ApiEndpoint))
		}
	}

	// Add HTTP Proxy outputs
	if len(a.HttpProxies) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("HTTP Proxies:\n──────────────"))
		for proxyName, proxy := range a.HttpProxies {
			outputs = append(outputs, pulumi.Sprintf("%s: %s", proxyName, proxy.ApiEndpoint))
		}
	}

	// Add Websocket outputs
	if len(a.Websockets) > 0 {
		if len(outputs) > 0 {
			outputs = append(outputs, "\n")
		}
		outputs = append(outputs, pulumi.Sprintf("Websockets:\n──────────────"))
		for wsName, ws := range a.Websockets {
			outputs = append(outputs, pulumi.Sprintf("%s: %s/%s", wsName, ws.ApiEndpoint, common.DefaultWsStageName))
		}
	}

	output, ok := pulumi.All(outputs...).ApplyT(func(deets []interface{}) string {
		stringyOutputs := make([]string, len(deets))
		for i, d := range deets {
			stringyOutputs[i] = d.(string)
		}

		return strings.Join(stringyOutputs, "\n")
	}).(pulumi.StringOutput)

	if !ok {
		return pulumi.StringOutput{}, fmt.Errorf("failed to generate pulumi output")
	}

	return output, nil
}

func NewNitricAwsProvider() *NitricAwsPulumiProvider {
	return &NitricAwsPulumiProvider{
		Lambdas:             make(map[string]*lambda.Function),
		LambdaRoles:         make(map[string]*iam.Role),
		Apis:                make(map[string]*apigatewayv2.Api),
		HttpProxies:         make(map[string]*apigatewayv2.Api),
		Secrets:             make(map[string]*secretsmanager.Secret),
		Buckets:             make(map[string]*s3.Bucket),
		BucketNotifications: make(map[string]*s3.BucketNotification),
		Websockets:          make(map[string]*apigatewayv2.Api),
		Topics:              make(map[string]*topic),
		Queues:              make(map[string]*sqs.Queue),
		KeyValueStores:      make(map[string]*dynamodb.Table),
	}
}
