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

func (j *jsiiProxy_Service) validateSetContainerAppEnvironmentIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetCpuParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetEnvParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetImageUriParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetMaxReplicasParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetMemoryParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetMinReplicasParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetRegistryLoginServerParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetRegistryPasswordParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetRegistryUsernameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetResourceGroupNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetStackNameParameters(val *string) error {
	return nil
}

func validateNewServiceParameters(scope constructs.Construct, id *string, config *ServiceConfig) error {
	return nil
}

