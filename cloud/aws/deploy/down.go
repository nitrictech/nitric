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
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/events"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *DeployServer) Down(request *deploy.DeployDownRequest, stream deploy.DeployService_DownServer) error {
	details, err := commonDeploy.CommonStackDetailsFromAttributes(request.Attributes.AsMap())
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// If we're interactive then we want to provide
	outputStream := &pulumiutils.DownStreamMessageWriter{
		Stream: stream,
	}

	pulumiDestroyOpts := []optdestroy.Option{
		optdestroy.ProgressStreams(noninteractive.NewNonInterativeOutput(outputStream)),
	}

	if request.Interactive {
		pulumiEventChan := make(chan events.EngineEvent)
		deployModel := interactive.NewOutputModel(make(chan tea.Msg), pulumiEventChan)
		teaProgram := tea.NewProgram(deployModel, tea.WithOutput(outputStream))
		pulumiDestroyOpts = []optdestroy.Option{
			optdestroy.ProgressStreams(deployModel),
			optdestroy.EventStreams(pulumiEventChan),
		}

		//nolint:errcheck
		go teaProgram.Run()
		// Close the program when we're done
		defer teaProgram.Quit()
	}

	s, err := auto.UpsertStackInlineSource(context.TODO(), details.FullStackName, details.Project, nil)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	// destroy the stack
	_, err = s.Destroy(context.TODO(), pulumiDestroyOpts...)
	if err != nil {
		return err
	}

	return nil
}
