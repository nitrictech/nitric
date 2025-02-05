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
	"time"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awscodebuild "github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/nitrictech/nitric/cloud/aws/deploy/embeds"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/codebuild"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/rds"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
)

func checkBuildStatus(client *awscodebuild.CodeBuild, buildId string) func() error {
	return func() error {
		resp, err := client.BatchGetBuilds(&awscodebuild.BatchGetBuildsInput{
			Ids: []*string{
				aws.String(buildId),
			},
		})
		if err != nil {
			return err
		}

		status := aws.StringValue(resp.Builds[0].BuildStatus)
		if status != awscodebuild.StatusTypeInProgress {
			if status != awscodebuild.StatusTypeSucceeded {
				return retry.Unrecoverable(fmt.Errorf("codebuild job %s failed", buildId))
			}

			return nil
		}

		fmt.Printf("Waiting for codebuild job %s to finish\n", buildId)
		return fmt.Errorf("build %s still in progress", buildId)
	}
}

func (a *NitricAwsPulumiProvider) rds(ctx *pulumi.Context) error {
	var err error

	a.DbMasterPassword, err = random.NewRandomPassword(ctx, "db-master-password", &random.RandomPasswordArgs{
		Length:  pulumi.Int(16),
		Special: pulumi.Bool(false),
	})
	if err != nil {
		return err
	}

	dbSubnetGroup, err := rds.NewSubnetGroup(ctx, "dbsubnetgroup", &rds.SubnetGroupArgs{
		SubnetIds: a.Vpc.PrivateSubnetIds,
		Tags:      pulumi.ToStringMap(tags.Tags(a.StackId, "database-subnet-group", "database-subnet-group")),
	})
	if err != nil {
		return err
	}

	scaling := rds.ClusterServerlessv2ScalingConfigurationArgs{
		MinCapacity: pulumi.Float64(a.AwsConfig.AuroraRdsClusterConfig.MinCapacity),
		MaxCapacity: pulumi.Float64(a.AwsConfig.AuroraRdsClusterConfig.MaxCapacity),
	}

	if a.AwsConfig.AuroraRdsClusterConfig.SecondsUntilAutoPause != nil && *a.AwsConfig.AuroraRdsClusterConfig.SecondsUntilAutoPause > 0 {
		if a.AwsConfig.AuroraRdsClusterConfig.MinCapacity != 0 {
			return fmt.Errorf("seconds-until-auto-pause can only be set when min-capacity is 0")
		}
		scaling.SecondsUntilAutoPause = pulumi.Int(*a.AwsConfig.AuroraRdsClusterConfig.SecondsUntilAutoPause)
	}

	a.DatabaseCluster, err = rds.NewCluster(ctx, "postgresql", &rds.ClusterArgs{
		ApplyImmediately: pulumi.Bool(true),
		Engine:           pulumi.String(rds.EngineTypeAuroraPostgresql),
		EngineVersion:    pulumi.String("13.16"),
		// TODO: limit number of availability zones
		AvailabilityZones:                pulumi.ToStringArray(a.VpcAzs),
		DatabaseName:                     pulumi.String("nitric"),
		MasterUsername:                   pulumi.String("nitric"),
		MasterPassword:                   a.DbMasterPassword.Result,
		EngineMode:                       pulumi.String(rds.EngineModeProvisioned),
		Serverlessv2ScalingConfiguration: &scaling,
		VpcSecurityGroupIds:              pulumi.StringArray{a.RdsSecurityGroup.ID()},
		DbSubnetGroupName:                dbSubnetGroup.Name,
		SkipFinalSnapshot:                pulumi.Bool(true),
		Tags:                             pulumi.ToStringMap(tags.Tags(a.StackId, "database-cluster", "DatabaseCluster")),
		// NOTE: Workaround for https://github.com/pulumi/pulumi-aws/issues/2426
		// Aurora instances don't support StorageType so we need to ignore changes otherwise we'll get unsolicited replacements
	}, pulumi.IgnoreChanges([]string{"storageType"}))
	if err != nil {
		return err
	}

	dbInstance, err := rds.NewClusterInstance(ctx, "example", &rds.ClusterInstanceArgs{
		ClusterIdentifier: a.DatabaseCluster.ID(),
		InstanceClass:     pulumi.String("db.serverless"),
		Engine:            a.DatabaseCluster.Engine,
		EngineVersion:     a.DatabaseCluster.EngineVersion,
		DbSubnetGroupName: a.DatabaseCluster.DbSubnetGroupName,
		Tags:              pulumi.ToStringMap(tags.Tags(a.StackId, "database-cluster-instance", "DatabaseInstance")),
	})
	if err != nil {
		return err
	}

	a.CodeBuildRole, err = iam.NewRole(ctx, "codeBuildRole", &iam.RoleArgs{
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
		"ecrReadonly":    iam.ManagedPolicyAmazonEC2ContainerRegistryReadOnly,
	}

	for name, policy := range codebuildManagedPolicies {
		_, err = iam.NewRolePolicyAttachment(ctx, name+"PolicyAttachment", &iam.RolePolicyAttachmentArgs{
			Role:      a.CodeBuildRole.Name,
			PolicyArn: policy,
		})
		if err != nil {
			return err
		}
	}

	// Attach the AWSCodeBuildDeveloperAccess policy to the role
	_, err = iam.NewRolePolicyAttachment(ctx, "codeBuildPolicyAttachment", &iam.RolePolicyAttachmentArgs{
		Role:      a.CodeBuildRole.Name,
		PolicyArn: iam.ManagedPolicyAWSCodeBuildAdminAccess,
	})
	if err != nil {
		return err
	}

	// Attach the VPC access policy to the role
	_, err = iam.NewRolePolicyAttachment(ctx, "codeBuildRdsPolicyAttachment", &iam.RolePolicyAttachmentArgs{
		Role:      a.CodeBuildRole.Name,
		PolicyArn: iam.ManagedPolicyAmazonRDSFullAccess,
	})
	if err != nil {
		return err
	}

	// Attach the VPC access policy to the role
	_, err = iam.NewRolePolicyAttachment(ctx, "codeBuildEc2PolicyAttachment", &iam.RolePolicyAttachmentArgs{
		Role:      a.CodeBuildRole.Name,
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
					Value: a.DbMasterPassword.Result,
				},
			},
		},
		ServiceRole: a.CodeBuildRole.Arn,
		Source: &codebuild.ProjectSourceArgs{
			Type:      pulumi.String("NO_SOURCE"),
			Buildspec: embeds.GetCodeBuildCreateDatabaseConfig(),
		},
		VpcConfig: &codebuild.ProjectVpcConfigArgs{
			SecurityGroupIds: a.DatabaseCluster.VpcSecurityGroupIds,
			Subnets:          a.Vpc.PrivateSubnetIds,
			VpcId:            a.Vpc.VpcId,
		},
		Tags: pulumi.ToStringMap(tags.Tags(a.StackId, "database-migration-job", "Job")),
	}, pulumi.DependsOn([]pulumi.Resource{a.DatabaseCluster, dbInstance}))
	if err != nil {
		return err
	}

	return nil
}

