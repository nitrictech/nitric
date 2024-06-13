package provider

import (
	"reflect"

	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
)

func init() {
	_jsii_.RegisterClass(
		"docker.provider.DockerProvider",
		reflect.TypeOf((*DockerProvider)(nil)).Elem(),
		[]_jsii_.Member{
			_jsii_.MemberMethod{JsiiMethod: "addOverride", GoMethod: "AddOverride"},
			_jsii_.MemberProperty{JsiiProperty: "alias", GoGetter: "Alias"},
			_jsii_.MemberProperty{JsiiProperty: "aliasInput", GoGetter: "AliasInput"},
			_jsii_.MemberProperty{JsiiProperty: "caMaterial", GoGetter: "CaMaterial"},
			_jsii_.MemberProperty{JsiiProperty: "caMaterialInput", GoGetter: "CaMaterialInput"},
			_jsii_.MemberProperty{JsiiProperty: "cdktfStack", GoGetter: "CdktfStack"},
			_jsii_.MemberProperty{JsiiProperty: "certMaterial", GoGetter: "CertMaterial"},
			_jsii_.MemberProperty{JsiiProperty: "certMaterialInput", GoGetter: "CertMaterialInput"},
			_jsii_.MemberProperty{JsiiProperty: "certPath", GoGetter: "CertPath"},
			_jsii_.MemberProperty{JsiiProperty: "certPathInput", GoGetter: "CertPathInput"},
			_jsii_.MemberProperty{JsiiProperty: "constructNodeMetadata", GoGetter: "ConstructNodeMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "fqn", GoGetter: "Fqn"},
			_jsii_.MemberProperty{JsiiProperty: "friendlyUniqueId", GoGetter: "FriendlyUniqueId"},
			_jsii_.MemberProperty{JsiiProperty: "host", GoGetter: "Host"},
			_jsii_.MemberProperty{JsiiProperty: "hostInput", GoGetter: "HostInput"},
			_jsii_.MemberProperty{JsiiProperty: "keyMaterial", GoGetter: "KeyMaterial"},
			_jsii_.MemberProperty{JsiiProperty: "keyMaterialInput", GoGetter: "KeyMaterialInput"},
			_jsii_.MemberProperty{JsiiProperty: "metaAttributes", GoGetter: "MetaAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "node", GoGetter: "Node"},
			_jsii_.MemberMethod{JsiiMethod: "overrideLogicalId", GoMethod: "OverrideLogicalId"},
			_jsii_.MemberProperty{JsiiProperty: "rawOverrides", GoGetter: "RawOverrides"},
			_jsii_.MemberProperty{JsiiProperty: "registryAuth", GoGetter: "RegistryAuth"},
			_jsii_.MemberProperty{JsiiProperty: "registryAuthInput", GoGetter: "RegistryAuthInput"},
			_jsii_.MemberMethod{JsiiMethod: "resetAlias", GoMethod: "ResetAlias"},
			_jsii_.MemberMethod{JsiiMethod: "resetCaMaterial", GoMethod: "ResetCaMaterial"},
			_jsii_.MemberMethod{JsiiMethod: "resetCertMaterial", GoMethod: "ResetCertMaterial"},
			_jsii_.MemberMethod{JsiiMethod: "resetCertPath", GoMethod: "ResetCertPath"},
			_jsii_.MemberMethod{JsiiMethod: "resetHost", GoMethod: "ResetHost"},
			_jsii_.MemberMethod{JsiiMethod: "resetKeyMaterial", GoMethod: "ResetKeyMaterial"},
			_jsii_.MemberMethod{JsiiMethod: "resetOverrideLogicalId", GoMethod: "ResetOverrideLogicalId"},
			_jsii_.MemberMethod{JsiiMethod: "resetRegistryAuth", GoMethod: "ResetRegistryAuth"},
			_jsii_.MemberMethod{JsiiMethod: "resetSshOpts", GoMethod: "ResetSshOpts"},
			_jsii_.MemberProperty{JsiiProperty: "sshOpts", GoGetter: "SshOpts"},
			_jsii_.MemberProperty{JsiiProperty: "sshOptsInput", GoGetter: "SshOptsInput"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeAttributes", GoMethod: "SynthesizeAttributes"},
			_jsii_.MemberMethod{JsiiMethod: "synthesizeHclAttributes", GoMethod: "SynthesizeHclAttributes"},
			_jsii_.MemberProperty{JsiiProperty: "terraformGeneratorMetadata", GoGetter: "TerraformGeneratorMetadata"},
			_jsii_.MemberProperty{JsiiProperty: "terraformProviderSource", GoGetter: "TerraformProviderSource"},
			_jsii_.MemberProperty{JsiiProperty: "terraformResourceType", GoGetter: "TerraformResourceType"},
			_jsii_.MemberMethod{JsiiMethod: "toHclTerraform", GoMethod: "ToHclTerraform"},
			_jsii_.MemberMethod{JsiiMethod: "toMetadata", GoMethod: "ToMetadata"},
			_jsii_.MemberMethod{JsiiMethod: "toString", GoMethod: "ToString"},
			_jsii_.MemberMethod{JsiiMethod: "toTerraform", GoMethod: "ToTerraform"},
		},
		func() interface{} {
			j := jsiiProxy_DockerProvider{}
			_jsii_.InitJsiiProxy(&j.Type__cdktfTerraformProvider)
			return &j
		},
	)
	_jsii_.RegisterStruct(
		"docker.provider.DockerProviderConfig",
		reflect.TypeOf((*DockerProviderConfig)(nil)).Elem(),
	)
	_jsii_.RegisterStruct(
		"docker.provider.DockerProviderRegistryAuth",
		reflect.TypeOf((*DockerProviderRegistryAuth)(nil)).Elem(),
	)
}
