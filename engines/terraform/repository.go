package terraform

type PlatformRepository interface {
	// terraform/nitric-aws
	GetPlatform(string) PlatformSpec
}

type MockPlatformRepository struct {
}

func (MockPlatformRepository) GetPlatform(name string) *PlatformSpec {
	return &PlatformSpec{
		Name: "aws",
		ServicesSpec: NitricResourceSpec{
			ResourceSpec: ResourceSpec{
				PluginId: "nitric-aws-lambda",
				Properties: map[string]interface{}{
					"timeout":                "${var.lambda_timeout}",
					"function_url_auth_type": "${var.function_url_auth_type}",
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
	}
}

func NewMockPlatformRepository() *MockPlatformRepository {
	return &MockPlatformRepository{}
}
