package schema

type WebsiteIntent struct {
	Resource     `json:",inline" yaml:",inline"`
	BaseDir      string      `json:"base_dir" yaml:"base_dir"`
	AssetDir     string      `json:"asset_dir" yaml:"asset_dir"`
	ErrorPage    string      `json:"error_page" yaml:"error_page"`
	BuildCommand *string     `json:"build,omitempty" yaml:"build,omitempty"`
	Dev          *WebsiteDev `json:"dev,omitempty" yaml:"dev,omitempty"`
}

func (w *WebsiteIntent) GetType() string {
	return "website"
}

type WebsiteDev struct {
	Url     string `json:"url" yaml:"url"`
	Command string `json:"command" yaml:"command"`
}
