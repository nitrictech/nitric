package service


type ServiceEndpointSpecPorts struct {
	// The port inside the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#target_port Service#target_port}
	TargetPort *float64 `field:"required" json:"targetPort" yaml:"targetPort"`
	// A random name for the port.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#name Service#name}
	Name *string `field:"optional" json:"name" yaml:"name"`
	// Rrepresents the protocol of a port: `tcp`, `udp` or `sctp`. Defaults to `tcp`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#protocol Service#protocol}
	Protocol *string `field:"optional" json:"protocol" yaml:"protocol"`
	// The port on the swarm hosts.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#published_port Service#published_port}
	PublishedPort *float64 `field:"optional" json:"publishedPort" yaml:"publishedPort"`
	// Represents the mode in which the port is to be published: 'ingress' or 'host'. Defaults to `ingress`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#publish_mode Service#publish_mode}
	PublishMode *string `field:"optional" json:"publishMode" yaml:"publishMode"`
}

