//go:build no_runtime_type_checking

package website

// Building without runtime type checking enabled, so all the below just return nil

func (w *jsiiProxy_Website) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (w *jsiiProxy_Website) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (w *jsiiProxy_Website) validateGetStringParameters(output *string) error {
	return nil
}

func (w *jsiiProxy_Website) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (w *jsiiProxy_Website) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateWebsite_IsConstructParameters(x interface{}) error {
	return nil
}

func validateWebsite_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Website) validateSetBasePathParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Website) validateSetLocalDirectoryParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Website) validateSetLocationParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Website) validateSetResourceGroupNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Website) validateSetStackNameParameters(val *string) error {
	return nil
}

func validateNewWebsiteParameters(scope constructs.Construct, id *string, config *WebsiteConfig) error {
	return nil
}

