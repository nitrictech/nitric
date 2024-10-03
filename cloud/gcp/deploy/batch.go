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
	"fmt"
	"regexp"
	"strings"

	"cloud.google.com/go/batch/apiv1/batchpb"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/serviceaccount"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/storage"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"google.golang.org/protobuf/encoding/protojson"
)

var projectPermissions = map[string]string{
	"ar-reader":      "roles/artifactregistry.reader",
	"storage-viewer": "roles/storage.objectViewer",
	"batch-agent":    "roles/batch.agentReporter",
	"log-writer":     "roles/logging.logWriter",
}

func (p *NitricGcpPulumiProvider) Batch(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.Batch, runtimeProvider provider.RuntimeProvider) error {
	// Create a GCP batch task specification and push it to a GCP bucket for later use to run tasks
	defaultResourceOpts := []pulumi.ResourceOption{pulumi.Parent(parent), pulumi.Provider(p.DockerProvider)}

	// Deploy the image for the Batch instaces to GCR
	imageUriSplit := strings.Split(config.GetImage().GetUri(), "/")
	imageName := imageUriSplit[len(imageUriSplit)-1]

	image, err := image.NewImage(ctx, fmt.Sprintf("batch-image-%s", name), &image.ImageArgs{
		SourceImage:   config.GetImage().Uri,
		RepositoryUrl: pulumi.Sprintf("%s-docker.pkg.dev/%s/%s/%s", p.Region, p.GcpConfig.ProjectId, p.ContainerRegistry.Name, imageName),
		RegistryArgs:  p.RegistryArgs,
		Runtime:       runtimeProvider(),
	}, p.WithDefaultResourceOptions(defaultResourceOpts...)...)
	if err != nil {
		return err
	}

	invalidChars := regexp.MustCompile(`[^a-z0-9\-]`)
	gcpBatchName := invalidChars.ReplaceAllString(name, "-")

	p.BatchServiceAccounts[name], err = NewServiceAccount(ctx, gcpBatchName+"-batch-cloudrun-exec-acct", &GcpIamServiceAccountArgs{
		AccountId: gcpBatchName + "-exec",
	}, p.WithDefaultResourceOptions(defaultResourceOpts...)...)
	if err != nil {
		return err
	}

	// Apply base compute permissions required for nitric runtime to work
	_, err = projects.NewIAMMember(ctx, gcpBatchName+"-project-member", &projects.IAMMemberArgs{
		Project: pulumi.String(p.GcpConfig.ProjectId),
		Member:  pulumi.Sprintf("serviceAccount:%s", p.BatchServiceAccounts[name].ServiceAccount.Email),
		Role:    p.BaseComputeRole.Name,
	})
	if err != nil {
		return errors.WithMessage(err, "function project membership "+name)
	}

	// Apply additional project level permissions for the service account to interact with GCP batch
	for permName, role := range projectPermissions {
		// give the service account permission to pull images from GCR
		_, err = projects.NewIAMMember(ctx, fmt.Sprintf("%s-%s", gcpBatchName, permName), &projects.IAMMemberArgs{
			Project: pulumi.String(p.GcpConfig.ProjectId),
			Member:  pulumi.Sprintf("serviceAccount:%s", p.BatchServiceAccounts[name].ServiceAccount.Email),
			Role:    pulumi.String(role),
		}, p.WithDefaultResourceOptions(defaultResourceOpts...)...)
		if err != nil {
			return errors.WithMessage(err, "service account artifact registry membership "+name)
		}
	}

	// give the service account permission to act as itself so it may delegate delayed operations with its own permissions
	_, err = serviceaccount.NewIAMMember(ctx, gcpBatchName+"-acct-member", &serviceaccount.IAMMemberArgs{
		ServiceAccountId: p.BatchServiceAccounts[name].ServiceAccount.Name,
		Member:           pulumi.Sprintf("serviceAccount:%s", p.BatchServiceAccounts[name].ServiceAccount.Email),
		Role:             pulumi.String("roles/iam.serviceAccountUser"),
	}, p.WithDefaultResourceOptions(defaultResourceOpts...)...)
	if err != nil {
		return errors.WithMessage(err, "service account self membership "+name)
	}

	// for each job in the spec create a batch job protobuf specification
	for _, j := range config.Jobs {
		var accelerators []*batchpb.AllocationPolicy_Accelerator = nil
		if j.Requirements.Gpus > 0 {
			accelerators = []*batchpb.AllocationPolicy_Accelerator{
				{
					Type:  p.GcpConfig.GcpBatchCompute.AcceleratorType,
					Count: j.Requirements.Gpus,
				},
			}
		}

		containerOptions := []string{}
		if j.Requirements.Gpus > 0 {
			// TODO: Add support for additional accelerator types
			containerOptions = append(containerOptions, "--runtime=nvidia")
		}

		var dbUrl pulumi.StringOutput = pulumi.String("").ToStringOutput()
		if p.masterDb != nil {
			dbUrl = pulumi.Sprintf("postgresql://postgres:%s@%s:5432", p.dbMasterPassword.Result, p.masterDb.PrivateIpAddress)
		}

		var privateSubnet pulumi.StringOutput = pulumi.String("").ToStringOutput()
		if p.privateSubnet != nil {
			privateSubnet = p.privateSubnet.SelfLink
		}

		var privateNetwork pulumi.StringOutput = pulumi.String("").ToStringOutput()
		if p.privateNetwork != nil {
			privateNetwork = p.privateNetwork.SelfLink
		}

		jobDefinitionContents := pulumi.All(
			image.URI(), p.BatchServiceAccounts[name].ServiceAccount.Email, dbUrl, privateNetwork, privateSubnet, p.StackId,
		).ApplyT(func(args []interface{}) (string, error) {
			uri := args[0].(string)
			saEmail := args[1].(string)
			dbUrl := args[2].(string)
			privateNetwork := args[3].(string)
			privateSubnet := args[4].(string)
			stackId := args[5].(string)

			envVars := map[string]string{
				"NITRIC_JOB_NAME":       j.Name,
				"NITRIC_STACK_ID":       stackId,
				"GOOGLE_PROJECT_ID":     p.GcpConfig.ProjectId,
				"SERVICE_ACCOUNT_EMAIL": saEmail,
				"GCP_REGION":            p.Region,
				"MIN_WORKERS":           fmt.Sprintf("%d", len(config.Jobs)),
			}

			if dbUrl != "" {
				envVars["NITRIC_DATABASE_BASE_URL"] = dbUrl
			}

			var networkInterfaces *batchpb.AllocationPolicy_NetworkPolicy = nil
			if p.privateNetwork != nil && p.privateSubnet != nil {
				networkInterfaces = &batchpb.AllocationPolicy_NetworkPolicy{
					NetworkInterfaces: []*batchpb.AllocationPolicy_NetworkInterface{
						{
							Network:    privateNetwork,
							Subnetwork: privateSubnet,
						},
					},
				}
			}

			job := &batchpb.Job{
				TaskGroups: []*batchpb.TaskGroup{
					{
						TaskSpec: &batchpb.TaskSpec{
							Runnables: []*batchpb.Runnable{
								{
									Executable: &batchpb.Runnable_Container_{
										Container: &batchpb.Runnable_Container{
											// TODO: Configure the image uri
											ImageUri: uri,
											Options:  strings.Join(containerOptions, " "),
										},
									},
								},
							},
							Environment: &batchpb.Environment{
								Variables: envVars,
							},
							ComputeResource: &batchpb.ComputeResource{
								CpuMilli:  int64(j.Requirements.Cpus * 1000),
								MemoryMib: j.Requirements.Memory,
							},
						},
					},
				},
				AllocationPolicy: &batchpb.AllocationPolicy{
					ServiceAccount: &batchpb.ServiceAccount{
						Email: saEmail,
						Scopes: []string{
							"https://www.googleapis.com/auth/cloud-platform",
						},
					},
					Network: networkInterfaces,
					Instances: []*batchpb.AllocationPolicy_InstancePolicyOrTemplate{
						{
							PolicyTemplate: &batchpb.AllocationPolicy_InstancePolicyOrTemplate_Policy{
								Policy: &batchpb.AllocationPolicy_InstancePolicy{
									// MachineType:       machineType.name,
									// ProvisioningModel: batchpb.AllocationPolicy_STANDARD,
									Accelerators: accelerators,
								},
							},
						},
					},
				},
				LogsPolicy: &batchpb.LogsPolicy{
					Destination: batchpb.LogsPolicy_CLOUD_LOGGING,
				},
			}

			jobDefinitionJson, err := protojson.Marshal(job)
			if err != nil {
				return "", err
			}

			return string(jobDefinitionJson), nil
		}).(pulumi.StringOutput)

		// Store the job in the bucket for retrieval at runtime
		p.JobDefinitions[j.Name], err = storage.NewBucketObject(ctx, j.Name, &storage.BucketObjectArgs{
			Name:    pulumi.Sprintf("%s.json", j.Name),
			Bucket:  p.JobDefinitionBucket.Name,
			Content: jobDefinitionContents,
		}, pulumi.Parent(parent))
		if err != nil {
			return err
		}
	}

	return nil
}
