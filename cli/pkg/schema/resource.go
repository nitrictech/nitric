package schema

type Resource struct {
	SubType string `json:"subtype,omitempty" yaml:"subtype,omitempty"`
}

func (r *Resource) GetSubType() string {
	return r.SubType
}

func (r *Resource) GetAccess() (map[string][]string, bool) {
	return nil, false
}

type IResource interface {
	GetType() string
	GetSubType() string
	GetAccess() (map[string][]string, bool)
}
