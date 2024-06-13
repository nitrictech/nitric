//go:build no_runtime_type_checking

package service

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsList) validateGetParameters(index *float64) error {
	return nil
}

func (s *jsiiProxy_ServiceTaskSpecContainerSpecConfigsList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_ServiceTaskSpecContainerSpecConfigsList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewServiceTaskSpecContainerSpecConfigsListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

