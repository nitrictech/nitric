package client

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
)

func GeneratePython(fs afero.Fs, appSpec schema.Application, outputDir string) error {
	return fmt.Errorf("python SDK generation not implemented")
}

func SpecHasClientResources(appSpec schema.Application) bool {
	for _, intent := range appSpec.GetResourceIntents() {
		// TODO: Add other adaptable resources here.
		if intent.GetType() == "bucket" {
			return true
		}
	}

	return false
}
