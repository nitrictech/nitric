package container


type ContainerVolumes struct {
	// The path in the container where the volume will be mounted.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#container_path Container#container_path}
	ContainerPath *string `field:"optional" json:"containerPath" yaml:"containerPath"`
	// The container where the volume is coming from.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#from_container Container#from_container}
	FromContainer *string `field:"optional" json:"fromContainer" yaml:"fromContainer"`
	// The path on the host where the volume is coming from.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#host_path Container#host_path}
	HostPath *string `field:"optional" json:"hostPath" yaml:"hostPath"`
	// If `true`, this volume will be readonly. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#read_only Container#read_only}
	ReadOnly interface{} `field:"optional" json:"readOnly" yaml:"readOnly"`
	// The name of the docker volume which should be mounted.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#volume_name Container#volume_name}
	VolumeName *string `field:"optional" json:"volumeName" yaml:"volumeName"`
}

