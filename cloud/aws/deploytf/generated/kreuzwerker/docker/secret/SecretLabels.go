package secret


type SecretLabels struct {
	// Name of the label.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/secret#label Secret#label}
	Label *string `field:"required" json:"label" yaml:"label"`
	// Value of the label.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/secret#value Secret#value}
	Value *string `field:"required" json:"value" yaml:"value"`
}

