package provider


type DockerProviderRegistryAuth struct {
	// Address of the registry.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#address DockerProvider#address}
	Address *string `field:"required" json:"address" yaml:"address"`
	// Setting this to `true` will tell the provider that this registry does not need authentication.
	//
	// Due to the docker internals, the provider will use dummy credentials (see https://github.com/kreuzwerker/terraform-provider-docker/issues/470 for more information). Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#auth_disabled DockerProvider#auth_disabled}
	AuthDisabled interface{} `field:"optional" json:"authDisabled" yaml:"authDisabled"`
	// Path to docker json file for registry auth.
	//
	// Defaults to `~/.docker/config.json`. If `DOCKER_CONFIG` is set, the value of `DOCKER_CONFIG` is used as the path. `config_file` has predencen over all other options.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#config_file DockerProvider#config_file}
	ConfigFile *string `field:"optional" json:"configFile" yaml:"configFile"`
	// Plain content of the docker json file for registry auth. `config_file_content` has precedence over username/password.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#config_file_content DockerProvider#config_file_content}
	ConfigFileContent *string `field:"optional" json:"configFileContent" yaml:"configFileContent"`
	// Password for the registry. Defaults to `DOCKER_REGISTRY_PASS` env variable if set.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#password DockerProvider#password}
	Password *string `field:"optional" json:"password" yaml:"password"`
	// Username for the registry. Defaults to `DOCKER_REGISTRY_USER` env variable if set.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs#username DockerProvider#username}
	Username *string `field:"optional" json:"username" yaml:"username"`
}

