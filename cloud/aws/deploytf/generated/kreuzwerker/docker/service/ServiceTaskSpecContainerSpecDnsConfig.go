package service


type ServiceTaskSpecContainerSpecDnsConfig struct {
	// The IP addresses of the name servers.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#nameservers Service#nameservers}
	Nameservers *[]*string `field:"required" json:"nameservers" yaml:"nameservers"`
	// A list of internal resolver variables to be modified (e.g., `debug`, `ndots:3`, etc.).
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#options Service#options}
	Options *[]*string `field:"optional" json:"options" yaml:"options"`
	// A search list for host-name lookup.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#search Service#search}
	Search *[]*string `field:"optional" json:"search" yaml:"search"`
}

