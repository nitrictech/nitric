package container


type ContainerMounts struct {
	// Container path.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#target Container#target}
	Target *string `field:"required" json:"target" yaml:"target"`
	// The mount type.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#type Container#type}
	Type *string `field:"required" json:"type" yaml:"type"`
	// bind_options block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#bind_options Container#bind_options}
	BindOptions *ContainerMountsBindOptions `field:"optional" json:"bindOptions" yaml:"bindOptions"`
	// Whether the mount should be read-only.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#read_only Container#read_only}
	ReadOnly interface{} `field:"optional" json:"readOnly" yaml:"readOnly"`
	// Mount source (e.g. a volume name, a host path).
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#source Container#source}
	Source *string `field:"optional" json:"source" yaml:"source"`
	// tmpfs_options block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#tmpfs_options Container#tmpfs_options}
	TmpfsOptions *ContainerMountsTmpfsOptions `field:"optional" json:"tmpfsOptions" yaml:"tmpfsOptions"`
	// volume_options block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#volume_options Container#volume_options}
	VolumeOptions *ContainerMountsVolumeOptions `field:"optional" json:"volumeOptions" yaml:"volumeOptions"`
}

