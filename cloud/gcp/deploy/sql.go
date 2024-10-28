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
	"strings"

	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/cloudrunv2"
	"github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp/sql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type CloudBuild struct {
	pulumi.ResourceState

	ID pulumi.StringOutput
}

func (a *NitricGcpPulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {
	dbConfig := &sql.DatabaseArgs{
		Name:     pulumi.String(name),
		Instance: a.masterDb.Name,
		Project:  pulumi.String(a.GcpConfig.ProjectId),
	}

	if a.GcpConfig.Databases[name] != nil {
		dbConfig.DeletionPolicy = pulumi.String(a.GcpConfig.Databases[name].DeletionPolicy)
	}

	_, err := sql.NewDatabase(ctx, name, dbConfig, pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{a.masterDb}))
	if err != nil {
		return err
	}

	if a.DatabaseMigrationBuild[name] == nil && config.GetImageUri() != "" {
		// clientContext := context.TODO()

		imageUriSplit := strings.Split(config.GetImageUri(), "/")
		imageName := imageUriSplit[len(imageUriSplit)-1]

		inspect, err := image.CommandFromImageInspect(config.GetImageUri(), " ")
		if err != nil {
			return err
		}

		repo, err := artifactregistry.NewRepository(ctx, fmt.Sprintf("%s-migration-repo-%s", a.StackId, name), &artifactregistry.RepositoryArgs{
			Location:     pulumi.String(a.Region),
			RepositoryId: pulumi.Sprintf("%s-migration-repo-%s", a.StackId, name),
			Format:       pulumi.String("DOCKER"),
		})
		if err != nil {
			return err
		}

		imageUrl := pulumi.Sprintf("%s-docker.pkg.dev/%s/%s/%s", a.Region, a.GcpConfig.ProjectId, repo.Name, imageName)

		newTag, err := docker.NewTag(ctx, name+"-tag", &docker.TagArgs{
			SourceImage: pulumi.String(inspect.ID),
			TargetImage: imageUrl,
		}, pulumi.Parent(parent))
		if err != nil {
			return err
		}

		image, err := docker.NewRegistryImage(ctx, name+"-remote", &docker.RegistryImageArgs{
			Name: imageUrl,
			Triggers: pulumi.Map{
				"imageSha": pulumi.String(inspect.ID),
			},
		}, pulumi.Parent(parent), pulumi.Provider(a.DockerProvider), pulumi.DependsOn([]pulumi.Resource{newTag}))
		if err != nil {
			return err
		}

		databaseUrl := pulumi.Sprintf("postgres://%s:%s@%s:%s/%s", "postgres", a.dbMasterPassword.Result, a.masterDb.PrivateIpAddress, "5432", name)

		// Run as google cloud run jobs instead of using cloud build
		// This way we don't need to configre private worker pools (can share VPC connectory with cloud run services)

		a.DatabaseMigrationBuild[name], err = cloudrunv2.NewJob(ctx, name+"-migration", &cloudrunv2.JobArgs{
			Location:            pulumi.String(a.Region),
			StartExecutionToken: pulumi.String("start-once-created"),
			DeletionProtection:  pulumi.Bool(false),
			Template: &cloudrunv2.JobTemplateArgs{
				Template: &cloudrunv2.JobTemplateTemplateArgs{
					VpcAccess: &cloudrunv2.JobTemplateTemplateVpcAccessArgs{
						Connector: a.vpcConnector.SelfLink,
						Egress:    pulumi.String("PRIVATE_RANGES_ONLY"),
						// TODO: Re-enable when pulumi network interface support is fixed for tear down
						// NetworkInterfaces: cloudrunv2.JobTemplateTemplateVpcAccessNetworkInterfaceArray{
						// 	&cloudrunv2.JobTemplateTemplateVpcAccessNetworkInterfaceArgs{
						// 		Network:    a.privateNetwork.ID(),
						// 		Subnetwork: a.privateSubnet.ID(),
						// 	},
						// },
					},
					Containers: cloudrunv2.JobTemplateTemplateContainerArray{
						&cloudrunv2.JobTemplateTemplateContainerArgs{
							Image: image.Name,
							Envs: cloudrunv2.JobTemplateTemplateContainerEnvArray{
								&cloudrunv2.JobTemplateTemplateContainerEnvArgs{
									Name:  pulumi.String("NITRIC_DB_NAME"),
									Value: pulumi.String(name),
								},
								&cloudrunv2.JobTemplateTemplateContainerEnvArgs{
									Name:  pulumi.String("DB_URL"),
									Value: databaseUrl,
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// buildId := pulumi.All(databaseUrl, a.cloudBuildWorkerPool.ID().ToStringOutput(), image.Name).ApplyT(func(args []interface{}) (string, error) {
		// 	creds, err := google.FindDefaultCredentials(clientContext)
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	client, err := cloudbuild.NewClient(clientContext, option.WithCredentials(creds), option.WithQuotaProject(a.GcpConfig.ProjectId))
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	defer client.Close()

		// 	url := args[0].(string)
		// 	workerPoolId := args[1].(string)
		// 	imageUri := args[2].(string)

		// 	build, err := client.CreateBuild(clientContext, &cloudbuildpb.CreateBuildRequest{
		// 		Parent:    fmt.Sprintf("projects/%s/locations/%s", a.GcpConfig.ProjectId, a.Region),
		// 		ProjectId: a.GcpConfig.ProjectId,
		// 		Build: &cloudbuildpb.Build{
		// 			Substitutions: map[string]string{
		// 				"_DATABASE_NAME": name,
		// 				"_DATABASE_URL":  url,
		// 			},
		// 			Steps: []*cloudbuildpb.BuildStep{
		// 				{
		// 					Name: imageUri,
		// 					Dir:  "/",
		// 					Env: []string{
		// 						"NITRIC_DB_NAME=${_DATABASE_NAME}",
		// 						"DB_URL=${_DATABASE_URL}",
		// 					},
		// 				},
		// 			},
		// 			Options: &cloudbuildpb.BuildOptions{
		// 				Pool: &cloudbuildpb.BuildOptions_PoolOption{
		// 					Name: workerPoolId,
		// 				},
		// 			},
		// 		},
		// 	})
		// 	if err != nil {
		// 		return "", fmt.Errorf("error creating build for db %s: %w", name, err)
		// 	}

		// 	err = retry.Do(func() error {
		// 		metadata, err := build.Metadata()
		// 		if err != nil {
		// 			return retry.Unrecoverable(fmt.Errorf("failed to retrieve metadata: %w", err))
		// 		}

		// 		if metadata == nil {
		// 			return fmt.Errorf("unable to retrieve metadata")
		// 		}

		// 		currBuild, err := client.GetBuild(clientContext, &cloudbuildpb.GetBuildRequest{
		// 			Id:        metadata.Build.Id,
		// 			Name:      fmt.Sprintf("projects/%s/locations/%s/builds/%s", a.GcpConfig.ProjectId, a.Region, metadata.Build.Id),
		// 			ProjectId: a.GcpConfig.ProjectId,
		// 		})
		// 		if err != nil {
		// 			if strings.Contains(err.Error(), "not found") {
		// 				return fmt.Errorf("build operation not found: %w", err)
		// 			}
		// 			return fmt.Errorf("failed to poll build: %w", err)
		// 		}

		// 		if currBuild.Status == cloudbuildpb.Build_PENDING || currBuild.Status == cloudbuildpb.Build_WORKING || currBuild.Status == cloudbuildpb.Build_QUEUED {
		// 			return fmt.Errorf("build still in progress with status: %s", currBuild.Status)
		// 		}

		// 		if currBuild.Status == cloudbuildpb.Build_SUCCESS {
		// 			return nil
		// 		}

		// 		return retry.Unrecoverable(fmt.Errorf("build failed with status: %s", currBuild.Status))
		// 	}, retry.Attempts(10), retry.Delay(10*time.Second))
		// 	if err != nil {
		// 		return "", err
		// 	}

		// 	return build.Name(), nil
		// }).(pulumi.StringOutput)

		// res := &CloudBuild{
		// 	ID: buildId,
		// }

		// err = ctx.RegisterComponentResource("nitricgcp:cloudbuild:Build", name, res, pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{a.masterDb}))
		// if err != nil {
		// 	return err
		// }

		// res.ID = buildId

		// err = ctx.RegisterResourceOutputs(res, pulumi.Map{
		// 	"ID": buildId,
		// })
		// if err != nil {
		// 	return err
		// }

		// ctx.Export(fmt.Sprintf("migration-%s-build-id", name), buildId)
		// a.DatabaseMigrationBuild[name] = res
	}

	return nil
}
