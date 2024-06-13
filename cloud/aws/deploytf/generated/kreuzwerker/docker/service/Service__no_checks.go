//go:build no_runtime_type_checking

package service

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_Service) validateAddMoveTargetParameters(moveTarget *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetAnyMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetBooleanAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetBooleanMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetNumberAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetNumberListAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetNumberMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetStringAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateGetStringMapAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateImportFromParameters(id *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateInterpolationForAttributeParameters(terraformAttribute *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateMoveFromIdParameters(id *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateMoveToParameters(moveTarget *string, index interface{}) error {
	return nil
}

func (s *jsiiProxy_Service) validateMoveToIdParameters(id *string) error {
	return nil
}

func (s *jsiiProxy_Service) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutAuthParameters(value *ServiceAuth) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutConvergeConfigParameters(value *ServiceConvergeConfig) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutEndpointSpecParameters(value *ServiceEndpointSpec) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutLabelsParameters(value interface{}) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutModeParameters(value *ServiceMode) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutRollbackConfigParameters(value *ServiceRollbackConfig) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutTaskSpecParameters(value *ServiceTaskSpec) error {
	return nil
}

func (s *jsiiProxy_Service) validatePutUpdateConfigParameters(value *ServiceUpdateConfig) error {
	return nil
}

func validateService_GenerateConfigForImportParameters(scope constructs.Construct, importToId *string, importFromId *string) error {
	return nil
}

func validateService_IsConstructParameters(x interface{}) error {
	return nil
}

func validateService_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func validateService_IsTerraformResourceParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetConnectionParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetCountParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetLifecycleParameters(val *cdktf.TerraformResourceLifecycle) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Service) validateSetProvisionersParameters(val *[]interface{}) error {
	return nil
}

func validateNewServiceParameters(scope constructs.Construct, id *string, config *ServiceConfig) error {
	return nil
}

