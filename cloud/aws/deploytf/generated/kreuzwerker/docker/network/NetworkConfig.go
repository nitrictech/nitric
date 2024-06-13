package network

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type NetworkConfig struct {
	// Experimental.
	Connection interface{} `field:"optional" json:"connection" yaml:"connection"`
	// Experimental.
	Count interface{} `field:"optional" json:"count" yaml:"count"`
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Lifecycle *cdktf.TerraformResourceLifecycle `field:"optional" json:"lifecycle" yaml:"lifecycle"`
	// Experimental.
	Provider cdktf.TerraformProvider `field:"optional" json:"provider" yaml:"provider"`
	// Experimental.
	Provisioners *[]interface{} `field:"optional" json:"provisioners" yaml:"provisioners"`
	// The name of the Docker network.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#name Network#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// Enable manual container attachment to the network.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#attachable Network#attachable}
	Attachable interface{} `field:"optional" json:"attachable" yaml:"attachable"`
	// Requests daemon to check for networks with same name.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#check_duplicate Network#check_duplicate}
	CheckDuplicate interface{} `field:"optional" json:"checkDuplicate" yaml:"checkDuplicate"`
	// The driver of the Docker network. Possible values are `bridge`, `host`, `overlay`, `macvlan`. See [network docs](https://docs.docker.com/network/#network-drivers) for more details.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#driver Network#driver}
	Driver *string `field:"optional" json:"driver" yaml:"driver"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#id Network#id}.
	//
	// Please be aware that the id field is automatically added to all resources in Terraform providers using a Terraform provider SDK version below 2.
	// If you experience problems setting this value it might not be settable. Please take a look at the provider documentation to ensure it should be settable.
	Id *string `field:"optional" json:"id" yaml:"id"`
	// Create swarm routing-mesh network. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#ingress Network#ingress}
	Ingress interface{} `field:"optional" json:"ingress" yaml:"ingress"`
	// Whether the network is internal.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#internal Network#internal}
	Internal interface{} `field:"optional" json:"internal" yaml:"internal"`
	// ipam_config block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#ipam_config Network#ipam_config}
	IpamConfig interface{} `field:"optional" json:"ipamConfig" yaml:"ipamConfig"`
	// Driver used by the custom IP scheme of the network. Defaults to `default`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#ipam_driver Network#ipam_driver}
	IpamDriver *string `field:"optional" json:"ipamDriver" yaml:"ipamDriver"`
	// Provide explicit options to the IPAM driver.
	//
	// Valid options vary with `ipam_driver` and refer to that driver's documentation for more details.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#ipam_options Network#ipam_options}
	IpamOptions *map[string]*string `field:"optional" json:"ipamOptions" yaml:"ipamOptions"`
	// Enable IPv6 networking. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#ipv6 Network#ipv6}
	Ipv6 interface{} `field:"optional" json:"ipv6" yaml:"ipv6"`
	// labels block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#labels Network#labels}
	Labels interface{} `field:"optional" json:"labels" yaml:"labels"`
	// Only available with bridge networks. See [bridge options docs](https://docs.docker.com/engine/reference/commandline/network_create/#bridge-driver-options) for more details.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network#options Network#options}
	Options *map[string]*string `field:"optional" json:"options" yaml:"options"`
}

