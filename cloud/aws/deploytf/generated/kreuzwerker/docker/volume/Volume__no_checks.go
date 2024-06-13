//go:build no_runtime_type_checking

package volume

// Building without runtime type checking enabled, so all the below just return nil

func (v *jsiiProxy_Volume) validateAddMoveTargetParameters(moveTarget *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetAnyMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetBooleanAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetBooleanMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetNumberAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetNumberListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetNumberMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetStringAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateGetStringMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateImportFromParameters(id *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateInterpolationForAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateMoveFromIdParameters(id *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateMoveToParameters(moveTarget *string, index interface{}) error {
	return nil
}

func (v *jsiiProxy_Volume) validateMoveToIdParameters(id *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func (v *jsiiProxy_Volume) validatePutLabelsParameters(value interface{}) error {
	return nil
}

func validateVolume_GenerateConfigForImportParameters(scope constructs.Construct, importToId *string, importFromId *string) error {
	return nil
}

func validateVolume_IsConstructParameters(x interface{}) error {
	return nil
}

func validateVolume_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validateVolume_IsTerraformResourceParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetConnectionParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetCountParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetDriverParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetDriverOptsParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetLifecycleParameters(val *cdktf.TerraformResourceLifecycle) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Volume) validateSetProvisionersParameters(val *[]interface{}) error {
	return nil
}

func validateNewVolumeParameters(scope constructs.Construct, id *string, config *VolumeConfig) error {
	return nil
}

