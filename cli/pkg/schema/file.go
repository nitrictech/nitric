package schema

import (
	"fmt"
	"strings"

	"github.com/nitrictech/nitric/cli/internal/version"
	"github.com/spf13/afero"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

func LoadFromFile(fs afero.Fs, path string, validate bool) (*Application, error) {
	if exists, err := afero.Exists(fs, path); err != nil {
		return nil, fmt.Errorf("%s application file could not be loaded at path: %s", version.ProductName, path)
	} else if !exists {
		return nil, fmt.Errorf("%s application file not found at path: %s", version.ProductName, path)
	}

	fileData, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, fmt.Errorf("error opening %s application file at path %s: %s", version.ProductName, path, err)
	}

	var appSpec *Application
	var results *gojsonschema.Result

	if strings.HasSuffix(path, ".json") {
		appSpec, results, err = ApplicationFromJson(string(fileData))
	} else if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		appSpec, results, err = ApplicationFromYaml(string(fileData))
	} else {
		return nil, fmt.Errorf("%s application file must be a .json or .yaml/.yml file: %s", version.ProductName, path)
	}

	if err != nil {
		return nil, fmt.Errorf("error parsing %s application file %s: %s", version.ProductName, path, err)
	}

	if !validate {
		return appSpec, nil
	}

	if results != nil && !results.Valid() {
		errs := ""
		for _, err := range results.Errors() {
			errs += fmt.Sprintf(" - %s\n", err)
		}
		return nil, fmt.Errorf("invalid %s application file %s:\n%s", version.ProductName, path, errs)
	}

	// Perform additional validation checks on the application
	if err := appSpec.IsValid(); err != nil {
		return nil, fmt.Errorf("invalid %s application file %s:\n%s", version.ProductName, path, err)
	}

	return appSpec, nil
}

func SaveToYaml(fs afero.Fs, path string, appSpec *Application) error {
	yaml, err := yaml.Marshal(appSpec)
	if err != nil {
		return fmt.Errorf("error marshalling %s application file %s: %s", version.ProductName, path, err)
	}

	return afero.WriteFile(fs, path, yaml, 0644)
}
