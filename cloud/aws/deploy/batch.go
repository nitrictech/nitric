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

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
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
		Server:        pulumi.String(p.EcrAuthToken.ProxyEndpoint),
		Username:      pulumi.String(p.EcrAuthToken.UserName),
		Password:      pulumi.String(p.EcrAuthToken.Password),
		Runtime:       runtime(),
	}, pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{repo}))
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

	listActions := []string{
		// TODO: test that all resources still work without these permissions
		"sns:ListTopics",
		"sqs:ListQueues",
		"dynamodb:ListTables",
		"s3:ListAllMyBuckets",
		"tag:GetResources",
		"apigateway:GET",
		// Allow batch job submission
		// TODO: Limit this to batches available within the nitric stack
		"batch:SubmitJob",
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
				"Action":   append(listActions),
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

	// Create a new Iam Role for the job

	// Deploy one job for each job that a batch handles
	// The job that it executes is defined by the job name provided in its env variables

	for _, job := range config.Jobs {
		// get the stack config for the job
		jobConfig, ok := p.AwsConfig.Jobs[job]
		if !ok {
			jobConfig = p.AwsConfig.Jobs["default"]
		}

		jobName := job
		containerProperties := pulumi.All(wrappedImage.URI(), p.BatchRoles[name].Arn).ApplyT(func(args []interface{}) (string, error) {
			imageName := args[0].(string)
			jobRoleArn := args[1].(string)

			jobDefinitionContainerProperties := JobDefinitionContainerProperties{
				Image: imageName,
				// Command:          []string{""},
				ResourceRequirements: []ResourceRequirement{
					// TODO: Make these configurable options
					// Or template parameters that can be set at runtime
					{
						Type:  "MEMORY",
						Value: fmt.Sprintf("%d", jobConfig.Memory),
					},
					{
						Type:  "VCPU",
						Value: fmt.Sprintf("%d", jobConfig.Cpus),
					},
				},
				Environment: []EnvironmentVariable{
					{
						Name:  "NITRIC_JOB_NAME",
						Value: jobName,
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
				// ExecutionRoleArn: batchRoleArn,
			}

			if jobConfig.Gpus > 0 {
				jobDefinitionContainerProperties.ResourceRequirements = append(jobDefinitionContainerProperties.ResourceRequirements, ResourceRequirement{
					Type:  "GPU",
					Value: fmt.Sprintf("%d", jobConfig.Gpus),
				})
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
