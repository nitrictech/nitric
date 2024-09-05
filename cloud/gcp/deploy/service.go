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
	"fmt"
	"regexp"
	"strings"

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumix"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/cloudrun"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NitricCloudRunService - A wrapper that encapsulates all important information about a cloud run service deployed by nitric
type NitricCloudRunService struct {
	Name           string
	Service        *cloudrun.Service
	ServiceAccount *serviceaccount.Account
	Url            pulumi.StringInput
	Invoker        *serviceaccount.Account
	EventToken     pulumi.StringOutput
}

func (p *NitricGcpPulumiProvider) Service(ctx *pulumi.Context, parent pulumi.Resource, name string, config *pulumix.NitricPulumiServiceConfig, runtime provider.RuntimeProvider) error {
	opts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.Provider(p.DockerProvider)}

	res := &NitricCloudRunService{
		Name: name,
	}

	invalidChars := regexp.MustCompile(`[^a-z0-9\-]`)
	gcpServiceName := invalidChars.ReplaceAllString(name, "-")

	if config.GetImage() == nil || config.GetImage().Uri == "" {
		return fmt.Errorf("gcp provider can only deploy service with an image source")
	}

	if config.Type == "" {
		config.Type = "default"
	}

	// get config for service
	unitConfig, hasConfig := p.GcpConfig.Config[config.Type]
	if !hasConfig {
		return status.Errorf(codes.InvalidArgument, "unable to find config %s in stack config %+v", config.Type, p.GcpConfig.Config)
	}

	if unitConfig.CloudRun == nil {
		return status.Errorf(codes.InvalidArgument, "unable to find cloud run config in stack config %+v", p.GcpConfig.Config)
	}

	// Get the image name:tag from the uri
	imageUriSplit := strings.Split(config.GetImage().GetUri(), "/")
	imageName := imageUriSplit[len(imageUriSplit)-1]

	image, err := image.NewImage(ctx, gcpServiceName, &image.ImageArgs{
		SourceImage:   config.GetImage().Uri,
		RepositoryUrl: pulumi.Sprintf("gcr.io/%s/%s", p.GcpConfig.ProjectId, imageName),
		Runtime:       runtime(),
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return err
	}

	sa, err := NewServiceAccount(ctx, gcpServiceName+"-cloudrun-exec-acct", &GcpIamServiceAccountArgs{
		AccountId: gcpServiceName + "-exec",
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return err
	}

	res.ServiceAccount = sa.ServiceAccount

	// generate a token for internal application events to authenticate themeselves
	// https://cloud.google.com/appengine/docs/flexible/writing-and-responding-to-pub-sub-messages?tab=go#top
	token, err := random.NewRandomPassword(ctx, gcpServiceName+"-event-token", &random.RandomPasswordArgs{
		Special: pulumi.Bool(false),
		Length:  pulumi.Int(32),
		Keepers: pulumi.ToMap(map[string]interface{}{
			"name": gcpServiceName,
		}),
	})
	if err != nil {
		return errors.WithMessage(err, "service event token")
	}

	res.EventToken = token.Result

	_, err = projects.NewIAMMember(ctx, gcpServiceName+"-project-member", &projects.IAMMemberArgs{
		Project: pulumi.String(p.GcpConfig.ProjectId),
		Member:  pulumi.Sprintf("serviceAccount:%s", sa.ServiceAccount.Email),
		Role:    p.BaseComputeRole.Name,
	})
	if err != nil {
		return errors.WithMessage(err, "function project membership "+name)
	}

	// give the service account permission to use itself
	_, err = serviceaccount.NewIAMMember(ctx, gcpServiceName+"-acct-member", &serviceaccount.IAMMemberArgs{
		ServiceAccountId: sa.ServiceAccount.Name,
		Member:           pulumi.Sprintf("serviceAccount:%s", sa.ServiceAccount.Email),
		Role:             pulumi.String("roles/iam.serviceAccountUser"),
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return errors.WithMessage(err, "service account self membership "+name)
	}

	env := cloudrun.ServiceTemplateSpecContainerEnvArray{
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_ENVIRONMENT"),
			Value: pulumi.String("cloud"),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("MIN_WORKERS"),
			Value: pulumi.String(fmt.Sprintf("%d", config.Workers)),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_STACK_ID"),
			Value: pulumi.String(p.StackId),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("SERVICE_ACCOUNT_EMAIL"),
			Value: sa.ServiceAccount.Email,
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("GCP_REGION"),
			Value: pulumi.String(p.Region),
		},
		cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("NITRIC_HTTP_PROXY_PORT"),
			Value: pulumi.String(fmt.Sprint(3000)),
		},
	}

	env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
		Name:  pulumi.String("EVENT_TOKEN"),
		Value: res.EventToken,
	})

	if p.DelayQueue != nil {
		env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String("DELAY_QUEUE_NAME"),
			Value: pulumi.Sprintf("projects/%s/locations/%s/queues/%s", p.DelayQueue.Project, p.DelayQueue.Location, p.DelayQueue.Name),
		})
	}

	for k, v := range config.Env() {
		env = append(env, cloudrun.ServiceTemplateSpecContainerEnvArgs{
			Name:  pulumi.String(k),
			Value: v,
		})
	}

	res.Service, err = cloudrun.NewService(ctx, gcpServiceName, &cloudrun.ServiceArgs{
		AutogenerateRevisionName: pulumi.BoolPtr(true),
		Location:                 pulumi.String(p.Region),
		Project:                  pulumi.String(p.GcpConfig.ProjectId),
		Template: cloudrun.ServiceTemplateArgs{
			Metadata: cloudrun.ServiceTemplateMetadataArgs{
				Annotations: pulumi.StringMap{
					"autoscaling.knative.dev/minScale": pulumi.Sprintf("%d", unitConfig.CloudRun.MinInstances),
					"autoscaling.knative.dev/maxScale": pulumi.Sprintf("%d", unitConfig.CloudRun.MaxInstances),
				},
			},
			Spec: cloudrun.ServiceTemplateSpecArgs{
				ServiceAccountName:   sa.ServiceAccount.Email,
				ContainerConcurrency: pulumi.Int(unitConfig.CloudRun.Concurrency),
				TimeoutSeconds:       pulumi.Int(unitConfig.CloudRun.Timeout),
				Containers: cloudrun.ServiceTemplateSpecContainerArray{
					cloudrun.ServiceTemplateSpecContainerArgs{
						Envs:  env,
						Image: image.URI(),
						Ports: cloudrun.ServiceTemplateSpecContainerPortArray{
							cloudrun.ServiceTemplateSpecContainerPortArgs{
								ContainerPort: pulumi.Int(9001),
							},
						},
						Resources: cloudrun.ServiceTemplateSpecContainerResourcesArgs{
							Limits: pulumi.StringMap{
								"cpu":    pulumi.Sprintf("%2f", unitConfig.CloudRun.Cpus),
								"memory": pulumi.Sprintf("%dMi", unitConfig.CloudRun.Memory),
							},
						},
					},
				},
			},
		},
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return errors.WithMessage(err, "cloud run "+name)
	}

	svcAcct, err := NewServiceAccount(ctx, gcpServiceName+"-cloudrun-invoker", &GcpIamServiceAccountArgs{
		AccountId: gcpServiceName,
	})
	if err != nil {
		return errors.WithMessage(err, "invokerAccount "+name)
	}

	res.Invoker = svcAcct.ServiceAccount

	_, err = cloudrun.NewIamMember(ctx, gcpServiceName+"-invoker", &cloudrun.IamMemberArgs{
		Member:   pulumi.Sprintf("serviceAccount:%s", res.Invoker.Email),
		Role:     pulumi.String("roles/run.invoker"),
		Service:  res.Service.Name,
		Location: res.Service.Location,
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return errors.WithMessage(err, "iam member "+name)
	}

	res.Url = res.Service.Statuses.ApplyT(func(ss []cloudrun.ServiceStatus) (string, error) {
		if len(ss) == 0 {
			return "", errors.New("serviceStatus is empty")
		}

		return *ss[0].Url, nil
	}).(pulumi.StringInput)

	p.CloudRunServices[name] = res

	return nil
}
