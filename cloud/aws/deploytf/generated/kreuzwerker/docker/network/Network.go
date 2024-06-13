package network

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/network/internal"
)

// Represents a {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network docker_network}.
type Network interface {
	cdktf.TerraformResource
	Attachable() interface{}
	SetAttachable(val interface{})
	AttachableInput() interface{}
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	CheckDuplicate() interface{}
	SetCheckDuplicate(val interface{})
	CheckDuplicateInput() interface{}
	// Experimental.
	Connection() interface{}
	// Experimental.
	SetConnection(val interface{})
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
	Driver() *string
	SetDriver(val *string)
	DriverInput() *string
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
	Ingress() interface{}
	SetIngress(val interface{})
	IngressInput() interface{}
	Internal() interface{}
	SetInternal(val interface{})
	InternalInput() interface{}
	IpamConfig() NetworkIpamConfigList
	IpamConfigInput() interface{}
	IpamDriver() *string
	SetIpamDriver(val *string)
	IpamDriverInput() *string
	IpamOptions() *map[string]*string
	SetIpamOptions(val *map[string]*string)
	IpamOptionsInput() *map[string]*string
	Ipv6() interface{}
	SetIpv6(val interface{})
	Ipv6Input() interface{}
	Labels() NetworkLabelsList
	LabelsInput() interface{}
	// Experimental.
	Lifecycle() *cdktf.TerraformResourceLifecycle
	// Experimental.
	SetLifecycle(val *cdktf.TerraformResourceLifecycle)
	Name() *string
	SetName(val *string)
	NameInput() *string
	// The tree node.
	Node() constructs.Node
	Options() *map[string]*string
	SetOptions(val *map[string]*string)
	OptionsInput() *map[string]*string
	// Experimental.
	Provider() cdktf.TerraformProvider
	// Experimental.
	SetProvider(val cdktf.TerraformProvider)
	// Experimental.
	Provisioners() *[]interface{}
	// Experimental.
	SetProvisioners(val *[]interface{})
	// Experimental.
	RawOverrides() interface{}
	Scope() *string
	// Experimental.
	TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata
	// Experimental.
	TerraformMetaArguments() *map[string]interface{}
	// Experimental.
	TerraformResourceType() *string
	// Adds a user defined moveTarget string to this resource to be later used in .moveTo(moveTarget) to resolve the location of the move.
	// Experimental.
	AddMoveTarget(moveTarget *string)
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
	HasResourceMove() interface{}
	// Experimental.
	ImportFrom(id *string, provider cdktf.TerraformProvider)
	// Experimental.
	InterpolationForAttribute(terraformAttribute *string) cdktf.IResolvable
	// Move the resource corresponding to "id" to this resource.
	//
	// Note that the resource being moved from must be marked as moved using it's instance function.
	// Experimental.
	MoveFromId(id *string)
	// Moves this resource to the target resource given by moveTarget.
	// Experimental.
	MoveTo(moveTarget *string, index interface{})
	// Moves this resource to the resource corresponding to "id".
	// Experimental.
	MoveToId(id *string)
	// Overrides the auto-generated logical ID with a specific ID.
	// Experimental.
	OverrideLogicalId(newLogicalId *string)
	PutIpamConfig(value interface{})
	PutLabels(value interface{})
	ResetAttachable()
	ResetCheckDuplicate()
	ResetDriver()
	ResetId()
	ResetIngress()
	ResetInternal()
	ResetIpamConfig()
	ResetIpamDriver()
	ResetIpamOptions()
	ResetIpv6()
	ResetLabels()
	ResetOptions()
	// Resets a previously passed logical Id to use the auto-generated logical id again.
	// Experimental.
	ResetOverrideLogicalId()
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

// The jsii proxy struct for Network
type jsiiProxy_Network struct {
	internal.Type__cdktfTerraformResource
}

func (j *jsiiProxy_Network) Attachable() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"attachable",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) AttachableInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"attachableInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) CheckDuplicate() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"checkDuplicate",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) CheckDuplicateInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"checkDuplicateInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Connection() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"connection",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Count() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"count",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Driver() *string {
	var returns *string
	_jsii_.Get(
		j,
		"driver",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) DriverInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"driverInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Id() *string {
	var returns *string
	_jsii_.Get(
		j,
		"id",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IdInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"idInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Ingress() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ingress",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IngressInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ingressInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Internal() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"internal",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) InternalInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"internalInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IpamConfig() NetworkIpamConfigList {
	var returns NetworkIpamConfigList
	_jsii_.Get(
		j,
		"ipamConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IpamConfigInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ipamConfigInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IpamDriver() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ipamDriver",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IpamDriverInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"ipamDriverInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IpamOptions() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"ipamOptions",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) IpamOptionsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"ipamOptionsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Ipv6() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ipv6",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Ipv6Input() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ipv6Input",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Labels() NetworkLabelsList {
	var returns NetworkLabelsList
	_jsii_.Get(
		j,
		"labels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) LabelsInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"labelsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Lifecycle() *cdktf.TerraformResourceLifecycle {
	var returns *cdktf.TerraformResourceLifecycle
	_jsii_.Get(
		j,
		"lifecycle",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) NameInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"nameInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Options() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"options",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) OptionsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"optionsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Provider() cdktf.TerraformProvider {
	var returns cdktf.TerraformProvider
	_jsii_.Get(
		j,
		"provider",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Provisioners() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"provisioners",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) Scope() *string {
	var returns *string
	_jsii_.Get(
		j,
		"scope",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) TerraformGeneratorMetadata() *cdktf.TerraformProviderGeneratorMetadata {
	var returns *cdktf.TerraformProviderGeneratorMetadata
	_jsii_.Get(
		j,
		"terraformGeneratorMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) TerraformMetaArguments() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"terraformMetaArguments",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Network) TerraformResourceType() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformResourceType",
		&returns,
	)
	return returns
}


// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network docker_network} Resource.
func NewNetwork(scope constructs.Construct, id *string, config *NetworkConfig) Network {
	_init_.Initialize()

	if err := validateNewNetworkParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Network{}

	_jsii_.Create(
		"docker.network.Network",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

// Create a new {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/network docker_network} Resource.
func NewNetwork_Override(n Network, scope constructs.Construct, id *string, config *NetworkConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.network.Network",
		[]interface{}{scope, id, config},
		n,
	)
}

func (j *jsiiProxy_Network)SetAttachable(val interface{}) {
	if err := j.validateSetAttachableParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"attachable",
		val,
	)
}

func (j *jsiiProxy_Network)SetCheckDuplicate(val interface{}) {
	if err := j.validateSetCheckDuplicateParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"checkDuplicate",
		val,
	)
}

func (j *jsiiProxy_Network)SetConnection(val interface{}) {
	if err := j.validateSetConnectionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"connection",
		val,
	)
}

func (j *jsiiProxy_Network)SetCount(val interface{}) {
	if err := j.validateSetCountParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"count",
		val,
	)
}

func (j *jsiiProxy_Network)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Network)SetDriver(val *string) {
	if err := j.validateSetDriverParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"driver",
		val,
	)
}

func (j *jsiiProxy_Network)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Network)SetId(val *string) {
	if err := j.validateSetIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"id",
		val,
	)
}

