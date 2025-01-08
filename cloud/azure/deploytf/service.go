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
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/service"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricAzureTerraformProvider) Service(stack cdktf.TerraformStack, name string, config *deploymentspb.Service, runtimeProvider provider.RuntimeProvider) error {
	imageId, err := image.BuildWrappedImage(&image.BuildWrappedImageArgs{
		ServiceName: name,
		SourceImage: config.GetImage().Uri,
		// TODO: Use correct image uri
		TargetImage: name,
		Runtime:     runtimeProvider(),
	})
	if err != nil {
		return err
	}

	jsiiEnv := map[string]*string{
		"MIN_WORKERS": jsii.String(fmt.Sprint(config.Workers)),
	}

	for k, v := range config.GetEnv() {
		jsiiEnv[k] = jsii.String(v)
	}

	// If the database is enabled, set the database connection string
	// if a.Rds != nil {
	// 	jsiiEnv["NITRIC_DATABASE_BASE_URL"] = jsii.Sprintf("postgres://%s:%s@%s:%s", *a.Rds.ClusterUsernameOutput(), *a.Rds.ClusterPasswordOutput(),
	// 		*a.Rds.ClusterEndpointOutput(), "5432")
	// }

	a.Services[name] = service.NewService(stack, jsii.String(name), &service.ServiceConfig{
		Name:     jsii.String(name),
		ImageUri: jsii.String(imageId),
	})

	return fmt.Errorf("not implemented")
}
