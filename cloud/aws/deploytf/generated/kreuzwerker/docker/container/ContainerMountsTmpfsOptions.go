package container


type ContainerMountsTmpfsOptions struct {
	// The permission mode for the tmpfs mount in an integer.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#mode Container#mode}
	Mode *float64 `field:"optional" json:"mode" yaml:"mode"`
	// The size for the tmpfs mount in bytes.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#size_bytes Container#size_bytes}
	SizeBytes *float64 `field:"optional" json:"sizeBytes" yaml:"sizeBytes"`
}

