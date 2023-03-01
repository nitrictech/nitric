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

	commonDeploy "github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/pulumi"
	deploy "github.com/nitrictech/nitric/core/pkg/api/nitric/deploy/v1"
)

type DeployServer struct {
	deploy.UnimplementedDeployServiceServer
}

// Embeds the runtime directly into the deploytime binary
// This way the versions will always match as they're always built and versioned together (as a single artifact)
// This should also help with docker build speeds as the runtime has already been "downloaded"
//
//go:embed runtime-azure
var runtime []byte

type StackDetails struct {
	*commonDeploy.CommonStackDetails
	Org        string
	AdminEmail string
}

// Read nitric attributes from the provided deployment attributes
func getStackDetailsFromAttributes(attributes map[string]string) (*StackDetails, error) {
	commonDetails, err := commonDeploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return nil, err
	}

	org, ok := attributes["org"]
	if !ok || org == "" {
		// need a valid stack name
		return nil, fmt.Errorf("org is not set or invalid")
	}

	adminEmail, ok := attributes["adminemail"]
	if !ok || adminEmail == "" {
		// need a valid stack name
		return nil, fmt.Errorf("adminemail is not set or invalid")
	}

	return &StackDetails{
		CommonStackDetails: commonDetails,
		Org:                org,
		AdminEmail:         adminEmail,
	}, nil
}

func NewServer() (*DeployServer, error) {
	err := pulumi.InstallResources()
	if err != nil {
		return nil, err
	}

	return &DeployServer{}, nil
}
