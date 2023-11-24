// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
	"github.com/nitrictech/nitric/cloud/common/deploy/output/interactive"
	"github.com/nitrictech/nitric/cloud/common/deploy/output/noninteractive"
	pulumiutils "github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optpreview"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/nitrictech/nitric/cloud/azure/deploy/config"
	commonDeploy "github.com/nitrictech/nitric/cloud/common/deploy"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

// Preview - Dry run deployment of requested infrastructure for a stack
func (d *DeployServer) Preview(request *deploy.DeployPreviewRequest, stream deploy.DeployService_PreviewServer) error {
	details, err := getStackDetailsFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	config, err := config.ConfigFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// If we're interactive then we want to provide
	outputStream := &pulumiutils.PreviewStreamMessageWriter{
		Stream: stream,
	}

	// Default to the non-interactive writer
	pulumiPreviewOpts := []optpreview.Option{
		optpreview.ProgressStreams(noninteractive.NewNonInterativeOutput(outputStream)),
	}

	var interactiveProgram *interactive.Program
	if request.Interactive {
		pulumiEventChan := make(chan events.EngineEvent)
		deployModel, err := interactive.NewOutputModel(make(chan tea.Msg), pulumiEventChan)
		if err != nil {
			return err
		}

		pulumiPreviewOpts = []optpreview.Option{
			optpreview.ProgressStreams(deployModel),
			optpreview.EventStreams(pulumiEventChan),
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
		"azure-native:location": auto.ConfigValue{Value: details.Region},
		"azure:location":        auto.ConfigValue{Value: details.Region},
		"azure-native:version":  auto.ConfigValue{Value: pulumiAzureNativeVersion},
		"azure:version":         auto.ConfigValue{Value: pulumiAzureVersion},
		"docker:version":        auto.ConfigValue{Value: commonDeploy.PulumiDockerVersion},
	})
	if err != nil {
		return err
	}

	if config.Refresh {
		_ = stream.Send(&deploy.DeployPreviewEvent{
			Content: &deploy.DeployPreviewEvent_Message{
				Message: &deploy.DeployEventMessage{
					Message: "refreshing pulumi stack",
				},
			},
		})

		// TODO: Handle refresh logs
		_, err := pulumiStack.Refresh(context.TODO())
		if err != nil {
			return err
		}
	}

	_, err = pulumiStack.Preview(context.TODO(), pulumiPreviewOpts...)
	// Run the program
	if err != nil {
		return err
	}

	return err
}