func (j *jsiiProxy_Network)SetIngress(val interface{}) {
	if err := j.validateSetIngressParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"ingress",
		val,
	)
}

func (j *jsiiProxy_Network)SetInternal(val interface{}) {
	if err := j.validateSetInternalParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internal",
		val,
	)
}

func (j *jsiiProxy_Network)SetIpamDriver(val *string) {
	if err := j.validateSetIpamDriverParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"ipamDriver",
		val,
	)
}

func (j *jsiiProxy_Network)SetIpamOptions(val *map[string]*string) {
	if err := j.validateSetIpamOptionsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"ipamOptions",
		val,
	)
}

func (j *jsiiProxy_Network)SetIpv6(val interface{}) {
	if err := j.validateSetIpv6Parameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"ipv6",
		val,
	)
}

func (j *jsiiProxy_Network)SetLifecycle(val *cdktf.TerraformResourceLifecycle) {
	if err := j.validateSetLifecycleParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"lifecycle",
		val,
	)
}

func (j *jsiiProxy_Network)SetName(val *string) {
	if err := j.validateSetNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_Network)SetOptions(val *map[string]*string) {
	if err := j.validateSetOptionsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"options",
		val,
	)
}

func (j *jsiiProxy_Network)SetProvider(val cdktf.TerraformProvider) {
	_jsii_.Set(
		j,
		"provider",
		val,
	)
}

func (j *jsiiProxy_Network)SetProvisioners(val *[]interface{}) {
	if err := j.validateSetProvisionersParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"provisioners",
		val,
	)
}

