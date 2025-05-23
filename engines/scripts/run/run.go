package main

import (
	"bytes"
	"encoding/json"
	"log"

	app_spec_schema "github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/nitrictech/nitric/engines/terraform"
)

type MockTerraformPluginRepository struct {
	plugins map[string]*terraform.PluginManifest
}

func (r *MockTerraformPluginRepository) GetPlugin(name string) (*terraform.PluginManifest, error) {
	return r.plugins[name], nil
}

func createMockTerraformPluginRepository() *MockTerraformPluginRepository {
	return &MockTerraformPluginRepository{
		plugins: map[string]*terraform.PluginManifest{
			"nitric-aws-lambda": {
				Name: "nitric-aws-lambda",
				Deployment: terraform.DeploymentModule{
					Terraform: "terraform-aws-modules/lambda/aws",
				},
			},
			"nitric-aws-cloudfront": {
				Name: "nitric-aws-cloudfront",
				Deployment: terraform.DeploymentModule{
					Terraform: "terraform-aws-modules/cloudfront/aws",
				},
			},
			"nitric-aws-vpc": {
				Name: "nitric-aws-vpc",
				Deployment: terraform.DeploymentModule{
					Terraform: "terraform-aws-modules/vpc/aws",
				},
			},
		},
	}
}

func main() {
	platformConfig := terraform.PlatformSpec{
		Name: "aws",
		ServicesSpec: terraform.NitricResourceSpec{
			ResourceSpec: terraform.ResourceSpec{
				PluginId: "nitric-aws-lambda",
				Properties: map[string]interface{}{
					"timeout":                "${var.lambda_timeout}",
					"function_url_auth_type": "${var.function_url_auth_type}",
				},
			},
		},
		EntrypointsSpec: terraform.NitricResourceSpec{
			ResourceSpec: terraform.ResourceSpec{
				PluginId:   "nitric-aws-cloudfront",
				Properties: map[string]interface{}{},
			},
		},
		Infra: map[string]terraform.InfraResourceSpec{
			"vpc": {
				ResourceSpec: terraform.ResourceSpec{
					PluginId:   "nitric-aws-vpc",
					Properties: map[string]interface{}{},
				},
			},
		},
	}

	// serialize the platform config to json
	platformConfigJSON, err := json.Marshal(platformConfig)
	if err != nil {
		log.Fatalf("failed to marshal platform config: %v", err)
	}

	mockRepository := terraform.NewNitricTerraformPluginRepository()

	// provide a bytes reader to the terraform engine
	platform := terraform.NewFromFile(bytes.NewReader(platformConfigJSON), terraform.WithRepository(mockRepository))

	err = platform.Apply(&app_spec_schema.Application{
		Name: "test",
		Resources: map[string]app_spec_schema.Resource{
			"service": {
				Type: "service",
				ServiceResource: &app_spec_schema.ServiceResource{
					Port: 3000,
					Env: map[string]string{
						"TEST": "test",
						"PORT": "3000",
					},
					Container: app_spec_schema.Container{
						Image: &app_spec_schema.DockerImage{
							ID: "ealen/echo-server:latest",
						},
					},
				},
			},
			"ingress": {
				Type: "entrypoint",
				EntrypointResource: &app_spec_schema.EntrypointResource{
					Routes: map[string]app_spec_schema.Route{
						"/": {
							TargetName: "service",
						},
					},
				},
			},
		},
	})

	if err != nil {
		log.Fatalf("failed to apply platform: %v", err)
	}
}
