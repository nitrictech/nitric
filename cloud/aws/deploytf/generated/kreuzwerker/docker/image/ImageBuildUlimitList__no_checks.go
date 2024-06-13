//go:build no_runtime_type_checking

package image

// Building without runtime type checking enabled, so all the below just return nil

func (i *jsiiProxy_ImageBuildUlimitList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (i *jsiiProxy_ImageBuildUlimitList) validateGetParameters(index *float64) error {
	return nil
}

func (i *jsiiProxy_ImageBuildUlimitList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_ImageBuildUlimitList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_ImageBuildUlimitList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_ImageBuildUlimitList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_ImageBuildUlimitList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewImageBuildUlimitListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

