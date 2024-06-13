package service


type ServiceTaskSpecContainerSpecPrivileges struct {
	// credential_spec block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#credential_spec Service#credential_spec}
	CredentialSpec *ServiceTaskSpecContainerSpecPrivilegesCredentialSpec `field:"optional" json:"credentialSpec" yaml:"credentialSpec"`
	// se_linux_context block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#se_linux_context Service#se_linux_context}
	SeLinuxContext *ServiceTaskSpecContainerSpecPrivilegesSeLinuxContext `field:"optional" json:"seLinuxContext" yaml:"seLinuxContext"`
}

