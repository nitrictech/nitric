package container

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type ContainerConfig struct {
	// Experimental.
	Connection interface{} `field:"optional" json:"connection" yaml:"connection"`
	// Experimental.
	Count interface{} `field:"optional" json:"count" yaml:"count"`
	// Experimental.
	DependsOn *[]cdktf.ITerraformDependable `field:"optional" json:"dependsOn" yaml:"dependsOn"`
	// Experimental.
	ForEach cdktf.ITerraformIterator `field:"optional" json:"forEach" yaml:"forEach"`
	// Experimental.
	Lifecycle *cdktf.TerraformResourceLifecycle `field:"optional" json:"lifecycle" yaml:"lifecycle"`
	// Experimental.
	Provider cdktf.TerraformProvider `field:"optional" json:"provider" yaml:"provider"`
	// Experimental.
	Provisioners *[]interface{} `field:"optional" json:"provisioners" yaml:"provisioners"`
	// The ID of the image to back this container.
	//
	// The easiest way to get this value is to use the `docker_image` resource as is shown in the example.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#image Container#image}
	Image *string `field:"required" json:"image" yaml:"image"`
	// The name of the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#name Container#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// If `true` attach to the container after its creation and waits the end of its execution. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#attach Container#attach}
	Attach interface{} `field:"optional" json:"attach" yaml:"attach"`
	// capabilities block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#capabilities Container#capabilities}
	Capabilities *ContainerCapabilities `field:"optional" json:"capabilities" yaml:"capabilities"`
	// Cgroup namespace mode to use for the container. Possible values are: `private`, `host`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#cgroupns_mode Container#cgroupns_mode}
	CgroupnsMode *string `field:"optional" json:"cgroupnsMode" yaml:"cgroupnsMode"`
	// The command to use to start the container.
	//
	// For example, to run `/usr/bin/myprogram -f baz.conf` set the command to be `["/usr/bin/myprogram","-f","baz.con"]`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#command Container#command}
	Command *[]*string `field:"optional" json:"command" yaml:"command"`
	// The total number of milliseconds to wait for the container to reach status 'running'.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#container_read_refresh_timeout_milliseconds Container#container_read_refresh_timeout_milliseconds}
	ContainerReadRefreshTimeoutMilliseconds *float64 `field:"optional" json:"containerReadRefreshTimeoutMilliseconds" yaml:"containerReadRefreshTimeoutMilliseconds"`
	// A comma-separated list or hyphen-separated range of CPUs a container can use, e.g. `0-1`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#cpu_set Container#cpu_set}
	CpuSet *string `field:"optional" json:"cpuSet" yaml:"cpuSet"`
	// CPU shares (relative weight) for the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#cpu_shares Container#cpu_shares}
	CpuShares *float64 `field:"optional" json:"cpuShares" yaml:"cpuShares"`
	// If defined will attempt to stop the container before destroying.
	//
	// Container will be destroyed after `n` seconds or on successful stop.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#destroy_grace_seconds Container#destroy_grace_seconds}
	DestroyGraceSeconds *float64 `field:"optional" json:"destroyGraceSeconds" yaml:"destroyGraceSeconds"`
	// devices block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#devices Container#devices}
	Devices interface{} `field:"optional" json:"devices" yaml:"devices"`
	// DNS servers to use.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#dns Container#dns}
	Dns *[]*string `field:"optional" json:"dns" yaml:"dns"`
	// DNS options used by the DNS provider(s), see `resolv.conf` documentation for valid list of options.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#dns_opts Container#dns_opts}
	DnsOpts *[]*string `field:"optional" json:"dnsOpts" yaml:"dnsOpts"`
	// DNS search domains that are used when bare unqualified hostnames are used inside of the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#dns_search Container#dns_search}
	DnsSearch *[]*string `field:"optional" json:"dnsSearch" yaml:"dnsSearch"`
	// Domain name of the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#domainname Container#domainname}
	Domainname *string `field:"optional" json:"domainname" yaml:"domainname"`
	// The command to use as the Entrypoint for the container.
	//
	// The Entrypoint allows you to configure a container to run as an executable. For example, to run `/usr/bin/myprogram` when starting a container, set the entrypoint to be `"/usr/bin/myprogra"]`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#entrypoint Container#entrypoint}
	Entrypoint *[]*string `field:"optional" json:"entrypoint" yaml:"entrypoint"`
	// Environment variables to set in the form of `KEY=VALUE`, e.g. `DEBUG=0`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#env Container#env}
	Env *[]*string `field:"optional" json:"env" yaml:"env"`
	// GPU devices to add to the container.
	//
	// Currently, only the value `all` is supported. Passing any other value will result in unexpected behavior.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#gpus Container#gpus}
	Gpus *string `field:"optional" json:"gpus" yaml:"gpus"`
	// Additional groups for the container user.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#group_add Container#group_add}
	GroupAdd *[]*string `field:"optional" json:"groupAdd" yaml:"groupAdd"`
	// healthcheck block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#healthcheck Container#healthcheck}
	Healthcheck *ContainerHealthcheck `field:"optional" json:"healthcheck" yaml:"healthcheck"`
	// host block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#host Container#host}
	Host interface{} `field:"optional" json:"host" yaml:"host"`
	// Hostname of the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#hostname Container#hostname}
	Hostname *string `field:"optional" json:"hostname" yaml:"hostname"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#id Container#id}.
	//
	// Please be aware that the id field is automatically added to all resources in Terraform providers using a Terraform provider SDK version below 2.
	// If you experience problems setting this value it might not be settable. Please take a look at the provider documentation to ensure it should be settable.
	Id *string `field:"optional" json:"id" yaml:"id"`
	// Configured whether an init process should be injected for this container.
	//
	// If unset this will default to the `dockerd` defaults.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#init Container#init}
	Init interface{} `field:"optional" json:"init" yaml:"init"`
	// IPC sharing mode for the container. Possible values are: `none`, `private`, `shareable`, `container:<name|id>` or `host`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#ipc_mode Container#ipc_mode}
	IpcMode *string `field:"optional" json:"ipcMode" yaml:"ipcMode"`
	// labels block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#labels Container#labels}
	Labels interface{} `field:"optional" json:"labels" yaml:"labels"`
	// The logging driver to use for the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#log_driver Container#log_driver}
	LogDriver *string `field:"optional" json:"logDriver" yaml:"logDriver"`
	// Key/value pairs to use as options for the logging driver.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#log_opts Container#log_opts}
	LogOpts *map[string]*string `field:"optional" json:"logOpts" yaml:"logOpts"`
	// Save the container logs (`attach` must be enabled). Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#logs Container#logs}
	Logs interface{} `field:"optional" json:"logs" yaml:"logs"`
	// The maximum amount of times to an attempt a restart when `restart` is set to 'on-failure'.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#max_retry_count Container#max_retry_count}
	MaxRetryCount *float64 `field:"optional" json:"maxRetryCount" yaml:"maxRetryCount"`
	// The memory limit for the container in MBs.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#memory Container#memory}
	Memory *float64 `field:"optional" json:"memory" yaml:"memory"`
	// The total memory limit (memory + swap) for the container in MBs.
	//
	// This setting may compute to `-1` after `terraform apply` if the target host doesn't support memory swap, when that is the case docker will use a soft limitation.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#memory_swap Container#memory_swap}
	MemorySwap *float64 `field:"optional" json:"memorySwap" yaml:"memorySwap"`
	// mounts block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#mounts Container#mounts}
	Mounts interface{} `field:"optional" json:"mounts" yaml:"mounts"`
	// If `true`, then the Docker container will be kept running.
	//
	// If `false`, then as long as the container exists, Terraform assumes it is successful. Defaults to `true`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#must_run Container#must_run}
	MustRun interface{} `field:"optional" json:"mustRun" yaml:"mustRun"`
	// Network mode of the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#network_mode Container#network_mode}
	NetworkMode *string `field:"optional" json:"networkMode" yaml:"networkMode"`
	// networks_advanced block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#networks_advanced Container#networks_advanced}
	NetworksAdvanced interface{} `field:"optional" json:"networksAdvanced" yaml:"networksAdvanced"`
	// he PID (Process) Namespace mode for the container. Either `container:<name|id>` or `host`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#pid_mode Container#pid_mode}
	PidMode *string `field:"optional" json:"pidMode" yaml:"pidMode"`
	// ports block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#ports Container#ports}
	Ports interface{} `field:"optional" json:"ports" yaml:"ports"`
	// If `true`, the container runs in privileged mode.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#privileged Container#privileged}
	Privileged interface{} `field:"optional" json:"privileged" yaml:"privileged"`
	// Publish all ports of the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#publish_all_ports Container#publish_all_ports}
	PublishAllPorts interface{} `field:"optional" json:"publishAllPorts" yaml:"publishAllPorts"`
	// If `true`, the container will be started as readonly. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#read_only Container#read_only}
	ReadOnly interface{} `field:"optional" json:"readOnly" yaml:"readOnly"`
	// If `true`, it will remove anonymous volumes associated with the container. Defaults to `true`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#remove_volumes Container#remove_volumes}
	RemoveVolumes interface{} `field:"optional" json:"removeVolumes" yaml:"removeVolumes"`
	// The restart policy for the container. Must be one of 'no', 'on-failure', 'always', 'unless-stopped'. Defaults to `no`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#restart Container#restart}
	Restart *string `field:"optional" json:"restart" yaml:"restart"`
	// If `true`, then the container will be automatically removed when it exits. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#rm Container#rm}
	Rm interface{} `field:"optional" json:"rm" yaml:"rm"`
	// Runtime to use for the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#runtime Container#runtime}
	Runtime *string `field:"optional" json:"runtime" yaml:"runtime"`
	// List of string values to customize labels for MLS systems, such as SELinux. See https://docs.docker.com/engine/reference/run/#security-configuration.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#security_opts Container#security_opts}
	SecurityOpts *[]*string `field:"optional" json:"securityOpts" yaml:"securityOpts"`
	// Size of `/dev/shm` in MBs.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#shm_size Container#shm_size}
	ShmSize *float64 `field:"optional" json:"shmSize" yaml:"shmSize"`
	// If `true`, then the Docker container will be started after creation.
	//
	// If `false`, then the container is only created. Defaults to `true`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#start Container#start}
	Start interface{} `field:"optional" json:"start" yaml:"start"`
	// If `true`, keep STDIN open even if not attached (`docker run -i`). Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#stdin_open Container#stdin_open}
	StdinOpen interface{} `field:"optional" json:"stdinOpen" yaml:"stdinOpen"`
	// Signal to stop a container (default `SIGTERM`).
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#stop_signal Container#stop_signal}
	StopSignal *string `field:"optional" json:"stopSignal" yaml:"stopSignal"`
	// Timeout (in seconds) to stop a container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#stop_timeout Container#stop_timeout}
	StopTimeout *float64 `field:"optional" json:"stopTimeout" yaml:"stopTimeout"`
	// Key/value pairs for the storage driver options, e.g. `size`: `120G`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#storage_opts Container#storage_opts}
	StorageOpts *map[string]*string `field:"optional" json:"storageOpts" yaml:"storageOpts"`
	// A map of kernel parameters (sysctls) to set in the container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#sysctls Container#sysctls}
	Sysctls *map[string]*string `field:"optional" json:"sysctls" yaml:"sysctls"`
	// A map of container directories which should be replaced by `tmpfs mounts`, and their corresponding mount options.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#tmpfs Container#tmpfs}
	Tmpfs *map[string]*string `field:"optional" json:"tmpfs" yaml:"tmpfs"`
	// If `true`, allocate a pseudo-tty (`docker run -t`). Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#tty Container#tty}
	Tty interface{} `field:"optional" json:"tty" yaml:"tty"`
	// ulimit block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#ulimit Container#ulimit}
	Ulimit interface{} `field:"optional" json:"ulimit" yaml:"ulimit"`
	// upload block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#upload Container#upload}
	Upload interface{} `field:"optional" json:"upload" yaml:"upload"`
	// User used for run the first process.
	//
	// Format is `user` or `user:group` which user and group can be passed literraly or by name.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#user Container#user}
	User *string `field:"optional" json:"user" yaml:"user"`
	// Sets the usernamespace mode for the container when usernamespace remapping option is enabled.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#userns_mode Container#userns_mode}
	UsernsMode *string `field:"optional" json:"usernsMode" yaml:"usernsMode"`
	// volumes block.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#volumes Container#volumes}
	Volumes interface{} `field:"optional" json:"volumes" yaml:"volumes"`
	// If `true`, then the Docker container is waited for being healthy state after creation.
	//
	// If `false`, then the container health state is not checked. Defaults to `false`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#wait Container#wait}
	Wait interface{} `field:"optional" json:"wait" yaml:"wait"`
	// The timeout in seconds to wait the container to be healthy after creation. Defaults to `60`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#wait_timeout Container#wait_timeout}
	WaitTimeout *float64 `field:"optional" json:"waitTimeout" yaml:"waitTimeout"`
	// The working directory for commands to run in.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/container#working_dir Container#working_dir}
	WorkingDir *string `field:"optional" json:"workingDir" yaml:"workingDir"`
}

