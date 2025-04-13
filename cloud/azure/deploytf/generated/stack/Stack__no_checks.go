//go:build no_runtime_type_checking

package stack

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_Stack) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (s *jsiiProxy_Stack) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (s *jsiiProxy_Stack) validateGetStringParameters(output *string) error {
	return nil
}

func (s *jsiiProxy_Stack) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (s *jsiiProxy_Stack) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateStack_IsConstructParameters(x interface{}) error {
	return nil
}

func validateStack_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Stack) validateSetLocationParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Stack) validateSetPrivateEndpointsParameters(val *bool) error {
	return nil
}

func (j *jsiiProxy_Stack) validateSetResourceGroupNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Stack) validateSetStackNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Stack) validateSetSubnetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Stack) validateSetTagsParameters(val *map[string]*string) error {
	return nil
}

func (j *jsiiProxy_Stack) validateSetVnetNameParameters(val *string) error {
	return nil
}

func validateNewStackParameters(scope constructs.Construct, id *string, config *StackConfig) error {
	return nil
}

