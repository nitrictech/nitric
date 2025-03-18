//go:build no_runtime_type_checking

package domain

// Building without runtime type checking enabled, so all the below just return nil

func (d *jsiiProxy_Domain) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (d *jsiiProxy_Domain) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (d *jsiiProxy_Domain) validateGetStringParameters(output *string) error {
	return nil
}

func (d *jsiiProxy_Domain) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (d *jsiiProxy_Domain) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateDomain_IsConstructParameters(x interface{}) error {
	return nil
}

func validateDomain_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Domain) validateSetApiIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Domain) validateSetApiStageNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Domain) validateSetDomainNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Domain) validateSetStackIdParameters(val *string) error {
	return nil
}

func validateNewDomainParameters(scope constructs.Construct, id *string, config *DomainConfig) error {
	return nil
}

