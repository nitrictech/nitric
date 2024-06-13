package service


type ServiceTaskSpecLogDriver struct {
	// The logging driver to use.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#name Service#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// The options for the logging driver.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#options Service#options}
	Options *map[string]*string `field:"optional" json:"options" yaml:"options"`
}

