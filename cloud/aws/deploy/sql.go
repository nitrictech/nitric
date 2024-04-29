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
	"time"

	"github.com/avast/retry-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Sqldatabase - Implements PostgresSql database deployments use AWS Aurora
func (a *NitricAwsPulumiProvider) SqlDatabase(ctx *pulumi.Context, parent pulumi.Resource, name string, config *deploymentspb.SqlDatabase) error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(a.Region), // replace with your AWS region
	})
	if err != nil {
		return err
	}

	client := codebuild.New(sess)

	// Run the database creation step
	a.CreateDatabaseProject.Name.ApplyT(func(projectName string) (bool, error) {
		fmt.Printf("Starting database creation build %s\n", name)
		out, err := client.StartBuild(&codebuild.StartBuildInput{
			ProjectName: aws.String(projectName),
			EnvironmentVariablesOverride: []*codebuild.EnvironmentVariable{
				{
					Name:  aws.String("DB_NAME"),
					Value: aws.String(name),
				},
			},
		})
		if err != nil {
			return false, err
		}

		var finalErr error
		err = retry.Do(func() error {
			fmt.Printf("Waiting for database creation build %s\n", name)
			resp, err := client.BatchGetBuilds(&codebuild.BatchGetBuildsInput{
				Ids: []*string{out.Build.Id},
			})
			if err != nil {
				return err
			}

			status := aws.StringValue(resp.Builds[0].BuildStatus)
			if status != codebuild.StatusTypeInProgress {
				if status == codebuild.StatusTypeFailed {
					finalErr = fmt.Errorf("database creation build %s failed", name)
				}

				fmt.Printf("Complete database creation build %s\n", name)
				return nil
			}

			fmt.Printf("Still waiting database creation build %s\n", name)
			return fmt.Errorf("build still in progress")
		}, retry.Attempts(10), retry.Delay(time.Second*15))
		if err != nil {
			return false, err
		}

		if finalErr != nil {
			return false, finalErr
		}

		return true, nil
	})

	return nil
}
