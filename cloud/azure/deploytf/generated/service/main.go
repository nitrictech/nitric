// service
package service

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"service.Service",
		reflect.TypeOf((*Service)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addProvider", GoMethod: "AddProvider"},
			_jsii_.MemberProperty{JsiiProperty: "cdktfStack", GoGetter: "CdktfStack"},
			_jsii_.MemberProperty{JsiiProperty: "clientIdOutput", GoGetter: "ClientIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "clientSecretOutput", GoGetter: "ClientSecretOutput"},
			_jsii_.MemberProperty{JsiiProperty: "constructNodeMetadata", GoGetter: "ConstructNodeMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "containerAppEnvironmentId", GoGetter: "ContainerAppEnvironmentId"},
			_jsii_.MemberProperty{JsiiProperty: "containerAppIdOutput", GoGetter: "ContainerAppIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "cpu", GoGetter: "Cpu"},
			_jsii_.MemberProperty{JsiiProperty: "daprAppIdOutput", GoGetter: "DaprAppIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "dependsOn", GoGetter: "DependsOn"},
			_jsii_.MemberProperty{JsiiProperty: "endpointOutput", GoGetter: "EndpointOutput"},
			_jsii_.MemberProperty{JsiiProperty: "env", GoGetter: "Env"},
			_jsii_.MemberProperty{JsiiProperty: "eventTokenOutput", GoGetter: "EventTokenOutput"},
			_jsii_.MemberProperty{JsiiProperty: "forEach", GoGetter: "ForEach"},
			_jsii_.MemberProperty{JsiiProperty: "fqdnOutput", GoGetter: "FqdnOutput"},
			_jsii_.MemberProperty{JsiiProperty: "fqn", GoGetter: "Fqn"},
			_jsii_.MemberProperty{JsiiProperty: "friendlyUniqueId", GoGetter: "FriendlyUniqueId"},
			_jsii_.MemberMethod{JsiiMethod: "getString", GoMethod: "GetString"},
			_jsii_.MemberProperty{JsiiProperty: "imageUri", GoGetter: "ImageUri"},
			_jsii_.MemberMethod{JsiiMethod: "interpolationForOutput", GoMethod: "InterpolationForOutput"},
			_jsii_.MemberProperty{JsiiProperty: "maxReplicas", GoGetter: "MaxReplicas"},
			_jsii_.MemberProperty{JsiiProperty: "memory", GoGetter: "Memory"},
			_jsii_.MemberProperty{JsiiProperty: "minReplicas", GoGetter: "MinReplicas"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "pollUrlOutput", GoGetter: "PollUrlOutput"},
			_jsii_.MemberProperty{JsiiProperty: "providers", GoGetter: "Providers"},
			_jsii_.MemberProperty{JsiiProperty: "rawOverrides", GoGetter: "RawOverrides"},
			_jsii_.MemberProperty{JsiiProperty: "registryLoginServer", GoGetter: "RegistryLoginServer"},
			_jsii_.MemberProperty{JsiiProperty: "registryPassword", GoGetter: "RegistryPassword"},
			_jsii_.MemberProperty{JsiiProperty: "registryUsername", GoGetter: "RegistryUsername"},
			_jsii_.MemberMethod{JsiiMethod: "resetOverrideLogicalId", GoMethod: "ResetOverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "resourceGroupName", GoGetter: "ResourceGroupName"},
			_jsii_.MemberProperty{JsiiProperty: "servicePrincipalIdOutput", GoGetter: "ServicePrincipalIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "skipAssetCreationFromLocalModules", GoGetter: "SkipAssetCreationFromLocalModules"},
			_jsii_.MemberProperty{JsiiProperty: "source", GoGetter: "Source"},
			_jsii_.MemberProperty{JsiiProperty: "stackId", GoGetter: "StackId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeAttributes", GoMethod: "SynthesizeAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeHclAttributes", GoMethod: "SynthesizeHclAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberProperty{JsiiProperty: "tenantIdOutput", GoGetter: "TenantIdOutput"},
			_jsii_.MemberMethod{JsiiMethod: "toHclTerraform", GoMethod: "ToHclTerraform"},
			_jsii_.MemberMethod{JsiiMethod: "toMetadata", GoMethod: "ToMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "toTerraform", GoMethod: "ToTerraform"},
			_jsii_.MemberProperty{JsiiProperty: "version", GoGetter: "Version"},
		},
		func() interface{} {
			j := jsiiProxy_Service{}
			_jsii_.InitJsiiProxy(&j.Type__cdktfTerraformModule)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"service.ServiceConfig",
		reflect.TypeOf((*ServiceConfig)(nil)).Elem(),
	)
}
