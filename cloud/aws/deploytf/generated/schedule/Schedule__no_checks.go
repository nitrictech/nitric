//go:build no_runtime_type_checking

package schedule

// Building without runtime type checking enabled, so all the below just return nil

func (s *jsiiProxy_Schedule) validateAddOverrideParameters(path *string, value interface{}) error {
	return nil
}

func (s *jsiiProxy_Schedule) validateAddProviderParameters(provider interface{}) error {
	return nil
}

func (s *jsiiProxy_Schedule) validateGetStringParameters(output *string) error {
	return nil
}

func (s *jsiiProxy_Schedule) validateInterpolationForOutputParameters(moduleOutput *string) error {
	return nil
}

func (s *jsiiProxy_Schedule) validateOverrideLogicalIdParameters(newLogicalId *string) error {
	return nil
}

func validateSchedule_IsConstructParameters(x interface{}) error {
	return nil
}

func validateSchedule_IsTerraformElementParameters(x interface{}) error {
	return nil
}

func (j *jsiiProxy_Schedule) validateSetScheduleExpressionParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Schedule) validateSetScheduleNameParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Schedule) validateSetScheduleTimezoneParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Schedule) validateSetStackIdParameters(val *string) error {
	return nil
}

func (j *jsiiProxy_Schedule) validateSetTargetLambdaArnParameters(val *string) error {
	return nil
}

func validateNewScheduleParameters(scope constructs.Construct, id *string, config *ScheduleConfig) error {
	return nil
}

