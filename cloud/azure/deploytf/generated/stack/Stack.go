package stack

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/azure/deploytf/generated/stack/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/stack/internal"
)

// Defines an Stack based on a Terraform module.
//
// Source at ./.nitric/modules/stack
type Stack interface {
	cdktf.TerraformModule
	AppIdentityClientIdOutput() *string
	AppIdentityOutput() *string
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	ContainerAppEnvironmentIdOutput() *string
	ContainerAppSubnetIdOutput() *string
	DatabaseMasterPasswordOutput() *string
	DatabaseServerFqdnOutput() *string
	DatabaseServerIdOutput() *string
	DatabaseServerNameOutput() *string
	// Experimental.
	DependsOn() *[]*string
	// Experimental.
	SetDependsOn(val *[]*string)
	EnableDatabase() *bool
	SetEnableDatabase(val *bool)
	EnableKeyvault() *bool
	SetEnableKeyvault(val *bool)
	EnableStorage() *bool
	SetEnableStorage(val *bool)
	// Experimental.
	ForEach() cdktf.ITerraformIterator
	// Experimental.
	SetForEach(val cdktf.ITerraformIterator)
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	InfrastructureSubnetId() *string
	SetInfrastructureSubnetId(val *string)
	InfrastructureSubnetIdOutput() *string
	KeyvaultNameOutput() *string
	Location() *string
	SetLocation(val *string)
	// The tree node.
	Node() constructs.Node
	// Experimental.
	Providers() *[]interface{}
	// Experimental.
	RawOverrides() interface{}
	RegistryLoginServerOutput() *string
	RegistryPasswordOutput() *string
	RegistryUsernameOutput() *string
	ResourceGroupName() *string
	SetResourceGroupName(val *string)
	ResourceGroupNameOutput() *string
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	StackIdOutput() *string
	StackName() *string
	SetStackName(val *string)
	StackNameOutput() *string
	StorageAccountBlobEndpointOutput() *string
	StorageAccountConnectionStringOutput() *string
	StorageAccountIdOutput() *string
	StorageAccountNameOutput() *string
	StorageAccountQueueEndpointOutput() *string
	SubscriptionIdOutput() *string
	Tags() *map[string]*string
	SetTags(val *map[string]*string)
	// Experimental.
	Version() *string
	// Experimental.
	AddOverride(path *string, value interface{})
	// Experimental.
	AddProvider(provider interface{})
	// Experimental.
	GetString(output *string) *string
	// Experimental.
	InterpolationForOutput(moduleOutput *string) cdktf.IResolvable
	// Overrides the auto-generated logical ID with a specific ID.
	// Experimental.
	OverrideLogicalId(newLogicalId *string)
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
	// Experimental.
	ToTerraform() interface{}
}

// The jsii proxy struct for Stack
type jsiiProxy_Stack struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Stack) AppIdentityClientIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"appIdentityClientIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) AppIdentityOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"appIdentityOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ContainerAppEnvironmentIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"containerAppEnvironmentIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ContainerAppSubnetIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"containerAppSubnetIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) DatabaseMasterPasswordOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"databaseMasterPasswordOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) DatabaseServerFqdnOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"databaseServerFqdnOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) DatabaseServerIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"databaseServerIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) DatabaseServerNameOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"databaseServerNameOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) EnableDatabase() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"enableDatabase",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) EnableKeyvault() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"enableKeyvault",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) EnableStorage() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"enableStorage",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) InfrastructureSubnetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"infrastructureSubnetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) InfrastructureSubnetIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"infrastructureSubnetIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) KeyvaultNameOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"keyvaultNameOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Location() *string {
	var returns *string
	_jsii_.Get(
		j,
		"location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) RegistryLoginServerOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"registryLoginServerOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) RegistryPasswordOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"registryPasswordOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) RegistryUsernameOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"registryUsernameOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ResourceGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) ResourceGroupNameOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceGroupNameOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StackIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StackName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StackNameOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackNameOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StorageAccountBlobEndpointOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"storageAccountBlobEndpointOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StorageAccountConnectionStringOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"storageAccountConnectionStringOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StorageAccountIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"storageAccountIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StorageAccountNameOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"storageAccountNameOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) StorageAccountQueueEndpointOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"storageAccountQueueEndpointOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) SubscriptionIdOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"subscriptionIdOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Tags() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"tags",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Stack) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


func NewStack(scope constructs.Construct, id *string, config *StackConfig) Stack {
	_init_.Initialize()

	if err := validateNewStackParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Stack{}

	_jsii_.Create(
		"stack.Stack",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewStack_Override(s Stack, scope constructs.Construct, id *string, config *StackConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"stack.Stack",
		[]interface{}{scope, id, config},
		s,
	)
}

func (j *jsiiProxy_Stack)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Stack)SetEnableDatabase(val *bool) {
	_jsii_.Set(
		j,
		"enableDatabase",
		val,
	)
}

func (j *jsiiProxy_Stack)SetEnableKeyvault(val *bool) {
	_jsii_.Set(
		j,
		"enableKeyvault",
		val,
	)
}

func (j *jsiiProxy_Stack)SetEnableStorage(val *bool) {
	_jsii_.Set(
		j,
		"enableStorage",
		val,
	)
}

func (j *jsiiProxy_Stack)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Stack)SetInfrastructureSubnetId(val *string) {
	_jsii_.Set(
		j,
		"infrastructureSubnetId",
		val,
	)
}

func (j *jsiiProxy_Stack)SetLocation(val *string) {
	if err := j.validateSetLocationParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"location",
		val,
	)
}

func (j *jsiiProxy_Stack)SetResourceGroupName(val *string) {
	_jsii_.Set(
		j,
		"resourceGroupName",
		val,
	)
}

func (j *jsiiProxy_Stack)SetStackName(val *string) {
	if err := j.validateSetStackNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stackName",
		val,
	)
}

func (j *jsiiProxy_Stack)SetTags(val *map[string]*string) {
	if err := j.validateSetTagsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"tags",
		val,
	)
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
func Stack_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateStack_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"stack.Stack",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Stack_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateStack_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"stack.Stack",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) AddOverride(path *string, value interface{}) {
	if err := s.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (s *jsiiProxy_Stack) AddProvider(provider interface{}) {
	if err := s.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addProvider",
		[]interface{}{provider},
	)
}

func (s *jsiiProxy_Stack) GetString(output *string) *string {
	if err := s.validateGetStringParameters(output); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		s,
		"getString",
		[]interface{}{output},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
	if err := s.validateInterpolationForOutputParameters(moduleOutput); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		s,
		"interpolationForOutput",
		[]interface{}{moduleOutput},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) OverrideLogicalId(newLogicalId *string) {
	if err := s.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (s *jsiiProxy_Stack) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		s,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (s *jsiiProxy_Stack) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Stack) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

