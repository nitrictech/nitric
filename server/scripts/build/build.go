package main

import (
	_ "embed"
	"log"
	"os"
	"text/template"
)

//go:embed main.tmpl
var mainTmpl string

type goPlugin struct {
	Alias  string `json:"Alias"`
	Name   string `json:"Name"`
	Import string `json:"Import"`
}

func main() {
	// template our main.go by injecting the plugin name and known plugin constructor
	tmpl, err := template.New("main").Parse(mainTmpl)
	if err != nil {
		log.Fatalf("error parsing template: %v", err)
	}

	// NOTE: The plugin definitions will come from an externally provided config file
	// This is hardcoded here as a demonstration
	err = tmpl.Execute(os.Stdout, map[string][]goPlugin{
		"Storage": {
			{
				Alias:  "s3",
				Name:   "default",
				Import: "github.com/nitrictech/plugins-poc/plugins/storage/s3",
			},
			// {
			// 	Alias:  "gcloud",
			// 	Name:   "default",
			// 	Import: "github.com/nitrictech/plugins-poc/plugins/storage/gcloud",
			// },
		},
		"PubSub": {
			{
				Alias:  "sns",
				Name:   "default",
				Import: "github.com/nitrictech/plugins-poc/plugins/pubsub/sns",
			},
			// {
			// 	Alias:  "gcloudpubsub",
			// 	Name:   "default",
			// 	Import: "github.com/nitrictech/plugins-poc/plugins/pubsub/gcloudpubsub",
			// },
		},
	})
	if err != nil {
		log.Fatalf("error executing template: %v", err)
	}
}
