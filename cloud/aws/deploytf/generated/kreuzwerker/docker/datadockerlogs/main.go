package datadockerlogs

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"docker.dataDockerLogs.DataDockerLogs",
		reflect.TypeOf((*DataDockerLogs)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberProperty{JsiiProperty: "cdktfStack", GoGetter: "CdktfStack"},
			_jsii_.MemberProperty{JsiiProperty: "constructNodeMetadata", GoGetter: "ConstructNodeMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "count", GoGetter: "Count"},
			_jsii_.MemberProperty{JsiiProperty: "dependsOn", GoGetter: "DependsOn"},
			_jsii_.MemberProperty{JsiiProperty: "details", GoGetter: "Details"},
			_jsii_.MemberProperty{JsiiProperty: "detailsInput", GoGetter: "DetailsInput"},
			_jsii_.MemberProperty{JsiiProperty: "discardHeaders", GoGetter: "DiscardHeaders"},
			_jsii_.MemberProperty{JsiiProperty: "discardHeadersInput", GoGetter: "DiscardHeadersInput"},
			_jsii_.MemberProperty{JsiiProperty: "follow", GoGetter: "Follow"},
			_jsii_.MemberProperty{JsiiProperty: "followInput", GoGetter: "FollowInput"},
			_jsii_.MemberProperty{JsiiProperty: "forEach", GoGetter: "ForEach"},
			_jsii_.MemberProperty{JsiiProperty: "fqn", GoGetter: "Fqn"},
			_jsii_.MemberProperty{JsiiProperty: "friendlyUniqueId", GoGetter: "FriendlyUniqueId"},
			_jsii_.MemberMethod{JsiiMethod: "getAnyMapAttribute", GoMethod: "GetAnyMapAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getBooleanAttribute", GoMethod: "GetBooleanAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getBooleanMapAttribute", GoMethod: "GetBooleanMapAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getListAttribute", GoMethod: "GetListAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getNumberAttribute", GoMethod: "GetNumberAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getNumberListAttribute", GoMethod: "GetNumberListAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getNumberMapAttribute", GoMethod: "GetNumberMapAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getStringAttribute", GoMethod: "GetStringAttribute"},
			_jsii_.MemberMethod{JsiiMethod: "getStringMapAttribute", GoMethod: "GetStringMapAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "id", GoGetter: "Id"},
			_jsii_.MemberProperty{JsiiProperty: "idInput", GoGetter: "IdInput"},
			_jsii_.MemberMethod{JsiiMethod: "interpolationForAttribute", GoMethod: "InterpolationForAttribute"},
			_jsii_.MemberProperty{JsiiProperty: "lifecycle", GoGetter: "Lifecycle"},
			_jsii_.MemberProperty{JsiiProperty: "logsListString", GoGetter: "LogsListString"},
			_jsii_.MemberProperty{JsiiProperty: "logsListStringEnabled", GoGetter: "LogsListStringEnabled"},
			_jsii_.MemberProperty{JsiiProperty: "logsListStringEnabledInput", GoGetter: "LogsListStringEnabledInput"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "nameInput", GoGetter: "NameInput"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "provider", GoGetter: "Provider"},
			_jsii_.MemberProperty{JsiiProperty: "rawOverrides", GoGetter: "RawOverrides"},
			_jsii_.MemberMethod{JsiiMethod: "resetDetails", GoMethod: "ResetDetails"},
			_jsii_.MemberMethod{JsiiMethod: "resetDiscardHeaders", GoMethod: "ResetDiscardHeaders"},
			_jsii_.MemberMethod{JsiiMethod: "resetFollow", GoMethod: "ResetFollow"},
			_jsii_.MemberMethod{JsiiMethod: "resetId", GoMethod: "ResetId"},
			_jsii_.MemberMethod{JsiiMethod: "resetLogsListStringEnabled", GoMethod: "ResetLogsListStringEnabled"},
			_jsii_.MemberMethod{JsiiMethod: "resetOverrideLogicalId", GoMethod: "ResetOverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "resetShowStderr", GoMethod: "ResetShowStderr"},
			_jsii_.MemberMethod{JsiiMethod: "resetShowStdout", GoMethod: "ResetShowStdout"},
			_jsii_.MemberMethod{JsiiMethod: "resetSince", GoMethod: "ResetSince"},
			_jsii_.MemberMethod{JsiiMethod: "resetTail", GoMethod: "ResetTail"},
			_jsii_.MemberMethod{JsiiMethod: "resetTimestamps", GoMethod: "ResetTimestamps"},
			_jsii_.MemberMethod{JsiiMethod: "resetUntil", GoMethod: "ResetUntil"},
			_jsii_.MemberProperty{JsiiProperty: "showStderr", GoGetter: "ShowStderr"},
			_jsii_.MemberProperty{JsiiProperty: "showStderrInput", GoGetter: "ShowStderrInput"},
			_jsii_.MemberProperty{JsiiProperty: "showStdout", GoGetter: "ShowStdout"},
			_jsii_.MemberProperty{JsiiProperty: "showStdoutInput", GoGetter: "ShowStdoutInput"},
			_jsii_.MemberProperty{JsiiProperty: "since", GoGetter: "Since"},
			_jsii_.MemberProperty{JsiiProperty: "sinceInput", GoGetter: "SinceInput"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeAttributes", GoMethod: "SynthesizeAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeHclAttributes", GoMethod: "SynthesizeHclAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "tail", GoGetter: "Tail"},
			_jsii_.MemberProperty{JsiiProperty: "tailInput", GoGetter: "TailInput"},
			_jsii_.MemberProperty{JsiiProperty: "terraformGeneratorMetadata", GoGetter: "TerraformGeneratorMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "terraformMetaArguments", GoGetter: "TerraformMetaArguments"},
			_jsii_.MemberProperty{JsiiProperty: "terraformResourceType", GoGetter: "TerraformResourceType"},
			_jsii_.MemberProperty{JsiiProperty: "timestamps", GoGetter: "Timestamps"},
			_jsii_.MemberProperty{JsiiProperty: "timestampsInput", GoGetter: "TimestampsInput"},
			_jsii_.MemberMethod{JsiiMethod: "toHclTerraform", GoMethod: "ToHclTerraform"},
			_jsii_.MemberMethod{JsiiMethod: "toMetadata", GoMethod: "ToMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "toTerraform", GoMethod: "ToTerraform"},
			_jsii_.MemberProperty{JsiiProperty: "until", GoGetter: "Until"},
			_jsii_.MemberProperty{JsiiProperty: "untilInput", GoGetter: "UntilInput"},
		},
		func() interface{} {
			j := jsiiProxy_DataDockerLogs{}
			_jsii_.InitJsiiProxy(&j.Type__cdktfTerraformDataSource)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"docker.dataDockerLogs.DataDockerLogsConfig",
		reflect.TypeOf((*DataDockerLogsConfig)(nil)).Elem(),
	)
}
