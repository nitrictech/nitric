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

package exec

import (
	"encoding/json"
	"fmt"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/lambda/lambdaiface"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	awslambda "github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type LambdaExecUnitArgs struct {
	Client  lambdaiface.LambdaAPI
	StackID pulumi.StringInput
	// Image needs to be built and uploaded first
	DockerImage *image.Image
	Compute     *v1.ExecutionUnit
	EnvMap      map[string]string
}

type LambdaExecUnit struct {
	pulumi.ResourceState

	Name     string
	Function *awslambda.Function
	Role     *iam.Role
}

func NewLambdaExecutionUnit(ctx *pulumi.Context, name string, args *LambdaExecUnitArgs, opts ...pulumi.ResourceOption) (*LambdaExecUnit, error) {
	res := &LambdaExecUnit{Name: name}

	err := ctx.RegisterComponentResource("nitric:exec:AWSLambda", name, res, opts...)
	if err != nil {
		return nil, err
	}

	opts = append(opts, pulumi.Parent(res))

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
		return nil, err
	}

	res.Role, err = iam.NewRole(ctx, name+"LambdaRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(tmpJSON),
		Tags:             common.Tags(ctx, args.StackID, name+"LambdaRole"),
	}, opts...)
	if err != nil {
		return nil, err
	}

	_, err = iam.NewRolePolicyAttachment(ctx, name+"LambdaBasicExecution", &iam.RolePolicyAttachmentArgs{
		PolicyArn: iam.ManagedPolicyAWSLambdaBasicExecutionRole,
		Role:      res.Role.ID(),
	}, opts...)
	if err != nil {
		return nil, err
	}

	telemeteryActions := []string{
		"xray:PutTraceSegments",
		"xray:PutTelemetryRecords",
		"xray:GetSamplingRules",
		"xray:GetSamplingTargets",
		"xray:GetSamplingStatisticSummaries",
		"ssm:GetParameters",
		//		"logs:CreateLogGroup",
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
				"Action":   append(listActions, telemeteryActions...),
				"Effect":   "Allow",
				"Resource": "*",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	// TODO: Lock this SNS topics for which this function has pub definitions
	// FIXME: Limit to known resources
	_, err = iam.NewRolePolicy(ctx, name+"ListAccess", &iam.RolePolicyArgs{
		Role:   res.Role.ID(),
		Policy: pulumi.String(tmpJSON),
	}, opts...)
	if err != nil {
		return nil, err
	}

	// allow lambda to execute step function

	envVars := pulumi.StringMap{
		"NITRIC_ENVIRONMENT": pulumi.String("cloud"),
		"NITRIC_STACK":       args.StackID,
		"MIN_WORKERS":        pulumi.String(fmt.Sprint(args.Compute.Workers)),
	}
	for k, v := range args.EnvMap {
		envVars[k] = pulumi.String(v)
	}

	// Set defaults if not already provided
	if args.Compute.Memory == 0 {
		args.Compute.Memory = 128
	}

	if args.Compute.Timeout == 0 {
		args.Compute.Timeout = 10
	}

	res.Function, err = awslambda.NewFunction(ctx, name, &awslambda.FunctionArgs{
		ImageUri:    args.DockerImage.URI(),
		MemorySize:  pulumi.IntPtr(int(args.Compute.Memory)),
		Timeout:     pulumi.IntPtr(int(args.Compute.Timeout)),
		PackageType: pulumi.String("Image"),
		Role:        res.Role.Arn,
		Tags:        common.Tags(ctx, args.StackID, name),
		Environment: awslambda.FunctionEnvironmentArgs{Variables: envVars},
	}, opts...)
	if err != nil {
		return nil, err
	}

	// ensure that the lambda was deployed successfully
	isHealthy := res.Function.Arn.ApplyT(func(arn string) (bool, error) {
		payload, _ := json.Marshal(map[string]interface{}{
			"x-nitric-healthcheck": true,
		})

		err := retry.Do(func() error {
			_, err := args.Client.Invoke(&lambda.InvokeInput{
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

	// Deploy subscription separately
	// for _, t := range args.Compute.Unit().Triggers.Topics {
	// 	topic, ok := args.Topics[t]
	// 	if ok {
	// 		_, err = awslambda.NewPermission(ctx, name+t+"Permission", &awslambda.PermissionArgs{
	// 			SourceArn: topic.Sns.Arn,
	// 			Function:  res.Function.Name,
	// 			Principal: pulumi.String("sns.amazonaws.com"),
	// 			Action:    pulumi.String("lambda:InvokeFunction"),
	// 		}, opts...)
	// 		if err != nil {
	// 			return nil, err
	// 		}

	// 		_, err = sns.NewTopicSubscription(ctx, name+t+"Subscription", &sns.TopicSubscriptionArgs{
	// 			Endpoint: res.Function.Arn,
	// 			Protocol: pulumi.String("lambda"),
	// 			Topic:    topic.Sns.ID(),
	// 		}, opts...)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 	} else {
	// 		fmt.Printf("WARNING: Function %s has a Trigger %s, but the topic is missing", name, t)
	// 	}
	// }

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":    pulumi.String(res.Name),
		"lambda":  res.Function,
		"healthy": isHealthy,
	})
}
