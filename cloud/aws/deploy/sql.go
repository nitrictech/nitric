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
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/codebuild"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
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
			if status == awscodebuild.StatusTypeFailed {
				return retry.Unrecoverable(fmt.Errorf("codebuild job %s failed", buildId))
			}

			return nil
		}

		fmt.Printf("Waiting for codebuild job %s to finish\n", buildId)
		return fmt.Errorf("build %s still in progress", buildId)
	}
}

// Sqldatabase - Implements PostgresSql database deployments use AWS Aurora
func (a *NitricAwsPulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {
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

		cmd, _, err := image.CommandFromImageInspect(config.GetImageUri())
		if err != nil {
			return err
		}

		// Push the migration image to ECR using pulumi
		_, err = docker.NewTag(ctx, name, &docker.TagArgs{
			SourceImage: pulumi.String(config.GetImageUri()),
			TargetImage: repo.RepositoryUrl,
		})
		if err != nil {
			return err
		}

		// push to the ECR repository
		image, err := docker.NewRegistryImage(ctx, name, &docker.RegistryImageArgs{
			Name: repo.RepositoryUrl,
		})
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
				ComputeType: pulumi.String("BUILD_GENERAL1_SMALL"),
				Image:       image.Name,
				Type:        pulumi.String("LINUX_CONTAINER"),
			},
			ServiceRole: a.CodeBuildRole.Arn,
			Source: &codebuild.ProjectSourceArgs{
				Type:      pulumi.String("NO_SOURCE"),
				Buildspec: embeds.GetCodeBuildMigrateDatabaseConfig(cmd),
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

	pulumi.All(a.CreateDatabaseProject.Name, a.DatabaseMigrationJobs[config.GetImageUri()].Name, dbUrl).ApplyT(func(args []interface{}) (bool, error) {
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

		// Run the database migration step
		out, err = client.StartBuild(&awscodebuild.StartBuildInput{
			ProjectName: aws.String(migrateDatabaseJob),
			EnvironmentVariablesOverride: []*awscodebuild.EnvironmentVariable{
				{
					Name:  aws.String("DB_URL"),
					Value: aws.String(databaseUrl),
				},
			},
		})

		err = retry.Do(checkBuildStatus(client, *out.Build.Id), retry.Attempts(10), retry.Delay(time.Second*15))
		if err != nil {
			return false, err
		}

		return true, nil
	})

	// Run the database creation step
	// a.CreateDatabaseProject.Name.ApplyT(func(projectName string) (bool, error) {
	// 	fmt.Printf("Starting database creation build %s\n", name)
	// 	out, err := client.StartBuild(&awscodebuild.StartBuildInput{
	// 		ProjectName: aws.String(projectName),
	// 		EnvironmentVariablesOverride: []*awscodebuild.EnvironmentVariable{
	// 			{
	// 				Name:  aws.String("DB_NAME"),
	// 				Value: aws.String(name),
	// 			},
	// 		},
	// 	})
	// 	if err != nil {
	// 		return false, err
	// 	}

	// 	var finalErr error
	// 	err = retry.Do(checkBuildStatus(client, out.Build.Id), retry.Attempts(10), retry.Delay(time.Second*15))
	// 	if err != nil {
	// 		return false, err
	// 	}

	// 	if finalErr != nil {
	// 		return false, finalErr
	// 	}

	// 	return true, nil
	// })

	return nil
}
