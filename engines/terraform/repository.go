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
		ServiceBlueprints: map[string]*ServiceBlueprint{
			"default": {
				ResourceBlueprint: &ResourceBlueprint{
					PluginId: "nitric-aws-lambda",
					Properties: map[string]interface{}{
						"timeout":                "${var.lambda_timeout}",
						"function_url_auth_type": "${var.function_url_auth_type}",
					},
				},
				IdentitiesBlueprint: &IdentitiesBlueprint{
					Identities: []ResourceBlueprint{
						ResourceBlueprint{
							PluginId:   "nitric-aws-iam-role",
							Properties: map[string]interface{}{},
						},
					},
				},
			},
		},
		EntrypointBlueprints: map[string]*ResourceBlueprint{
			"default": {
				PluginId:   "nitric-aws-cloudfront",
				Properties: map[string]interface{}{},
			},
		},
		InfraSpecs: map[string]*ResourceBlueprint{
			"vpc": {
				PluginId:   "nitric-aws-vpc",
				Properties: map[string]interface{}{},
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
