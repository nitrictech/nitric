//go:build no_runtime_type_checking

package policy

// Building without runtime type checking enabled, so all the below just return nil

func (p *jsiiProxy_Policy) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (p *jsiiProxy_Policy) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (p *jsiiProxy_Policy) validateGetStringParameters(output *string) error {
	return nil
}

func (p *jsiiProxy_Policy) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (p *jsiiProxy_Policy) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validatePolicy_IsConstructParameters(x interface{}) error {
	return nil
}

func validatePolicy_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Policy) validateSetRoleDefinitionIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Policy) validateSetScopeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Policy) validateSetServicePrincipalIdParameters(val *string) error {
	return nil
}

func validateNewPolicyParameters(scope constructs.Construct, id *string, config *PolicyConfig) error {
	return nil
}

