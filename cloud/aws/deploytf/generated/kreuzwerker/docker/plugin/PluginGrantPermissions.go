package plugin


type PluginGrantPermissions struct {
	// The name of the permission.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#name Plugin#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// The value of the permission.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/plugin#value Plugin#value}
	Value *[]*string `field:"required" json:"value" yaml:"value"`
}

