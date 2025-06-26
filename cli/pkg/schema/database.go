package schema

type DatabaseIntent struct {
	Resource `json:",inline" yaml:",inline"`

	Access map[string][]string `json:"access,omitempty" yaml:"access,omitempty"`
}

func (d *DatabaseIntent) GetAccess() (map[string][]string, bool) {
	return d.Access, true
}

func (d *DatabaseIntent) GetType() string {
	return "database"
}
