package service


type ServiceTaskSpecContainerSpecMountsVolumeOptions struct {
	// Name of the driver to use to create the volume.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#driver_name Service#driver_name}
	DriverName *string `field:"optional" json:"driverName" yaml:"driverName"`
	// key/value map of driver specific options.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#driver_options Service#driver_options}
	DriverOptions *map[string]*string `field:"optional" json:"driverOptions" yaml:"driverOptions"`
	// labels block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#labels Service#labels}
	Labels interface{} `field:"optional" json:"labels" yaml:"labels"`
	// Populate volume with data from the target.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#no_copy Service#no_copy}
	NoCopy interface{} `field:"optional" json:"noCopy" yaml:"noCopy"`
}

