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
	"github.com/nitrictech/nitric/cloud/common/deploy/image"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/service"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

func (a *NitricGcpTerraformProvider) Service(stack cdktf.TerraformStack, name string, config *deploymentspb.Service, runtimeProvider provider.RuntimeProvider) error {
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

	if config.Type == "" {
		config.Type = "default"
	}

	typeConfig, hasConfig := a.GcpConfig.Config[config.Type]
	if !hasConfig {
		return fmt.Errorf("could not find config for type %s in %+v", config.Type, a.GcpConfig)
	}

	jsiiEnv := map[string]*string{
		"NITRIC_STACK_ID":        a.Stack.StackIdOutput(),
		"NITRIC_ENVIRONMENT":     jsii.String("cloud"),
		"MIN_WORKERS":            jsii.String(fmt.Sprint(config.Workers)),
		"NITRIC_HTTP_PROXY_PORT": jsii.String(fmt.Sprint(3000)),
	}
	for k, v := range config.GetEnv() {
		jsiiEnv[k] = jsii.String(v)
	}

	a.Services[name] = service.NewService(stack, jsii.Sprintf("service_%s", name), &service.ServiceConfig{
		ProjectId:       jsii.String(a.GcpConfig.ProjectId),
		Region:          jsii.String(a.Region),
		ServiceName:     jsii.String(name[:min(len(name), 63)]),
		Image:           jsii.String(imageId),
		Environment:     &jsiiEnv,
		StackId:         a.Stack.StackIdOutput(),
		Cpus:            jsii.Number(typeConfig.CloudRun.Cpus),
		Gpus:            jsii.Number(typeConfig.CloudRun.Gpus),
		MemoryMb:        jsii.Number(typeConfig.CloudRun.Memory),
		TimeoutSeconds:  jsii.Number(typeConfig.CloudRun.Timeout),
		BaseComputeRole: a.Stack.BaseComputeRoleOutput(),
	})

	return nil
}

// func normalizeServiceName(name string) string {
// 	name = strings.ToLower(name)

// 	// Replace all non-alphanumeric characters with `-`
// 	replace := regexp.MustCompile("[^a-zA-Z0-9]")
// 	name = replace.ReplaceAllString(name, "-")

// 	// Replace all multiple `-` with single `-`
// 	replace = regexp.MustCompile("-+")
// 	name = replace.ReplaceAllString(name, "-")

// 	// Replace all leading and trailing `-` with empty string
// 	replace = regexp.MustCompile("^-|-$")
// 	name = replace.ReplaceAllString(name, "")

// 	// Truncate the name to 63 characters
// 	if len(name) > 23 {
// 		name = name[:23]
// 	}

// 	return name
// }
