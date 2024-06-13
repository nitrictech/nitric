//go:build no_runtime_type_checking

package plugin

// Building without runtime type checking enabled, so all the below just return nil

func (p *jsiiProxy_PluginGrantPermissionsList) validateAllWithMapKeyParameters(mapKeyAttributeName *string) error {
	return nil
}

func (p *jsiiProxy_PluginGrantPermissionsList) validateGetParameters(index *float64) error {
	return nil
}

func (p *jsiiProxy_PluginGrantPermissionsList) validateResolveParameters(_context cdktf.IResolveContext) error {
	return nil
}

func (j *jsiiProxy_PluginGrantPermissionsList) validateSetInternalValueParameters(val interface{}) error {
	return nil
}

func (j *jsiiProxy_PluginGrantPermissionsList) validateSetTerraformAttributeParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_PluginGrantPermissionsList) validateSetTerraformResourceParameters(val cdktf.IInterpolatingParent) error {
	return nil
}

func (j *jsiiProxy_PluginGrantPermissionsList) validateSetWrapsSetParameters(val *bool) error {
	return nil
}

func validateNewPluginGrantPermissionsListParameters(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string, wrapsSet *bool) error {
	return nil
}

