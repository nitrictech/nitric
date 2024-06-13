package image


type ImageBuild struct {
	// Value to specify the build context.
	//
	// Currently, only a `PATH` context is supported. You can use the helper function '${path.cwd}/context-dir'. Please see https://docs.docker.com/build/building/context/ for more information about build contexts.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#context Image#context}
	Context *string `field:"required" json:"context" yaml:"context"`
	// auth_config block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#auth_config Image#auth_config}
	AuthConfig interface{} `field:"optional" json:"authConfig" yaml:"authConfig"`
	// Set build-time variables.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#build_arg Image#build_arg}
	BuildArg *map[string]*string `field:"optional" json:"buildArg" yaml:"buildArg"`
	// Pairs for build-time variables in the form TODO.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#build_args Image#build_args}
	BuildArgs *map[string]*string `field:"optional" json:"buildArgs" yaml:"buildArgs"`
	// BuildID is an optional identifier that can be passed together with the build request.
	//
	// The same identifier can be used to gracefully cancel the build with the cancel request.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#build_id Image#build_id}
	BuildId *string `field:"optional" json:"buildId" yaml:"buildId"`
	// Images to consider as cache sources.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#cache_from Image#cache_from}
	CacheFrom *[]*string `field:"optional" json:"cacheFrom" yaml:"cacheFrom"`
	// Optional parent cgroup for the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#cgroup_parent Image#cgroup_parent}
	CgroupParent *string `field:"optional" json:"cgroupParent" yaml:"cgroupParent"`
	// The length of a CPU period in microseconds.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#cpu_period Image#cpu_period}
	CpuPeriod *float64 `field:"optional" json:"cpuPeriod" yaml:"cpuPeriod"`
	// Microseconds of CPU time that the container can get in a CPU period.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#cpu_quota Image#cpu_quota}
	CpuQuota *float64 `field:"optional" json:"cpuQuota" yaml:"cpuQuota"`
	// CPUs in which to allow execution (e.g., `0-3`, `0`, `1`).
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#cpu_set_cpus Image#cpu_set_cpus}
	CpuSetCpus *string `field:"optional" json:"cpuSetCpus" yaml:"cpuSetCpus"`
	// MEMs in which to allow execution (`0-3`, `0`, `1`).
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#cpu_set_mems Image#cpu_set_mems}
	CpuSetMems *string `field:"optional" json:"cpuSetMems" yaml:"cpuSetMems"`
	// CPU shares (relative weight).
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#cpu_shares Image#cpu_shares}
	CpuShares *float64 `field:"optional" json:"cpuShares" yaml:"cpuShares"`
	// Name of the Dockerfile. Defaults to `Dockerfile`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#dockerfile Image#dockerfile}
	Dockerfile *string `field:"optional" json:"dockerfile" yaml:"dockerfile"`
	// A list of hostnames/IP mappings to add to the container’s /etc/hosts file. Specified in the form ["hostname:IP"].
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#extra_hosts Image#extra_hosts}
	ExtraHosts *[]*string `field:"optional" json:"extraHosts" yaml:"extraHosts"`
	// Always remove intermediate containers.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#force_remove Image#force_remove}
	ForceRemove interface{} `field:"optional" json:"forceRemove" yaml:"forceRemove"`
	// Isolation represents the isolation technology of a container. The supported values are.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#isolation Image#isolation}
	Isolation *string `field:"optional" json:"isolation" yaml:"isolation"`
	// Set metadata for an image.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#label Image#label}
	Label *map[string]*string `field:"optional" json:"label" yaml:"label"`
	// User-defined key/value metadata.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#labels Image#labels}
	Labels *map[string]*string `field:"optional" json:"labels" yaml:"labels"`
	// Set memory limit for build.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#memory Image#memory}
	Memory *float64 `field:"optional" json:"memory" yaml:"memory"`
	// Total memory (memory + swap), -1 to enable unlimited swap.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#memory_swap Image#memory_swap}
	MemorySwap *float64 `field:"optional" json:"memorySwap" yaml:"memorySwap"`
	// Set the networking mode for the RUN instructions during build.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#network_mode Image#network_mode}
	NetworkMode *string `field:"optional" json:"networkMode" yaml:"networkMode"`
	// Do not use the cache when building the image.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#no_cache Image#no_cache}
	NoCache interface{} `field:"optional" json:"noCache" yaml:"noCache"`
	// Set platform if server is multi-platform capable.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#platform Image#platform}
	Platform *string `field:"optional" json:"platform" yaml:"platform"`
	// Attempt to pull the image even if an older image exists locally.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#pull_parent Image#pull_parent}
	PullParent interface{} `field:"optional" json:"pullParent" yaml:"pullParent"`
	// A Git repository URI or HTTP/HTTPS context URI.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#remote_context Image#remote_context}
	RemoteContext *string `field:"optional" json:"remoteContext" yaml:"remoteContext"`
	// Remove intermediate containers after a successful build. Defaults to `true`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#remove Image#remove}
	Remove interface{} `field:"optional" json:"remove" yaml:"remove"`
	// The security options.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#security_opt Image#security_opt}
	SecurityOpt *[]*string `field:"optional" json:"securityOpt" yaml:"securityOpt"`
	// Set an ID for the build session.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#session_id Image#session_id}
	SessionId *string `field:"optional" json:"sessionId" yaml:"sessionId"`
	// Size of /dev/shm in bytes. The size must be greater than 0.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#shm_size Image#shm_size}
	ShmSize *float64 `field:"optional" json:"shmSize" yaml:"shmSize"`
	// If true the new layers are squashed into a new image with a single new layer.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#squash Image#squash}
	Squash interface{} `field:"optional" json:"squash" yaml:"squash"`
	// Suppress the build output and print image ID on success.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#suppress_output Image#suppress_output}
	SuppressOutput interface{} `field:"optional" json:"suppressOutput" yaml:"suppressOutput"`
	// Name and optionally a tag in the 'name:tag' format.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#tag Image#tag}
	Tag *[]*string `field:"optional" json:"tag" yaml:"tag"`
	// Set the target build stage to build.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#target Image#target}
	Target *string `field:"optional" json:"target" yaml:"target"`
	// ulimit block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#ulimit Image#ulimit}
	Ulimit interface{} `field:"optional" json:"ulimit" yaml:"ulimit"`
	// Version of the underlying builder to use.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/image#version Image#version}
	Version *string `field:"optional" json:"version" yaml:"version"`
}

