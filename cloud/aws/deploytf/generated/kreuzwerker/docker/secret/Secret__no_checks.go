//go:build no_runtime_type_checking

package secret

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_Secret) validateAddMoveTargetParameters(moveTarget *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetAnyMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetBooleanAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetBooleanMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetNumberAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetNumberListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetNumberMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetStringAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateGetStringMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateImportFromParameters(id *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateInterpolationForAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateMoveFromIdParameters(id *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateMoveToParameters(moveTarget *string, index interface{}) error {
	return nil
}

func (s *jsiiProxy_Secret) validateMoveToIdParameters(id *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func (s *jsiiProxy_Secret) validatePutLabelsParameters(value interface{}) error {
	return nil
}

func validateSecret_GenerateConfigForImportParameters(scope constructs.Construct, importToId *string, importFromId *string) error {
	return nil
}

func validateSecret_IsConstructParameters(x interface{}) error {
	return nil
}

func validateSecret_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validateSecret_IsTerraformResourceParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetConnectionParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetCountParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetDataParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetLifecycleParameters(val *cdktf.TerraformResourceLifecycle) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Secret) validateSetProvisionersParameters(val *[]interface{}) error {
	return nil
}

func validateNewSecretParameters(scope constructs.Construct, id *string, config *SecretConfig) error {
	return nil
}

