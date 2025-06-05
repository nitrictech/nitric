package sdk

import (
	"bytes"
	"fmt"
	"go/format"
	"path/filepath"
	"text/template"

	_ "embed"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
)

//go:embed go_client_template
var clientTemplate string

type GoSDKTemplateData struct {
	Package    string
	ImportPath string
	Buckets    []ResourceNameNormalizer
}

func AppSpecToGoTemplateData(appSpec schema.Application, goPackageName string) (GoSDKTemplateData, error) {
	buckets := []ResourceNameNormalizer{}
	for name, resource := range appSpec.ResourceIntents {

		if resource.Type != "bucket" {
			continue
		}

		normalized, err := NewResourceNameNormalizer(name)
		if err != nil {
			return GoSDKTemplateData{}, fmt.Errorf("failed to normalize resource name: %w", err)
		}

		buckets = append(buckets, normalized)
	}

	return GoSDKTemplateData{
		Package: goPackageName,
		Buckets: buckets,
	}, nil
}

// GenerateGoSDK generates Go SDK
func GenerateGoSDK(fs afero.Fs, appSpec schema.Application, outputDir string, goPackageName string) error {
	if outputDir == "" {
		outputDir = "nitric/go/client"
	}

	if goPackageName == "" {
		goPackageName = filepath.Base(outputDir)
	}

	tmpl := template.Must(template.New("client").Parse(clientTemplate))
	data, err := AppSpecToGoTemplateData(appSpec, goPackageName)
	if err != nil {
		return fmt.Errorf("failed to convert nitric application spec into Go SDK template data: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	err = fs.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "client.go")
	err = afero.WriteFile(fs, filePath, formatted, 0644)
	if err != nil {
		return fmt.Errorf("failed to write generated file: %w", err)
	}

	fmt.Printf("Go SDK generated at %s\n", filePath)

	return nil
}
