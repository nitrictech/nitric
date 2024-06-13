package service


type ServiceTaskSpecResourcesReservation struct {
	// generic_resources block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#generic_resources Service#generic_resources}
	GenericResources *ServiceTaskSpecResourcesReservationGenericResources `field:"optional" json:"genericResources" yaml:"genericResources"`
	// The amounf of memory in bytes the container allocates.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#memory_bytes Service#memory_bytes}
	MemoryBytes *float64 `field:"optional" json:"memoryBytes" yaml:"memoryBytes"`
	// CPU shares in units of 1/1e9 (or 10^-9) of the CPU. Should be at least `1000000`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#nano_cpus Service#nano_cpus}
	NanoCpus *float64 `field:"optional" json:"nanoCpus" yaml:"nanoCpus"`
}

