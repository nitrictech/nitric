package datadockerlogs

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/datadockerlogs/internal"
)

// Represents a {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs docker_logs}.
type DataDockerLogs interface {
	cdktf.TerraformDataSource
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	// Experimental.
	Count() interface{}
	// Experimental.
	SetCount(val interface{})
	// Experimental.
	DependsOn() *[]*string
	// Experimental.
	SetDependsOn(val *[]*string)
	Details() interface{}
	SetDetails(val interface{})
	DetailsInput() interface{}
	DiscardHeaders() interface{}
	SetDiscardHeaders(val interface{})
	DiscardHeadersInput() interface{}
	Follow() interface{}
	SetFollow(val interface{})
	FollowInput() interface{}
	// Experimental.
	ForEach() cdktf.ITerraformIterator
	// Experimental.
	SetForEach(val cdktf.ITerraformIterator)
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	Id() *string
	SetId(val *string)
	IdInput() *string
	// Experimental.
	Lifecycle() *cdktf.TerraformResourceLifecycle
	// Experimental.
	SetLifecycle(val *cdktf.TerraformResourceLifecycle)
	LogsListString() *[]*string
	LogsListStringEnabled() interface{}
	SetLogsListStringEnabled(val interface{})
	LogsListStringEnabledInput() interface{}
	Name() *string
	SetName(val *string)
	NameInput() *string
	// The tree node.
	Node() constructs.Node
	// Experimental.
	Provider() cdktf.TerraformProvider
	// Experimental.
	SetProvider(val cdktf.TerraformProvider)
	// Experimental.
	RawOverrides() interface{}
	ShowStderr() interface{}
	SetShowStderr(val interface{})
	ShowStderrInput() interface{}
	ShowStdout() interface{}
	SetShowStdout(val interface{})
	ShowStdoutInput() interface{}
	Since() *string
	SetSince(val *string)
	SinceInput() *string
	Tail() *string
	SetTail(val *string)
	TailInput() *string
	// Experimental.
	TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata
	// Experimental.
	TerraformMetaArguments() *map[string]interface{}
	// Experimental.
	TerraformResourceType() *string
	Timestamps() interface{}
	SetTimestamps(val interface{})
	TimestampsInput() interface{}
	Until() *string
	SetUntil(val *string)
	UntilInput() *string
	// Experimental.
	AddOverride(path *string, value interface{})
	// Experimental.
	GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{}
	// Experimental.
	GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable
	// Experimental.
	GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool
	// Experimental.
	GetListAttribute(terraformAttribute *string) *[]*string
	// Experimental.
	GetNumberAttribute(terraformAttribute *string) *float64
	// Experimental.
	GetNumberListAttribute(terraformAttribute *string) *[]*float64
	// Experimental.
	GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64
	// Experimental.
	GetStringAttribute(terraformAttribute *string) *string
	// Experimental.
	GetStringMapAttribute(terraformAttribute *string) *map[string]*string
	// Experimental.
	InterpolationForAttribute(terraformAttribute *string) cdktf.IResolvable
	// Overrides the auto-generated logical ID with a specific ID.
	// Experimental.
	OverrideLogicalId(newLogicalId *string)
	ResetDetails()
	ResetDiscardHeaders()
	ResetFollow()
	ResetId()
	ResetLogsListStringEnabled()
	// Resets a previously passed logical Id to use the auto-generated logical id again.
	// Experimental.
	ResetOverrideLogicalId()
	ResetShowStderr()
	ResetShowStdout()
	ResetSince()
	ResetTail()
	ResetTimestamps()
	ResetUntil()
	SynthesizeAttributes() *map[string]interface{}
	SynthesizeHclAttributes() *map[string]interface{}
	// Adds this resource to the terraform JSON output.
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

// The jsii proxy struct for DataDockerLogs
type jsiiProxy_DataDockerLogs struct {
	internal.Type__cdktfTerraformDataSource
}

func (j *jsiiProxy_DataDockerLogs) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Count() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"count",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Details() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"details",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) DetailsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"detailsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) DiscardHeaders() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"discardHeaders",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) DiscardHeadersInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"discardHeadersInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Follow() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"follow",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) FollowInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"followInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) IdInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"idInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Lifecycle() *cdktf.TerraformResourceLifecycle {
	var returns *cdktf.TerraformResourceLifecycle
	_jsii_.Get(
		j,
		"lifecycle",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) LogsListString() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"logsListString",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) LogsListStringEnabled() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"logsListStringEnabled",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) LogsListStringEnabledInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"logsListStringEnabledInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) NameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"nameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Provider() cdktf.TerraformProvider {
	var returns cdktf.TerraformProvider
	_jsii_.Get(
		j,
		"provider",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) ShowStderr() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"showStderr",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) ShowStderrInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"showStderrInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) ShowStdout() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"showStdout",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) ShowStdoutInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"showStdoutInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Since() *string {
	var returns *string
	_jsii_.Get(
		j,
		"since",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) SinceInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sinceInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Tail() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tail",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) TailInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"tailInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata {
	var returns *cdktf.TerraformProviderGeneratorMetadata
	_jsii_.Get(
		j,
		"terraformGeneratorMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) TerraformMetaArguments() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"terraformMetaArguments",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) TerraformResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformResourceType",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Timestamps() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"timestamps",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) TimestampsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"timestampsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) Until() *string {
	var returns *string
	_jsii_.Get(
		j,
		"until",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_DataDockerLogs) UntilInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"untilInput",
		&returns,
	)
	return returns
}


// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs docker_logs} Data Source.
func NewDataDockerLogs(scope constructs.Construct, id *string, config *DataDockerLogsConfig) DataDockerLogs {
	_init_.Initialize()

	if err := validateNewDataDockerLogsParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_DataDockerLogs{}

	_jsii_.Create(
		"docker.dataDockerLogs.DataDockerLogs",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs docker_logs} Data Source.
func NewDataDockerLogs_Override(d DataDockerLogs, scope constructs.Construct, id *string, config *DataDockerLogsConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.dataDockerLogs.DataDockerLogs",
		[]interface{}{scope, id, config},
		d,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetCount(val interface{}) {
	if err := j.validateSetCountParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"count",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetDetails(val interface{}) {
	if err := j.validateSetDetailsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"details",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetDiscardHeaders(val interface{}) {
	if err := j.validateSetDiscardHeadersParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"discardHeaders",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetFollow(val interface{}) {
	if err := j.validateSetFollowParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"follow",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetId(val *string) {
	if err := j.validateSetIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"id",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetLifecycle(val *cdktf.TerraformResourceLifecycle) {
	if err := j.validateSetLifecycleParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"lifecycle",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetLogsListStringEnabled(val interface{}) {
	if err := j.validateSetLogsListStringEnabledParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"logsListStringEnabled",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetName(val *string) {
	if err := j.validateSetNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetProvider(val cdktf.TerraformProvider) {
	_jsii_.Set(
		j,
		"provider",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetShowStderr(val interface{}) {
	if err := j.validateSetShowStderrParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"showStderr",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetShowStdout(val interface{}) {
	if err := j.validateSetShowStdoutParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"showStdout",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetSince(val *string) {
	if err := j.validateSetSinceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"since",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetTail(val *string) {
	if err := j.validateSetTailParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"tail",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetTimestamps(val interface{}) {
	if err := j.validateSetTimestampsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"timestamps",
		val,
	)
}

func (j *jsiiProxy_DataDockerLogs)SetUntil(val *string) {
	if err := j.validateSetUntilParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"until",
		val,
	)
}

// Generates CDKTF code for importing a DataDockerLogs resource upon running "cdktf plan <stack-name>".
func DataDockerLogs_GenerateConfigForImport(scope constructs.Construct, importToId *string, importFromId *string, provider cdktf.TerraformProvider) cdktf.ImportableResource {
	_init_.Initialize()

	if err := validateDataDockerLogs_GenerateConfigForImportParameters(scope, importToId, importFromId); err != nil {
		panic(err)
	}
	var returns cdktf.ImportableResource

	_jsii_.StaticInvoke(
		"docker.dataDockerLogs.DataDockerLogs",
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
func DataDockerLogs_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDataDockerLogs_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.dataDockerLogs.DataDockerLogs",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func DataDockerLogs_IsTerraformDataSource(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDataDockerLogs_IsTerraformDataSourceParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.dataDockerLogs.DataDockerLogs",
		"isTerraformDataSource",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func DataDockerLogs_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateDataDockerLogs_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.dataDockerLogs.DataDockerLogs",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func DataDockerLogs_TfResourceType() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"docker.dataDockerLogs.DataDockerLogs",
		"tfResourceType",
		&returns,
	)
	return returns
}

func (d *jsiiProxy_DataDockerLogs) AddOverride(path *string, value interface{}) {
	if err := d.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		d,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (d *jsiiProxy_DataDockerLogs) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
	if err := d.validateGetAnyMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]interface{}

	_jsii_.Invoke(
		d,
		"getAnyMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := d.validateGetBooleanAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		d,
		"getBooleanAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
	if err := d.validateGetBooleanMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*bool

	_jsii_.Invoke(
		d,
		"getBooleanMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetListAttribute(terraformAttribute *string) *[]*string {
	if err := d.validateGetListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*string

	_jsii_.Invoke(
		d,
		"getListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetNumberAttribute(terraformAttribute *string) *float64 {
	if err := d.validateGetNumberAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		d,
		"getNumberAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
	if err := d.validateGetNumberListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*float64

	_jsii_.Invoke(
		d,
		"getNumberListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
	if err := d.validateGetNumberMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*float64

	_jsii_.Invoke(
		d,
		"getNumberMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetStringAttribute(terraformAttribute *string) *string {
	if err := d.validateGetStringAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		d,
		"getStringAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
	if err := d.validateGetStringMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*string

	_jsii_.Invoke(
		d,
		"getStringMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) InterpolationForAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := d.validateInterpolationForAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		d,
		"interpolationForAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) OverrideLogicalId(newLogicalId *string) {
	if err := d.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		d,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetDetails() {
	_jsii_.InvokeVoid(
		d,
		"resetDetails",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetDiscardHeaders() {
	_jsii_.InvokeVoid(
		d,
		"resetDiscardHeaders",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetFollow() {
	_jsii_.InvokeVoid(
		d,
		"resetFollow",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetId() {
	_jsii_.InvokeVoid(
		d,
		"resetId",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetLogsListStringEnabled() {
	_jsii_.InvokeVoid(
		d,
		"resetLogsListStringEnabled",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		d,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetShowStderr() {
	_jsii_.InvokeVoid(
		d,
		"resetShowStderr",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetShowStdout() {
	_jsii_.InvokeVoid(
		d,
		"resetShowStdout",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetSince() {
	_jsii_.InvokeVoid(
		d,
		"resetSince",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetTail() {
	_jsii_.InvokeVoid(
		d,
		"resetTail",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetTimestamps() {
	_jsii_.InvokeVoid(
		d,
		"resetTimestamps",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) ResetUntil() {
	_jsii_.InvokeVoid(
		d,
		"resetUntil",
		nil, // no parameters
	)
}

func (d *jsiiProxy_DataDockerLogs) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		d,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		d,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		d,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (d *jsiiProxy_DataDockerLogs) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		d,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

