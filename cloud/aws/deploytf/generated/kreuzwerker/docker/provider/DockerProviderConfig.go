package provider


type DockerProviderConfig struct {
	// Alias name.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#alias DockerProvider#alias}
	Alias *string `field:"optional" json:"alias" yaml:"alias"`
	// PEM-encoded content of Docker host CA certificate.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#ca_material DockerProvider#ca_material}
	CaMaterial *string `field:"optional" json:"caMaterial" yaml:"caMaterial"`
	// PEM-encoded content of Docker client certificate.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#cert_material DockerProvider#cert_material}
	CertMaterial *string `field:"optional" json:"certMaterial" yaml:"certMaterial"`
	// Path to directory with Docker TLS config.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#cert_path DockerProvider#cert_path}
	CertPath *string `field:"optional" json:"certPath" yaml:"certPath"`
	// The Docker daemon address.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#host DockerProvider#host}
	Host *string `field:"optional" json:"host" yaml:"host"`
	// PEM-encoded content of Docker client private key.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#key_material DockerProvider#key_material}
	KeyMaterial *string `field:"optional" json:"keyMaterial" yaml:"keyMaterial"`
	// registry_auth block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#registry_auth DockerProvider#registry_auth}
	RegistryAuth interface{} `field:"optional" json:"registryAuth" yaml:"registryAuth"`
	// Additional SSH option flags to be appended when using `ssh://` protocol.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#ssh_opts DockerProvider#ssh_opts}
	SshOpts *[]*string `field:"optional" json:"sshOpts" yaml:"sshOpts"`
}

