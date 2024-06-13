package datadockerlogs

import (
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

type DataDockerLogsConfig struct {
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
	// The name of the Docker Container.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#name DataDockerLogs#name}
	Name *string `field:"required" json:"name" yaml:"name"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#details DataDockerLogs#details}.
	Details interface{} `field:"optional" json:"details" yaml:"details"`
	// Discard headers that docker appends to each log entry.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#discard_headers DataDockerLogs#discard_headers}
	DiscardHeaders interface{} `field:"optional" json:"discardHeaders" yaml:"discardHeaders"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#follow DataDockerLogs#follow}.
	Follow interface{} `field:"optional" json:"follow" yaml:"follow"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#id DataDockerLogs#id}.
	//
	// Please be aware that the id field is automatically added to all resources in Terraform providers using a Terraform provider SDK version below 2.
	// If you experience problems setting this value it might not be settable. Please take a look at the provider documentation to ensure it should be settable.
	Id *string `field:"optional" json:"id" yaml:"id"`
	// If true populate computed value `logs_list_string`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#logs_list_string_enabled DataDockerLogs#logs_list_string_enabled}
	LogsListStringEnabled interface{} `field:"optional" json:"logsListStringEnabled" yaml:"logsListStringEnabled"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#show_stderr DataDockerLogs#show_stderr}.
	ShowStderr interface{} `field:"optional" json:"showStderr" yaml:"showStderr"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#show_stdout DataDockerLogs#show_stdout}.
	ShowStdout interface{} `field:"optional" json:"showStdout" yaml:"showStdout"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#since DataDockerLogs#since}.
	Since *string `field:"optional" json:"since" yaml:"since"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#tail DataDockerLogs#tail}.
	Tail *string `field:"optional" json:"tail" yaml:"tail"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#timestamps DataDockerLogs#timestamps}.
	Timestamps interface{} `field:"optional" json:"timestamps" yaml:"timestamps"`
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/data-sources/logs#until DataDockerLogs#until}.
	Until *string `field:"optional" json:"until" yaml:"until"`
}

