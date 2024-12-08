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
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/service"
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricAwsTerraformProvider) Service(stack cdktf.TerraformStack, name string, config *deploymentspb.Service, runtimeProvider provider.RuntimeProvider) error {
	targetImageName := fmt.Sprintf("%s-%s", name, a.StackName)
	imageId, err := image.BuildWrappedImage(&image.BuildWrappedImageArgs{
		ServiceName: name,
		SourceImage: config.GetImage().Uri,
		// TODO: Use correct image uri
		TargetImage: targetImageName,
		Runtime:     runtimeProvider(),
	})
	if err != nil {
		return err
	}

	if config.Type == "" {
		config.Type = "default"
	}

	typeConfig, hasConfig := a.AwsConfig.Config[config.Type]
	if !hasConfig {
		return fmt.Errorf("could not find config for type %s in %+v", config.Type, a.AwsConfig)
	}

	jsiiEnv := map[string]*string{
		"NITRIC_STACK_ID":              a.Stack.StackIdOutput(),
		"NITRIC_ENVIRONMENT":           jsii.String("cloud"),
		"MIN_WORKERS":                  jsii.String(fmt.Sprint(config.Workers)),
		"NITRIC_HTTP_PROXY_PORT":       jsii.String(fmt.Sprint(3000)),
		"NITRIC_AWS_RESOURCE_RESOLVER": jsii.String("tagging"),
	}

	// TODO: Only apply to requesting services
	if a.Rds != nil {
		jsiiEnv["NITRIC_DATABASE_BASE_URL"] = jsii.Sprintf("postgres://%s:%s@%s:%s", *a.Rds.ClusterUsernameOutput(), *a.Rds.ClusterPasswordOutput(),
			*a.Rds.ClusterEndpointOutput(), "5432")
	}

	for k, v := range config.GetEnv() {
		jsiiEnv[k] = jsii.String(v)
	}

	serviceConfig := &service.ServiceConfig{
		ServiceName:      jsii.String(name),
		Image:            jsii.String(imageId),
		Environment:      &jsiiEnv,
		StackId:          a.Stack.StackIdOutput(),
		Memory:           jsii.Number(typeConfig.Lambda.Memory),
		Timeout:          jsii.Number(typeConfig.Lambda.Timeout),
		EphemeralStorage: jsii.Number(typeConfig.Lambda.EphemeralStorage),
	}

	if a.Vpc != nil && a.Rds != nil {
		serviceConfig.SecurityGroupIds = &[]*string{(a.Rds.SecurityGroupIdOutput())}
		serviceConfig.SubnetIds = cdktf.Token_AsList(a.Vpc.PrivateSubnetIdsOutput(), &cdktf.EncodingOptions{})
	}

	a.Services[name] = service.NewService(stack, jsii.Sprintf("service_%s", name), serviceConfig)

	return nil
}