// Generates CDKTF code for importing a Network resource upon running "cdktf plan <stack-name>".
func Network_GenerateConfigForImport(scope constructs.Construct, importToId *string, importFromId *string, provider cdktf.TerraformProvider) cdktf.ImportableResource {
	_init_.Initialize()

	if err := validateNetwork_GenerateConfigForImportParameters(scope, importToId, importFromId); err != nil {
		panic(err)
	}
	var returns cdktf.ImportableResource

	_jsii_.StaticInvoke(
		"docker.network.Network",
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
func Network_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateNetwork_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.network.Network",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Network_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateNetwork_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.network.Network",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Network_IsTerraformResource(x interface{}) *bool {
	_init_.Initialize()

	if err := validateNetwork_IsTerraformResourceParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"docker.network.Network",
		"isTerraformResource",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func Network_TfResourceType() *string {
	_init_.Initialize()
	var returns *string
	_jsii_.StaticGet(
		"docker.network.Network",
		"tfResourceType",
		&returns,
	)
	return returns
}

func (n *jsiiProxy_Network) AddMoveTarget(moveTarget *string) {
	if err := n.validateAddMoveTargetParameters(moveTarget); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"addMoveTarget",
		[]interface{}{moveTarget},
	)
}

func (n *jsiiProxy_Network) AddOverride(path *string, value interface{}) {
	if err := n.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (n *jsiiProxy_Network) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
	if err := n.validateGetAnyMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]interface{}

	_jsii_.Invoke(
		n,
		"getAnyMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := n.validateGetBooleanAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		n,
		"getBooleanAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
	if err := n.validateGetBooleanMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*bool

	_jsii_.Invoke(
		n,
		"getBooleanMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetListAttribute(terraformAttribute *string) *[]*string {
	if err := n.validateGetListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*string

	_jsii_.Invoke(
		n,
		"getListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetNumberAttribute(terraformAttribute *string) *float64 {
	if err := n.validateGetNumberAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		n,
		"getNumberAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
	if err := n.validateGetNumberListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*float64

	_jsii_.Invoke(
		n,
		"getNumberListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
	if err := n.validateGetNumberMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*float64

	_jsii_.Invoke(
		n,
		"getNumberMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetStringAttribute(terraformAttribute *string) *string {
	if err := n.validateGetStringAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		n,
		"getStringAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
	if err := n.validateGetStringMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*string

	_jsii_.Invoke(
		n,
		"getStringMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) HasResourceMove() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		n,
		"hasResourceMove",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) ImportFrom(id *string, provider cdktf.TerraformProvider) {
	if err := n.validateImportFromParameters(id); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"importFrom",
		[]interface{}{id, provider},
	)
}

func (n *jsiiProxy_Network) InterpolationForAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := n.validateInterpolationForAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		n,
		"interpolationForAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) MoveFromId(id *string) {
	if err := n.validateMoveFromIdParameters(id); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"moveFromId",
		[]interface{}{id},
	)
}

func (n *jsiiProxy_Network) MoveTo(moveTarget *string, index interface{}) {
	if err := n.validateMoveToParameters(moveTarget, index); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"moveTo",
		[]interface{}{moveTarget, index},
	)
}

func (n *jsiiProxy_Network) MoveToId(id *string) {
	if err := n.validateMoveToIdParameters(id); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"moveToId",
		[]interface{}{id},
	)
}

func (n *jsiiProxy_Network) OverrideLogicalId(newLogicalId *string) {
	if err := n.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (n *jsiiProxy_Network) PutIpamConfig(value interface{}) {
	if err := n.validatePutIpamConfigParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"putIpamConfig",
		[]interface{}{value},
	)
}

func (n *jsiiProxy_Network) PutLabels(value interface{}) {
	if err := n.validatePutLabelsParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		n,
		"putLabels",
		[]interface{}{value},
	)
}

func (n *jsiiProxy_Network) ResetAttachable() {
	_jsii_.InvokeVoid(
		n,
		"resetAttachable",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetCheckDuplicate() {
	_jsii_.InvokeVoid(
		n,
		"resetCheckDuplicate",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetDriver() {
	_jsii_.InvokeVoid(
		n,
		"resetDriver",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetId() {
	_jsii_.InvokeVoid(
		n,
		"resetId",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetIngress() {
	_jsii_.InvokeVoid(
		n,
		"resetIngress",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetInternal() {
	_jsii_.InvokeVoid(
		n,
		"resetInternal",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetIpamConfig() {
	_jsii_.InvokeVoid(
		n,
		"resetIpamConfig",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetIpamDriver() {
	_jsii_.InvokeVoid(
		n,
		"resetIpamDriver",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetIpamOptions() {
	_jsii_.InvokeVoid(
		n,
		"resetIpamOptions",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetIpv6() {
	_jsii_.InvokeVoid(
		n,
		"resetIpv6",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetLabels() {
	_jsii_.InvokeVoid(
		n,
		"resetLabels",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetOptions() {
	_jsii_.InvokeVoid(
		n,
		"resetOptions",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		n,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (n *jsiiProxy_Network) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		n,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		n,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		n,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		n,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		n,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (n *jsiiProxy_Network) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		n,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

