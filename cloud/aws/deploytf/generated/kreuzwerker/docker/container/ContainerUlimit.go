package container


type ContainerUlimit struct {
	// The hard limit.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#hard Container#hard}
	Hard *float64 `field:"required" json:"hard" yaml:"hard"`
	// The name of the ulimit.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#name Container#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// The soft limit.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#soft Container#soft}
	Soft *float64 `field:"required" json:"soft" yaml:"soft"`
}

