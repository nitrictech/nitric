package service


type ServiceAuth struct {
	// The address of the server for the authentication.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#server_address Service#server_address}
	ServerAddress *string `field:"required" json:"serverAddress" yaml:"serverAddress"`
	// The password.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#password Service#password}
	Password *string `field:"optional" json:"password" yaml:"password"`
	// The username.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#username Service#username}
	Username *string `field:"optional" json:"username" yaml:"username"`
}

