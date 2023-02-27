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

	pulumiutils "github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/auto/optdestroy"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *DeployServer) Down(request *deploy.DeployDownRequest, stream deploy.DeployService_DownServer) error {
	details, err := getStackDetailsFromAttributes(request.Attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	dsMessageWriter := &pulumiutils.DownStreamMessageWriter{
		Stream: stream,
	}

	s, err := auto.UpsertStackInlineSource(context.TODO(), details.Stack, details.Project, nil)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	_, err = s.Destroy(context.TODO(), optdestroy.ProgressStreams(dsMessageWriter))
	if err != nil {
		return err
	}

	return nil
}
