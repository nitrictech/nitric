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

	_ "embed"

	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type DeployServer struct {
	deploy.UnimplementedDeployServiceServer
}

// Embeds the runtime directly into the deploytime binary
// This way the versions will always match as they're always built and versioned together (as a single artifact)
// This should also help with docker build speeds as the runtime has already been "downloaded"
//
//go:embed runtime-gcp
var runtime []byte

type StackDetails struct {
	Project   string
	ProjectId string
	Stack     string
	Region    string
}

// Read nitric attributes from the provided deployment attributes
func getStackDetailsFromAttributes(attributes map[string]string) (*StackDetails, error) {
	project, ok := attributes["project"]
	if !ok || project == "" {
		// need a valid project name
		return nil, fmt.Errorf("project is not set of invalid")
	}

	projectId, ok := attributes["gcp-project-id"]
	if !ok || projectId == "" {
		// need a valid project name
		return nil, fmt.Errorf("gcp-project-id is not set of invalid")
	}

	stack, ok := attributes["stack"]
	if !ok || stack == "" {
		// need a valid stack name
		return nil, fmt.Errorf("stack is not set of invalid")
	}

	region, ok := attributes["region"]
	if !ok || stack == "" {
		// need a valid stack name
		return nil, fmt.Errorf("region is not set of invalid")
	}

	return &StackDetails{
		Project:   project,
		ProjectId: projectId,
		Stack:     stack,
		Region:    region,
	}, nil
}

func NewServer() (*DeployServer, error) {
	return &DeployServer{}, nil
}
