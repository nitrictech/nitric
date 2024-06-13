package service


type ServiceUpdateConfig struct {
	// Delay between task updates `(ns|us|ms|s|m|h)`. Defaults to `0s`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#delay Service#delay}
	Delay *string `field:"optional" json:"delay" yaml:"delay"`
	// Action on update failure: `pause`, `continue` or `rollback`. Defaults to `pause`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#failure_action Service#failure_action}
	FailureAction *string `field:"optional" json:"failureAction" yaml:"failureAction"`
	// Failure rate to tolerate during an update. Defaults to `0.0`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#max_failure_ratio Service#max_failure_ratio}
	MaxFailureRatio *string `field:"optional" json:"maxFailureRatio" yaml:"maxFailureRatio"`
	// Duration after each task update to monitor for failure (ns|us|ms|s|m|h). Defaults to `5s`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#monitor Service#monitor}
	Monitor *string `field:"optional" json:"monitor" yaml:"monitor"`
	// Update order: either 'stop-first' or 'start-first'. Defaults to `stop-first`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#order Service#order}
	Order *string `field:"optional" json:"order" yaml:"order"`
	// Maximum number of tasks to be updated in one iteration. Defaults to `1`.
	//
	// Docs at Terraform Registry: {@link https://registry.terraform.io/providers/kreuzwerker/docker/3.0.2/docs/resources/service#parallelism Service#parallelism}
	Parallelism *float64 `field:"optional" json:"parallelism" yaml:"parallelism"`
}

