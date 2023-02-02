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
	"fmt"

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	v1 "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudtasks"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/pubsub"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudRunner struct {
	pulumi.ResourceState

	Name    string
	Service *cloudrun.Service
	Url     pulumi.StringInput
}

type CloudRunnerArgs struct {
	Location        pulumi.StringInput
	ProjectId       string
	Compute         *v1.ExecutionUnit
	Image           *image.Image
	EnvMap          map[string]string
	Topics          map[string]*pubsub.Topic
	DelayQueue      *cloudtasks.Queue
	BaseComputeRole *projects.IAMCustomRole
	ServiceAccount  *serviceaccount.Account

	StackID pulumi.StringInput
}

var defaultConcurrency = 300

func GetPerms() []string {
	return []string{
		"storage.buckets.list",
		"storage.buckets.get",
		"cloudtasks.queues.get",
		"cloudtasks.tasks.create",
		"cloudtrace.traces.patch",
		"monitoring.timeSeries.create",
		// permission for blob signing
		// this is safe as only permissions this account has are delegated
		"iam.serviceAccounts.signBlob",
	}
}

func NewCloudRunner(ctx *pulumi.Context, name string, args *CloudRunnerArgs, opts ...pulumi.ResourceOption) (*CloudRunner, error) {
	res := &CloudRunner{
		Name: name,
	}

	err := ctx.RegisterComponentResource("nitric:exec:GCPCloudRunner", name, res, opts...)
	if err != nil {
		return nil, err
	}

	// apply basic project level permissions for nitric resource discovery
	_, err = projects.NewIAMMember(ctx, res.Name+"-project-member", &projects.IAMMemberArgs{
		Project: pulumi.String(args.ProjectId),
		Member:  pulumi.Sprintf("serviceAccount:%s", args.ServiceAccount.Email),
		Role:    args.BaseComputeRole.Name,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "function project membership "+res.Name)
	}

	// give the service account permission to use itself
	_, err = serviceaccount.NewIAMMember(ctx, res.Name+"-acct-member", &serviceaccount.IAMMemberArgs{
		ServiceAccountId: args.ServiceAccount.Name,
		Member:           pulumi.Sprintf("serviceAccount:%s", args.ServiceAccount.Email),
		Role:             pulumi.String("roles/iam.serviceAccountUser"),
	})
	if err != nil {
		return nil, errors.WithMessage(err, "service account self membership "+res.Name)
	}

	env := getCloudRunnerEnvs(args)

	if args.DelayQueue != nil {
		env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("DELAY_QUEUE_NAME"),
			Value: pulumi.Sprintf("projects/%s/locations/%s/queues/%s", args.DelayQueue.Project, args.DelayQueue.Location, args.DelayQueue.Name),
		})
	}

	for k, v := range args.EnvMap {
		env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String(k),
			Value: pulumi.String(v),
		})
	}

	// Deploy the func
	maxScale := 10
	minScale := 0

	res.Service, err = cloudrun.NewService(ctx, name, &cloudrun.ServiceArgs{
		AutogenerateRevisionName: pulumi.BoolPtr(true),
		Location:                 args.Location,
		Project:                  pulumi.String(args.ProjectId),
		Template: cloudrun.ServiceTemplateArgs{
			Metadata: cloudrun.ServiceTemplateMetadataArgs{
				Annotations: pulumi.StringMap{
					"autoscaling.knative.dev/minScale": pulumi.Sprintf("%d", minScale),
					"autoscaling.knative.dev/maxScale": pulumi.Sprintf("%d", maxScale),
				},
			},
			Spec: cloudrun.ServiceTemplateSpecArgs{
				ServiceAccountName:   args.ServiceAccount.Email,
				ContainerConcurrency: pulumi.Int(defaultConcurrency),
				Containers: cloudrun.ServiceTemplateSpecContainerArray{
					cloudrun.ServiceTemplateSpecContainerArgs{
						Envs:  env,
						Image: args.Image.URI(),
						Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
							cloudrun.ServiceTemplateSpecContainerPortArgs{
								ContainerPort: pulumi.Int(9001),
							},
						},
						Resources: cloudrun.ServiceTemplateSpecContainerResourcesArgs{
							// Limits: pulumi.StringMap{"memory": pulumi.Sprintf("%dMi", args.Compute.Unit().Memory)},
						},
					},
				},
			},
		},
	}, append(opts, pulumi.Parent(res))...)
	if err != nil {
		return nil, errors.WithMessage(err, "iam member "+name)
	}

	res.Url = res.Service.Statuses.ApplyT(func(ss []cloudrun.ServiceStatus) (string, error) {
		if len(ss) == 0 {
			return "", errors.New("serviceStatus is empty")
		}

		return *ss[0].Url, nil
	}).(pulumi.StringInput)

	// wire up its subscriptions
	// if len(args.Topics) > 0 {
	// 	// Create an account for invoking this func via subscriptions
	// 	// TODO: Do we want to make this one account for subscription in future
	// 	// TODO: We will likely configure this via eventarc in the future
	// 	invokerAccount, err := serviceaccount.NewAccount(ctx, name+"subacct", &serviceaccount.AccountArgs{
	// 		// accountId accepts a max of 30 chars, limit our generated name to this length
	// 		AccountId: pulumi.String(utils.StringTrunc(name, 30-8) + "subacct"),
	// 	}, append(opts, pulumi.Parent(res))...)
	// 	if err != nil {
	// 		return nil, errors.WithMessage(err, "invokerAccount "+name)
	// 	}

	// 	// Apply permissions for the above account to the newly deployed cloud run service
	// 	_, err = cloudrun.NewIamMember(ctx, name+"-subrole", &cloudrun.IamMemberArgs{
	// 		Member:   pulumi.Sprintf("serviceAccount:%s", invokerAccount.Email),
	// 		Role:     pulumi.String("roles/run.invoker"),
	// 		Service:  res.Service.Name,
	// 		Location: res.Service.Location,
	// 	}, append(opts, pulumi.Parent(res))...)
	// 	if err != nil {
	// 		return nil, errors.WithMessage(err, "iam member "+name)
	// 	}

	// 	for _, t := range args.Topics {
	// 		t.Name.ApplyT(func (topicName string) (*events.PubSubSubscription, error) {
	// 			s, err := events.NewPubSubSubscription(ctx, name, &events.PubSubSubscriptionArgs{
	// 				Topic: topicName,
	// 				Url: res.Url,
	// 				ServiceAccount: invokerAccount.Email,
	// 				Function: res.Service,
	// 			})
	// 			if err != nil {
	// 				return nil, err
	// 			}

	// 			return s, nil
	// 		})
	// 	}
	// }

	return res, ctx.RegisterResourceOutputs(res, pulumi.Map{
		"name":    pulumi.String(res.Name),
		"service": res.Service,
		"url":     res.Url,
	})
}

func getCloudRunnerEnvs(args *CloudRunnerArgs) cloudrun.ServiceTemplateSpecContainerEnvArray {
	return cloudrun.ServiceTemplateSpecContainerEnvArray{
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_ENVIRONMENT"),
			Value: pulumi.String("cloud"),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("MIN_WORKERS"),
			Value: pulumi.String(fmt.Sprintf("%d", args.Compute.Workers)),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_STACK"),
			Value: args.StackID,
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("SERVICE_ACCOUNT_EMAIL"),
			Value: args.ServiceAccount.Email,
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("GCP_REGION"),
			Value: args.Location,
		},
	}
}
