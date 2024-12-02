//go:build no_runtime_type_checking

package secret

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_Secret) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (s *jsiiProxy_Secret) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetStringParameters(output *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateSecret_IsConstructParameters(x interface{}) error {
	return nil
}

func validateSecret_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetLocationParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetSecretNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetStackIdParameters(val *string) error {
	return nil
}

func validateNewSecretParameters(scope constructs.Construct, id *string, config *SecretConfig) error {
	return nil
}

