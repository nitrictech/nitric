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

	pulumiutils "github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DownStreamMessageWriter struct {
	stream deploy.DeployService_DownServer
}

func (s *DownStreamMessageWriter) Write(bytes []byte) (int, error) {
	err := s.stream.Send(&deploy.DeployDownEvent{
		Content: &deploy.DeployDownEvent_Message{
			Message: &deploy.DeployEventMessage{
				Message: string(bytes),
			},
		},
	})
	if err != nil {
		return 0, err
	}

	return len(bytes), nil
}

func (d *DeployServer) Down(request *deploy.DeployDownRequest, stream deploy.DeployService_DownServer) error {
	details, err := getStackDetailsFromAttributes(request.Attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	// TODO: Tear down the requested stack
	dsMessageWriter := &pulumiutils.DownStreamMessageWriter{
		Stream: stream,
	}

	s, err := auto.UpsertStackInlineSource(context.TODO(), details.FullStackName, details.Project, nil)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	// destroy the stack
	_, err = s.Destroy(context.TODO(), optdestroy.ProgressStreams(dsMessageWriter))
	if err != nil {
		return err
	}

	return nil
}
