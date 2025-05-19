package schema

type Application struct {
	Platform    string `json:"platform"`
	Name        string `json:"name"`
	Description string `json:"description"`

	Resources map[string]Resource `json:"resources,omitempty"`
}
