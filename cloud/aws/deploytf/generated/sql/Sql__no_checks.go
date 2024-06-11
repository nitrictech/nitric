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

func (j *jsiiProxy_Sql) validateSetCodebuildRoleArnParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetCreateDatabaseProjectNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetDbNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetImageUriParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetMigrateCommandParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetRdsClusterEndpointParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetRdsClusterPasswordParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetRdsClusterUsernameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetSecurityGroupIdsParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetSubnetIdsParameters(val *[]*string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetVpcIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Sql) validateSetWorkDirParameters(val *string) error {
	return nil
}

func validateNewSqlParameters(scope constructs.Construct, id *string, config *SqlConfig) error {
	return nil
}

