package container


type ContainerCapabilities struct {
	// List of linux capabilities to add.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#add Container#add}
	Add *[]*string `field:"optional" json:"add" yaml:"add"`
	// List of linux capabilities to drop.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#drop Container#drop}
	Drop *[]*string `field:"optional" json:"drop" yaml:"drop"`
}

