package service


type ServiceTaskSpecContainerSpecSecrets struct {
	// Represents the final filename in the filesystem.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#file_name Service#file_name}
	FileName *string `field:"required" json:"fileName" yaml:"fileName"`
	// ID of the specific secret that we're referencing.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#secret_id Service#secret_id}
	SecretId *string `field:"required" json:"secretId" yaml:"secretId"`
	// Represents the file GID. Defaults to `0`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#file_gid Service#file_gid}
	FileGid *string `field:"optional" json:"fileGid" yaml:"fileGid"`
	// Represents represents the FileMode of the file. Defaults to `0o444`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#file_mode Service#file_mode}
	FileMode *float64 `field:"optional" json:"fileMode" yaml:"fileMode"`
	// Represents the file UID. Defaults to `0`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#file_uid Service#file_uid}
	FileUid *string `field:"optional" json:"fileUid" yaml:"fileUid"`
	// Name of the secret that this references, but this is just provided for lookup/display purposes.
	//
	// The config in the reference will be identified by its ID
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#secret_name Service#secret_name}
	SecretName *string `field:"optional" json:"secretName" yaml:"secretName"`
}

