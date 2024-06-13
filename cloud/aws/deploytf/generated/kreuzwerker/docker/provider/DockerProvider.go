package provider

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/provider/internal"
)

// Represents a {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs docker}.
type DockerProvider interface {
	cdktf.TerraformProvider
	Alias() *string
	SetAlias(val *string)
	AliasInput() *string
	CaMaterial() *string
	SetCaMaterial(val *string)
	CaMaterialInput() *string
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	CertMaterial() *string
	SetCertMaterial(val *string)
	CertMaterialInput() *string
	CertPath() *string
	SetCertPath(val *string)
	CertPathInput() *string
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	Host() *string
	SetHost(val *string)
	HostInput() *string
	KeyMaterial() *string
	SetKeyMaterial(val *string)
	KeyMaterialInput() *string
	// Experimental.
	MetaAttributes() *map[string]interface{}
	// The tree node.
	Node() constructs.Node
	// Experimental.
	RawOverrides() interface{}
	RegistryAuth() interface{}
	SetRegistryAuth(val interface{})
	RegistryAuthInput() interface{}
	SshOpts() *[]*string
	SetSshOpts(val *[]*string)
	SshOptsInput() *[]*string
	// Experimental.
	TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata
	// Experimental.
	TerraformProviderSource() *string
	// Experimental.
	TerraformResourceType() *string
	// Experimental.
	AddOverride(path *string, value interface{})
	// Overrides the auto-generated logical ID with a specific ID.
	// Experimental.
	OverrideLogicalId(newLogicalId *string)
	ResetAlias()
	ResetCaMaterial()
	ResetCertMaterial()
	ResetCertPath()
	ResetHost()
	ResetKeyMaterial()
	// Resets a previously passed logical Id to use the auto-generated logical id again.
	// Experimental.
	ResetOverrideLogicalId()
	ResetRegistryAuth()
	ResetSshOpts()
	SynthesizeAttributes() *map[string]interface{}
	SynthesizeHclAttributes() *map[string]interface{}
	// Experimental.
	ToHclTerraform() interface{}
	// Experimental.
	ToMetadata() interface{}
	// Returns a string representation of this construct.
	ToString() *string
	// Adds this resource to the terraform JSON output.
	// Experimental.
	ToTerraform() interface{}
}

// The jsii proxy struct for DockerProvider
type jsiiProxy_DockerProvider struct {
	internal.Type__cdktfTerraformProvider
}

func (j *jsiiProxy_DockerProvider) Alias() *string {
	var returns *string
	_jsii_.Get(
		j,
		"alias",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) AliasInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"aliasInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) CaMaterial() *string {
	var returns *string
	_jsii_.Get(
		j,
		"caMaterial",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) CaMaterialInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"caMaterialInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) CertMaterial() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certMaterial",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) CertMaterialInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certMaterialInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) CertPath() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certPath",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) CertPathInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"certPathInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) Host() *string {
	var returns *string
	_jsii_.Get(
		j,
		"host",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) HostInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"hostInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) KeyMaterial() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyMaterial",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) KeyMaterialInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyMaterialInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) MetaAttributes() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"metaAttributes",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) RegistryAuth() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"registryAuth",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) RegistryAuthInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"registryAuthInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) SshOpts() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"sshOpts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) SshOptsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"sshOptsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata {
	var returns *cdktf.TerraformProviderGeneratorMetadata
	_jsii_.Get(
		j,
		"terraformGeneratorMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) TerraformProviderSource() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformProviderSource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DockerProvider) TerraformResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformResourceType",
		&returns,
	)
	return returns
}


// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs docker} Resource.
func NewDockerProvider(scope constructs.Construct, id *string, config *DockerProviderConfig) DockerProvider {
	_init_.Initialize()

	if err := validateNewDockerProviderParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_DockerProvider{}

	_jsii_.Create(
		"docker.provider.DockerProvider",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs docker} Resource.
func NewDockerProvider_Override(d DockerProvider, scope constructs.Construct, id *string, config *DockerProviderConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.provider.DockerProvider",
		[]interface{}{scope, id, config},
		d,
	)
}

func (j *jsiiProxy_DockerProvider)SetAlias(val *string) {
	_jsii_.Set(
		j,
		"alias",
		val,
	)
}

