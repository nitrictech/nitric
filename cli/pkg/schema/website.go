package schema

type WebsiteIntent struct {
	Resource     `json:",inline" yaml:",inline"`
	BaseDir      string      `json:"baseDir" yaml:"baseDir"`
	AssetDir     string      `json:"assetDir" yaml:"assetDir"`
	ErrorPage    string      `json:"errorPage" yaml:"errorPage"`
	BuildCommand *string     `json:"build,omitempty" yaml:"build,omitempty"`
	Dev          *WebsiteDev `json:"dev,omitempty" yaml:"dev,omitempty"`
}

func (w *WebsiteIntent) GetType() string {
	return "website"
}

type WebsiteDev struct {
	Url     string `json:"url"`
	Command string `json:"command"`
}
