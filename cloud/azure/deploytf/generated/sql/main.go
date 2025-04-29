// sql
package sql

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"sql.Sql",
		reflect.TypeOf((*Sql)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberMethod{JsiiMethod: "addProvider", GoMethod: "AddProvider"},
			_jsii_.MemberProperty{JsiiProperty: "cdktfStack", GoGetter: "CdktfStack"},
			_jsii_.MemberProperty{JsiiProperty: "constructNodeMetadata", GoGetter: "ConstructNodeMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "databaseMasterPassword", GoGetter: "DatabaseMasterPassword"},
			_jsii_.MemberProperty{JsiiProperty: "databaseServerFqdn", GoGetter: "DatabaseServerFqdn"},
			_jsii_.MemberProperty{JsiiProperty: "dependsOn", GoGetter: "DependsOn"},
			_jsii_.MemberProperty{JsiiProperty: "forEach", GoGetter: "ForEach"},
			_jsii_.MemberProperty{JsiiProperty: "fqn", GoGetter: "Fqn"},
			_jsii_.MemberProperty{JsiiProperty: "friendlyUniqueId", GoGetter: "FriendlyUniqueId"},
			_jsii_.MemberMethod{JsiiMethod: "getString", GoMethod: "GetString"},
			_jsii_.MemberProperty{JsiiProperty: "imageRegistryPassword", GoGetter: "ImageRegistryPassword"},
			_jsii_.MemberProperty{JsiiProperty: "imageRegistryServer", GoGetter: "ImageRegistryServer"},
			_jsii_.MemberProperty{JsiiProperty: "imageRegistryUsername", GoGetter: "ImageRegistryUsername"},
			_jsii_.MemberMethod{JsiiMethod: "interpolationForOutput", GoMethod: "InterpolationForOutput"},
			_jsii_.MemberProperty{JsiiProperty: "location", GoGetter: "Location"},
			_jsii_.MemberProperty{JsiiProperty: "migrationContainerSubnetId", GoGetter: "MigrationContainerSubnetId"},
			_jsii_.MemberProperty{JsiiProperty: "migrationImageUri", GoGetter: "MigrationImageUri"},
			_jsii_.MemberProperty{JsiiProperty: "name", GoGetter: "Name"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "providers", GoGetter: "Providers"},
			_jsii_.MemberProperty{JsiiProperty: "rawOverrides", GoGetter: "RawOverrides"},
			_jsii_.MemberMethod{JsiiMethod: "resetOverrideLogicalId", GoMethod: "ResetOverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "resourceGroupName", GoGetter: "ResourceGroupName"},
			_jsii_.MemberProperty{JsiiProperty: "serverId", GoGetter: "ServerId"},
			_jsii_.MemberProperty{JsiiProperty: "skipAssetCreationFromLocalModules", GoGetter: "SkipAssetCreationFromLocalModules"},
			_jsii_.MemberProperty{JsiiProperty: "source", GoGetter: "Source"},
			_jsii_.MemberProperty{JsiiProperty: "stackId", GoGetter: "StackId"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeAttributes", GoMethod: "SynthesizeAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeHclAttributes", GoMethod: "SynthesizeHclAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "toHclTerraform", GoMethod: "ToHclTerraform"},
			_jsii_.MemberMethod{JsiiMethod: "toMetadata", GoMethod: "ToMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "toTerraform", GoMethod: "ToTerraform"},
			_jsii_.MemberProperty{JsiiProperty: "version", GoGetter: "Version"},
		},
		func() interface{} {
			j := jsiiProxy_Sql{}
			_jsii_.InitJsiiProxy(&j.Type__cdktfTerraformModule)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"sql.SqlConfig",
		reflect.TypeOf((*SqlConfig)(nil)).Elem(),
	)
}