func (j *jsiiProxy_DockerProvider)SetCaMaterial(val *string) {
	_jsii_.Set(
		j,
		"caMaterial",
		val,
	)
}

func (j *jsiiProxy_DockerProvider)SetCertMaterial(val *string) {
	_jsii_.Set(
		j,
		"certMaterial",
		val,
	)
}

func (j *jsiiProxy_DockerProvider)SetCertPath(val *string) {
	_jsii_.Set(
		j,
		"certPath",
		val,
	)
}

func (j *jsiiProxy_DockerProvider)SetHost(val *string) {
	_jsii_.Set(
		j,
		"host",
		val,
	)
}

func (j *jsiiProxy_DockerProvider)SetKeyMaterial(val *string) {
	_jsii_.Set(
		j,
		"keyMaterial",
		val,
	)
}

func (j *jsiiProxy_DockerProvider)SetRegistryAuth(val interface{}) {
	if err := j.validateSetRegistryAuthParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"registryAuth",
		val,
	)
}

func (j *jsiiProxy_DockerProvider)SetSshOpts(val *[]*string) {
	_jsii_.Set(
		j,
		"sshOpts",
		val,
	)
}

// Generates CDKTF code for importing a DockerProvider resource upon running "cdktf plan <stack-name>".
func DockerProvider_GenerateConfigForImport(scope constructs.Construct, importToId *string, importFromId *string, provider cdktf.TerraformProvider) cdktf.ImportableResource {
	_init_.Initialize()

	if err := validateDockerProvider_GenerateConfigForImportParameters(scope, importToId, importFromId); err != nil {
		panic(err)
	}
	var returns cdktf.ImportableResource

	_jsii_.StaticInvoke(
		"docker.provider.DockerProvider",
		"generateConfigForImport",
		[]interface{}{scope, importToId, importFromId, provider},
		&returns,
	)

	return returns
}

// Checks if `x` is a construct.
//
// Use this method instead of `instanceof` to properly detect `Construct`
// instances, even when the construct library is symlinked.
//
// Explanation: in JavaScript, multiple copies of the `constructs` library on
// disk are seen as independent, completely different libraries. As a
// consequence, the class `Construct` in each copy of the `constructs` library
// is seen as a different class, and an instance of one class will not test as
// `instanceof` the other class. `npm install` will not create installations
// like this, but users may manually symlink construct libraries together or
// use a monorepo tool: in those cases, multiple copies of the `constructs`
// library can be accidentally installed, and `instanceof` will behave
// unpredictably. It is safest to avoid using `instanceof`, and using
// this type-testing method instead.
//
// Returns: true if `x` is an object created from a class which extends `Construct`.
func DockerProvider_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDockerProvider_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.provider.DockerProvider",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func DockerProvider_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDockerProvider_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.provider.DockerProvider",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func DockerProvider_IsTerraformProvider(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDockerProvider_IsTerraformProviderParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.provider.DockerProvider",
		"isTerraformProvider",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func DockerProvider_TfResourceType() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"docker.provider.DockerProvider",
		"tfResourceType",
		&returns,
	)
	return returns
}

func (d *jsiiProxy_DockerProvider) AddOverride(path *string, value interface{}) {
	if err := d.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		d,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (d *jsiiProxy_DockerProvider) OverrideLogicalId(newLogicalId *string) {
	if err := d.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		d,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (d *jsiiProxy_DockerProvider) ResetAlias() {
	_jsii_.InvokeVoid(
		d,
		"resetAlias",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetCaMaterial() {
	_jsii_.InvokeVoid(
		d,
		"resetCaMaterial",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetCertMaterial() {
	_jsii_.InvokeVoid(
		d,
		"resetCertMaterial",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetCertPath() {
	_jsii_.InvokeVoid(
		d,
		"resetCertPath",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetHost() {
	_jsii_.InvokeVoid(
		d,
		"resetHost",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetKeyMaterial() {
	_jsii_.InvokeVoid(
		d,
		"resetKeyMaterial",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		d,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetRegistryAuth() {
	_jsii_.InvokeVoid(
		d,
		"resetRegistryAuth",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) ResetSshOpts() {
	_jsii_.InvokeVoid(
		d,
		"resetSshOpts",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DockerProvider) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		d,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DockerProvider) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		d,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DockerProvider) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DockerProvider) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DockerProvider) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DockerProvider) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

