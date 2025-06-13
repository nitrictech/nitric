package client

import (
	"bytes"
	"fmt"
	"path/filepath"
	"text/template"

	_ "embed"

	"github.com/nitrictech/nitric/cli/pkg/schema"
	"github.com/spf13/afero"
)

//go:embed ts_client_template
var tsClientTemplate string

type TSSDKTemplateData struct {
	Package string
	Buckets []ResourceNameNormalizer
}

func AppSpecToTSTemplateData(appSpec schema.Application) (TSSDKTemplateData, error) {
	buckets := []ResourceNameNormalizer{}
	for name, resource := range appSpec.ResourceIntents {
		if resource.Type != "bucket" {
			continue
		}

		normalized, err := NewResourceNameNormalizer(name)
		if err != nil {
			return TSSDKTemplateData{}, fmt.Errorf("failed to normalize resource name: %w", err)
		}

		buckets = append(buckets, normalized)
	}

	return TSSDKTemplateData{
		Package: "client",
		Buckets: buckets,
	}, nil
}

// GenerateTypeScript generates TypeScript SDK
func GenerateTypeScript(fs afero.Fs, appSpec schema.Application, outputDir string) error {
	if outputDir == "" {
		outputDir = "nitric/ts/client"
	}

	tmpl := template.Must(template.New("client").Parse(tsClientTemplate))
	data, err := AppSpecToTSTemplateData(appSpec)
	if err != nil {
		return fmt.Errorf("failed to convert nitric application spec into TypeScript SDK template data: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	err = fs.MkdirAll(outputDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	filePath := filepath.Join(outputDir, "client.ts")
	err = afero.WriteFile(fs, filePath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("failed to write generated file: %w", err)
	}

	fmt.Printf("TypeScript SDK generated at %s\n", filePath)

	return nil
}
