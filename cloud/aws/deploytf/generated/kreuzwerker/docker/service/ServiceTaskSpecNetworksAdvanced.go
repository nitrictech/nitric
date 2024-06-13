package service


type ServiceTaskSpecNetworksAdvanced struct {
	// The name/id of the network.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#name Service#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// The network aliases of the container in the specific network.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#aliases Service#aliases}
	Aliases *[]*string `field:"optional" json:"aliases" yaml:"aliases"`
	// An array of driver options for the network, e.g. `opts1=value`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#driver_opts Service#driver_opts}
	DriverOpts *[]*string `field:"optional" json:"driverOpts" yaml:"driverOpts"`
}

