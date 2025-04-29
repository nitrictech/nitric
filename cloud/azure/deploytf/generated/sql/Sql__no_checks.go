//go:build no_runtime_type_checking

package sql

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_Sql) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (s *jsiiProxy_Sql) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (s *jsiiProxy_Sql) validateGetStringParameters(output *string) error {
	return nil
}

func (s *jsiiProxy_Sql) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (s *jsiiProxy_Sql) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateSql_IsConstructParameters(x interface{}) error {
	return nil
}

func validateSql_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetDatabaseMasterPasswordParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetDatabaseServerFqdnParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetImageRegistryPasswordParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetImageRegistryServerParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetImageRegistryUsernameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetLocationParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetMigrationContainerSubnetIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetMigrationImageUriParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetResourceGroupNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetServerIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetStackIdParameters(val *string) error {
	return nil
}

func validateNewSqlParameters(scope constructs.Construct, id *string, config *SqlConfig) error {
	return nil
}

