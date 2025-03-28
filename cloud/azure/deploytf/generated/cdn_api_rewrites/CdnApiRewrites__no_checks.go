//go:build no_runtime_type_checking

package cdn_api_rewrites

// Building without runtime type checking enabled, so all the below just return nil

func (c *jsiiProxy_CdnApiRewrites) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (c *jsiiProxy_CdnApiRewrites) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (c *jsiiProxy_CdnApiRewrites) validateGetStringParameters(output *string) error {
	return nil
}

func (c *jsiiProxy_CdnApiRewrites) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (c *jsiiProxy_CdnApiRewrites) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateCdnApiRewrites_IsConstructParameters(x interface{}) error {
	return nil
}

func validateCdnApiRewrites_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_CdnApiRewrites) validateSetApiHostNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_CdnApiRewrites) validateSetCdnFrontdoorProfileIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_CdnApiRewrites) validateSetCdnFrontdoorRuleSetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_CdnApiRewrites) validateSetNameParameters(val *string) error {
	return nil
}

func validateNewCdnApiRewritesParameters(scope constructs.Construct, id *string, config *CdnApiRewritesConfig) error {
	return nil
}

