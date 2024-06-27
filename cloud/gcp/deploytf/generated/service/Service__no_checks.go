//go:build no_runtime_type_checking

package service

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_Service) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (s *jsiiProxy_Service) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetStringParameters(output *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateService_IsConstructParameters(x interface{}) error {
	return nil
}

func validateService_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetBaseComputeRoleParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetEnvironmentParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetImageParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetProjectIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetRegionParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetServiceNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetStackIdParameters(val *string) error {
	return nil
}

func validateNewServiceParameters(scope constructs.Construct, id *string, config *ServiceConfig) error {
	return nil
}

