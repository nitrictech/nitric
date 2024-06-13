//go:build no_runtime_type_checking

package image

// Building without runtime type checking enabled, so all the below just return nil

func (i *jsiiProxy_ImageBuildAuthConfigList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (i *jsiiProxy_ImageBuildAuthConfigList) validateGetParameters(index *float64) error {
	return nil
}

func (i *jsiiProxy_ImageBuildAuthConfigList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_ImageBuildAuthConfigList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_ImageBuildAuthConfigList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_ImageBuildAuthConfigList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_ImageBuildAuthConfigList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewImageBuildAuthConfigListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

