package container


type ContainerDevices struct {
	// The path on the host where the device is located.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#host_path Container#host_path}
	HostPath *string `field:"required" json:"hostPath" yaml:"hostPath"`
	// The path in the container where the device will be bound.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#container_path Container#container_path}
	ContainerPath *string `field:"optional" json:"containerPath" yaml:"containerPath"`
	// The cgroup permissions given to the container to access the device. Defaults to `rwm`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#permissions Container#permissions}
	Permissions *string `field:"optional" json:"permissions" yaml:"permissions"`
}

