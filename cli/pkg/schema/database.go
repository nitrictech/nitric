package schema

type DatabaseIntent struct {
	Resource

	Access map[string][]string `json:"access,omitempty" yaml:"access,omitempty"`
}
