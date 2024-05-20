package deploytf

import (
	"fmt"

	"github.com/aws/jsii-runtime-go"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/service"
	"github.com/nitrictech/nitric/cloud/common/deploy/provider"
	deploymentspb "github.com/nitrictech/nitric/core/pkg/proto/deployments/v1"
)

// def create_provider_runtime_image(source_img: str, image_name: str):
//     """Wrap the source image with with the provider runtime binary"""
//     client = docker.from_env()

//     if not os.path.isdir("/tmp/build"):
//         os.makedirs("/tmp/build")

//     shutil.copy("/app/runtime/Dockerfile", "/tmp/build/Dockerfile")
//     shutil.copy("/app/runtime/runtime", "/tmp/build/runtime")

//     try:
//         # Run the docker build
//         build = client.images.build(
//             path="/tmp/build", tag=image_name, buildargs={"BASE_IMAGE": source_img}
//         )

//         # TODO: print status of the build
//         # for line in build:
//         #     print(line)

//     except Exception as e:
//         raise RuntimeError(f"Could not build image: {image_name}", e)

func (a *NitricAwsTerraformProvider) RuntimeImage(sourceImg string, imageName string) error {
	// Create a new docker api clent
	// shutil.copy("/app/runtime/Dockerfile", "/tmp/build/Dockerfile")
	// shutil.copy("/app/runtime/runtime", "/tmp/build/runtime")

}

func (a *NitricAwsTerraformProvider) Service(stack cdktf.TerraformStack, name string, config *deploymentspb.Service, runtimeProvider provider.RuntimeProvider) error {

	a.Services[name] = service.NewService(stack, jsii.String(name), &service.ServiceConfig{
		ServiceName: jsii.String(name),
		// TODO: Add wrapped image uri here
		// TODO: Determine how repeatable builds can be achieved without regenerating source before deployment
	})
	// Wrap the image ready for deployment

	return fmt.Errorf("nitric AWS terraform provider does not support Service deployment")
}
