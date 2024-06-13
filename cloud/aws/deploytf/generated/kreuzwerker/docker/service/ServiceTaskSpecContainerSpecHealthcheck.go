package service


type ServiceTaskSpecContainerSpecHealthcheck struct {
	// The test to perform as list.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#test Service#test}
	Test *[]*string `field:"required" json:"test" yaml:"test"`
	// Time between running the check (ms|s|m|h). Defaults to `0s`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#interval Service#interval}
	Interval *string `field:"optional" json:"interval" yaml:"interval"`
	// Consecutive failures needed to report unhealthy. Defaults to `0`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#retries Service#retries}
	Retries *float64 `field:"optional" json:"retries" yaml:"retries"`
	// Start period for the container to initialize before counting retries towards unstable (ms|s|m|h). Defaults to `0s`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#start_period Service#start_period}
	StartPeriod *string `field:"optional" json:"startPeriod" yaml:"startPeriod"`
	// Maximum time to allow one check to run (ms|s|m|h). Defaults to `0s`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#timeout Service#timeout}
	Timeout *string `field:"optional" json:"timeout" yaml:"timeout"`
}

