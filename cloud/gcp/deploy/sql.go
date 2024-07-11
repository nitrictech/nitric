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
	"strings"
	"time"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	"github.com/avast/retry-go"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi-docker/sdk/v4/go/docker"
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/sql"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/samber/lo"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type CloudBuild struct {
	pulumi.ResourceState

	ID pulumi.StringOutput
}

func (a *NitricGcpPulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {
	imageUriSplit := strings.Split(config.GetImageUri(), "/")
	imageName := imageUriSplit[len(imageUriSplit)-1]

	inspect, err := image.CommandFromImageInspect(config.GetImageUri(), " ")
	if err != nil {
		return err
	}

	repoUrl := pulumi.Sprintf("gcr.io/%s/%s", a.GcpConfig.ProjectId, imageName)

	newTag, err := docker.NewTag(ctx, name+"-tag", &docker.TagArgs{
		SourceImage: pulumi.String(inspect.ID),
		TargetImage: repoUrl,
	}, pulumi.Parent(parent))
	if err != nil {
		return err
	}

	image, err := docker.NewRegistryImage(ctx, name+"-remote", &docker.RegistryImageArgs{
		Name: repoUrl,
		Triggers: pulumi.Map{
			"imageSha": pulumi.String(inspect.ID),
		},
	}, pulumi.Parent(parent), pulumi.Provider(a.DockerProvider), pulumi.DependsOn([]pulumi.Resource{newTag}))
	if err != nil {
		return err
	}

	_, err = sql.NewDatabase(ctx, name, &sql.DatabaseArgs{
		Name:           pulumi.String(name),
		Instance:       a.masterDb.Name,
		DeletionPolicy: pulumi.String("DELETE"),
		Project:        pulumi.String(a.GcpConfig.ProjectId),
	}, pulumi.Parent(parent), pulumi.DependsOn([]pulumi.Resource{a.masterDb}))
	if err != nil {
		return err
	}

	if a.DatabaseMigrationBuild[name] == nil && config.GetImageUri() != "" {
		clientContext := context.TODO()

		databaseUrl := pulumi.Sprintf("postgres://%s:%s@%s:%s/%s", "postgres", a.dbMasterPassword.Result, a.masterDb.PrivateIpAddress, "5432", name)

		buildId := pulumi.All(databaseUrl, a.cloudBuildWorkerPool.ID().ToStringOutput(), image.Name, a.masterDb.ToDatabaseInstanceOutput()).ApplyT(func(args []interface{}) (string, error) {
			creds, err := google.FindDefaultCredentials(clientContext)
			if err != nil {
				return "", err
			}

			client, err := cloudbuild.NewClient(clientContext, option.WithCredentials(creds), option.WithQuotaProject(a.GcpConfig.ProjectId))
			if err != nil {
				return "", err
			}

			defer client.Close()

			url := args[0].(string)
			workerPoolId := args[1].(string)
			imageUri := args[2].(string)

			build, err := client.CreateBuild(clientContext, &cloudbuildpb.CreateBuildRequest{
				Parent:    fmt.Sprintf("projects/%s/locations/%s", a.GcpConfig.ProjectId, a.Region),
				ProjectId: a.GcpConfig.ProjectId,
				Build: &cloudbuildpb.Build{
					Substitutions: map[string]string{
						"_DATABASE_NAME": name,
						"_DATABASE_URL":  url,
					},
					Steps: []*cloudbuildpb.BuildStep{
						{
							Name: imageUri,
							Dir:  "/",
							Env: []string{
								"NITRIC_DB_NAME=${_DATABASE_NAME}",
								"DB_URL=${_DATABASE_URL}",
							},
						},
					},
					Options: &cloudbuildpb.BuildOptions{
						Pool: &cloudbuildpb.BuildOptions_PoolOption{
							Name: workerPoolId,
						},
					},
				},
			})
			if err != nil {
				return "", fmt.Errorf("error creating build for db %s: %w", name, err)
			}

			err = retry.Do(func() error {
				metadata, err := build.Metadata()
				if err != nil {
					return retry.Unrecoverable(err)
				}

				if metadata.Build.Status == cloudbuildpb.Build_SUCCESS {
					return nil
				} else if lo.Contains([]cloudbuildpb.Build_Status{
					cloudbuildpb.Build_PENDING,
					cloudbuildpb.Build_WORKING,
					cloudbuildpb.Build_QUEUED,
					cloudbuildpb.Build_STATUS_UNKNOWN,
				}, metadata.Build.Status) {
					return fmt.Errorf("build still in progress with status: %s", metadata.Build.Status)
				} else {
					return retry.Unrecoverable(fmt.Errorf("build failed with status: %s", metadata.Build.Status))
				}
			}, retry.Attempts(10), retry.Delay(10*time.Second))
			if err != nil {
				return "", err
			}

			return build.Name(), nil
		}).(pulumi.StringOutput)

		res := &CloudBuild{
			ID: buildId,
		}

		err = ctx.RegisterComponentResource("nitricgcp:cloudbuild:Build", name, res, pulumi.Parent(parent))
		if err != nil {
			return err
		}

		res.ID = buildId

		err = ctx.RegisterResourceOutputs(res, pulumi.Map{
			"ID": buildId,
		})
		if err != nil {
			return err
		}

		ctx.Export(fmt.Sprintf("migration-%s-build-id", name), buildId)
		a.DatabaseMigrationBuild[name] = res
	}

	return nil
}
