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
	"encoding/json"
	"fmt"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
)

func createEcrRepository(ctx *pulumi.Context, parent pulumi.Resource, stackId string, name string) (*ecr.Repository, error) {
	return ecr.NewRepository(ctx, name, &ecr.RepositoryArgs{
		ForceDelete: pulumi.BoolPtr(true),
		Tags:        pulumi.ToStringMap(tags.Tags(stackId, name, resources.Service)),
	}, pulumi.Parent(parent))
}

func (a *NitricAwsPulumiProvider) Service(ctx *pulumi.Context, parent pulumi.Resource, name string, config *pulumix.NitricPulumiServiceConfig, runtime provider.RuntimeProvider) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	// Create the ECR repository to push the image to
	repo, err := createEcrRepository(ctx, parent, a.StackId, name)
	if err != nil {
		return err
	}

	if config.GetImage() == nil {
		return fmt.Errorf("aws provider can only deploy service with an image source")
	}

	if config.GetImage().GetUri() == "" {
		return fmt.Errorf("aws provider can only deploy service with an image source")
	}

	if config.Type == "" {
		config.Type = "default"
	}

	typeConfig, hasConfig := a.AwsConfig.Config[config.Type]
	if !hasConfig {
		return fmt.Errorf("could not find config for type %s in %+v", config.Type, a.AwsConfig)
	}

	image, err := image.NewImage(ctx, name, &image.ImageArgs{
		SourceImage:   config.GetImage().GetUri(),
		RepositoryUrl: repo.RepositoryUrl,
		Runtime:       runtime(),
	}, pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{repo}), pulumi.Provider(a.DockerProvider))
	// image, err := createImage(ctx, parent, name, repo, config, a.DockerProvider, runtime)
	if err != nil {
		return err
	}

	opts = append(opts, pulumi.Parent(parent))

	tmpJSON, err := json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Sid":    "",
				"Effect": "Allow",
				"Principal": map[string]interface{}{
					"Service": "lambda.amazonaws.com",
				},
				"Action": "sts:AssumeRole",
			},
		},
	})
	if err != nil {
		return err
	}

	a.LambdaRoles[name], err = iam.NewRole(ctx, name+"LambdaRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(tmpJSON),
		Tags:             pulumi.ToStringMap(tags.Tags(a.StackId, name+"LambdaRole", resources.Service)),
	}, opts...)
	if err != nil {
		return err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, name+"LambdaBasicExecution", &iam.RolePolicyAttachmentArgs{
		PolicyArn: iam.ManagedPolicyAWSLambdaBasicExecutionRole,
		Role:      a.LambdaRoles[name].ID(),
	}, opts...)
	if err != nil {
		return err
	}

	telemetryActions := []string{
		"xray:PutTraceSegments",
		"xray:PutTelemetryRecords",
		"xray:GetSamplingRules",
		"xray:GetSamplingTargets",
		"xray:GetSamplingStatisticSummaries",
		"ssm:GetParameters",
		"logs:CreateLogStream",
		"logs:PutLogEvents",
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
	// Currently the nitric server will use list operations
	tmpJSON, err = json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Action":   append(listActions, telemetryActions...),
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
		Role:   a.LambdaRoles[name].ID(),
		Policy: pulumi.String(tmpJSON),
	}, opts...)
	if err != nil {
		return err
	}

	// allow lambda to execute step function

	envVars := pulumi.StringMap{
		"NITRIC_ENVIRONMENT":     pulumi.String("cloud"),
		"NITRIC_STACK_ID":        pulumi.String(a.StackId),
		"MIN_WORKERS":            pulumi.String(fmt.Sprint(config.Workers)),
		"NITRIC_HTTP_PROXY_PORT": pulumi.String(fmt.Sprint(3000)),
	}
	for k, v := range config.Env() {
		envVars[k] = v
	}

	if a.JobQueue != nil {
		envVars["NITRIC_JOB_QUEUE_ARN"] = a.JobQueue.Arn
	}

	if a.DatabaseCluster != nil {
		// Include the base database cluster URI for the runtime to resolve databases based on their name
		envVars["NITRIC_DATABASE_BASE_URL"] = pulumi.Sprintf("postgres://%s:%s@%s:%s", "nitric", a.DbMasterPassword.Result,
			a.DatabaseCluster.Endpoint, "5432")

		// Include database migrations to ensure a pulumi dependency is created for the lambda
		//	the migrations need to complete before the lambda is deployed
		sqlDatabasesMigrated := lo.Map(lo.Values(a.SqlDatabases), func(item *RdsDatabase, idx int) interface{} {
			return item.Migrated
		})

		databasesMigrated := pulumi.All(sqlDatabasesMigrated...).ApplyT(func(migrated []interface{}) bool {
			for _, m := range migrated {
				if b, ok := m.(bool); !ok || !b {
					return false
				}
			}

			return true
		})

		envVars["DATABASES_MIGRATED"] = pulumi.Sprintf("%t", databasesMigrated)
	}

	var vpcConfig *awslambda.FunctionVpcConfigArgs = nil
	if typeConfig.Lambda.Vpc != nil {
		if a.Vpc != nil {
			return fmt.Errorf("this stack has configuration that requires its own VPC configuration and cannot accept custom lambda VPC configurations")
		}

		vpcConfig = &awslambda.FunctionVpcConfigArgs{
			SubnetIds:        pulumi.ToStringArray(typeConfig.Lambda.Vpc.SubnetIds),
			SecurityGroupIds: pulumi.ToStringArray(typeConfig.Lambda.Vpc.SecurityGroupIds),
		}
	} else if a.Vpc != nil {
		vpcConfig = &awslambda.FunctionVpcConfigArgs{
			SubnetIds:        a.Vpc.PrivateSubnetIds,
			SecurityGroupIds: pulumi.StringArray{a.RdsSecurityGroup.ID()},
			// VpcId:            a.Vpc.VpcId,
		}
	}

	if vpcConfig != nil {
		// Create a policy attachment for VPC access
		_, err = iam.NewRolePolicyAttachment(ctx, name+"VPCAccessExecutionRole", &iam.RolePolicyAttachmentArgs{
			PolicyArn: iam.ManagedPolicyAWSLambdaVPCAccessExecutionRole,
			Role:      a.LambdaRoles[name].ID(),
		}, opts...)
		if err != nil {
			return err
		}
	}

	dependsOn := []pulumi.ResourceOption{pulumi.DependsOn([]pulumi.Resource{image})}
	// Add Sql database migration dependencies
	for _, db := range a.SqlDatabases {
		dependsOn = append(dependsOn, pulumi.DependsOn([]pulumi.Resource{db}))
	}

	a.Lambdas[name], err = awslambda.NewFunction(ctx, name, &awslambda.FunctionArgs{
		// Use repository to generate the URI, instead of the image, using the image results in errors when the same project is torn down and redeployed.
		// This appears to be because the local image ends up with multiple repositories and the wrong one is selected.
		// XXX: Reverted change for the above comment as lambda image deployments were not rolling forward (under tag latest)
		// causing intermittent deployment and runtime failures
		ImageUri:         image.URI(),
		MemorySize:       pulumi.Int(typeConfig.Lambda.Memory),
		Timeout:          pulumi.Int(typeConfig.Lambda.Timeout),
		EphemeralStorage: awslambda.FunctionEphemeralStorageArgs{Size: pulumi.Int(typeConfig.Lambda.EphemeralStorage)},
		Publish:          pulumi.Bool(true),
		PackageType:      pulumi.String("Image"),
		Role:             a.LambdaRoles[name].Arn,
		Tags:             pulumi.ToStringMap(tags.Tags(a.StackId, name, resources.Service)),
		VpcConfig:        vpcConfig,
		Environment:      awslambda.FunctionEnvironmentArgs{Variables: envVars},
		// since we only rely on the repository to determine the ImageUri, the image must be added as a dependency to avoid a race.
	}, append(dependsOn, opts...)...)
	if err != nil {
		return err
	}

	if typeConfig.Lambda.ProvisionedConcurreny > 0 {
		_, err = awslambda.NewProvisionedConcurrencyConfig(ctx, name, &awslambda.ProvisionedConcurrencyConfigArgs{
			FunctionName:                    a.Lambdas[name].Arn,
			ProvisionedConcurrentExecutions: pulumi.Int(typeConfig.Lambda.ProvisionedConcurreny),
			Qualifier:                       a.Lambdas[name].Version,
		}, pulumi.DependsOn([]pulumi.Resource{a.Lambdas[name]}))
		if err != nil {
			return err
		}
	}

	// ensure that the lambda was deployed successfully
	_ = a.Lambdas[name].Arn.ApplyT(func(arn string) (bool, error) {
		payload, _ := json.Marshal(map[string]interface{}{
			"x-nitric-healthcheck": true,
		})

		err := retry.Do(func() error {
			_, err := a.LambdaClient.Invoke(&lambda.InvokeInput{
				FunctionName: aws.String(arn),
				Payload:      payload,
			})

			return err
		}, retry.Attempts(3))
		if err != nil {
			return false, err
		}

		return true, nil
	})

	return nil
}
