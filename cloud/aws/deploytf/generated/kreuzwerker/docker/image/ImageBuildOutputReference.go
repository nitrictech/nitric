package image

import (
	_jsii_ "github.com/aws/jsii-runtime-go/runtime"
	_init_ "github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/jsii"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
	"github.com/nitrictech/nitric/cloud/aws/deploytf/generated/kreuzwerker/docker/image/internal"
)

type ImageBuildOutputReference interface {
	cdktf.ComplexObject
	AuthConfig() ImageBuildAuthConfigList
	AuthConfigInput() interface{}
	BuildArg() *map[string]*string
	SetBuildArg(val *map[string]*string)
	BuildArgInput() *map[string]*string
	BuildArgs() *map[string]*string
	SetBuildArgs(val *map[string]*string)
	BuildArgsInput() *map[string]*string
	BuildId() *string
	SetBuildId(val *string)
	BuildIdInput() *string
	CacheFrom() *[]*string
	SetCacheFrom(val *[]*string)
	CacheFromInput() *[]*string
	CgroupParent() *string
	SetCgroupParent(val *string)
	CgroupParentInput() *string
	// the index of the complex object in a list.
	// Experimental.
	ComplexObjectIndex() interface{}
	// Experimental.
	SetComplexObjectIndex(val interface{})
	// set to true if this item is from inside a set and needs tolist() for accessing it set to "0" for single list items.
	// Experimental.
	ComplexObjectIsFromSet() *bool
	// Experimental.
	SetComplexObjectIsFromSet(val *bool)
	Context() *string
	SetContext(val *string)
	ContextInput() *string
	CpuPeriod() *float64
	SetCpuPeriod(val *float64)
	CpuPeriodInput() *float64
	CpuQuota() *float64
	SetCpuQuota(val *float64)
	CpuQuotaInput() *float64
	CpuSetCpus() *string
	SetCpuSetCpus(val *string)
	CpuSetCpusInput() *string
	CpuSetMems() *string
	SetCpuSetMems(val *string)
	CpuSetMemsInput() *string
	CpuShares() *float64
	SetCpuShares(val *float64)
	CpuSharesInput() *float64
	// The creation stack of this resolvable which will be appended to errors thrown during resolution.
	//
	// If this returns an empty array the stack will not be attached.
	// Experimental.
	CreationStack() *[]*string
	Dockerfile() *string
	SetDockerfile(val *string)
	DockerfileInput() *string
	ExtraHosts() *[]*string
	SetExtraHosts(val *[]*string)
	ExtraHostsInput() *[]*string
	ForceRemove() interface{}
	SetForceRemove(val interface{})
	ForceRemoveInput() interface{}
	// Experimental.
	Fqn() *string
	InternalValue() *ImageBuild
	SetInternalValue(val *ImageBuild)
	Isolation() *string
	SetIsolation(val *string)
	IsolationInput() *string
	Label() *map[string]*string
	SetLabel(val *map[string]*string)
	LabelInput() *map[string]*string
	Labels() *map[string]*string
	SetLabels(val *map[string]*string)
	LabelsInput() *map[string]*string
	Memory() *float64
	SetMemory(val *float64)
	MemoryInput() *float64
	MemorySwap() *float64
	SetMemorySwap(val *float64)
	MemorySwapInput() *float64
	NetworkMode() *string
	SetNetworkMode(val *string)
	NetworkModeInput() *string
	NoCache() interface{}
	SetNoCache(val interface{})
	NoCacheInput() interface{}
	Platform() *string
	SetPlatform(val *string)
	PlatformInput() *string
	PullParent() interface{}
	SetPullParent(val interface{})
	PullParentInput() interface{}
	RemoteContext() *string
	SetRemoteContext(val *string)
	RemoteContextInput() *string
	Remove() interface{}
	SetRemove(val interface{})
	RemoveInput() interface{}
	SecurityOpt() *[]*string
	SetSecurityOpt(val *[]*string)
	SecurityOptInput() *[]*string
	SessionId() *string
	SetSessionId(val *string)
	SessionIdInput() *string
	ShmSize() *float64
	SetShmSize(val *float64)
	ShmSizeInput() *float64
	Squash() interface{}
	SetSquash(val interface{})
	SquashInput() interface{}
	SuppressOutput() interface{}
	SetSuppressOutput(val interface{})
	SuppressOutputInput() interface{}
	Tag() *[]*string
	SetTag(val *[]*string)
	TagInput() *[]*string
	Target() *string
	SetTarget(val *string)
	TargetInput() *string
	// Experimental.
	TerraformAttribute() *string
	// Experimental.
	SetTerraformAttribute(val *string)
	// Experimental.
	TerraformResource() cdktf.IInterpolatingParent
	// Experimental.
	SetTerraformResource(val cdktf.IInterpolatingParent)
	Ulimit() ImageBuildUlimitList
	UlimitInput() interface{}
	Version() *string
	SetVersion(val *string)
	VersionInput() *string
	// Experimental.
	ComputeFqn() *string
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
	InterpolationAsList() cdktf.IResolvable
	// Experimental.
	InterpolationForAttribute(property *string) cdktf.IResolvable
	PutAuthConfig(value interface{})
	PutUlimit(value interface{})
	ResetAuthConfig()
	ResetBuildArg()
	ResetBuildArgs()
	ResetBuildId()
	ResetCacheFrom()
	ResetCgroupParent()
	ResetCpuPeriod()
	ResetCpuQuota()
	ResetCpuSetCpus()
	ResetCpuSetMems()
	ResetCpuShares()
	ResetDockerfile()
	ResetExtraHosts()
	ResetForceRemove()
	ResetIsolation()
	ResetLabel()
	ResetLabels()
	ResetMemory()
	ResetMemorySwap()
	ResetNetworkMode()
	ResetNoCache()
	ResetPlatform()
	ResetPullParent()
	ResetRemoteContext()
	ResetRemove()
	ResetSecurityOpt()
	ResetSessionId()
	ResetShmSize()
	ResetSquash()
	ResetSuppressOutput()
	ResetTag()
	ResetTarget()
	ResetUlimit()
	ResetVersion()
	// Produce the Token's value at resolution time.
	// Experimental.
	Resolve(_context cdktf.IResolveContext) interface{}
	// Return a string representation of this resolvable object.
	//
	// Returns a reversible string representation.
	// Experimental.
	ToString() *string
}

