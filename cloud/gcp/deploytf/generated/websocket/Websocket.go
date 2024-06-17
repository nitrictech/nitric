package websocket

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/websocket/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/gcp/deploytf/generated/websocket/internal"
)

// Defines an Websocket based on a Terraform module.
//
// Source at ./.nitric/modules/websocket
type Websocket interface {
	cdktf.TerraformModule
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
	LambdaConnectTarget() *string
	SetLambdaConnectTarget(val *string)
	LambdaDisconnectTarget() *string
	SetLambdaDisconnectTarget(val *string)
	LambdaMessageTarget() *string
	SetLambdaMessageTarget(val *string)
	// The tree node.
	Node() constructs.Node
	// Experimental.
	Providers() *[]interface{}
	// Experimental.
	RawOverrides() interface{}
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	StackId() *string
	SetStackId(val *string)
	// Experimental.
	Version() *string
	WebsocketArnOutput() *string
	WebsocketExecArnOutput() *string
	WebsocketName() *string
	SetWebsocketName(val *string)
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

// The jsii proxy struct for Websocket
type jsiiProxy_Websocket struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Websocket) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) LambdaConnectTarget() *string {
	var returns *string
	_jsii_.Get(
		j,
		"lambdaConnectTarget",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) LambdaDisconnectTarget() *string {
	var returns *string
	_jsii_.Get(
		j,
		"lambdaDisconnectTarget",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) LambdaMessageTarget() *string {
	var returns *string
	_jsii_.Get(
		j,
		"lambdaMessageTarget",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) StackId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) WebsocketArnOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"websocketArnOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) WebsocketExecArnOutput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"websocketExecArnOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Websocket) WebsocketName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"websocketName",
		&returns,
	)
	return returns
}


func NewWebsocket(scope constructs.Construct, id *string, config *WebsocketConfig) Websocket {
	_init_.Initialize()

	if err := validateNewWebsocketParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Websocket{}

	_jsii_.Create(
		"websocket.Websocket",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewWebsocket_Override(w Websocket, scope constructs.Construct, id *string, config *WebsocketConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"websocket.Websocket",
		[]interface{}{scope, id, config},
		w,
	)
}

func (j *jsiiProxy_Websocket)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Websocket)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Websocket)SetLambdaConnectTarget(val *string) {
	if err := j.validateSetLambdaConnectTargetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"lambdaConnectTarget",
		val,
	)
}

func (j *jsiiProxy_Websocket)SetLambdaDisconnectTarget(val *string) {
	if err := j.validateSetLambdaDisconnectTargetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"lambdaDisconnectTarget",
		val,
	)
}

func (j *jsiiProxy_Websocket)SetLambdaMessageTarget(val *string) {
	if err := j.validateSetLambdaMessageTargetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"lambdaMessageTarget",
		val,
	)
}

func (j *jsiiProxy_Websocket)SetStackId(val *string) {
	if err := j.validateSetStackIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stackId",
		val,
	)
}

func (j *jsiiProxy_Websocket)SetWebsocketName(val *string) {
	if err := j.validateSetWebsocketNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"websocketName",
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
func Websocket_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateWebsocket_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"websocket.Websocket",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Websocket_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateWebsocket_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"websocket.Websocket",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) AddOverride(path *string, value interface{}) {
	if err := w.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		w,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (w *jsiiProxy_Websocket) AddProvider(provider interface{}) {
	if err := w.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		w,
		"addProvider",
		[]interface{}{provider},
	)
}

func (w *jsiiProxy_Websocket) GetString(output *string) *string {
	if err := w.validateGetStringParameters(output); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		w,
		"getString",
		[]interface{}{output},
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
	if err := w.validateInterpolationForOutputParameters(moduleOutput); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		w,
		"interpolationForOutput",
		[]interface{}{moduleOutput},
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) OverrideLogicalId(newLogicalId *string) {
	if err := w.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		w,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (w *jsiiProxy_Websocket) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		w,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (w *jsiiProxy_Websocket) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		w,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		w,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		w,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (w *jsiiProxy_Websocket) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		w,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

