package container


type ContainerPorts struct {
	// Port within the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#internal Container#internal}
	Internal *float64 `field:"required" json:"internal" yaml:"internal"`
	// Port exposed out of the container. If not given a free random port `>= 32768` will be used.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#external Container#external}
	External *float64 `field:"optional" json:"external" yaml:"external"`
	// IP address/mask that can access this port. Defaults to `0.0.0.0`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#ip Container#ip}
	Ip *string `field:"optional" json:"ip" yaml:"ip"`
	// Protocol that can be used over this port. Defaults to `tcp`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#protocol Container#protocol}
	Protocol *string `field:"optional" json:"protocol" yaml:"protocol"`
}

