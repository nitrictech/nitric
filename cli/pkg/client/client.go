package client

import (
	"fmt"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
)

func GeneratePythonSDK(fs afero.Fs, appSpec schema.Application, outputDir string, pythonPackageName string) error {
	return fmt.Errorf("python SDK generation not implemented")
}

func GenerateJavaScriptSDK(fs afero.Fs, appSpec schema.Application, outputDir string, javascriptPackageName string) error {
	return fmt.Errorf("javascript SDK generation not implemented")
}

func GenerateTypeScriptSDK(fs afero.Fs, appSpec schema.Application, outputDir string, typescriptPackageName string) error {
	return fmt.Errorf("typescript SDK generation not implemented")
}
