// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deploytf

import (
	"embed"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/common"
	"github.com/nitrictech/nitric/cloud/common/deploy"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NitricAzureTerraformProvider struct {
	*deploy.CommonStackDetails

	AzureConfig *common.AzureConfig

	provider.NitricDefaultOrder
}

var _ provider.NitricTerraformProvider = (*NitricAzureTerraformProvider)(nil)

func (a *NitricAzureTerraformProvider) Init(attributes map[string]interface{}) error {
	var err error

	a.CommonStackDetails, err = deploy.CommonStackDetailsFromAttributes(attributes)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	a.AzureConfig, err = common.ConfigFromAttributes(attributes)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "Bad stack configuration: %s", err)
	}

	return nil
}

// embed the modules directory here
//
//go:embed .nitric/modules/**/*
var modules embed.FS

func (a *NitricAzureTerraformProvider) CdkTfModules() ([]provider.ModuleDirectory, error) {
	return []provider.ModuleDirectory{
		{
			ParentDir: ".nitric/modules",
			Modules:   modules,
		},
	}, nil
}

func (a *NitricAzureTerraformProvider) RequiredProviders() map[string]interface{} {
	return map[string]interface{}{}
}

func (a *NitricAzureTerraformProvider) Pre(stack cdktf.TerraformStack, resources []*deploymentspb.Resource) error {
	return nil
}

func (a *NitricAzureTerraformProvider) Post(stack cdktf.TerraformStack) error {
	return nil
}

func NewNitricAzureProvider() *NitricAzureTerraformProvider {
	return &NitricAzureTerraformProvider{}
}
