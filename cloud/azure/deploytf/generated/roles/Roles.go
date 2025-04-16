package roles

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/azure/deploytf/generated/roles/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/roles/internal"
)

// Defines an Roles based on a Terraform module.
//
// Source at ./.nitric/modules/roles
type Roles interface {
	cdktf.TerraformModule
	AllowUserDelegationKeyGenerationOutput() *string
	BucketDeleteOutput() *string
	BucketListOutput() *string
	BucketReadOutput() *string
	BucketWriteOutput() *string
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	// Experimental.
	DependsOn() *[]*string
	// Experimental.
	SetDependsOn(val *[]*string)
	// Experimental.
	ForEach() cdktf.ITerraformIterator
	// Experimental.
	SetForEach(val cdktf.ITerraformIterator)
	// Experimental.
	Fqn() *string
	// Experimental.
	FriendlyUniqueId() *string
	KvDeleteOutput() *string
	KvReadOutput() *string
	KvWriteOutput() *string
	// The tree node.
	Node() constructs.Node
	// Experimental.
	Providers() *[]interface{}
	QueueDequeueOutput() *string
	QueueEnqueueOutput() *string
	// Experimental.
	RawOverrides() interface{}
	ResourceGroupName() *string
	SetResourceGroupName(val *string)
	SecretAccessOutput() *string
	SecretPutOutput() *string
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	StackName() *string
	SetStackName(val *string)
	TopicPublishOutput() *string
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

// The jsii proxy struct for Roles
type jsiiProxy_Roles struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Roles) AllowUserDelegationKeyGenerationOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"allowUserDelegationKeyGenerationOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) BucketDeleteOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketDeleteOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) BucketListOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketListOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) BucketReadOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketReadOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) BucketWriteOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"bucketWriteOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) KvDeleteOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kvDeleteOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) KvReadOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kvReadOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) KvWriteOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"kvWriteOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) QueueDequeueOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueDequeueOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) QueueEnqueueOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"queueEnqueueOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) ResourceGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) SecretAccessOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretAccessOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) SecretPutOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"secretPutOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) StackName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) TopicPublishOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"topicPublishOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Roles) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


func NewRoles(scope constructs.Construct, id *string, config *RolesConfig) Roles {
	_init_.Initialize()

	if err := validateNewRolesParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Roles{}

	_jsii_.Create(
		"roles.Roles",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewRoles_Override(r Roles, scope constructs.Construct, id *string, config *RolesConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"roles.Roles",
		[]interface{}{scope, id, config},
		r,
	)
}

func (j *jsiiProxy_Roles)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Roles)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Roles)SetResourceGroupName(val *string) {
	if err := j.validateSetResourceGroupNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"resourceGroupName",
		val,
	)
}

func (j *jsiiProxy_Roles)SetStackName(val *string) {
	if err := j.validateSetStackNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stackName",
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
func Roles_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateRoles_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"roles.Roles",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Roles_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateRoles_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"roles.Roles",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) AddOverride(path *string, value interface{}) {
	if err := r.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		r,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (r *jsiiProxy_Roles) AddProvider(provider interface{}) {
	if err := r.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		r,
		"addProvider",
		[]interface{}{provider},
	)
}

func (r *jsiiProxy_Roles) GetString(output *string) *string {
	if err := r.validateGetStringParameters(output); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		r,
		"getString",
		[]interface{}{output},
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
	if err := r.validateInterpolationForOutputParameters(moduleOutput); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		r,
		"interpolationForOutput",
		[]interface{}{moduleOutput},
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) OverrideLogicalId(newLogicalId *string) {
	if err := r.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		r,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (r *jsiiProxy_Roles) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		r,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (r *jsiiProxy_Roles) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		r,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		r,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		r,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		r,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		r,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (r *jsiiProxy_Roles) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		r,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

