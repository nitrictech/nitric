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

	tea "github.com/charmbracelet/bubbletea"
	commonDeploy "github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/output/interactive"
	"github.com/nitrictech/nitric/cloud/common/deploy/output/noninteractive"
	pulumiutils "github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	"github.com/nitrictech/nitric/cloud/gcp/deploy/config"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optup"
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

	var interactiveProgram *interactive.Program
	if request.Interactive {
		pulumiEventChan := make(chan events.EngineEvent)
		deployModel, err := interactive.NewOutputModel(make(chan tea.Msg), pulumiEventChan)
		if err != nil {
			return err
		}

		pulumiUpOpts = []optup.Option{
			optup.ProgressStreams(deployModel),
			optup.EventStreams(pulumiEventChan),
		}

		interactiveProgram = interactive.NewProgram(deployModel, &interactive.ProgramArgs{
			Writer: outputStream,
		})

		go interactiveProgram.Run()

		defer interactiveProgram.Stop()
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