// A customer SQL database pulumi resource
type RdsDatabase struct {
	pulumi.ResourceState

	Name string

	Migrated pulumi.BoolOutput
}

// Sqldatabase - Implements PostgresSql database deployments use AWS Aurora
func (a *NitricAwsPulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {
	a.SqlDatabases[name] = &RdsDatabase{
		Name: name,
	}
	err := ctx.RegisterComponentResource("nitricaws:SqlDatabase", name, a.SqlDatabases[name], pulumi.Parent(parent))
	if err != nil {
		return err
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.Region), // replace with your AWS region
	})
	if err != nil {
		return err
	}

	client := awscodebuild.New(sess)

	if config.GetImageUri() != "" && a.DatabaseMigrationJobs[config.GetImageUri()] == nil {
		repo, err := ecr.NewRepository(ctx, name, &ecr.RepositoryArgs{
			ForceDelete: pulumi.BoolPtr(true),
			Tags:        pulumi.ToStringMap(tags.Tags(a.StackId, name, resources.SqlDatabase)),
		}, pulumi.Parent(parent))
		if err != nil {
			return err
		}

		inspect, err := image.CommandFromImageInspect(config.GetImageUri(), " ")
		if err != nil {
			return err
		}

		newTag, err := docker.NewTag(ctx, name+"-tag", &docker.TagArgs{
			SourceImage: pulumi.String(inspect.ID),
			TargetImage: repo.RepositoryUrl,
		}, pulumi.Parent(parent))
		if err != nil {
			return err
		}

		image, err := docker.NewRegistryImage(ctx, name+"-remote", &docker.RegistryImageArgs{
			Name: repo.RepositoryUrl,
			Triggers: pulumi.Map{
				"imageSha": pulumi.String(inspect.ID),
			},
		}, pulumi.Parent(parent), pulumi.Provider(a.DockerProvider), pulumi.DependsOn([]pulumi.Resource{newTag}))
		if err != nil {
			return err
		}

		// Create a new codebuild job for the image
		a.DatabaseMigrationJobs[config.GetImageUri()], err = codebuild.NewProject(ctx, name, &codebuild.ProjectArgs{
			Name: pulumi.String(name),
			Artifacts: &codebuild.ProjectArtifactsArgs{
				Type: pulumi.String("NO_ARTIFACTS"),
			},
			Environment: &codebuild.ProjectEnvironmentArgs{
				ComputeType:              pulumi.String("BUILD_GENERAL1_SMALL"),
				Image:                    pulumi.Sprintf("%s@%s", repo.RepositoryUrl, image.Sha256Digest),
				ImagePullCredentialsType: pulumi.String("SERVICE_ROLE"),
				Type:                     pulumi.String("LINUX_CONTAINER"),
			},
			ServiceRole: a.CodeBuildRole.Arn,
			Source: &codebuild.ProjectSourceArgs{
				Type:      pulumi.String("NO_SOURCE"),
				Buildspec: embeds.GetCodeBuildMigrateDatabaseConfig(lo.Ternary(inspect.WorkDir != "", inspect.WorkDir, "/"), fmt.Sprintf("'%s'", inspect.Cmd)),
			},
			VpcConfig: &codebuild.ProjectVpcConfigArgs{
				SecurityGroupIds: a.DatabaseCluster.VpcSecurityGroupIds,
				Subnets:          a.Vpc.PrivateSubnetIds,
				VpcId:            a.Vpc.VpcId,
			},
		})
		if err != nil {
			return err
		}
	}

	dbUrl := pulumi.Sprintf("postgres://%s:%s@%s:%s/%s", "nitric", a.DbMasterPassword.Result, a.DatabaseCluster.Endpoint, "5432", name)

	databaseMigrationJobName := pulumi.String("").ToStringOutput()
	if config.GetImageUri() != "" && a.DatabaseMigrationJobs[config.GetImageUri()] != nil {
		databaseMigrationJobName = a.DatabaseMigrationJobs[config.GetImageUri()].Name
	}

	a.SqlDatabases[name].Migrated = pulumi.All(a.CreateDatabaseProject.Name, databaseMigrationJobName, dbUrl).ApplyT(func(args []interface{}) (bool, error) {
		createDatabaseProject := args[0].(string)
		migrateDatabaseJob := args[1].(string)
		databaseUrl := args[2].(string)

		// Run the database creation step
		out, err := client.StartBuild(&awscodebuild.StartBuildInput{
			ProjectName: aws.String(createDatabaseProject),
			EnvironmentVariablesOverride: []*awscodebuild.EnvironmentVariable{
				{
					Name:  aws.String("DB_NAME"),
					Value: aws.String(name),
				},
			},
		})
		if err != nil {
			return false, err
		}

		err = retry.Do(checkBuildStatus(client, *out.Build.Id), retry.Attempts(10), retry.Delay(time.Second*15))
		if err != nil {
			return false, err
		}

		// Run the database migration step if the migration job exists
		if migrateDatabaseJob != "" {
			out, err = client.StartBuild(&awscodebuild.StartBuildInput{
				ProjectName: aws.String(migrateDatabaseJob),
				EnvironmentVariablesOverride: []*awscodebuild.EnvironmentVariable{
					{
						Name:  aws.String("NITRIC_DB_NAME"),
						Value: aws.String(name),
					},
					{
						Name:  aws.String("DB_URL"),
						Value: aws.String(databaseUrl),
					},
				},
			})
			if err != nil {
				return false, err
			}

			err = retry.Do(checkBuildStatus(client, *out.Build.Id), retry.Attempts(10), retry.Delay(time.Second*15), retry.LastErrorOnly(true))
			if err != nil {
				return false, fmt.Errorf("database migrations failed for %s: %s", name, err.Error())
			}
		}

		return true, nil
	}).(pulumi.BoolOutput)

	err = ctx.RegisterResourceOutputs(a.SqlDatabases[name], pulumi.Map{
		"migrated": a.SqlDatabases[name].Migrated,
	})
	if err != nil {
		return err
	}

	// if the result isn't exported the error will be ignored, causing failed migrations to be missed
	ctx.Export(name+"migrationsRun", a.SqlDatabases[name].Migrated)

	return nil
}
