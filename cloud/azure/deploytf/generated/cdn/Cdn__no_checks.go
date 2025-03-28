//go:build no_runtime_type_checking

package cdn

// Building without runtime type checking enabled, so all the below just return nil

func (c *jsiiProxy_Cdn) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (c *jsiiProxy_Cdn) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (c *jsiiProxy_Cdn) validateGetStringParameters(output *string) error {
	return nil
}

func (c *jsiiProxy_Cdn) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (c *jsiiProxy_Cdn) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateCdn_IsConstructParameters(x interface{}) error {
	return nil
}

func validateCdn_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Cdn) validateSetCustomDomainHostNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Cdn) validateSetDomainNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Cdn) validateSetPrimaryWebHostParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Cdn) validateSetResourceGroupNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Cdn) validateSetStackNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Cdn) validateSetZoneNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Cdn) validateSetZoneResourceGroupNameParameters(val *string) error {
	return nil
}

func validateNewCdnParameters(scope constructs.Construct, id *string, config *CdnConfig) error {
	return nil
}

