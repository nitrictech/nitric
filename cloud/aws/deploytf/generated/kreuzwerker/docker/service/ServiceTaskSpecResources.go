package service


type ServiceTaskSpecResources struct {
	// limits block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#limits Service#limits}
	Limits *ServiceTaskSpecResourcesLimits `field:"optional" json:"limits" yaml:"limits"`
	// reservation block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#reservation Service#reservation}
	Reservation *ServiceTaskSpecResourcesReservation `field:"optional" json:"reservation" yaml:"reservation"`
}

