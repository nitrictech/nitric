package client

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
)

func GeneratePython(fs afero.Fs, appSpec schema.Application, outputDir string) error {
	return fmt.Errorf("python SDK generation not implemented")
}
