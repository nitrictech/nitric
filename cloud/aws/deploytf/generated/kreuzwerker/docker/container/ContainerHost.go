package container


type ContainerHost struct {
	// Hostname to add.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#host Container#host}
	Host *string `field:"required" json:"host" yaml:"host"`
	// IP address this hostname should resolve to.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#ip Container#ip}
	Ip *string `field:"required" json:"ip" yaml:"ip"`
}

