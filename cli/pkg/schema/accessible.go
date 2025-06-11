package schema

type Access interface {
	GetAccess() map[string][]string
}

type Accessible struct {
	Access map[string][]string `json:"access,omitempty" yaml:"access,omitempty"`
}

func (a Accessible) GetAccess() map[string][]string {
	return a.Access
}

func IsAccessible(iface any) (Access, bool) {
	access, ok := iface.(Access)
	return access, ok
}
