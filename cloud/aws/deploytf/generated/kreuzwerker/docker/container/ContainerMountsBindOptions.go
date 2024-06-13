package container


type ContainerMountsBindOptions struct {
	// A propagation mode with the value.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#propagation Container#propagation}
	Propagation *string `field:"optional" json:"propagation" yaml:"propagation"`
}

