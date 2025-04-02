// stack
package stack

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"stack.Stack",
		reflect.TypeOf((*Stack)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addProvider", GoMethod: "AddProvider"},
			_jsii_.MemberProperty{JsiiProperty: "appIdentityClientIdOutput", GoGetter: "AppIdentityClientIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "appIdentityOutput", GoGetter: "AppIdentityOutput"},
			_jsii_.MemberProperty{JsiiProperty: "cdktfStack", GoGetter: "CdktfStack"},
			_jsii_.MemberProperty{JsiiProperty: "constructNodeMetadata", GoGetter: "ConstructNodeMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "containerAppEnvironmentIdOutput", GoGetter: "ContainerAppEnvironmentIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "containerAppSubnetIdOutput", GoGetter: "ContainerAppSubnetIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "databaseMasterPasswordOutput", GoGetter: "DatabaseMasterPasswordOutput"},
			_jsii_.MemberProperty{JsiiProperty: "databaseServerFqdnOutput", GoGetter: "DatabaseServerFqdnOutput"},
			_jsii_.MemberProperty{JsiiProperty: "databaseServerIdOutput", GoGetter: "DatabaseServerIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "databaseServerNameOutput", GoGetter: "DatabaseServerNameOutput"},
			_jsii_.MemberProperty{JsiiProperty: "dependsOn", GoGetter: "DependsOn"},
			_jsii_.MemberProperty{JsiiProperty: "enableDatabase", GoGetter: "EnableDatabase"},
			_jsii_.MemberProperty{JsiiProperty: "enableKeyvault", GoGetter: "EnableKeyvault"},
			_jsii_.MemberProperty{JsiiProperty: "enableStorage", GoGetter: "EnableStorage"},
			_jsii_.MemberProperty{JsiiProperty: "forEach", GoGetter: "ForEach"},
			_jsii_.MemberProperty{JsiiProperty: "fqn", GoGetter: "Fqn"},
			_jsii_.MemberProperty{JsiiProperty: "friendlyUniqueId", GoGetter: "FriendlyUniqueId"},
			_jsii_.MemberMethod{JsiiMethod: "getString", GoMethod: "GetString"},
			_jsii_.MemberProperty{JsiiProperty: "infrastructureSubnetId", GoGetter: "InfrastructureSubnetId"},
			_jsii_.MemberProperty{JsiiProperty: "infrastructureSubnetIdOutput", GoGetter: "InfrastructureSubnetIdOutput"},
			_jsii_.MemberMethod{JsiiMethod: "interpolationForOutput", GoMethod: "InterpolationForOutput"},
			_jsii_.MemberProperty{JsiiProperty: "keyvaultNameOutput", GoGetter: "KeyvaultNameOutput"},
			_jsii_.MemberProperty{JsiiProperty: "location", GoGetter: "Location"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "providers", GoGetter: "Providers"},
			_jsii_.MemberProperty{JsiiProperty: "rawOverrides", GoGetter: "RawOverrides"},
			_jsii_.MemberProperty{JsiiProperty: "registryLoginServerOutput", GoGetter: "RegistryLoginServerOutput"},
			_jsii_.MemberProperty{JsiiProperty: "registryPasswordOutput", GoGetter: "RegistryPasswordOutput"},
			_jsii_.MemberProperty{JsiiProperty: "registryUsernameOutput", GoGetter: "RegistryUsernameOutput"},
			_jsii_.MemberMethod{JsiiMethod: "resetOverrideLogicalId", GoMethod: "ResetOverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "resourceGroupNameOutput", GoGetter: "ResourceGroupNameOutput"},
			_jsii_.MemberProperty{JsiiProperty: "skipAssetCreationFromLocalModules", GoGetter: "SkipAssetCreationFromLocalModules"},
			_jsii_.MemberProperty{JsiiProperty: "source", GoGetter: "Source"},
			_jsii_.MemberProperty{JsiiProperty: "stackIdOutput", GoGetter: "StackIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "stackName", GoGetter: "StackName"},
			_jsii_.MemberProperty{JsiiProperty: "stackNameOutput", GoGetter: "StackNameOutput"},
			_jsii_.MemberProperty{JsiiProperty: "storageAccountBlobEndpointOutput", GoGetter: "StorageAccountBlobEndpointOutput"},
			_jsii_.MemberProperty{JsiiProperty: "storageAccountConnectionStringOutput", GoGetter: "StorageAccountConnectionStringOutput"},
			_jsii_.MemberProperty{JsiiProperty: "storageAccountIdOutput", GoGetter: "StorageAccountIdOutput"},
			_jsii_.MemberProperty{JsiiProperty: "storageAccountNameOutput", GoGetter: "StorageAccountNameOutput"},
			_jsii_.MemberProperty{JsiiProperty: "storageAccountQueueEndpointOutput", GoGetter: "StorageAccountQueueEndpointOutput"},
			_jsii_.MemberProperty{JsiiProperty: "subscriptionIdOutput", GoGetter: "SubscriptionIdOutput"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeAttributes", GoMethod: "SynthesizeAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeHclAttributes", GoMethod: "SynthesizeHclAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "tags", GoGetter: "Tags"},
			_jsii_.MemberMethod{JsiiMethod: "toHclTerraform", GoMethod: "ToHclTerraform"},
			_jsii_.MemberMethod{JsiiMethod: "toMetadata", GoMethod: "ToMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "toTerraform", GoMethod: "ToTerraform"},
			_jsii_.MemberProperty{JsiiProperty: "version", GoGetter: "Version"},
		},
		func() interface{} {
			j := jsiiProxy_Stack{}
			_jsii_.InitJsiiProxy(&j.Type__cdktfTerraformModule)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"stack.StackConfig",
		reflect.TypeOf((*StackConfig)(nil)).Elem(),
	)
}
