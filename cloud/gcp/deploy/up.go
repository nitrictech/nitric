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
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/utils"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/bucket"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/events"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/exec"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/gateway"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/policy"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/queue"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/secret"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudtasks"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
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

	pulumiStack, err := auto.UpsertStackInlineSource(context.TODO(), details.Stack, details.Project, func(ctx *pulumi.Context) error {
		project, err := organizations.LookupProject(ctx, &organizations.LookupProjectArgs{
			ProjectId: &details.ProjectId,
		}, nil)
		if err != nil {
			return err
		}

		nitricProj, err := NewProject(ctx, "project", &ProjectArgs{
			ProjectId:     details.ProjectId,
			ProjectNumber: project.Number,
		})
		if err != nil {
			return err
		}

		defaultResourceOptions := pulumi.DependsOn([]pulumi.Resource{nitricProj})

		stackRandId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-stack-name", ctx.Stack()), &random.RandomStringArgs{
			Special: pulumi.Bool(false),
			Length:  pulumi.Int(8),
			Upper:   pulumi.Bool(false),
			Keepers: pulumi.ToMap(map[string]interface{}{
				"stack-name": ctx.Stack(),
			}),
		}, defaultResourceOptions)
		if err != nil {
			return err
		}

		stackID := pulumi.Sprintf("%s-%s", ctx.Stack(), stackRandId.ID())

		// Deploy all buckets
		buckets := map[string]*bucket.CloudStorageBucket{}
		for _, res := range request.Spec.Resources {
			switch b := res.Config.(type) {
			case *deploy.Resource_Bucket:
				buckets[res.Name], err = bucket.NewCloudStorageBucket(ctx, res.Name, &bucket.CloudStorageBucketArgs{
					StackID:  stackID,
					Bucket:   b.Bucket,
					Location: details.Region,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}
			}
		}

		// Deploy all queues
		queues := map[string]*queue.PubSubTopic{}
		for _, res := range request.Spec.Resources {
			switch q := res.Config.(type) {
			case *deploy.Resource_Queue:
				queues[res.Name], err = queue.NewPubSubTopic(ctx, res.Name, &queue.PubSubTopicArgs{
					StackID:  stackID,
					Queue:    q.Queue,
					Location: details.Region,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}
			}
		}

		// create a shared queue for enabling delayed messaging
		// cloud run functions will create OIDC tokens for their own service accounts
		// to apply to push actions to pubsub, so their scope should still be limited to that
		topicDelayQueue, err := cloudtasks.NewQueue(ctx, "delay-queue", &cloudtasks.QueueArgs{
			Location: pulumi.String(details.Region),
		}, defaultResourceOptions)
		if err != nil {
			return err
		}

		// Deploy all execution units
		authToken, err := getGCPToken(ctx)
		if err != nil {
			return err
		}

		execs := map[string]*exec.CloudRunner{}

		baseCustomRoleId, err := random.NewRandomString(ctx, fmt.Sprintf("%s-base-role", details.Stack), &random.RandomStringArgs{
			Special: pulumi.Bool(false),
			Length:  pulumi.Int(8),
			Keepers: pulumi.ToMap(map[string]interface{}{
				"stack-name": details.Stack,
			}),
		})
		if err != nil {
			return errors.WithMessage(err, "base customRole id")
		}

		// Telemetry permissions
		// for _, fc := range g.sc.Config {
		// 	if fc.Telemetry != nil && *fc.Telemetry > 0 {
		// 		perms = append(perms, []string{
		// 			"monitoring.metricDescriptors.create",
		// 			"monitoring.metricDescriptors.get",
		// 			"monitoring.metricDescriptors.list",
		// 			"monitoring.monitoredResourceDescriptors.get",
		// 			"monitoring.monitoredResourceDescriptors.list",
		// 			"monitoring.timeSeries.create",
		// 		}...)

		// 		break
		// 	}
		// }

		principalMap := make(policy.PrincipalMap)
		principalMap[v1.ResourceType_Function] = make(map[string]*serviceaccount.Account)

		// setup a basic IAM role for general access and resource discovery
		baseComputeRole, err := projects.NewIAMCustomRole(ctx, "base-role", &projects.IAMCustomRoleArgs{
			Title:       pulumi.String(details.Stack + "-functions-base-role"),
			Permissions: pulumi.ToStringArray(exec.GetPerms()),
			RoleId:      baseCustomRoleId.ID(),
		}, defaultResourceOptions)
		if err != nil {
			return errors.WithMessage(err, "base customRole")
		}

		for _, res := range request.Spec.Resources {
			switch eu := res.Config.(type) {
			case *deploy.Resource_ExecutionUnit:
				if eu.ExecutionUnit.GetImage() == nil {
					return fmt.Errorf("gcp provider can only deploy execution with an image source")
				}

				if eu.ExecutionUnit.GetImage().GetUri() == "" {
					return fmt.Errorf("gcp provider can only deploy execution with an image source")
				}

				// Get the image name:tag from the uri
				imageUriSplit := strings.Split(eu.ExecutionUnit.GetImage().GetUri(), "/")
				imageName := imageUriSplit[len(imageUriSplit)-1]

				image, err := image.NewImage(ctx, res.Name, &image.ImageArgs{
					SourceImage:   eu.ExecutionUnit.GetImage().GetUri(),
					RepositoryUrl: pulumi.Sprintf("gcr.io/%s/%s", details.ProjectId, imageName),
					Username:      pulumi.String("oauth2accesstoken"),
					Password:      pulumi.String(authToken.AccessToken),
					Server:        pulumi.String("https://gcr.io"),
					Runtime:       runtime,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				// Create a service account for this cloud run instance
				sa, err := serviceaccount.NewAccount(ctx, res.Name+"acct", &serviceaccount.AccountArgs{
					AccountId: pulumi.String(utils.StringTrunc(res.Name, 30-5) + "-acct"),
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				execs[res.Name], err = exec.NewCloudRunner(ctx, res.Name, &exec.CloudRunnerArgs{
					Location:        pulumi.String(details.Region),
					ProjectId:       details.ProjectId,
					Topics:          map[string]*pubsub.Topic{},
					Compute:         res.GetExecutionUnit(),
					Image:           image,
					EnvMap:          map[string]string{},
					DelayQueue:      topicDelayQueue,
					ServiceAccount:  sa,
					BaseComputeRole: baseComputeRole,
					StackID:         stackID,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				principalMap[v1.ResourceType_Function][res.Name] = sa
			}
		}

		apis := map[string]*gateway.ApiGateway{}
		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Api:
				if t.Api.GetOpenapi() == "" {
					return fmt.Errorf("gcp provider can only deploy OpenAPI specs")
				}

				doc := &openapi3.T{}
				err := doc.UnmarshalJSON([]byte(t.Api.GetOpenapi()))
				if err != nil {
					return fmt.Errorf("invalid document suppled for api: %s", res.Name)
				}

				apis[res.Name], err = gateway.NewApiGateway(ctx, res.Name, &gateway.ApiGatewayArgs{
					StackID:     stackID,
					ProjectId:   details.ProjectId,
					Functions:   execs,
					OpenAPISpec: doc,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}
			}
		}

		secrets := map[string]*secret.SecretManagerSecret{}
		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Secret:
				secrets[res.Name], err = secret.NewSecretManagerSecret(ctx, res.Name, &secret.SecretManagerSecretArgs{
					StackID: stackID,
					Secret:  t.Secret,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}
			}
		}

		topics := map[string]*events.PubSubTopic{}
		for _, res := range request.Spec.Resources {
			switch t := res.Config.(type) {
			case *deploy.Resource_Topic:
				topics[res.Name], err = events.NewPubSubTopic(ctx, res.Name, &events.PubSubTopicArgs{
					Topic:     t.Topic,
					Location:  details.Region,
					ProjectId: details.ProjectId,
					StackID:   stackID,
				}, defaultResourceOptions)
				if err != nil {
					return err
				}

				for _, sub := range t.Topic.Subscriptions {
					subName := fmt.Sprintf("%s-%s-sub", sub.GetExecutionUnit(), res.Name)

					// Get the deployed execution unit
					unit, ok := execs[sub.GetExecutionUnit()]
					if !ok {
						return fmt.Errorf("invalid execution unit %s given for topic subscription", sub.GetExecutionUnit())
					}

					_, err = events.NewPubSubSubscription(ctx, subName, &events.PubSubSubscriptionArgs{
						Topic:    res.Name,
						Function: unit,
					}, defaultResourceOptions)
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
				_, err := policy.NewIAMPolicy(ctx, res.Name, &policy.PolicyArgs{
					Policy: t.Policy,
					Resources: &policy.StackResources{
						Buckets: buckets,
						Topics:  topics,
						Queues:  queues,
						Secrets: secrets,
					},
					Principals: principalMap,
					ProjectID:  pulumi.String(details.ProjectId),
					StackID:    stackID,
				})
				if err != nil {
					return err
				}
			}
		}

		return nil
	})

	err = pulumiStack.SetConfig(context.TODO(), "gcp:region", auto.ConfigValue{Value: details.Region})
	if err != nil {
		return err
	}

	err = pulumiStack.SetConfig(context.TODO(), "gcp:project", auto.ConfigValue{Value: details.ProjectId})
	if err != nil {
		return err
	}

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

func getGCPToken(ctx *pulumi.Context) (*oauth2.Token, error) {
	// If the user is attempting to impersonate a gcp service account using pulumi using the GOOGLE_IMPERSONATE_SERVICE_ACCOUNT env var
	// Read more: (https://www.pulumi.com/registry/packages/gcp/installation-configuration/#configuration-reference)
	targetSA := os.Getenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT")

	var token *oauth2.Token

	if targetSA != "" {
		service, err := iamcredentials.NewService(ctx.Context())
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("Unable to impersonate service account: %s", targetSA))
		}

		accessToken, err := service.Projects.ServiceAccounts.GenerateAccessToken(fmt.Sprintf("projects/-/serviceAccounts/%s", targetSA), &iamcredentials.GenerateAccessTokenRequest{
			Scope: []string{
				"https://www.googleapis.com/auth/cloud-platform",
				"https://www.googleapis.com/auth/trace.append",
			},
		}).Do()
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("Unable to impersonate service account: %s", targetSA))
		}

		if accessToken == nil {
			return nil, fmt.Errorf("Unable to impersonate service account.")
		}

		token = &oauth2.Token{AccessToken: token.AccessToken}
	}

	if token == nil { // for unit testing
		creds, err := google.FindDefaultCredentialsWithParams(ctx.Context(), google.CredentialsParams{
			Scopes: []string{
				"https://www.googleapis.com/auth/cloud-platform",
				"https://www.googleapis.com/auth/trace.append",
			},
		})
		if err != nil {
			return nil, errors.WithMessage(err, "Unable to find credentials, try 'gcloud auth application-default login'")
		}

		token, err = creds.TokenSource.Token()
		if err != nil {
			return nil, errors.WithMessage(err, "Unable to acquire token source")
		}
	}

	return token, nil
}
