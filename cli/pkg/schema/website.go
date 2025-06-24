package schema

type WebsiteIntent struct {
	BaseDir      string      `json:"baseDir" yaml:"baseDir"`
	AssetDir     string      `json:"assetDir" yaml:"assetDir"`
	ErrorPage    string      `json:"errorPage" yaml:"errorPage"`
	BuildCommand *string     `json:"build,omitempty" yaml:"build,omitempty"`
	Dev          *WebsiteDev `json:"dev,omitempty" yaml:"dev,omitempty"`
}

type WebsiteDev struct {
	Url     string `json:"url"`
	Command string `json:"command"`
}
