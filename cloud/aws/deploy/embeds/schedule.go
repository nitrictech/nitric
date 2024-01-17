package embeds

import (
	_ "embed"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

//go:embed scheduler-execution-permissions.json
var schedule_PermissionsTemplate string

func GetSchedulePermissionDocument(targetArn pulumi.StringInput) pulumi.StringOutput {
	return pulumi.Sprintf(schedule_PermissionsTemplate, targetArn)
}

//go:embed scheduler-execution-role.json
var schedule_RoleTemplate string

func GetScheduleRoleDocument() pulumi.StringOutput {
	return pulumi.Sprintf(schedule_RoleTemplate)
}

//go:embed scheduler-input.json
var schedule_InputTemplate string

func GetScheduleInputDocument(scheduleName pulumi.StringInput) pulumi.StringOutput {
	return pulumi.Sprintf(schedule_InputTemplate, scheduleName)
}
