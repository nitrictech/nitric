package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"text/template"

	"github.com/nitrictech/nitric/server/plugin"
)

//go:embed main.tmpl
var mainTmpl string

func main() {
	// template our main.go by injecting the plugin name and known plugin constructor
	tmpl, err := template.New("main").Parse(mainTmpl)
	if err != nil {
		log.Fatalf("error parsing template: %v", err)
	}

	pluginDefEnv := os.Getenv("PLUGIN_DEFINITION")
	if pluginDefEnv == "" {
		log.Fatalf("PLUGIN_DEFINITION is not set")
	}

	var pluginDef plugin.PluginDefintion
	if err := json.Unmarshal([]byte(pluginDefEnv), &pluginDef); err != nil {
		log.Fatalf("error unmarshaling plugin definition: %v", err)
	}

	for _, get := range pluginDef.Gets {
		// go get -u <get>
		cmd := exec.Command("go", "get", "-u", get)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("error running go get: %v", err)
		}
	}

	// NOTE: The plugin definitions will come from an externally provided config file
	// This is hardcoded here as a demonstration
	err = tmpl.Execute(os.Stdout, pluginDef)
	if err != nil {
		log.Fatalf("error executing template: %v", err)
	}
}
