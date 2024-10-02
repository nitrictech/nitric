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
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-awsx/sdk/go/awsx/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/batch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	awsec2 "github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
)

type ResourceRequirement struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (a *NitricAwsPulumiProvider) batch(ctx *pulumi.Context) error {
	var subnets pulumi.StringArrayOutput
	var vpcId pulumi.StringOutput
	var err error

	if a.Vpc != nil {
		allSubnets := allVpcSubnetIds(a.Vpc)

		subnets = allSubnets
		vpcId = a.Vpc.VpcId
	} else {
		vpc, err := ec2.NewDefaultVpc(ctx, "default-vpc", nil)
		if err != nil {
			return fmt.Errorf("could not resolve default VPC")
		}

		subnets = vpc.PublicSubnetIds
		vpcId = vpc.VpcId
	}

	a.BatchSecurityGroup, err = awsec2.NewSecurityGroup(ctx, "batch-sg", &awsec2.SecurityGroupArgs{
		VpcId: vpcId,
		Egress: awsec2.SecurityGroupEgressArray{
			awsec2.SecurityGroupEgressArgs{
				Protocol: pulumi.String("-1"),
				FromPort: pulumi.Int(0),
				ToPort:   pulumi.Int(0),
				// Still need public internet access for batch jobs by default
				// TODO: Allow restriction via config
				CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
			},
		},
	})
	if err != nil {
		return err
	}

	ecsInstanceRole, err := iam.NewRole(ctx, "EcsInstanceRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": "ec2.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`),
	})
	if err != nil {
		return err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, "EcsInstanceRolePolicyAttachment", &iam.RolePolicyAttachmentArgs{
		Role:      ecsInstanceRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"),
	})
	if err != nil {
		return err
	}

	batchServiceRole, err := iam.NewRole(ctx, "BatchExecutionRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Principal": {
						"Service": "batch.amazonaws.com"
					},
					"Effect": "Allow",
				}
			]
		}`),
	})
	if err != nil {
		return err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, "BatchExecutionRoleAttachment", &iam.RolePolicyAttachmentArgs{
		Role:      batchServiceRole.Name,
		PolicyArn: pulumi.String("arn:aws:iam::aws:policy/service-role/AWSBatchServiceRole"),
	})
	if err != nil {
		return err
	}

	instanceProfile, err := iam.NewInstanceProfile(ctx, "BatchInstanceProfile", &iam.InstanceProfileArgs{
		Role: ecsInstanceRole.Name,
	})
	if err != nil {
		return err
	}

	computeResourceOptions := &batch.ComputeEnvironmentComputeResourcesArgs{
		MinVcpus:         pulumi.Int(a.AwsConfig.BatchComputeEnvConfig.MinCpus),
		MaxVcpus:         pulumi.Int(a.AwsConfig.BatchComputeEnvConfig.MaxCpus),
		DesiredVcpus:     pulumi.Int(0),
		InstanceTypes:    pulumi.ToStringArray(a.AwsConfig.BatchComputeEnvConfig.InstanceTypes),
		Type:             pulumi.String("EC2"),
		Subnets:          subnets,
		SecurityGroupIds: pulumi.StringArray{a.BatchSecurityGroup.ID()},
		InstanceRole:     instanceProfile.Arn,
	}

	a.ComputeEnvironment, err = batch.NewComputeEnvironment(ctx, "compute-environment", &batch.ComputeEnvironmentArgs{
		ComputeEnvironmentName: pulumi.Sprintf("%s-compute-environment", a.StackId),
		ComputeResources:       computeResourceOptions,
		Type:                   pulumi.String("MANAGED"),
		ServiceRole:            batchServiceRole.Arn,
	})
	if err != nil {
		return fmt.Errorf("error creating compute environment: %w", err)
	}

	a.JobQueue, err = batch.NewJobQueue(ctx, "job-queue", &batch.JobQueueArgs{
		ComputeEnvironments: pulumi.StringArray{
			a.ComputeEnvironment.Arn,
		},
		State:    pulumi.String("ENABLED"),
		Priority: pulumi.Int(1),
		Tags:     pulumi.ToStringMap(tags.Tags(a.StackId, "job-queue", "job-queue")),
	})
	if err != nil {
		return err
	}

	return nil
}

// Docs: https://docs.aws.amazon.com/batch/latest/userguide/job_definition_parameters.html
type JobDefinitionContainerProperties struct {
	Image                string                `json:"image"`
	ResourceRequirements []ResourceRequirement `json:"resourceRequirements"`
	Command              []string              `json:"command"`
	JobRoleArn           string                `json:"jobRoleArn"`
	ExecutionRoleArn     string                `json:"executionRoleArn"`
	Environment          []EnvironmentVariable `json:"environment"`
}

func (p *NitricAwsPulumiProvider) Batch(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Batch, runtime provider.RuntimeProvider) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	// Tag the image
	repo, err := ecr.NewRepository(ctx, name, &ecr.RepositoryArgs{
		ForceDelete: pulumi.BoolPtr(true),
		Tags:        pulumi.ToStringMap(tags.Tags(p.StackId, name, "batch")),
	}, opts...)
	if err != nil {
		return err
	}

	wrappedImage, err := image.NewImage(ctx, name, &image.ImageArgs{
		SourceImage:   config.GetImage().GetUri(),
		RepositoryUrl: repo.RepositoryUrl,
		Runtime:       runtime(),
		RegistryArgs:  p.RegistryArgs,
	}, pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{repo}), pulumi.Provider(p.DockerProvider))
	if err != nil {
		return err
	}

	p.BatchRoles[name], err = iam.NewRole(ctx, "BatchJobRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Principal": {
						"Service": "ecs-tasks.amazonaws.com"
					},
					"Effect": "Allow",
				}
			]
		}`),
	}, opts...)
	if err != nil {
		return err
	}

	listActions := []string{
		// TODO: test that all resources still work without these permissions
		"sns:ListTopics",
		"sqs:ListQueues",
		"dynamodb:ListTables",
		"s3:ListAllMyBuckets",
		"tag:GetResources",
		"apigateway:GET",
	}

	// This is a tag key unique to this instance of the deployed stack.
	// Any resource with this unique tag will inherently be scoped to this stack.
	// This is used to scope the permissions of the lambda to only resources created by this stack.
	// stackScopedNameKey := tags.GetResourceNameKey(a.stackId)

	// Add resource list permissions
	// Currently the membrane will use list operations
	tmpJSON, err := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Action":   listActions,
				"Effect":   "Allow",
				"Resource": "*",
				// "Condition": map[string]map[string]string{
				// 	// Only apply this to resources who have a resource name key that matches this stack
				// 	"Null": {
				// 		fmt.Sprintf("aws:ResourceTag/%s", stackScopedNameKey): "false",
				// 	},
				// },
			},
		},
	})
	if err != nil {
		return err
	}

	_, err = iam.NewRolePolicy(ctx, name+"ListAccess", &iam.RolePolicyArgs{
		Role:   p.BatchRoles[name].ID(),
		Policy: pulumi.String(tmpJSON),
	}, opts...)
	// Deploy one job for each job that a batch handles
	// The job that it executes is defined by the job name provided in its env variables

	for _, job := range config.Jobs {
		jobName := job.GetName()

		dbEndpoint := pulumi.String("").ToStringOutput()
		if p.DatabaseCluster != nil {
			dbEndpoint = p.DatabaseCluster.Endpoint
		}

		dbPassword := pulumi.String("").ToStringOutput()
		if p.DbMasterPassword != nil {
			dbPassword = p.DbMasterPassword.Result
		}

		containerProperties := pulumi.All(wrappedImage.URI(), p.BatchRoles[name].Arn, dbEndpoint, dbPassword).ApplyT(func(args []interface{}) (string, error) {
			imageName := args[0].(string)
			jobRoleArn := args[1].(string)
			nitricDbEndpoint := args[2].(string)
			nitricDbPassword := args[3].(string)

			jobDefinitionContainerProperties := JobDefinitionContainerProperties{
				Image: imageName,
				ResourceRequirements: []ResourceRequirement{
					{
						Type:  "MEMORY",
						Value: fmt.Sprintf("%d", job.Requirements.Memory),
					},
					{
						Type:  "VCPU",
						Value: strconv.FormatFloat(float64(job.Requirements.Cpus), 'G', -1, 32),
					},
				},
				Environment: []EnvironmentVariable{
					{
						Name:  "NITRIC_JOB_NAME",
						Value: jobName,
					},
					{
						Name:  "MIN_WORKERS",
						Value: fmt.Sprintf("%d", len(config.Jobs)),
					},
					{
						Name:  "NITRIC_STACK_ID",
						Value: p.StackId,
					},
					{
						Name:  "AWS_REGION",
						Value: p.Region,
					},
				},
				JobRoleArn: jobRoleArn,
			}

			if nitricDbEndpoint != "" {
				jobDefinitionContainerProperties.Environment = append(jobDefinitionContainerProperties.Environment, EnvironmentVariable{
					Name:  "NITRIC_DATABASE_BASE_URL",
					Value: fmt.Sprintf("postgres://%s:%s@%s:%s", "nitric", nitricDbPassword, nitricDbEndpoint, "5432"),
				})
			}

			if job.Requirements.Gpus > 0 {
				jobDefinitionContainerProperties.ResourceRequirements = append(jobDefinitionContainerProperties.ResourceRequirements, ResourceRequirement{
					Type:  "GPU",
					Value: fmt.Sprintf("%d", job.Requirements.Gpus),
				})
			}

			containerPropertiesJson, err := json.Marshal(jobDefinitionContainerProperties)
			if err != nil {
				return "", err
			}

			return string(containerPropertiesJson), nil
		}).(pulumi.StringOutput)

		p.JobDefinitions[jobName], err = batch.NewJobDefinition(ctx, jobName, &batch.JobDefinitionArgs{
			Name:                pulumi.Sprintf("%s-job-%s", p.StackId, job.Name),
			ContainerProperties: containerProperties,
			Type:                pulumi.String("container"),
			Tags:                pulumi.ToStringMap(tags.Tags(p.StackId, jobName, "job")),
		}, opts...)
	}

	return err
}
