package plugin

type GoPlugin struct {
	Alias  string `json:"Alias"`
	Name   string `json:"Name"`
	Import string `json:"Import"`
}

type PluginDefintion struct {
	Gets    []string   `json:"Gets"`
	Pubsub  []GoPlugin `json:"Pubsub"`
	Storage []GoPlugin `json:"Storage"`
	Service GoPlugin   `json:"Service"`
}
