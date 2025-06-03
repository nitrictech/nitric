package terraform

import "fmt"

type PlatformRepository interface {
	// terraform/nitric-aws
	GetPlatform(string) (*PlatformSpec, error)
}

type MockPlatformRepository struct {
}

var platformSpecs = map[string]*PlatformSpec{
	"aws": {
		Name: "aws",
		ServicesSpec: NitricServiceSpec{
			ServiceSpec: ServiceSpec{
				ResourceSpec: ResourceSpec{
					PluginId: "nitric-aws-lambda",
					Properties: map[string]interface{}{
						"timeout":                "${var.lambda_timeout}",
						"function_url_auth_type": "${var.function_url_auth_type}",
					},
				},
				Identities: map[string]InfraResourceSpec{
					"aws:iam:role": InfraResourceSpec{
						ResourceSpec: ResourceSpec{
							PluginId:   "nitric-aws-iam-role",
							Properties: map[string]interface{}{},
						},
					},
				},
			},
		},
		EntrypointsSpec: NitricResourceSpec{
			ResourceSpec: ResourceSpec{
				PluginId:   "nitric-aws-cloudfront",
				Properties: map[string]interface{}{},
			},
		},
		Infra: map[string]InfraResourceSpec{
			"vpc": {
				ResourceSpec: ResourceSpec{
					PluginId:   "nitric-aws-vpc",
					Properties: map[string]interface{}{},
				},
			},
		},
	},
}

func (MockPlatformRepository) GetPlatform(name string) (*PlatformSpec, error) {
	platform, ok := platformSpecs[name]
	if !ok {
		return nil, fmt.Errorf("no platform %s available in repository")
	}

	return platform, nil
}

func NewMockPlatformRepository() *MockPlatformRepository {
	return &MockPlatformRepository{}
}