// The jsii proxy struct for ImageBuildOutputReference
type jsiiProxy_ImageBuildOutputReference struct {
	internal.Type__cdktfComplexObject
}

func (j *jsiiProxy_ImageBuildOutputReference) AuthConfig() ImageBuildAuthConfigList {
	var returns ImageBuildAuthConfigList
	_jsii_.Get(
		j,
		"authConfig",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) AuthConfigInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"authConfigInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) BuildArg() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"buildArg",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) BuildArgInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"buildArgInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) BuildArgs() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"buildArgs",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) BuildArgsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"buildArgsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) BuildId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"buildId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) BuildIdInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"buildIdInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CacheFrom() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"cacheFrom",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CacheFromInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"cacheFromInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CgroupParent() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cgroupParent",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CgroupParentInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cgroupParentInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ComplexObjectIndex() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"complexObjectIndex",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ComplexObjectIsFromSet() *bool {
	var returns *bool
	_jsii_.Get(
		j,
		"complexObjectIsFromSet",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Context() *string {
	var returns *string
	_jsii_.Get(
		j,
		"context",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ContextInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"contextInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuPeriod() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuPeriod",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuPeriodInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuPeriodInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuQuota() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuQuota",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuQuotaInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuQuotaInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuSetCpus() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cpuSetCpus",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuSetCpusInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cpuSetCpusInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuSetMems() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cpuSetMems",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuSetMemsInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"cpuSetMemsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuShares() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuShares",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CpuSharesInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"cpuSharesInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) CreationStack() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"creationStack",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Dockerfile() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dockerfile",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) DockerfileInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"dockerfileInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ExtraHosts() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"extraHosts",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ExtraHostsInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"extraHostsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ForceRemove() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"forceRemove",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ForceRemoveInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"forceRemoveInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Fqn() *string {
	var returns *string
	_jsii_.Get(
		j,
		"fqn",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) InternalValue() *ImageBuild {
	var returns *ImageBuild
	_jsii_.Get(
		j,
		"internalValue",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Isolation() *string {
	var returns *string
	_jsii_.Get(
		j,
		"isolation",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) IsolationInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"isolationInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Label() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"label",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) LabelInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"labelInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Labels() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"labels",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) LabelsInput() *map[string]*string {
	var returns *map[string]*string
	_jsii_.Get(
		j,
		"labelsInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Memory() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memory",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) MemoryInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memoryInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) MemorySwap() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memorySwap",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) MemorySwapInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"memorySwapInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) NetworkMode() *string {
	var returns *string
	_jsii_.Get(
		j,
		"networkMode",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) NetworkModeInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"networkModeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) NoCache() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"noCache",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) NoCacheInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"noCacheInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Platform() *string {
	var returns *string
	_jsii_.Get(
		j,
		"platform",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) PlatformInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"platformInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) PullParent() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"pullParent",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) PullParentInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"pullParentInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) RemoteContext() *string {
	var returns *string
	_jsii_.Get(
		j,
		"remoteContext",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) RemoteContextInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"remoteContextInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Remove() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"remove",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) RemoveInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"removeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) SecurityOpt() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityOpt",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) SecurityOptInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"securityOptInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) SessionId() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sessionId",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) SessionIdInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"sessionIdInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ShmSize() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"shmSize",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) ShmSizeInput() *float64 {
	var returns *float64
	_jsii_.Get(
		j,
		"shmSizeInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Squash() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"squash",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) SquashInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"squashInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) SuppressOutput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"suppressOutput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) SuppressOutputInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"suppressOutputInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Tag() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"tag",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) TagInput() *[]*string {
	var returns *[]*string
	_jsii_.Get(
		j,
		"tagInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Target() *string {
	var returns *string
	_jsii_.Get(
		j,
		"target",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) TargetInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"targetInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) TerraformAttribute() *string {
	var returns *string
	_jsii_.Get(
		j,
		"terraformAttribute",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) TerraformResource() cdktf.IInterpolatingParent {
	var returns cdktf.IInterpolatingParent
	_jsii_.Get(
		j,
		"terraformResource",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Ulimit() ImageBuildUlimitList {
	var returns ImageBuildUlimitList
	_jsii_.Get(
		j,
		"ulimit",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) UlimitInput() interface{} {
	var returns interface{}
	_jsii_.Get(
		j,
		"ulimitInput",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) Version() *string {
	var returns *string
	_jsii_.Get(
		j,
		"version",
		&returns,
	)
	return returns
}

func (j *jsiiProxy_ImageBuildOutputReference) VersionInput() *string {
	var returns *string
	_jsii_.Get(
		j,
		"versionInput",
		&returns,
	)
	return returns
}


func NewImageBuildOutputReference(terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) ImageBuildOutputReference {
	_init_.Initialize()

	if err := validateNewImageBuildOutputReferenceParameters(terraformResource, terraformAttribute); err != nil {
		panic(err)
	}
	j := jsiiProxy_ImageBuildOutputReference{}

	_jsii_.Create(
		"docker.image.ImageBuildOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		&j,
	)

	return &j
}

func NewImageBuildOutputReference_Override(i ImageBuildOutputReference, terraformResource cdktf.IInterpolatingParent, terraformAttribute *string) {
	_init_.Initialize()

	_jsii_.Create(
		"docker.image.ImageBuildOutputReference",
		[]interface{}{terraformResource, terraformAttribute},
		i,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetBuildArg(val *map[string]*string) {
	if err := j.validateSetBuildArgParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"buildArg",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetBuildArgs(val *map[string]*string) {
	if err := j.validateSetBuildArgsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"buildArgs",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetBuildId(val *string) {
	if err := j.validateSetBuildIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"buildId",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetCacheFrom(val *[]*string) {
	if err := j.validateSetCacheFromParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cacheFrom",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetCgroupParent(val *string) {
	if err := j.validateSetCgroupParentParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cgroupParent",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetComplexObjectIndex(val interface{}) {
	if err := j.validateSetComplexObjectIndexParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIndex",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetComplexObjectIsFromSet(val *bool) {
	if err := j.validateSetComplexObjectIsFromSetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"complexObjectIsFromSet",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetContext(val *string) {
	if err := j.validateSetContextParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"context",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetCpuPeriod(val *float64) {
	if err := j.validateSetCpuPeriodParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cpuPeriod",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetCpuQuota(val *float64) {
	if err := j.validateSetCpuQuotaParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cpuQuota",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetCpuSetCpus(val *string) {
	if err := j.validateSetCpuSetCpusParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cpuSetCpus",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetCpuSetMems(val *string) {
	if err := j.validateSetCpuSetMemsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cpuSetMems",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetCpuShares(val *float64) {
	if err := j.validateSetCpuSharesParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"cpuShares",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetDockerfile(val *string) {
	if err := j.validateSetDockerfileParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"dockerfile",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetExtraHosts(val *[]*string) {
	if err := j.validateSetExtraHostsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"extraHosts",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetForceRemove(val interface{}) {
	if err := j.validateSetForceRemoveParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"forceRemove",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetInternalValue(val *ImageBuild) {
	if err := j.validateSetInternalValueParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"internalValue",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetIsolation(val *string) {
	if err := j.validateSetIsolationParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"isolation",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetLabel(val *map[string]*string) {
	if err := j.validateSetLabelParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"label",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetLabels(val *map[string]*string) {
	if err := j.validateSetLabelsParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"labels",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetMemory(val *float64) {
	if err := j.validateSetMemoryParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"memory",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetMemorySwap(val *float64) {
	if err := j.validateSetMemorySwapParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"memorySwap",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetNetworkMode(val *string) {
	if err := j.validateSetNetworkModeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"networkMode",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetNoCache(val interface{}) {
	if err := j.validateSetNoCacheParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"noCache",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetPlatform(val *string) {
	if err := j.validateSetPlatformParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"platform",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetPullParent(val interface{}) {
	if err := j.validateSetPullParentParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"pullParent",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetRemoteContext(val *string) {
	if err := j.validateSetRemoteContextParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"remoteContext",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetRemove(val interface{}) {
	if err := j.validateSetRemoveParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"remove",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetSecurityOpt(val *[]*string) {
	if err := j.validateSetSecurityOptParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"securityOpt",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetSessionId(val *string) {
	if err := j.validateSetSessionIdParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"sessionId",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetShmSize(val *float64) {
	if err := j.validateSetShmSizeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"shmSize",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetSquash(val interface{}) {
	if err := j.validateSetSquashParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"squash",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetSuppressOutput(val interface{}) {
	if err := j.validateSetSuppressOutputParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"suppressOutput",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetTag(val *[]*string) {
	if err := j.validateSetTagParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"tag",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetTarget(val *string) {
	if err := j.validateSetTargetParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"target",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetTerraformAttribute(val *string) {
	if err := j.validateSetTerraformAttributeParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformAttribute",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetTerraformResource(val cdktf.IInterpolatingParent) {
	if err := j.validateSetTerraformResourceParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"terraformResource",
		val,
	)
}

func (j *jsiiProxy_ImageBuildOutputReference)SetVersion(val *string) {
	if err := j.validateSetVersionParameters(val); err != nil {
		panic(err)
	}
	_jsii_.Set(
		j,
		"version",
		val,
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ComputeFqn() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"computeFqn",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetAnyMapAttribute(terraformAttribute *string) *map[string]interface{} {
	if err := i.validateGetAnyMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]interface{}

	_jsii_.Invoke(
		i,
		"getAnyMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetBooleanAttribute(terraformAttribute *string) cdktf.IResolvable {
	if err := i.validateGetBooleanAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		i,
		"getBooleanAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetBooleanMapAttribute(terraformAttribute *string) *map[string]*bool {
	if err := i.validateGetBooleanMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*bool

	_jsii_.Invoke(
		i,
		"getBooleanMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetListAttribute(terraformAttribute *string) *[]*string {
	if err := i.validateGetListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*string

	_jsii_.Invoke(
		i,
		"getListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetNumberAttribute(terraformAttribute *string) *float64 {
	if err := i.validateGetNumberAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *float64

	_jsii_.Invoke(
		i,
		"getNumberAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetNumberListAttribute(terraformAttribute *string) *[]*float64 {
	if err := i.validateGetNumberListAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *[]*float64

	_jsii_.Invoke(
		i,
		"getNumberListAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetNumberMapAttribute(terraformAttribute *string) *map[string]*float64 {
	if err := i.validateGetNumberMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*float64

	_jsii_.Invoke(
		i,
		"getNumberMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetStringAttribute(terraformAttribute *string) *string {
	if err := i.validateGetStringAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *string

	_jsii_.Invoke(
		i,
		"getStringAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) GetStringMapAttribute(terraformAttribute *string) *map[string]*string {
	if err := i.validateGetStringMapAttributeParameters(terraformAttribute); err != nil {
		panic(err)
	}
	var returns *map[string]*string

	_jsii_.Invoke(
		i,
		"getStringMapAttribute",
		[]interface{}{terraformAttribute},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) InterpolationAsList() cdktf.IResolvable {
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		i,
		"interpolationAsList",
		nil, // no parameters
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) InterpolationForAttribute(property *string) cdktf.IResolvable {
	if err := i.validateInterpolationForAttributeParameters(property); err != nil {
		panic(err)
	}
	var returns cdktf.IResolvable

	_jsii_.Invoke(
		i,
		"interpolationForAttribute",
		[]interface{}{property},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) PutAuthConfig(value interface{}) {
	if err := i.validatePutAuthConfigParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		i,
		"putAuthConfig",
		[]interface{}{value},
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) PutUlimit(value interface{}) {
	if err := i.validatePutUlimitParameters(value); err != nil {
		panic(err)
	}
	_jsii_.InvokeVoid(
		i,
		"putUlimit",
		[]interface{}{value},
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetAuthConfig() {
	_jsii_.InvokeVoid(
		i,
		"resetAuthConfig",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetBuildArg() {
	_jsii_.InvokeVoid(
		i,
		"resetBuildArg",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetBuildArgs() {
	_jsii_.InvokeVoid(
		i,
		"resetBuildArgs",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetBuildId() {
	_jsii_.InvokeVoid(
		i,
		"resetBuildId",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetCacheFrom() {
	_jsii_.InvokeVoid(
		i,
		"resetCacheFrom",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetCgroupParent() {
	_jsii_.InvokeVoid(
		i,
		"resetCgroupParent",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetCpuPeriod() {
	_jsii_.InvokeVoid(
		i,
		"resetCpuPeriod",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetCpuQuota() {
	_jsii_.InvokeVoid(
		i,
		"resetCpuQuota",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetCpuSetCpus() {
	_jsii_.InvokeVoid(
		i,
		"resetCpuSetCpus",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetCpuSetMems() {
	_jsii_.InvokeVoid(
		i,
		"resetCpuSetMems",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetCpuShares() {
	_jsii_.InvokeVoid(
		i,
		"resetCpuShares",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetDockerfile() {
	_jsii_.InvokeVoid(
		i,
		"resetDockerfile",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetExtraHosts() {
	_jsii_.InvokeVoid(
		i,
		"resetExtraHosts",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetForceRemove() {
	_jsii_.InvokeVoid(
		i,
		"resetForceRemove",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetIsolation() {
	_jsii_.InvokeVoid(
		i,
		"resetIsolation",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetLabel() {
	_jsii_.InvokeVoid(
		i,
		"resetLabel",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetLabels() {
	_jsii_.InvokeVoid(
		i,
		"resetLabels",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetMemory() {
	_jsii_.InvokeVoid(
		i,
		"resetMemory",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetMemorySwap() {
	_jsii_.InvokeVoid(
		i,
		"resetMemorySwap",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetNetworkMode() {
	_jsii_.InvokeVoid(
		i,
		"resetNetworkMode",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetNoCache() {
	_jsii_.InvokeVoid(
		i,
		"resetNoCache",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetPlatform() {
	_jsii_.InvokeVoid(
		i,
		"resetPlatform",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetPullParent() {
	_jsii_.InvokeVoid(
		i,
		"resetPullParent",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetRemoteContext() {
	_jsii_.InvokeVoid(
		i,
		"resetRemoteContext",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetRemove() {
	_jsii_.InvokeVoid(
		i,
		"resetRemove",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetSecurityOpt() {
	_jsii_.InvokeVoid(
		i,
		"resetSecurityOpt",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetSessionId() {
	_jsii_.InvokeVoid(
		i,
		"resetSessionId",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetShmSize() {
	_jsii_.InvokeVoid(
		i,
		"resetShmSize",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetSquash() {
	_jsii_.InvokeVoid(
		i,
		"resetSquash",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetSuppressOutput() {
	_jsii_.InvokeVoid(
		i,
		"resetSuppressOutput",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetTag() {
	_jsii_.InvokeVoid(
		i,
		"resetTag",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetTarget() {
	_jsii_.InvokeVoid(
		i,
		"resetTarget",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetUlimit() {
	_jsii_.InvokeVoid(
		i,
		"resetUlimit",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) ResetVersion() {
	_jsii_.InvokeVoid(
		i,
		"resetVersion",
		nil, // no parameters
	)
}

func (i *jsiiProxy_ImageBuildOutputReference) Resolve(_context cdktf.IResolveContext) interface{} {
	if err := i.validateResolveParameters(_context); err != nil {
		panic(err)
	}
	var returns interface{}

	_jsii_.Invoke(
		i,
		"resolve",
		[]interface{}{_context},
		&returns,
	)

	return returns
}

func (i *jsiiProxy_ImageBuildOutputReference) ToString() *string {
	var returns *string

	_jsii_.Invoke(
		i,
		"toString",
		nil, // no parameters
		&returns,
	)

	return returns
}

