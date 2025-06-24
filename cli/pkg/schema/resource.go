package schema

type Resource struct {
	SubType string `json:"sub-type,omitempty" yaml:"sub-type,omitempty" jsonschema:"-"`
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
