package schema

type Website struct {
	Name         string      `json:"name"`
	BaseDir      string      `json:"baseDir"`
	AssetDir     string      `json:"assetDir"`
	ErrorPage    string      `json:"errorPage"`
	BuildCommand *string     `json:"build,omitempty"`
	Dev          *WebsiteDev `json:"dev,omitempty"`
}

type WebsiteDev struct {
	Url     string `json:"url"`
	Command string `json:"command"`
}
