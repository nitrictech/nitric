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
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NitricCloudRunService - A wrapper that encapsulates all important information about a cloud run service deployed by nitric
type NitricCloudRunService struct {
	Name           string
	Service        *cloudrunv2.Service
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
		RepositoryUrl: pulumi.Sprintf("%s-docker.pkg.dev/%s/%s/%s", p.Region, p.GcpConfig.ProjectId, p.ContainerRegistry.Name, imageName),
		RegistryArgs:  p.RegistryArgs,
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

	env := cloudrunv2.ServiceTemplateContainerEnvArray{
		cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("NITRIC_ENVIRONMENT"),
			Value: pulumi.String("cloud"),
		},
		cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("MIN_WORKERS"),
			Value: pulumi.String(fmt.Sprintf("%d", config.Workers)),
		},
		cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("NITRIC_STACK_ID"),
			Value: pulumi.String(p.StackId),
		},
		cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("GOOGLE_PROJECT_ID"),
			Value: pulumi.String(p.GcpConfig.ProjectId),
		},
		cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("SERVICE_ACCOUNT_EMAIL"),
			Value: sa.ServiceAccount.Email,
		},
		cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("GCP_REGION"),
			Value: pulumi.String(p.Region),
		},
		cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("NITRIC_HTTP_PROXY_PORT"),
			Value: pulumi.String(fmt.Sprint(3000)),
		},
	}

	if p.JobDefinitionBucket != nil {
		env = append(env, cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("NITRIC_JOBS_BUCKET_NAME"),
			Value: p.JobDefinitionBucket.Name,
		})
	}

	env = append(env, cloudrunv2.ServiceTemplateContainerEnvArgs{
		Name:  pulumi.String("EVENT_TOKEN"),
		Value: res.EventToken,
	})

	if p.DelayQueue != nil {
		env = append(env, cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("DELAY_QUEUE_NAME"),
			Value: pulumi.Sprintf("projects/%s/locations/%s queues/%s", p.DelayQueue.Project, p.DelayQueue.Location, p.DelayQueue.Name),
		})
	}

	if p.masterDb != nil {
		env = append(env, cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String("NITRIC_DATABASE_BASE_URL"),
			Value: pulumi.Sprintf("postgresql://postgres:%s@%s:5432", p.dbMasterPassword.Result, p.masterDb.PrivateIpAddress),
		})
	}

	for k, v := range config.Env() {
		env = append(env, cloudrunv2.ServiceTemplateContainerEnvArgs{
			Name:  pulumi.String(k),
			Value: v,
		})
	}

	limits := map[string]string{
		"cpu":    fmt.Sprintf("%2f", unitConfig.CloudRun.Cpus),
		"memory": fmt.Sprintf("%dMi", unitConfig.CloudRun.Memory),
	}

	// Configuration will still break if the user specifies a 0 GPU count but their account does not support GPUs
	// only add if requested
	var nodeSelector cloudrunv2.ServiceTemplateNodeSelectorPtrInput = nil
	if unitConfig.CloudRun.Gpus > 0 {
		nodeSelector = &cloudrunv2.ServiceTemplateNodeSelectorArgs{
			Accelerator: pulumi.String("nvidia-l4"),
		}
		limits["nvidia.com/gpu"] = fmt.Sprintf("%d", unitConfig.CloudRun.Gpus)
	}

	serviceTemplate := cloudrunv2.ServiceTemplateArgs{
		ServiceAccount:                sa.ServiceAccount.Email,
		MaxInstanceRequestConcurrency: pulumi.Int(unitConfig.CloudRun.MaxInstances),
		Scaling: &cloudrunv2.ServiceTemplateScalingArgs{
			MinInstanceCount: pulumi.Int(unitConfig.CloudRun.MinInstances),
			MaxInstanceCount: pulumi.Int(unitConfig.CloudRun.MaxInstances),
		},
		Timeout: pulumi.Sprintf("%ds", unitConfig.CloudRun.Timeout),
		Containers: cloudrunv2.ServiceTemplateContainerArray{
			cloudrunv2.ServiceTemplateContainerArgs{
				Envs:  env,
				Image: image.URI(),
				Ports: cloudrunv2.ServiceTemplateContainerPortsArgs{
					ContainerPort: pulumi.Int(9001),
				},
				Resources: cloudrunv2.ServiceTemplateContainerResourcesArgs{
					Limits: pulumi.ToStringMap(limits),
				},
			},
		},
		NodeSelector: nodeSelector,
	}

	// Add vpc egress if there is a sql database
	if p.masterDb != nil {
		serviceTemplate.VpcAccess = &cloudrunv2.ServiceTemplateVpcAccessArgs{
			Connector: p.vpcConnector.SelfLink,
			Egress:    pulumi.String("PRIVATE_RANGES_ONLY"),
			// TODO: Re-enable when pulumi network interface support is fixed for tear down
			// NetworkInterfaces: &cloudrunv2.ServiceTemplateVpcAccessNetworkInterfaceArray{
			// 	&cloudrunv2.ServiceTemplateVpcAccessNetworkInterfaceArgs{
			// 		Network:    p.privateNetwork.ID(),
			// 		Subnetwork: p.privateSubnet.ID(),
			// 	},
			// },
		}

		dependsOn := []pulumi.Resource{p.privateNetwork, p.privateSubnet}
		for _, db := range p.DatabaseMigrationBuild {
			dependsOn = append(dependsOn, db)
		}

		opts = append(opts, pulumi.DependsOn(dependsOn))

		serviceTemplate.Annotations = pulumi.ToStringMapOutput(map[string]pulumi.StringOutput{"run.googleapis.com/cloudsql-instances": p.masterDb.ConnectionName})
	}

	migrationBuilds := lo.Values(p.DatabaseMigrationBuild)

	migrationBuildResources := []pulumi.Resource{}
	for _, migration := range migrationBuilds {
		migrationBuildResources = append(migrationBuildResources, migration)
	}

	res.Service, err = cloudrunv2.NewService(ctx, gcpServiceName, &cloudrunv2.ServiceArgs{
		Location: pulumi.String(p.Region),
		Project:  pulumi.String(p.GcpConfig.ProjectId),
		Template: serviceTemplate,
		Ingress:  pulumi.String("INGRESS_TRAFFIC_ALL"),
		Traffics: cloudrunv2.ServiceTrafficArray{
			&cloudrunv2.ServiceTrafficArgs{
				Percent: pulumi.Int(100),
				Type:    pulumi.String("TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"),
			},
		},
	}, p.WithDefaultResourceOptions(append([]pulumi.ResourceOption{pulumi.DependsOn(migrationBuildResources)}, opts...)...)...)
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

	_, err = cloudrunv2.NewServiceIamMember(ctx, gcpServiceName+"-invoker", &cloudrunv2.ServiceIamMemberArgs{
		Member:   pulumi.Sprintf("serviceAccount:%s", res.Invoker.Email),
		Role:     pulumi.String("roles/run.invoker"),
		Name:     res.Service.Name,
		Location: res.Service.Location,
	}, p.WithDefaultResourceOptions(opts...)...)
	if err != nil {
		return errors.WithMessage(err, "iam member "+name)
	}

	res.Url = res.Service.Uri

	p.CloudRunServices[name] = res

	return nil
}
