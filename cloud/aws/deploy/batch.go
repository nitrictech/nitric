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

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/batch"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
)

type ResourceRequirement struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type EnvironmentVariable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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

func (p *NitricAwsPulumiProvider) Batch(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Batch) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	// Tag the image
	repo, err := ecr.NewRepository(ctx, name, &ecr.RepositoryArgs{
		ForceDelete: pulumi.BoolPtr(true),
		Tags:        pulumi.ToStringMap(tags.Tags(p.StackId, name, "batch")),
	}, opts...)
	if err != nil {
		return err
	}

	inspect, err := image.CommandFromImageInspect(config.GetImage().Uri, " ")
	if err != nil {
		return err
	}

	newTag, err := docker.NewTag(ctx, name+"-tag", &docker.TagArgs{
		SourceImage: pulumi.String(inspect.ID),
		TargetImage: repo.RepositoryUrl,
	}, opts...)
	if err != nil {
		return err
	}

	image, err := docker.NewRegistryImage(ctx, name+"-remote", &docker.RegistryImageArgs{
		Name: repo.RepositoryUrl,
		Triggers: pulumi.Map{
			"imageSha": pulumi.String(inspect.ID),
		},
	}, pulumi.Parent(parent), pulumi.Provider(p.DockerProvider), pulumi.DependsOn([]pulumi.Resource{newTag}))
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
					"Sid": ""
				}
			]
		}`),
	}, opts...)
	if err != nil {
		return err
	}

	// create a job role for the task definition
	// jobRole, err := iam.NewRole(ctx, name+"-job-role", &iam.RoleArgs{})

	// Create a new Iam Role for the job

	// Deploy one job for each job that a batch handles
	// The job that it executes is defined by the job name provided in its env variables

	for _, job := range config.Jobs {
		jobName := job
		containerProperties := pulumi.All(image.Name, p.BatchExecutionRole.Arn, p.BatchRoles[name].Arn).ApplyT(func(args []interface{}) (string, error) {
			imageName := args[0].(string)
			batchRoleArn := args[1].(string)
			jobRoleArn := args[2].(string)

			jobDefinitionContainerProperties := JobDefinitionContainerProperties{
				Image: imageName,
				// Command:          []string{""},
				ResourceRequirements: []ResourceRequirement{
					// TODO: Make these configurable options
					// Or template parameters that can be set at runtime
					{
						Type:  "MEMORY",
						Value: "512",
					},
					{
						Type:  "VCPU",
						Value: "0.25",
					},
				},
				Environment: []EnvironmentVariable{
					{
						Name:  "NITRIC_JOB_NAME",
						Value: jobName,
					},
				},
				JobRoleArn:       jobRoleArn,
				ExecutionRoleArn: batchRoleArn,
			}

			containerPropertiesJson, err := json.Marshal(jobDefinitionContainerProperties)
			if err != nil {
				return "", err
			}

			return string(containerPropertiesJson), nil
		}).(pulumi.StringOutput)

		_, err = batch.NewJobDefinition(ctx, name, &batch.JobDefinitionArgs{
			Name:                pulumi.Sprintf("%s-job-%s", p.StackId, job),
			ContainerProperties: containerProperties,

			// TODO: Set tags for job definition discovery
			Type: pulumi.String("container"),
			Tags: pulumi.ToStringMap(tags.Tags(p.StackId, job, "job")),
		}, opts...)
	}

	return err
}
