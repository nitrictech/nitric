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
					"vpc_security_group_ids": "${infra.vpc.default_security_group_id}",
					"vpc_subnet_ids":         "${infra.vpc.infra_subnets}",
					"timeout":                "${var.lambda_timeout}",
				},
			},
		},
		EntrypointsSpec: NitricResourceSpec{
			ResourceSpec: ResourceSpec{
				PluginId: "nitric-aws-cloudfront",
				Properties: map[string]interface{}{
					"region": "${var.region}",
				},
			},
		},
		Infra: map[string]InfraResourceSpec{
			"vpc": {
				ResourceSpec: ResourceSpec{
					PluginId: "nitric-aws-vpc",
					Properties: map[string]interface{}{
						"region": "${var.region}",
					},
				},
			},
		},
	}
}

func NewMockPlatformRepository() *MockPlatformRepository {
	return &MockPlatformRepository{}
}
