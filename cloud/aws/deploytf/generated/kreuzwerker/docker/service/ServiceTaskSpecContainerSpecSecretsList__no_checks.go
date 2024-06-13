//go:build no_runtime_type_checking

package service

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_ServiceTaskSpecContainerSpecSecretsList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecSecretsList) validateGetParameters(index *float64) error {
	return nil
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecSecretsList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecSecretsList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecSecretsList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecSecretsList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecSecretsList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewServiceTaskSpecContainerSpecSecretsListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

