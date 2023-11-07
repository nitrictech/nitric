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

	tea "github.com/charmbracelet/bubbletea"
	commonDeploy "github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/output/interactive"
	"github.com/nitrictech/nitric/cloud/common/deploy/output/noninteractive"
	pulumiutils "github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/config"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iamcredentials/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Up - Deploy requested infrastructure for a stack
func (d *DeployServer) Up(request *deploy.DeployUpRequest, stream deploy.DeployService_UpServer) error {
	details, err := getStackDetailsFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	config, err := config.ConfigFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// If we're interactive then we want to provide
	outputStream := &pulumiutils.UpStreamMessageWriter{
		Stream: stream,
	}

	// Default to the non-interactive writer
	pulumiUpOpts := []optup.Option{
		optup.ProgressStreams(noninteractive.NewNonInterativeOutput(outputStream)),
	}

	if request.Interactive {
		pulumiEventChan := make(chan events.EngineEvent)
		deployModel := interactive.NewOutputModel(make(chan tea.Msg), pulumiEventChan)
		teaProgram := tea.NewProgram(deployModel, tea.WithOutput(&pulumiutils.UpStreamMessageWriter{
			Stream: stream,
		}))

		pulumiUpOpts = []optup.Option{
			optup.ProgressStreams(deployModel),
			optup.EventStreams(pulumiEventChan),
		}

		//nolint:errcheck
		go teaProgram.Run()
		// Close the program when we're done
		defer teaProgram.Quit()
	}

	pulumiStack, err := NewUpProgram(context.TODO(), details, config, request.Spec)
	if err != nil {
		return err
	}

	err = pulumiStack.SetAllConfig(context.TODO(), auto.ConfigMap{
		"gcp:region":     auto.ConfigValue{Value: details.Region},
		"gcp:project":    auto.ConfigValue{Value: details.ProjectId},
		"gcp:version":    auto.ConfigValue{Value: pulumiGcpVersion},
		"docker:version": auto.ConfigValue{Value: commonDeploy.PulumiDockerVersion},
	})
	if err != nil {
		return err
	}

	err = pulumiStack.SetConfig(context.TODO(), "gcp:project", auto.ConfigValue{Value: details.ProjectId})
	if err != nil {
		return err
	}

	if config.Refresh {
		// refresh the stack first
		_, err := pulumiStack.Refresh(context.TODO())
		if err != nil {
			return err
		}
	}

	// Run the program
	res, err := pulumiStack.Up(context.TODO(), pulumiUpOpts...)
	if err != nil {
		return err
	}

	// Send terminal message
	err = stream.Send(pulumiutils.PulumiOutputsToResult(res.Outputs))

	return err
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
			return nil, fmt.Errorf("unable to impersonate service account")
		}

		token = &oauth2.Token{AccessToken: accessToken.AccessToken}
	}

	if token == nil {
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
