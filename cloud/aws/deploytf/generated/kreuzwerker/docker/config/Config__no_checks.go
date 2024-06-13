//go:build no_runtime_type_checking

package config

// Building without runtime type checking enabled, so all the below just return nil

func (c *jsiiProxy_Config) validateAddMoveTargetParameters(moveTarget *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetAnyMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetBooleanAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetBooleanMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetNumberAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetNumberListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetNumberMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetStringAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateGetStringMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateImportFromParameters(id *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateInterpolationForAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateMoveFromIdParameters(id *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateMoveToParameters(moveTarget *string, index interface{}) error {
	return nil
}

func (c *jsiiProxy_Config) validateMoveToIdParameters(id *string) error {
	return nil
}

func (c *jsiiProxy_Config) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateConfig_GenerateConfigForImportParameters(scope constructs.Construct, importToId *string, importFromId *string) error {
	return nil
}

func validateConfig_IsConstructParameters(x interface{}) error {
	return nil
}

func validateConfig_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validateConfig_IsTerraformResourceParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Config) validateSetConnectionParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Config) validateSetCountParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Config) validateSetDataParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Config) validateSetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Config) validateSetLifecycleParameters(val *cdktf.TerraformResourceLifecycle) error {
	return nil
}

func (j *jsiiProxy_Config) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Config) validateSetProvisionersParameters(val *[]interface{}) error {
	return nil
}

func validateNewConfigParameters(scope constructs.Construct, id *string, config *ConfigConfig) error {
	return nil
}

