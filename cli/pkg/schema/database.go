package schema

type DatabaseIntent struct {
	Resource `json:",inline" yaml:",inline"`

	Access map[string][]string `json:"access,omitempty" yaml:"access,omitempty"`
	// The key that this database will export to services that depend on it.
	EnvVarKey string `json:"env_var_key" yaml:"env_var_key"`
}

func (d *DatabaseIntent) GetAccess() (map[string][]string, bool) {
	return d.Access, true
}

func (d *DatabaseIntent) GetType() string {
	return "database"
}
