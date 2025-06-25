package schema

type BucketIntent struct {
	Resource

	ContentPath string `json:"content_path,omitempty" yaml:"content_path,omitempty"`

	Access map[string][]string `json:"access,omitempty" yaml:"access,omitempty"`
}

func (b *BucketIntent) GetAccess() (map[string][]string, bool) {
	return b.Access, true
}

func (b *BucketIntent) GetType() string {
	return "bucket"
}
