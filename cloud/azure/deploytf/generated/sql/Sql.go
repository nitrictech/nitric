package sql

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/azure/deploytf/generated/sql/jsii"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/azure/deploytf/generated/sql/internal"
)

// Defines an Sql based on a Terraform module.
//
// Source at ./.nitric/modules/sql
type Sql interface {
	cdktf.TerraformModule
	// Experimental.
	CdktfStack() cdktf.TerraformStack
	// Experimental.
	ConstructNodeMetadata() *map[string]interface{}
	DatabaseMasterPassword() *string
	SetDatabaseMasterPassword(val *string)
	DatabaseServerFqdn() *string
	SetDatabaseServerFqdn(val *string)
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
	ImageRegistryPassword() *string
	SetImageRegistryPassword(val *string)
	ImageRegistryServer() *string
	SetImageRegistryServer(val *string)
	ImageRegistryUsername() *string
	SetImageRegistryUsername(val *string)
	Location() *string
	SetLocation(val *string)
	MigrationContainerSubnetId() *string
	SetMigrationContainerSubnetId(val *string)
	MigrationImageUri() *string
	SetMigrationImageUri(val *string)
	Name() *string
	SetName(val *string)
	// The tree node.
	Node() constructs.Node
	// Experimental.
	Providers() *[]interface{}
	// Experimental.
	RawOverrides() interface{}
	ResourceGroupName() *string
	SetResourceGroupName(val *string)
	ServerId() *string
	SetServerId(val *string)
	// Experimental.
	SkipAssetCreationFromLocalModules() *bool
	// Experimental.
	Source() *string
	StackId() *string
	SetStackId(val *string)
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

// The jsii proxy struct for Sql
type jsiiProxy_Sql struct {
	internal.Type__cdktfTerraformModule
}

func (j *jsiiProxy_Sql) CdktfStack() cdktf.TerraformStack {
	var returns cdktf.TerraformStack
	_jsii_.Get(
		j,
		"cdktfStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ConstructNodeMetadata() *map[string]interface{} {
	var returns *map[string]interface{}
	_jsii_.Get(
		j,
		"constructNodeMetadata",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) DatabaseMasterPassword() *string {
	var returns *string
	_jsii_.Get(
		j,
		"databaseMasterPassword",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) DatabaseServerFqdn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"databaseServerFqdn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) DependsOn() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"dependsOn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ForEach() cdktf.ITerraformIterator {
	var returns cdktf.ITerraformIterator
	_jsii_.Get(
		j,
		"forEach",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) FriendlyUniqueId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"friendlyUniqueId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ImageRegistryPassword() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageRegistryPassword",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ImageRegistryServer() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageRegistryServer",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ImageRegistryUsername() *string {
	var returns *string
	_jsii_.Get(
		j,
		"imageRegistryUsername",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Location() *string {
	var returns *string
	_jsii_.Get(
		j,
		"location",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) MigrationContainerSubnetId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"migrationContainerSubnetId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) MigrationImageUri() *string {
	var returns *string
	_jsii_.Get(
		j,
		"migrationImageUri",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Name() *string {
	var returns *string
	_jsii_.Get(
		j,
		"name",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Node() constructs.Node {
	var returns constructs.Node
	_jsii_.Get(
		j,
		"node",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Providers() *[]interface{} {
	var returns *[]interface{}
	_jsii_.Get(
		j,
		"providers",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) RawOverrides() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"rawOverrides",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ResourceGroupName() *string {
	var returns *string
	_jsii_.Get(
		j,
		"resourceGroupName",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) ServerId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"serverId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) SkipAssetCreationFromLocalModules() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"skipAssetCreationFromLocalModules",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Source() *string {
	var returns *string
	_jsii_.Get(
		j,
		"source",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) StackId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"stackId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_Sql) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}


func NewSql(scope constructs.Construct, id *string, config *SqlConfig) Sql {
	_init_.Initialize()

	if err := validateNewSqlParameters(scope, id, config); err != nil {
		panic(err)
	}
	j := jsiiProxy_Sql{}

	_jsii_.Create(
		"sql.Sql",
		[]interface{}{scope, id, config},
		&j,
	)

	return &j
}

func NewSql_Override(s Sql, scope constructs.Construct, id *string, config *SqlConfig) {
	_init_.Initialize()

	_jsii_.Create(
		"sql.Sql",
		[]interface{}{scope, id, config},
		s,
	)
}

func (j *jsiiProxy_Sql)SetDatabaseMasterPassword(val *string) {
	if err := j.validateSetDatabaseMasterPasswordParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"databaseMasterPassword",
		val,
	)
}

func (j *jsiiProxy_Sql)SetDatabaseServerFqdn(val *string) {
	if err := j.validateSetDatabaseServerFqdnParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"databaseServerFqdn",
		val,
	)
}

func (j *jsiiProxy_Sql)SetDependsOn(val *[]*string) {
	_jsii_.Set(
		j,
		"dependsOn",
		val,
	)
}

func (j *jsiiProxy_Sql)SetForEach(val cdktf.ITerraformIterator) {
	_jsii_.Set(
		j,
		"forEach",
		val,
	)
}

func (j *jsiiProxy_Sql)SetImageRegistryPassword(val *string) {
	if err := j.validateSetImageRegistryPasswordParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"imageRegistryPassword",
		val,
	)
}

func (j *jsiiProxy_Sql)SetImageRegistryServer(val *string) {
	if err := j.validateSetImageRegistryServerParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"imageRegistryServer",
		val,
	)
}

func (j *jsiiProxy_Sql)SetImageRegistryUsername(val *string) {
	if err := j.validateSetImageRegistryUsernameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"imageRegistryUsername",
		val,
	)
}

func (j *jsiiProxy_Sql)SetLocation(val *string) {
	if err := j.validateSetLocationParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"location",
		val,
	)
}

func (j *jsiiProxy_Sql)SetMigrationContainerSubnetId(val *string) {
	if err := j.validateSetMigrationContainerSubnetIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"migrationContainerSubnetId",
		val,
	)
}

func (j *jsiiProxy_Sql)SetMigrationImageUri(val *string) {
	if err := j.validateSetMigrationImageUriParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"migrationImageUri",
		val,
	)
}

func (j *jsiiProxy_Sql)SetName(val *string) {
	if err := j.validateSetNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"name",
		val,
	)
}

func (j *jsiiProxy_Sql)SetResourceGroupName(val *string) {
	if err := j.validateSetResourceGroupNameParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"resourceGroupName",
		val,
	)
}

func (j *jsiiProxy_Sql)SetServerId(val *string) {
	if err := j.validateSetServerIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"serverId",
		val,
	)
}

func (j *jsiiProxy_Sql)SetStackId(val *string) {
	if err := j.validateSetStackIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"stackId",
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
func Sql_IsConstruct(x interface{}) *bool {
	_init_.Initialize()

	if err := validateSql_IsConstructParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"sql.Sql",
		"isConstruct",
		[]interface{}{x},
		&returns,
	)

	return returns
}

// Experimental.
func Sql_IsTerraformElement(x interface{}) *bool {
	_init_.Initialize()

	if err := validateSql_IsTerraformElementParameters(x); err != nil {
		panic(err)
	}
	var returns *bool

	_jsii_.StaticInvoke(
		"sql.Sql",
		"isTerraformElement",
		[]interface{}{x},
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) AddOverride(path *string, value interface{}) {
	if err := s.validateAddOverrideParameters(path, value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addOverride",
		[]interface{}{path, value},
	)
}

func (s *jsiiProxy_Sql) AddProvider(provider interface{}) {
	if err := s.validateAddProviderParameters(provider); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"addProvider",
		[]interface{}{provider},
	)
}

func (s *jsiiProxy_Sql) GetString(output *string) *string {
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

func (s *jsiiProxy_Sql) InterpolationForOutput(moduleOutput *string) cdktf.IResolvable {
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

func (s *jsiiProxy_Sql) OverrideLogicalId(newLogicalId *string) {
	if err := s.validateOverrideLogicalIdParameters(newLogicalId); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		s,
		"overrideLogicalId",
		[]interface{}{newLogicalId},
	)
}

func (s *jsiiProxy_Sql) ResetOverrideLogicalId() {
	_jsii_.InvokeVoid(
		s,
		"resetOverrideLogicalId",
		nil, // no parameters
	)
}

func (s *jsiiProxy_Sql) SynthesizeAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) SynthesizeHclAttributes() *map[string]interface{} {
	var returns *map[string]interface{}

	_jsii_.Invoke(
		s,
		"synthesizeHclAttributes",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToHclTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toHclTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToMetadata() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toMetadata",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		s,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (s *jsiiProxy_Sql) ToTerraform() interface{} {
	var returns interface{}

	_jsii_.Invoke(
		s,
		"toTerraform",
		nil, // no parameters
		&returns,
	)

	return returns
}

