package service


type ServiceTaskSpecContainerSpec struct {
	// The image name to use for the containers of the service, like `nginx:1.17.6`. Also use the data-source or resource of `docker_image` with the `repo_digest` or `docker_registry_image` with the `name` attribute for this, as shown in the examples.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#image Service#image}
	Image *string `field:"required" json:"image" yaml:"image"`
	// Arguments to the command.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#args Service#args}
	Args *[]*string `field:"optional" json:"args" yaml:"args"`
	// The command/entrypoint to be run in the image.
	//
	// According to the [docker cli](https://github.com/docker/cli/blob/v20.10.7/cli/command/service/opts.go#L705) the override of the entrypoint is also passed to the `command` property and there is no `entrypoint` attribute in the `ContainerSpec` of the service.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#command Service#command}
	Command *[]*string `field:"optional" json:"command" yaml:"command"`
	// configs block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#configs Service#configs}
	Configs interface{} `field:"optional" json:"configs" yaml:"configs"`
	// The working directory for commands to run in.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#dir Service#dir}
	Dir *string `field:"optional" json:"dir" yaml:"dir"`
	// dns_config block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#dns_config Service#dns_config}
	DnsConfig *ServiceTaskSpecContainerSpecDnsConfig `field:"optional" json:"dnsConfig" yaml:"dnsConfig"`
	// A list of environment variables in the form VAR="value".
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#env Service#env}
	Env *map[string]*string `field:"optional" json:"env" yaml:"env"`
	// A list of additional groups that the container process will run as.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#groups Service#groups}
	Groups *[]*string `field:"optional" json:"groups" yaml:"groups"`
	// healthcheck block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#healthcheck Service#healthcheck}
	Healthcheck *ServiceTaskSpecContainerSpecHealthcheck `field:"optional" json:"healthcheck" yaml:"healthcheck"`
	// The hostname to use for the container, as a valid RFC 1123 hostname.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#hostname Service#hostname}
	Hostname *string `field:"optional" json:"hostname" yaml:"hostname"`
	// hosts block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#hosts Service#hosts}
	Hosts interface{} `field:"optional" json:"hosts" yaml:"hosts"`
	// Isolation technology of the containers running the service. (Windows only). Defaults to `default`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#isolation Service#isolation}
	Isolation *string `field:"optional" json:"isolation" yaml:"isolation"`
	// labels block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#labels Service#labels}
	Labels interface{} `field:"optional" json:"labels" yaml:"labels"`
	// mounts block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#mounts Service#mounts}
	Mounts interface{} `field:"optional" json:"mounts" yaml:"mounts"`
	// privileges block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#privileges Service#privileges}
	Privileges *ServiceTaskSpecContainerSpecPrivileges `field:"optional" json:"privileges" yaml:"privileges"`
	// Mount the container's root filesystem as read only.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#read_only Service#read_only}
	ReadOnly interface{} `field:"optional" json:"readOnly" yaml:"readOnly"`
	// secrets block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#secrets Service#secrets}
	Secrets interface{} `field:"optional" json:"secrets" yaml:"secrets"`
	// Amount of time to wait for the container to terminate before forcefully removing it (ms|s|m|h).
	//
	// If not specified or '0s' the destroy will not check if all tasks/containers of the service terminate.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#stop_grace_period Service#stop_grace_period}
	StopGracePeriod *string `field:"optional" json:"stopGracePeriod" yaml:"stopGracePeriod"`
	// Signal to stop the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#stop_signal Service#stop_signal}
	StopSignal *string `field:"optional" json:"stopSignal" yaml:"stopSignal"`
	// Sysctls config (Linux only).
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#sysctl Service#sysctl}
	Sysctl *map[string]*string `field:"optional" json:"sysctl" yaml:"sysctl"`
	// The user inside the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#user Service#user}
	User *string `field:"optional" json:"user" yaml:"user"`
}

