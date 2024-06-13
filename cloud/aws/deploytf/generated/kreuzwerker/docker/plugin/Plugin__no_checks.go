//go:build no_runtime_type_checking

package plugin

// Building without runtime type checking enabled, so all the below just return nil

func (p *jsiiProxy_Plugin) validateAddMoveTargetParameters(moveTarget *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetAnyMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetBooleanAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetBooleanMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetNumberAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetNumberListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetNumberMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetStringAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateGetStringMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateImportFromParameters(id *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateInterpolationForAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateMoveFromIdParameters(id *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateMoveToParameters(moveTarget *string, index interface{}) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateMoveToIdParameters(id *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func (p *jsiiProxy_Plugin) validatePutGrantPermissionsParameters(value interface{}) error {
	return nil
}

func validatePlugin_GenerateConfigForImportParameters(scope constructs.Construct, importToId *string, importFromId *string) error {
	return nil
}

func validatePlugin_IsConstructParameters(x interface{}) error {
	return nil
}

func validatePlugin_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validatePlugin_IsTerraformResourceParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetAliasParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetConnectionParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetCountParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetEnabledParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetEnableTimeoutParameters(val *float64) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetEnvParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetForceDestroyParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetForceDisableParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetGrantAllPermissionsParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetLifecycleParameters(val *cdktf.TerraformResourceLifecycle) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Plugin) validateSetProvisionersParameters(val *[]interface{}) error {
	return nil
}

func validateNewPluginParameters(scope constructs.Construct, id *string, config *PluginConfig) error {
	return nil
}

