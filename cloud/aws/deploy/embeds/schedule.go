// Copyright 2021 Nitric Technologies Pty Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package embeds

import (
	_ "embed"
	"fmt"

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

func GetScheduleInputDocumentString(scheduleName string) string {
	return fmt.Sprintf(schedule_InputTemplate, scheduleName)
}
