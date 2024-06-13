package network


type NetworkIpamConfig struct {
	// Auxiliary IPv4 or IPv6 addresses used by Network driver.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#aux_address Network#aux_address}
	AuxAddress *map[string]*string `field:"optional" json:"auxAddress" yaml:"auxAddress"`
	// The IP address of the gateway.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#gateway Network#gateway}
	Gateway *string `field:"optional" json:"gateway" yaml:"gateway"`
	// The ip range in CIDR form.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#ip_range Network#ip_range}
	IpRange *string `field:"optional" json:"ipRange" yaml:"ipRange"`
	// The subnet in CIDR form.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#subnet Network#subnet}
	Subnet *string `field:"optional" json:"subnet" yaml:"subnet"`
}

