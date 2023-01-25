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
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/nitrictech/nitric/cloud/aws/deploy/api"
	"github.com/nitrictech/nitric/cloud/aws/deploy/bucket"
	"github.com/nitrictech/nitric/cloud/aws/deploy/collection"
	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	"github.com/nitrictech/nitric/cloud/aws/deploy/policy"
	"github.com/nitrictech/nitric/cloud/aws/deploy/queue"
	"github.com/nitrictech/nitric/cloud/aws/deploy/schedule"
	"github.com/nitrictech/nitric/cloud/aws/deploy/stack"
	"github.com/nitrictech/nitric/cloud/aws/deploy/topic"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpStreamMessageWriter struct {
	stream deploy.DeployService_UpServer
}

func (s *UpStreamMessageWriter) Write(bytes []byte) (int, error) {
	err := s.stream.Send(&deploy.DeployUpEvent{
		Content: &deploy.DeployUpEvent_Message{
			Message: &deploy.DeployEventMessage{
				Message: string(bytes),
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}

// Up - Deploy requested infrastructure for a stack
func (d *DeployServer) Up(request *deploy.DeployUpRequest, stream deploy.DeployService_UpServer) error {
	details, err := getStackDetailsFromAttributes(request.Attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	lambdaClient := lambda.New(sess, &aws.Config{Region: aws.String(details.Region)})

	pulumiStack, err := auto.UpsertStackInlineSource(context.TODO(), details.Stack, details.Project, func(ctx *pulumi.Context) error {
		principals := map[v1.ResourceType]map[string]*iam.Role{}

		// Calculate unique stackID
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
		stackID := pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.ID())

		_, err = stack.NewAwsResourceGroup(ctx, details.Stack, &stack.AwsResourceGroupArgs{
			StackID: stackID,
		})
		if err != nil {
			return err
		}

		// Deploy all buckets
		buckets := map[string]*bucket.S3Bucket{}
		for _, res := range request.Spec.Resources {
			switch b := res.Config.(type) {
			case *deploy.Resource_Bucket:
				buckets[res.Name], err = bucket.NewS3Bucket(ctx, res.Name, &bucket.S3BucketArgs{
					// TODO: Calculate stack ID
					StackID: stackID,
					Bucket:  b.Bucket,
				})
				if err != nil {
					return err
				}
			}
		}

		// Deploy all collections
		collections := map[string]*collection.DynamodbCollection{}
		for _, res := range request.Spec.Resources {
			switch c := res.Config.(type) {
			case *deploy.Resource_Collection:
				collections[res.Name], err = collection.NewDynamodbCollection(ctx, res.Name, &collection.DynamodbCollectionArgs{
					StackID:    stackID,
					Collection: c.Collection,
				})
				if err != nil {
					return err
				}
			}
		}

		// Deploy all queues
		queues := map[string]*queue.SQSQueue{}
		for _, res := range request.Spec.Resources {
			switch q := res.Config.(type) {
			case *deploy.Resource_Queue:
				queues[res.Name], err = queue.NewSQSQueue(ctx, res.Name, &queue.SQSQueueArgs{
					// TODO: Calculate stack ID
					StackID: stackID,
					Queue:   q.Queue,
				})
				if err != nil {
					return err
				}
			}
		}

		authToken, err := ecr.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenArgs{})
		if err != nil {
			return err
		}

		// Deploy all execution units
		execs := map[string]*exec.LambdaExecUnit{}
		execPrincipals := map[string]*iam.Role{}
		for _, res := range request.Spec.Resources {
			switch eu := res.Config.(type) {
			case *deploy.Resource_ExecutionUnit:
				repo, err := ecr.NewRepository(ctx, res.Name, &ecr.RepositoryArgs{
					ForceDelete: pulumi.BoolPtr(true),
					Tags:        common.Tags(ctx, stackID, res.Name),
				})
				if err != nil {
					return err
				}

				if eu.ExecutionUnit.GetImage() == nil {
					return fmt.Errorf("aws provider can only deploy execution with an image source")
				}

				if eu.ExecutionUnit.GetImage().GetUri() == "" {
					return fmt.Errorf("aws provider can only deploy execution with an image source")
				}

				fmt.Println(eu.ExecutionUnit.GetImage().GetUri())

				image, err := image.NewImage(ctx, res.Name, &image.ImageArgs{
					SourceImage:   eu.ExecutionUnit.GetImage().GetUri(),
					RepositoryUrl: repo.RepositoryUrl,
					Server:        pulumi.String(authToken.ProxyEndpoint),
					Username:      pulumi.String(authToken.UserName),
					Password:      pulumi.String(authToken.Password),
					Runtime:       runtime,
				}, pulumi.DependsOn([]pulumi.Resource{repo}))
				if err != nil {
					return err
				}

				// Create execution unit
				execs[res.Name], err = exec.NewLambdaExecutionUnit(ctx, res.Name, &exec.LambdaExecUnitArgs{
					DockerImage: image,
					StackID:     stackID,
					Compute:     eu.ExecutionUnit,
					EnvMap:      map[string]string{},
					Client:      lambdaClient,
				})
				execPrincipals[res.Name] = execs[res.Name].Role

				if err != nil {
					return err
				}
			}
		}
		principals[v1.ResourceType_Function] = execPrincipals

		// deploy API Gateways
		// gws := map[string]
		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Api:
				// Deserialize the OpenAPI document

				if t.Api.GetOpenapi() == "" {
					return fmt.Errorf("aws provider can only deploy OpenAPI specs")
				}

				doc := &openapi3.T{}
				err := doc.UnmarshalJSON([]byte(t.Api.GetOpenapi()))
				if err != nil {
					return fmt.Errorf("invalid document suppled for api: %s", res.Name)
				}

				_, err = api.NewAwsApiGateway(ctx, res.Name, &api.AwsApiGatewayArgs{
					LambdaFunctions: execs,
					StackID:         stackID,
					OpenAPISpec:     doc,
				})
				if err != nil {
					return err
				}
			}
		}

		// Deploy all schedules
		schedules := map[string]*schedule.AwsCloudwatchSchedule{}
		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Schedule:
				// get the target of the schedule

				execUnitName := t.Schedule.Target.GetExecutionUnit()
				execUnit, ok := execs[execUnitName]
				if !ok {
					return fmt.Errorf("no execution unit with name %s", execUnitName)
				}

				// Create schedule targeting a given lambda
				schedules[res.Name], err = schedule.NewAwsCloudwatchSchedule(ctx, res.Name, &schedule.AwsCloudwatchScheduleArgs{
					StackID: stackID,
					Exec:    execUnit,
					Cron:    t.Schedule.Cron,
				})
				if err != nil {
					return err
				}
			}
		}

		// Deploy all topics
		topics := map[string]*topic.SNSTopic{}
		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Topic:
				// Create topics
				topics[res.Name], err = topic.NewSNSTopic(ctx, res.Name, &topic.SNSTopicArgs{
					StackID: stackID,
					Topic:   t.Topic,
				})
				if err != nil {
					return err
				}

				// Create subscriptions for the topic
				for _, sub := range t.Topic.Subscriptions {
					subName := fmt.Sprintf("%s-%s-sub", sub.GetExecutionUnit(), res.Name)
					// Get the deployed execution unit
					unit, ok := execs[sub.GetExecutionUnit()]
					if !ok {
						return fmt.Errorf("invalid execution unit %s given for topic subscription", sub.GetExecutionUnit())
					}

					_, err = topic.NewSNSTopicSubscription(ctx, subName, &topic.SNSTopicSubscriptionArgs{
						Lambda: unit,
						Topic:  topics[res.Name],
					})
					if err != nil {
						return err
					}
				}
			}
		}

		// Create policies
		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Policy:
				_, err = policy.NewIAMPolicy(ctx, res.Name, &policy.PolicyArgs{
					Policy: t.Policy,
					Resources: &policy.StackResources{
						Buckets: buckets,
						Topics:  topics,
						Queues:  queues,
					},
					Principals: principals,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	_ = pulumiStack.SetConfig(context.TODO(), "aws:region", auto.ConfigValue{Value: details.Region})

	messageWriter := &UpStreamMessageWriter{
		stream: stream,
	}

	// Run the program
	_, err = pulumiStack.Up(context.TODO(), optup.ProgressStreams(messageWriter))
	if err != nil {
		return err
	}

	return nil
}
