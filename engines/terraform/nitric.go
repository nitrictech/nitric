package terraform

type NitricVariables struct {
	Name *string `json:"name"`
}

type NitricServiceVariables struct {
	NitricVariables `json:",inline"`
	ImageId         *string                           `json:"image_id"`
	Env             interface{}                       `json:"env"`
	Identities      *map[string]interface{}           `json:"identities"`
	Schedules       *map[string]NitricServiceSchedule `json:"schedules,omitempty"`
	StackId         *string                           `json:"stack_id"`
}

type NitricServiceSchedule struct {
	CronExpression *string `json:"cron_expression"`
	Path           *string `json:"path"`
}

type NitricOutputs struct {
	Id *string `json:"id"`
}

type NitricServiceOutputs struct {
	NitricOutputs `json:",inline"`
	HttpEndpoint  *string `json:"http_endpoint"`
}
