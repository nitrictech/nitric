//go:build no_runtime_type_checking

package provider

// Building without runtime type checking enabled, so all the below just return nil

func (d *jsiiProxy_DockerProvider) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (d *jsiiProxy_DockerProvider) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateDockerProvider_GenerateConfigForImportParameters(scope constructs.Construct, importToId *string, importFromId *string) error {
	return nil
}

func validateDockerProvider_IsConstructParameters(x interface{}) error {
	return nil
}

func validateDockerProvider_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validateDockerProvider_IsTerraformProviderParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_DockerProvider) validateSetRegistryAuthParameters(val interface{}) error {
	return nil
}

func validateNewDockerProviderParameters(scope constructs.Construct, id *string, config *DockerProviderConfig) error {
	return nil
}

