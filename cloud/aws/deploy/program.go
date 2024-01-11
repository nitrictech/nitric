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
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/nitrictech/nitric/cloud/aws/deploy/api"
	"github.com/nitrictech/nitric/cloud/aws/deploy/bucket"
	"github.com/nitrictech/nitric/cloud/aws/deploy/collection"
	"github.com/nitrictech/nitric/cloud/aws/deploy/config"
	"github.com/nitrictech/nitric/cloud/aws/deploy/exec"
	"github.com/nitrictech/nitric/cloud/aws/deploy/policy"
	"github.com/nitrictech/nitric/cloud/aws/deploy/schedule"
	"github.com/nitrictech/nitric/cloud/aws/deploy/secret"
	"github.com/nitrictech/nitric/cloud/aws/deploy/stack"
	"github.com/nitrictech/nitric/cloud/aws/deploy/topic"
	"github.com/nitrictech/nitric/cloud/aws/deploy/websocket"
	commonDeploy "github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/resources"
	common "github.com/nitrictech/nitric/cloud/common/deploy/tags"
	"github.com/nitrictech/nitric/cloud/common/deploy/telemetry"
	deploy "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/proto/resources/v1"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ecr"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func NewUpProgram(ctx context.Context, details *commonDeploy.CommonStackDetails, config *config.AwsConfig, spec *deploy.Spec) (auto.Stack, error) {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	lambdaClient := lambda.New(sess, &aws.Config{Region: aws.String(details.Region)})
	resourceTaggingClient := resourcegroupstaggingapi.New(sess)

	return auto.UpsertStackInlineSource(context.TODO(), details.FullStackName, details.Project, func(ctx *pulumi.Context) error {
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

		stackIdChan := make(chan string)
		pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.Result).ApplyT(func(id string) string {
			stackIdChan <- id
			return id
		})

		stackID := <-stackIdChan
		_, err = stack.NewAwsResourceGroup(ctx, details.FullStackName, &stack.AwsResourceGroupArgs{
			StackID: stackID,
		})

		if err != nil {
			return err
		}

		// Deploy all secrets
		secrets := map[string]*secret.SecretsManagerSecret{}
		for _, res := range spec.Resources {
			switch c := res.Config.(type) {
			case *deploy.Resource_Secret:
				importArn := ""

				if config.Import.Secrets != nil {
					importArn = config.Import.Secrets[res.Name]
				}

				secrets[res.Name], err = secret.NewSecretsManagerSecret(ctx, res.Name, &secret.SecretsManagerSecretArgs{
					StackID: stackID,
					Secret:  c.Secret,
					Import:  importArn,
					Client:  resourceTaggingClient,
				})
				if err != nil {
					return err
				}
			}
		}

		// Deploy all collections
		collections := map[string]*collection.DynamodbCollection{}
		for _, res := range spec.Resources {
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

		authToken, err := ecr.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenArgs{})
		if err != nil {
			return err
		}

		// Deploy all execution units
		execs := map[string]*exec.LambdaExecUnit{}
		execPrincipals := map[string]*iam.Role{}
		for _, res := range spec.Resources {
			switch eu := res.Config.(type) {
			case *deploy.Resource_ExecutionUnit:
				repo, err := ecr.NewRepository(ctx, res.Name, &ecr.RepositoryArgs{
					ForceDelete: pulumi.BoolPtr(true),
					Tags:        pulumi.ToStringMap(common.Tags(stackID, res.Name, resources.ExecutionUnit)),
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

				if eu.ExecutionUnit.Type == "" {
					eu.ExecutionUnit.Type = "default"
				}

				typeConfig, hasConfig := config.Config[eu.ExecutionUnit.Type]
				if !hasConfig {
					return fmt.Errorf("could not find config for type %s in %+v", eu.ExecutionUnit.Type, config.Config)
				}

				image, err := image.NewImage(ctx, res.Name, &image.ImageArgs{
					SourceImage:   eu.ExecutionUnit.GetImage().GetUri(),
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
				if err != nil {
					return err
				}

				if typeConfig.Lambda != nil {
					execs[res.Name], err = exec.NewLambdaExecutionUnit(ctx, res.Name, &exec.LambdaExecUnitArgs{
						DockerImage: image,
						StackID:     stackID,
						Compute:     eu.ExecutionUnit,
						EnvMap:      pulumi.ToStringMap(eu.ExecutionUnit.Env),
						Client:      lambdaClient,
						Config:      *typeConfig.Lambda,
					})
					execPrincipals[res.Name] = execs[res.Name].Role
				} else {
					return fmt.Errorf("no target execution unit specified for %s", res.Name)
				}

				if err != nil {
					return err
				}
			}
		}
		principals[v1.ResourceType_Function] = execPrincipals

		// Deploy all buckets
		buckets := map[string]*bucket.S3Bucket{}
		for _, res := range spec.Resources {
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

				if len(b.Bucket.Notifications) > 0 {
					_, err = bucket.NewS3Notification(ctx, fmt.Sprintf("notification-%s", res.Name), &bucket.S3NotificationArgs{
						StackID:      stackID,
						Location:     details.Region,
						Bucket:       buckets[res.Name],
						Functions:    execs,
						Notification: b.Bucket.Notifications,
					})
					if err != nil {
						return err
					}
				}
			}
		}

		// deploy API Gateways
		// gws := map[string]
		for _, res := range spec.Resources {
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

				config, _ := config.Apis[res.Name]

				_, err = api.NewAwsApiGateway(ctx, res.Name, &api.AwsApiGatewayArgs{
					LambdaFunctions: execs,
					StackID:         stackID,
					OpenAPISpec:     doc,
					Config:          config,
				})
				if err != nil {
					return err
				}
			}
		}

		// Add all HTTP proxies
		httpProxies := map[string]*api.AwsHttpProxy{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Http:
				fun := execs[t.Http.Target.GetExecutionUnit()]

				httpProxies[res.Name], err = api.NewAwsHttpProxy(ctx, res.Name, &api.AwsHttpProxyArgs{
					StackID:        stackID,
					LambdaFunction: fun,
				})
				if err != nil {
					return err
				}
			}
		}

		// deploy websockets
		websockets := map[string]*websocket.AwsWebsocketApiGateway{}
		for _, res := range spec.Resources {
			switch ws := res.Config.(type) {
			case *deploy.Resource_Websocket:
				websockets[res.Name], err = websocket.NewAwsWebsocketApiGateway(ctx, res.Name, &websocket.AwsWebsocketApiGatewayArgs{
					DefaultTarget:    execs[ws.Websocket.MessageTarget.GetExecutionUnit()],
					ConnectTarget:    execs[ws.Websocket.ConnectTarget.GetExecutionUnit()],
					DisconnectTarget: execs[ws.Websocket.DisconnectTarget.GetExecutionUnit()],
					StackID:          stackID,
				})
				if err != nil {
					return err
				}
			}
		}

		// Deploy all schedules
		schedules := map[string]*schedule.AwsEventbridgeSchedule{}
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Schedule:
				// get the target of the schedule

				execUnitName := t.Schedule.Target.GetExecutionUnit()
				execUnit, ok := execs[execUnitName]
				if !ok {
					return fmt.Errorf("no execution unit with name %s", execUnitName)
				}

				// Create schedule targeting a given lambda
				schedules[res.Name], err = schedule.NewAwsEventbridgeSchedule(ctx, res.Name, &schedule.AwsEventbridgeScheduleArgs{
					Exec:     execUnit,
					Schedule: t.Schedule,
					Tz:       config.ScheduleTimezone,
				})
				if err != nil {
					return err
				}
			}
		}

		// Deploy all topics
		topics := map[string]*topic.SNSTopic{}
		for _, res := range spec.Resources {
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
		for _, res := range spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Policy:
				_, err = policy.NewIAMPolicy(ctx, res.Name, &policy.PolicyArgs{
					Policy: t.Policy,
					Resources: &policy.StackResources{
						Buckets:     buckets,
						Topics:      topics,
						Collections: collections,
						Secrets:     secrets,
						Websockets:  websockets,
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
}
