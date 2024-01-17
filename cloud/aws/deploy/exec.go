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
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	"github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/common/deploy/telemetry"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createEcrRepository(ctx *pulumi.Context, parent pulumi.Resource, stackId string, name string) (*ecr.Repository, error) {
	return ecr.NewRepository(ctx, name, &ecr.RepositoryArgs{
		ForceDelete: pulumi.BoolPtr(true),
		Tags:        pulumi.ToStringMap(tags.Tags(stackId, name, resources.ExecutionUnit)),
	}, pulumi.Parent(parent))
}

func createImage(ctx *pulumi.Context, parent pulumi.Resource, name string, authToken *ecr.GetAuthorizationTokenResult, repo *ecr.Repository, typeConfig *AwsConfigItem, config *deploymentspb.ExecutionUnit) (*image.Image, error) {
	if config.GetImage() == nil {
		return nil, fmt.Errorf("aws provider can only deploy execution with an image source")
	}

	if config.GetImage().GetUri() == "" {
		return nil, fmt.Errorf("aws provider can only deploy execution with an image source")
	}

	if config.Type == "" {
		config.Type = "default"
	}

	return image.NewImage(ctx, name, &image.ImageArgs{
		SourceImage:   config.GetImage().GetUri(),
		RepositoryUrl: repo.RepositoryUrl,
		Server:        pulumi.String(authToken.ProxyEndpoint),
		Username:      pulumi.String(authToken.UserName),
		Password:      pulumi.String(authToken.Password),
		Runtime:       runtime,
		Telemetry: &telemetry.TelemetryConfigArgs{
			TraceSampling: typeConfig.Telemetry,
			TraceName:     "awsxray",
			MetricName:    "awsemf",
			Extensions:    []string{},
		},
	}, pulumi.DependsOn([]pulumi.Resource{repo}))
}

func (a *NitricAwsPulumiProvider) ExecUnit(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.ExecutionUnit) error {

	opts := []pulumi.ResourceOption{pulumi.Parent(parent)}

	// Create the ECR repository to push the image to
	repo, err := createEcrRepository(ctx, parent, a.stackId, name)
	if err != nil {
		return err
	}

	if config.GetImage() == nil {
		return fmt.Errorf("aws provider can only deploy execution with an image source")
	}

	if config.GetImage().GetUri() == "" {
		return fmt.Errorf("aws provider can only deploy execution with an image source")
	}

	if config.Type == "" {
		config.Type = "default"
	}

	typeConfig, hasConfig := a.config.Config[config.Type]
	if !hasConfig {
		return fmt.Errorf("could not find config for type %s in %+v", config.Type, a.config)
	}

	image, err := createImage(ctx, parent, name, a.ecrAuthToken, repo, typeConfig, config)
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

	a.lambdaRoles[name], err = iam.NewRole(ctx, name+"LambdaRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(tmpJSON),
		Tags:             pulumi.ToStringMap(tags.Tags(a.stackId, name+"LambdaRole", resources.ExecutionUnit)),
	}, opts...)
	if err != nil {
		return err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, name+"LambdaBasicExecution", &iam.RolePolicyAttachmentArgs{
		PolicyArn: iam.ManagedPolicyAWSLambdaBasicExecutionRole,
		Role:      a.lambdaRoles[name].ID(),
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
		"sns:ListTopics",
		"sqs:ListQueues",
		"dynamodb:ListTables",
		"s3:ListAllMyBuckets",
		"tag:GetResources",
		"apigateway:GET",
	}

	// Add resource list permissions
	// Currently the membrane will use list operations
	tmpJSON, err = json.Marshal(map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Action":   append(listActions, telemetryActions...),
				"Effect":   "Allow",
				"Resource": "*",
			},
		},
	})
	if err != nil {
		return err
	}

	// TODO: Lock this SNS topics for which this function has pub definitions
	// FIXME: Limit to known resources
	_, err = iam.NewRolePolicy(ctx, name+"ListAccess", &iam.RolePolicyArgs{
		Role:   a.lambdaRoles[name].ID(),
		Policy: pulumi.String(tmpJSON),
	}, opts...)
	if err != nil {
		return err
	}

	// allow lambda to execute step function

	envVars := pulumi.StringMap{
		"NITRIC_ENVIRONMENT":     pulumi.String("cloud"),
		"NITRIC_STACK_ID":        pulumi.String(a.stackId),
		"MIN_WORKERS":            pulumi.String(fmt.Sprint(config.Workers)),
		"NITRIC_HTTP_PROXY_PORT": pulumi.String(fmt.Sprint(3000)),
	}
	for k, v := range config.Env {
		envVars[k] = pulumi.String(v)
	}

	var vpcConfig *awslambda.FunctionVpcConfigArgs = nil
	if typeConfig.Lambda.Vpc != nil {
		vpcConfig = &awslambda.FunctionVpcConfigArgs{
			SubnetIds:        pulumi.ToStringArray(typeConfig.Lambda.Vpc.SubnetIds),
			SecurityGroupIds: pulumi.ToStringArray(typeConfig.Lambda.Vpc.SecurityGroupIds),
		}

		// Create a policy attachment for VPC access
		_, err = iam.NewRolePolicyAttachment(ctx, name+"VPCAccessExecutionRole", &iam.RolePolicyAttachmentArgs{
			PolicyArn: iam.ManagedPolicyAWSLambdaVPCAccessExecutionRole,
			Role:      a.lambdaRoles[name].ID(),
		}, opts...)
		if err != nil {
			return err
		}
	}

	a.lambdas[name], err = awslambda.NewFunction(ctx, name, &awslambda.FunctionArgs{
		ImageUri:    image.URI(),
		MemorySize:  pulumi.IntPtr(typeConfig.Lambda.Memory),
		Timeout:     pulumi.IntPtr(typeConfig.Lambda.Timeout),
		PackageType: pulumi.String("Image"),
		Role:        a.lambdaRoles[name].Arn,
		Tags:        pulumi.ToStringMap(tags.Tags(a.stackId, name, resources.ExecutionUnit)),
		VpcConfig:   vpcConfig,
		Environment: awslambda.FunctionEnvironmentArgs{Variables: envVars},
	}, opts...)
	if err != nil {
		return err
	}

	if typeConfig.Lambda.ProvisionedConcurreny > 0 {
		_, err := awslambda.NewProvisionedConcurrencyConfig(ctx, name, &awslambda.ProvisionedConcurrencyConfigArgs{
			FunctionName:                    a.lambdas[name].Arn,
			ProvisionedConcurrentExecutions: pulumi.Int(typeConfig.Lambda.ProvisionedConcurreny),
			Qualifier:                       a.lambdas[name].Name,
		})
		if err != nil {
			return err
		}
	}

	// ensure that the lambda was deployed successfully
	_ = a.lambdas[name].Arn.ApplyT(func(arn string) (bool, error) {
		payload, _ := json.Marshal(map[string]interface{}{
			"x-nitric-healthcheck": true,
		})

		err := retry.Do(func() error {
			_, err := a.lambdaClient.Invoke(&lambda.InvokeInput{
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