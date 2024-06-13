//go:build no_runtime_type_checking

package volume

// Building without runtime type checking enabled, so all the below just return nil

func (v *jsiiProxy_VolumeLabelsList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (v *jsiiProxy_VolumeLabelsList) validateGetParameters(index *float64) error {
	return nil
}

func (v *jsiiProxy_VolumeLabelsList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_VolumeLabelsList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_VolumeLabelsList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_VolumeLabelsList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_VolumeLabelsList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewVolumeLabelsListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

