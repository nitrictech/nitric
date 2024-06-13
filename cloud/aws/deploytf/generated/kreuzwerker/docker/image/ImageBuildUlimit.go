package image


type ImageBuildUlimit struct {
	// soft limit.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#hard Image#hard}
	Hard *float64 `field:"required" json:"hard" yaml:"hard"`
	// type of ulimit, e.g. `nofile`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#name Image#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// hard limit.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#soft Image#soft}
	Soft *float64 `field:"required" json:"soft" yaml:"soft"`
}

