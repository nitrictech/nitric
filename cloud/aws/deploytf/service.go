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
	err := image.BuildWrappedImage(&image.BuildWrappedImageArgs{
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
		"NITRIC_STACK_ID":        a.Stack.StackIdOutput(),
		"NITRIC_ENVIRONMENT":     jsii.String("cloud"),
		"MIN_WORKERS":            jsii.String(fmt.Sprint(config.Workers)),
		"NITRIC_HTTP_PROXY_PORT": jsii.String(fmt.Sprint(3000)),
	}
	for k, v := range config.GetEnv() {
		jsiiEnv[k] = jsii.String(v)
	}

	a.Services[name] = service.NewService(stack, jsii.String(name), &service.ServiceConfig{
		ServiceName: jsii.String(name),
		// TODO: Match to target image property above
		Image: jsii.String(name),

		// Convert env to pointer type
		Environment: &jsiiEnv,
	})

	return nil
}
